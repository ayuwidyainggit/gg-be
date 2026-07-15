package service

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"master/entity"
	"master/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type salesTargetDistributorRepositoryStub struct {
	list                                  []model.SalesTargetDistributorYearly
	total                                 int
	lastPage                              int
	listErr                               error
	lastListFilter                        entity.SalesTargetDistributorQueryFilter
	lastListCustID                        string
	yearly                                model.SalesTargetDistributorYearly
	yearlyErr                             error
	monthly                               []model.SalesTargetDistributorMonthly
	monthlyErr                            error
	distributorChildCustIDByDistributorID map[int]string
	lastResolvedDistributorID             int
	resolveDistributorChildCustIDErr      error
	allocationSummary                     model.SalesTargetAllocationSummary
	allocationErr                         error
	lastUpdateID                          int
	lastUpdateRequest                     entity.UpdateSalesTargetDistributorRequest
	updateYearlyErr                       error
	deleteMonthlyErr                      error
	storeMonthlyErr                       error
	storeMonthlyID                        int
	updateMonthlyErr                      error
	tx                                    *sqlx.Tx
	beginTxErr                            error
}

func (s *salesTargetDistributorRepositoryStub) FindAllByCustId(filter entity.SalesTargetDistributorQueryFilter, custID string) ([]model.SalesTargetDistributorYearly, int, int, error) {
	s.lastListFilter = filter
	s.lastListCustID = custID
	if s.listErr != nil {
		return nil, 0, 0, s.listErr
	}
	return s.list, s.total, s.lastPage, nil
}

func (s *salesTargetDistributorRepositoryStub) FindOneByIdAndCustId(_ int, _ string) (model.SalesTargetDistributorYearly, error) {
	if s.yearlyErr != nil {
		return model.SalesTargetDistributorYearly{}, s.yearlyErr
	}
	return s.yearly, nil
}

func (s *salesTargetDistributorRepositoryStub) FindMonthlyDetailsByYearlyId(_ int) ([]model.SalesTargetDistributorMonthly, error) {
	if s.monthlyErr != nil {
		return nil, s.monthlyErr
	}
	return s.monthly, nil
}

func (s *salesTargetDistributorRepositoryStub) FindChildCustIDByDistributorID(distributorID int) (string, error) {
	s.lastResolvedDistributorID = distributorID
	if s.resolveDistributorChildCustIDErr != nil {
		return "", s.resolveDistributorChildCustIDErr
	}
	if s.distributorChildCustIDByDistributorID == nil {
		return "", nil
	}
	return s.distributorChildCustIDByDistributorID[distributorID], nil
}

func (s *salesTargetDistributorRepositoryStub) FindAllocationSummaryByYearlyId(_ int, _ string) (model.SalesTargetAllocationSummary, error) {
	if s.allocationErr != nil {
		return model.SalesTargetAllocationSummary{}, s.allocationErr
	}
	return s.allocationSummary, nil
}

func (s *salesTargetDistributorRepositoryStub) StoreYearly(_ *sqlx.Tx, _ model.SalesTargetDistributorYearly) (int, error) {
	return 0, nil
}

func (s *salesTargetDistributorRepositoryStub) StoreMonthly(_ *sqlx.Tx, _ model.SalesTargetDistributorMonthly) (int, error) {
	if s.storeMonthlyErr != nil {
		return 0, s.storeMonthlyErr
	}
	if s.storeMonthlyID == 0 {
		s.storeMonthlyID = 999
	}
	return s.storeMonthlyID, nil
}

func (s *salesTargetDistributorRepositoryStub) UpdateMonthlyTarget(_ *sqlx.Tx, _ int, _ int, _ int64) error {
	if s.updateMonthlyErr != nil {
		return s.updateMonthlyErr
	}

	return nil
}

func (s *salesTargetDistributorRepositoryStub) UpdateYearly(_ *sqlx.Tx, id int, request entity.UpdateSalesTargetDistributorRequest) error {
	s.lastUpdateID = id
	s.lastUpdateRequest = request
	if s.updateYearlyErr != nil {
		return s.updateYearlyErr
	}
	return nil
}

func (s *salesTargetDistributorRepositoryStub) DeleteMonthlyByYearlyId(_ *sqlx.Tx, _ int, _ int64) error {
	if s.deleteMonthlyErr != nil {
		return s.deleteMonthlyErr
	}
	return nil
}

func (s *salesTargetDistributorRepositoryStub) BeginTx() (*sqlx.Tx, error) {
	if s.beginTxErr != nil {
		return nil, s.beginTxErr
	}

	return s.tx, nil
}

