package service

import (
	"context"
	"errors"
	"master/entity"
	"master/model"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

const surveyTitleOverlapConflictError = "survey title already used in overlapping active period"

type surveyRepositoryRedStub struct {
	findOneSurvey      model.Survey
	findOneErr         error
	overlapExists      bool
	overlapErr         error
	findCustIdsResult  []string
	findCustIdsErr     error
	findAreasResult    []model.SurveyArea
	findAreasErr       error
	findAreasInput     []int
	findCustIdsInput   []int
	storeAreasInput    []model.SurveyArea
	storeSalesmenInput []model.SurveySalesman
	storeDistributorInput []int
	storeDistributorCustId string
	storeDistributorSurveyId int
	storeDistributorCreatedBy int64
	storeDistributorErr       error
	deleteDistributorCalled  bool
	findDistributorsResult    []model.SurveyDistributor
	deleteAreasCalled  bool
	storeAreasErr      error
	storeSalesmenErr   error
	findDetailAreas    []model.SurveyArea
	findDetailSurvey   model.Survey
	findDetailErr      error
	findOutletsResult  []model.SurveyOutlet
	findSalesmenResult []model.SurveySalesman
	findAllResult      []model.Survey
}

type salesmanRepositoryStub struct {
	validEmpIds     map[int64]bool
	validEmpCustIds map[int64]map[string]bool
	empNames        map[int64]string
}

func (s *salesmanRepositoryStub) FindOneByEmpIdAndCustId(params entity.DetailSalesmanParams) (model.SalesmanList, error) {
	if len(s.validEmpCustIds) > 0 {
		if validCustIds, exists := s.validEmpCustIds[params.EmpId]; exists && validCustIds[params.CustId] {
			return model.SalesmanList{}, nil
		}
		return model.SalesmanList{}, errors.New("not found")
	}
	if s.validEmpIds == nil || !s.validEmpIds[params.EmpId] {
		return model.SalesmanList{}, errors.New("not found")
	}
	return model.SalesmanList{}, nil
}

func (s *salesmanRepositoryStub) FindOneSalesmanCanvasByEmpIdAndCustId(_ int64, _ string) (model.Salesman, error) {
	return model.Salesman{}, nil
}

func (s *salesmanRepositoryStub) FindSalesmanNameByEmpID(empID int64) (string, error) {
	if s.empNames != nil {
		if name, ok := s.empNames[empID]; ok {
			return name, nil
		}
	}
	return "", errors.New("not found")
}

func (s *surveyRepositoryRedStub) FindAllByCustId(_ entity.SurveyQueryFilter, _ string) ([]model.Survey, int, int, error) {
	if len(s.findAllResult) == 0 {
		return nil, 0, 0, nil
	}
	return append([]model.Survey(nil), s.findAllResult...), len(s.findAllResult), 1, nil
}

func (s *surveyRepositoryRedStub) FindOneById(_ int, _ string) (model.Survey, error) {
	if s.findDetailErr != nil {
		return model.Survey{}, s.findDetailErr
	}
	if s.findOneErr != nil {
		return model.Survey{}, s.findOneErr
	}
	if s.findDetailSurvey.SurveyId != 0 {
		return s.findDetailSurvey, nil
	}
	return s.findOneSurvey, nil
}

func (s *surveyRepositoryRedStub) FindCustIdsByDistributorIds(_ string, distributorIds []int) ([]string, error) {
	s.findCustIdsInput = append([]int(nil), distributorIds...)
	if s.findCustIdsErr != nil {
		return nil, s.findCustIdsErr
	}
	return append([]string(nil), s.findCustIdsResult...), nil
}

func (s *surveyRepositoryRedStub) FindSurveyAreasByDistributorIds(distributorIds []int) ([]model.SurveyArea, error) {
	s.findAreasInput = append([]int(nil), distributorIds...)
	if s.findAreasErr != nil {
		return nil, s.findAreasErr
	}
	return append([]model.SurveyArea(nil), s.findAreasResult...), nil
}

func (s *surveyRepositoryRedStub) Store(_ *sqlx.Tx, _ model.Survey) (int, error) {
	return 1, nil
}

func (s *surveyRepositoryRedStub) Update(_ *sqlx.Tx, _ int, _ string, _ model.Survey) error {
	return nil
}

func (s *surveyRepositoryRedStub) ExistsActiveTitleOverlap(_ string, _ string, _ time.Time, _ time.Time, _ *int) (bool, error) {
	if s.overlapErr != nil {
		return false, s.overlapErr
	}
	return s.overlapExists, nil
}

func (s *surveyRepositoryRedStub) Deactivate(_ *sqlx.Tx, _ int, _ string, _ bool, _ int64) error {
	return nil
}

func (s *surveyRepositoryRedStub) StoreAreas(_ *sqlx.Tx, surveyAreas []model.SurveyArea) error {
	s.storeAreasInput = append([]model.SurveyArea(nil), surveyAreas...)
	return s.storeAreasErr
}

func (s *surveyRepositoryRedStub) DeleteAreasBySurveyId(_ *sqlx.Tx, _ int) error {
	s.deleteAreasCalled = true
	return nil
}

func (s *surveyRepositoryRedStub) FindAreasBySurveyId(_ int) ([]model.SurveyArea, error) {
	return append([]model.SurveyArea(nil), s.findDetailAreas...), nil
}

func (s *surveyRepositoryRedStub) StoreSalesmen(_ *sqlx.Tx, surveySalesmen []model.SurveySalesman) error {
	s.storeSalesmenInput = append([]model.SurveySalesman(nil), surveySalesmen...)
	return s.storeSalesmenErr
}

func (s *surveyRepositoryRedStub) DeleteSalesmenBySurveyId(_ *sqlx.Tx, _ int) error {
	return nil
}

func (s *surveyRepositoryRedStub) FindSalesmenBySurveyId(_ int) ([]model.SurveySalesman, error) {
	return append([]model.SurveySalesman(nil), s.findSalesmenResult...), nil
}

func (s *surveyRepositoryRedStub) StoreSurveyDistributors(_ *sqlx.Tx, distributorIds []int, surveyId int, custId string, createdBy int64) error {
	s.storeDistributorInput = append([]int(nil), distributorIds...)
	s.storeDistributorSurveyId = surveyId
	s.storeDistributorCustId = custId
	s.storeDistributorCreatedBy = createdBy
	return s.storeDistributorErr
}

func (s *surveyRepositoryRedStub) DeleteSurveyDistributorsBySurveyId(_ *sqlx.Tx, _ int) error {
	s.deleteDistributorCalled = true
	return nil
}

func (s *surveyRepositoryRedStub) FindSurveyDistributorsBySurveyId(_ int) ([]model.SurveyDistributor, error) {
	return append([]model.SurveyDistributor(nil), s.findDistributorsResult...), nil
}

func (s *surveyRepositoryRedStub) StoreOutlets(_ *sqlx.Tx, _ int, _ []int) error {
	return nil
}

func (s *surveyRepositoryRedStub) DeleteOutletsBySurveyId(_ *sqlx.Tx, _ int) error {
	return nil
}

func (s *surveyRepositoryRedStub) FindOutletsBySurveyId(_ int) ([]model.SurveyOutlet, error) {
	return append([]model.SurveyOutlet(nil), s.findOutletsResult...), nil
}

func (s *surveyRepositoryRedStub) StoreDetails(_ *sqlx.Tx, _ int, _ []int) error {
	return nil
}

func (s *surveyRepositoryRedStub) DeleteDetailsBySurveyId(_ *sqlx.Tx, _ int) error {
	return nil
}

func (s *surveyRepositoryRedStub) FindDetailsBySurveyId(_ int) ([]model.SurveyDetail, error) {
	return nil, nil
}

func (s *surveyRepositoryRedStub) BeginTx() (*sqlx.Tx, error) {
	return nil, nil
}

type transactionManagerStub struct{}

func (t *transactionManagerStub) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(context.WithValue(ctx, "tx", (*sqlx.Tx)(nil)))
}

