package repository

import (
	"context"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryGamificationImpl struct {
		*gorm.DB
	}
)
type GamificationRepository interface {
	FindAllByCustId(dataFilter entity.GamificationQueryFilter) ([]model.GamificationList, int64, int, error)
	Store(c context.Context, data *model.Gamification) error
	StoreParticipant(c context.Context, data *model.GamificationParticipant) error
	StoreProduct(c context.Context, data *model.GamificationProduct) error
	StoreOutlet(c context.Context, data *model.GamificationOutlet) error
	StoreAnnouncement(c context.Context, data *model.GamificationAnnouncement) error
	StoreFilter(c context.Context, data *model.GamificationFilter) error
	FindOneByGamificationNo(gamificationNo string, custId string, parentCustId string) (gamification model.GamificationRead, err error)
	FindGamificationParticipants(gamificationNo string, custId string, parentCustId string) (participants []model.GamificationParticipantRead, err error)
	FindGamificationProducts(gamificationNo string, custId string, parentCustId string) (products []model.GamificationProductRead, err error)
	FindGamificationOutlets(gamificationNo string, custId string, parentCustId string) (outlets []model.GamificationOutletRead, err error)
	FindGamificationAnnouncements(gamificationNo string, custId string, parentCustId string) (announcements []model.GamificationAnnouncementRead, err error)
	FindGamificationRankingsByRevenue(gamification model.GamificationRead, custId string, parentCustId string) (rankings []model.GamificationRankingRead, err error)
	FindGamificationRankingsByQuantity(gamification model.GamificationRead, custId string, parentCustId string) (rankings []model.GamificationRankingRead, err error)
	FindGamificationRankingsByTotal(gamification model.GamificationRead, custId string, parentCustId string) (rankings []model.GamificationRankingRead, err error)
	FindGamificationFilter(gamificationNo string, custId string, parentCustId string) (filter model.GamificationFilterRead, err error)
	Update(c context.Context, data *model.Gamification) error
	DeleteParticipant(c context.Context, gamificationNo string) error
	DeleteProduct(c context.Context, gamificationNo string) error
	DeleteOutlet(c context.Context, gamificationNo string) error
	DeleteAnnouncement(c context.Context, gamificationNo string) error
	DeleteFilter(c context.Context, gamificationNo string) error

	FindAllGamificationStatusesLookupMode(dataFilter entity.GeneralQueryFilter) (gamificationStatuses []model.GamificationStatusesFilter, total int64, lastPage int, err error)
	FindAllEmployeeGroupByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (employeeGroups []model.EmployeeGroups, total int64, lastPage int, err error)
	FindAllSalesTeamByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (salesTeams []model.SalesTeams, total int64, lastPage int, err error)
	FindAllOperationTypeByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (operationTypes []model.OperationTypes, total int64, lastPage int, err error)
	FindAllWarehouseByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (warehouses []model.Warehouses, total int64, lastPage int, err error)
	FindAllSalesmanByCustIdLookupMode(dataFilter entity.SalesmanQueryFilter) (salesmans []model.Salesmans, total int64, lastPage int, err error)
	FindAllProductCategoryByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (productCategories []model.ProductCategories, total int64, lastPage int, err error)
	FindAllProductLineByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (productLines []model.ProductLines, total int64, lastPage int, err error)
	FindAllBrandByCustIdLookupMode(dataFilter entity.BrandQueryFilter) (brands []model.Brands, total int64, lastPage int, err error)
	FindAllSubBrand1ByCustIdLookupMode(dataFilter entity.SubBrand1QueryFilter) (subBrands1 []model.SubBrands1, total int64, lastPage int, err error)
	FindAllSubBrand2ByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (subBrands2 []model.SubBrands2, total int64, lastPage int, err error)
	FindAllProductByCustIdLookupMode(dataFilter entity.ProductQueryFilter) (products []model.Products, total int64, lastPage int, err error)
	FindAllOutletLocationByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (outletLocations []model.OutletLocations, total int64, lastPage int, err error)
	FindAllOutletGroupByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (outletGroups []model.OutletGroups, total int64, lastPage int, err error)
	FindAllOutletClassByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (outletClasses []model.OutletClasses, total int64, lastPage int, err error)
	FindAllMarketByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (markets []model.Markets, total int64, lastPage int, err error)
	FindAllDistrictByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (districts []model.Districts, total int64, lastPage int, err error)
	FindAllIndustryByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (industries []model.Industries, total int64, lastPage int, err error)
	FindAllOutletByCustIdLookupMode(dataFilter entity.OutletQueryFilter) (outlets []model.Outlets, total int64, lastPage int, err error)
}

