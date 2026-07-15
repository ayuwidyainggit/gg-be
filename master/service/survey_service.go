package service

import (
	"context"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/repository"
	"sort"
	"time"
)

var ErrSurveyTitleConflict = errors.New("survey title already used in overlapping active period")
var ErrSurveyAreaDistributorRequired = errors.New("distributor_id is required when area_id is provided")
var ErrSurveyAreaDistributorMismatch = errors.New("distributor_id must contain exactly one value or match the number of area_id values")
var ErrSurveyInvalidDateFormat = errors.New("effective date must use format YYYY-MM-DD")
var ErrSurveySalesmanNotFound = errors.New("one or more salesman_id values are invalid for the current cust_id")
var ErrSurveyInvalidLevelTarget = errors.New("level_target must be one of: Salesman, Outlet, Distributor")

type SurveyInvalidSalesmenError struct {
	InvalidEmpID    []int    `json:"invalid_emp_id"`
	InvalidSalesman []string `json:"invalid_salesman"`
}

type invalidSalesmanDetail struct {
	empID int
	name  string
}

func (e *SurveyInvalidSalesmenError) Error() string {
	return "one or more salesman_id values are invalid for the current cust_id"
}

func (e *SurveyInvalidSalesmenError) Unwrap() error {
	return ErrSurveySalesmanNotFound
}

type SurveyService interface {
	List(filter entity.SurveyQueryFilter, custId string) ([]entity.SurveyListResponse, int, int, error)
	Detail(surveyId int, custId string) (entity.SurveyDetailResponse, error)
	Store(request entity.CreateSurveyBody) error
	Update(surveyId int, request entity.UpdateSurveyBody) error
	Deactivate(surveyId int, request entity.DeactivateSurveyBody) error
}

type surveySalesmanValidator interface {
	FindOneByEmpIdAndCustId(params entity.DetailSalesmanParams) (model.SalesmanList, error)
	FindSalesmanNameByEmpID(empID int64) (string, error)
}

func NewSurveyService(
	txManager repository.TransactionManager,
	surveyRepo repository.SurveyRepository,
	salesmanRepo surveySalesmanValidator,
	surveyTemplateRepo repository.SurveyTemplateRepository,
	questionTemplateRepo repository.QuestionTemplateRepository,
	qOptionTemplateRepo repository.QOptionTemplateRepository,
) SurveyService {
	return &surveyServiceImpl{
		txManager:            txManager,
		surveyRepo:           surveyRepo,
		salesmanRepo:         salesmanRepo,
		surveyTemplateRepo:   surveyTemplateRepo,
		questionTemplateRepo: questionTemplateRepo,
		qOptionTemplateRepo:  qOptionTemplateRepo,
	}
}

type surveyServiceImpl struct {
	txManager            repository.TransactionManager
	surveyRepo           repository.SurveyRepository
	salesmanRepo         surveySalesmanValidator
	surveyTemplateRepo   repository.SurveyTemplateRepository
	questionTemplateRepo repository.QuestionTemplateRepository
	qOptionTemplateRepo  repository.QOptionTemplateRepository
}

func resolveSurveyCustIds(custId, parentCustId string, distributorIds []int, surveyRepo repository.SurveyRepository) ([]string, error) {
	if len(distributorIds) == 0 {
		return []string{custId}, nil
	}

	scopeParentCustId := parentCustId
	if scopeParentCustId == "" {
		scopeParentCustId = custId
	}

	custIds, err := surveyRepo.FindCustIdsByDistributorIds(scopeParentCustId, distributorIds)
	if err != nil {
		return nil, err
	}

	if len(custIds) == 0 {
		return []string{custId}, nil
	}

	return custIds, nil
}