func TestSurveyService_Store_ShouldReturnConflictError_WhenNormalizedTitleOverlapsActivePeriod(t *testing.T) {
	repo := &surveyRepositoryRedStub{overlapExists: true}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "  SALES VISIT  ",
		EfectiveDateStart: "2026-01-10",
		EfectiveDateEnd:   "2026-01-20",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Mandatory",
		SurveyTemplateId:  entity.FlexibleIntArray{1},
		CustId:            "C1001",
		CreatedBy:         10,
	})

	if err == nil {
		t.Fatalf("expected conflict error, got nil")
	}

	if err.Error() != surveyTitleOverlapConflictError {
		t.Fatalf("expected error %q, got %q", surveyTitleOverlapConflictError, err.Error())
	}
}

func TestSurveyService_Store_ShouldReturnConflictError_WhenBoundaryPeriodOverlaps(t *testing.T) {
	repo := &surveyRepositoryRedStub{overlapExists: true}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Sales Visit",
		EfectiveDateStart: "2026-01-20",
		EfectiveDateEnd:   "2026-01-30",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Mandatory",
		SurveyTemplateId:  entity.FlexibleIntArray{1},
		CustId:            "C1001",
		CreatedBy:         10,
	})

	if err == nil {
		t.Fatalf("expected conflict error, got nil")
	}

	if !errors.Is(err, ErrSurveyTitleConflict) {
		t.Fatalf("expected ErrSurveyTitleConflict, got %v", err)
	}
}

func TestSurveyService_Store_ShouldNotConflict_WhenPeriodDoesNotOverlap(t *testing.T) {
	repo := &surveyRepositoryRedStub{overlapExists: false}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Sales Visit",
		EfectiveDateStart: "2026-02-01",
		EfectiveDateEnd:   "2026-02-10",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Mandatory",
		SurveyTemplateId:  entity.FlexibleIntArray{1},
		CustId:            "C1001",
		CreatedBy:         10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSurveyService_Store_ShouldNotConflict_WhenExistingSurveyIsInactive(t *testing.T) {
	repo := &surveyRepositoryRedStub{overlapExists: false}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Sales Visit",
		EfectiveDateStart: "2026-01-10",
		EfectiveDateEnd:   "2026-01-20",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Mandatory",
		SurveyTemplateId:  entity.FlexibleIntArray{1},
		CustId:            "C1001",
		CreatedBy:         10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSurveyService_Update_ShouldReturnConflictError_WhenChangingToDuplicateNormalizedTitleAndOverlapsExcludingSelf(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{
			SurveyId:          101,
			CustId:            "C1001",
			SurveyTitle:       "Current Survey",
			EfectiveDateStart: &now,
			EfectiveDateEnd:   &now,
		},
		overlapExists: true,
	}

	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Update(101, entity.UpdateSurveyBody{
		SurveyTitle:       " sales visit ",
		EfectiveDateStart: "2026-01-15",
		EfectiveDateEnd:   "2026-01-25",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Mandatory",
		SurveyTemplateId:  entity.FlexibleIntArray{1},
		CustId:            "C1001",
		UpdatedBy:         10,
	})

	if err == nil {
		t.Fatalf("expected conflict error, got nil")
	}

	if err.Error() != surveyTitleOverlapConflictError {
		t.Fatalf("expected error %q, got %q", surveyTitleOverlapConflictError, err.Error())
	}
}

func ptrString(value string) *string { return &value }

func TestSurveyService_Store_ShouldMapDistributorIdIntoSurveyAreasAndAllEmpIdsIntoSurveySalesmen(t *testing.T) {
	repo := &surveyRepositoryRedStub{
		findCustIdsResult: []string{"C260020001"},
		findAreasResult: []model.SurveyArea{
			{DistributorId: 102, AreaId: 88},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{421: true, 422: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Salesman",
		EfectiveDateStart: "2026-04-01",
		EfectiveDateEnd:   "2026-04-18",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{102},
		AreaId:            []int{88, 89},
		SurveyTemplateId:  entity.FlexibleIntArray{40},
		EmpId:             entity.FlexibleIntArray{421, 422},
		CustId:            "C260020001",
		CreatedBy:         10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.storeAreasInput) != 1 {
		t.Fatalf("expected 1 survey area from distributor lookup, got %d", len(repo.storeAreasInput))
	}
	for _, surveyArea := range repo.storeAreasInput {
		if surveyArea.DistributorId != 102 {
			t.Fatalf("expected distributor_id 102, got %d", surveyArea.DistributorId)
		}
		if surveyArea.SurveyId != 1 {
			t.Fatalf("expected survey_id 1, got %d", surveyArea.SurveyId)
		}
	}

	if len(repo.storeSalesmenInput) != 2 {
		t.Fatalf("expected 2 survey salesmen, got %d", len(repo.storeSalesmenInput))
	}
	if repo.storeSalesmenInput[0].CustId != "C260020001" {
		t.Fatalf("expected cust_id C260020001, got %s", repo.storeSalesmenInput[0].CustId)
	}
	if repo.storeSalesmenInput[0].SalesmanId != 421 || repo.storeSalesmenInput[1].SalesmanId != 422 {
		t.Fatalf("expected salesman ids [421 422], got [%d %d]", repo.storeSalesmenInput[0].SalesmanId, repo.storeSalesmenInput[1].SalesmanId)
	}
}

func TestSurveyService_Store_ShouldReturnError_WhenAreaProvidedWithoutDistributor(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C260020001"}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Salesman",
		EfectiveDateStart: "2026-04-01",
		EfectiveDateEnd:   "2026-04-18",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		AreaId:            []int{88},
		SurveyTemplateId:  entity.FlexibleIntArray{40},
		CustId:            "C260020001",
		CreatedBy:         10,
	})

	if !errors.Is(err, ErrSurveyAreaDistributorRequired) {
		t.Fatalf("expected ErrSurveyAreaDistributorRequired, got %v", err)
	}
}

func TestSurveyService_Store_ShouldAcceptEmptyDistributorIdWithTargetCustId_AsPrincipal(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C22001"}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{369: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Principal Empty Distributor",
		EfectiveDateStart: "2026-07-01",
		EfectiveDateEnd:   "2026-07-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		AreaId:            []int{82},
		TargetCustId:      "C22001",
		EmpId:             entity.FlexibleIntArray{369},
		SurveyTemplateId:  entity.FlexibleIntArray{53},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.findAreasInput) != 0 {
		t.Fatalf("expected no distributor area lookup, got %+v", repo.findAreasInput)
	}
	if len(repo.storeAreasInput) != 1 || repo.storeAreasInput[0].DistributorId != 0 || repo.storeAreasInput[0].AreaId != 82 || repo.storeAreasInput[0].TargetCustId != "C22001" {
		t.Fatalf("expected principal area row with target_cust_id, got %+v", repo.storeAreasInput)
	}
}