func NewGamificationRepo(db *gorm.DB) *RepositoryGamificationImpl {
	return &RepositoryGamificationImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryGamificationImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryGamificationImpl) FindAllByCustId(dataFilter entity.GamificationQueryFilter) ([]model.GamificationList, int64, int, error) {
	var rtn []model.GamificationList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gamification_no")
	query := repository.Select(`
			sls.gamification.gamification_no, 
			sls.gamification.gamification_date, 
			sls.gamification.title, 
			sls.gamification.start_date, 
			sls.gamification.end_date, 
			sls.gamification.finished_date,
			sls.gamification.gamification_status
		`)

	queryCount.Where("sls.gamification.cust_id=?", dataFilter.CustId)
	query.Where("sls.gamification.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.gamification.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.gamification.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.gamification.gamification_status in ?", dataFilter.Status)
		query.Where("sls.gamification.gamification_status in ?", dataFilter.Status)
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.gamification.gamification_no LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("sls.gamification.gamification_no LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("gamification_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&rtn).Error
	if err != nil {
		return rtn, total, 0, err
	}
	err = queryCount.Model(&rtn).Count(&total).Error
	if err != nil {
		return rtn, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return rtn, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllGamificationStatusesLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.GamificationStatusesFilter, int64, int, error) {

	var gamificationStatuses []model.GamificationStatusesFilter

	var total int64

	queryCount := repository.Select("sls.gamification.gamification_status")
	query := repository.Select(`sls.gamification.gamification_status as gamification_status`)

	queryCount.Where("sls.gamification.cust_id=?", dataFilter.CustId)
	query.Where("sls.gamification.cust_id=?", dataFilter.CustId)

	queryCount.Where("sls.gamification.gamification_status IS NOT NULL")
	query.Where("sls.gamification.gamification_status IS NOT NULL")

	queryCount.Group("sls.gamification.gamification_status")
	query.Group("sls.gamification.gamification_status")

	query.Order("sls.gamification.gamification_status ASC")

	err := query.Find(&gamificationStatuses).Error
	if err != nil {
		return gamificationStatuses, total, 0, err
	}

	total = int64(len(gamificationStatuses))
	lastPage := 1
	return gamificationStatuses, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllEmployeeGroupByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.EmployeeGroups, int64, int, error) {

	var empGroups []model.EmployeeGroups

	var total int64

	queryCount := repository.Select("emp_grp_id")
	query := repository.Select(`mst.m_emp_group.emp_grp_id, mst.m_emp_group.emp_grp_code, mst.m_emp_group.emp_grp_name`)

	queryCount.Where("lower(mst.m_emp_group.emp_grp_name) = 'salesman'")
	query.Where("lower(mst.m_emp_group.emp_grp_name) = 'salesman'")

	queryCount.Where("mst.m_emp_group.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_emp_group.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_emp_group.is_active=?", true)
	query.Where("mst.m_emp_group.is_active=?", true)

	queryCount.Where("mst.m_emp_group.is_del=?", false)
	query.Where("mst.m_emp_group.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_emp_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_emp_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_emp_group.emp_grp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_emp_group.emp_grp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("emp_grp_id DESC")
	}

	err := query.Find(&empGroups).Error
	if err != nil {
		return empGroups, total, 0, err
	}
	err = queryCount.Model(&empGroups).Count(&total).Error
	if err != nil {
		return empGroups, total, 0, err
	}

	lastPage := 1
	return empGroups, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllSalesTeamByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.SalesTeams, int64, int, error) {

	var salesTeams []model.SalesTeams

	var total int64

	queryCount := repository.Select("sales_team_id")
	query := repository.Select(`mst.m_sales_team.sales_team_id, mst.m_sales_team.sales_team_code, mst.m_sales_team.sales_team_name`)

	queryCount.Where("mst.m_sales_team.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_sales_team.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_sales_team.is_active=?", true)
	query.Where("mst.m_sales_team.is_active=?", true)

	queryCount.Where("mst.m_sales_team.is_del=?", false)
	query.Where("mst.m_sales_team.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_sales_team.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_sales_team.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_sales_team.sales_team_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_sales_team.sales_team_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("sales_team_id DESC")
	}

	err := query.Find(&salesTeams).Error
	if err != nil {
		return salesTeams, total, 0, err
	}
	err = queryCount.Model(&salesTeams).Count(&total).Error
	if err != nil {
		return salesTeams, total, 0, err
	}

	lastPage := 1
	return salesTeams, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllOperationTypeByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OperationTypes, int64, int, error) {

	var operationTypes []model.OperationTypes

	var total int64

	queryCount := repository.Select("operation_type_code")
	query := repository.Select(`mst.m_operation_type.operation_type_code, mst.m_operation_type.operation_type_name`)

	// queryCount.Where("mst.m_operation_type.cust_id=?", dataFilter.ParentCustId)
	// query.Where("mst.m_operation_type.cust_id=?", dataFilter.ParentCustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_operation_type.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_operation_type.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_operation_type.operation_type_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_operation_type.operation_type_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("operation_type_code DESC")
	}

	err := query.Find(&operationTypes).Error
	if err != nil {
		return operationTypes, total, 0, err
	}
	err = queryCount.Model(&operationTypes).Count(&total).Error
	if err != nil {
		return operationTypes, total, 0, err
	}

	lastPage := 1
	return operationTypes, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllWarehouseByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.Warehouses, int64, int, error) {

	var warehouses []model.Warehouses

	var total int64

	queryCount := repository.Select("wh_id")
	query := repository.Select(`mst.m_warehouse.wh_id, mst.m_warehouse.wh_code, mst.m_warehouse.wh_name`)

	queryCount.Where("mst.m_warehouse.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_warehouse.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_warehouse.is_active=?", true)
	query.Where("mst.m_warehouse.is_active=?", true)

	queryCount.Where("mst.m_warehouse.is_del=?", false)
	query.Where("mst.m_warehouse.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_warehouse.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_warehouse.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_warehouse.wh_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_warehouse.wh_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("wh_id DESC")
	}

	err := query.Find(&warehouses).Error
	if err != nil {
		return warehouses, total, 0, err
	}
	err = queryCount.Model(&warehouses).Count(&total).Error
	if err != nil {
		return warehouses, total, 0, err
	}

	lastPage := 1
	return warehouses, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllSalesmanByCustIdLookupMode(dataFilter entity.SalesmanQueryFilter) ([]model.Salesmans, int64, int, error) {

	var salesmans []model.Salesmans

	var total int64

	queryCount := repository.Select("employee.emp_id").
		Joins("left join mst.m_employee employee on employee.emp_id = mst.m_salesman.emp_id AND employee.cust_id = ?", dataFilter.CustId)
	query := repository.Select(`employee.emp_id as salesman_id, employee.emp_code as salesman_code, employee.emp_name as salesman_name`).
		Joins("left join mst.m_employee employee on employee.emp_id = mst.m_salesman.emp_id AND employee.cust_id = ?", dataFilter.CustId)

	if len(dataFilter.SalesTeamId) > 0 {
		queryCount.Where("mst.m_salesman.sales_team_id in ?", dataFilter.SalesTeamId)
		query.Where("mst.m_salesman.sales_team_id in ?", dataFilter.SalesTeamId)
	}

	if len(dataFilter.OprType) > 0 {
		queryCount.Where("mst.m_salesman.opr_type in ?", dataFilter.OprType)
		query.Where("mst.m_salesman.opr_type in ?", dataFilter.OprType)
	}

	if len(dataFilter.WhId) > 0 {
		queryCount.Where("mst.m_salesman.wh_id in ?", dataFilter.WhId)
		query.Where("mst.m_salesman.wh_id in ?", dataFilter.WhId)
	}

	queryCount.Where("mst.m_salesman.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_salesman.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_salesman.is_active=?", true)
	query.Where("mst.m_salesman.is_active=?", true)

	queryCount.Where("mst.m_salesman.is_del=?", false)
	query.Where("mst.m_salesman.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_salesman.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_salesman.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_salesman.salesman_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_salesman.wh_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("salesman_id DESC")
	}

	err := query.Find(&salesmans).Error
	if err != nil {
		return salesmans, total, 0, err
	}
	err = queryCount.Model(&salesmans).Count(&total).Error
	if err != nil {
		return salesmans, total, 0, err
	}

	lastPage := 1
	return salesmans, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllProductCategoryByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.ProductCategories, int64, int, error) {

	var productCategories []model.ProductCategories

	var total int64

	queryCount := repository.Select("pcat_id")
	query := repository.Select(`mst.m_product_cat.pcat_id, mst.m_product_cat.pcat_code, mst.m_product_cat.pcat_name`)

	queryCount.Where("mst.m_product_cat.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_product_cat.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_product_cat.is_active=?", true)
	query.Where("mst.m_product_cat.is_active=?", true)

	queryCount.Where("mst.m_product_cat.is_del=?", false)
	query.Where("mst.m_product_cat.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_product_cat.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_product_cat.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_product_cat.pcat_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_product_cat.pcat_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("pcat_id DESC")
	}

	err := query.Find(&productCategories).Error
	if err != nil {
		return productCategories, total, 0, err
	}
	err = queryCount.Model(&productCategories).Count(&total).Error
	if err != nil {
		return productCategories, total, 0, err
	}

	lastPage := 1
	return productCategories, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllProductLineByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.ProductLines, int64, int, error) {

	var productLines []model.ProductLines

	var total int64

	queryCount := repository.Select("pl_id")
	query := repository.Select(`mst.m_product_line.pl_id, mst.m_product_line.pl_code, mst.m_product_line.pl_name`)

	queryCount.Where("mst.m_product_line.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_product_line.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_product_line.is_active=?", true)
	query.Where("mst.m_product_line.is_active=?", true)

	queryCount.Where("mst.m_product_line.is_del=?", false)
	query.Where("mst.m_product_line.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_product_line.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_product_line.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_product_line.pl_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_product_line.pl_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("pl_id DESC")
	}

	err := query.Find(&productLines).Error
	if err != nil {
		return productLines, total, 0, err
	}
	err = queryCount.Model(&productLines).Count(&total).Error
	if err != nil {
		return productLines, total, 0, err
	}

	lastPage := 1
	return productLines, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllBrandByCustIdLookupMode(dataFilter entity.BrandQueryFilter) ([]model.Brands, int64, int, error) {

	var brands []model.Brands

	var total int64

	queryCount := repository.Select("brand_id")
	query := repository.Select(`mst.m_brand.brand_id, mst.m_brand.brand_code, mst.m_brand.brand_name`)

	queryCount.Where("mst.m_brand.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_brand.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_brand.is_active=?", true)
	query.Where("mst.m_brand.is_active=?", true)

	queryCount.Where("mst.m_brand.is_del=?", false)
	query.Where("mst.m_brand.is_del=?", false)

	if len(dataFilter.PLId) > 0 {
		queryCount.Where("mst.m_brand.pl_id in ?", dataFilter.PLId)
		query.Where("mst.m_brand.pl_id in ?", dataFilter.PLId)
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_brand.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_brand.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_brand.brand_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_brand.brand_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("brand_id DESC")
	}

	err := query.Find(&brands).Error
	if err != nil {
		return brands, total, 0, err
	}
	err = queryCount.Model(&brands).Count(&total).Error
	if err != nil {
		return brands, total, 0, err
	}

	lastPage := 1
	return brands, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllSubBrand1ByCustIdLookupMode(dataFilter entity.SubBrand1QueryFilter) ([]model.SubBrands1, int64, int, error) {

	var subBrands1 []model.SubBrands1

	var total int64

	queryCount := repository.Select("sbrand1_id")
	query := repository.Select(`mst.m_sub_brand1.sbrand1_id, mst.m_sub_brand1.sbrand1_code, mst.m_sub_brand1.sbrand1_name`)

	queryCount.Where("mst.m_sub_brand1.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_sub_brand1.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_sub_brand1.is_active=?", true)
	query.Where("mst.m_sub_brand1.is_active=?", true)

	queryCount.Where("mst.m_sub_brand1.is_del=?", false)
	query.Where("mst.m_sub_brand1.is_del=?", false)

	if len(dataFilter.BrandId) > 0 {
		queryCount.Where("mst.m_sub_brand1.brand_id in ?", dataFilter.BrandId)
		query.Where("mst.m_sub_brand1.brand_id in ?", dataFilter.BrandId)
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_sub_brand1.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_sub_brand1.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_sub_brand1.sbrand1_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_sub_brand1.sbrand1_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("sbrand1_id DESC")
	}

	err := query.Find(&subBrands1).Error
	if err != nil {
		return subBrands1, total, 0, err
	}
	err = queryCount.Model(&subBrands1).Count(&total).Error
	if err != nil {
		return subBrands1, total, 0, err
	}

	lastPage := 1
	return subBrands1, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllSubBrand2ByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.SubBrands2, int64, int, error) {

	var subBrands2 []model.SubBrands2

	var total int64

	queryCount := repository.Select("sbrand2_id")
	query := repository.Select(`mst.m_sub_brand2.sbrand2_id, mst.m_sub_brand2.sbrand2_code, mst.m_sub_brand2.sbrand2_name`)

	queryCount.Where("mst.m_sub_brand2.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_sub_brand2.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_sub_brand2.is_active=?", true)
	query.Where("mst.m_sub_brand2.is_active=?", true)

	queryCount.Where("mst.m_sub_brand2.is_del=?", false)
	query.Where("mst.m_sub_brand2.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_sub_brand2.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_sub_brand2.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_sub_brand2.sbrand2_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_sub_brand2.sbrand2_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("sbrand2_id DESC")
	}

	err := query.Find(&subBrands2).Error
	if err != nil {
		return subBrands2, total, 0, err
	}
	err = queryCount.Model(&subBrands2).Count(&total).Error
	if err != nil {
		return subBrands2, total, 0, err
	}

	lastPage := 1
	return subBrands2, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllProductByCustIdLookupMode(dataFilter entity.ProductQueryFilter) ([]model.Products, int64, int, error) {

	var products []model.Products

	var total int64

	queryCount := repository.Select("pro_id")
	query := repository.Select(`mst.m_product.pro_id, mst.m_product.pro_code, mst.m_product.pro_name`)

	queryCount.Where("mst.m_product.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_product.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_product.is_active=?", true)
	query.Where("mst.m_product.is_active=?", true)

	queryCount.Where("mst.m_product.is_del=?", false)
	query.Where("mst.m_product.is_del=?", false)

	if len(dataFilter.PCatId) > 0 {
		queryCount.Where("mst.m_product.pcat_id in ?", dataFilter.PCatId)
		query.Where("mst.m_product.pcat_id in ?", dataFilter.PCatId)
	}

	if len(dataFilter.SBrand1Id) > 0 {
		queryCount.Where("mst.m_product.sbrand1_id in ?", dataFilter.SBrand1Id)
		query.Where("mst.m_product.sbrand1_id in ?", dataFilter.SBrand1Id)
	}

	if len(dataFilter.SBrand2Id) > 0 {
		queryCount.Where("mst.m_product.sbrand2_id in ?", dataFilter.SBrand2Id)
		query.Where("mst.m_product.sbrand2_id in ?", dataFilter.SBrand2Id)
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_product.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_product.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_product.pro_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_product.pro_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("pro_id DESC")
	}

	err := query.Find(&products).Error
	if err != nil {
		return products, total, 0, err
	}
	err = queryCount.Model(&products).Count(&total).Error
	if err != nil {
		return products, total, 0, err
	}

	lastPage := 1
	return products, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllOutletLocationByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OutletLocations, int64, int, error) {

	var outletLocations []model.OutletLocations

	var total int64

	queryCount := repository.Select("ot_loc_id")
	query := repository.Select(`mst.m_outlet_loc.ot_loc_id, mst.m_outlet_loc.ot_loc_code, mst.m_outlet_loc.ot_loc_name`)

	queryCount.Where("mst.m_outlet_loc.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_outlet_loc.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_outlet_loc.is_active=?", true)
	query.Where("mst.m_outlet_loc.is_active=?", true)

	queryCount.Where("mst.m_outlet_loc.is_del=?", false)
	query.Where("mst.m_outlet_loc.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_outlet_loc.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_outlet_loc.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet_loc.ot_loc_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet_loc.ot_loc_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("ot_loc_id DESC")
	}

	err := query.Find(&outletLocations).Error
	if err != nil {
		return outletLocations, total, 0, err
	}
	err = queryCount.Model(&outletLocations).Count(&total).Error
	if err != nil {
		return outletLocations, total, 0, err
	}

	lastPage := 1
	return outletLocations, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllOutletGroupByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OutletGroups, int64, int, error) {

	var outletGroups []model.OutletGroups

	var total int64

	queryCount := repository.Select("ot_grp_id")
	query := repository.Select(`mst.m_outlet_group.ot_grp_id, mst.m_outlet_group.ot_grp_code, mst.m_outlet_group.ot_grp_name`)

	queryCount.Where("mst.m_outlet_group.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_outlet_group.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_outlet_group.is_active=?", true)
	query.Where("mst.m_outlet_group.is_active=?", true)

	queryCount.Where("mst.m_outlet_group.is_del=?", false)
	query.Where("mst.m_outlet_group.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_outlet_group.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet_group.ot_grp_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("ot_grp_id DESC")
	}

	err := query.Find(&outletGroups).Error
	if err != nil {
		return outletGroups, total, 0, err
	}
	err = queryCount.Model(&outletGroups).Count(&total).Error
	if err != nil {
		return outletGroups, total, 0, err
	}

	lastPage := 1
	return outletGroups, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllOutletClassByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.OutletClasses, int64, int, error) {

	var outletClasses []model.OutletClasses

	var total int64

	queryCount := repository.Select("ot_class_id")
	query := repository.Select(`mst.m_outlet_class.ot_class_id, mst.m_outlet_class.ot_class_code, mst.m_outlet_class.ot_class_name`)

	queryCount.Where("mst.m_outlet_class.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_outlet_class.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_outlet_class.is_active=?", true)
	query.Where("mst.m_outlet_class.is_active=?", true)

	queryCount.Where("mst.m_outlet_class.is_del=?", false)
	query.Where("mst.m_outlet_class.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_outlet_class.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_outlet_class.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet_class.ot_class_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet_class.ot_class_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("ot_class_id DESC")
	}

	err := query.Find(&outletClasses).Error
	if err != nil {
		return outletClasses, total, 0, err
	}
	err = queryCount.Model(&outletClasses).Count(&total).Error
	if err != nil {
		return outletClasses, total, 0, err
	}

	lastPage := 1
	return outletClasses, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllMarketByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.Markets, int64, int, error) {

	var markets []model.Markets

	var total int64

	queryCount := repository.Select("market_id")
	query := repository.Select(`mst.m_market.market_id, mst.m_market.market_code, mst.m_market.market_name`)

	queryCount.Where("mst.m_market.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_market.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_market.is_active=?", true)
	query.Where("mst.m_market.is_active=?", true)

	queryCount.Where("mst.m_market.is_del=?", false)
	query.Where("mst.m_market.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_market.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_market.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_market.market_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_market.market_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("market_id DESC")
	}

	err := query.Find(&markets).Error
	if err != nil {
		return markets, total, 0, err
	}
	err = queryCount.Model(&markets).Count(&total).Error
	if err != nil {
		return markets, total, 0, err
	}

	lastPage := 1
	return markets, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllDistrictByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.Districts, int64, int, error) {

	var districts []model.Districts

	var total int64

	queryCount := repository.Select("district_id")
	query := repository.Select(`mst.m_district.district_id, mst.m_district.district_code, mst.m_district.district_name`)

	queryCount.Where("mst.m_district.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_district.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_district.is_active=?", true)
	query.Where("mst.m_district.is_active=?", true)

	queryCount.Where("mst.m_district.is_del=?", false)
	query.Where("mst.m_district.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_district.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_district.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_district.district_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_district.district_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("district_id DESC")
	}

	err := query.Find(&districts).Error
	if err != nil {
		return districts, total, 0, err
	}
	err = queryCount.Model(&districts).Count(&total).Error
	if err != nil {
		return districts, total, 0, err
	}

	lastPage := 1
	return districts, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllIndustryByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.Industries, int64, int, error) {

	var industries []model.Industries

	var total int64

	queryCount := repository.Select("industry_id")
	query := repository.Select(`mst.m_industry.industry_id, mst.m_industry.industry_code, mst.m_industry.industry_name`)

	queryCount.Where("mst.m_industry.cust_id=?", dataFilter.ParentCustId)
	query.Where("mst.m_industry.cust_id=?", dataFilter.ParentCustId)

	queryCount.Where("mst.m_industry.is_active=?", true)
	query.Where("mst.m_industry.is_active=?", true)

	queryCount.Where("mst.m_industry.is_del=?", false)
	query.Where("mst.m_industry.is_del=?", false)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_industry.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_industry.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_industry.industry_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_industry.industry_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("industry_id DESC")
	}

	err := query.Find(&industries).Error
	if err != nil {
		return industries, total, 0, err
	}
	err = queryCount.Model(&industries).Count(&total).Error
	if err != nil {
		return industries, total, 0, err
	}

	lastPage := 1
	return industries, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) FindAllOutletByCustIdLookupMode(dataFilter entity.OutletQueryFilter) ([]model.Outlets, int64, int, error) {

	var outlets []model.Outlets

	var total int64

	queryCount := repository.Select("outlet_id")
	query := repository.Select(`mst.m_outlet.outlet_id, mst.m_outlet.outlet_code, mst.m_outlet.outlet_name`)

	queryCount.Where("mst.m_outlet.cust_id=?", dataFilter.CustId)
	query.Where("mst.m_outlet.cust_id=?", dataFilter.CustId)

	queryCount.Where("mst.m_outlet.is_active=?", true)
	query.Where("mst.m_outlet.is_active=?", true)

	queryCount.Where("mst.m_outlet.is_del=?", false)
	query.Where("mst.m_outlet.is_del=?", false)

	if len(dataFilter.OtLocId) > 0 {
		queryCount.Where("mst.m_outlet.ot_loc_id in ?", dataFilter.OtLocId)
		query.Where("mst.m_outlet.ot_loc_id in ?", dataFilter.OtLocId)
	}

	if len(dataFilter.OtGrpId) > 0 {
		queryCount.Where("mst.m_outlet.ot_grp_id in ?", dataFilter.OtGrpId)
		query.Where("mst.m_outlet.ot_grp_id in ?", dataFilter.OtGrpId)
	}

	if len(dataFilter.OtClassId) > 0 {
		queryCount.Where("mst.m_outlet.ot_class_id in ?", dataFilter.OtClassId)
		query.Where("mst.m_outlet.ot_class_id in ?", dataFilter.OtClassId)
	}

	if len(dataFilter.MarketId) > 0 {
		queryCount.Where("mst.m_outlet.market_id in ?", dataFilter.MarketId)
		query.Where("mst.m_outlet.market_id in ?", dataFilter.MarketId)
	}

	if len(dataFilter.DistrictId) > 0 {
		queryCount.Where("mst.m_outlet.district_id in ?", dataFilter.DistrictId)
		query.Where("mst.m_outlet.district_id in ?", dataFilter.DistrictId)
	}

	if len(dataFilter.IndustryId) > 0 {
		queryCount.Where("mst.m_outlet.industry_id in ?", dataFilter.IndustryId)
		query.Where("mst.m_outlet.industry_id in ?", dataFilter.IndustryId)
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("mst.m_outlet.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("mst.m_outlet.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("mst.m_outlet.outlet_name LIKE ?", "%"+dataFilter.Query+"%")
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("outlet_id DESC")
	}

	err := query.Find(&outlets).Error
	if err != nil {
		return outlets, total, 0, err
	}
	err = queryCount.Model(&outlets).Count(&total).Error
	if err != nil {
		return outlets, total, 0, err
	}

	lastPage := 1
	return outlets, total, lastPage, nil
}

func (repository *RepositoryGamificationImpl) Store(c context.Context, data *model.Gamification) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGamificationImpl) StoreParticipant(c context.Context, data *model.GamificationParticipant) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGamificationImpl) StoreProduct(c context.Context, data *model.GamificationProduct) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGamificationImpl) StoreOutlet(c context.Context, data *model.GamificationOutlet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGamificationImpl) StoreAnnouncement(c context.Context, data *model.GamificationAnnouncement) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGamificationImpl) StoreFilter(c context.Context, data *model.GamificationFilter) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGamificationImpl) FindOneByGamificationNo(gamificationNo string, custId string, parentCustId string) (gamification model.GamificationRead, err error) {
	err = repository.Select(`
			sls.gamification.cust_id,  
			sls.gamification.gamification_no, 
			sls.gamification.gamification_date, 
			sls.gamification.title, 
			sls.gamification.description, 
			sls.gamification.emp_grp_id, 
			sls.gamification.product_filter, 
			sls.gamification.outlet_filter, 
			emp_grp.emp_grp_code, 
			emp_grp.emp_grp_name, 
			sls.gamification.source_id, 
			sls.gamification.measurement_id, 
			sls.gamification.start_date, 
			sls.gamification.end_date, 
			sls.gamification.finished_date, 
			sls.gamification.announcement_id, 
			sls.gamification.subannouncement_id, 
			sls.gamification.announcement_target, 
			sls.gamification.gamification_status
		`).
		Joins("left join mst.m_emp_group emp_grp on emp_grp.emp_grp_id = sls.gamification.emp_grp_id AND emp_grp.cust_id = ?", parentCustId).
		Where("sls.gamification.cust_id=? AND sls.gamification.gamification_no = ? ", custId, gamificationNo).
		Take(&gamification).Error

	return gamification, err
}

