package service

import (
	"context"
	"database/sql"
	"errors"
	"mobile/entity"
	"mobile/model"
	"mobile/repository"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type SurveyService interface {
	List(ctx context.Context, filter entity.SurveyQueryFilter) ([]model.SurveyResponse, int64, int64, error)
	Detail(ctx context.Context, surveyID int64) (*model.SurveyDetailResponse, error)
	SubmitSurvey(ctx context.Context, req entity.SubmitSurveyRequest) error
	GetSubmittedSurvey(ctx context.Context, surveyID, surveyAnswerID int64) (*model.SurveyDetailResponse, error)
	ListSurveyAnswer(ctx context.Context, surveyID, outletID int64, page, limit int) ([]model.SurveyAnswerListItem, int64, int64, error)
}

type surveyService struct {
	Repo                  repository.SurveyRepository
	DistributorRepository repository.DistributorRepository
	MOutletRepository     repository.MOutletRepository
}

func NewSurveyService(repo repository.SurveyRepository, distributorRepository repository.DistributorRepository, mOutletRepository repository.MOutletRepository) SurveyService {
	return &surveyService{
		Repo:                  repo,
		DistributorRepository: distributorRepository,
		MOutletRepository:     mOutletRepository,
	}
}

func (s *surveyService) List(ctx context.Context, filter entity.SurveyQueryFilter) ([]model.SurveyResponse, int64, int64, error) {
	total, err := s.Repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int64(1)
	if filter.Limit > 0 {
		lastPage = (total + filter.Limit - 1) / filter.Limit
	}

	if lastPage == 0 {
		lastPage = 1
	}

	data, err := s.Repo.List(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	return data, total, lastPage, nil
}

func (s *surveyService) ListSurveyAnswer(ctx context.Context, surveyID, outletID int64, page, limit int) ([]model.SurveyAnswerListItem, int64, int64, error) {
	offset := 0
	if page > 0 && limit > 0 {
		offset = (page - 1) * limit
	}

	total, err := s.Repo.CountSurveyAnswer(ctx, surveyID, outletID)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int64(1)
	if limit > 0 {
		lastPage = (total + int64(limit) - 1) / int64(limit)
	}

	if lastPage == 0 {
		lastPage = 1
	}

	data, err := s.Repo.ListSurveyAnswer(ctx, surveyID, outletID, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	return data, total, lastPage, nil
}

func (s *surveyService) Detail(ctx context.Context, surveyID int64) (*model.SurveyDetailResponse, error) {
	data, err := s.Repo.Detail(ctx, surveyID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *surveyService) SubmitSurvey(ctx context.Context, req entity.SubmitSurveyRequest) error {
	// 1. Fetch survey details to validate its existence & extract metadata
	survey, err := s.Repo.Detail(ctx, req.SurveyID)
	if err != nil {
		return err
	}

	//
	// exists, err := s.Repo.CheckExistingSubmission(ctx, req.SurveyID, req.EmpID, req.OutletID)
	// if err != nil {
	// 	return err
	// }
	// if exists {
	// 	return errors.New("survey already submitted")
	// }

	var (
		now             = time.Now()
		questionsSurvey = make(map[int64]map[int64]string) // question_template_id:[q_option_template_id]:option
	)

	// validate survey with user
	_, err = s.Repo.List(ctx, entity.SurveyQueryFilter{
		OutletID: req.OutletID,
		EmpID:    req.EmpID,
	})
	if err != nil {
		return errors.New("outlet customer (" + req.CustID + ") with id: " + strconv.FormatInt(req.OutletID, 10) + " not found")
	}

	// prepare before validate
	for _, sq := range survey.TemplateData.Questions {
		options := make(map[int64]string)
		for _, sqo := range sq.Options {
			options[sqo.QOptionTemplateID] = sqo.Option
		}

		questionsSurvey[sq.QuestionTemplateID] = options
	}

	// validate params request
	if req.SurveyTemplateID != survey.TemplateData.SurveyTemplateID {
		return errors.New("survey_template_id: " + strconv.FormatInt(req.SurveyTemplateID, 10) + " not found")
	}

	for _, question := range req.Questions {
		opt, exists := questionsSurvey[question.QuestionTemplateID]
		if !exists {
			return errors.New("question_template_id: " + strconv.FormatInt(question.QuestionTemplateID, 10) + " not found")
		}

		for _, option := range question.Options {
			optionText, exists := opt[option.QOptionTemplateID]
			if !exists {
				return errors.New("q_option_template_id: " + strconv.FormatInt(option.QOptionTemplateID, 10) + " not found")
			}

			if option.Option != optionText {
				return errors.New("option: " + option.Option + " not found")
			}
		}
	}

	answer := &model.SurveyAnswer{
		CustID:           req.CustID,
		SurveyTemplateID: req.SurveyTemplateID,
		SurveyID:         req.SurveyID,
		EmpID:            req.EmpID,
		OutletID:         req.OutletID,
		AnswerDate:       &now,
		Status:           "Submitted",
		CreatedBy:        req.UserID,
		CreatedAt:        &now,
		UpdatedBy:        req.UserID,
		UpdatedAt:        &now,
		IsDel:            false,
		Details:          make([]model.SurveyAnswerDetail, 0),
	}
	if req.DistributorID > 0 {
		distributor, err := s.DistributorRepository.GetDistributorByID(ctx, req.DistributorID)
		if err != nil {
			return err
		}

		if distributor != nil && distributor.AreaID != nil {
			answer.AreaID = distributor.AreaID
		}
	}

	for _, q := range req.Questions {
		detail := model.SurveyAnswerDetail{
			CustID:             req.CustID,
			QuestionTemplateID: q.QuestionTemplateID,
			InputType:          q.InputType,
			AnswerType:         q.AnswerType,
			Seq:                q.Seq,
			IsAnswered:         q.IsAnswered,
			FreeTextAnswer:     q.FreeTextAnswer,
			CreatedBy:          req.UserID,
			CreatedAt:          &now,
			UpdatedBy:          req.UserID,
			UpdatedAt:          &now,
			IsDel:              false,
			Files:              make([]model.SurveyAnswerFile, 0),
			Options:            make([]model.SurveyAnswerOption, 0),
		}

		for _, f := range q.Files {
			detail.Files = append(detail.Files, model.SurveyAnswerFile{
				CustID:   req.CustID,
				FileName: f.ClientFileName,
			})
		}

		for _, o := range q.Options {
			detail.Options = append(detail.Options, model.SurveyAnswerOption{
				CustID:            req.CustID,
				QOptionTemplateID: o.QOptionTemplateID,
				OptionLabel:       o.Option,
				CreatedBy:         req.UserID,
				CreatedAt:         &now,
				UpdatedBy:         req.UserID,
				UpdatedAt:         &now,
				IsDel:             false,
			})
		}

		answer.Details = append(answer.Details, detail)
	}

	return s.Repo.SubmitSurvey(ctx, answer)
}

func (s *surveyService) GetSubmittedSurvey(ctx context.Context, surveyID, surveyAnswerID int64) (*model.SurveyDetailResponse, error) {
	sumitedSurvey, err := s.Repo.GetSubmittedSurvey(ctx, surveyID, surveyAnswerID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("survey not found")
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("survey not found")
	}

	if err != nil {
		return nil, err
	}
	return sumitedSurvey, nil
}
