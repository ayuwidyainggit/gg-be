package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryHierarchyApprovalImpl struct {
		*gorm.DB
	}
)
type HierarchyApprovalRepository interface {
	FindCompanies(dataFilter entity.CompaniesQueryFilter) ([]model.SmcMCustomer, int64, int, error)
	Store(c context.Context, data *model.HierarchyApproval) error
	StoreDetail(c context.Context, data *model.HierarchyApprovalDet) error
	StoreDetailEmp(c context.Context, data []*model.HierarchyApprovalDetEmp) error
	FindOneByEmployeeIdAndCustId(employeeID int64, custID string) (model.HierarcyApprovalEmployee, error)
	FindAllByCustID(dataFilter entity.HierarcyApprovalQueryFilter) ([]model.HierarchyApprovalList, int64, int, error)
	FindByCustID(hierarchyApprovalID int64, custId string, parentCustID string) (details model.HierarchyApprovalRead, err error)
	FindDetail(hierarchyApprovalID int64) (details []model.HierarchyApprovalDetRead, err error)
	FindDetailEmp(hierarchyApprovalDetID int64) (details []model.HierarchyApprovalDetEmpRead, err error)
	Delete(c context.Context, custId string, parentCustID string, hierarchyApprovalID int64, deletedBy int64) error
	Update(c context.Context, hierarchyApprovalID int64, data model.HierarchyApproval) error
	UpdateDetail(c context.Context, hierarchyApprovalDetID int64, data model.HierarchyApprovalUpdate) error
	UpdateDetailEmp(c context.Context, hierarchyApprovalDetEmpID int64, data model.HierarchyApprovalDetEmpUpdate) error
	DeleteDetailEmp(c context.Context, hierarchyApprovalDetID int64, seq int) error
	GetUser(custID string) (model.SmcMCustomer, error)
	FindAllEmployeeByCustIdLookupMode(dataFilter entity.EmployeeHierarchyQueryFilter) ([]model.Employee, int64, int, error)
	FindBySetupFor(setupFor string) (details model.HierarchyApprovalRead, err error)
	FindBySetupForOnly(setupFor string) (details model.HierarchyApprovalRead, err error)
}

func NewHierarchyApprovalRepo(db *gorm.DB) *RepositoryHierarchyApprovalImpl {
	return &RepositoryHierarchyApprovalImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryHierarchyApprovalImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryHierarchyApprovalImpl) GetUser(custID string) (model.SmcMCustomer, error) {
	var cust model.SmcMCustomer
	err := repository.Select("smc.m_customer.*, CASE WHEN smc.m_customer.cust_id = smc.m_customer.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END AS company_code").
		Joins("left join mst.m_distributor dist on smc.m_customer.distributor_id = dist.distributor_id").
		Joins("left join mst.m_supplier supp on smc.m_customer.distributor_id = supp.sup_id").
		Where("smc.m_customer.cust_id = ?", custID).Take(&cust).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return cust, errors.New(fmt.Sprintf("cust id : %v not found", custID))
		}

		return cust, err
	}
	return cust, nil
}