type syncTargetsToMonthlyCall struct {
	custID        string
	yearlyID      int
	month         int
	monthlyID     int
	monthlyTarget int
	updatedBy     int64
}

type syncRecordingSalesTargetRepositoryStub struct {
	salesTargetRepositoryStub
	monthlyAllocationRowsByCustID map[string][]model.SalesTargetMonthlyAllocation
	lastMonthlyAllocationYearlyID int
	lastMonthlyAllocationCustID   string
	syncCalls                     []syncTargetsToMonthlyCall
	syncErr                       error
}

type custAwareSalesTargetRepositoryStub struct {
	salesTargetRepositoryStub
	monthlyAllocationRowsByCustID map[string][]model.SalesTargetMonthlyAllocation
	lastMonthlyAllocationYearlyID int
	lastMonthlyAllocationCustID   string
}

func (s *custAwareSalesTargetRepositoryStub) FindMonthlyAllocationByYearlyId(yearlyID int, custID string) ([]model.SalesTargetMonthlyAllocation, error) {
	s.lastMonthlyAllocationYearlyID = yearlyID
	s.lastMonthlyAllocationCustID = custID
	if s.monthlyAllocationRowsByCustID == nil {
		return nil, nil
	}
	return s.monthlyAllocationRowsByCustID[custID], nil
}

func (s *syncRecordingSalesTargetRepositoryStub) FindMonthlyAllocationByYearlyId(yearlyID int, custID string) ([]model.SalesTargetMonthlyAllocation, error) {
	s.lastMonthlyAllocationYearlyID = yearlyID
	s.lastMonthlyAllocationCustID = custID
	if s.monthlyAllocationRowsByCustID != nil {
		return s.monthlyAllocationRowsByCustID[custID], nil
	}
	return s.salesTargetRepositoryStub.FindMonthlyAllocationByYearlyId(yearlyID, custID)
}

func (s *syncRecordingSalesTargetRepositoryStub) SyncTargetsToMonthly(_ *sqlx.Tx, custID string, yearlyID int, month int, monthlyID int, monthlyTarget int, updatedBy int64) error {
	s.syncCalls = append(s.syncCalls, syncTargetsToMonthlyCall{
		custID:        custID,
		yearlyID:      yearlyID,
		month:         month,
		monthlyID:     monthlyID,
		monthlyTarget: monthlyTarget,
		updatedBy:     updatedBy,
	})

	if s.syncErr != nil {
		return s.syncErr
	}

	return nil
}

func setupSalesTargetDistributorUpdateTx(t *testing.T, expectCommit bool) (*sqlx.Tx, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	mock.ExpectBegin()
	if expectCommit {
		mock.ExpectCommit()
	} else {
		mock.ExpectRollback()
	}

	tx, err := sqlxDB.Beginx()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}
		_ = db.Close()
	}

	return tx, cleanup
}

func buildSalesTargetDistributorUpdateRequest(monthlyTarget int) entity.UpdateSalesTargetDistributorRequest {
	return entity.UpdateSalesTargetDistributorRequest{
		CustId:    "C22001",
		UpdatedBy: 321,
		Data: []entity.UpdateSalesTargetDistributorMonthly{{
			Month:         4,
			MonthlyTarget: monthlyTarget,
		}},
	}
}

func buildSalesTargetDistributorUpdateRequestWithDistributor(monthlyTarget int, distributorID int) entity.UpdateSalesTargetDistributorRequest {
	request := buildSalesTargetDistributorUpdateRequest(monthlyTarget)
	request.DistributorId = &distributorID
	return request
}

