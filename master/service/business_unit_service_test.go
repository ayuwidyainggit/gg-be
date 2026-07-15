package service

import (
	"database/sql"
	"errors"
	"testing"

	"master/entity"
	"master/model"
)

type businessUnitRepositoryStub struct {
	userInfo            model.UserInfo
	userInfoErr         error
	principalData       []model.BusinessUnitDistributor
	principalTotal      int
	principalLastPage   int
	principalErr        error
	customerName        string
	customerErr         error
	findCustomerCustID  string
	findCustomerCalls   int
	distributorData     model.BusinessUnitDistributor
	distributorErr      error
	scopeEmployee       model.Employee
	scopeErr            error
	findPrincipalFilter entity.BusinessUnitQueryFilter
	findDistributorID   int
	findDistributorCust string
}

func (s *businessUnitRepositoryStub) FindUserByUsername(username string) (model.UserInfo, error) {
	if s.userInfoErr != nil {
		return model.UserInfo{}, s.userInfoErr
	}

	return s.userInfo, nil
}

func (s *businessUnitRepositoryStub) FindDistributorsByCustId(dataFilter entity.BusinessUnitQueryFilter) ([]model.BusinessUnitDistributor, int, int, error) {
	s.findPrincipalFilter = dataFilter
	return s.principalData, s.principalTotal, s.principalLastPage, s.principalErr
}

func (s *businessUnitRepositoryStub) FindCustomerNameByCustId(custId string) (string, error) {
	s.findCustomerCalls++
	s.findCustomerCustID = custId
	if s.customerErr != nil {
		return "", s.customerErr
	}
	return s.customerName, nil
}

func (s *businessUnitRepositoryStub) FindDistributorByDistributorId(distributorId int, custId string) (model.BusinessUnitDistributor, error) {
	s.findDistributorID = distributorId
	s.findDistributorCust = custId
	if s.distributorErr != nil {
		return model.BusinessUnitDistributor{}, s.distributorErr
	}

	return s.distributorData, nil
}

func (s *businessUnitRepositoryStub) FindEmployeeDropdownScope(empID int, custID string) (model.Employee, error) {
	if s.scopeErr != nil {
		return model.Employee{}, s.scopeErr
	}
	return s.scopeEmployee, nil
}

func TestBusinessUnitService_GetBusinessUnit_PrincipalIncludesCustID(t *testing.T) {
	areaID, regionID := 12, 1
	repo := &businessUnitRepositoryStub{
		userInfo:     model.UserInfo{UserId: 10, UserFullname: "Principal User"},
		customerName: "Principal Customer",
		principalData: []model.BusinessUnitDistributor{
			{
				CustId:          "C260020001",
				DistributorId:   23,
				DistributorCode: "DST001",
				DistributorName: "PT Sumber Makmur",
				AreaId:          &areaID,
				RegionId:        &regionID,
			},
		},
		principalTotal:    1,
		principalLastPage: 1,
	}

	svc := NewBusinessUnitService(repo, repo)
	filter := entity.BusinessUnitQueryFilter{CustId: "C26002", UserName: "principal@example.com", EmployeeId: 99}
	repo.scopeEmployee = model.Employee{RegionScope: "all", AreaScope: "all", DistributorScope: "all"}

	data, total, lastPage, err := svc.GetBusinessUnit(filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}

	resp, ok := data.(entity.BusinessUnitPrincipalResponse)
	if !ok {
		t.Fatalf("expected principal response type, got %T", data)
	}

	if resp.CustId != "C26002" {
		t.Fatalf("expected principal cust_id C26002, got %s", resp.CustId)
	}
	if resp.UserFullname != "Principal Customer" || resp.CustName != "Principal Customer" {
		t.Fatalf("expected principal customer display Principal Customer, got user_fullname=%s cust_name=%s", resp.UserFullname, resp.CustName)
	}

	if len(resp.DistributorData) != 1 {
		t.Fatalf("expected one distributor data, got %d", len(resp.DistributorData))
	}

	if resp.DistributorData[0].CustId != "C260020001" {
		t.Fatalf("expected distributor cust_id C260020001, got %s", resp.DistributorData[0].CustId)
	}

	if repo.findPrincipalFilter.CustId != "C26002" {
		t.Fatalf("expected repository filter cust_id C26002, got %s", repo.findPrincipalFilter.CustId)
	}
	if repo.findCustomerCustID != "C26002" || repo.findCustomerCalls != 1 {
		t.Fatalf("expected customer lookup once with C26002, got cust_id=%s calls=%d", repo.findCustomerCustID, repo.findCustomerCalls)
	}
	if repo.findPrincipalFilter.Scope.RegionScope != "all" || repo.findPrincipalFilter.Scope.AreaScope != "all" || repo.findPrincipalFilter.Scope.DistributorScope != "all" {
		t.Fatalf("expected normalized all scopes, got %+v", repo.findPrincipalFilter.Scope)
	}
}

func TestBusinessUnitService_GetBusinessUnit_PrincipalUsesCustomerNameForPrincessa(t *testing.T) {
	repo := &businessUnitRepositoryStub{
		userInfo:          model.UserInfo{UserId: 10, UserFullname: "Princessa Ahsani Taqwim"},
		customerName:      "PT. Madura Sejahtera",
		principalData:     []model.BusinessUnitDistributor{},
		principalTotal:    0,
		principalLastPage: 0,
		scopeEmployee:     model.Employee{RegionScope: "all", AreaScope: "all", DistributorScope: "all"},
	}
	svc := NewBusinessUnitService(repo, repo)

	data, _, _, err := svc.GetBusinessUnit(entity.BusinessUnitQueryFilter{CustId: "C26002", UserName: "princessa@example.com", EmployeeId: 99})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	resp := data.(entity.BusinessUnitPrincipalResponse)
	if resp.UserFullname != "PT. Madura Sejahtera" || resp.CustName != "PT. Madura Sejahtera" {
		t.Fatalf("expected customer name display, got user_fullname=%s cust_name=%s", resp.UserFullname, resp.CustName)
	}
	if repo.findCustomerCustID != "C26002" || repo.findCustomerCalls != 1 {
		t.Fatalf("expected customer lookup once with C26002, got cust_id=%s calls=%d", repo.findCustomerCustID, repo.findCustomerCalls)
	}
}

