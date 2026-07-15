package service

import (
	"reflect"
	"testing"
	"time"

	"master/entity"
	"master/model"

	"github.com/jmoiron/sqlx"
)

func createSalesTargetRequestWithFlexibleStatus(status int, data []entity.SalesAllocatedItemRequest) entity.CreateSalesTargetRequest {
	var allocatedTotal int64
	remaining := int64(100)
	for _, item := range data {
		allocatedTotal += item.Allocated
	}
	remaining = 100 - allocatedTotal

	request := entity.CreateSalesTargetRequest{
		CustId:                          "C22001",
		CreatedBy:                       12,
		SalesTargetDistributorYearlyId:  1,
		SalesTargetDistributorMonthlyId: 2,
		Month:                           2,
		Year:                            2025,
		AllocatedTotal:                  allocatedTotal,
		MonthlyTarget:                   100,
		Remaining:                       &remaining,
		Data:                            data,
	}

	statusField := reflect.ValueOf(&request).Elem().FieldByName("Status")
	if statusField.IsValid() && statusField.CanSet() {
		switch statusField.Kind() {
		case reflect.Int:
			statusField.SetInt(int64(status))
		case reflect.Ptr:
			if statusField.Type().Elem().Kind() == reflect.Int {
				statusValue := status
				statusField.Set(reflect.ValueOf(&statusValue))
			}
		}
	}

	return request
}

type salesTargetRepositoryStub struct {
	findAllData       []model.SalesTargetList
	findAllTotal      int
	findAllLastPage   int
	findAllErr        error
	findOneData       model.SalesTargetList
	findOneErr        error
	storeID           int64
	storeErr          error
	storeAllocatedErr error
	updatePartialErr  error

	trxBeginCalled    bool
	trxCommitCalled   bool
	trxRollbackCalled bool

	storedSalesTarget     model.SalesTarget
	storedAllocatedData   []model.SalesAllocated
	updatedFields         map[string]interface{}
	updatedSalesTargetID  int64
	updatedCustID         string
	monthlyAllocationRows []model.SalesTargetMonthlyAllocation
}

func (s *salesTargetRepositoryStub) FindAll(_ entity.SalesTargetQueryFilter, _ string) ([]model.SalesTargetList, int, int, error) {
	if s.findAllErr != nil {
		return nil, 0, 0, s.findAllErr
	}
	return s.findAllData, s.findAllTotal, s.findAllLastPage, nil
}

func (s *salesTargetRepositoryStub) FindOneById(_ int64, _ string) (model.SalesTargetList, error) {
	if s.findOneErr != nil {
		return model.SalesTargetList{}, s.findOneErr
	}
	return s.findOneData, nil
}

func (s *salesTargetRepositoryStub) FindDetailsBySalesTargetId(_ int64, _ string) ([]model.SalesAllocatedDetail, error) {
	return []model.SalesAllocatedDetail{}, nil
}

func (s *salesTargetRepositoryStub) FindMonthlyDistributor(_ entity.SalesTargetMonthlyDistQuery) (model.SalesTargetMonthlyDist, error) {
	return model.SalesTargetMonthlyDist{}, nil
}

func (s *salesTargetRepositoryStub) Store(salesTarget model.SalesTarget) (int64, error) {
	s.storedSalesTarget = salesTarget
	if s.storeErr != nil {
		return 0, s.storeErr
	}
	if s.storeID == 0 {
		s.storeID = 1
	}
	return s.storeID, nil
}

func (s *salesTargetRepositoryStub) StoreAllocated(salesAllocated model.SalesAllocated) error {
	s.storedAllocatedData = append(s.storedAllocatedData, salesAllocated)
	if s.storeAllocatedErr != nil {
		return s.storeAllocatedErr
	}
	return nil
}

func (s *salesTargetRepositoryStub) Update(_ int64, _ model.SalesTarget) error {
	return nil
}

func (s *salesTargetRepositoryStub) UpdatePartial(salesTargetId int64, custId string, updates map[string]interface{}) error {
	s.updatedSalesTargetID = salesTargetId
	s.updatedCustID = custId
	s.updatedFields = updates
	if s.updatePartialErr != nil {
		return s.updatePartialErr
	}
	return nil
}