func TestSurveyService_Store_ShouldIgnoreEmptyDistributorIdWithTargetCustId_WhenCallerIsDistributorChild(t *testing.T) {
	repo := &surveyRepositoryRedStub{}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Child Distributor",
		EfectiveDateStart: "2026-07-01",
		EfectiveDateEnd:   "2026-07-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		AreaId:            []int{82},
		TargetCustId:      "C22001",
		SurveyTemplateId:  entity.FlexibleIntArray{53},
		CustId:            "C220010001",
		ParentCustId:      "C22001",
		CreatedBy:         10,
	})

	if !errors.Is(err, ErrSurveyAreaDistributorRequired) {
		t.Fatalf("expected ErrSurveyAreaDistributorRequired, got %v", err)
	}
}

func TestSurveyService_Store_ShouldStillRequireDistributorId_WhenAreaIdProvidedWithoutTargetCustId(t *testing.T) {
	repo := &surveyRepositoryRedStub{}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Principal Missing TargetCust",
		EfectiveDateStart: "2026-07-01",
		EfectiveDateEnd:   "2026-07-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{53},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         10,
	})

	if !errors.Is(err, ErrSurveyAreaDistributorRequired) {
		t.Fatalf("expected ErrSurveyAreaDistributorRequired, got %v", err)
	}
}