func (repository *RepositoryHierarchyApprovalImpl) FindCompanies(dataFilter entity.CompaniesQueryFilter) ([]model.SmcMCustomer, int64, int, error) {
	var customers []model.SmcMCustomer
	var total int64

	queryCount := repository.Select("smc.m_customer.cust_id").
		Joins("left join mst.m_distributor dist on smc.m_customer.distributor_id = dist.distributor_id").
		Joins("left join mst.m_supplier supp on smc.m_customer.distributor_id = supp.sup_id")

	query := repository.Select("smc.m_customer.*, CASE WHEN smc.m_customer.cust_id = smc.m_customer.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END AS company_code, CASE WHEN smc.m_customer.cust_id = smc.m_customer.parent_cust_id THEN supp.sup_name ELSE dist.distributor_name END AS company_name ").
		Joins("left join mst.m_distributor dist on smc.m_customer.distributor_id = dist.distributor_id").
		Joins("left join mst.m_supplier supp on smc.m_customer.distributor_id = supp.sup_id")

	queryCount.Where("smc.m_customer.parent_cust_id = ?", dataFilter.CustId)
	query.Where("smc.m_customer.parent_cust_id=?", dataFilter.CustId)

	if len(dataFilter.CompanyIds) > 0 {
		queryCount.Where("smc.m_customer.cust_id in ?", dataFilter.CompanyIds)
		query.Where("smc.m_customer.cust_id in ?", dataFilter.CompanyIds)
	}

	if dataFilter.Query != "" {
		queryCount.Where("CASE WHEN smc.m_customer.cust_id = smc.m_customer.parent_cust_id THEN lower(supp.sup_name) ELSE lower(dist.distributor_name) END LIKE ? OR CASE WHEN smc.m_customer.cust_id = smc.m_customer.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END LIKE ?", "%"+strings.ToLower(dataFilter.Query)+"%", "%"+dataFilter.Query+"%")
		query.Where("CASE WHEN smc.m_customer.cust_id = smc.m_customer.parent_cust_id THEN lower(supp.sup_name) ELSE lower(dist.distributor_name) END LIKE ? OR CASE WHEN smc.m_customer.cust_id = smc.m_customer.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END LIKE ?", "%"+strings.ToLower(dataFilter.Query)+"%", "%"+dataFilter.Query+"%")
	}
	err := query.Find(&customers).Error
	if err != nil {
		return customers, total, 0, err
	}
	err = queryCount.Model(&customers).Count(&total).Error
	if err != nil {
		return customers, total, 0, err
	}

	lastPage := 1

	return customers, total, lastPage, nil
}

func (repository *RepositoryHierarchyApprovalImpl) Store(c context.Context, data *model.HierarchyApproval) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) StoreDetail(c context.Context, data *model.HierarchyApprovalDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) StoreDetailEmp(c context.Context, data []*model.HierarchyApprovalDetEmp) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) FindOneByEmployeeIdAndCustId(employeeID int64, custID string) (model.HierarcyApprovalEmployee, error) {
	employee := model.HierarcyApprovalEmployee{}
	err := repository.Select("*").Where("emp_id = ? AND cust_id = ?", employeeID, custID).Take(&employee).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return employee, errors.New(fmt.Sprintf("employee id : %v not found on cust_id %v", employeeID, custID))
		}

		return employee, err
	}
	return employee, nil
}

func (repository *RepositoryHierarchyApprovalImpl) FindAllByCustID(dataFilter entity.HierarcyApprovalQueryFilter) ([]model.HierarchyApprovalList, int64, int, error) {
	var hieararchyApprovals []model.HierarchyApprovalList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}
	queryCount := repository.Select("hierarchy_approval_id").
		Joins("LEFT JOIN smc.m_customer mu ON sls.hierarchy_approvals.setup_for = mu.cust_id")

	query := repository.Select(`sls.hierarchy_approvals.*, mu.cust_name AS setup_for_name, us.user_fullname AS updated_by_name, mu.*, CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END AS company_code, CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_name ELSE dist.distributor_name END AS company_name `).
		Joins("left join sys.m_user us on us.user_id = sls.hierarchy_approvals.updated_by").
		Joins("LEFT JOIN smc.m_customer mu ON sls.hierarchy_approvals.setup_for = mu.cust_id").
		Joins("left join mst.m_distributor dist on mu.distributor_id = dist.distributor_id").
		Joins("left join mst.m_supplier supp on mu.distributor_id = supp.sup_id")

	if dataFilter.CustId != dataFilter.ParentCustId {
		queryCount.Where("sls.hierarchy_approvals.setup_for = ?", dataFilter.CustId)
		query.Where("sls.hierarchy_approvals.setup_for = ?", dataFilter.CustId)
	} else {
		queryCount.Where("sls.hierarchy_approvals.setup_for in (select cust_id from smc.m_customer where parent_cust_id = ?)", dataFilter.ParentCustId)
		query.Where("sls.hierarchy_approvals.setup_for in (select cust_id from smc.m_customer where parent_cust_id = ?)", dataFilter.ParentCustId)
	}

	if dataFilter.Company != "" {
		queryCount.Where("mu.cust_id = ?", dataFilter.Company)
		query.Where("mu.cust_id = ?", dataFilter.Company)
	}

	if dataFilter.ApprovalType != nil {
		queryCount.Where("sls.hierarchy_approvals.hierarchy_approval_type = ?", dataFilter.ApprovalType)
		query.Where("sls.hierarchy_approvals.hierarchy_approval_type = ?", dataFilter.ApprovalType)
	}

	err := query.Find(&hieararchyApprovals).Error
	if err != nil {
		return hieararchyApprovals, total, 0, err
	}
	err = queryCount.Model(&hieararchyApprovals).Count(&total).Error
	if err != nil {
		return hieararchyApprovals, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return hieararchyApprovals, total, lastPage, nil
}