func TestSalesTargetDistributorService_List_StatusDerivationMatrix(t *testing.T) {
	currentYear := time.Now().Year()
	updatedBy := int64(77)
	now := time.Now().UTC()

	repo := &salesTargetDistributorRepositoryStub{
		list: []model.SalesTargetDistributorYearly{
			{
				SalesTargetDistributorYearlyId: 1,
				Status:                         int(entity.SALES_TARGET_STATUS_DRAFT),
				Year:                           currentYear,
				IsActive:                       true,
				CreatedBy:                      10,
				CreatedAt:                      now,
			},
			{
				SalesTargetDistributorYearlyId: 2,
				Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
				Year:                           currentYear,
				IsActive:                       true,
				UpdatedBy:                      &updatedBy,
				UpdatedAt:                      &now,
			},
			{
				SalesTargetDistributorYearlyId: 3,
				Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
				Year:                           currentYear + 1,
				IsActive:                       true,
				CreatedBy:                      11,
				CreatedAt:                      now,
			},
			{
				SalesTargetDistributorYearlyId: 4,
				Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
				Year:                           currentYear,
				IsActive:                       false,
				CreatedBy:                      12,
				CreatedAt:                      now,
			},
		},
		total:    4,
		lastPage: 1,
	}

	svc := NewSalesTargetDistributorService(repo, nil)
	data, total, lastPage, err := svc.List(entity.SalesTargetDistributorQueryFilter{}, "C22001")
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
		t.Fatalf("expected draft status, got %s", data[0].Status)
	}
	if data[1].Status != "Active" {
		t.Fatalf("expected active status, got %s", data[1].Status)
	}
	if data[2].Status != "Inactive" {
		t.Fatalf("expected inactive from future year, got %s", data[2].Status)
	}
	if data[3].Status != "Inactive" {
		t.Fatalf("expected inactive from is_active=false, got %s", data[3].Status)
	}
}

func TestSalesTargetDistributorService_List_PassesStatusArrayFilterToRepository(t *testing.T) {
	status := []int{int(entity.SALES_TARGET_STATUS_ACTIVE), int(entity.SALES_TARGET_STATUS_INACTIVE)}
	filter := entity.SalesTargetDistributorQueryFilter{
		Year:   ptrInt(2026),
		Status: &status,
	}

	repo := &salesTargetDistributorRepositoryStub{}
	svc := NewSalesTargetDistributorService(repo, nil)

	_, _, _, err := svc.List(filter, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.lastListCustID != "C22001" {
		t.Fatalf("expected cust id C22001, got %s", repo.lastListCustID)
	}

	if repo.lastListFilter.Status == nil || len(*repo.lastListFilter.Status) != 2 {
		t.Fatalf("expected status filter with 2 values, got %+v", repo.lastListFilter.Status)
	}
}

func TestSalesTargetDistributorService_Detail_WithAllocationSummary(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedBy := int64(99)

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 1,
			Year:                           time.Now().Year(),
			YearlyTarget:                   500000,
			AreaId:                         10,
			AreaCode:                       "AR-01",
			AreaName:                       "Area Jakarta",
			RegionId:                       5,
			RegionCode:                     "RG-01",
			RegionName:                     "Region Barat",
			DistributorId:                  2,
			DistributorCode:                "D001",
			DistributorName:                "Distributor A",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			UpdatedBy:                      &updatedBy,
			UpdatedAt:                      &now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{
			{
				SalesTargetDistributorMonthlyId: 101,
				Month:                           2,
				MonthlyTarget:                   300000,
				IsActive:                        true,
			},
		},
		allocationSummary: model.SalesTargetAllocationSummary{
			IsAllocated:    true,
			AllocatedTotal: 12000000,
		},
	}

	svc := NewSalesTargetDistributorService(repo, &salesTargetRepositoryStub{})
	resp, err := svc.Detail(1, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !resp.IsAllocated {
		t.Fatalf("expected is_allocated true, got false")
	}

	if resp.AllocationTotal <= 0 {
		t.Fatalf("expected allocation_total > 0, got %d", resp.AllocationTotal)
	}

	if len(resp.Details) != 1 {
		t.Fatalf("expected details length 1, got %d", len(resp.Details))
	}

	if !resp.Details[0].IsActive {
		t.Fatalf("expected details[0].is_active true, got false")
	}
}

func TestSalesTargetDistributorService_Detail_AllocationSummaryErrorFallback(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 1,
			Year:                           time.Now().Year(),
			YearlyTarget:                   500000,
			AreaId:                         10,
			AreaCode:                       "AR-01",
			AreaName:                       "Area Jakarta",
			RegionId:                       5,
			RegionCode:                     "RG-01",
			RegionName:                     "Region Barat",
			DistributorId:                  2,
			DistributorCode:                "D001",
			DistributorName:                "Distributor A",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			CreatedBy:                      12,
			CreatedAt:                      now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{
			{
				SalesTargetDistributorMonthlyId: 101,
				Month:                           2,
				MonthlyTarget:                   300000,
				IsActive:                        false,
			},
		},
		allocationErr: errors.New("allocation query timeout"),
	}

	svc := NewSalesTargetDistributorService(repo, &salesTargetRepositoryStub{})
	resp, err := svc.Detail(1, "C22001")
	if err != nil {
		t.Fatalf("expected no error on non-fatal allocation error, got %v", err)
	}

	if resp.IsAllocated {
		t.Fatalf("expected is_allocated false fallback, got true")
	}

	if resp.AllocationTotal != 0 {
		t.Fatalf("expected allocation_total fallback 0, got %d", resp.AllocationTotal)
	}

	if len(resp.Details) != 1 {
		t.Fatalf("expected details length 1, got %d", len(resp.Details))
	}

	if resp.Details[0].IsActive {
		t.Fatalf("expected details[0].is_active false, got true")
	}
}