func TestSurveyService_Store_ShouldCreateDistinctSurveyAreasFromDistributorLookup(t *testing.T) {
	repo := &surveyRepositoryRedStub{
		findCustIdsResult: []string{"C260020001"},
		findAreasResult: []model.SurveyArea{
			{DistributorId: 102, AreaId: 88},
			{DistributorId: 104, AreaId: 90},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Salesman",
		EfectiveDateStart: "2026-04-01",
		EfectiveDateEnd:   "2026-04-18",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		AreaId:            []int{88, 89},
		DistributorId:     entity.FlexibleIntArray{102, 103, 104},
		SurveyTemplateId:  entity.FlexibleIntArray{40},
		CustId:            "C260020001",
		CreatedBy:         10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.storeAreasInput) != 2 {
		t.Fatalf("expected 2 distinct survey area mappings, got %d", len(repo.storeAreasInput))
	}

	areaPairs := make(map[[2]int]bool, len(repo.storeAreasInput))
	for _, surveyArea := range repo.storeAreasInput {
		areaPair := [2]int{surveyArea.DistributorId, surveyArea.AreaId}
		if areaPairs[areaPair] {
			t.Fatalf("expected no duplicate distributor-area mapping, got duplicate %+v", areaPair)
		}
		areaPairs[areaPair] = true
	}

	for _, expectedAreaPair := range [][2]int{{102, 88}, {104, 90}} {
		if !areaPairs[expectedAreaPair] {
			t.Fatalf("expected distributor-area mapping %+v to be stored, got %+v", expectedAreaPair, repo.storeAreasInput)
		}
	}
}

func TestSurveyService_Store_ShouldPreserveTwoDistributorsInSameArea(t *testing.T) {
	repo := &surveyRepositoryRedStub{
		findCustIdsResult: []string{"C220010001", "C220010002"},
		findAreasResult: []model.SurveyArea{
			{DistributorId: 67, AreaId: 82},
			{DistributorId: 68, AreaId: 82},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{370: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "survey S",
		EfectiveDateStart: "2026-04-28",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{67, 68},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{59},
		EmpId:             entity.FlexibleIntArray{370},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	areaPairs := map[[2]int]bool{}
	for _, surveyArea := range repo.storeAreasInput {
		areaPairs[[2]int{surveyArea.DistributorId, surveyArea.AreaId}] = true
	}

	if len(repo.storeAreasInput) != 2 || !areaPairs[[2]int{67, 82}] || !areaPairs[[2]int{68, 82}] {
		t.Fatalf("expected distinct mappings [67 82] and [68 82], got %+v", repo.storeAreasInput)
	}
}

func TestSurveyService_Store_ShouldPersistPrincipalAndResolveAreas(t *testing.T) {
	repo := &surveyRepositoryRedStub{
		findCustIdsResult: []string{"C220010001"},
		findAreasResult: []model.SurveyArea{
			{DistributorId: 67, AreaId: 82},
			{DistributorId: 68, AreaId: 82},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{370: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "survey S",
		EfectiveDateStart: "2026-04-28",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{67, 0, 68},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{59},
		EmpId:             entity.FlexibleIntArray{370},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedDistributorIds := []int{67, 68}
	if len(repo.findAreasInput) != len(expectedDistributorIds) || len(repo.findCustIdsInput) != len(expectedDistributorIds) {
		t.Fatalf("expected distributor_id 0 to be ignored, got area lookup %+v and cust lookup %+v", repo.findAreasInput, repo.findCustIdsInput)
	}
	for index, distributorId := range expectedDistributorIds {
		if repo.findAreasInput[index] != distributorId || repo.findCustIdsInput[index] != distributorId {
			t.Fatalf("expected distributor lookup ids %+v, got area lookup %+v and cust lookup %+v", expectedDistributorIds, repo.findAreasInput, repo.findCustIdsInput)
		}
	}

	if len(repo.storeAreasInput) != 3 {
		t.Fatalf("expected 3 distinct survey area mappings, got %d", len(repo.storeAreasInput))
	}

	areaPairs := map[[2]int]bool{}
	for _, surveyArea := range repo.storeAreasInput {
		areaPairs[[2]int{surveyArea.DistributorId, surveyArea.AreaId}] = true
	}

	if !areaPairs[[2]int{0, 82}] || !areaPairs[[2]int{67, 82}] || !areaPairs[[2]int{68, 82}] {
		t.Fatalf("expected resolved distributor-area pairs [0 82], [67 82] and [68 82], got %+v", repo.storeAreasInput)
	}
}

func TestSurveyService_Store_ShouldPersistPrincipalBusinessUnitAndDistributorMappings(t *testing.T) {
	repo := &surveyRepositoryRedStub{
		findCustIdsResult: []string{"C220010001"},
		findAreasResult:   []model.SurveyArea{{DistributorId: 67, AreaId: 82}},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{370: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "survey principal",
		EfectiveDateStart: "2026-04-28",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{0, 67, 67},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{59},
		EmpId:             entity.FlexibleIntArray{370, 370},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.findAreasInput) != 1 || repo.findAreasInput[0] != 67 {
		t.Fatalf("expected only positive distributor lookup [67], got %+v", repo.findAreasInput)
	}

	if len(repo.storeAreasInput) != 2 {
		t.Fatalf("expected principal and distributor mappings, got %+v", repo.storeAreasInput)
	}
	areaPairs := map[[2]int]bool{}
	for _, surveyArea := range repo.storeAreasInput {
		areaPairs[[2]int{surveyArea.DistributorId, surveyArea.AreaId}] = true
	}
	if !areaPairs[[2]int{0, 82}] || !areaPairs[[2]int{67, 82}] {
		t.Fatalf("expected mappings [0 82] and [67 82], got %+v", repo.storeAreasInput)
	}

	if len(repo.storeSalesmenInput) != 1 || repo.storeSalesmenInput[0].SalesmanId != 370 {
		t.Fatalf("expected unique salesman ids to be persisted once, got %+v", repo.storeSalesmenInput)
	}
}

func TestSurveyService_Detail_ShouldReturnPrincipalBusinessUnitAndDistributorNames(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{SurveyId: 101, CustId: "C22001", SurveyTitle: "survey S", EfectiveDateStart: &now, EfectiveDateEnd: &now},
		findDetailAreas: []model.SurveyArea{
			{DistributorId: 0, AreaId: 70, AreaName: ptrString("Area 70")},
			{DistributorId: 0, AreaId: 82, AreaName: ptrString("Area 82")},
			{DistributorId: 67, AreaId: 82, AreaName: ptrString("Area 82"), DistributorName: ptrString("Distributor 67")},
		},
		findOutletsResult:  []model.SurveyOutlet{},
		findSalesmenResult: []model.SurveySalesman{},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	resp, err := svc.Detail(101, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.DistributorId) != 2 || resp.DistributorId[0] != 0 || resp.DistributorId[1] != 67 {
		t.Fatalf("expected distributor list to include principal and positive ids, got %+v", resp.DistributorId)
	}
	if len(resp.BusinessUnits) != 2 {
		t.Fatalf("expected principal once plus distributor, got %+v", resp.BusinessUnits)
	}
	if resp.BusinessUnits[0].DistributorId != 0 || resp.BusinessUnits[0].AreaId != 0 || resp.BusinessUnits[0].Name != "Principal" || resp.BusinessUnits[0].Type != "principal" {
		t.Fatalf("expected principal business unit, got %+v", resp.BusinessUnits[0])
	}
	if resp.BusinessUnits[1].DistributorId != 67 || resp.BusinessUnits[1].BusinessUnitName != "Distributor 67" || resp.BusinessUnits[1].Type != "distributor" {
		t.Fatalf("expected distributor business unit, got %+v", resp.BusinessUnits[1])
	}
	if len(resp.TargetSurvey.Area) != 2 || resp.TargetSurvey.Area[0].AreaId != 70 || resp.TargetSurvey.Area[1].AreaId != 82 {
		t.Fatalf("expected principal area to remain in target survey area, got %+v", resp.TargetSurvey.Area)
	}
}

func TestSurveyService_Update_ShouldKeepSalesmanMappingUniqueAcrossSecondEdit(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{
		findOneSurvey:     model.Survey{SurveyId: 101, CustId: "C22001", SurveyTitle: "survey S", EfectiveDateStart: &now, EfectiveDateEnd: &now},
		findCustIdsResult: []string{"C220010001"},
		findAreasResult:   []model.SurveyArea{{DistributorId: 67, AreaId: 82}},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{370: true}}, txManager: &transactionManagerStub{}}

	request := entity.UpdateSurveyBody{
		SurveyTitle:       "survey S Updated",
		EfectiveDateStart: "2026-04-28",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{67},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{59},
		EmpId:             entity.FlexibleIntArray{370, 370},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		UpdatedBy:         1,
	}

	if err := svc.Update(101, request); err != nil {
		t.Fatalf("first update expected no error, got %v", err)
	}
	if err := svc.Update(101, request); err != nil {
		t.Fatalf("second update expected no error, got %v", err)
	}
	if len(repo.storeSalesmenInput) != 1 || repo.storeSalesmenInput[0].SalesmanId != 370 {
		t.Fatalf("expected unique salesman mapping to remain intact, got %+v", repo.storeSalesmenInput)
	}
}

func TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman(t *testing.T) {
	repo := &surveyRepositoryRedStub{}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{369: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Testing Survey Principal",
		EfectiveDateStart: "2026-03-05",
		EfectiveDateEnd:   "2026-05-05",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{0},
		AreaId:            []int{89, 86, 85, 84, 83, 82, 70},
		OutletId:          []int{},
		SurveyTemplateId:  entity.FlexibleIntArray{63},
		EmpId:             entity.FlexibleIntArray{369},
		CustId:            "C26002",
		ParentCustId:      "C26002",
		CreatedBy:         140,
	})

	if err != nil {
		t.Fatalf("expected principal-only survey to be created, got %v", err)
	}
	if len(repo.findAreasInput) != 0 {
		t.Fatalf("expected no distributor area lookup for principal-only sentinel, got %+v", repo.findAreasInput)
	}
	if len(repo.findCustIdsInput) != 0 {
		t.Fatalf("expected no distributor cust lookup for principal-only sentinel, got %+v", repo.findCustIdsInput)
	}

	expectedAreaIds := []int{89, 86, 85, 84, 83, 82, 70}
	if len(repo.storeAreasInput) != len(expectedAreaIds) {
		t.Fatalf("expected %d principal survey area rows, got %d: %+v", len(expectedAreaIds), len(repo.storeAreasInput), repo.storeAreasInput)
	}
	for index, expectedAreaId := range expectedAreaIds {
		storedArea := repo.storeAreasInput[index]
		if storedArea.DistributorId != 0 {
			t.Fatalf("expected principal-only survey area to use distributor_id sentinel 0, got %+v", storedArea)
		}
		if storedArea.AreaId != expectedAreaId {
			t.Fatalf("expected area_id %d at index %d, got %+v", expectedAreaId, index, storedArea)
		}
	}

	if len(repo.storeSalesmenInput) != 1 {
		t.Fatalf("expected 1 survey salesman row, got %d", len(repo.storeSalesmenInput))
	}
	if repo.storeSalesmenInput[0].CustId != "C26002" || repo.storeSalesmenInput[0].SalesmanId != 369 {
		t.Fatalf("expected salesman 369 scoped to principal cust C26002, got %+v", repo.storeSalesmenInput[0])
	}
}

func TestSurveyService_Update_ShouldReplaceSurveyAreaMappings(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{
		findOneSurvey:     model.Survey{SurveyId: 101, CustId: "C22001", SurveyTitle: "survey S", EfectiveDateStart: &now, EfectiveDateEnd: &now},
		findCustIdsResult: []string{"C220010001", "C220010002"},
		findAreasResult: []model.SurveyArea{
			{DistributorId: 67, AreaId: 82},
			{DistributorId: 68, AreaId: 82},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{370: true}}, txManager: &transactionManagerStub{}}

	err := svc.Update(101, entity.UpdateSurveyBody{
		SurveyTitle:       "survey S Updated",
		EfectiveDateStart: "2026-04-28",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{67, 68},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{59},
		EmpId:             entity.FlexibleIntArray{370},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		UpdatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.deleteAreasCalled {
		t.Fatalf("expected existing survey area mappings to be deleted")
	}

	areaPairs := map[[2]int]bool{}
	for _, surveyArea := range repo.storeAreasInput {
		areaPairs[[2]int{surveyArea.DistributorId, surveyArea.AreaId}] = true
	}
	if len(repo.storeAreasInput) != 2 || !areaPairs[[2]int{67, 82}] || !areaPairs[[2]int{68, 82}] {
		t.Fatalf("expected replacement mappings [67 82] and [68 82], got %+v", repo.storeAreasInput)
	}
}

func TestSurveyService_Store_ShouldRollback_WhenStoreSalesmenFails(t *testing.T) {
	repo := &surveyRepositoryRedStub{storeSalesmenErr: errors.New("insert salesman failed"), findCustIdsResult: []string{"C260020001"}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{421: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Salesman",
		EfectiveDateStart: "2026-04-01",
		EfectiveDateEnd:   "2026-04-18",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		DistributorId:     entity.FlexibleIntArray{102},
		AreaId:            []int{88},
		SurveyTemplateId:  entity.FlexibleIntArray{40},
		EmpId:             entity.FlexibleIntArray{421},
		CustId:            "C260020001",
		CreatedBy:         10,
	})

	if err == nil || err.Error() != "failed to create survey salesmen: insert salesman failed" {
		t.Fatalf("expected wrapped create survey salesmen error, got %v", err)
	}
}

func TestSurveyService_Store_ShouldStayCompatible_WhenEmpIdAndAreaAreEmpty(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C260020001"}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "General Survey",
		EfectiveDateStart: "2026-04-01",
		EfectiveDateEnd:   "2026-04-18",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		SurveyTemplateId:  entity.FlexibleIntArray{40},
		CustId:            "C260020001",
		CreatedBy:         10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.storeAreasInput) != 0 {
		t.Fatalf("expected no survey area rows, got %d", len(repo.storeAreasInput))
	}
	if len(repo.storeSalesmenInput) != 0 {
		t.Fatalf("expected no survey salesman rows, got %d", len(repo.storeSalesmenInput))
	}
}

func TestSurveyService_Store_ShouldReturnError_WhenSalesmanIsInvalid(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C260020001"}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{421: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Salesman Invalid",
		EfectiveDateStart: "2026-04-01",
		EfectiveDateEnd:   "2026-04-18",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		DistributorId:     entity.FlexibleIntArray{102},
		AreaId:            []int{88},
		SurveyTemplateId:  entity.FlexibleIntArray{40},
		EmpId:             entity.FlexibleIntArray{99999999},
		CustId:            "C260020001",
		CreatedBy:         10,
	})

	if !errors.Is(err, ErrSurveySalesmanNotFound) {
		t.Fatalf("expected ErrSurveySalesmanNotFound, got %v", err)
	}
	var invalidErr *SurveyInvalidSalesmenError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("expected SurveyInvalidSalesmenError, got %T", err)
	}
	if len(invalidErr.InvalidEmpID) != 1 || invalidErr.InvalidEmpID[0] != 99999999 {
		t.Fatalf("expected invalid_emp_id [99999999], got %+v", invalidErr.InvalidEmpID)
	}
	if len(invalidErr.InvalidSalesman) != 1 || invalidErr.InvalidSalesman[0] != "99999999" {
		t.Fatalf("expected invalid_salesman [99999999], got %+v", invalidErr.InvalidSalesman)
	}
	if len(repo.storeSalesmenInput) != 0 {
		t.Fatalf("expected no survey salesman rows, got %d", len(repo.storeSalesmenInput))
	}
}

func TestSurveyService_Store_ShouldAllowPrincipalOwnedAndChildDistributorSalesmen_WhenDistributorSelectionContainsZero(t *testing.T) {
	repo := &surveyRepositoryRedStub{
		findCustIdsResult: []string{"C220010001", "C220010002"},
		findAreasResult: []model.SurveyArea{
			{DistributorId: 67, AreaId: 82},
			{DistributorId: 68, AreaId: 82},
		},
	}
	svc := &surveyServiceImpl{
		surveyRepo: repo,
		salesmanRepo: &salesmanRepositoryStub{validEmpCustIds: map[int64]map[string]bool{
			369: {"C22001": true},
			219: {"C220010001": true},
		}},
		txManager: &transactionManagerStub{},
	}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Parent + Child",
		EfectiveDateStart: "2026-04-14",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Optional",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{0, 67, 68},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{39},
		EmpId:             entity.FlexibleIntArray{369, 219},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.storeSalesmenInput) != 2 {
		t.Fatalf("expected 2 survey salesman rows, got %+v", repo.storeSalesmenInput)
	}
	if repo.storeSalesmenInput[0].SalesmanId != 369 || repo.storeSalesmenInput[0].CustId != "C22001" {
		t.Fatalf("expected principal-owned salesman mapping for emp 369, got %+v", repo.storeSalesmenInput[0])
	}
	if repo.storeSalesmenInput[1].SalesmanId != 219 || repo.storeSalesmenInput[1].CustId != "C220010001" {
		t.Fatalf("expected child distributor salesman mapping for emp 219, got %+v", repo.storeSalesmenInput[1])
	}
}

func TestSurveyService_Store_ShouldReturnDetailedInvalidSalesmen_WhenAnyEmpIdIsOutOfScope(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C220010001"}}
	svc := &surveyServiceImpl{
		surveyRepo: repo,
		salesmanRepo: &salesmanRepositoryStub{validEmpCustIds: map[int64]map[string]bool{
			219: {"C220010001": true},
		}, empNames: map[int64]string{99999999: "Ghost Salesman"}},
		txManager: &transactionManagerStub{},
	}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey invalid detail",
		EfectiveDateStart: "2026-04-14",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Optional",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{67},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{39},
		EmpId:             entity.FlexibleIntArray{219, 99999999},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         1,
	})

	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
	var invalidErr *SurveyInvalidSalesmenError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("expected SurveyInvalidSalesmenError, got %T", err)
	}
	if len(invalidErr.InvalidEmpID) != 1 || invalidErr.InvalidEmpID[0] != 99999999 {
		t.Fatalf("expected invalid_emp_id [99999999], got %+v", invalidErr.InvalidEmpID)
	}
	if len(invalidErr.InvalidSalesman) != 1 || invalidErr.InvalidSalesman[0] != "Ghost Salesman" {
		t.Fatalf("expected invalid_salesman [Ghost Salesman], got %+v", invalidErr.InvalidSalesman)
	}
	if len(repo.storeSalesmenInput) != 0 {
		t.Fatalf("expected no survey salesman rows, got %d", len(repo.storeSalesmenInput))
	}
}