func resolveSalesmanCustIds(custIds []string, parentCustId string, salesmanIds []int, salesmanRepo surveySalesmanValidator) (map[int]string, error) {
	uniqueSalesmanIds := normalizeUniqueInts(salesmanIds)
	resolved := make(map[int]string, len(uniqueSalesmanIds))
	invalidDetails := make([]invalidSalesmanDetail, 0)
	for _, salesmanId := range uniqueSalesmanIds {
		isValid := false
		for _, custId := range custIds {
			_, err := salesmanRepo.FindOneByEmpIdAndCustId(entity.DetailSalesmanParams{
				CustId:       custId,
				ParentCustId: parentCustId,
				EmpId:        int64(salesmanId),
			})
			if err == nil {
				resolved[salesmanId] = custId
				isValid = true
				break
			}
		}

		if !isValid {
			salesmanName, lookupErr := salesmanRepo.FindSalesmanNameByEmpID(int64(salesmanId))
			if lookupErr != nil || salesmanName == "" {
				salesmanName = fmt.Sprintf("%d", salesmanId)
			}
			invalidDetails = append(invalidDetails, invalidSalesmanDetail{empID: salesmanId, name: salesmanName})
		}
	}

	if len(invalidDetails) > 0 {
		sort.SliceStable(invalidDetails, func(i, j int) bool {
			return invalidDetails[i].empID < invalidDetails[j].empID
		})
		invalidEmpIDs := make([]int, 0, len(invalidDetails))
		invalidSalesmen := make([]string, 0, len(invalidDetails))
		for _, invalidDetail := range invalidDetails {
			invalidEmpIDs = append(invalidEmpIDs, invalidDetail.empID)
			invalidSalesmen = append(invalidSalesmen, invalidDetail.name)
		}
		return nil, &SurveyInvalidSalesmenError{
			InvalidEmpID:    invalidEmpIDs,
			InvalidSalesman: invalidSalesmen,
		}
	}

	return resolved, nil
}

func normalizeBusinessUnitSelection(distributorIds []int) (bool, []int) {
	seen := make(map[int]bool, len(distributorIds))
	positiveIds := make([]int, 0, len(distributorIds))
	hasPrincipal := false
	for _, distributorId := range distributorIds {
		if distributorId == 0 {
			hasPrincipal = true
			continue
		}
		if distributorId < 0 || seen[distributorId] {
			continue
		}
		seen[distributorId] = true
		positiveIds = append(positiveIds, distributorId)
	}
	return hasPrincipal, positiveIds
}

func inferPrincipalScopeFromTargetCustId(distributorIds []int, targetCustId, custId, parentCustId string) bool {
	if len(distributorIds) > 0 || targetCustId == "" || custId == "" {
		return false
	}
	effectiveParent := parentCustId
	if effectiveParent == "" {
		effectiveParent = custId
	}
	return custId == effectiveParent
}

func normalizeUniqueInts(values []int) []int {
	result := make([]int, 0, len(values))
	seen := make(map[int]bool, len(values))
	for _, value := range values {
		if value <= 0 || seen[value] {
			continue
		}
		result = append(result, value)
		seen[value] = true
	}
	return result
}

func normalizeBusinessUnitIds(values []int) []int {
	result := make([]int, 0, len(values))
	seen := make(map[int]bool, len(values))
	for _, value := range values {
		if value < 0 || seen[value] {
			continue
		}
		result = append(result, value)
		seen[value] = true
	}
	return result
}

// normalizeTargetDistributorIds drops non-positive ids (0/negative) and
// deduplicates while preserving order. The DOCX (Enhance_Create_Survey_BE,
// field table for target_distributor_id) is explicit that 0/negatives are
// not real distributor ids and must be filtered before any DB call.
func normalizeTargetDistributorIds(values []int) []int {
	result := make([]int, 0, len(values))
	seen := make(map[int]bool, len(values))
	for _, value := range values {
		if value <= 0 || seen[value] {
			continue
		}
		result = append(result, value)
		seen[value] = true
	}
	return result
}

// isPrincipalScope returns true when the request targets the principal's own
// business unit (the sentinel distributor_id = 0 used by the existing service).
func isPrincipalScope(distributorIds []int) bool {
	for _, id := range distributorIds {
		if id == 0 {
			return true
		}
	}
	return false
}

// buildSurveyDistributors decides whether the request should produce rows in
// m_survey_distributor. The DOCX (Enhance_Create_Survey_BE, impact database
// sections) only writes target distributor mappings when level_target =
// "Distributor"; for "Salesman"/"Outlet" the payload's target_distributor_id
// is treated as no-op (ignored, no error).
func buildSurveyDistributors(levelTarget string, distributorIds []int) []int {
	if levelTarget != "Distributor" {
		return nil
	}
	return normalizeTargetDistributorIds(distributorIds)
}