func TestSalesTargetDistributorResponseContract_JSONKeys(t *testing.T) {
	listBytes, err := json.Marshal(entity.SalesTargetDistributorListResponse{})
	if err != nil {
		t.Fatalf("marshal list response failed: %v", err)
	}

	detailBytes, err := json.Marshal(entity.SalesTargetDistributorDetailResponse{})
	if err != nil {
		t.Fatalf("marshal detail response failed: %v", err)
	}

	if !containsJSONKey(listBytes, "status") {
		t.Fatalf("expected list response to contain key status")
	}
	if !containsJSONKey(detailBytes, "is_allocated") {
		t.Fatalf("expected detail response to contain key is_allocated")
	}
	if !containsJSONKey(detailBytes, "allocation_total") {
		t.Fatalf("expected detail response to contain key allocation_total")
	}
	if containsJSONKey(detailBytes, "allocated_total") {
		t.Fatalf("expected detail response to not contain key allocated_total")
	}

	monthlyBytes, err := json.Marshal(entity.SalesTargetDistributorMonthlyDetail{})
	if err != nil {
		t.Fatalf("marshal monthly detail failed: %v", err)
	}

	for _, key := range []string{"allocated_total", "remaining", "is_allocated", "is_past_month", "is_editable", "disable_reason"} {
		if !containsJSONKey(monthlyBytes, key) {
			t.Fatalf("expected monthly detail to contain key %s", key)
		}
	}
}

func TestApplyStatusTransitionMetadata_StatusOneToTwo_SetsInactiveMetadata(t *testing.T) {
	now := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	status := int(entity.SALES_TARGET_STATUS_INACTIVE)
	request := entity.UpdateSalesTargetDistributorRequest{
		UpdatedBy: 99,
		Status:    &status,
	}

	applyStatusTransitionMetadata(&request, now)

	if request.UserInactive == nil || *request.UserInactive != 99 {
		t.Fatalf("expected user_inactive to be set to 99")
	}
	if request.InactiveAt == nil || !request.InactiveAt.Equal(now) {
		t.Fatalf("expected inactive_at to be set to %s", now)
	}
	if request.IsActive == nil || *request.IsActive {
		t.Fatalf("expected is_active=false when status inactive")
	}
}

func TestApplyStatusTransitionMetadata_StatusZeroToTwo_SetsInactiveMetadata(t *testing.T) {
	now := time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)
	status := int(entity.SALES_TARGET_STATUS_INACTIVE)
	request := entity.UpdateSalesTargetDistributorRequest{
		UpdatedBy: 88,
		Status:    &status,
	}

	applyStatusTransitionMetadata(&request, now)

	if request.UserInactive == nil || *request.UserInactive != 88 {
		t.Fatalf("expected user_inactive to be set to 88")
	}
	if request.InactiveAt == nil || !request.InactiveAt.Equal(now) {
		t.Fatalf("expected inactive_at to be set to %s", now)
	}
	if request.IsActive == nil || *request.IsActive {
		t.Fatalf("expected is_active=false when status inactive")
	}
}

func ptrInt(v int) *int {
	return &v
}

func containsJSONKey(data []byte, key string) bool {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}
	_, ok := m[key]
	return ok
}