func TestSurveyService_Update_ShouldMirrorCreateValidationRules_ForPrincipalAndDistributorScope(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{
		findOneSurvey:     model.Survey{SurveyId: 101, CustId: "C22001", SurveyTitle: "survey S", EfectiveDateStart: &now, EfectiveDateEnd: &now},
		findCustIdsResult: []string{"C220010001", "C220010002"},
		findAreasResult: []model.SurveyArea{
			{DistributorId: 67, AreaId: 82},
			{DistributorId: 68, AreaId: 82},
		},
	}
	svc := &surveyServiceImpl{
		surveyRepo: repo,
		salesmanRepo: &salesmanRepositoryStub{validEmpCustIds: map[int64]map[string]bool{
			369: {"C22001": true},
			219: {"C220010001": true},
		}},
		txManager: &transactionManagerStub{},
	}

	err := svc.Update(101, entity.UpdateSurveyBody{
		SurveyTitle:       "survey S Updated",
		EfectiveDateStart: "2026-04-28",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{0, 67, 68},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{59},
		EmpId:             entity.FlexibleIntArray{369, 219},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		UpdatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.findAreasInput) != 2 || repo.findAreasInput[0] != 67 || repo.findAreasInput[1] != 68 {
		t.Fatalf("expected positive distributor lookup only [67 68], got %+v", repo.findAreasInput)
	}
	if len(repo.storeSalesmenInput) != 2 {
		t.Fatalf("expected 2 survey salesman rows, got %+v", repo.storeSalesmenInput)
	}
}

func TestSurveyService_Update_ShouldAcceptEmptyDistributorIdWithTargetCustId_AsPrincipal(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{SurveyId: 101, CustId: "C22001", SurveyTitle: "survey S", EfectiveDateStart: &now, EfectiveDateEnd: &now},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Update(101, entity.UpdateSurveyBody{
		SurveyTitle:       "Survey Principal Update",
		EfectiveDateStart: "2026-07-01",
		EfectiveDateEnd:   "2026-07-30",
		AnswerFrequency:   "One Time",
		ResponseType:      "Mandatory",
		AreaId:            []int{82},
		TargetCustId:      "C22001",
		SurveyTemplateId:  entity.FlexibleIntArray{53},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		UpdatedBy:         10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.deleteAreasCalled {
		t.Fatalf("expected existing area rows to be deleted before replace")
	}
	if len(repo.findAreasInput) != 0 {
		t.Fatalf("expected no distributor area lookup, got %+v", repo.findAreasInput)
	}
	if len(repo.storeAreasInput) != 1 || repo.storeAreasInput[0].DistributorId != 0 || repo.storeAreasInput[0].AreaId != 82 || repo.storeAreasInput[0].TargetCustId != "C22001" {
		t.Fatalf("expected principal area row with target_cust_id, got %+v", repo.storeAreasInput)
	}
}

func TestSurveyService_Store_ShouldValidateSalesmanAgainstDistributorChildCustId(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C220010001"}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{validEmpIds: map[int64]bool{360: true}}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Parent To Child",
		EfectiveDateStart: "2026-04-14",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Optional",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{67},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{39},
		EmpId:             entity.FlexibleIntArray{360},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.storeSalesmenInput) != 1 {
		t.Fatalf("expected 1 survey salesman row, got %d", len(repo.storeSalesmenInput))
	}

	if repo.storeSalesmenInput[0].CustId != "C220010001" {
		t.Fatalf("expected child cust_id C220010001, got %s", repo.storeSalesmenInput[0].CustId)
	}

	if repo.storeSalesmenInput[0].SalesmanId != 360 {
		t.Fatalf("expected salesman id 360, got %d", repo.storeSalesmenInput[0].SalesmanId)
	}
}

func TestSurveyService_Store_ShouldPersistSalesmanWithResolvedCustIdPerEmployee(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C220010001", "C220010002"}}
	svc := &surveyServiceImpl{
		surveyRepo: repo,
		salesmanRepo: &salesmanRepositoryStub{validEmpCustIds: map[int64]map[string]bool{
			219: {"C220010001": true},
			22:  {"C220010002": true},
		}},
		txManager: &transactionManagerStub{},
	}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:       "Survey Parent Multi Child",
		EfectiveDateStart: "2026-04-14",
		EfectiveDateEnd:   "2026-04-30",
		AnswerFrequency:   "Multiple",
		ResponseType:      "Optional",
		TargetType:        "Specific",
		DistributorId:     entity.FlexibleIntArray{67, 68},
		AreaId:            []int{82},
		SurveyTemplateId:  entity.FlexibleIntArray{39},
		EmpId:             entity.FlexibleIntArray{219, 22},
		CustId:            "C22001",
		ParentCustId:      "C22001",
		CreatedBy:         1,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.storeSalesmenInput) != 2 {
		t.Fatalf("expected 2 survey salesman rows, got %+v", repo.storeSalesmenInput)
	}
	if repo.storeSalesmenInput[0].SalesmanId != 219 || repo.storeSalesmenInput[0].CustId != "C220010001" {
		t.Fatalf("expected emp 219 with cust C220010001, got %+v", repo.storeSalesmenInput[0])
	}
	if repo.storeSalesmenInput[1].SalesmanId != 22 || repo.storeSalesmenInput[1].CustId != "C220010002" {
		t.Fatalf("expected emp 22 with cust C220010002, got %+v", repo.storeSalesmenInput[1])
	}
}

func TestSurveyService_Detail_ShouldReturnDistributorAndSalesmanArray(t *testing.T) {
	distributorCode := "DIST001"
	distributorName := "Distributor One"
	distributorId := 102
	areaName := "Area West"
	salesTeamId := 30
	salesTeamName := "Team Alpha"
	salesName := "Budi Sales"
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{
			SurveyId:        80,
			CustId:          "C260020001",
			SurveyTitle:     "Survey Salesman",
			AnswerFrequency: "One Time",
			ResponseType:    "Mandatory",
			DistributorId:   &distributorId,
			DistributorCode: &distributorCode,
			DistributorName: &distributorName,
		},
		findDetailAreas: []model.SurveyArea{
			{SurveyAreaId: 1, SurveyId: 80, DistributorId: 102, AreaId: 88, AreaName: &areaName},
		},
		findSalesmenResult: []model.SurveySalesman{
			{
				MSurveySalesmanId: 1,
				SurveyId:          80,
				SalesmanId:        20,
				SalesTeamId:       &salesTeamId,
				SalesTeamName:     &salesTeamName,
				SalesName:         &salesName,
			},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, err := svc.Detail(80, "C260020001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(response.DistributorId) != 1 || response.DistributorId[0] != 102 {
		t.Fatalf("expected distributor_id 102, got %v", response.DistributorId)
	}
	if len(response.AreaId) != 1 || response.AreaId[0] != 88 {
		t.Fatalf("expected area_id 88, got %v", response.AreaId)
	}
	if len(response.BusinessUnits) != 1 || response.BusinessUnits[0].DistributorId != 102 || response.BusinessUnits[0].AreaId != 88 {
		t.Fatalf("expected business unit [102 88], got %+v", response.BusinessUnits)
	}
	if response.DistributorCode != "DIST001" || response.DistributorName != "Distributor One" {
		t.Fatalf("expected distributor code/name, got %s/%s", response.DistributorCode, response.DistributorName)
	}
	if response.Outlet == nil || len(response.Outlet) != 0 {
		t.Fatalf("expected empty outlet array, got %+v", response.Outlet)
	}
	if len(response.Salesman) != 1 {
		t.Fatalf("expected 1 salesman, got %+v", response.Salesman)
	}
	if response.Salesman[0].MSurveySalesmanId != 1 || response.Salesman[0].SalesId != 20 || response.Salesman[0].SalesName != "Budi Sales" {
		t.Fatalf("unexpected salesman response: %+v", response.Salesman[0])
	}
	if response.Salesman[0].SalesTeamId == nil || *response.Salesman[0].SalesTeamId != 30 || response.Salesman[0].SalesTeamName != "Team Alpha" {
		t.Fatalf("unexpected salesman team response: %+v", response.Salesman[0])
	}
}

func TestSurveyService_Detail_ShouldReturnSalesmanName_WhenRepositoryResolvesLegacyCustIdMismatch(t *testing.T) {
	salesTeamId := 58
	salesTeamName := "GT"
	salesName := "Mariana"
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{SurveyId: 107, CustId: "C22001", SurveyTitle: "legacy survey"},
		findSalesmenResult: []model.SurveySalesman{
			{
				MSurveySalesmanId: 712,
				SurveyId:          107,
				SalesmanId:        361,
				SalesTeamId:       &salesTeamId,
				SalesTeamName:     &salesTeamName,
				SalesName:         &salesName,
			},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, err := svc.Detail(107, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(response.Salesman) != 1 || response.Salesman[0].SalesId != 361 || response.Salesman[0].SalesName != "Mariana" {
		t.Fatalf("expected resolved salesman name, got %+v", response.Salesman)
	}
	if response.Salesman[0].SalesTeamId == nil || *response.Salesman[0].SalesTeamId != 58 || response.Salesman[0].SalesTeamName != "GT" {
		t.Fatalf("expected resolved salesman team, got %+v", response.Salesman[0])
	}
}

func TestSurveyService_Detail_ShouldReturnAllDistributorIdsFromSurveyAreas(t *testing.T) {
	areaName := "Area Same"
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{SurveyId: 90, CustId: "C22001", SurveyTitle: "survey S"},
		findDetailAreas: []model.SurveyArea{
			{SurveyAreaId: 1, SurveyId: 90, DistributorId: 67, AreaId: 82, AreaName: &areaName},
			{SurveyAreaId: 2, SurveyId: 90, DistributorId: 68, AreaId: 82, AreaName: &areaName},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, err := svc.Detail(90, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(response.DistributorId) != 2 || response.DistributorId[0] != 67 || response.DistributorId[1] != 68 {
		t.Fatalf("expected distributor_id [67 68], got %v", response.DistributorId)
	}
	if len(response.AreaId) != 1 || response.AreaId[0] != 82 {
		t.Fatalf("expected deduped area_id [82], got %v", response.AreaId)
	}
	if len(response.BusinessUnits) != 2 {
		t.Fatalf("expected 2 business units, got %+v", response.BusinessUnits)
	}
	for _, expectedBusinessUnit := range []entity.SurveyBusinessUnit{{DistributorId: 67, AreaId: 82, BusinessUnitName: "Distributor 67", Name: "Distributor 67", Type: "distributor"}, {DistributorId: 68, AreaId: 82, BusinessUnitName: "Distributor 68", Name: "Distributor 68", Type: "distributor"}} {
		found := false
		for _, businessUnit := range response.BusinessUnits {
			if businessUnit.DistributorId == expectedBusinessUnit.DistributorId && businessUnit.AreaId == expectedBusinessUnit.AreaId && businessUnit.Type == expectedBusinessUnit.Type {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected business unit %+v, got %+v", expectedBusinessUnit, response.BusinessUnits)
		}
	}
}

func TestSurveyService_Detail_ShouldReturnOutletArray(t *testing.T) {
	distributorId := 102
	distributorCode := "DIST001"
	distributorName := "Distributor One"
	outletCode := "OUT001"
	outletName := "Outlet One"
	otClassId := 1
	otClassName := "Class A"
	otGrpId := 2
	otGrpName := "Group B"
	otTypeId := 3
	otTypeName := "Type C"
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{
			SurveyId:        81,
			CustId:          "C260020001",
			SurveyTitle:     "Survey Outlet",
			DistributorId:   &distributorId,
			DistributorCode: &distributorCode,
			DistributorName: &distributorName,
		},
		findOutletsResult: []model.SurveyOutlet{
			{
				SurveyOutletId: 1,
				SurveyId:       81,
				OutletId:       10,
				OutletCode:     &outletCode,
				OutletName:     &outletName,
				OtClassId:      &otClassId,
				OtClassName:    &otClassName,
				OtGrpId:        &otGrpId,
				OtGrpName:      &otGrpName,
				OtTypeId:       &otTypeId,
				OtTypeName:     &otTypeName,
			},
		},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, err := svc.Detail(81, "C260020001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(response.Outlet) != 1 {
		t.Fatalf("expected 1 outlet, got %+v", response.Outlet)
	}
	outlet := response.Outlet[0]
	if outlet.SurveyOutletId != 1 || outlet.OutletId != 10 || outlet.OutletCode != "OUT001" || outlet.OutletName != "Outlet One" {
		t.Fatalf("unexpected outlet response: %+v", outlet)
	}
	if outlet.OtClassId == nil || *outlet.OtClassId != 1 || outlet.OtClassName != "Class A" {
		t.Fatalf("unexpected outlet class response: %+v", outlet)
	}
	if outlet.OtGrpId == nil || *outlet.OtGrpId != 2 || outlet.OtGrpName != "Group B" {
		t.Fatalf("unexpected outlet group response: %+v", outlet)
	}
	if outlet.OtTypeId == nil || *outlet.OtTypeId != 3 || outlet.OtTypeName != "Type C" {
		t.Fatalf("unexpected outlet type response: %+v", outlet)
	}
	if response.Salesman == nil || len(response.Salesman) != 0 {
		t.Fatalf("expected empty salesman array, got %+v", response.Salesman)
	}
}

func TestSurveyService_Store_ShouldInsertTargetDistributorsForDistributorLevel_AndMarkTargetCustIdOnPrincipalArea(t *testing.T) {
	repo := &surveyRepositoryRedStub{}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:         "Survey Distributor Principal",
		EfectiveDateStart:   "2026-07-01",
		EfectiveDateEnd:     "2026-07-30",
		AnswerFrequency:     "One Time",
		ResponseType:        "Mandatory",
		LevelTarget:         "Distributor",
		DistributorId:       entity.FlexibleIntArray{0},
		AreaId:              []int{82},
		TargetCustId:        "C22001",
		TargetDistributorId: entity.FlexibleIntArray{120, 0, -1, 120},
		SurveyTemplateId:    entity.FlexibleIntArray{53},
		CustId:              "C22001",
		ParentCustId:        "C22001",
		CreatedBy:           1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.storeAreasInput) != 1 || repo.storeAreasInput[0].TargetCustId != "C22001" {
		t.Fatalf("expected principal target_cust_id C22001, got %+v", repo.storeAreasInput)
	}
	if len(repo.storeDistributorInput) != 1 || repo.storeDistributorInput[0] != 120 {
		t.Fatalf("expected filtered target_distributor_id [120], got %+v", repo.storeDistributorInput)
	}
	if repo.storeDistributorCustId != "C22001" || repo.storeDistributorSurveyId != 1 {
		t.Fatalf("unexpected distributor insert metadata: cust=%s survey=%d", repo.storeDistributorCustId, repo.storeDistributorSurveyId)
	}
}

func TestSurveyService_Store_ShouldIgnoreTargetDistributors_WhenLevelTargetIsOutlet(t *testing.T) {
	repo := &surveyRepositoryRedStub{findCustIdsResult: []string{"C220010001"}, findAreasResult: []model.SurveyArea{{DistributorId: 67, AreaId: 82}}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Store(entity.CreateSurveyBody{
		SurveyTitle:         "Survey Outlet Distributor",
		EfectiveDateStart:   "2026-07-01",
		EfectiveDateEnd:     "2026-07-30",
		AnswerFrequency:     "Multiple Times, One Day",
		ResponseType:        "Mandatory",
		LevelTarget:         "Outlet",
		DistributorId:       entity.FlexibleIntArray{67},
		AreaId:              []int{82},
		TargetDistributorId: entity.FlexibleIntArray{120},
		OutletId:            []int{3489},
		SurveyTemplateId:    entity.FlexibleIntArray{53},
		CustId:              "C22001",
		ParentCustId:        "C22001",
		CreatedBy:           1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.storeDistributorInput) != 0 {
		t.Fatalf("expected no distributor inserts for Outlet level_target, got %+v", repo.storeDistributorInput)
	}
}

func TestSurveyService_Update_ShouldReplaceSurveyDistributors_NoDuplicates(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{findOneSurvey: model.Survey{SurveyId: 101, CustId: "C22001", SurveyTitle: "survey S", EfectiveDateStart: &now, EfectiveDateEnd: &now}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	err := svc.Update(101, entity.UpdateSurveyBody{
		SurveyTitle:         "Survey Distributor",
		EfectiveDateStart:   "2026-07-01",
		EfectiveDateEnd:     "2026-07-30",
		AnswerFrequency:     "One Time",
		ResponseType:        "Mandatory",
		LevelTarget:         "Distributor",
		DistributorId:       entity.FlexibleIntArray{0},
		AreaId:              []int{82},
		TargetCustId:        "C22001",
		TargetDistributorId: entity.FlexibleIntArray{120, 120, -1},
		SurveyTemplateId:    entity.FlexibleIntArray{53},
		CustId:              "C22001",
		ParentCustId:        "C22001",
		UpdatedBy:           1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.deleteDistributorCalled {
		t.Fatalf("expected existing survey distributors to be deleted before replace")
	}
	if len(repo.storeDistributorInput) != 1 || repo.storeDistributorInput[0] != 120 {
		t.Fatalf("expected deduped distributor insert [120], got %+v", repo.storeDistributorInput)
	}
}

func TestSurveyService_Detail_ShouldExposeLevelTarget_TargetCustId_AndTargetDistributorList(t *testing.T) {
	now := time.Now().UTC()
	levelTarget := "Distributor"
	distributorCode := "DIST120"
	distributorName := "PT Makmur"
	targetCustName := "Principal Customer"
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{SurveyId: 107, CustId: "C22001", SurveyTitle: "legacy survey", LevelTarget: &levelTarget, CreatedAt: &now},
		findDetailAreas: []model.SurveyArea{{DistributorId: 0, AreaId: 82, TargetCustId: "C22001", CustName: &targetCustName}},
		findDistributorsResult: []model.SurveyDistributor{{MSurveyDistributorId: 9, DistributorId: 120, DistributorCode: &distributorCode, DistributorName: &distributorName}},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, err := svc.Detail(107, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.LevelTarget != "Distributor" {
		t.Fatalf("expected level_target Distributor, got %q", response.LevelTarget)
	}
	if len(response.BusinessUnits) != 1 || response.BusinessUnits[0].TargetCustId != "C22001" || response.BusinessUnits[0].TargetCustName != "Principal Customer" {
		t.Fatalf("expected target_cust fields on business_units, got %+v", response.BusinessUnits)
	}
	if len(response.TargetDistributor) != 1 || response.TargetDistributor[0].DistributorId != 120 || response.TargetDistributor[0].DistributorCode != "DIST120" || response.TargetDistributor[0].DistributorName != "PT Makmur" {
		t.Fatalf("expected target_distributor list, got %+v", response.TargetDistributor)
	}
}

func TestSurveyService_List_ShouldKeepLegacyAnswerFrequencyRaw(t *testing.T) {
	now := time.Now().UTC()
	repo := &surveyRepositoryRedStub{findAllResult: []model.Survey{{SurveyId: 1, AnswerFrequency: "Multiple", CreatedAt: &now}}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, _, _, err := svc.List(entity.SurveyQueryFilter{Page: 1, Limit: 10}, "C1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(response) != 1 || response[0].AnswerFrequency != "Multiple" {
		t.Fatalf("expected raw legacy answer_frequency Multiple, got %+v", response)
	}
}

func TestSurveyService_Detail_ShouldKeepLegacyAnswerFrequencyRaw(t *testing.T) {
	repo := &surveyRepositoryRedStub{findOneSurvey: model.Survey{SurveyId: 82, CustId: "C260020001", SurveyTitle: "General Survey", AnswerFrequency: "Multiple"}}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, err := svc.Detail(82, "C260020001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.AnswerFrequency != "Multiple" {
		t.Fatalf("expected raw legacy answer_frequency Multiple, got %q", response.AnswerFrequency)
	}
}

func TestSurveyService_Detail_ShouldReturnEmptyOutletAndSalesmanArrays(t *testing.T) {
	repo := &surveyRepositoryRedStub{
		findOneSurvey: model.Survey{SurveyId: 82, CustId: "C260020001", SurveyTitle: "General Survey"},
	}
	svc := &surveyServiceImpl{surveyRepo: repo, salesmanRepo: &salesmanRepositoryStub{}, txManager: &transactionManagerStub{}}

	response, err := svc.Detail(82, "C260020001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.Outlet == nil || len(response.Outlet) != 0 {
		t.Fatalf("expected empty outlet array, got %+v", response.Outlet)
	}
	if response.Salesman == nil || len(response.Salesman) != 0 {
		t.Fatalf("expected empty salesman array, got %+v", response.Salesman)
	}
	if response.TargetSurvey == nil || response.TargetSurvey.Outlet == nil || response.TargetSurvey.Salesman == nil {
		t.Fatalf("expected target_survey outlet and salesman empty arrays, got %+v", response.TargetSurvey)
	}
}
