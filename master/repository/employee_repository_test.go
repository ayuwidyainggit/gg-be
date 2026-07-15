package repository

import (
	"database/sql/driver"
	"master/entity"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"master/model"
)

func setupEmployeeRepositoryTest(t *testing.T) (*EmployeeRepositoryImpl, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &EmployeeRepositoryImpl{DB: sqlxDB}, mock, func() { _ = db.Close() }
}

func anyEmployeeStoreArgs() []driver.Value {
	args := make([]driver.Value, 32)
	for i := range args {
		args[i] = sqlmock.AnyArg()
	}
	return args
}

func TestEmployeeRepositoryStorePersistsTerritoryMappings(t *testing.T) {
	repo, mock, cleanup := setupEmployeeRepositoryTest(t)
	defer cleanup()

	now := time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC)
	address := "Jl. Merdeka"
	empType := "S"
	empGrp := 7
	provinceID := "31"
	cityID := "3174"
	subDistrictID := "317401"
	wardID := "31740101"
	identityNo := "317400000001"
	postCode := "12345"
	divisionID := 3
	userID := int64(10)

	employee := model.Employee{
		CustId:           "C22001",
		EmployeeCode:     "EMP001",
		EmployeeName:     "Employee One",
		Address:          &address,
		EmpTypeId:        &empType,
		EmpGrpId:         &empGrp,
		WorkDate:         &now,
		LastEducation:    nil,
		Dob:              &now,
		IsActive:         true,
		CreatedBy:        &userID,
		CreatedAt:        &now,
		UpdatedBy:        &userID,
		UpdatedAt:        &now,
		ProvinceId:       &provinceID,
		CityId:           &cityID,
		SubDistrictId:    &subDistrictID,
		WardId:           &wardID,
		IdentityNo:       &identityNo,
		PostCode:         &postCode,
		DivisionId:       &divisionID,
		RegionScope:      "SELECTED",
		AreaScope:        "SELECTED",
		DistributorScope: "SELECTED",
		RegionIds:        []int{1},
		AreaIds:          []int{10},
		DistributorIds:   []int{100},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO mst\.m_employee`).
		WithArgs(anyEmployeeStoreArgs()...).
		WillReturnRows(sqlmock.NewRows([]string{"emp_id"}).AddRow(99))
	mock.ExpectExec(`INSERT INTO mst\.m_employee_region_mapping`).
		WithArgs("C22001", 99, 1, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO mst\.m_employee_area_mapping`).
		WithArgs("C22001", 99, 10, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO mst\.m_employee_distributor_mapping`).
		WithArgs("C22001", 99, 100, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	employeeID, err := repo.Store(employee)
	if err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if employeeID != 99 {
		t.Fatalf("expected returned employee id 99, got %d", employeeID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestEmployeeRepositoryValidateEmployeeTerritoryMappingRejectsAreaOutsideRegion(t *testing.T) {
	repo, mock, cleanup := setupEmployeeRepositoryTest(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT COUNT\(DISTINCT region_id\)`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT area_id\)`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

	err := repo.ValidateEmployeeTerritoryMapping("C22001", model.EmployeeTerritoryMapping{
		RegionIds: []int{1},
		AreaIds:   []int{10, 11},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "area_ids must belong to selected region_ids") {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestEmployeeRepositoryFindEmployeeTerritoryDetail(t *testing.T) {
	repo, mock, cleanup := setupEmployeeRepositoryTest(t)
	defer cleanup()

	params := entity.DetailEmployeeParams{
		CustId:       "C22001001",
		ParentCustId: "C22001",
		EmployeeId:   99,
	}

	mock.ExpectQuery(`FROM mst\.m_employee_region_mapping`).
		WithArgs(params.CustId, params.EmployeeId, params.ParentCustId).
		WillReturnRows(sqlmock.NewRows([]string{"region_id", "region_code", "region_name"}).
			AddRow(1, "R001", "Region One"))
	mock.ExpectQuery(`FROM mst\.m_employee_area_mapping`).
		WithArgs(params.CustId, params.EmployeeId, params.ParentCustId).
		WillReturnRows(sqlmock.NewRows([]string{"area_id", "area_code", "area_name", "region_id", "region_code", "region_name"}).
			AddRow(10, "A001", "Area One", 1, "R001", "Region One"))
	mock.ExpectQuery(`FROM mst\.m_employee_distributor_mapping`).
		WithArgs(params.CustId, params.EmployeeId, params.ParentCustId, params.ParentCustId+"%").
		WillReturnRows(sqlmock.NewRows([]string{"distributor_id", "distributor_code", "distributor_name", "area_id", "area_code", "area_name", "region_id", "region_code", "region_name"}).
			AddRow(100, "D001", "Distributor One", 10, "A001", "Area One", 1, "R001", "Region One"))

	detail, err := repo.FindEmployeeTerritoryDetail(params)
	if err != nil {
		t.Fatalf("FindEmployeeTerritoryDetail returned error: %v", err)
	}
	if len(detail.Regions) != 1 || detail.Regions[0].RegionId != 1 {
		t.Fatalf("unexpected regions: %+v", detail.Regions)
	}
	if len(detail.Areas) != 1 || detail.Areas[0].AreaId != 10 || detail.Areas[0].RegionId != 1 {
		t.Fatalf("unexpected areas: %+v", detail.Areas)
	}
	if len(detail.Distributors) != 1 || detail.Distributors[0].DistributorId != 100 || detail.Distributors[0].AreaId != 10 {
		t.Fatalf("unexpected distributors: %+v", detail.Distributors)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestEmployeeRepositoryUpdateReplacesTerritoryMappings(t *testing.T) {
	repo, mock, cleanup := setupEmployeeRepositoryTest(t)
	defer cleanup()

	request := entity.UpdateEmployeeRequest{
		CustId:                   "C22001",
		UpdatedBy:                10,
		EmployeeCode:             "EMP001",
		RegionScope:              "SELECTED",
		AreaScope:                "SELECTED",
		DistributorScope:         "SELECTED",
		RegionIds:                []int{1},
		AreaIds:                  []int{10},
		DistributorIds:           []int{100},
		TerritoryMappingProvided: true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE mst\.m_employee`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE mst\.m_employee_region_mapping`).
		WithArgs(request.CustId, 99, request.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE mst\.m_employee_area_mapping`).
		WithArgs(request.CustId, 99, request.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`UPDATE mst\.m_employee_distributor_mapping`).
		WithArgs(request.CustId, 99, request.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO mst\.m_employee_region_mapping`).
		WithArgs(request.CustId, 99, 1, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO mst\.m_employee_area_mapping`).
		WithArgs(request.CustId, 99, 10, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO mst\.m_employee_distributor_mapping`).
		WithArgs(request.CustId, 99, 100, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.Update(99, request); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
