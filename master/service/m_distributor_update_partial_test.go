package service

import (
	"database/sql"
	"testing"

	"master/entity"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestDistributorService_Update_WithoutLocationFields_Success(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:          "CUST001",
		UpdatedBy:       serviceInt64Ptr(1002),
		DistributorName: serviceStrPtr("Distributor Partial Update"),
	}
	expectResolvedDistributorCustID(mock, 103, "CUST001", "CUST001")

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(103, request)
	if err != nil {
		t.Fatalf("expected update success without location fields, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_WithPartialLocationField_Success(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:                "CUST001",
		UpdatedBy:             serviceInt64Ptr(1002),
		ProvinceId:            serviceStrPtr("11"),
		RegionId:              serviceIntPtr(1),
		Address:               serviceStrPtr("Jl Updated"),
		IsActive:              serviceBoolPtr(true),
		Latitude:              serviceStrPtr("-6.2100"),
		Longitude:             serviceStrPtr("106.8266"),
		ChannelId:             serviceIntPtr(3),
		AreaId:                serviceIntPtr(2),
		DistPriceGrpId:        serviceIntPtr(5),
		SubDistributorGroupId: serviceIntPtr(4),
	}
	expectResolvedDistributorCustID(mock, 103, "CUST001", "CUST001")

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(103, request)
	if err != nil {
		t.Fatalf("expected update success with partial location fields, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_WithAllLocationFields_Success(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:                "CUST001",
		UpdatedBy:             serviceInt64Ptr(1002),
		DistributorCode:       serviceStrPtr("DIST001"),
		DistributorName:       serviceStrPtr("Distributor Full Location Update"),
		RegionId:              serviceIntPtr(1),
		AreaId:                serviceIntPtr(2),
		ChannelId:             serviceIntPtr(3),
		SubDistributorGroupId: serviceIntPtr(4),
		DistPriceGrpId:        serviceIntPtr(5),
		Address:               serviceStrPtr("Jl Updated"),
		ProvinceId:            serviceStrPtr("11"),
		RegencyId:             serviceStrPtr("1101"),
		SubDistrictId:         serviceStrPtr("1101010"),
		WardId:                serviceStrPtr("1101010001"),
		Latitude:              serviceStrPtr("-6.2100"),
		Longitude:             serviceStrPtr("106.8266"),
		IsActive:              serviceBoolPtr(true),
	}
	expectResolvedDistributorCustID(mock, 103, "CUST001", "CUST001")

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST001", "CUST001").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(103, request)
	if err != nil {
		t.Fatalf("expected update success with all location fields, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDistributorService_Update_AcceptsAlphanumericCode(t *testing.T) {
	svc, mock, cleanup := setupDistributorServiceTest(t)
	defer cleanup()

	request := entity.UpdateDistributorRequest{
		CustId:          "CUST001",
		UpdatedBy:       serviceInt64Ptr(1002),
		DistributorCode: serviceStrPtr("DIST-15676761A"),
		DistributorName: serviceStrPtr("Distributor Alpha"),
	}
	expectResolvedDistributorCustID(mock, 104, "CUST001", "CUST001")

	mock.ExpectQuery("SELECT sp\\.distributor_id, sp\\.distributor_code").
		WithArgs("DIST-15676761A", "CUST001").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE mst\\.m_distributor").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := svc.Update(104, request)
	if err != nil {
		t.Fatalf("expected alphanumeric code accepted, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