func (repository *RepositoryGamificationImpl) FindGamificationParticipants(gamificationNo string, custId string, parentCustId string) (participants []model.GamificationParticipantRead, err error) {
	err = repository.Select(`
			sls.gamification_participants.gamification_participant_id, 
			sls.gamification_participants.emp_id as salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name
		`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.gamification_participants.emp_id AND employee.cust_id = ?", custId).
		Where("sls.gamification_participants.cust_id=? AND sls.gamification_participants.gamification_no = ? ", custId, gamificationNo).
		Order("sls.gamification_participants.gamification_participant_id ASC").
		Find(&participants).Error

	return participants, err
}

func (repository *RepositoryGamificationImpl) FindGamificationProducts(gamificationNo string, custId string, parentCustId string) (products []model.GamificationProductRead, err error) {
	err = repository.Select(`
			sls.gamification_products.gamification_product_id, 
			sls.gamification_products.pro_id, 
			product.pro_code, 
			product.pro_name
		`).
		Joins("left join mst.m_product product on product.pro_id = sls.gamification_products.pro_id AND product.cust_id = ?", custId).
		Where("sls.gamification_products.cust_id=? AND sls.gamification_products.gamification_no = ? ", custId, gamificationNo).
		Order("sls.gamification_products.gamification_product_id ASC").
		Find(&products).Error

	return products, err
}

func (repository *RepositoryGamificationImpl) FindGamificationOutlets(gamificationNo string, custId string, parentCustId string) (outlets []model.GamificationOutletRead, err error) {
	err = repository.Select(`
			sls.gamification_outlets.gamification_outlet_id, 
			sls.gamification_outlets.outlet_id, 
			outlet.outlet_code, 
			outlet.outlet_name
		`).
		Joins("left join mst.m_outlet outlet on outlet.outlet_id = sls.gamification_outlets.outlet_id AND outlet.cust_id = ?", custId).
		Where("sls.gamification_outlets.cust_id=? AND sls.gamification_outlets.gamification_no = ? ", custId, gamificationNo).
		Order("sls.gamification_outlets.gamification_outlet_id ASC").
		Find(&outlets).Error

	return outlets, err
}

func (repository *RepositoryGamificationImpl) FindGamificationAnnouncements(gamificationNo string, custId string, parentCustId string) (announcements []model.GamificationAnnouncementRead, err error) {
	err = repository.Select(`
			sls.gamification_announcements.gamification_announcement_id, 
			sls.gamification_announcements.announcement_date
		`).
		Where("sls.gamification_announcements.cust_id=? AND sls.gamification_announcements.gamification_no = ? ", custId, gamificationNo).
		Order("sls.gamification_announcements.gamification_announcement_id ASC").
		Find(&announcements).Error

	return announcements, err
}

func (repository *RepositoryGamificationImpl) FindGamificationRankingsByRevenue(gamification model.GamificationRead, custId string, parentCustId string) (rankings []model.GamificationRankingRead, err error) {
	checkingDate := "ro_date"
	invoiceFlag := "ISNULL"
	productJoin := "left"
	outletJoin := "left"

	if gamification.SourceID == 2 {
		checkingDate = "invoice_date"
		invoiceFlag = "IS NOT NULL"
	}

	if gamification.ProductFilter == "selected" {
		productJoin = "inner"
	}

	if gamification.OutletFilter == "selected" {
		outletJoin = "inner"
	}

	queryRanking := `left join (
			SELECT 
				po.salesman_id,
				sum(sls.order_detail.amount) as total
			FROM "sls"."order_detail"
			inner join sls.order po on po.ro_no = sls.order_detail.ro_no and po.cust_id='` + custId + `' AND po."` + checkingDate + `" >= '` + gamification.StartDate.Format("2006-01-02") + `' AND po."` + checkingDate + `" <= '` + gamification.EndDate.Format("2006-01-02") + `' and po.invoice_no ` + invoiceFlag + ` 
			inner join sls.gamification_participants participant on participant.emp_id = po.salesman_id AND participant.gamification_no = '` + gamification.GamificationNo + `' AND participant.cust_id = '` + custId + `'
			` + outletJoin + ` join sls.gamification_outlets outlet on outlet.outlet_id = po.outlet_id AND outlet.gamification_no = '` + gamification.GamificationNo + `' AND outlet.cust_id = '` + custId + `'  
			` + productJoin + ` join sls.gamification_products product on product.pro_id = sls.order_detail.pro_id AND product.gamification_no = '` + gamification.GamificationNo + `' AND product.cust_id = '` + custId + `'  
			group by po.salesman_id
		) gamification_rankings on gamification_rankings.salesman_id = sls.gamification_participants.emp_id`

	err = repository.Select(`
			sls.gamification_participants.emp_id as salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name,
			coalesce(gamification_rankings.total, 0) as total
		`).
		Joins(queryRanking).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.gamification_participants.emp_id AND employee.cust_id = ?", custId).
		Where("sls.gamification_participants.gamification_no = ?", gamification.GamificationNo).
		Order("coalesce(gamification_rankings.total, 0) DESC, sls.gamification_participants.gamification_participant_id ASC").
		Find(&rankings).Error

	return rankings, err
}

func (repository *RepositoryGamificationImpl) FindGamificationRankingsByQuantity(gamification model.GamificationRead, custId string, parentCustId string) (rankings []model.GamificationRankingRead, err error) {
	checkingDate := "ro_date"
	invoiceFlag := "ISNULL"
	productJoin := "left"
	outletJoin := "left"

	if gamification.SourceID == 2 {
		checkingDate = "invoice_date"
		invoiceFlag = "IS NOT NULL"
	}

	if gamification.ProductFilter == "selected" {
		productJoin = "inner"
	}

	if gamification.OutletFilter == "selected" {
		outletJoin = "inner"
	}

	queryRanking := `left join (
			SELECT 
				po.salesman_id,
				sum(((sls.order_detail.qty1 * sls.order_detail.conv_unit2 * sls.order_detail.conv_unit3) + (sls.order_detail.qty2 * sls.order_detail.conv_unit3) + sls.order_detail.qty3)) as qty
			FROM "sls"."order_detail"
			inner join sls.order po on po.ro_no = sls.order_detail.ro_no and po.cust_id='` + custId + `' AND po."` + checkingDate + `" >= '` + gamification.StartDate.Format("2006-01-02") + `' AND po."` + checkingDate + `" <= '` + gamification.EndDate.Format("2006-01-02") + `' and po.invoice_no ` + invoiceFlag + ` 
			inner join sls.gamification_participants participant on participant.emp_id = po.salesman_id AND participant.gamification_no = '` + gamification.GamificationNo + `' AND participant.cust_id = '` + custId + `'
			` + outletJoin + ` join sls.gamification_outlets outlet on outlet.outlet_id = po.outlet_id AND outlet.gamification_no = '` + gamification.GamificationNo + `' AND outlet.cust_id = '` + custId + `'  
			` + productJoin + ` join sls.gamification_products product on product.pro_id = sls.order_detail.pro_id AND product.gamification_no = '` + gamification.GamificationNo + `' AND product.cust_id = '` + custId + `'   
			group by po.salesman_id
		) gamification_rankings on gamification_rankings.salesman_id = sls.gamification_participants.emp_id`

	err = repository.Select(`
			sls.gamification_participants.emp_id as salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name,
			coalesce(gamification_rankings.qty, 0) as total
		`).
		Joins(queryRanking).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.gamification_participants.emp_id AND employee.cust_id = ?", custId).
		Where("sls.gamification_participants.gamification_no = ?", gamification.GamificationNo).
		Order("coalesce(gamification_rankings.qty, 0) DESC, sls.gamification_participants.gamification_participant_id ASC").
		Find(&rankings).Error

	return rankings, err
}

func (repository *RepositoryGamificationImpl) FindGamificationRankingsByTotal(gamification model.GamificationRead, custId string, parentCustId string) (rankings []model.GamificationRankingRead, err error) {
	checkingDate := "ro_date"
	invoiceFlag := "ISNULL"
	productJoin := "left"
	outletJoin := "left"

	if gamification.SourceID == 2 {
		checkingDate = "invoice_date"
		invoiceFlag = "IS NOT NULL"
	}

	if gamification.ProductFilter == "selected" {
		productJoin = "inner"
	}

	if gamification.OutletFilter == "selected" {
		outletJoin = "inner"
	}

	queryRanking := `left join (
			SELECT 
				salesman_id,
				count(ro_no) as total
			FROM (
				SELECT 
					po.salesman_id,
					po.ro_no
				FROM "sls"."order_detail"
				inner join sls.order po on po.ro_no = sls.order_detail.ro_no and po.cust_id='` + custId + `' AND po."` + checkingDate + `" >= '` + gamification.StartDate.Format("2006-01-02") + `' AND po."` + checkingDate + `" <= '` + gamification.EndDate.Format("2006-01-02") + `' and po.invoice_no ` + invoiceFlag + ` 
				inner join sls.gamification_participants participant on participant.emp_id = po.salesman_id AND participant.gamification_no = '` + gamification.GamificationNo + `' AND participant.cust_id = '` + custId + `'
				` + outletJoin + ` join sls.gamification_outlets outlet on outlet.outlet_id = po.outlet_id AND outlet.gamification_no = '` + gamification.GamificationNo + `' AND outlet.cust_id = '` + custId + `'  
				` + productJoin + ` join sls.gamification_products product on product.pro_id = sls.order_detail.pro_id AND product.gamification_no = '` + gamification.GamificationNo + `' AND product.cust_id = '` + custId + `'   
				group by po.salesman_id, po.ro_no
			) rekap_po group by salesman_id	
		) gamification_rankings on gamification_rankings.salesman_id = sls.gamification_participants.emp_id`

	err = repository.Select(`
			sls.gamification_participants.emp_id as salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name,
			coalesce(gamification_rankings.total, 0) as total
		`).
		Joins(queryRanking).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.gamification_participants.emp_id AND employee.cust_id = ?", custId).
		Where("sls.gamification_participants.gamification_no = ?", gamification.GamificationNo).
		Order("coalesce(gamification_rankings.total, 0) DESC, sls.gamification_participants.gamification_participant_id ASC").
		Find(&rankings).Error

	return rankings, err
}

func (repository *RepositoryGamificationImpl) FindGamificationFilter(gamificationNo string, custId string, parentCustId string) (filter model.GamificationFilterRead, err error) {
	err = repository.Select(`
			sls.gamification_filters.gamification_filter_id, 
			sls.gamification_filters.sales_team,
			sls.gamification_filters.operation_type,
			sls.gamification_filters.warehouse,
			sls.gamification_filters.product_category,
			sls.gamification_filters.product_line,
			sls.gamification_filters.brand,
			sls.gamification_filters.sub_brand1,
			sls.gamification_filters.sub_brand2,
			sls.gamification_filters.outlet_location,
			sls.gamification_filters.outlet_group,
			sls.gamification_filters.outlet_class,
			sls.gamification_filters.market,
			sls.gamification_filters.district,
			sls.gamification_filters.industry
		`).
		Where("sls.gamification_filters.gamification_no = ? ", gamificationNo).
		Take(&filter).Error

	return filter, err
}

func (repository *RepositoryGamificationImpl) DeleteParticipant(c context.Context, gamificationNo string) error {
	var participants model.GamificationParticipant
	err := repository.model(c).Where("gamification_no=?", gamificationNo).Delete(&participants).Error

	return err
}

func (repository *RepositoryGamificationImpl) DeleteProduct(c context.Context, gamificationNo string) error {
	var products model.GamificationProduct
	err := repository.model(c).Where("gamification_no=?", gamificationNo).Delete(&products).Error

	return err
}

func (repository *RepositoryGamificationImpl) DeleteOutlet(c context.Context, gamificationNo string) error {
	var outlets model.GamificationOutlet
	err := repository.model(c).Where("gamification_no=?", gamificationNo).Delete(&outlets).Error

	return err
}

func (repository *RepositoryGamificationImpl) DeleteAnnouncement(c context.Context, gamificationNo string) error {
	var announcements model.GamificationAnnouncement
	err := repository.model(c).Where("gamification_no=?", gamificationNo).Delete(&announcements).Error

	return err
}

func (repository *RepositoryGamificationImpl) DeleteFilter(c context.Context, gamificationNo string) error {
	var filter model.GamificationFilter
	err := repository.model(c).Where("gamification_no=?", gamificationNo).Delete(&filter).Error

	return err
}

func (repository *RepositoryGamificationImpl) Update(c context.Context, data *model.Gamification) error {
	result := repository.model(c).Model(&data).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