func mergeUniqueCustIDs(base []string, values ...string) []string {
	seen := make(map[string]bool, len(base)+len(values))
	result := make([]string, 0, len(base)+len(values))
	for _, custID := range base {
		if custID == "" || seen[custID] {
			continue
		}
		seen[custID] = true
		result = append(result, custID)
	}
	for _, custID := range values {
		if custID == "" || seen[custID] {
			continue
		}
		seen[custID] = true
		result = append(result, custID)
	}
	return result
}

func (s *surveyServiceImpl) List(filter entity.SurveyQueryFilter, custId string) ([]entity.SurveyListResponse, int, int, error) {
	surveys, total, lastPage, err := s.surveyRepo.FindAllByCustId(filter, custId)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []entity.SurveyListResponse
	for _, sv := range surveys {
		responses = append(responses, entity.SurveyListResponse{
			SurveyId:          sv.SurveyId,
			CreatedAt:         sv.CreatedAt,
			AnswerFrequency:   sv.AnswerFrequency,
			SurveyTitle:       sv.SurveyTitle,
			ResponseType:      sv.ResponseType,
			EfectiveDateStart: sv.EfectiveDateStart,
			EfectiveDateEnd:   sv.EfectiveDateEnd,
			Status:            sv.Status,
		})
	}

	return responses, total, lastPage, nil
}

