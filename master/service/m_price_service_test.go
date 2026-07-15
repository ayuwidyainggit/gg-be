package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"master/entity"
	"master/model"
	"master/repository"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

func mPriceTestInt64Ptr(value int64) *int64 { return &value }

type mPriceRepositoryStub struct {
	findDetail             model.MPriceDetail
	findDetailErr          error
	findList               []model.MPrice
	findListTotal          int
	findListLastPage       int
	findListErr            error
	snapshotByID           model.MPriceProductSnapshot
	snapshotByIDErr        error
	snapshotByCode         model.MPriceProductSnapshot
	snapshotByCodeErr      error
	stored                 *model.MPrice
	storeErr               error
	publishReq             *entity.PublishByRmqMPriceReq
	publishErr             error
	principalUpdateProID   int64
	principalUpdateDistID  []int64
	principalUpdateDetail  model.MPriceDetail
	principalUpdateErr     error
	distributorUpdateCust  string
	distributorUpdateDist  int64
	distributorUpdatePro   int64
	distributorUpdateErr   error
	affectedDistributorIDs []int64
	affectedErr            error
	updatedDistributorIDs  []int64
	updatedErr             error
	cancelReq              *entity.CancelMPriceParams
	deleteCustID           string
	deletePriceID          string
	deleteUserID           int64
	brokenLinks            []int64
	brokenLinksErr         error
	brokenLinksCallArgs    *struct {
		parentProID    int64
		parentCustID   string
		distributorIDs []int64
	}
}

func (s *mPriceRepositoryStub) FindOneByMPriceIDAndCustID(_ entity.DetailMPriceParams) (model.MPriceDetail, error) {
	if s.findDetailErr != nil {
		return model.MPriceDetail{}, s.findDetailErr
	}
	return s.findDetail, nil
}

func (s *mPriceRepositoryStub) FindAllByCustID(_ entity.MPriceQueryFilter, _ string) ([]model.MPrice, int, int, error) {
	if s.findListErr != nil {
		return nil, 0, 0, s.findListErr
	}
	return append([]model.MPrice(nil), s.findList...), s.findListTotal, s.findListLastPage, nil
}

func (s *mPriceRepositoryStub) Store(mprice *model.MPrice) error {
	if s.storeErr != nil {
		return s.storeErr
	}
	copied := *mprice
	s.stored = &copied
	return nil
}

func (s *mPriceRepositoryStub) Update(_ string, _ entity.UpdateMPriceRequest) error {
	return nil
}

func (s *mPriceRepositoryStub) Delete(custID string, priceID string, deletedBy int64) error {
	s.deleteCustID = custID
	s.deletePriceID = priceID
	s.deleteUserID = deletedBy
	return nil
}

func (s *mPriceRepositoryStub) Cancel(detail entity.CancelMPriceParams) error {
	s.cancelReq = &detail
	return nil
}

func (s *mPriceRepositoryStub) FindOneProductByProID(_ int64, _ string) (model.Product, error) {
	return model.Product{}, nil
}

func (s *mPriceRepositoryStub) FindOneProductSnapshotByProID(_ int64, _ string) (model.MPriceProductSnapshot, error) {
	if s.snapshotByIDErr != nil {
		return model.MPriceProductSnapshot{}, s.snapshotByIDErr
	}
	return s.snapshotByID, nil
}

func (s *mPriceRepositoryStub) FindOneProductSnapshotByCode(_ string, _ string) (model.MPriceProductSnapshot, error) {
	if s.snapshotByCodeErr != nil {
		return model.MPriceProductSnapshot{}, s.snapshotByCodeErr
	}
	return s.snapshotByCode, nil
}

func (s *mPriceRepositoryStub) FindAffectedDistributorProductIDs(_ model.MPriceDetail, _ string) ([]int64, error) {
	if s.affectedErr != nil {
		return nil, s.affectedErr
	}
	return append([]int64(nil), s.affectedDistributorIDs...), nil
}

func (s *mPriceRepositoryStub) FindUpdatedDistributorProductIDs(_ model.MPriceDetail, _ string, _ []int64) ([]int64, error) {
	if s.updatedErr != nil {
		return nil, s.updatedErr
	}
	return append([]int64(nil), s.updatedDistributorIDs...), nil
}

func (s *mPriceRepositoryStub) UpdatePrincipalAssignedProductPrices(parentProID int64, distributorIDs []int64, detail model.MPriceDetail) error {
	s.principalUpdateProID = parentProID
	s.principalUpdateDistID = append([]int64(nil), distributorIDs...)
	s.principalUpdateDetail = detail
	return s.principalUpdateErr
}

func (s *mPriceRepositoryStub) UpdateDistributorProductPrices(custID string, distributorID, proID int64, _ model.MPriceDetail) error {
	s.distributorUpdateCust = custID
	s.distributorUpdateDist = distributorID
	s.distributorUpdatePro = proID
	return s.distributorUpdateErr
}

