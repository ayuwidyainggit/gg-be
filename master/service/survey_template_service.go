package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"time"
)

type SurveyTemplateService interface {
	List(filter entity.SurveyTemplateQueryFilter, custId string) ([]entity.SurveyTemplateListResponse, int, int, error)
	Detail(surveyTemplateId int, custId string) (entity.SurveyTemplateDetailResponse, error)
	Store(request entity.CreateSurveyTemplateBody) error
	Update(surveyTemplateId int, request entity.UpdateSurveyTemplateBody) error
	Delete(custId string, surveyTemplateId int, userId int64) error
}

func NewSurveyTemplateService(
	surveyTemplateRepo repository.SurveyTemplateRepository,
	questionTemplateRepo repository.QuestionTemplateRepository,
	qOptionTemplateRepo repository.QOptionTemplateRepository,
) SurveyTemplateService {
	return &surveyTemplateServiceImpl{
		surveyTemplateRepo:   surveyTemplateRepo,
		questionTemplateRepo: questionTemplateRepo,
		qOptionTemplateRepo:  qOptionTemplateRepo,
	}
}

type surveyTemplateServiceImpl struct {
	surveyTemplateRepo   repository.SurveyTemplateRepository
	questionTemplateRepo repository.QuestionTemplateRepository
	qOptionTemplateRepo  repository.QOptionTemplateRepository
}

func (s *surveyTemplateServiceImpl) List(filter entity.SurveyTemplateQueryFilter, custId string) ([]entity.SurveyTemplateListResponse, int, int, error) {
	templates, total, lastPage, err := s.surveyTemplateRepo.FindAllByCustId(filter, custId)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []entity.SurveyTemplateListResponse
	for _, t := range templates {
		responses = append(responses, entity.SurveyTemplateListResponse{
			SurveyTemplateId: t.SurveyTemplateId,
			TemplateCode:     t.TemplateCode,
			TemplateTitle:    t.TemplateTitle,
			QuestionTotal:    t.QuestionTotal,
			UseImage:         t.UseImage,
			IsActive:         t.IsActive,
			CreatedAt:        t.CreatedAt,
		})
	}

	return responses, total, lastPage, nil
}

func (s *surveyTemplateServiceImpl) Detail(surveyTemplateId int, custId string) (entity.SurveyTemplateDetailResponse, error) {
	var response entity.SurveyTemplateDetailResponse

	// Get template
	template, err := s.surveyTemplateRepo.FindOneById(surveyTemplateId, custId)
	if err != nil {
		return response, err
	}

	response = entity.SurveyTemplateDetailResponse{
		SurveyTemplateId: template.SurveyTemplateId,
		TemplateCode:     template.TemplateCode,
		TemplateTitle:    template.TemplateTitle,
		QuestionTotal:    template.QuestionTotal,
		UseImage:         template.UseImage,
		IsActive:         template.IsActive,
		CreatedAt:        template.CreatedAt,
	}

	// Get questions
	questions, err := s.questionTemplateRepo.FindAllBySurveyTemplateId(surveyTemplateId)
	if err != nil {
		return response, err
	}

	if len(questions) == 0 {
		response.QuestionTemplate = []entity.QuestionTemplateResponse{}
		return response, nil
	}

	// Collect question IDs for batch fetching options
	questionIds := make([]int, len(questions))
	for i, q := range questions {
		questionIds[i] = q.QuestionTemplateId
	}

	// Get all options in one query
	options, err := s.qOptionTemplateRepo.FindAllByQuestionTemplateIds(questionIds)
	if err != nil {
		return response, err
	}

	// Group options by question ID
	optionsMap := make(map[int][]entity.QOptionTemplateResponse)
	for _, o := range options {
		optionsMap[o.QuestionTemplateId] = append(optionsMap[o.QuestionTemplateId], entity.QOptionTemplateResponse{
			QOptionTemplateId: o.QOptionTemplateId,
			Option:            o.Option,
		})
	}

	// Build question responses
	for _, q := range questions {
		qResponse := entity.QuestionTemplateResponse{
			QuestionTemplateId: q.QuestionTemplateId,
			SurveyTemplateId:   q.SurveyTemplateId,
			Question:           q.Question,
			InputType:          q.InputType,
			AnswerType:         q.AnswerType,
			UseImage:           q.UseImage,
			MQOptionTemplate:   optionsMap[q.QuestionTemplateId],
		}
		if qResponse.MQOptionTemplate == nil {
			qResponse.MQOptionTemplate = []entity.QOptionTemplateResponse{}
		}
		response.QuestionTemplate = append(response.QuestionTemplate, qResponse)
	}

	return response, nil
}