func TestSalesTargetDistributorService_Detail_MonthlyFlags(t *testing.T) {
	now := currentTimeUTC()
	currentYear := now.Year()
	pastMonth := int(now.Month()) - 1
	if pastMonth < 1 {
		pastMonth = 1
	}
	futureMonth := int(now.Month()) + 1
	if futureMonth > 12 {
		futureMonth = 12
	}

	salesRepo := &salesTargetRepositoryStub{
		monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{
			{Month: int(now.Month()), AllocatedTotal: 100, TargetCount: 1},
		},
	}

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 1,
			Year:                           currentYear,
			YearlyTarget:                   500000,
			DistributorId:                  2,
			DistributorCode:                "D001",
			DistributorName:                "Distributor A",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			CreatedBy:                      12,
			CreatedAt:                      now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{
			{SalesTargetDistributorMonthlyId: 11, Month: pastMonth, MonthlyTarget: 1000, IsActive: true},
			{SalesTargetDistributorMonthlyId: 12, Month: int(now.Month()), MonthlyTarget: 1000, IsActive: true},
			{SalesTargetDistributorMonthlyId: 13, Month: futureMonth, MonthlyTarget: 1000, IsActive: true},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	resp, err := svc.Detail(1, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Details) != 3 {
		t.Fatalf("expected 3 details, got %d", len(resp.Details))
	}

	if !resp.Details[0].IsPastMonth || resp.Details[0].IsEditable {
		t.Fatalf("expected past month to be non-editable")
	}

	if !resp.Details[1].IsAllocated || resp.Details[1].IsEditable {
		t.Fatalf("expected allocated current month to be non-editable")
	}

	if resp.Details[2].IsPastMonth || resp.Details[2].IsAllocated || !resp.Details[2].IsEditable {
		t.Fatalf("expected future month without allocation to be editable")
	}
}

func TestSalesTargetDistributorService_Detail_ComputesAllocatedEvenWhenChildPointsToOldMonthlyRow(t *testing.T) {
	now := currentTimeUTC()

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			Year:                           now.Year(),
			YearlyTarget:                   4250000,
			DistributorId:                  2,
			DistributorCode:                "D001",
			DistributorName:                "Distributor A",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			CreatedBy:                      12,
			CreatedAt:                      now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   4250000,
			IsActive:                        true,
		}},
	}

	salesRepo := &salesTargetRepositoryStub{
		monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{{
			Month:          4,
			AllocatedTotal: 1250000,
			TargetCount:    1,
		}},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	resp, err := svc.Detail(41, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(resp.Details))
	}

	if !resp.Details[0].IsAllocated {
		t.Fatalf("expected detail month 4 to stay allocated from month anchor")
	}

	if resp.Details[0].AllocatedTotal <= 0 {
		t.Fatalf("expected allocated_total > 0, got %d", resp.Details[0].AllocatedTotal)
	}
}

func TestSalesTargetDistributorService_Detail_UsesDistributorChildCustIDForMonthlyAllocation(t *testing.T) {
	now := currentTimeUTC()

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			CustId:                         "C22001",
			Year:                           now.Year(),
			YearlyTarget:                   4250000,
			DistributorId:                  67,
			DistributorCode:                "D067",
			DistributorName:                "Distributor Child",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			CreatedBy:                      12,
			CreatedAt:                      now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{
			{SalesTargetDistributorMonthlyId: 301, SalesTargetDistributorYearlyId: 41, Month: 3, MonthlyTarget: 1000000, IsActive: true},
			{SalesTargetDistributorMonthlyId: 302, SalesTargetDistributorYearlyId: 41, Month: 4, MonthlyTarget: 1000000, IsActive: true},
		},
		distributorChildCustIDByDistributorID: map[int]string{
			67: "C220010001",
		},
	}

	salesRepo := &custAwareSalesTargetRepositoryStub{
		monthlyAllocationRowsByCustID: map[string][]model.SalesTargetMonthlyAllocation{
			"C220010001": []model.SalesTargetMonthlyAllocation{
				{Month: 3, AllocatedTotal: 250000, TargetCount: 1},
				{Month: 4, AllocatedTotal: 500000, TargetCount: 1},
			},
			"C22001": []model.SalesTargetMonthlyAllocation{},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	resp, err := svc.Detail(41, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Details) != 2 {
		t.Fatalf("expected 2 details, got %d", len(resp.Details))
	}

	if salesRepo.lastMonthlyAllocationYearlyID != 41 {
		t.Fatalf("expected yearly id 41, got %d", salesRepo.lastMonthlyAllocationYearlyID)
	}

	if salesRepo.lastMonthlyAllocationCustID != "C220010001" {
		t.Fatalf("expected monthly allocation to use child cust id C220010001, got %q", salesRepo.lastMonthlyAllocationCustID)
	}

	if !resp.Details[0].IsAllocated || resp.Details[0].AllocatedTotal <= 0 {
		t.Fatalf("expected month 3 allocated detail, got is_allocated=%v allocated_total=%d", resp.Details[0].IsAllocated, resp.Details[0].AllocatedTotal)
	}

	if !resp.Details[1].IsAllocated || resp.Details[1].AllocatedTotal <= 0 {
		t.Fatalf("expected month 4 allocated detail, got is_allocated=%v allocated_total=%d", resp.Details[1].IsAllocated, resp.Details[1].AllocatedTotal)
	}
}

func TestSalesTargetDistributorService_Detail_IsAllocatedTrueWhenTargetCountPositiveEvenAllocatedTotalZero(t *testing.T) {
	now := currentTimeUTC()
	targetMonth := 4

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			Year:                           now.Year(),
			YearlyTarget:                   4250000,
			DistributorId:                  67,
			DistributorCode:                "D067",
			DistributorName:                "Distributor Child",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			CreatedBy:                      12,
			CreatedAt:                      now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 302,
			SalesTargetDistributorYearlyId:  41,
			Month:                           targetMonth,
			MonthlyTarget:                   1000000,
			IsActive:                        true,
		}},
	}

	salesRepo := &salesTargetRepositoryStub{
		monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{{
			Month:          targetMonth,
			AllocatedTotal: 0,
			TargetCount:    1,
		}},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	resp, err := svc.Detail(41, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(resp.Details))
	}

	if !resp.Details[0].IsAllocated {
		t.Fatalf("expected month 4 is_allocated true when target_count > 0 and allocated_total = 0")
	}
}