func (s *mPriceRepositoryStub) PublishByRMQ(request entity.PublishByRmqMPriceReq) error {
	copied := request
	s.publishReq = &copied
	return s.publishErr
}

func (s *mPriceRepositoryStub) FindBrokenDistributorChildLinks(parentProID int64, parentCustID string, distributorIDs []int64) ([]int64, error) {
	s.brokenLinksCallArgs = &struct {
		parentProID    int64
		parentCustID   string
		distributorIDs []int64
	}{
		parentProID:    parentProID,
		parentCustID:   parentCustID,
		distributorIDs: append([]int64(nil), distributorIDs...),
	}
	if s.brokenLinksErr != nil {
		return nil, s.brokenLinksErr
	}
	return append([]int64(nil), s.brokenLinks...), nil
}

type mTransactionPriceRepositoryStub struct {
	deleteByStartDateCalled bool
	deleteByStartDateCustID string
	deleteByStartDateProID  int64
	deleteByStartDateValue  time.Time
	storeCalled             bool
	storeInput              *model.MTransactionPrice
	storeErr                error
	storeBatchCalled        bool
	storeBatchInput         []model.MTransactionPrice
	updateBatchCalled       bool
	updateBatchInput        []model.MTransactionPrice
}

func (s *mTransactionPriceRepositoryStub) Store(mTransPrice *model.MTransactionPrice) error {
	s.storeCalled = true
	copied := *mTransPrice
	s.storeInput = &copied
	return s.storeErr
}

func (s *mTransactionPriceRepositoryStub) StoreBatch(items []model.MTransactionPrice) error {
	s.storeBatchCalled = true
	s.storeBatchInput = append([]model.MTransactionPrice(nil), items...)
	return nil
}

func (s *mTransactionPriceRepositoryStub) GetByCustPro(_ string, _, _ int64, _ string) (*model.MTransactionPrice, error) {
	return nil, errors.New("not found")
}

func (s *mTransactionPriceRepositoryStub) Delete(_ string, _ string) error {
	return nil
}

func (s *mTransactionPriceRepositoryStub) FindAllByCustID(_ entity.MTransactionPriceQueryFilter) ([]model.MTransactionPrice, int, int, error) {
	return nil, 0, 0, nil
}

func (s *mTransactionPriceRepositoryStub) DeleteByStartDate(custID string, proID int64, effectiveDate time.Time) error {
	s.deleteByStartDateCalled = true
	s.deleteByStartDateCustID = custID
	s.deleteByStartDateProID = proID
	s.deleteByStartDateValue = effectiveDate
	return nil
}

func (s *mTransactionPriceRepositoryStub) UpdateBatch(items []model.MTransactionPrice) error {
	s.updateBatchCalled = true
	s.updateBatchInput = append([]model.MTransactionPrice(nil), items...)
	return nil
}

func (s *mTransactionPriceRepositoryStub) StoreIfNotExists(_ *model.MTransactionPrice) error {
	return nil
}