func (repository *RepositoryHierarchyApprovalImpl) FindByCustID(hierarchyApprovalID int64, custId string, parentCustID string) (details model.HierarchyApprovalRead, err error) {
	query := repository.Select(`sls.hierarchy_approvals.*, mu.cust_name AS setup_for_name,  CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END AS company_code, CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_name ELSE dist.distributor_name END AS company_name `).
		Joins("LEFT JOIN smc.m_customer mu ON sls.hierarchy_approvals.setup_for = mu.cust_id").
		Joins("left join mst.m_distributor dist on mu.distributor_id = dist.distributor_id").
		Joins("left join mst.m_supplier supp on mu.distributor_id = supp.sup_id").
		Where("sls.hierarchy_approvals.hierarchy_approval_id = ?", hierarchyApprovalID)

	if custId != parentCustID {
		query.Where("sls.hierarchy_approvals.setup_for = ?", custId)
	} else {
		query.Where("sls.hierarchy_approvals.setup_for in (select cust_id from smc.m_customer where parent_cust_id = ?)", parentCustID)
	}

	query.Take(&details)
	if query.Error != nil {
		return details, query.Error
	}

	return details, err
}

func (repository *RepositoryHierarchyApprovalImpl) FindBySetupFor(setupFor string) (details model.HierarchyApprovalRead, err error) {
	err = repository.Select(`sls.hierarchy_approvals.*, mu.cust_name AS setup_for_name,  CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END AS company_code, CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_name ELSE dist.distributor_name END AS company_name `).
		Joins("LEFT JOIN smc.m_customer mu ON sls.hierarchy_approvals.setup_for = mu.cust_id").
		Joins("left join mst.m_distributor dist on mu.distributor_id = dist.distributor_id").
		Joins("left join mst.m_supplier supp on mu.distributor_id = supp.sup_id").
		Where("sls.hierarchy_approvals.setup_for=?", setupFor).
		Take(&details).Error
	return details, err
}

func (repository *RepositoryHierarchyApprovalImpl) FindBySetupForOnly(setupFor string) (details model.HierarchyApprovalRead, err error) {
	err = repository.Select(`sls.hierarchy_approvals.*, mu.cust_name AS setup_for_name,  CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_code ELSE dist.distributor_code END AS company_code, CASE WHEN mu.cust_id = mu.parent_cust_id THEN supp.sup_name ELSE dist.distributor_name END AS company_name `).
		Joins("LEFT JOIN smc.m_customer mu ON sls.hierarchy_approvals.setup_for = mu.cust_id").
		Joins("left join mst.m_distributor dist on mu.distributor_id = dist.distributor_id").
		Joins("left join mst.m_supplier supp on mu.distributor_id = supp.sup_id").
		Where("sls.hierarchy_approvals.setup_for=?", setupFor).
		Take(&details).Error
	return details, err
}

func (repository *RepositoryHierarchyApprovalImpl) FindDetail(hierarchyApprovalID int64) (details []model.HierarchyApprovalDetRead, err error) {
	err = repository.Select(`sls.hierarchy_approvals_details.*, mu.cust_name AS setup_for_name`).
		Joins("LEFT JOIN smc.m_customer mu ON sls.hierarchy_approvals_details.hierarchy_approval_detail_cust_id = mu.cust_id").
		Where("sls.hierarchy_approvals_details.hierarchy_approval_id=?", hierarchyApprovalID).
		Order("level ASC").
		Find(&details).Error
	return details, err
}

