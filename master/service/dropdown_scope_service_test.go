package service

import (
	"errors"
	"testing"

	"master/entity"
	"master/model"
)

type regionRepositoryScopeStub struct {
	rows      []model.Region
	total     int
	lastPage  int
	err       error
	lastInput entity.RegionQueryFilter
	scopeRow  model.Employee
	scopeErr  error
}

func (s *regionRepositoryScopeStub) FindOneByRegionIdAndCustId(regionId int, custId string) (model.Region, error) {
	return model.Region{}, nil
}
func (s *regionRepositoryScopeStub) FindOneByRegionCodeAndCustId(regionCode string, custId string) (model.Region, error) {
	return model.Region{}, nil
}
func (s *regionRepositoryScopeStub) FindAllByCustId(filter entity.RegionQueryFilter) ([]model.Region, int, int, error) {
	s.lastInput = filter
	return s.rows, s.total, s.lastPage, s.err
}
func (s *regionRepositoryScopeStub) FindAllByCustIdLookupMode(filter entity.RegionQueryFilter) ([]model.Region, int, int, error) {
	s.lastInput = filter
	return s.rows, s.total, s.lastPage, s.err
}
func (s *regionRepositoryScopeStub) Store(region model.Region) (int, error) { return 0, nil }
func (s *regionRepositoryScopeStub) Update(regionId int, request entity.UpdateRegionRequest) error { return nil }
func (s *regionRepositoryScopeStub) Delete(custId string, regionId int, deletedBy int64) error { return nil }
func (s *regionRepositoryScopeStub) FindEmployeeDropdownScope(empID int, custID string) (model.Employee, error) {
	if s.scopeErr != nil { return model.Employee{}, s.scopeErr }
	return s.scopeRow, nil
}

type areaRepositoryScopeStub struct {
	rows      []model.AreaList
	total     int
	lastPage  int
	err       error
	lastInput entity.AreaQueryFilter
	scopeRow  model.Employee
	scopeErr  error
}

func (s *areaRepositoryScopeStub) FindOneByAreaIdAndCustId(areaId int, custId string) (model.AreaList, error) {
	return model.AreaList{}, nil
}
func (s *areaRepositoryScopeStub) FindOneByAreaCodeAndCustId(areaCode string, custId string) (model.Area, error) {
	return model.Area{}, nil
}
func (s *areaRepositoryScopeStub) FindAllByCustId(filter entity.AreaQueryFilter) ([]model.AreaList, int, int, error) {
	s.lastInput = filter
	return s.rows, s.total, s.lastPage, s.err
}
func (s *areaRepositoryScopeStub) FindAllByCustIdLookupMode(filter entity.AreaQueryFilter) ([]model.AreaList, int, int, error) {
	s.lastInput = filter
	return s.rows, s.total, s.lastPage, s.err
}
func (s *areaRepositoryScopeStub) Store(area model.Area) (int, error) { return 0, nil }
func (s *areaRepositoryScopeStub) Update(areaId int, request entity.UpdateAreaRequest) error { return nil }
func (s *areaRepositoryScopeStub) Delete(custId string, areaId int, deletedBy int64) error { return nil }
func (s *areaRepositoryScopeStub) FindEmployeeDropdownScope(empID int, custID string) (model.Employee, error) {
	if s.scopeErr != nil { return model.Employee{}, s.scopeErr }
	return s.scopeRow, nil
}

func TestRegionService_List_PrincipalLoadsScope(t *testing.T) {
	repo := &regionRepositoryScopeStub{
		rows:     []model.Region{{RegionId: 1, RegionCode: "R1", RegionName: "Region 1"}},
		total:    1,
		lastPage: 1,
		scopeRow: model.Employee{RegionScope: "SELECTED", AreaScope: "ALL", DistributorScope: "ALL"},
	}
	svc := NewRegionService(repo, repo)

	_, _, _, err := svc.List(entity.RegionQueryFilter{CustId: "C22001", DistributorId: 0, EmployeeId: 77})
	if err != nil { t.Fatalf("expected no error, got %v", err) }
	if repo.lastInput.Scope.RegionScope != "specific" {
		t.Fatalf("expected normalized specific scope, got %+v", repo.lastInput.Scope)
	}
}

func TestRegionService_List_NonPrincipalSkipsScopeLookup(t *testing.T) {
	repo := &regionRepositoryScopeStub{rows: []model.Region{}, total: 0, lastPage: 0, scopeErr: errors.New("should not be called")}
	svc := NewRegionService(repo, repo)

	distributorID := 99
	_, _, _, err := svc.List(entity.RegionQueryFilter{CustId: "C220010001", DistributorId: distributorID})
	if err != nil { t.Fatalf("expected no error, got %v", err) }
}

func TestAreaService_List_PrincipalLoadsScope(t *testing.T) {
	repo := &areaRepositoryScopeStub{
		rows:     []model.AreaList{{AreaId: 10, AreaCode: "A1", AreaName: "Area 1"}},
		total:    1,
		lastPage: 1,
		scopeRow: model.Employee{RegionScope: "SPESIFIC", AreaScope: "ALL", DistributorScope: "ALL"},
	}
	svc := NewAreaService(repo, repo)

	_, _, _, err := svc.List(entity.AreaQueryFilter{CustId: "C22001", DistributorId: 0, EmployeeId: 77})
	if err != nil { t.Fatalf("expected no error, got %v", err) }
	if repo.lastInput.Scope.RegionScope != "specific" {
		t.Fatalf("expected normalized specific region scope, got %+v", repo.lastInput.Scope)
	}
}

func TestAreaService_List_PrincipalRequiresEmployeeID(t *testing.T) {
	repo := &areaRepositoryScopeStub{}
	svc := NewAreaService(repo, repo)

	_, _, _, err := svc.List(entity.AreaQueryFilter{CustId: "C22001", DistributorId: 0})
	if err == nil {
		t.Fatal("expected error for missing employee_id")
	}
}

func TestAreaService_List_NonPrincipalSkipsScopeLookup(t *testing.T) {
	repo := &areaRepositoryScopeStub{rows: []model.AreaList{}, total: 0, lastPage: 0, scopeErr: errors.New("should not be called")}
	svc := NewAreaService(repo, repo)

	_, _, _, err := svc.List(entity.AreaQueryFilter{CustId: "C220010001", DistributorId: 99})
	if err != nil { t.Fatalf("expected no error, got %v", err) }
}