func TestSalesTargetDistributorService_Detail_SetsDisableReasonAllocated(t *testing.T) {
	now := currentTimeUTC()
	targetMonth := int(now.Month())

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			Year:                           now.Year(),
			YearlyTarget:                   4250000,
			DistributorId:                  2,
			DistributorCode:                "D001",
			DistributorName:                "Distributor A",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			CreatedBy:                      12,
			CreatedAt:                      now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           targetMonth,
			MonthlyTarget:                   4250000,
			IsActive:                        true,
		}},
	}

	salesRepo := &salesTargetRepositoryStub{
		monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{{
			Month:          targetMonth,
			AllocatedTotal: 1250000,
			TargetCount:    1,
		}},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	resp, err := svc.Detail(41, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(resp.Details))
	}

	if resp.Details[0].DisableReason != "allocated" {
		t.Fatalf("expected disable_reason allocated, got %q", resp.Details[0].DisableReason)
	}

	if resp.Details[0].IsEditable {
		t.Fatalf("expected allocated month to be non-editable")
	}
}

func TestSalesTargetDistributorService_Detail_SetsDisableReasonPastMonth(t *testing.T) {
	now := currentTimeUTC()
	targetYear := now.Year()
	targetMonth := int(now.Month()) - 1
	if targetMonth < 1 {
		targetMonth = 12
		targetYear = now.Year() - 1
	}

	repo := &salesTargetDistributorRepositoryStub{
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			Year:                           targetYear,
			YearlyTarget:                   4250000,
			DistributorId:                  2,
			DistributorCode:                "D001",
			DistributorName:                "Distributor A",
			Status:                         int(entity.SALES_TARGET_STATUS_ACTIVE),
			IsActive:                       true,
			CreatedBy:                      12,
			CreatedAt:                      now,
			UpdatedByName:                  "Admin",
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           targetMonth,
			MonthlyTarget:                   4250000,
			IsActive:                        true,
		}},
	}

	svc := NewSalesTargetDistributorService(repo, &salesTargetRepositoryStub{})
	resp, err := svc.Detail(41, "C22001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Details) != 1 {
		t.Fatalf("expected 1 detail, got %d", len(resp.Details))
	}

	if resp.Details[0].DisableReason != "past_month" {
		t.Fatalf("expected disable_reason past_month, got %q", resp.Details[0].DisableReason)
	}

	if resp.Details[0].IsEditable {
		t.Fatalf("expected past month to be non-editable")
	}
}

func TestEnsureMonthlyTargetsConsistent_RejectsOverAllocation(t *testing.T) {
	err := ensureMonthlyTargetsConsistent(100, 150, 7)
	if err == nil {
		t.Fatalf("expected validation error when allocation exceeds monthly target")
	}
}

func TestComputeMonthlyFlags(t *testing.T) {
	currentTime := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	isPastMonth, isEditable := computeMonthlyFlags(2026, 6, currentTime, 0)
	if !isPastMonth || isEditable {
		t.Fatalf("expected past month to be non-editable")
	}

	isPastMonth, isEditable = computeMonthlyFlags(2026, 7, currentTime, 100)
	if isPastMonth || isEditable {
		t.Fatalf("expected allocated current month to be non-editable")
	}

	isPastMonth, isEditable = computeMonthlyFlags(2026, 8, currentTime, 0)
	if isPastMonth || !isEditable {
		t.Fatalf("expected future month without allocation to be editable")
	}
}

