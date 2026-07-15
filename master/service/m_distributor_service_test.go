package service

import (
	"database/sql"
	"errors"
	"testing"

	"master/entity"
	"master/pkg/constant"
	"master/repository"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func setupDistributorServiceTest(t *testing.T) (*distributorServiceImpl, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewDistributorRepository(sqlxDB)
	svc := NewDistributorService(repo)

	cleanup := func() {
		_ = db.Close()
	}

	return svc, mock, cleanup
}

func serviceStrPtr(v string) *string { return &v }
func serviceIntPtr(v int) *int       { return &v }
func serviceInt64Ptr(v int64) *int64 { return &v }
func serviceBoolPtr(v bool) *bool    { return &v }

func expectResolvedDistributorCustID(mock sqlmock.Sqlmock, distributorID int64, scopeCustID, resolvedCustID string) {
	detailRows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "longitude",
		"phone", "fax_number", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "updated_by_name", "is_active", "is_del",
		"created_by", "created_at", "updated_by", "updated_at", "deleted_by", "deleted_at", "region_code", "region_name",
		"channel_code", "channel_name", "sub_distributor_group_code", "sub_distributor_group_name", "dist_price_grp_code", "dist_price_grp_name",
		"allow_add_product", "allow_edit_product", "allow_manage_pricing", "allow_upload_secondary_sales",
	}).AddRow(
		resolvedCustID, scopeCustID, distributorID, "DIST001", "Distributor Existing", nil,
		1, 2, 3, 4, 5, "Alamat Existing",
		"31", "3171", "3171010", "3171010001", "12345", 10, "-6.2", "106.8",
		nil, nil, "AR01", "Area", "31", "DKI Jakarta", "3171", "Jakarta Selatan",
		"3171010", "Kebayoran Baru", "3171010001", "Senayan", "Updater", true, false,
		nil, nil, nil, nil, nil, nil, "RG01", "Region",
		"CH01", "Channel", "SDG01", "Sub Group", "DPG01", "Price Group",
		false, false, false, false,
	)

	mock.ExpectQuery(`(?s)SELECT .*mdist\.parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*mdist\.cust_id LIKE \$2.*`).
		WithArgs(distributorID, scopeCustID+"%").
		WillReturnRows(detailRows)
}

func validStoreRequest() entity.CreateDistributorBody {
	return entity.CreateDistributorBody{
		CustId:                "CUST001",
		ParentCustId:          "CUST001",
		CreatedBy:             serviceInt64Ptr(1001),
		DistributorCode:       "DIST001",
		DistributorName:       "Distributor One",
		RegionId:              1,
		AreaId:                2,
		ChannelId:             3,
		SubDistributorGroupId: 4,
		DistPriceGrpId:        5,
		Address:               "Jl Test",
		Latitude:              "-6.2000",
		Longitude:             "106.8166",
		IsActive:              serviceBoolPtr(true),
	}
}

func validUpdateRequest() entity.UpdateDistributorRequest {
	return entity.UpdateDistributorRequest{
		CustId:                "CUST001",
		UpdatedBy:             serviceInt64Ptr(1002),
		DistributorCode:       serviceStrPtr("DIST001"),
		DistributorName:       serviceStrPtr("Distributor Updated"),
		RegionId:              serviceIntPtr(1),
		AreaId:                serviceIntPtr(2),
		ChannelId:             serviceIntPtr(3),
		SubDistributorGroupId: serviceIntPtr(4),
		DistPriceGrpId:        serviceIntPtr(5),
		Address:               serviceStrPtr("Jl Updated"),
		Latitude:              serviceStrPtr("-6.2100"),
		Longitude:             serviceStrPtr("106.8266"),
		IsActive:              serviceBoolPtr(true),
	}
}