func (s *salesTargetRepositoryStub) DeleteAllocatedByTargetId(_ int64, _ string) error {
	return nil
}

func (s *salesTargetRepositoryStub) FindMonthlyAllocationByYearlyId(_ int, _ string) ([]model.SalesTargetMonthlyAllocation, error) {
	return s.monthlyAllocationRows, nil
}

func (s *salesTargetRepositoryStub) SyncTargetsToMonthly(_ *sqlx.Tx, _ string, _ int, _ int, _ int, _ int, _ int64) error {
	return nil
}

func (s *salesTargetRepositoryStub) TrxBegin() {
	s.trxBeginCalled = true
}

func (s *salesTargetRepositoryStub) TrxCommit() error {
	s.trxCommitCalled = true
	return nil
}

func (s *salesTargetRepositoryStub) TrxRollback() error {
	s.trxRollbackCalled = true
	return nil
}

func TestSalesTargetService_List_StatusDerivationMatrix(t *testing.T) {
	now := time.Now().UTC()
	createdBy := "creator"
	updatedBy := "editor"
	updatedAt := now.Add(-time.Hour)

	repo := &salesTargetRepositoryStub{
		findAllData: []model.SalesTargetList{
			{Year: now.Year(), Month: int(now.Month()), Status: entity.StatusDraft, CreatedBy: &createdBy, CreatedAt: &updatedAt},
			{Year: now.Year(), Month: int(now.Month()), Status: entity.StatusActive, UpdatedBy: &updatedBy, UpdatedAt: &updatedAt, CreatedBy: &createdBy, CreatedAt: &updatedAt},
			{Year: now.Year(), Month: int(now.Month()), Status: entity.StatusInactive, CreatedBy: &createdBy, CreatedAt: &updatedAt},
			{Year: now.Year() + 1, Month: int(now.Month()), Status: entity.StatusActive, CreatedBy: &createdBy, CreatedAt: &updatedAt},
		},
		findAllTotal:    4,
		findAllLastPage: 1,
	}

	svc := NewSalesTargetService(repo)
	data, total, lastPage, err := svc.List(entity.SalesTargetQueryFilter{}, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 4 || lastPage != 1 {
		t.Fatalf("expected total=4,lastPage=1 got total=%d,lastPage=%d", total, lastPage)
	}
	if len(data) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(data))
	}

	if data[0].Status != "Draft" {
		t.Fatalf("expected Draft, got %s", data[0].Status)
	}
	if data[1].Status != "Active" {
		t.Fatalf("expected Active, got %s", data[1].Status)
	}
	if data[2].Status != "Inactive" {
		t.Fatalf("expected Inactive, got %s", data[2].Status)
	}
	if data[3].Status != "Inactive" {
		t.Fatalf("expected Inactive for different year/month, got %s", data[3].Status)
	}

	if data[0].UpdatedBy != "creator" {
		t.Fatalf("expected fallback updated_by from creator, got %s", data[0].UpdatedBy)
	}
	if data[1].UpdatedBy != "editor" {
		t.Fatalf("expected updated_by editor, got %s", data[1].UpdatedBy)
	}
}

func TestSalesTargetService_Store_PersistRequestStatus(t *testing.T) {
	repo := &salesTargetRepositoryStub{storeID: 10}
	svc := NewSalesTargetService(repo)

	err := svc.Store(createSalesTargetRequestWithFlexibleStatus(entity.StatusDraft, []entity.SalesAllocatedItemRequest{
		{SalesmanId: 1, SalesTeamId: 10, Allocated: 100},
	}))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.trxBeginCalled || !repo.trxCommitCalled {
		t.Fatalf("expected transaction begin and commit to be called")
	}
	if repo.storedSalesTarget.Status != entity.StatusDraft {
		t.Fatalf("expected stored status=%d, got %d", entity.StatusDraft, repo.storedSalesTarget.Status)
	}
}