func (s *surveyTemplateServiceImpl) Store(request entity.CreateSurveyTemplateBody) error {
	// Generate template code
	templateCode, err := s.surveyTemplateRepo.GenerateTemplateCode(request.CustId)
	if err != nil {
		return errors.New("failed to generate template code")
	}

	// Begin transaction
	tx, err := s.surveyTemplateRepo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	// Store survey template
	template := model.SurveyTemplate{
		CustId:        request.CustId,
		TemplateCode:  templateCode,
		TemplateTitle: request.TemplateTitle,
		QuestionTotal: request.QuestionTotal,
		UseImage:      request.UseImage,
		IsActive:      request.IsActive,
		CreatedAt:     &now,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &now,
		UpdatedBy:     &request.CreatedBy,
	}

	surveyTemplateId, err := s.surveyTemplateRepo.Store(tx, template)
	if err != nil {
		return errors.New("failed to create survey template")
	}

	// Store questions and options
	for seq, q := range request.Question {
		question := model.QuestionTemplate{
			SurveyTemplateId: surveyTemplateId,
			Question:         q.Question,
			InputType:        q.InputType,
			AnswerType:       q.AnswerType,
			UseImage:         q.UseImage,
			Seq:              seq + 1,
			CreatedAt:        &now,
			CreatedBy:        &request.CreatedBy,
			UpdatedAt:        &now,
			UpdatedBy:        &request.CreatedBy,
		}

		questionId, err := s.questionTemplateRepo.Store(tx, question)
		if err != nil {
			return errors.New("failed to create question template")
		}

		// Store options
		for optSeq, opt := range q.QOption {
			if opt.Option == "" {
				continue
			}
			option := model.QOptionTemplate{
				QuestionTemplateId: questionId,
				Option:             opt.Option,
				Seq:                optSeq + 1,
				CreatedAt:          &now,
				CreatedBy:          &request.CreatedBy,
				UpdatedAt:          &now,
				UpdatedBy:          &request.CreatedBy,
			}
			_, err := s.qOptionTemplateRepo.Store(tx, option)
			if err != nil {
				return errors.New("failed to create option template")
			}
		}
	}

	return tx.Commit()
}

func (s *surveyTemplateServiceImpl) Update(surveyTemplateId int, request entity.UpdateSurveyTemplateBody) error {
	// Verify template exists
	_, err := s.surveyTemplateRepo.FindOneById(surveyTemplateId, request.CustId)
	if err != nil {
		return errors.New("survey template not found")
	}

	// Begin transaction
	tx, err := s.surveyTemplateRepo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	// Update survey template
	template := model.SurveyTemplate{
		TemplateTitle: request.TemplateTitle,
		QuestionTotal: request.QuestionTotal,
		UseImage:      request.UseImage,
		IsActive:      request.IsActive,
		UpdatedAt:     &now,
		UpdatedBy:     &request.UpdatedBy,
	}

	err = s.surveyTemplateRepo.Update(tx, surveyTemplateId, request.CustId, template)
	if err != nil {
		return errors.New("failed to update survey template")
	}

	// Get existing questions to delete their options
	existingQuestions, _ := s.questionTemplateRepo.FindAllBySurveyTemplateId(surveyTemplateId)
	existingQuestionIds := make([]int, len(existingQuestions))
	for i, q := range existingQuestions {
		existingQuestionIds[i] = q.QuestionTemplateId
	}

	// Delete existing options
	if len(existingQuestionIds) > 0 {
		err = s.qOptionTemplateRepo.DeleteByQuestionTemplateIds(tx, existingQuestionIds, request.UpdatedBy)
		if err != nil {
			return errors.New("failed to delete existing options")
		}
	}

	// Delete existing questions
	err = s.questionTemplateRepo.DeleteBySurveyTemplateId(tx, surveyTemplateId, request.UpdatedBy)
	if err != nil {
		return errors.New("failed to delete existing questions")
	}

	// Insert new questions and options
	for seq, q := range request.Question {
		question := model.QuestionTemplate{
			SurveyTemplateId: surveyTemplateId,
			Question:         q.Question,
			InputType:        q.InputType,
			AnswerType:       q.AnswerType,
			UseImage:         q.UseImage,
			Seq:              seq + 1,
			CreatedAt:        &now,
			CreatedBy:        &request.UpdatedBy,
			UpdatedAt:        &now,
			UpdatedBy:        &request.UpdatedBy,
		}

		questionId, err := s.questionTemplateRepo.Store(tx, question)
		if err != nil {
			return errors.New("failed to create question template")
		}

		// Store options
		for optSeq, opt := range q.QOption {
			if opt.Option == "" {
				continue
			}
			option := model.QOptionTemplate{
				QuestionTemplateId: questionId,
				Option:             opt.Option,
				Seq:                optSeq + 1,
				CreatedAt:          &now,
				CreatedBy:          &request.UpdatedBy,
				UpdatedAt:          &now,
				UpdatedBy:          &request.UpdatedBy,
			}
			_, err := s.qOptionTemplateRepo.Store(tx, option)
			if err != nil {
				return errors.New("failed to create option template")
			}
		}
	}

	return tx.Commit()
}

func (s *surveyTemplateServiceImpl) Delete(custId string, surveyTemplateId int, userId int64) error {
	// Begin transaction
	tx, err := s.surveyTemplateRepo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get existing questions
	existingQuestions, _ := s.questionTemplateRepo.FindAllBySurveyTemplateId(surveyTemplateId)
	existingQuestionIds := make([]int, len(existingQuestions))
	for i, q := range existingQuestions {
		existingQuestionIds[i] = q.QuestionTemplateId
	}

	// Delete options
	if len(existingQuestionIds) > 0 {
		err = s.qOptionTemplateRepo.DeleteByQuestionTemplateIds(tx, existingQuestionIds, userId)
		if err != nil {
			return errors.New("failed to delete options")
		}
	}

	// Delete questions
	err = s.questionTemplateRepo.DeleteBySurveyTemplateId(tx, surveyTemplateId, userId)
	if err != nil {
		return errors.New("failed to delete questions")
	}

	// Delete template
	err = s.surveyTemplateRepo.Delete(tx, custId, surveyTemplateId, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