func TestDistributorService_Store_DuplicateMessageStandardized(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST001", "CUST001").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id", "distributor_code"}).AddRow(10, "DIST001"))

	_, err := svc.Store(validStoreRequest())
	if err == nil || err.Error() != "Distributor code already exists. Please use a different distributor code." {
		t.Fatalf("expected standardized duplicate message, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Store_DBUniqueViolationMapped(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST001", "CUST001").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO mst\\.m_distributor").
		WillReturnError(&pq.Error{Code: "23505"})
	mock.ExpectRollback()

	_, err := svc.Store(validStoreRequest())
	if err == nil || err.Error() != "Distributor code already exists. Please use a different distributor code." {
		t.Fatalf("expected unique violation mapped message, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_ContactDeleteGuardUsesContactIDs(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	contactID := int64(101)
	request := validUpdateRequest()
	request.Contacts = []entity.DistributorContactUpdate{{
		DistributorContactId: &contactID,
		ContactName:          serviceStrPtr("Jane"),
	}}
	request.Tax = nil
	expectResolvedDistributorCustID(mock, 1, "CUST001", "CUST001")

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST001", "CUST001").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id", "distributor_code"}).AddRow(1, "DIST001"))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM mst\\.m_distributor_contact").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE mst\\.m_distributor_contact").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(1, request)
	if err != nil {
		t.Fatalf("expected update success, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_NoRowsAffectedReturnsNotFoundError(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:          "CUST001",
		UpdatedBy:       serviceInt64Ptr(1002),
		DistributorName: serviceStrPtr("Distributor Missing"),
	}

	mock.ExpectQuery(`(?s)SELECT .*mdist\.parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*mdist\.cust_id LIKE \$2.*`).
		WithArgs(999, "CUST001%").
		WillReturnError(sql.ErrNoRows)

	err := svc.Update(999, request)
	if !errors.Is(err, constant.ErrNoRowsAffected) {
		t.Fatalf("expected no rows affected error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_ParentScopeUsesChildCustIDForUpdate(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := validUpdateRequest()
	request.CustId = "C26002"
	request.DistributorCode = serviceStrPtr("162612")

	detailRows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "longitude",
		"phone", "fax_number", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "updated_by_name", "is_active", "is_del",
		"created_by", "created_at", "updated_by", "updated_at", "deleted_by", "deleted_at", "region_code", "region_name",
		"channel_code", "channel_name", "sub_distributor_group_code", "sub_distributor_group_name", "dist_price_grp_code", "dist_price_grp_name",
		"allow_add_product", "allow_edit_product", "allow_manage_pricing", "allow_upload_secondary_sales",
	}).AddRow(
		"C260020001", "C26002", int64(102), "162612", "PT. Besi Makmur", nil,
		1, 2, 3, 4, 5, "Alamat Child",
		"31", "3171", "3171010", "3171010001", "12345", 10, "-6.2", "106.8",
		nil, nil, "AR01", "Area", "31", "DKI Jakarta", "3171", "Jakarta Selatan",
		"3171010", "Kebayoran Baru", "3171010001", "Senayan", "Updater", true, false,
		nil, nil, nil, nil, nil, nil, "RG01", "Region",
		"CH01", "Channel", "SDG01", "Sub Group", "DPG01", "Price Group",
		false, false, false, false,
	)

	mock.ExpectQuery(`(?s)SELECT .*mdist\.parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*mdist\.cust_id LIKE \$2.*`).
		WithArgs(102, "C26002%").
		WillReturnRows(detailRows)
	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("162612", "C260020001").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id", "distributor_code"}).AddRow(102, "162612"))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE mst\.m_distributor`).
		WithArgs(
			"162612",
			"Distributor Updated",
			1,
			2,
			3,
			4,
			5,
			"Jl Updated",
			"-6.2100",
			"106.8266",
			true,
			int64(1002),
			"C260020001",
			int64(102),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(102, request)
	if err != nil {
		t.Fatalf("expected update success for parent scope, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_EmptyBarcodeAndZipCodeAreIncludedInPatch(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:          "CUST001",
		UpdatedBy:       serviceInt64Ptr(1002),
		DistributorName: serviceStrPtr("Distributor Nullable Update"),
		Barcode:         serviceStrPtr(""),
		ZipCode:         serviceStrPtr(""),
		BarcodeProvided: true,
		ZipCodeProvided: true,
	}
	expectResolvedDistributorCustID(mock, 103, "CUST001", "CUST001")

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE mst\.m_distributor`).
		WithArgs(
			"Distributor Nullable Update",
			nil,
			nil,
			int64(1002),
			"CUST001",
			int64(103),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(103, request)
	if err != nil {
		t.Fatalf("expected update success with nullable barcode and zip_code, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_AllOptionalNullableFieldsAreIncludedInPatch(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:                "CUST001",
		UpdatedBy:             serviceInt64Ptr(1002),
		DistributorName:       serviceStrPtr("Distributor Nullable Full Update"),
		ProvinceId:            serviceStrPtr(""),
		RegencyId:             serviceStrPtr(""),
		SubDistrictId:         serviceStrPtr(""),
		WardId:                serviceStrPtr(""),
		OtLocId:               serviceIntPtr(0),
		Phone:                 serviceStrPtr(""),
		FaxNumber:             serviceStrPtr(""),
		ProvinceIdProvided:    true,
		RegencyIdProvided:     true,
		SubDistrictIdProvided: true,
		WardIdProvided:        true,
		OtLocIdProvided:       true,
		PhoneProvided:         true,
		FaxNumberProvided:     true,
	}
	expectResolvedDistributorCustID(mock, 103, "CUST001", "CUST001")

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE mst\.m_distributor`).
		WithArgs(
			"Distributor Nullable Full Update",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			int64(1002),
			"CUST001",
			int64(103),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(103, request)
	if err != nil {
		t.Fatalf("expected update success with all optional nullable fields, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_PrincipalDistributorWithNullParentCustIDStillResolvesCustID(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:          "C26002",
		UpdatedBy:       serviceInt64Ptr(1002),
		DistributorName: serviceStrPtr("Principal Distributor Update"),
	}

	detailRows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "longitude",
		"phone", "fax_number", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "updated_by_name", "is_active", "is_del",
		"created_by", "created_at", "updated_by", "updated_at", "deleted_by", "deleted_at", "region_code", "region_name",
		"channel_code", "channel_name", "sub_distributor_group_code", "sub_distributor_group_name", "dist_price_grp_code", "dist_price_grp_name",
		"allow_add_product", "allow_edit_product", "allow_manage_pricing", "allow_upload_secondary_sales",
	}).AddRow(
		"C26002", "", int64(123), "AUTO26002LOCAL01", "Principal Distributor", nil,
		1, 2, 3, 4, 5, "Alamat Principal",
		"31", "3171", "3171010", "3171010001", "12345", 10, "-6.2", "106.8",
		nil, nil, "AR01", "Area", "31", "DKI Jakarta", "3171", "Jakarta Selatan",
		"3171010", "Kebayoran Baru", "3171010001", "Senayan", "Updater", true, false,
		nil, nil, nil, nil, nil, nil, "RG01", "Region",
		"CH01", "Channel", "SDG01", "Sub Group", "DPG01", "Price Group",
		false, false, false, false,
	)

	mock.ExpectQuery(`(?s)SELECT .*COALESCE\(mdist\.parent_cust_id, ''\) AS parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*mdist\.cust_id LIKE \$2.*`).
		WithArgs(123, "C26002%").
		WillReturnRows(detailRows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(123, request)
	if err != nil {
		t.Fatalf("expected update success for principal distributor with null parent_cust_id, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_EmptyContactsDeletesAllExistingContacts(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := validUpdateRequest()
	request.Contacts = []entity.DistributorContactUpdate{}
	request.Tax = nil
	expectResolvedDistributorCustID(mock, 1, "CUST001", "CUST001")

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST001", "CUST001").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id", "distributor_code"}).AddRow(1, "DIST001"))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM mst\\.m_distributor_contact WHERE distributor_id = \\$1;").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	err := svc.Update(1, request)
	if err != nil {
		t.Fatalf("expected update success, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_EmptyTaxDeletesAllExistingTax(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := validUpdateRequest()
	request.Contacts = nil
	request.Tax = []entity.DistributorTaxUpdate{}
	expectResolvedDistributorCustID(mock, 1, "CUST001", "CUST001")

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST001", "CUST001").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id", "distributor_code"}).AddRow(1, "DIST001"))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM mst\\.m_distributor_tax WHERE distributor_id = \\$1;").WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(1, request)
	if err != nil {
		t.Fatalf("expected update success, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_RollsBackWhenTaxUpdateFails(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	taxID := int64(77)
	request := validUpdateRequest()
	request.Contacts = nil
	request.Tax = []entity.DistributorTaxUpdate{{
		DistributorTaxId: &taxID,
		TaxName:          serviceStrPtr("Updated Tax"),
	}}
	expectResolvedDistributorCustID(mock, 1, "CUST001", "CUST001")

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST001", "CUST001").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id", "distributor_code"}).AddRow(1, "DIST001"))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM mst\\.m_distributor_tax").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE mst\\.m_distributor_tax").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := svc.Update(1, request)
	if !errors.Is(err, constant.ErrNoRowsAffected) {
		t.Fatalf("expected rollback with no rows affected, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_List_NullLocationIDs_DoesNotError(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10, ParentCustId: "C22001"}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude",
		"area_code", "area_name",
		"province_code", "province_name",
		"regency_code", "regency_name",
		"sub_district_code", "sub_district_name",
		"ward_code", "ward_name",
		"longitude", "region_code", "region_name", "channel_code", "channel_name",
		"dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name",
	}).AddRow(
		"C22001", int64(96), "98760", "Distributor Buah Segar", nil,
		0, 0, 0, 0, 0, "Alamat Test",
		"", "", "", "", "", 0, "-6.2",
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		"106.8", nil, nil, nil, nil,
		nil, nil, true, nil,
	)

	mock.ExpectQuery("(?s)SELECT .*COALESCE\\(sp\\.region_id, 0\\) AS region_id.*COALESCE\\(sp\\.area_id, 0\\) AS area_id.*COALESCE\\(sp\\.channel_id, 0\\) AS channel_id.*COALESCE\\(sp\\.sub_distributor_group_id, 0\\) AS sub_distributor_group_id.*COALESCE\\(sp\\.dist_price_grp_id, 0\\) AS dist_price_grp_id.*COALESCE\\(sp\\.address, ''\\) AS address.*COALESCE\\(sp\\.province_id, ''\\) AS province_id.*COALESCE\\(sp\\.regency_id, ''\\) AS regency_id.*COALESCE\\(sp\\.sub_district_id, ''\\) AS sub_district_id.*COALESCE\\(sp\\.ward_id, ''\\) AS ward_id.*COALESCE\\(sp\\.zip_code, ''\\) AS zip_code.*COALESCE\\(sp\\.ot_loc_id, 0\\) AS ot_loc_id.*COALESCE\\(sp\\.latitude, ''\\) AS latitude.*COALESCE\\(sp\\.longitude, ''\\) AS longitude.*COALESCE\\(sp\\.is_active, false\\) AS is_active.*FROM mst\\.m_distributor sp.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.List(filter, "C22001")
	if err != nil {
		t.Fatalf("expected list success with NULL location ids, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}
	if data[0].RegionId != 0 || data[0].AreaId != 0 || data[0].ChannelId != 0 || data[0].SubDistributorGroupId != 0 || data[0].DistPriceGrpId != 0 || data[0].OtLocId != 0 {
		t.Fatalf("expected zero value for nullable integer ids, got region=%d area=%d channel=%d sub_group=%d dist_price_group=%d ot_loc=%d", data[0].RegionId, data[0].AreaId, data[0].ChannelId, data[0].SubDistributorGroupId, data[0].DistPriceGrpId, data[0].OtLocId)
	}
	if data[0].ProvinceId != "" || data[0].RegencyId != "" || data[0].SubDistrictId != "" || data[0].WardId != "" {
		t.Fatalf("expected empty string for nullable location ids, got province=%q regency=%q sub_district=%q ward=%q", data[0].ProvinceId, data[0].RegencyId, data[0].SubDistrictId, data[0].WardId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_List_NullStringAndBoolFields_DoesNotError(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10, ParentCustId: "C22001"}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude",
		"area_code", "area_name",
		"province_code", "province_name",
		"regency_code", "regency_name",
		"sub_district_code", "sub_district_name",
		"ward_code", "ward_name",
		"longitude", "region_code", "region_name", "channel_code", "channel_name",
		"dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name",
	}).AddRow(
		"C22001", int64(97), "98761", "Distributor Null Strings", nil,
		0, 0, 0, 0, 0, "",
		"", "", "", "", "", 0, "",
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		"", nil, nil, nil, nil,
		nil, nil, false, nil,
	)

	mock.ExpectQuery("(?s)SELECT .*COALESCE\\(sp\\.address, ''\\) AS address.*COALESCE\\(sp\\.zip_code, ''\\) AS zip_code.*COALESCE\\(sp\\.latitude, ''\\) AS latitude.*COALESCE\\(sp\\.longitude, ''\\) AS longitude.*COALESCE\\(sp\\.is_active, false\\) AS is_active.*FROM mst\\.m_distributor sp.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.List(filter, "C22001")
	if err != nil {
		t.Fatalf("expected list success with NULL string and bool fields, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}
	if data[0].Address != "" || data[0].ZipCode != "" || data[0].Latitude != "" || data[0].Longitude != "" {
		t.Fatalf("expected empty string fallback for nullable string fields, got address=%q zip_code=%q latitude=%q longitude=%q", data[0].Address, data[0].ZipCode, data[0].Latitude, data[0].Longitude)
	}
	if data[0].IsActive {
		t.Fatalf("expected false fallback for nullable is_active, got true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_ListWithCustomer_NullIntegerIDs_DoesNotError(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10}

	mock.ExpectQuery(`SELECT\s+COUNT\(DISTINCT sp\.distributor_id\)\s+AS total`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "ot_loc_id",
		"is_active", "customer_id",
	}).AddRow(
		"C22001", int64(96), "98760", "Distributor Customer",
		0, 0, 0, 0, 0, 0,
		true, "CUS-001",
	)

	mock.ExpectQuery("(?s)SELECT .*COALESCE\\(sp\\.region_id, 0\\) AS region_id.*COALESCE\\(sp\\.area_id, 0\\) AS area_id.*COALESCE\\(sp\\.channel_id, 0\\) AS channel_id.*COALESCE\\(sp\\.sub_distributor_group_id, 0\\) AS sub_distributor_group_id.*COALESCE\\(sp\\.dist_price_grp_id, 0\\) AS dist_price_grp_id.*COALESCE\\(sp\\.address, ''\\) AS address.*COALESCE\\(sp\\.province_id, ''\\) AS province_id.*COALESCE\\(sp\\.regency_id, ''\\) AS regency_id.*COALESCE\\(sp\\.sub_district_id, ''\\) AS sub_district_id.*COALESCE\\(sp\\.ward_id, ''\\) AS ward_id.*COALESCE\\(sp\\.zip_code, ''\\) AS zip_code.*COALESCE\\(sp\\.ot_loc_id, 0\\) AS ot_loc_id.*COALESCE\\(sp\\.latitude, ''\\) AS latitude.*COALESCE\\(sp\\.longitude, ''\\) AS longitude.*COALESCE\\(sp\\.is_active, false\\) AS is_active.*COALESCE\\(mc\\.cust_id, ''\\) AS customer_id.*FROM mst\\.m_distributor sp.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.ListWithCustomer(filter, "C22001")
	if err != nil {
		t.Fatalf("expected list with customer success with NULL integer ids, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}
	if data[0].RegionId != 0 || data[0].AreaId != 0 || data[0].ChannelId != 0 || data[0].SubDistributorGroupId != 0 || data[0].DistPriceGrpId != 0 || data[0].OtLocId != 0 {
		t.Fatalf("expected zero value for nullable integer ids, got region=%d area=%d channel=%d sub_group=%d dist_price_group=%d ot_loc=%d", data[0].RegionId, data[0].AreaId, data[0].ChannelId, data[0].SubDistributorGroupId, data[0].DistPriceGrpId, data[0].OtLocId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_ListWithCustomer_NullStringFields_DoesNotError(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10}

	mock.ExpectQuery(`SELECT\s+COUNT\(DISTINCT sp\.distributor_id\)\s+AS total`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name",
		"address", "zip_code", "latitude", "longitude", "is_active", "customer_id",
	}).AddRow(
		"C22001", int64(98), "98762", "Distributor Customer Null",
		"", "", "", "", false, "",
	)

	mock.ExpectQuery("(?s)SELECT .*COALESCE\\(sp\\.address, ''\\) AS address.*COALESCE\\(sp\\.zip_code, ''\\) AS zip_code.*COALESCE\\(sp\\.latitude, ''\\) AS latitude.*COALESCE\\(sp\\.longitude, ''\\) AS longitude.*COALESCE\\(sp\\.is_active, false\\) AS is_active.*COALESCE\\(mc\\.cust_id, ''\\) AS customer_id.*FROM mst\\.m_distributor sp.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.ListWithCustomer(filter, "C22001")
	if err != nil {
		t.Fatalf("expected list with customer success with NULL string fields, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}
	if data[0].Address != "" || data[0].ZipCode != "" || data[0].Latitude != "" || data[0].Longitude != "" || data[0].CustomerID != "" {
		t.Fatalf("expected empty string fallback for nullable customer list fields, got address=%q zip_code=%q latitude=%q longitude=%q customer_id=%q", data[0].Address, data[0].ZipCode, data[0].Latitude, data[0].Longitude, data[0].CustomerID)
	}
	if data[0].IsActive {
		t.Fatalf("expected false fallback for nullable is_active, got true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_ListWithCustomer_DeduplicatesDistributorRows(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10}

	mock.ExpectQuery(`SELECT\s+COUNT\(DISTINCT sp\.distributor_id\)\s+AS total`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id",
		"address", "province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id",
		"latitude", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "longitude", "region_code", "region_name",
		"channel_code", "channel_name", "dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name", "customer_id",
	}).AddRow(
		"C260020001", "C26002", int64(102), "162612", "PT Besi Makmur",
		81, 88, 46, 37, 46,
		"jalan bm", "1", "01", "10101", "1010101", "11", 0,
		"-6.2088", "101", "JAVA", "1", "JAWA BARAT", "01", "CIREBON",
		"10101", "HARJAMUKTI", "1010101", "LARANGAN", "106.8456", "1", "CENTRAL",
		"01", "NON", "1", "NON", true, "Princessa Ahsani Taqwim", "C260020001",
	)

	mock.ExpectQuery(`(?s)SELECT .*DISTINCT ON \(sp\.distributor_id\).*COALESCE\(mc\.cust_id, ''\) AS customer_id.*FROM mst\.m_distributor sp.*LIMIT 10 OFFSET 0`).
		WillReturnRows(rows)

	data, total, lastPage, err := svc.ListWithCustomer(filter, "C26002")
	if err != nil {
		t.Fatalf("expected deduplicated list with customer success, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row after deduplication, got %d", len(data))
	}
	if data[0].DistributorId != 102 {
		t.Fatalf("expected distributor_id 102, got %d", data[0].DistributorId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_List_DeduplicatesDistributorRows(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10}

	mock.ExpectQuery(`SELECT\s+COUNT\(DISTINCT sp\.distributor_id\)\s+AS total`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude",
		"area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "longitude", "region_code", "region_name",
		"channel_code", "channel_name", "dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name",
	}).AddRow(
		"C260020001", "C26002", int64(102), "162612", "PT Besi Makmur", "",
		81, 88, 46, 37, 46, "jalan besi makmur madura",
		"1", "01", "10101", "1010101", "11", 0, "-6.2088",
		"101", "JAVA", "1", "JAWA BARAT", "01", "CIREBON",
		"10101", "HARJAMUKTI", "1010101", "LARANGAN", "106.8456", "1", "CENTRAL",
		"01", "NON", "1", "NON", true, "Princessa Ahsani Taqwim",
	)

	mock.ExpectQuery(`(?s)SELECT .*DISTINCT ON \(sp\.distributor_id\).*COALESCE\(pv\.province, ''\) AS province_name.*FROM mst\.m_distributor sp.*ORDER BY sp\.distributor_id DESC.*LIMIT 10 OFFSET 0`).
		WillReturnRows(rows)

	data, total, lastPage, err := svc.List(filter, "C26002")
	if err != nil {
		t.Fatalf("expected deduplicated list success, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row after deduplication, got %d", len(data))
	}
	if data[0].DistributorId != 102 {
		t.Fatalf("expected distributor_id 102, got %d", data[0].DistributorId)
	}
	if data[0].ProvinceName != "JAWA BARAT" {
		t.Fatalf("expected province_name JAWA BARAT, got %q", data[0].ProvinceName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_List_PrincipalScopeUsesCustIDPrefixAndDefaultActiveTrue(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10, ParentCustId: "P22001", JwtDistributorId: 0}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude",
		"area_code", "area_name",
		"province_code", "province_name",
		"regency_code", "regency_name",
		"sub_district_code", "sub_district_name",
		"ward_code", "ward_name",
		"longitude", "region_code", "region_name", "channel_code", "channel_name",
		"dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name",
	}).AddRow(
		"C22002", int64(96), "98760", "Distributor Principal Scope", nil,
		1, 1, 1, 1, 1, "Alamat Test",
		"", "", "", "", "", 0, "-6.2",
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		"106.8", nil, nil, nil, nil,
		nil, nil, true, nil,
	)

	mock.ExpectQuery("(?s)SELECT .*FROM mst\\.m_distributor sp.*sp\\.cust_id LIKE 'P22001%'.*sp\\.is_active = true.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.List(filter, "P22001")
	if err != nil {
		t.Fatalf("expected list success for principal prefix scope, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_List_DistributorScopeUsesCustIDPrefixAndDefaultActiveTrue(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10, ParentCustId: "P22001", JwtDistributorId: 44}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude",
		"area_code", "area_name",
		"province_code", "province_name",
		"regency_code", "regency_name",
		"sub_district_code", "sub_district_name",
		"ward_code", "ward_name",
		"longitude", "region_code", "region_name", "channel_code", "channel_name",
		"dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name",
	}).AddRow(
		"C22001", int64(44), "98761", "Distributor Exact Scope", nil,
		1, 1, 1, 1, 1, "Alamat Test",
		"", "", "", "", "", 0, "-6.2",
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		"106.8", nil, nil, nil, nil,
		nil, nil, true, nil,
	)

	mock.ExpectQuery("(?s)SELECT .*FROM mst\\.m_distributor sp.*sp\\.cust_id LIKE 'C22001%'.*sp\\.is_active = true.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.List(filter, "C22001")
	if err != nil {
		t.Fatalf("expected list success for distributor prefix scope, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_List_DistributorSelfRowReturnedByDistributorIDWhenCustIDDiffers(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10, ParentCustId: "C22001", JwtDistributorId: 67}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude",
		"area_code", "area_name",
		"province_code", "province_name",
		"regency_code", "regency_name",
		"sub_district_code", "sub_district_name",
		"ward_code", "ward_name",
		"longitude", "region_code", "region_name", "channel_code", "channel_name",
		"dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name",
	}).AddRow(
		"C22001", int64(67), "3434", "Distributor iDetama", nil,
		1, 1, 1, 1, 1, "Alamat Test",
		"", "", "", "", "", 0, "-6.2",
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		"106.8", nil, nil, nil, nil,
		nil, nil, true, nil,
	)

	mock.ExpectQuery("(?s)SELECT .*FROM mst\\.m_distributor sp.*\\(sp\\.distributor_id = 67 OR sp\\.cust_id LIKE 'C220010001%'.*sp\\.is_active = true.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.List(filter, "C220010001")
	if err != nil {
		t.Fatalf("expected list success for distributor self row fallback, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}
	if data[0].CustId != "C22001" || data[0].DistributorId != 67 {
		t.Fatalf("expected self distributor row from distributor_id fallback, got cust_id=%s distributor_id=%d", data[0].CustId, data[0].DistributorId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_LookupList_UsesCustIDPrefixAndDefaultActiveTrue(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10, ParentCustId: "C22001", JwtDistributorId: 67}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "is_active",
		"longitude", "updated_by_name",
	}).AddRow(
		"C220010001", int64(67), "D067", "Distributor Self Child Cust", nil,
		1, 1, 1, 1, 1, "Alamat Test",
		"", "", "", "", "", 0, "-6.2", true,
		"106.8", nil,
	)

	mock.ExpectQuery("(?s)SELECT .*FROM mst\\.m_distributor sp.*\\(sp\\.distributor_id = 67 OR sp\\.cust_id LIKE 'C22001%'.*sp\\.is_active = true.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.LookupList(filter, "C22001")
	if err != nil {
		t.Fatalf("expected lookup list success for prefix scope, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}
	if data[0].DistributorId != 67 {
		t.Fatalf("expected distributor_id 67, got %d", data[0].DistributorId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_LookupList_AppliesPaginationAndZeroTotalReturnsZeroPageTotal(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	filter := entity.DistributorQueryFilter{Page: 2, Limit: 1, ParentCustId: "C22001"}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(0))

	mock.ExpectQuery("(?s)SELECT .*FROM mst\\.m_distributor sp.*sp\\.cust_id LIKE 'C22001%'.*sp\\.is_active = true.*LIMIT 1 OFFSET 1").
		WillReturnRows(sqlmock.NewRows([]string{
			"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
			"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
			"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "is_active",
			"longitude", "updated_by_name",
		}))

	data, total, lastPage, err := svc.LookupList(filter, "C22001")
	if err != nil {
		t.Fatalf("expected lookup list success with zero total, got error: %v", err)
	}
	if total != 0 || lastPage != 0 {
		t.Fatalf("expected pagination total=0 lastPage=0, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 0 {
		t.Fatalf("expected zero lookup rows, got %d", len(data))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_List_ExplicitInactiveFilterUsesCustIDPrefix(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	inactiveFilter := 2
	filter := entity.DistributorQueryFilter{Page: 1, Limit: 10, ParentCustId: "C22001", JwtDistributorId: 67, IsActive: &inactiveFilter}

	mock.ExpectQuery("SELECT\\s+COUNT\\(DISTINCT sp\\.distributor_id\\)\\s+AS total").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude",
		"area_code", "area_name",
		"province_code", "province_name",
		"regency_code", "regency_name",
		"sub_district_code", "sub_district_name",
		"ward_code", "ward_name",
		"longitude", "region_code", "region_name", "channel_code", "channel_name",
		"dist_price_grp_code", "dist_price_grp_name", "is_active", "updated_by_name",
	}).AddRow(
		"C220010001", int64(67), "D067", "Distributor Self Child Cust", nil,
		1, 1, 1, 1, 1, "Alamat Test",
		"", "", "", "", "", 0, "-6.2",
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		nil, nil,
		"106.8", nil, nil, nil, nil,
		nil, nil, false, nil,
	)

	mock.ExpectQuery("(?s)SELECT .*FROM mst\\.m_distributor sp.*sp\\.cust_id LIKE 'C22001%'.*sp\\.is_active = false.*LIMIT 10 OFFSET 0").
		WillReturnRows(rows)

	data, total, lastPage, err := svc.List(filter, "C22001")
	if err != nil {
		t.Fatalf("expected list success for explicit inactive filter, got error: %v", err)
	}
	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}
	if len(data) != 1 {
		t.Fatalf("expected one distributor row, got %d", len(data))
	}
	if data[0].CustId != "C220010001" || data[0].DistributorId != 67 {
		t.Fatalf("expected distributor self row for child cust_id, got cust_id=%s distributor_id=%d", data[0].CustId, data[0].DistributorId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Detail_ScanSupportsParentCustIDAndDistributorSelfScope(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	params := entity.DetailDistributorParams{
		CustId:           "C220010001",
		ParentCustId:     "C22001",
		JwtDistributorId: 67,
		DistributorId:    103,
	}

	detailRows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "longitude",
		"phone", "fax_number", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "updated_by_name", "is_active", "is_del",
		"created_by", "created_at", "updated_by", "updated_at", "deleted_by", "deleted_at", "region_code", "region_name",
		"channel_code", "channel_name", "sub_distributor_group_code", "sub_distributor_group_name", "dist_price_grp_code", "dist_price_grp_name",
		"allow_add_product", "allow_edit_product", "allow_manage_pricing", "allow_upload_secondary_sales",
	}).AddRow(
		"C220010001", "C22001", int64(103), "D103", "Distributor Detail", nil,
		1, 2, 3, 4, 5, "Alamat Detail",
		"31", "3171", "3171010", "3171010001", "12345", 10, "-6.2", "106.8",
		nil, nil, "AR01", "Area", "31", "DKI Jakarta", "3171", "Jakarta Selatan",
		"3171010", "Kebayoran Baru", "3171010001", "Senayan", "Updater", true, false,
		nil, nil, nil, nil, nil, nil, "RG01", "Region",
		"CH01", "Channel", "SDG01", "Sub Group", "DPG01", "Price Group",
		true, false, true, false,
	)

	mock.ExpectQuery(`(?s)SELECT .*mdist\.parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*\(mdist\.distributor_id = \$2 OR mdist\.cust_id LIKE \$3\).*`).
		WithArgs(103, int64(67), "C220010001%").
		WillReturnRows(detailRows)

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(103, "C220010001%").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_contact_id", "cust_id", "distributor_id", "contact_name", "job_title", "phone_no", "is_wa_no", "wa_no", "email", "identity_no", "identity_type"}))

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(103, "C220010001%").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_tax_id", "cust_id", "distributor_id", "tax_identifier_no_type", "tax_identifier_no", "nitku", "tax_name", "tax_address"}))

	resp, err := svc.Detail(params)
	if err != nil {
		t.Fatalf("expected detail success, got error: %v", err)
	}

	if resp.ParentCustId != "C22001" {
		t.Fatalf("expected parent_cust_id C22001, got %s", resp.ParentCustId)
	}

	if resp.DistributorId != 103 {
		t.Fatalf("expected distributor_id 103, got %d", resp.DistributorId)
	}

	if resp.DistributorSetup == nil {
		t.Fatalf("expected distributor_setup to be mapped")
	}

	if !resp.DistributorSetup.AllowAddProduct || resp.DistributorSetup.AllowEditProduct || !resp.DistributorSetup.AllowManagePricing || resp.DistributorSetup.AllowUploadSecondarySales {
		t.Fatalf("expected distributor_setup flags to match query result, got %+v", resp.DistributorSetup)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Detail_PrincipalScopeUsesCustIDPrefix(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	params := entity.DetailDistributorParams{
		CustId:        "C26002",
		ParentCustId:  "C26002",
		DistributorId: 103,
	}

	detailRows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "longitude",
		"phone", "fax_number", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "updated_by_name", "is_active", "is_del",
		"created_by", "created_at", "updated_by", "updated_at", "deleted_by", "deleted_at", "region_code", "region_name",
		"channel_code", "channel_name", "sub_distributor_group_code", "sub_distributor_group_name", "dist_price_grp_code", "dist_price_grp_name",
		"allow_add_product", "allow_edit_product", "allow_manage_pricing", "allow_upload_secondary_sales",
	}).AddRow(
		"C260020002", "C26002", int64(103), "787123", "CV. Abadi Nan Jaya", nil,
		1, 2, 3, 4, 5, "Alamat Detail",
		"31", "3171", "3171010", "3171010001", "12345", 10, "-6.2", "106.8",
		nil, nil, "AR01", "Area", "31", "DKI Jakarta", "3171", "Jakarta Selatan",
		"3171010", "Kebayoran Baru", "3171010001", "Senayan", "Updater", true, false,
		nil, nil, nil, nil, nil, nil, "RG01", "Region",
		"CH01", "Channel", "SDG01", "Sub Group", "DPG01", "Price Group",
		true, false, true, false,
	)

	mock.ExpectQuery(`(?s)SELECT .*mdist\.parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*mdist\.cust_id LIKE \$2.*`).
		WithArgs(103, "C26002%").
		WillReturnRows(detailRows)

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(103, "C26002%").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_contact_id", "cust_id", "distributor_id", "contact_name", "job_title", "phone_no", "is_wa_no", "wa_no", "email", "identity_no", "identity_type"}))

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(103, "C26002%").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_tax_id", "cust_id", "distributor_id", "tax_identifier_no_type", "tax_identifier_no", "nitku", "tax_name", "tax_address"}))

	resp, err := svc.Detail(params)
	if err != nil {
		t.Fatalf("expected detail success for principal scope, got error: %v", err)
	}

	if resp.ParentCustId != "C26002" {
		t.Fatalf("expected parent_cust_id C26002, got %s", resp.ParentCustId)
	}

	if resp.DistributorId != 103 {
		t.Fatalf("expected distributor_id 103, got %d", resp.DistributorId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Detail_LoadsContactAndTaxForChildCustID(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	params := entity.DetailDistributorParams{
		CustId:        "C26002",
		ParentCustId:  "C26002",
		DistributorId: 119,
	}

	detailRows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "longitude",
		"phone", "fax_number", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "updated_by_name", "is_active", "is_del",
		"created_by", "created_at", "updated_by", "updated_at", "deleted_by", "deleted_at", "region_code", "region_name",
		"channel_code", "channel_name", "sub_distributor_group_code", "sub_distributor_group_name", "dist_price_grp_code", "dist_price_grp_name",
		"allow_add_product", "allow_edit_product", "allow_manage_pricing", "allow_upload_secondary_sales",
	}).AddRow(
		"C260020003", "C26002", int64(119), "NEW119", "Distributor Child", nil,
		1, 2, 3, 4, 5, "Alamat Child",
		"31", "3171", "3171010", "3171010001", "12345", 10, "-6.2", "106.8",
		nil, nil, "AR01", "Area", "31", "DKI Jakarta", "3171", "Jakarta Selatan",
		"3171010", "Kebayoran Baru", "3171010001", "Senayan", "Updater", true, false,
		nil, nil, nil, nil, nil, nil, "RG01", "Region",
		"CH01", "Channel", "SDG01", "Sub Group", "DPG01", "Price Group",
		false, false, false, false,
	)

	mock.ExpectQuery(`(?s)SELECT .*mdist\.parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*mdist\.cust_id LIKE \$2.*`).
		WithArgs(119, "C26002%").
		WillReturnRows(detailRows)

	contactRows := sqlmock.NewRows([]string{"distributor_contact_id", "cust_id", "distributor_id", "contact_name", "job_title", "phone_no", "is_wa_no", "wa_no", "email", "identity_no", "identity_type"}).
		AddRow(int64(135), "C260020003", int64(119), "Nadia Putri", "PIC", "081260020119", true, "081260020119", "nadia.119@dummy.test", "317400000119", "National ID")

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(119, "C26002%").
		WillReturnRows(contactRows)

	taxRows := sqlmock.NewRows([]string{"distributor_tax_id", "cust_id", "distributor_id", "tax_identifier_no_type", "tax_identifier_no", "nitku", "tax_name", "tax_address"}).
		AddRow(int64(69), "C260020003", int64(119), "NPWP", "01.234.567.8-119.000", "NITKU-119", "New Distributor Tax", "Jl. Dummy Distributor 119 No. 3")

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(119, "C26002%").
		WillReturnRows(taxRows)

	resp, err := svc.Detail(params)
	if err != nil {
		t.Fatalf("expected detail success for child cust_id, got error: %v", err)
	}

	if len(resp.DistributorContact) != 1 {
		t.Fatalf("expected one distributor contact, got %d", len(resp.DistributorContact))
	}

	if resp.DistributorContact[0].ContactName != "Nadia Putri" {
		t.Fatalf("expected contact Nadia Putri, got %s", resp.DistributorContact[0].ContactName)
	}

	if len(resp.DistributorTax) != 1 {
		t.Fatalf("expected one distributor tax, got %d", len(resp.DistributorTax))
	}

	if resp.DistributorTax[0].TaxName != "New Distributor Tax" {
		t.Fatalf("expected tax name New Distributor Tax, got %s", resp.DistributorTax[0].TaxName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Detail_DistributorSelfScopeDoesNotSendUnusedQueryParameter(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	params := entity.DetailDistributorParams{
		CustId:           "C260020001",
		ParentCustId:     "C26002",
		JwtDistributorId: 102,
		DistributorId:    102,
	}

	detailRows := sqlmock.NewRows([]string{
		"cust_id", "parent_cust_id", "distributor_id", "distributor_code", "distributor_name", "barcode",
		"region_id", "area_id", "channel_id", "sub_distributor_group_id", "dist_price_grp_id", "address",
		"province_id", "regency_id", "sub_district_id", "ward_id", "zip_code", "ot_loc_id", "latitude", "longitude",
		"phone", "fax_number", "area_code", "area_name", "province_code", "province_name", "regency_code", "regency_name",
		"sub_district_code", "sub_district_name", "ward_code", "ward_name", "updated_by_name", "is_active", "is_del",
		"created_by", "created_at", "updated_by", "updated_at", "deleted_by", "deleted_at", "region_code", "region_name",
		"channel_code", "channel_name", "sub_distributor_group_code", "sub_distributor_group_name", "dist_price_grp_code", "dist_price_grp_name",
		"allow_add_product", "allow_edit_product", "allow_manage_pricing", "allow_upload_secondary_sales",
	}).AddRow(
		"C260020001", "C26002", int64(102), "162612", "PT Besi Makmur", nil,
		81, 88, 46, 37, 46, "jalan besi makmur madura",
		"1", "01", "10101", "1010101", "11", 0, "-6.2088", "106.8456",
		nil, nil, "101", "JAVA", "1", "JAWA BARAT", "01", "CIREBON",
		"10101", "HARJAMUKTI", "1010101", "LARANGAN", "Phill Jones", true, false,
		nil, nil, int64(140), nil, nil, nil, "1", "CENTRAL",
		"01", "NON", "01", "NON GROUP", "1", "NON",
		true, false, false, false,
	)

	mock.ExpectQuery(`(?s)SELECT .*mdist\.parent_cust_id.*FROM mst\.m_distributor mdist.*WHERE mdist\.distributor_id = \$1.*\(mdist\.distributor_id = \$2 OR mdist\.cust_id LIKE \$3\).*`).
		WithArgs(102, int64(102), "C260020001%").
		WillReturnRows(detailRows)

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(102, "C260020001%").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_contact_id", "cust_id", "distributor_id", "contact_name", "job_title", "phone_no", "is_wa_no", "wa_no", "email", "identity_no", "identity_type"}).AddRow(int64(131), "C260020001", int64(102), "Rina Pratama", "Manager", "081260020001", true, "081260020001", "rina.102@dummy.test", "317400000102", "National ID"))

	mock.ExpectQuery(`SELECT \*`).
		WithArgs(102, "C260020001%").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_tax_id", "cust_id", "distributor_id", "tax_identifier_no_type", "tax_identifier_no", "nitku", "tax_name", "tax_address"}).AddRow(int64(67), "C260020001", int64(102), "National ID", "0000000000000000", "000000", "PT Besi Makmur NPWP", "Jl. Dummy Distributor 102 No. 1"))

	resp, err := svc.Detail(params)
	if err != nil {
		t.Fatalf("expected detail success for distributor self scope, got error: %v", err)
	}

	if resp.DistributorId != 102 {
		t.Fatalf("expected distributor_id 102, got %d", resp.DistributorId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