func TestSalesTargetService_Store_AllAssigneesMustBeActive(t *testing.T) {
	repo := &salesTargetRepositoryStub{storeID: 10}
	svc := NewSalesTargetService(repo)

	err := svc.Store(createSalesTargetRequestWithFlexibleStatus(entity.StatusDraft, []entity.SalesAllocatedItemRequest{
		{SalesmanId: 1, SalesTeamId: 10, Allocated: 30},
		{SalesmanId: 2, SalesTeamId: 20, Allocated: 30},
		{SalesmanId: 3, SalesTeamId: 30, Allocated: 40},
	}))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.storedAllocatedData) != 3 {
		t.Fatalf("expected 3 allocated rows, got %d", len(repo.storedAllocatedData))
	}

	for i, row := range repo.storedAllocatedData {
		if !row.IsActive {
			t.Fatalf("expected allocated row %d to be active", i)
		}
	}
}

func TestSalesTargetService_Store_RollbackWhenOneAllocatedInsertFails(t *testing.T) {
	repo := &salesTargetRepositoryStub{
		storeID:           10,
		storeAllocatedErr: assertAnError{},
	}
	svc := NewSalesTargetService(repo)

	err := svc.Store(createSalesTargetRequestWithFlexibleStatus(entity.StatusDraft, []entity.SalesAllocatedItemRequest{
		{SalesmanId: 1, SalesTeamId: 10, Allocated: 100},
	}))
	if err == nil {
		t.Fatalf("expected error when allocated insert fails")
	}

	if !repo.trxBeginCalled {
		t.Fatalf("expected transaction begin to be called")
	}
	if !repo.trxRollbackCalled {
		t.Fatalf("expected rollback to be called")
	}
	if repo.trxCommitCalled {
		t.Fatalf("expected commit not to be called")
	}
}

func TestSalesTargetService_List_StatusUsesStoredStatusForFutureMonth(t *testing.T) {
	createdBy := "creator"
	createdAt := time.Now().UTC().Add(-time.Hour)

	repo := &salesTargetRepositoryStub{
		findAllData: []model.SalesTargetList{{
			Year:      time.Now().UTC().Year() + 1,
			Month:     7,
			Status:    entity.StatusDraft,
			CreatedBy: &createdBy,
			CreatedAt: &createdAt,
		}},
		findAllTotal:    1,
		findAllLastPage: 1,
	}

	svc := NewSalesTargetService(repo)
	data, _, _, err := svc.List(entity.SalesTargetQueryFilter{}, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("expected 1 row, got %d", len(data))
	}

	if data[0].Status != "Draft" {
		t.Fatalf("expected Draft, got %s", data[0].Status)
	}
}