func (s *surveyServiceImpl) Detail(surveyId int, custId string) (entity.SurveyDetailResponse, error) {
	var response entity.SurveyDetailResponse

	survey, err := s.surveyRepo.FindOneById(surveyId, custId)
	if err != nil {
		return response, err
	}

	response = entity.SurveyDetailResponse{
		SurveyId:          survey.SurveyId,
		CreatedAt:         survey.CreatedAt,
		AnswerFrequency:   survey.AnswerFrequency,
		SurveyTitle:       survey.SurveyTitle,
		ResponseType:      survey.ResponseType,
		LevelTarget:       stringValue(survey.LevelTarget),
		EfectiveDateStart: survey.EfectiveDateStart,
		EfectiveDateEnd:   survey.EfectiveDateEnd,
		Status:            survey.Status,
		DistributorId:     entity.FlexibleIntArray{},
		AreaId:            entity.FlexibleIntArray{},
		DistributorCode:   stringValue(survey.DistributorCode),
		DistributorName:   stringValue(survey.DistributorName),
		BusinessUnits:     []entity.SurveyBusinessUnit{},
		TargetDistributor: []entity.SurveyDistributorResponse{},
		Outlet:            []entity.SurveyOutletResponse{},
		Salesman:          []entity.SurveySalesmanResponse{},
	}

	response.TargetSurvey = &entity.SurveyTargetResponse{
		TargetType: survey.TargetType,
		EmpId:      survey.EmpId,
		SalesName:  survey.SalesName,
		Area:       []entity.SurveyAreaResponse{},
		Outlet:     []entity.SurveyOutletResponse{},
		Salesman:   []entity.SurveySalesmanResponse{},
	}

	areas, _ := s.surveyRepo.FindAreasBySurveyId(surveyId)
	distributorIds := make([]int, 0, len(areas))
	areaIds := make([]int, 0, len(areas))
	businessUnits := make([]entity.SurveyBusinessUnit, 0, len(areas))
	areaExists := make(map[int]bool, len(areas))
	businessUnitExists := make(map[[2]int]bool, len(areas))
	principalBusinessUnitExists := false
	for _, a := range areas {
		areaName := ""
		if a.AreaName != nil {
			areaName = *a.AreaName
		}
		if a.DistributorId == 0 && a.AreaId > 0 {
			distributorIds = append(distributorIds, 0)
			if !principalBusinessUnitExists {
				businessUnits = append(businessUnits, entity.SurveyBusinessUnit{
					DistributorId:    0,
					AreaId:           0,
					TargetCustId:     a.TargetCustId,
					TargetCustName:   stringValue(a.CustName),
					BusinessUnitName: "Principal",
					Name:             "Principal",
					Type:             "principal",
				})
				principalBusinessUnitExists = true
			}
			if !areaExists[a.AreaId] {
				areaIds = append(areaIds, a.AreaId)
				areaExists[a.AreaId] = true
				response.TargetSurvey.Area = append(response.TargetSurvey.Area, entity.SurveyAreaResponse{
					AreaId:   a.AreaId,
					AreaName: areaName,
				})
			}
			continue
		}
		if a.DistributorId > 0 {
			distributorIds = append(distributorIds, a.DistributorId)
		}
		areaAlreadyExists := areaExists[a.AreaId]
		if a.AreaId > 0 && !areaAlreadyExists {
			areaIds = append(areaIds, a.AreaId)
			areaExists[a.AreaId] = true
		}
		businessUnitPair := [2]int{a.DistributorId, a.AreaId}
		if a.DistributorId > 0 && a.AreaId > 0 && !businessUnitExists[businessUnitPair] {
			businessUnits = append(businessUnits, entity.SurveyBusinessUnit{
				DistributorId:    a.DistributorId,
				AreaId:           a.AreaId,
				TargetCustId:     a.TargetCustId,
				TargetCustName:   stringValue(a.CustName),
				BusinessUnitName: stringValue(a.DistributorName),
				Name:             stringValue(a.DistributorName),
				Type:             "distributor",
			})
			businessUnitExists[businessUnitPair] = true
		}

		if a.AreaId > 0 && !areaAlreadyExists {
			response.TargetSurvey.Area = append(response.TargetSurvey.Area, entity.SurveyAreaResponse{
				AreaId:   a.AreaId,
				AreaName: areaName,
			})
		}
	}
	response.DistributorId = entity.FlexibleIntArray(normalizeBusinessUnitIds(distributorIds))
	if len(response.DistributorId) == 0 && survey.DistributorId != nil && *survey.DistributorId > 0 {
		response.DistributorId = entity.FlexibleIntArray{*survey.DistributorId}
	}
	response.AreaId = entity.FlexibleIntArray(areaIds)
	response.BusinessUnits = businessUnits

	outlets, _ := s.surveyRepo.FindOutletsBySurveyId(surveyId)
	for _, o := range outlets {
		outletResponse := entity.SurveyOutletResponse{
			SurveyOutletId: o.SurveyOutletId,
			OutletId:       o.OutletId,
			OutletCode:     stringValue(o.OutletCode),
			OutletName:     stringValue(o.OutletName),
			OtClassId:      o.OtClassId,
			OtClassName:    stringValue(o.OtClassName),
			OtGrpId:        o.OtGrpId,
			OtGrpName:      stringValue(o.OtGrpName),
			OtTypeId:       o.OtTypeId,
			OtTypeName:     stringValue(o.OtTypeName),
		}
		response.Outlet = append(response.Outlet, outletResponse)
		response.TargetSurvey.Outlet = append(response.TargetSurvey.Outlet, outletResponse)
	}

	salesmen, _ := s.surveyRepo.FindSalesmenBySurveyId(surveyId)
	for _, salesman := range salesmen {
		salesmanResponse := entity.SurveySalesmanResponse{
			MSurveySalesmanId: salesman.MSurveySalesmanId,
			SalesId:           salesman.SalesmanId,
			SalesTeamId:       salesman.SalesTeamId,
			SalesTeamName:     stringValue(salesman.SalesTeamName),
			SalesName:         stringValue(salesman.SalesName),
		}
		response.Salesman = append(response.Salesman, salesmanResponse)
		response.TargetSurvey.Salesman = append(response.TargetSurvey.Salesman, salesmanResponse)
		if response.TargetSurvey.EmpId == nil {
			empId := salesman.SalesmanId
			response.TargetSurvey.EmpId = &empId
		}
		if response.TargetSurvey.SalesName == nil && salesman.SalesName != nil {
			response.TargetSurvey.SalesName = salesman.SalesName
		}
	}

	distributorMappings, _ := s.surveyRepo.FindSurveyDistributorsBySurveyId(surveyId)
	for _, distributor := range distributorMappings {
		response.TargetDistributor = append(response.TargetDistributor, entity.SurveyDistributorResponse{
			MSurveyDistributorId: distributor.MSurveyDistributorId,
			DistributorId:        distributor.DistributorId,
			DistributorCode:      stringValue(distributor.DistributorCode),
			DistributorName:      stringValue(distributor.DistributorName),
		})
	}

	details, _ := s.surveyRepo.FindDetailsBySurveyId(surveyId)
	for _, d := range details {
		template, err := s.surveyTemplateRepo.FindOneById(d.SurveyTemplateId, custId)
		if err != nil {
			continue
		}

		templateNested := entity.SurveyTemplateNested{
			SurveyTemplateId: template.SurveyTemplateId,
			TemplateCode:     template.TemplateCode,
			TemplateTitle:    template.TemplateTitle,
		}

		questions, _ := s.questionTemplateRepo.FindAllBySurveyTemplateId(template.SurveyTemplateId)
		for _, q := range questions {
			qResponse := entity.QuestionTemplateResponse{
				QuestionTemplateId: q.QuestionTemplateId,
				SurveyTemplateId:   q.SurveyTemplateId,
				Question:           q.Question,
				InputType:          q.InputType,
				AnswerType:         q.AnswerType,
			}

			options, _ := s.qOptionTemplateRepo.FindAllByQuestionTemplateId(q.QuestionTemplateId)
			for _, opt := range options {
				qResponse.MQOptionTemplate = append(qResponse.MQOptionTemplate, entity.QOptionTemplateResponse{
					QOptionTemplateId: opt.QOptionTemplateId,
					Option:            opt.Option,
				})
			}
			if qResponse.MQOptionTemplate == nil {
				qResponse.MQOptionTemplate = []entity.QOptionTemplateResponse{}
			}

			templateNested.QuestionTemplate = append(templateNested.QuestionTemplate, qResponse)
		}
		if templateNested.QuestionTemplate == nil {
			templateNested.QuestionTemplate = []entity.QuestionTemplateResponse{}
		}

		response.Template = append(response.Template, templateNested)
	}
	if response.Template == nil {
		response.Template = []entity.SurveyTemplateNested{}
	}

	return response, nil
}