func TestSalesTargetDistributorService_Update_SyncsChildMonthlyReferenceToCanonicalRow(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, true)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			CustId:                         "C22001",
			DistributorId:                  67,
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   2000000,
			IsActive:                        true,
		}},
		distributorChildCustIDByDistributorID: map[int]string{
			67: "C220010001",
		},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		salesTargetRepositoryStub: salesTargetRepositoryStub{
			monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{{
				Month:          4,
				AllocatedTotal: 1250000,
				TargetCount:    1,
			}},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	request := buildSalesTargetDistributorUpdateRequest(4250000)
	err := svc.Update(41, request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(salesRepo.syncCalls) != 1 {
		t.Fatalf("expected sync to be called once, got %d", len(salesRepo.syncCalls))
	}

	call := salesRepo.syncCalls[0]
	if call.custID != "C220010001" {
		t.Fatalf("expected sync cust id C220010001, got %s", call.custID)
	}
	if call.yearlyID != 41 || call.month != 4 || call.monthlyID != 745 || call.monthlyTarget != 4250000 || call.updatedBy != request.UpdatedBy {
		t.Fatalf("unexpected sync call captured: %+v", call)
	}
}

func TestSalesTargetDistributorService_Update_SyncsChildMonthlyTargetAndRemaining(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, true)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   2000000,
			IsActive:                        true,
		}},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		salesTargetRepositoryStub: salesTargetRepositoryStub{
			monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{{
				Month:          4,
				AllocatedTotal: 1250000,
				TargetCount:    1,
			}},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	err := svc.Update(41, buildSalesTargetDistributorUpdateRequest(4250000))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(salesRepo.syncCalls) != 1 {
		t.Fatalf("expected sync to be called once, got %d", len(salesRepo.syncCalls))
	}
	if salesRepo.syncCalls[0].monthlyTarget != 4250000 {
		t.Fatalf("expected sync monthly target 4250000, got %d", salesRepo.syncCalls[0].monthlyTarget)
	}
}

func TestSalesTargetDistributorService_Update_UsesDistributorChildCustIDForMonthlyAllocation(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, true)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			CustId:                         "C22001",
			DistributorId:                  67,
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   2000000,
			IsActive:                        true,
		}},
		distributorChildCustIDByDistributorID: map[int]string{
			67: "C220010001",
		},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		monthlyAllocationRowsByCustID: map[string][]model.SalesTargetMonthlyAllocation{
			"C220010001": {{Month: 4, AllocatedTotal: 1250000, TargetCount: 1}},
			"C22001":     {},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	err := svc.Update(41, buildSalesTargetDistributorUpdateRequest(4350000))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if salesRepo.lastMonthlyAllocationYearlyID != 41 {
		t.Fatalf("expected allocation lookup yearly id 41, got %d", salesRepo.lastMonthlyAllocationYearlyID)
	}
	if salesRepo.lastMonthlyAllocationCustID != "C220010001" {
		t.Fatalf("expected allocation lookup cust id C220010001, got %q", salesRepo.lastMonthlyAllocationCustID)
	}
	if len(salesRepo.syncCalls) != 1 {
		t.Fatalf("expected sync to be attempted once, got %d", len(salesRepo.syncCalls))
	}
}

func TestSalesTargetDistributorService_Update_SyncsChildMonthlyReferenceUsingDistributorChildCustID(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, true)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			CustId:                         "C22001",
			DistributorId:                  67,
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   2000000,
			IsActive:                        true,
		}},
		distributorChildCustIDByDistributorID: map[int]string{
			67: "C220010001",
		},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		monthlyAllocationRowsByCustID: map[string][]model.SalesTargetMonthlyAllocation{
			"C220010001": {{Month: 4, AllocatedTotal: 1250000, TargetCount: 1}},
			"C22001":     {},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	request := buildSalesTargetDistributorUpdateRequest(4350000)
	err := svc.Update(41, request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(salesRepo.syncCalls) != 1 {
		t.Fatalf("expected sync to be attempted once, got %d", len(salesRepo.syncCalls))
	}

	call := salesRepo.syncCalls[0]
	if call.custID != "C220010001" {
		t.Fatalf("expected sync cust id C220010001, got %q", call.custID)
	}
	if call.monthlyID != 745 {
		t.Fatalf("expected sync monthly id 745, got %d", call.monthlyID)
	}
	if call.month != 4 {
		t.Fatalf("expected sync month 4, got %d", call.month)
	}
}

func TestSalesTargetDistributorService_Update_UsesPatchedDistributorIDWhenResolvingChildCustID(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, true)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			CustId:                         "C22001",
			DistributorId:                  67,
		},
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   2000000,
			IsActive:                        true,
		}},
		distributorChildCustIDByDistributorID: map[int]string{
			88: "C220010088",
		},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		monthlyAllocationRowsByCustID: map[string][]model.SalesTargetMonthlyAllocation{
			"C220010088": {{Month: 4, AllocatedTotal: 1250000, TargetCount: 1}},
			"C22001":     {},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	err := svc.Update(41, buildSalesTargetDistributorUpdateRequestWithDistributor(4350000, 88))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.lastResolvedDistributorID != 88 {
		t.Fatalf("expected child cust id resolver to use patched distributor id 88, got %d", repo.lastResolvedDistributorID)
	}
}