func (repository *RepositoryHierarchyApprovalImpl) FindDetailEmp(hierarchyApprovalDetID int64) (details []model.HierarchyApprovalDetEmpRead, err error) {
	err = repository.Select(`sls.hierarchy_approvals_details_emp.*, emp.emp_name`).
		Joins("LEFT JOIN mst.m_employee emp ON sls.hierarchy_approvals_details_emp.emp_id = emp.emp_id").
		Where("sls.hierarchy_approvals_details_emp.hierarchy_approval_detail_id=?", hierarchyApprovalDetID).
		Order("seq ASC").
		Find(&details).Error
	return details, err
}

func (repository *RepositoryHierarchyApprovalImpl) Delete(c context.Context, custId string, parentCustID string, hierarchyApprovalID int64, deletedBy int64) error {
	var data model.HierarchyApproval
	result := repository.model(c).Model(&data)

	if custId != parentCustID {
		result.Where("hierarchy_approval_id=? AND setup_for = ?", hierarchyApprovalID, custId)
	} else {
		result.Where("hierarchy_approval_id=? AND setup_for in (select cust_id from smc.m_customer where parent_cust_id = ?)", hierarchyApprovalID, parentCustID)
	}

	result.Updates(map[string]interface{}{"deleted_by": deletedBy, "deleted_at": time.Now()})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) Update(c context.Context, hierarchyApprovalID int64, data model.HierarchyApproval) error {
	result := repository.model(c).Model(&data).Where("hierarchy_approval_id = ?", hierarchyApprovalID).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) UpdateDetail(c context.Context, hierarchyApprovalDetID int64, data model.HierarchyApprovalUpdate) error {
	result := repository.model(c).Model(&data).Where("hierarchy_approval_detail_id = ? ", hierarchyApprovalDetID).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) UpdateDetailEmp(c context.Context, hierarchyApprovalDetEmpID int64, data model.HierarchyApprovalDetEmpUpdate) error {
	result := repository.model(c).Model(&data).Where("hierarchy_approval_detail_emp_id = ? ", hierarchyApprovalDetEmpID).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) DeleteDetailEmp(c context.Context, hierarchyApprovalDetID int64, seq int) error {
	var data model.HierarchyApprovalDetEmp
	result := repository.model(c).Where("hierarchy_approval_detail_id=? AND seq = ?", hierarchyApprovalDetID, seq).Delete(&data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryHierarchyApprovalImpl) FindAllEmployeeByCustIdLookupMode(dataFilter entity.EmployeeHierarchyQueryFilter) ([]model.Employee, int64, int, error) {

	var salesmans []model.Employee

	var total int64

	queryCount := repository.Select("mst.m_employee.emp_id").
		Joins("join sys.m_user mu on mst.m_employee.emp_id = mu.emp_id AND mu.cust_id = ?", dataFilter.CustId)

	query := repository.Select(`mst.m_employee.cust_id,mst.m_employee.emp_id, mst.m_employee.emp_code, mst.m_employee.emp_name, mst.m_employee.emp_grp_id`).
		Joins("join sys.m_user mu on mst.m_employee.emp_id = mu.emp_id AND mu.cust_id = ?", dataFilter.CustId)

	// queryCount.Where("sls.return.cust_id=?", dataFilter.CustId)
	// query.Where("sls.return.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_employee.is_active=?", true)
	query.Where("mst.m_employee.is_active=?", true)

	queryCount.Where("mst.m_employee.is_del=?", false)
	query.Where("mst.m_employee.is_del=?", false)

	if dataFilter.CustId != "" {
		queryCount.Where("mst.m_employee.cust_id=?", dataFilter.CustId)
		query.Where("mst.m_employee.cust_id=?", dataFilter.CustId)
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_employee.emp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, "mst.m_employee."+colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("mst.m_employee.emp_id DESC")
	}

	err := query.Find(&salesmans).Error
	if err != nil {
		return salesmans, total, 0, err
	}

	total = int64(len(salesmans))
	lastPage := 1
	return salesmans, total, lastPage, nil
}