func (s *surveyServiceImpl) Store(request entity.CreateSurveyBody) error {
	dateStart, err := time.Parse("2006-01-02", request.EfectiveDateStart)
	if err != nil {
		return ErrSurveyInvalidDateFormat
	}
	dateEnd, err := time.Parse("2006-01-02", request.EfectiveDateEnd)
	if err != nil {
		return ErrSurveyInvalidDateFormat
	}

	hasOverlap, err := s.surveyRepo.ExistsActiveTitleOverlap(request.CustId, request.SurveyTitle, dateStart, dateEnd, nil)
	if err != nil {
		return err
	}
	if hasOverlap {
		return ErrSurveyTitleConflict
	}

	hasPrincipal, distributorIds := normalizeBusinessUnitSelection([]int(request.DistributorId))
	if !hasPrincipal && inferPrincipalScopeFromTargetCustId(distributorIds, request.TargetCustId, request.CustId, request.ParentCustId) {
		hasPrincipal = true
	}
	areasFromDistributor, err := s.surveyRepo.FindSurveyAreasByDistributorIds(distributorIds)
	if err != nil {
		return err
	}
	surveyAreas, err := buildSurveyAreas(request.AreaId, distributorIds, areasFromDistributor, hasPrincipal)
	if err != nil {
		return err
	}
	surveyAreas = resolveTargetCustIdForAreas(surveyAreas, hasPrincipal, request.CustId, request.TargetCustId)
	surveyDistributorIds := buildSurveyDistributors(request.LevelTarget, []int(request.TargetDistributorId))
	salesmanCustIds, err := resolveSurveyCustIds(request.CustId, request.ParentCustId, distributorIds, s.surveyRepo)
	if err != nil {
		return err
	}
	if hasPrincipal {
		salesmanCustIds = mergeUniqueCustIDs([]string{request.CustId}, salesmanCustIds...)
	}
	salesmanCustIdByEmpId, err := resolveSalesmanCustIds(salesmanCustIds, request.ParentCustId, []int(request.EmpId), s.salesmanRepo)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	var empId *int
	if len(request.EmpId) > 0 {
		id := request.EmpId[0]
		empId = &id
	}

	var levelTarget *string
	if request.LevelTarget != "" {
		levelValue := request.LevelTarget
		levelTarget = &levelValue
	}

	survey := model.Survey{
		CustId:            request.CustId,
		SurveyTitle:       request.SurveyTitle,
		AnswerFrequency:   request.AnswerFrequency,
		ResponseType:      request.ResponseType,
		TargetType:        request.TargetType,
		LevelTarget:       levelTarget,
		EmpId:             empId,
		EfectiveDateStart: &dateStart,
		EfectiveDateEnd:   &dateEnd,
		Status:            1,
		CreatedAt:         &now,
		CreatedBy:         &request.CreatedBy,
		UpdatedAt:         &now,
		UpdatedBy:         &request.CreatedBy,
	}

	ctx := context.Background()
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tx := repository.GetTxFromContext(txCtx)

		surveyId, err := s.surveyRepo.Store(tx, survey)
		if err != nil {
			return fmt.Errorf("failed to create survey: %w", err)
		}

		if len(surveyAreas) > 0 {
			for index := range surveyAreas {
				surveyAreas[index].SurveyId = surveyId
			}
			if err := s.surveyRepo.StoreAreas(tx, surveyAreas); err != nil {
				return fmt.Errorf("failed to create survey areas: %w", err)
			}
		}

		if len(surveyDistributorIds) > 0 {
			if err := s.surveyRepo.StoreSurveyDistributors(tx, surveyDistributorIds, surveyId, request.CustId, request.CreatedBy); err != nil {
				return fmt.Errorf("failed to create survey distributors: %w", err)
			}
		}

		if len(request.EmpId) > 0 {
			surveySalesmen := buildSurveySalesmen(salesmanCustIds, salesmanCustIdByEmpId, surveyId, []int(request.EmpId))
			if err := s.surveyRepo.StoreSalesmen(tx, surveySalesmen); err != nil {
				return fmt.Errorf("failed to create survey salesmen: %w", err)
			}
		}

		if len(request.OutletId) > 0 {
			if err := s.surveyRepo.StoreOutlets(tx, surveyId, request.OutletId); err != nil {
				return fmt.Errorf("failed to create survey outlets: %w", err)
			}
		}

		if len(request.SurveyTemplateId) > 0 {
			if err := s.surveyRepo.StoreDetails(tx, surveyId, []int(request.SurveyTemplateId)); err != nil {
				return fmt.Errorf("failed to create survey details: %w", err)
			}
		}

		return nil
	})
}