func TestSalesTargetService_Detail_UsesConsistentInactiveLabel(t *testing.T) {
	createdBy := "creator"
	createdAt := time.Now().UTC().Add(-time.Hour)

	repo := &salesTargetRepositoryStub{
		findOneData: model.SalesTargetList{
			SalesTargetId: 1,
			Year:          time.Now().UTC().Year(),
			Month:         int(time.Now().UTC().Month()),
			Status:        entity.StatusInactive,
			CreatedBy:     &createdBy,
			CreatedAt:     &createdAt,
		},
	}

	svc := NewSalesTargetService(repo)
	resp, err := svc.Detail(1, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Status != "Inactive" {
		t.Fatalf("expected Inactive, got %s", resp.Status)
	}
}

func TestSalesTargetService_List_UsesConsistentInactiveLabel(t *testing.T) {
	createdBy := "creator"
	createdAt := time.Now().UTC().Add(-time.Hour)

	repo := &salesTargetRepositoryStub{
		findAllData: []model.SalesTargetList{{
			SalesTargetId: 1,
			Year:          time.Now().UTC().Year(),
			Month:         int(time.Now().UTC().Month()),
			Status:        entity.StatusInactive,
			CreatedBy:     &createdBy,
			CreatedAt:     &createdAt,
		}},
		findAllTotal:    1,
		findAllLastPage: 1,
	}

	svc := NewSalesTargetService(repo)
	data, _, _, err := svc.List(entity.SalesTargetQueryFilter{}, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("expected 1 row, got %d", len(data))
	}

	if data[0].Status != "Inactive" {
		t.Fatalf("expected Inactive, got %s", data[0].Status)
	}
}

func TestSalesTargetService_Detail_StatusUsesStoredStatusForFutureMonth(t *testing.T) {
	createdBy := "creator"
	createdAt := time.Now().UTC().Add(-time.Hour)

	repo := &salesTargetRepositoryStub{
		findOneData: model.SalesTargetList{
			SalesTargetId: 1,
			Year:          time.Now().UTC().Year() + 1,
			Month:         7,
			Status:        entity.StatusDraft,
			CreatedBy:     &createdBy,
			CreatedAt:     &createdAt,
		},
	}

	svc := NewSalesTargetService(repo)
	resp, err := svc.Detail(1, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Status != "Draft" {
		t.Fatalf("expected Draft, got %s", resp.Status)
	}
}

func TestSalesTargetService_Update_NoData_PartialUpdateWithoutTransaction(t *testing.T) {
	createdBy := "creator"
	createdAt := time.Now().UTC().Add(-2 * time.Hour)
	repo := &salesTargetRepositoryStub{
		findOneData: model.SalesTargetList{CreatedBy: &createdBy, CreatedAt: &createdAt},
	}
	svc := NewSalesTargetService(repo)

	status := entity.StatusInactive
	remaining := int64(5)
	err := svc.Update(99, entity.UpdateSalesTargetRequest{
		CustId:    "C22001",
		UpdatedBy: 88,
		Status:    &status,
		Remaining: &remaining,
		Data:      nil,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.trxBeginCalled {
		t.Fatalf("expected transaction not started for partial update without data")
	}
	if repo.updatedSalesTargetID != 99 || repo.updatedCustID != "C22001" {
		t.Fatalf("expected update target id/cust id captured correctly")
	}
	if repo.updatedFields == nil {
		t.Fatalf("expected updated fields captured")
	}
	if _, ok := repo.updatedFields["updated_by"]; !ok {
		t.Fatalf("expected updated_by field in updates map")
	}
	if _, ok := repo.updatedFields["updated_at"]; !ok {
		t.Fatalf("expected updated_at field in updates map")
	}
	if got, ok := repo.updatedFields["status"]; !ok || got != status {
		t.Fatalf("expected status=%d in updates map, got %v", status, got)
	}
}

func TestSalesTargetService_Store_EmptyDataShouldFail(t *testing.T) {
	repo := &salesTargetRepositoryStub{}
	svc := NewSalesTargetService(repo)

	err := svc.Store(entity.CreateSalesTargetRequest{
		CustId:                          "C22001",
		CreatedBy:                       12,
		SalesTargetDistributorYearlyId:  1,
		SalesTargetDistributorMonthlyId: 2,
		Month:                           2,
		Year:                            2025,
		AllocatedTotal:                  0,
		MonthlyTarget:                   100,
		Remaining:                       func() *int64 { v := int64(100); return &v }(),
		Data:                            []entity.SalesAllocatedItemRequest{},
	})
	if err == nil {
		t.Fatalf("expected error when data is empty")
	}
}

func TestSalesTargetService_Store_AllowsZeroRemaining(t *testing.T) {
	repo := &salesTargetRepositoryStub{storeID: 10}
	svc := NewSalesTargetService(repo)
	remaining := int64(0)

	err := svc.Store(entity.CreateSalesTargetRequest{
		CustId:                          "C22001",
		CreatedBy:                       12,
		SalesTargetDistributorYearlyId:  1,
		SalesTargetDistributorMonthlyId: 2,
		Month:                           7,
		Year:                            2026,
		AllocatedTotal:                  100,
		MonthlyTarget:                   100,
		Remaining:                       &remaining,
		Status:                          func() *int { v := entity.StatusActive; return &v }(),
		Data: []entity.SalesAllocatedItemRequest{
			{SalesmanId: 1, SalesTeamId: 10, Allocated: 100},
		},
	})
	if err != nil {
		t.Fatalf("expected no error for zero remaining, got %v", err)
	}

	if repo.storedSalesTarget.Remaining != 0 {
		t.Fatalf("expected stored remaining 0, got %d", repo.storedSalesTarget.Remaining)
	}
}

type assertAnError struct{}

func (assertAnError) Error() string {
	return "assert error"
}