func TestBusinessUnitService_GetBusinessUnit_PrincipalUsesCustomerNameForAgung(t *testing.T) {
	repo := &businessUnitRepositoryStub{
		userInfo:          model.UserInfo{UserId: 11, UserFullname: "Agung Citra"},
		customerName:      "PT. Madura Sejahtera",
		principalData:     []model.BusinessUnitDistributor{},
		principalTotal:    0,
		principalLastPage: 0,
		scopeEmployee:     model.Employee{RegionScope: "all", AreaScope: "all", DistributorScope: "all"},
	}
	svc := NewBusinessUnitService(repo, repo)

	data, _, _, err := svc.GetBusinessUnit(entity.BusinessUnitQueryFilter{CustId: "C26002", UserName: "agung@example.com", EmployeeId: 99})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	resp := data.(entity.BusinessUnitPrincipalResponse)
	if resp.UserFullname != "PT. Madura Sejahtera" || resp.CustName != "PT. Madura Sejahtera" {
		t.Fatalf("expected customer name display, got user_fullname=%s cust_name=%s", resp.UserFullname, resp.CustName)
	}
	if repo.findCustomerCustID != "C26002" || repo.findCustomerCalls != 1 {
		t.Fatalf("expected customer lookup once with C26002, got cust_id=%s calls=%d", repo.findCustomerCustID, repo.findCustomerCalls)
	}
}

func TestBusinessUnitService_GetBusinessUnit_PrincipalMissingCustomerReturnsError(t *testing.T) {
	repo := &businessUnitRepositoryStub{
		userInfo:      model.UserInfo{UserId: 10, UserFullname: "Princessa Ahsani Taqwim"},
		customerErr:   sql.ErrNoRows,
		scopeEmployee: model.Employee{RegionScope: "all", AreaScope: "all", DistributorScope: "all"},
	}
	svc := NewBusinessUnitService(repo, repo)

	_, _, _, err := svc.GetBusinessUnit(entity.BusinessUnitQueryFilter{CustId: "C26002", UserName: "princessa@example.com", EmployeeId: 99})
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
	if repo.findCustomerCustID != "C26002" || repo.findCustomerCalls != 1 {
		t.Fatalf("expected customer lookup once with C26002, got cust_id=%s calls=%d", repo.findCustomerCustID, repo.findCustomerCalls)
	}
}

func TestBusinessUnitService_GetBusinessUnit_PrincipalRequiresEmployeeID(t *testing.T) {
	repo := &businessUnitRepositoryStub{userInfo: model.UserInfo{UserId: 10, UserFullname: "Principal User"}}
	svc := NewBusinessUnitService(repo, repo)

	_, _, _, err := svc.GetBusinessUnit(entity.BusinessUnitQueryFilter{CustId: "C26002", UserName: "principal@example.com"})
	if err == nil {
		t.Fatal("expected error for missing employee_id")
	}
}

func TestBusinessUnitService_GetBusinessUnit_DistributorIncludesCustID(t *testing.T) {
	areaID, regionID := 12, 1
	repo := &businessUnitRepositoryStub{
		userInfo: model.UserInfo{UserId: 12, UserFullname: "Distributor User"},
		scopeErr: errors.New("scope lookup should not run for distributor"),
		distributorData: model.BusinessUnitDistributor{
			CustId:          "C260020001",
			DistributorId:   23,
			DistributorCode: "DST001",
			DistributorName: "PT Sumber Makmur",
			AreaId:          &areaID,
			RegionId:        &regionID,
		},
	}

	svc := NewBusinessUnitService(repo, repo)
	distributorID := 23
	filter := entity.BusinessUnitQueryFilter{CustId: "C260020001", UserName: "dist@example.com", DistributorId: &distributorID}

	data, total, lastPage, err := svc.GetBusinessUnit(filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 1 || lastPage != 1 {
		t.Fatalf("expected pagination total=1 lastPage=1, got total=%d lastPage=%d", total, lastPage)
	}

	resp, ok := data.(entity.BusinessUnitDistributorResponse)
	if !ok {
		t.Fatalf("expected distributor response type, got %T", data)
	}

	if resp.CustId != "C260020001" {
		t.Fatalf("expected distributor cust_id C260020001, got %s", resp.CustId)
	}

	if repo.findDistributorID != 23 || repo.findDistributorCust != "C260020001" {
		t.Fatalf("expected repository call with distributor_id=23 cust_id=C260020001, got distributor_id=%d cust_id=%s", repo.findDistributorID, repo.findDistributorCust)
	}
	if repo.findCustomerCalls != 0 {
		t.Fatalf("expected distributor path not to call customer lookup, got %d calls", repo.findCustomerCalls)
	}
}

func TestBusinessUnitService_GetBusinessUnit_ReturnsUserLookupError(t *testing.T) {
	repo := &businessUnitRepositoryStub{userInfoErr: errors.New("lookup failed")}
	svc := NewBusinessUnitService(repo, repo)

	_, _, _, err := svc.GetBusinessUnit(entity.BusinessUnitQueryFilter{UserName: "principal@example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