func (s *surveyServiceImpl) Update(surveyId int, request entity.UpdateSurveyBody) error {
	_, err := s.surveyRepo.FindOneById(surveyId, request.CustId)
	if err != nil {
		return errors.New("survey not found")
	}

	dateStart, err := time.Parse("2006-01-02", request.EfectiveDateStart)
	if err != nil {
		return ErrSurveyInvalidDateFormat
	}
	dateEnd, err := time.Parse("2006-01-02", request.EfectiveDateEnd)
	if err != nil {
		return ErrSurveyInvalidDateFormat
	}

	hasOverlap, err := s.surveyRepo.ExistsActiveTitleOverlap(request.CustId, request.SurveyTitle, dateStart, dateEnd, &surveyId)
	if err != nil {
		return err
	}
	if hasOverlap {
		return ErrSurveyTitleConflict
	}

	hasPrincipal, distributorIds := normalizeBusinessUnitSelection([]int(request.DistributorId))
	if !hasPrincipal && inferPrincipalScopeFromTargetCustId(distributorIds, request.TargetCustId, request.CustId, request.ParentCustId) {
		hasPrincipal = true
	}
	areasFromDistributor, err := s.surveyRepo.FindSurveyAreasByDistributorIds(distributorIds)
	if err != nil {
		return err
	}
	surveyAreas, err := buildSurveyAreas(request.AreaId, distributorIds, areasFromDistributor, hasPrincipal)
	if err != nil {
		return err
	}
	surveyAreas = resolveTargetCustIdForAreas(surveyAreas, hasPrincipal, request.CustId, request.TargetCustId)
	surveyDistributorIds := buildSurveyDistributors(request.LevelTarget, []int(request.TargetDistributorId))
	salesmanCustIds, err := resolveSurveyCustIds(request.CustId, request.ParentCustId, distributorIds, s.surveyRepo)
	if err != nil {
		return err
	}
	if hasPrincipal {
		salesmanCustIds = mergeUniqueCustIDs([]string{request.CustId}, salesmanCustIds...)
	}
	salesmanCustIdByEmpId, err := resolveSalesmanCustIds(salesmanCustIds, request.ParentCustId, []int(request.EmpId), s.salesmanRepo)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	var empId *int
	if len(request.EmpId) > 0 {
		id := request.EmpId[0]
		empId = &id
	}

	var levelTarget *string
	if request.LevelTarget != "" {
		levelValue := request.LevelTarget
		levelTarget = &levelValue
	}

	survey := model.Survey{
		SurveyTitle:       request.SurveyTitle,
		AnswerFrequency:   request.AnswerFrequency,
		ResponseType:      request.ResponseType,
		TargetType:        request.TargetType,
		LevelTarget:       levelTarget,
		EmpId:             empId,
		EfectiveDateStart: &dateStart,
		EfectiveDateEnd:   &dateEnd,
		UpdatedAt:         &now,
		UpdatedBy:         &request.UpdatedBy,
	}

	ctx := context.Background()
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tx := repository.GetTxFromContext(txCtx)

		if err := s.surveyRepo.Update(tx, surveyId, request.CustId, survey); err != nil {
			return fmt.Errorf("failed to update survey: %w", err)
		}

		if err := s.surveyRepo.DeleteAreasBySurveyId(tx, surveyId); err != nil {
			return fmt.Errorf("failed to update survey areas: %w", err)
		}
		if len(surveyAreas) > 0 {
			for index := range surveyAreas {
				surveyAreas[index].SurveyId = surveyId
			}
			if err := s.surveyRepo.StoreAreas(tx, surveyAreas); err != nil {
				return fmt.Errorf("failed to update survey areas: %w", err)
			}
		}

		if err := s.surveyRepo.DeleteSurveyDistributorsBySurveyId(tx, surveyId); err != nil {
			return fmt.Errorf("failed to update survey distributors: %w", err)
		}
		if len(surveyDistributorIds) > 0 {
			if err := s.surveyRepo.StoreSurveyDistributors(tx, surveyDistributorIds, surveyId, request.CustId, request.UpdatedBy); err != nil {
				return fmt.Errorf("failed to update survey distributors: %w", err)
			}
		}

		if err := s.surveyRepo.DeleteSalesmenBySurveyId(tx, surveyId); err != nil {
			return fmt.Errorf("failed to update survey salesmen: %w", err)
		}
		if len(request.EmpId) > 0 {
			surveySalesmen := buildSurveySalesmen(salesmanCustIds, salesmanCustIdByEmpId, surveyId, []int(request.EmpId))
			if err := s.surveyRepo.StoreSalesmen(tx, surveySalesmen); err != nil {
				return fmt.Errorf("failed to update survey salesmen: %w", err)
			}
		}

		if err := s.surveyRepo.DeleteOutletsBySurveyId(tx, surveyId); err != nil {
			return fmt.Errorf("failed to update survey outlets: %w", err)
		}
		if len(request.OutletId) > 0 {
			if err := s.surveyRepo.StoreOutlets(tx, surveyId, request.OutletId); err != nil {
				return fmt.Errorf("failed to update survey outlets: %w", err)
			}
		}

		if err := s.surveyRepo.DeleteDetailsBySurveyId(tx, surveyId); err != nil {
			return fmt.Errorf("failed to update survey details: %w", err)
		}
		if len(request.SurveyTemplateId) > 0 {
			if err := s.surveyRepo.StoreDetails(tx, surveyId, []int(request.SurveyTemplateId)); err != nil {
				return fmt.Errorf("failed to update survey details: %w", err)
			}
		}

		return nil
	})
}