func TestSalesTargetDistributorService_Update_RejectsTargetBelowAllocatedTotal(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, false)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   2000000,
			IsActive:                        true,
		}},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		salesTargetRepositoryStub: salesTargetRepositoryStub{
			monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{{
				Month:          4,
				AllocatedTotal: 1250000,
				TargetCount:    1,
			}},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	err := svc.Update(41, buildSalesTargetDistributorUpdateRequest(1000000))
	if err == nil {
		t.Fatalf("expected error when monthly target is below allocated total")
	}
	if len(salesRepo.syncCalls) != 0 {
		t.Fatalf("expected sync not to run on validation error, got %d calls", len(salesRepo.syncCalls))
	}
}

func TestSalesTargetDistributorService_Update_RollbackWhenChildSyncFails(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, false)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		monthly: []model.SalesTargetDistributorMonthly{{
			SalesTargetDistributorMonthlyId: 745,
			SalesTargetDistributorYearlyId:  41,
			Month:                           4,
			MonthlyTarget:                   2000000,
			IsActive:                        true,
		}},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		salesTargetRepositoryStub: salesTargetRepositoryStub{
			monthlyAllocationRows: []model.SalesTargetMonthlyAllocation{{
				Month:          4,
				AllocatedTotal: 1250000,
				TargetCount:    1,
			}},
		},
		syncErr: errors.New("sync failed"),
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	err := svc.Update(41, buildSalesTargetDistributorUpdateRequest(4250000))
	if err == nil {
		t.Fatalf("expected error when child sync fails")
	}
	if len(salesRepo.syncCalls) != 1 {
		t.Fatalf("expected sync to be attempted once before rollback, got %d", len(salesRepo.syncCalls))
	}
}

func TestSalesTargetDistributorService_Update_SyncsNewMonthWithoutExistingAllocation(t *testing.T) {
	tx, cleanup := setupSalesTargetDistributorUpdateTx(t, true)
	defer cleanup()

	repo := &salesTargetDistributorRepositoryStub{
		tx: tx,
		yearly: model.SalesTargetDistributorYearly{
			SalesTargetDistributorYearlyId: 41,
			CustId:                         "C22001",
			DistributorId:                  67,
		},
		storeMonthlyID: 901,
		distributorChildCustIDByDistributorID: map[int]string{
			67: "C220010001",
		},
	}
	salesRepo := &syncRecordingSalesTargetRepositoryStub{
		monthlyAllocationRowsByCustID: map[string][]model.SalesTargetMonthlyAllocation{
			"C220010001": {},
		},
	}

	svc := NewSalesTargetDistributorService(repo, salesRepo)
	err := svc.Update(41, entity.UpdateSalesTargetDistributorRequest{
		CustId:    "C22001",
		UpdatedBy: 321,
		Data: []entity.UpdateSalesTargetDistributorMonthly{{
			Month:         5,
			MonthlyTarget: 3500000,
		}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(salesRepo.syncCalls) != 1 {
		t.Fatalf("expected sync to be called once for new month, got %d", len(salesRepo.syncCalls))
	}

	call := salesRepo.syncCalls[0]
	if call.custID != "C220010001" {
		t.Fatalf("expected sync cust id C220010001, got %q", call.custID)
	}
	if call.month != 5 {
		t.Fatalf("expected sync month 5, got %d", call.month)
	}
	if call.monthlyID != 901 {
		t.Fatalf("expected new monthly id 901, got %d", call.monthlyID)
	}
	if call.monthlyTarget != 3500000 {
		t.Fatalf("expected sync monthly target 3500000, got %d", call.monthlyTarget)
	}
	if call.updatedBy != 321 {
		t.Fatalf("expected sync updated_by 321, got %d", call.updatedBy)
	}
	if salesRepo.lastMonthlyAllocationCustID != "C220010001" {
		t.Fatalf("expected allocation lookup cust id C220010001, got %q", salesRepo.lastMonthlyAllocationCustID)
	}
}