func setupMPriceServiceDistributorRepository(t *testing.T) (repository.DistributorRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewDistributorRepository(sqlxDB)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func expectMPriceManagePricingPermission(mock sqlmock.Sqlmock, distributorID int64, allowed bool) {
	mock.ExpectQuery(`(?s)SELECT COALESCE\(mdist\.allow_manage_pricing, false\) AS allow_manage_pricing.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*mdist\.is_del = false`).
		WithArgs(distributorID).
		WillReturnRows(sqlmock.NewRows([]string{"allow_manage_pricing"}).AddRow(allowed))
}

func TestMPriceService_Detail_ShouldPopulateDetailsAndStatusDesc(t *testing.T) {
	distributorRepo, mock, cleanup := setupMPriceServiceDistributorRepository(t)
	defer cleanup()

	effectiveDate := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	priceRepo := &mPriceRepositoryStub{
		findDetail: model.MPriceDetail{
			CustID:         "CUST001",
			PriceID:        "PRICE-001",
			Coverage:       "D",
			EffectiveDate:  &effectiveDate,
			ProID:          99,
			ProCode:        "PRO-99",
			ProName:        "Produk 99",
			Status:         10,
			DistributorIDs: pq.Int64Array{11},
		},
		updatedDistributorIDs: []int64{11},
	}
	transRepo := &mTransactionPriceRepositoryStub{}

	mock.ExpectQuery(`(?s)SELECT.*FROM mst\.m_distributor d.*WHERE d\.distributor_id IN \(\?\).*AND d\.parent_cust_id = \?.*AND d\.is_del = false`).
		WithArgs("PARENT001", int64(11), "PARENT001").
		WillReturnRows(sqlmock.NewRows([]string{
			"distributor_id", "distributor_code", "distributor_name", "region_id", "region_name", "area_id", "area_name",
		}).AddRow(int64(11), "DIST-11", "Distributor 11", 2, "Region Barat", 4, "Jakarta"))

	svc := NewMPriceService(priceRepo, distributorRepo, transRepo)
	resp, err := svc.Detail(entity.DetailMPriceParams{
		CustID:       "CUST001",
		ParentCustID: "PARENT001",
		PriceID:      "PRICE-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.StatusDesc != "Published" {
		t.Fatalf("expected status desc Published, got %q", resp.StatusDesc)
	}
	if resp.EffectiveDate != "2026-05-05" {
		t.Fatalf("expected formatted effective date, got %q", resp.EffectiveDate)
	}
	if len(resp.Details) != 1 {
		t.Fatalf("expected 1 detail row, got %d", len(resp.Details))
	}
	if resp.Details[0].DistributorCode != "DIST-11" || resp.Details[0].AreaName != "Jakarta" {
		t.Fatalf("unexpected detail payload: %+v", resp.Details[0])
	}
	if resp.Details[0].RegionID != 2 || resp.Details[0].RegionName != "Region Barat" {
		t.Fatalf("expected region detail to be populated, got %+v", resp.Details[0])
	}
	if !resp.Details[0].IsUpdated || resp.Details[0].UpdateStatus != "updated" {
		t.Fatalf("expected updated distributor product status, got %+v", resp.Details[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceService_Detail_ShouldDerivePrincipalWideDistributorProducts(t *testing.T) {
	distributorRepo, mock, cleanup := setupMPriceServiceDistributorRepository(t)
	defer cleanup()

	effectiveDate := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	priceRepo := &mPriceRepositoryStub{
		findDetail: model.MPriceDetail{
			CustID:         "PARENT001",
			PriceID:        "PRICE-001",
			Coverage:       "N",
			EffectiveDate:  &effectiveDate,
			ProID:          99,
			ProCode:        "PRO-99",
			ProName:        "Produk 99",
			Status:         1,
			DistributorIDs: pq.Int64Array{},
		},
		affectedDistributorIDs: []int64{11},
	}
	transRepo := &mTransactionPriceRepositoryStub{}

	mock.ExpectQuery(`(?s)SELECT.*FROM mst\.m_distributor d.*WHERE d\.distributor_id IN \(\?\).*AND d\.parent_cust_id = \?.*AND d\.is_del = false`).
		WithArgs("PARENT001", int64(11), "PARENT001").
		WillReturnRows(sqlmock.NewRows([]string{
			"distributor_id", "distributor_code", "distributor_name", "region_id", "region_name", "area_id", "area_name",
		}).AddRow(int64(11), "DIST-11", "Distributor 11", 2, "Region Barat", 4, "Jakarta"))

	svc := NewMPriceService(priceRepo, distributorRepo, transRepo)
	resp, err := svc.Detail(entity.DetailMPriceParams{
		CustID:       "PARENT001",
		ParentCustID: "PARENT001",
		PriceID:      "PRICE-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Details) != 1 {
		t.Fatalf("expected 1 detail row, got %d", len(resp.Details))
	}
	if resp.Details[0].IsUpdated || resp.Details[0].UpdateStatus != "not_updated" {
		t.Fatalf("expected scheduled price to be not updated, got %+v", resp.Details[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceService_Store_ShouldPersistCustIDAndPublishImmediately(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	effectiveDate := time.Now().In(loc).Format(mPriceDateLayout)
	effectiveDateTime, err := time.ParseInLocation(mPriceDateLayout, effectiveDate, loc)
	if err != nil {
		t.Fatalf("failed to parse effective date: %v", err)
	}

	priceRepo := &mPriceRepositoryStub{
		snapshotByID: model.MPriceProductSnapshot{
			ProID:         77,
			UnitID1:       "PCS",
			UnitID2:       "BOX",
			UnitID3:       "CRT",
			ConvUnit2:     2,
			ConvUnit3:     4,
			PurchPrice1:   10,
			PurchPrice2:   20,
			PurchPrice3:   30,
			SellPrice1:    40,
			SellPrice2:    50,
			SellPrice3:    60,
			DistributorID: nil,
		},
		findDetail: model.MPriceDetail{
			CustID:         "CUST001",
			Coverage:       "N",
			EffectiveDate:  &effectiveDateTime,
			ProID:          77,
			Status:         1,
			NewPurchPrice1: 101,
			NewPurchPrice2: 202,
			NewPurchPrice3: 303,
			NewSellPrice1:  404,
			NewSellPrice2:  505,
			NewSellPrice3:  606,
			DistributorIDs: pq.Int64Array{},
		},
	}
	transRepo := &mTransactionPriceRepositoryStub{}

	svc := &MPriceServiceImpl{
		MPriceRepository:            priceRepo,
		MTransactionPriceRepository: transRepo,
	}

	request := entity.CreateMPriceBody{
		CustID:         "CUST001",
		ParentCustID:   "CUST001",
		CreatedBy:      "Tester",
		CreatedByID:    mPriceTestInt64Ptr(99),
		Coverage:       "N",
		EffectiveDate:  effectiveDate,
		ProID:          77,
		NewPurchPrice1: 101,
		NewPurchPrice2: 202,
		NewPurchPrice3: 303,
		NewSellPrice1:  404,
		NewSellPrice2:  505,
		NewSellPrice3:  606,
	}

	resp, err := svc.Store(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if priceRepo.stored == nil {
		t.Fatal("expected price to be stored")
	}
	if priceRepo.stored.CustID != "CUST001" {
		t.Fatalf("expected stored cust_id CUST001, got %q", priceRepo.stored.CustID)
	}
	if priceRepo.publishReq == nil {
		t.Fatal("expected publish request to be issued for same-day effective date")
	}
	if resp.Status != 10 || resp.StatusDesc != "Published" {
		t.Fatalf("expected immediate publish response, got status=%d desc=%q", resp.Status, resp.StatusDesc)
	}
	if priceRepo.principalUpdateProID != 77 {
		t.Fatalf("expected principal publish to update product 77, got %d", priceRepo.principalUpdateProID)
	}
	if !transRepo.deleteByStartDateCalled {
		t.Fatal("expected transaction price cleanup to run for principal coverage")
	}
	if !transRepo.storeCalled {
		t.Fatal("expected transaction price store to run for principal coverage")
	}
}

func TestMPriceService_Store_ShouldRejectDistributorWithoutManagePricingPermission(t *testing.T) {
	distributorRepo, mock, cleanup := setupMPriceServiceDistributorRepository(t)
	defer cleanup()
	expectMPriceManagePricingPermission(mock, 77, false)

	priceRepo := &mPriceRepositoryStub{
		snapshotByID: model.MPriceProductSnapshot{ProID: 77},
	}
	transRepo := &mTransactionPriceRepositoryStub{}
	svc := NewMPriceService(priceRepo, distributorRepo, transRepo)

	_, err := svc.Store(entity.CreateMPriceBody{
		CustID:         "CUST0010001",
		ParentCustID:   "PARENT001",
		CreatedBy:      "Tester",
		CreatedByID:    mPriceTestInt64Ptr(99),
		DistributorID:  77,
		Coverage:       "D",
		DistributorIDs: []int64{77},
		EffectiveDate:  "2026-05-12",
		ProID:          77,
		NewPurchPrice1: 101,
		NewPurchPrice2: 202,
		NewPurchPrice3: 303,
		NewSellPrice1:  404,
		NewSellPrice2:  505,
		NewSellPrice3:  606,
	})
	if err == nil || err.Error() != errManagePricingNotAllowed.Error() {
		t.Fatalf("expected manage pricing permission error, got %v", err)
	}
	if priceRepo.stored != nil {
		t.Fatalf("expected price not to be stored, got %+v", priceRepo.stored)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceService_Import_ShouldRejectDistributorWithoutManagePricingPermissionBeforeDownload(t *testing.T) {
	distributorRepo, mock, cleanup := setupMPriceServiceDistributorRepository(t)
	defer cleanup()
	expectMPriceManagePricingPermission(mock, 77, false)

	svc := NewMPriceService(&mPriceRepositoryStub{}, distributorRepo, &mTransactionPriceRepositoryStub{})

	resp, err := svc.Import(
		entity.MPriceImportRequest{FileURL: "https://example.test/price.csv"},
		"CUST0010001",
		"PARENT001",
		99,
		77,
		"Tester",
	)
	if err == nil || err.Error() != errManagePricingNotAllowed.Error() {
		t.Fatalf("expected manage pricing permission error, got %v", err)
	}
	if resp.FileURL != "https://example.test/price.csv" || resp.TotalRow != 0 {
		t.Fatalf("expected import to stop before processing rows, got %+v", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMPriceService_Publish_ShouldManuallyPublishDueScheduledPrice(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	effectiveDate := time.Now().In(loc)
	effectiveDate = time.Date(effectiveDate.Year(), effectiveDate.Month(), effectiveDate.Day(), 0, 0, 0, 0, loc)

	priceRepo := &mPriceRepositoryStub{
		findDetail: model.MPriceDetail{
			CustID:         "CUST001",
			PriceID:        "PRICE-001",
			Coverage:       "N",
			EffectiveDate:  &effectiveDate,
			ProID:          77,
			Status:         1,
			NewPurchPrice1: 101,
			NewPurchPrice2: 202,
			NewPurchPrice3: 303,
			NewSellPrice1:  404,
			NewSellPrice2:  505,
			NewSellPrice3:  606,
			DistributorIDs: pq.Int64Array{11},
		},
	}
	transRepo := &mTransactionPriceRepositoryStub{}
	svc := &MPriceServiceImpl{
		MPriceRepository:            priceRepo,
		MTransactionPriceRepository: transRepo,
	}

	userID := int64(99)
	err := svc.Publish(entity.PublishMPriceParams{
		CustID:       "CUST001",
		ParentCustID: "CUST001",
		PriceID:      "PRICE-001",
		UpdatedBy:    "Tester",
		UpdatedByID:  &userID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if priceRepo.publishReq == nil {
		t.Fatal("expected repository PublishByRMQ to be called")
	}
	if priceRepo.publishReq.PriceID != "PRICE-001" || priceRepo.publishReq.Status != 10 {
		t.Fatalf("unexpected publish request: %+v", priceRepo.publishReq)
	}
	if priceRepo.principalUpdateProID != 77 {
		t.Fatalf("expected product prices to be applied, got pro_id %d", priceRepo.principalUpdateProID)
	}
	if !transRepo.storeCalled {
		t.Fatal("expected transaction price sync to store principal transaction price")
	}
}

func TestMPriceService_Publish_ShouldRejectFutureEffectiveDate(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	effectiveDate := time.Now().In(loc).AddDate(0, 0, 1)
	effectiveDate = time.Date(effectiveDate.Year(), effectiveDate.Month(), effectiveDate.Day(), 0, 0, 0, 0, loc)

	priceRepo := &mPriceRepositoryStub{
		findDetail: model.MPriceDetail{
			CustID:        "CUST001",
			PriceID:       "PRICE-001",
			Coverage:      "D",
			EffectiveDate: &effectiveDate,
			ProID:         77,
			Status:        1,
		},
	}
	svc := &MPriceServiceImpl{MPriceRepository: priceRepo}

	err := svc.Publish(entity.PublishMPriceParams{
		CustID:       "CUST001",
		ParentCustID: "PARENT001",
		PriceID:      "PRICE-001",
	})
	if err == nil || err.Error() != "manage price effective date is not due yet" {
		t.Fatalf("expected future effective date error, got %v", err)
	}
	if priceRepo.publishReq != nil {
		t.Fatalf("expected publish not to run, got %+v", priceRepo.publishReq)
	}
}

func TestMPriceService_Publish_ShouldIncludeCurrentStatusWhenRejected(t *testing.T) {
	effectiveDate := time.Now()
	priceRepo := &mPriceRepositoryStub{
		findDetail: model.MPriceDetail{
			CustID:        "CUST001",
			PriceID:       "PRICE-001",
			Coverage:      "D",
			EffectiveDate: &effectiveDate,
			ProID:         77,
			Status:        5,
		},
	}
	svc := &MPriceServiceImpl{MPriceRepository: priceRepo}

	err := svc.Publish(entity.PublishMPriceParams{
		CustID:       "CUST001",
		ParentCustID: "PARENT001",
		PriceID:      "PRICE-001",
	})
	expected := "only scheduled manage prices can be published; current status is Cancelled"
	if err == nil || err.Error() != expected {
		t.Fatalf("expected %q, got %v", expected, err)
	}
}

func TestBuildMPriceHeaderIndex_ShouldSupportGroupedTemplateHeaders(t *testing.T) {
	rows := [][]string{
		{"Effective date", "Product Code", "New Purchase Price", "", "", "New Selling Price", "", "", "Distributor Code"},
		{"", "", "Largest Unit", "Middle Unit", "Smallest Unit", "Largest Unit", "Middle Unit", "Smallest Unit", ""},
	}

	headerMap, dataStartRow := buildMPriceHeaderIndex(rows)
	expected := map[string]int{
		"effective_date":   0,
		"pro_code":         1,
		"new_purch_price3": 2,
		"new_purch_price2": 3,
		"new_purch_price1": 4,
		"new_sell_price3":  5,
		"new_sell_price2":  6,
		"new_sell_price1":  7,
		"distributor_code": 8,
	}

	if dataStartRow != 2 {
		t.Fatalf("expected dataStartRow 2, got %d", dataStartRow)
	}
	if !reflect.DeepEqual(headerMap, expected) {
		t.Fatalf("unexpected header map: %+v", headerMap)
	}
}

func TestMPriceService_Template_ShouldReturnCSVGroupedHeaders(t *testing.T) {
	svc := &MPriceServiceImpl{}

	buf, contentType, filename, err := svc.Template("csv", "CUST001", "CUST001", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if contentType != "text/csv" {
		t.Fatalf("expected csv content type, got %q", contentType)
	}
	if filename != "manage_price_template_principal.csv" {
		t.Fatalf("expected principal csv filename, got %q", filename)
	}

	rows, err := csv.NewReader(bytes.NewReader(buf.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("expected valid csv template, got %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 header rows, got %d", len(rows))
	}
	if rows[0][0] != "Effective date" || rows[0][2] != "New Purchase Price" || rows[0][5] != "New Selling Price" {
		t.Fatalf("unexpected first row: %+v", rows[0])
	}
	if rows[1][2] != "Largest Unit" || rows[1][4] != "Smallest Unit" || rows[1][7] != "Smallest Unit" {
		t.Fatalf("unexpected second row: %+v", rows[1])
	}
}

func TestMPriceService_Template_ShouldReturnXLSXByDefault(t *testing.T) {
	svc := &MPriceServiceImpl{}

	buf, contentType, filename, err := svc.Template("", "CUST001", "CUST001", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Fatalf("unexpected content type: %q", contentType)
	}
	if filename != "manage_price_template.xlsx" {
		t.Fatalf("unexpected filename: %q", filename)
	}
	if len(buf.Bytes()) == 0 {
		t.Fatal("expected non-empty xlsx template")
	}

	workbook, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("expected valid xlsx template, got %v", err)
	}
	sheets := workbook.GetSheetList()
	if !reflect.DeepEqual(sheets, []string{"Template sebagai Principal"}) {
		t.Fatalf("expected only principal sheet, got %+v", sheets)
	}
}

func TestMPriceService_Template_ShouldReturnOnlyDistributorSheetForDistributorToken(t *testing.T) {
	svc := &MPriceServiceImpl{}

	buf, _, _, err := svc.Template("xlsx", "CUST001", "CUST001", 19)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	workbook, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("expected valid xlsx template, got %v", err)
	}
	sheets := workbook.GetSheetList()
	if !reflect.DeepEqual(sheets, []string{"Template sebagai Distributor"}) {
		t.Fatalf("expected only distributor sheet, got %+v", sheets)
	}
}

func TestMPriceService_Template_ShouldRejectUnsupportedFormat(t *testing.T) {
	svc := &MPriceServiceImpl{}

	_, _, _, err := svc.Template("pdf", "CUST001", "CUST001", 0)
	if err == nil || err.Error() != "invalid format" {
		t.Fatalf("expected invalid format error, got %v", err)
	}
}

func TestMPriceService_Template_ShouldReturnDistributorCSVLayoutForDistributorScope(t *testing.T) {
	svc := &MPriceServiceImpl{}

	buf, contentType, filename, err := svc.Template("csv", "CUST0010001", "CUST001", 19)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if contentType != "text/csv" {
		t.Fatalf("unexpected content type: %q", contentType)
	}
	if filename != "manage_price_template_distributor.csv" {
		t.Fatalf("unexpected filename: %q", filename)
	}

	rows, err := csv.NewReader(bytes.NewReader(buf.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("expected valid csv template, got %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if len(rows[0]) != 8 {
		t.Fatalf("expected distributor template without distributor code column, got %d columns", len(rows[0]))
	}
}

func TestExportManagePriceCSV_ShouldUseImprovementNextLayout(t *testing.T) {
	rows := []entity.MPriceResponse{
		{
			EffectiveDate:  "2026-04-30",
			ProCode:        "JY1-005",
			ProName:        "Jersey Persinga Ngawi",
			PurchPrice1:    100,
			PurchPrice2:    200,
			PurchPrice3:    300,
			NewPurchPrice1: 100,
			NewPurchPrice2: 200,
			NewPurchPrice3: 300,
			SellPrice1:     400,
			SellPrice2:     500,
			SellPrice3:     600,
			NewSellPrice1:  440,
			NewSellPrice2:  550,
			NewSellPrice3:  660,
			Details: []entity.DistributorAreaRegionData{
				{
					DistributorCode: "162612",
					DistributorName: "PT Besi Makmur",
					AreaID:          101,
					AreaName:        "JAVA",
					IsUpdated:       true,
				},
				{
					DistributorCode: "787123",
					DistributorName: "CV Abadi Nan Jaya",
					AreaID:          101,
					AreaName:        "JAVA",
				},
			},
		},
	}

	buf, contentType, filename, err := exportManagePriceCSV(rows)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if contentType != "text/csv" || filename != "manage_price_export.csv" {
		t.Fatalf("unexpected export metadata: %q %q", contentType, filename)
	}

	records, err := csv.NewReader(bytes.NewReader(buf.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("expected valid csv export, got %v", err)
	}
	if len(records) != 4 {
		t.Fatalf("expected 2 header rows and 2 detail rows, got %d rows", len(records))
	}
	if records[0][0] != "Effective Date" || records[0][5] != "Previous Purchase Price" || records[0][17] != "Updated" || records[0][19] != "Not Updated" {
		t.Fatalf("unexpected grouped header row: %+v", records[0])
	}
	if records[1][5] != "Largest" || records[1][10] != "Smallest" || records[1][18] != "Distributor Name" {
		t.Fatalf("unexpected unit header row: %+v", records[1])
	}
	if records[2][17] != "162612" || records[2][18] != "PT Besi Makmur" || records[2][19] != "-" || records[2][20] != "-" {
		t.Fatalf("expected first distributor in updated columns, got %+v", records[2])
	}
	if records[2][5] != "N.A" || records[2][8] != "N.A" {
		t.Fatalf("expected unchanged purchase prices to be exported as N.A, got %+v", records[2][5:11])
	}
	if records[3][17] != "-" || records[3][18] != "-" || records[3][19] != "787123" || records[3][20] != "CV Abadi Nan Jaya" {
		t.Fatalf("expected second distributor in not-updated columns, got %+v", records[3])
	}
}

func TestBuildCreateRequestFromImportRow_ShouldDefaultCoverageByDistributorScope(t *testing.T) {
	priceRepo := &mPriceRepositoryStub{
		snapshotByCode: model.MPriceProductSnapshot{ProID: 77},
	}
	svc := &MPriceServiceImpl{MPriceRepository: priceRepo}

	headerMap, _ := buildMPriceHeaderIndex([][]string{
		{"Effective date", "Product Code", "New Purchase Price", "", "", "New Selling Price", "", ""},
		{"", "", "Largest Unit", "Middle Unit", "Smallest Unit", "Largest Unit", "Middle Unit", "Smallest Unit"},
	})

	req, err := svc.buildCreateRequestFromImportRow(
		[]string{"46085", "PRO-77", "300", "200", "100", "600", "500", "400"},
		headerMap,
		"CUST0010001",
		"CUST001",
		99,
		19,
		"Tester",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if req.Coverage != "D" {
		t.Fatalf("expected distributor scope coverage D, got %q", req.Coverage)
	}
	if req.EffectiveDate != "2026-03-04" {
		t.Fatalf("expected excel serial date normalized to 2026-03-04, got %q", req.EffectiveDate)
	}
}

func TestBuildCreateRequestFromImportRow_ShouldFallbackBlankOrZeroPricesToCurrentSnapshot(t *testing.T) {
	priceRepo := &mPriceRepositoryStub{
		snapshotByCode: model.MPriceProductSnapshot{
			ProID:       77,
			PurchPrice1: 10,
			PurchPrice2: 20,
			PurchPrice3: 30,
			SellPrice1:  40,
			SellPrice2:  50,
			SellPrice3:  60,
		},
	}
	svc := &MPriceServiceImpl{MPriceRepository: priceRepo}

	headerMap, _ := buildMPriceHeaderIndex([][]string{
		{"Effective date", "Product Code", "New Purchase Price", "", "", "New Selling Price", "", ""},
		{"", "", "Largest Unit", "Middle Unit", "Smallest Unit", "Largest Unit", "Middle Unit", "Smallest Unit"},
	})

	req, err := svc.buildCreateRequestFromImportRow(
		[]string{"2026-05-20", "PRO-77", "", "0", "11", "0", "", "44"},
		headerMap,
		"CUST001",
		"CUST001",
		99,
		0,
		"Tester",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.NewPurchPrice3 != 30 || req.NewPurchPrice2 != 20 || req.NewPurchPrice1 != 11 {
		t.Fatalf("unexpected purchase prices: %+v", req)
	}
	if req.NewSellPrice3 != 60 || req.NewSellPrice2 != 50 || req.NewSellPrice1 != 44 {
		t.Fatalf("unexpected selling prices: %+v", req)
	}
}

func TestMPriceService_Import_ShouldReturnFailedRows(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte(
			"Effective date,Product Code,New Purchase Price,,,New Selling Price,,,\n" +
				",,Largest Unit,Middle Unit,Smallest Unit,Largest Unit,Middle Unit,Smallest Unit,\n" +
				"2026-05-20,PRO-MISSING,30,20,10,60,50,40,\n",
		))
	}))
	defer server.Close()

	priceRepo := &mPriceRepositoryStub{
		snapshotByCodeErr: errors.New("not found"),
	}
	svc := &MPriceServiceImpl{MPriceRepository: priceRepo}

	resp, err := svc.Import(entity.MPriceImportRequest{FileURL: server.URL}, "CUST001", "CUST001", 99, 0, "Tester")
	if err == nil {
		t.Fatalf("expected import error")
	}
	if resp.TotalRow != 1 || resp.SuccessRow != 0 || resp.FailedRow != 1 {
		t.Fatalf("unexpected import counters: %+v", resp)
	}
	if len(resp.FailedRows) != 1 {
		t.Fatalf("expected one failed row, got %+v", resp.FailedRows)
	}
	expected := "row 3: product PRO-MISSING not found"
	if resp.FailedRows[0] != expected {
		t.Fatalf("expected failed row %q, got %q", expected, resp.FailedRows[0])
	}
}

func TestPublishByRMQ_PrincipalScope_LogsAnomaliesButContinues(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	effectiveDate := time.Now().In(loc)
	effectiveDate = time.Date(effectiveDate.Year(), effectiveDate.Month(), effectiveDate.Day(), 0, 0, 0, 0, loc)

	priceRepo := &mPriceRepositoryStub{
		findDetail: model.MPriceDetail{
			CustID:         "PRINCIPAL001",
			PriceID:        "PRICE-T4",
			Coverage:       "N",
			EffectiveDate:  &effectiveDate,
			ProID:          55,
			Status:         1,
			NewPurchPrice1: 100,
			NewSellPrice1:  200,
			DistributorIDs: pq.Int64Array{11, 22},
		},
		// FindBrokenDistributorChildLinks returns distributor 22 as broken
		brokenLinks: []int64{22},
	}
	transRepo := &mTransactionPriceRepositoryStub{}

	svc := &MPriceServiceImpl{
		MPriceRepository:            priceRepo,
		MTransactionPriceRepository: transRepo,
	}

	err := svc.PublishByRMQ(entity.PublishByRmqMPriceReq{
		CustID:       "PRINCIPAL001",
		ParentCustID: "PRINCIPAL001",
		PriceID:      "PRICE-T4",
		Status:       10,
		UpdatedBy:    "Tester",
	})
	if err != nil {
		t.Fatalf("expected no error even with broken links, got %v", err)
	}

	// UpdatePrincipalAssignedProductPrices must still be called
	if priceRepo.principalUpdateProID != 55 {
		t.Fatalf("expected UpdatePrincipalAssignedProductPrices called with pro_id=55, got %d", priceRepo.principalUpdateProID)
	}

	// FindBrokenDistributorChildLinks must have been called with correct args
	if priceRepo.brokenLinksCallArgs == nil {
		t.Fatal("expected FindBrokenDistributorChildLinks to be called")
	}
	if priceRepo.brokenLinksCallArgs.parentProID != 55 {
		t.Fatalf("expected broken links query with parent_pro_id=55, got %d", priceRepo.brokenLinksCallArgs.parentProID)
	}
	if len(priceRepo.brokenLinksCallArgs.distributorIDs) != 2 {
		t.Fatalf("expected broken links query with 2 distributor IDs, got %v", priceRepo.brokenLinksCallArgs.distributorIDs)
	}

	// principal coverage=N → syncTransactionPrices takes the N path: DeleteByStartDate + Store
	if !transRepo.deleteByStartDateCalled {
		t.Fatal("expected DeleteByStartDate to be called for principal coverage N")
	}
	if !transRepo.storeCalled {
		t.Fatal("expected Store to be called for principal coverage N")
	}

	// repo PublishByRMQ must be called
	if priceRepo.publishReq == nil {
		t.Fatal("expected repository PublishByRMQ to be called")
	}
}

func TestPublishByRMQ_DoesNotDuplicateTransactionPriceToChildProID(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	effectiveDate := time.Now().In(loc)
	effectiveDate = time.Date(effectiveDate.Year(), effectiveDate.Month(), effectiveDate.Day(), 0, 0, 0, 0, loc)

	priceRepo := &mPriceRepositoryStub{
		findDetail: model.MPriceDetail{
			CustID:         "PRINCIPAL001",
			PriceID:        "PRICE-T4B",
			Coverage:       "N",
			EffectiveDate:  &effectiveDate,
			ProID:          77,
			Status:         1,
			NewPurchPrice1: 150,
			NewSellPrice1:  300,
			DistributorIDs: pq.Int64Array{},
		},
	}
	transRepo := &mTransactionPriceRepositoryStub{}

	svc := &MPriceServiceImpl{
		MPriceRepository:            priceRepo,
		MTransactionPriceRepository: transRepo,
	}

	err := svc.PublishByRMQ(entity.PublishByRmqMPriceReq{
		CustID:       "PRINCIPAL001",
		ParentCustID: "PRINCIPAL001",
		PriceID:      "PRICE-T4B",
		Status:       10,
		UpdatedBy:    "Tester",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !transRepo.storeCalled {
		t.Fatal("expected Store to be called")
	}
	// ProID in stored transaction price must equal price.ProID, not any child ID
	if transRepo.storeInput == nil {
		t.Fatal("expected storeInput to be captured")
	}
	if transRepo.storeInput.ProID != 77 {
		t.Fatalf("expected stored transaction price ProID=77 (parent), got %d", transRepo.storeInput.ProID)
	}
}