func (s *surveyServiceImpl) Deactivate(surveyId int, request entity.DeactivateSurveyBody) error {
	ctx := context.Background()
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tx := repository.GetTxFromContext(txCtx)
		return s.surveyRepo.Deactivate(tx, surveyId, request.CustId, request.IsActive, request.UpdatedBy)
	})
}

func buildSurveyAreas(areaIds []int, distributorIds []int, areasFromDistributor []model.SurveyArea, hasPrincipal bool) ([]model.SurveyArea, error) {
	payloadAreaIds := normalizeUniqueInts(areaIds)
	if len(payloadAreaIds) == 0 && len(areasFromDistributor) == 0 && !hasPrincipal {
		return nil, nil
	}
	if len(distributorIds) == 0 && !hasPrincipal {
		return nil, ErrSurveyAreaDistributorRequired
	}

	surveyAreas := make([]model.SurveyArea, 0, len(areasFromDistributor)+1)
	areaPairExists := make(map[[2]int]bool, len(payloadAreaIds)+len(areasFromDistributor)+1)
	if hasPrincipal {
		for _, areaId := range payloadAreaIds {
			areaPair := [2]int{0, areaId}
			if areaId <= 0 || areaPairExists[areaPair] {
				continue
			}
			surveyAreas = append(surveyAreas, model.SurveyArea{DistributorId: 0, AreaId: areaId})
			areaPairExists[areaPair] = true
		}
	}

	for _, areaFromDistributor := range areasFromDistributor {
		areaPair := [2]int{areaFromDistributor.DistributorId, areaFromDistributor.AreaId}
		if areaFromDistributor.AreaId <= 0 || areaFromDistributor.DistributorId <= 0 || areaPairExists[areaPair] {
			continue
		}
		surveyAreas = append(surveyAreas, model.SurveyArea{
			DistributorId: areaFromDistributor.DistributorId,
			AreaId:        areaFromDistributor.AreaId,
		})
		areaPairExists[areaPair] = true
	}

	return surveyAreas, nil
}

// resolveTargetCustIdForAreas stamps m_survey_area.target_cust_id per DOCX
// (Enhance_Create_Survey_BE line 1322-1600): the principal BU row uses the
// tenant's cust_id (request.CustId), and distributor BU rows leave it empty.
// The token-derived cust_id is the source of truth — payload target_cust_id
// is only allowed to override the principal row when the principal scope is
// active; otherwise the payload value is ignored to keep tenant scope safe.
func resolveTargetCustIdForAreas(areas []model.SurveyArea, principalScope bool, requestCustId, payloadTargetCustId string) []model.SurveyArea {
	if len(areas) == 0 {
		return areas
	}
	resolved := make([]model.SurveyArea, len(areas))
	copy(resolved, areas)
	for index := range resolved {
		if resolved[index].DistributorId == 0 {
			if payloadTargetCustId != "" {
				resolved[index].TargetCustId = payloadTargetCustId
			} else if principalScope {
				resolved[index].TargetCustId = requestCustId
			}
		}
	}
	return resolved
}

func buildSurveySalesmen(custIds []string, custIdBySalesmanId map[int]string, surveyId int, salesmanIds []int) []model.SurveySalesman {
	salesmanIds = normalizeUniqueInts(salesmanIds)
	surveySalesmen := make([]model.SurveySalesman, 0, len(salesmanIds))
	defaultCustId := ""
	if len(custIds) > 0 {
		defaultCustId = custIds[0]
	}

	for _, salesmanId := range salesmanIds {
		custId := defaultCustId
		if resolvedCustId := custIdBySalesmanId[salesmanId]; resolvedCustId != "" {
			custId = resolvedCustId
		}
		surveySalesmen = append(surveySalesmen, model.SurveySalesman{
			CustId:     custId,
			SurveyId:   surveyId,
			SalesmanId: salesmanId,
		})
	}

	return surveySalesmen
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (s *surveyServiceImpl) validateSalesmen(custIds []string, parentCustId string, salesmanIds []int) error {
	invalidEmpIDs := make([]int, 0)
	for _, salesmanId := range salesmanIds {
		isValid := false
		for _, custId := range custIds {
			_, err := s.salesmanRepo.FindOneByEmpIdAndCustId(entity.DetailSalesmanParams{
				CustId:       custId,
				ParentCustId: parentCustId,
				EmpId:        int64(salesmanId),
			})
			if err == nil {
				isValid = true
				break
			}
		}

		if !isValid {
			invalidEmpIDs = append(invalidEmpIDs, salesmanId)
		}
	}

	if len(invalidEmpIDs) > 0 {
		sort.Ints(invalidEmpIDs)
		invalidSalesmen := make([]string, 0, len(invalidEmpIDs))
		for _, empID := range invalidEmpIDs {
			invalidSalesmen = append(invalidSalesmen, fmt.Sprintf("%d", empID))
		}
		return &SurveyInvalidSalesmenError{
			InvalidEmpID:    invalidEmpIDs,
			InvalidSalesman: invalidSalesmen,
		}
	}

	return nil
}

func BuildInvalidSalesmanErrors(err error) (map[string]any, bool) {
	var invalidErr *SurveyInvalidSalesmenError
	if !errors.As(err, &invalidErr) {
		return nil, false
	}
	return map[string]any{
		"invalid_emp_id":    invalidErr.InvalidEmpID,
		"invalid_salesman": invalidErr.InvalidSalesman,
	}, true
}
