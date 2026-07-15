package repository

import (
	"context"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"time"

	// "sales/pkg/structs"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

const (
	// Common WHERE clause for outlet criteria queries
	outletCriteriaWhereClause = "promotion_outlet_criteria.promo_id = ? AND promotion_outlet_criteria.cust_id = ?"
	// Common ORDER BY clause for outlet criteria queries
	outletCriteriaOrderClause = "promotion_outlet_criteria.id ASC"
)

type (
	RepositoryPromotionV2Impl struct {
		*gorm.DB
	}
)
type PromotionV2Repository interface {
	Store(c context.Context, data *model.PromotionV2) error
	Update(c context.Context, promoID string, data *model.PromotionV2) error
	UpdateStatus(c context.Context, promoID string, status model.PromotionStatus) error
	StoreSlabs(c context.Context, data []model.PromotionV2Slabs) error
	StoreStrata(c context.Context, data []model.PromotionV2Strata) error
	StoreProductCriteria(c context.Context, data []model.PromotionProductCriteria) error
	StoreRewardProducts(c context.Context, data []model.PromotionRewardProduct) error
	StoreCoverageDistributors(c context.Context, data []model.PromotionCoverageDistributors) error
	StoreOutletCriteria(c context.Context, data *model.PromotionOutletCriteria) (string, error)
	StoreOutletsSelected(c context.Context, data []model.PromotionOutletsSelected) error
	StoreOutletAttributeType(c context.Context, data []model.PromotionOutletAttributeType) error
	StoreOutletAttributeSalesTeam(c context.Context, data []model.PromotionOutletAttributeSalesTeam) error
	StoreOutletAttributeGroup(c context.Context, data []model.PromotionOutletAttributeGroup) error
	StoreOutletAttributeClass(c context.Context, data []model.PromotionOutletAttributeClass) error
	DeleteSlabs(c context.Context, custID, promoID string) error
	DeleteStratas(c context.Context, custID, promoID string) error
	DeleteProductCriteria(c context.Context, custIDs []string, promoID string) error
	DeleteRewardProducts(c context.Context, custID, promoID string) error
	DeleteCoverageDistributors(c context.Context, custID, promoID string) error
	DeleteOutletCriteria(c context.Context, custID, promoID string) error
	FindByPromoID(params entity.DetailPromotionParams) (promotion model.PromotionV2, err error)
	FindPromoSlabsByPromoID(params entity.DetailPromotionParams) (slabs []model.PromotionV2Slabs, err error)
	FindPromoStratasByPromoID(params entity.DetailPromotionParams) (stratas []model.PromotionV2Strata, err error)
	FindPromoProductCriteriasByPromoID(params entity.DetailPromotionParams) (productCriterias []model.PromotionProductCriteria, err error)
	FindPromoRewardProductsByPromoID(params entity.DetailPromotionParams) (rewardProducts []model.PromotionRewardProduct, err error)
	FindCoverageDistributorsByPromoID(params entity.DetailPromotionParams) (coverageDistributors []model.PromotionCoverageDistributors, err error)
	FindOutletCriteriaByPromoID(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error)
	FindOutletCriteriaWithJoins(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error)
	FindOutletCriteriaWithMasterData(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error)
	FindOutletCriteriaWithInnerJoins(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error)
	FindOutletCriteriaWithConditionalJoins(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error)
	FindOutletCriteriaWithPreloads(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error)
	FindAllByCustID(dataFilter entity.PromotionV2QueryFilter) ([]model.PromotionV2, int64, int, error)
	ExistsPromo(custID, promoID string) (bool, error)
	FindPromoIDsByBaseName(custID, basePromoID string) ([]string, error)
	FindOutletByID(outletID int64, custID string) (outlet model.OutletPromo, err error)
	FindSalesmanByID(salesmanID int64, custID string) (salesman model.SalesmanPromo, err error)
	FindWarehouseByID(whID int64, custID string) (wh model.WarehousePromo, err error)
	FindActivePromotionsByOutletCriteria(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) (promotions []model.PromotionV2, err error)
	FindProductCriteriasByPromoIDs(promoIDs []string, custID string) (productCriterias []model.PromotionProductCriteria, err error)
	FindSlabsByPromoIDs(promoIDs []string, custID string) (slabs []model.PromotionV2Slabs, err error)
	FindStratasByPromoIDs(promoIDs []string, custID string) (strata []model.PromotionV2Strata, err error)
	GetAllRewardProductFromStockV2(req entity.ConsultPromoV2Req, ctx model.RewardContext) (rewardProducts []model.PromotionRewardProduct, err error)
	FindProductByID(productID int64) (product model.ProductRead, err error)
	CloseExpiredPromotionStatuses(expiredBefore time.Time) (int64, error)
}

func NewPromotionV2Repo(db *gorm.DB) *RepositoryPromotionV2Impl {
	return &RepositoryPromotionV2Impl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryPromotionV2Impl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryPromotionV2Impl) Store(c context.Context, data *model.PromotionV2) error {
	// log.Info("ExistingPromoID:", *data.ExistingPromoID)
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreStatusLog(c context.Context, data *model.PromoStatusLog) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) FindByPromoID(params entity.DetailPromotionParams) (promotion model.PromotionV2, err error) {
	err = repository.
		Select(`promo.promotions.*`).
		Where("promo.promotions.promo_id = ? AND promo.promotions.cust_id=?", params.PromoID, params.ParentCustId).
		Take(&promotion).Error
	return promotion, err
}

func (repository *RepositoryPromotionV2Impl) applyFilters(query, queryCount *gorm.DB, dataFilter entity.PromotionV2QueryFilter) (*gorm.DB, *gorm.DB) {
	if dataFilter.EffectiveFrom != nil && dataFilter.EffectiveTo != nil {
		effectiveFrom := str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom)
		effectiveTo := str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo)
		whereClause := `promo.promotions.effective_from BETWEEN ? AND ? OR promo.promotions.effective_to BETWEEN ? AND ?`
		query = query.Where(whereClause, effectiveFrom, effectiveTo, effectiveFrom, effectiveTo)
		queryCount = queryCount.Where(whereClause, effectiveFrom, effectiveTo, effectiveFrom, effectiveTo)
	}

	if dataFilter.Query != "" {
		searchPattern := "%" + dataFilter.Query + "%"
		whereClause := "promo.promotions.promo_id ILIKE ? OR promo.promotions.promo_desc ILIKE ?"
		query = query.Where(whereClause, searchPattern, searchPattern)
		queryCount = queryCount.Where(whereClause, searchPattern, searchPattern)
	}

	if dataFilter.PromoID != "" {
		searchPattern := "%" + dataFilter.PromoID + "%"
		whereClause := "promo.promotions.promo_id ILIKE ?"
		query = query.Where(whereClause, searchPattern)
		queryCount = queryCount.Where(whereClause, searchPattern)
	}

	if dataFilter.PromoDesc != "" {
		searchPattern := "%" + dataFilter.PromoDesc + "%"
		whereClause := "promo.promotions.promo_desc ILIKE ?"
		query = query.Where(whereClause, searchPattern)
		queryCount = queryCount.Where(whereClause, searchPattern)
	}

	if len(dataFilter.PromoStatus) > 0 {
		whereClause := "promo.promotions.promo_status IN ?"
		query = query.Where(whereClause, dataFilter.PromoStatus)
		queryCount = queryCount.Where(whereClause, dataFilter.PromoStatus)
	}

	return query, queryCount
}

func (repository *RepositoryPromotionV2Impl) FindAllByCustID(dataFilter entity.PromotionV2QueryFilter) ([]model.PromotionV2, int64, int, error) {
	var (
		promo []model.PromotionV2
		total int64
	)
	limit := 10
	if dataFilter.Limit != 0 {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("promo.promotions.promo_id")
	query := repository.Select(`promo.promotions.*`)

	// Apply customer filter conditions
	whereClause := "promo.promotions.cust_id=?"
	args := []interface{}{dataFilter.ParentCustID}
	if dataFilter.CustID != dataFilter.ParentCustID { // login as distributor
		if dataFilter.DistributorID > 0 {
			whereClause += ` AND promo.promotions.distributor_cust_id=?`
			args = append(args, dataFilter.CustID)
		} else {
			joinClause := `LEFT JOIN promo.promotion_coverage_distributors pcd ON pcd.cust_id = ? 
			AND promo.promotions.coverage = ? 
			AND pcd.promo_id = promo.promotions.promo_id`
			query = query.Joins(joinClause, dataFilter.ParentCustID, model.CoverageByDistributor)
			queryCount = queryCount.Joins(joinClause, dataFilter.ParentCustID, model.CoverageByDistributor)

			whereClause += ` AND promo.promotions.distributor_cust_id=? 
				OR (promo.promotions.distributor_cust_id=?
					AND promo.promotions.promo_status IN (?,?,?) 
					AND (
						promo.promotions.coverage = ? 
						OR (promo.promotions.coverage = ? AND pcd.distributor_id = ?)
					)
				)`
			args = append(args, dataFilter.CustID, dataFilter.ParentCustID, model.PromoStatusApproved, model.PromoStatusActive, model.PromoStatusInactive,
				model.CoverageNational, model.CoverageByDistributor, dataFilter.TokenDistID)
		}
	} else {
		whereClause += ` AND promo.promotions.distributor_cust_id=?`
		args = append(args, dataFilter.CustID)
	}
	query = query.Where(whereClause, args...)
	queryCount = queryCount.Where(whereClause, args...)

	query, queryCount = repository.applyFilters(query, queryCount, dataFilter)

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
		query = query.Order(sortBy)
	} else {
		query = query.Order("created_at DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&promo).Error
	if err != nil {
		return promo, total, 0, err
	}
	err = queryCount.Model(&promo).Count(&total).Error
	if err != nil {
		return promo, total, 0, err
	}

	promoIDs := make([]string, 0)
	for _, p := range promo {
		promoIDs = append(promoIDs, p.PromoID)
	}

	var coverages []struct {
		PromoID       string
		DistributorID int64
	}

	repository.Raw(`
		SELECT promo_id, distributor_id 
		FROM promo.promotion_coverage_distributors
		WHERE promo_id IN ?
	`, promoIDs).Scan(&coverages)

	// Convert to map
	coverageMap := make(map[string][]int64)
	for _, c := range coverages {
		coverageMap[c.PromoID] = append(coverageMap[c.PromoID], c.DistributorID)
	}

	for i := range promo {
		promoID := promo[i].PromoID
		promo[i].DistributorIDs = coverageMap[promoID]
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return promo, total, lastPage, nil
}

func (repository *RepositoryPromotionV2Impl) StoreSlabs(c context.Context, data []model.PromotionV2Slabs) error {
	err := repository.model(c).Create(&data).Error
	if err != nil {
		log.Error("Error storing slabs:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreStrata(c context.Context, data []model.PromotionV2Strata) error {
	// log.Info("data:", structs.StructToJson(data))
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing strata:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreProductCriteria(c context.Context, data []model.PromotionProductCriteria) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing product criteria:", err)
		return err
	}
	// log.Info("data.ID:", data.ID)
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreRewardProducts(c context.Context, data []model.PromotionRewardProduct) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing reward products:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreCoverageDistributors(c context.Context, data []model.PromotionCoverageDistributors) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing coverage distributors:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreOutletCriteria(c context.Context, data *model.PromotionOutletCriteria) (string, error) {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing outlet criteria:", err)
		return "", err
	}
	return data.ID, nil
}

func (repository *RepositoryPromotionV2Impl) StoreOutletsSelected(c context.Context, data []model.PromotionOutletsSelected) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing outlets selected:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreOutletAttributeType(c context.Context, data []model.PromotionOutletAttributeType) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing outlet attribute type:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreOutletAttributeSalesTeam(c context.Context, data []model.PromotionOutletAttributeSalesTeam) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing outlet attribute sales team:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreOutletAttributeGroup(c context.Context, data []model.PromotionOutletAttributeGroup) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing outlet attribute group:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) StoreOutletAttributeClass(c context.Context, data []model.PromotionOutletAttributeClass) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Error("Error storing outlet attribute class:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) FindPromoSlabsByPromoID(params entity.DetailPromotionParams) (promoSlabs []model.PromotionV2Slabs, err error) {
	err = repository.
		Select(`promotion_slabs.*`).
		Where("promotion_slabs.promo_id = ? AND promotion_slabs.cust_id = ?", params.PromoID, params.ParentCustId).
		Order("promotion_slabs.ordinal ASC").
		Find(&promoSlabs).Error
	return promoSlabs, err
}

func (repository *RepositoryPromotionV2Impl) FindPromoStratasByPromoID(params entity.DetailPromotionParams) (promoStratas []model.PromotionV2Strata, err error) {
	err = repository.
		Select(`promotion_strata.*`).
		Where("promotion_strata.promo_id = ? AND promotion_strata.cust_id = ?", params.PromoID, params.ParentCustId).
		Order("promotion_strata.ordinal ASC").
		Find(&promoStratas).Error
	return promoStratas, err
}

func (repository *RepositoryPromotionV2Impl) FindPromoProductCriteriasByPromoID(params entity.DetailPromotionParams) (promoProductCriterias []model.PromotionProductCriteria, err error) {
	custIDs := []string{params.ParentCustId}
	if params.CustID != "" && params.CustID != params.ParentCustId {
		custIDs = append(custIDs, params.CustID)
	}

	err = repository.
		Select(`promotion_product_criteria.*, mp.pro_code, mp.pro_name`).
		Joins("left join mst.m_product mp on mp.pro_id = promotion_product_criteria.pro_id AND mp.cust_id IN ?", custIDs).
		Where("promotion_product_criteria.promo_id = ? AND promotion_product_criteria.cust_id IN ?", params.PromoID, custIDs).
		Order("promotion_product_criteria.id ASC").
		Find(&promoProductCriterias).Error
	return promoProductCriterias, err
}

func (repository *RepositoryPromotionV2Impl) FindPromoRewardProductsByPromoID(params entity.DetailPromotionParams) (promoRewardProducts []model.PromotionRewardProduct, err error) {
	custIDs := []string{params.ParentCustId}
	if params.CustID != "" && params.CustID != params.ParentCustId {
		custIDs = append(custIDs, params.CustID)
	}

	err = repository.
		Select(`promotion_reward_products.*, mp.pro_code, mp.pro_name`).
		Joins("left join mst.m_product mp on mp.pro_id = promotion_reward_products.pro_id AND mp.cust_id IN ?", custIDs).
		Where("promotion_reward_products.promo_id = ? AND promotion_reward_products.cust_id = ?", params.PromoID, params.ParentCustId).
		Order("promotion_reward_products.id ASC").
		Find(&promoRewardProducts).Error
	return promoRewardProducts, err
}

func (repository *RepositoryPromotionV2Impl) FindCoverageDistributorsByPromoID(params entity.DetailPromotionParams) (coverageDistributors []model.PromotionCoverageDistributors, err error) {
	err = repository.
		Select(`promotion_coverage_distributors.*, md.distributor_code, md.distributor_name`).
		Joins(`LEFT JOIN smc.m_customer mc ON mc.distributor_id = promotion_coverage_distributors.distributor_id AND mc.parent_cust_id = ? AND mc.cust_id <> mc.parent_cust_id`, params.ParentCustId).
		Joins(`LEFT JOIN mst.m_distributor md ON md.distributor_id = promotion_coverage_distributors.distributor_id AND (md.cust_id = mc.cust_id OR md.cust_id = ?)`, params.ParentCustId).
		Where("promotion_coverage_distributors.promo_id = ? AND promotion_coverage_distributors.cust_id = ?", params.PromoID, params.ParentCustId).
		Order("promotion_coverage_distributors.id ASC").
		Find(&coverageDistributors).Error
	return coverageDistributors, err
}

func (repository *RepositoryPromotionV2Impl) FindOutletCriteriaByPromoID(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error) {
	err = repository.
		Select(`promotion_outlet_criteria.*, mo.outlet_code, mo.outlet_name, mot.outlet_type_code, mot.outlet_type_name, mst.sales_team_code, mst.sales_team_name, mog.outlet_group_code, mog.outlet_group_name, moc.outlet_class_code, moc.outlet_class_name`).
		Joins("left join mst.m_outlet mo ON mo.outlet_id = promotion_outlet_criteria.outlet_id AND mo.cust_id = ?", params.ParentCustId).
		Joins("left join mst.m_outlet_type mot ON mot.outlet_type_id = promotion_outlet_criteria.outlet_type_id AND mot.cust_id = ?", params.ParentCustId).
		Joins("left join mst.m_sales_team mst ON mst.sales_team_id = promotion_outlet_criteria.sales_team_id AND mst.cust_id = ?", params.ParentCustId).
		Joins("left join mst.m_outlet_group mog ON mog.outlet_group_id = promotion_outlet_criteria.outlet_group_id AND mog.cust_id = ?", params.ParentCustId).
		Joins("left join mst.m_outlet_class moc ON moc.outlet_class_id = promotion_outlet_criteria.outlet_class_id AND moc.cust_id = ?", params.ParentCustId).
		Where(outletCriteriaWhereClause, params.PromoID, params.ParentCustId).
		Order(outletCriteriaOrderClause).
		Find(&outletCriteria).Error
	return outletCriteria, err
}

func (repository *RepositoryPromotionV2Impl) CheckPromotionExists(c context.Context, promoID string, custID string) (exists bool, err error) {
	var count int64
	err = repository.WithContext(c).
		Model(&model.PromotionV2{}).
		Where("promo_id = ? AND cust_id = ?", promoID, custID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindOutletCriteriaWithJoins demonstrates how to query with joins for all related tables
func (repository *RepositoryPromotionV2Impl) FindOutletCriteriaWithJoins(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error) {
	err = repository.
		Select(`
			promotion_outlet_criteria.*,
			-- Selected outlets data
			pos.outlet_id as selected_outlet_id,
			-- Attribute types data  
			poat.outlet_type_id as attribute_type_id,
			-- Attribute sales teams data
			poast.sales_team_id as attribute_sales_team_id,
			-- Attribute groups data
			poag.outlet_group_id as attribute_group_id,
			-- Attribute classes data
			poac.outlet_class_id as attribute_class_id
		`).
		Joins("LEFT JOIN promotion_outlets_selected pos ON pos.criteria_id = promotion_outlet_criteria.id AND pos.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_type poat ON poat.criteria_id = promotion_outlet_criteria.id AND poat.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_sales_team poast ON poast.criteria_id = promotion_outlet_criteria.id AND poast.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_group poag ON poag.criteria_id = promotion_outlet_criteria.id AND poag.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_class poac ON poac.criteria_id = promotion_outlet_criteria.id AND poac.cust_id = ?", params.ParentCustId).
		Where(outletCriteriaWhereClause, params.PromoID, params.ParentCustId).
		Order(outletCriteriaOrderClause).
		Find(&outletCriteria).Error
	return outletCriteria, err
}

// FindOutletCriteriaWithMasterData demonstrates joins with master data tables
func (repository *RepositoryPromotionV2Impl) FindOutletCriteriaWithMasterData(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error) {
	err = repository.
		Select(`
			promotion_outlet_criteria.*,
			-- Selected outlets with master outlet data
			pos.outlet_id as selected_outlet_id,
			mo.outlet_code as selected_outlet_code, 
			mo.outlet_name as selected_outlet_name, 
			-- Attribute types with master outlet type data
			poat.outlet_type_id as attribute_type_id,
			mot.ot_type_code as attribute_type_code, 
  			mot.ot_type_name as attribute_type_name, 
			-- Attribute sales teams with master sales team data
			poast.sales_team_id as attribute_sales_team_id,
			mst.sales_team_code as attribute_sales_team_code, 
  			mst.sales_team_name as attribute_sales_team_name, 
			-- Attribute groups with master outlet group data
			poag.outlet_group_id as attribute_group_id,
			mog.ot_grp_code as attribute_group_code, 
  			mog.ot_grp_name as attribute_group_name, 
			-- Attribute classes with master outlet class data
			poac.outlet_class_id as attribute_class_id,
			moc.ot_class_code as attribute_class_code, 
  			moc.ot_class_name as attribute_class_name 
		`).
		// Join with promotion tables
		Joins("LEFT JOIN promo.promotion_outlets_selected pos ON pos.criteria_id = promotion_outlet_criteria.id AND pos.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promo.promotion_outlet_attribute_type poat ON poat.criteria_id = promotion_outlet_criteria.id AND poat.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promo.promotion_outlet_attribute_sales_team poast ON poast.criteria_id = promotion_outlet_criteria.id AND poast.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promo.promotion_outlet_attribute_group poag ON poag.criteria_id = promotion_outlet_criteria.id AND poag.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN promo.promotion_outlet_attribute_class poac ON poac.criteria_id = promotion_outlet_criteria.id AND poac.cust_id = ?", params.ParentCustId).
		// Join with master data tables (assuming these exist in your schema)
		Joins("LEFT JOIN mst.m_outlet mo ON mo.outlet_id = pos.outlet_id AND mo.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN mst.m_outlet_type mot ON mot.ot_type_id = poat.outlet_type_id AND mot.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN mst.m_sales_team mst ON mst.sales_team_id = poast.sales_team_id AND mst.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN mst.m_outlet_group mog ON mog.ot_grp_id = poag.outlet_group_id AND mog.cust_id = ?", params.ParentCustId).
		Joins("LEFT JOIN mst.m_outlet_class moc ON moc.ot_class_id = poac.outlet_class_id AND moc.cust_id = ?", params.ParentCustId).
		Where(outletCriteriaWhereClause, params.PromoID, params.ParentCustId).
		Order(outletCriteriaOrderClause).
		Find(&outletCriteria).Error
	return outletCriteria, err
}

// FindOutletCriteriaWithInnerJoins demonstrates INNER JOIN usage
func (repository *RepositoryPromotionV2Impl) FindOutletCriteriaWithInnerJoins(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error) {
	err = repository.
		Select(`
			promotion_outlet_criteria.*,
			pos.outlet_id as selected_outlet_id,
			poat.outlet_type_id as attribute_type_id,
			poast.sales_team_id as attribute_sales_team_id,
			poag.outlet_group_id as attribute_group_id,
			poac.outlet_class_id as attribute_class_id
		`).
		// Use INNER JOIN to only get records that have related data
		Joins("INNER JOIN promotion_outlets_selected pos ON pos.criteria_id = promotion_outlet_criteria.id AND pos.cust_id = ?", params.ParentCustId).
		Joins("INNER JOIN promotion_outlet_attribute_type poat ON poat.criteria_id = promotion_outlet_criteria.id AND poat.cust_id = ?", params.ParentCustId).
		Joins("INNER JOIN promotion_outlet_attribute_sales_team poast ON poast.criteria_id = promotion_outlet_criteria.id AND poast.cust_id = ?", params.ParentCustId).
		Joins("INNER JOIN promotion_outlet_attribute_group poag ON poag.criteria_id = promotion_outlet_criteria.id AND poag.cust_id = ?", params.ParentCustId).
		Joins("INNER JOIN promotion_outlet_attribute_class poac ON poac.criteria_id = promotion_outlet_criteria.id AND poac.cust_id = ?", params.ParentCustId).
		Where(outletCriteriaWhereClause, params.PromoID, params.ParentCustId).
		Order(outletCriteriaOrderClause).
		Find(&outletCriteria).Error
	return outletCriteria, err
}

// FindOutletCriteriaWithConditionalJoins demonstrates conditional joins based on selection type
func (repository *RepositoryPromotionV2Impl) FindOutletCriteriaWithConditionalJoins(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error) {
	err = repository.
		Select(`
			promotion_outlet_criteria.*,
			CASE 
				WHEN promotion_outlet_criteria.selection_type = 'by_selection' THEN pos.outlet_id
				ELSE NULL
			END as selected_outlet_id,
			CASE 
				WHEN promotion_outlet_criteria.selection_type = 'by_attribute' THEN poat.outlet_type_id
				ELSE NULL
			END as attribute_type_id,
			CASE 
				WHEN promotion_outlet_criteria.selection_type = 'by_attribute' THEN poast.sales_team_id
				ELSE NULL
			END as attribute_sales_team_id,
			CASE 
				WHEN promotion_outlet_criteria.selection_type = 'by_attribute' THEN poag.outlet_group_id
				ELSE NULL
			END as attribute_group_id,
			CASE 
				WHEN promotion_outlet_criteria.selection_type = 'by_attribute' THEN poac.outlet_class_id
				ELSE NULL
			END as attribute_class_id
		`).
		Joins("LEFT JOIN promotion_outlets_selected pos ON pos.criteria_id = promotion_outlet_criteria.id AND pos.cust_id = ? AND promotion_outlet_criteria.selection_type = 'by_selection'", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_type poat ON poat.criteria_id = promotion_outlet_criteria.id AND poat.cust_id = ? AND promotion_outlet_criteria.selection_type = 'by_attribute'", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_sales_team poast ON poast.criteria_id = promotion_outlet_criteria.id AND poast.cust_id = ? AND promotion_outlet_criteria.selection_type = 'by_attribute'", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_group poag ON poag.criteria_id = promotion_outlet_criteria.id AND poag.cust_id = ? AND promotion_outlet_criteria.selection_type = 'by_attribute'", params.ParentCustId).
		Joins("LEFT JOIN promotion_outlet_attribute_class poac ON poac.criteria_id = promotion_outlet_criteria.id AND poac.cust_id = ? AND promotion_outlet_criteria.selection_type = 'by_attribute'", params.ParentCustId).
		Where(outletCriteriaWhereClause, params.PromoID, params.ParentCustId).
		Order(outletCriteriaOrderClause).
		Find(&outletCriteria).Error
	return outletCriteria, err
}

// FindOutletCriteriaWithPreloads uses GORM preloading to load related data
func (repository *RepositoryPromotionV2Impl) FindOutletCriteriaWithPreloads(params entity.DetailPromotionParams) (outletCriteria []model.PromotionOutletCriteria, err error) {
	custIDs := []string{params.ParentCustId}
	if params.CustID != "" && params.CustID != params.ParentCustId {
		custIDs = append(custIDs, params.CustID)
	}

	err = repository.
		Preload("SelectedOutlets", func(db *gorm.DB) *gorm.DB {
			return db.Joins(`LEFT JOIN mst.m_outlet mo ON mo.outlet_id = promotion_outlets_selected.outlet_id AND mo.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)`, params.ParentCustId).
				Joins(`LEFT JOIN smc.m_customer mc ON mc.cust_id = mo.cust_id AND mc.parent_cust_id = ? AND mc.cust_id <> mc.parent_cust_id`, params.ParentCustId).
				Joins(`LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND (md.cust_id = mc.cust_id OR md.cust_id = ?)`, params.ParentCustId).
				Select("promotion_outlets_selected.*, mo.outlet_code AS outlet_code, mo.outlet_name AS outlet_name, md.distributor_code AS distributor_code, md.distributor_name AS distributor_name")
		}).
		Preload("AttributeTypes", func(db *gorm.DB) *gorm.DB {
			return db.Joins("LEFT JOIN mst.m_outlet_type mot ON mot.ot_type_id = promotion_outlet_attribute_type.outlet_type_id AND mot.cust_id = ?", params.ParentCustId).
				Select("promotion_outlet_attribute_type.*, mot.ot_type_code AS outlet_type_code, mot.ot_type_name AS outlet_type_name")
		}).
		Preload("AttributeSalesTeams", func(db *gorm.DB) *gorm.DB {
			return db.Joins("LEFT JOIN mst.m_sales_team mst ON mst.sales_team_id = promotion_outlet_attribute_sales_team.sales_team_id AND mst.cust_id IN ?", custIDs).
				Select("promotion_outlet_attribute_sales_team.*, mst.sales_team_code AS sales_team_code, mst.sales_team_name AS sales_team_name")
		}).
		Preload("AttributeGroups", func(db *gorm.DB) *gorm.DB {
			return db.Joins("LEFT JOIN mst.m_outlet_group mog ON mog.ot_grp_id = promotion_outlet_attribute_group.outlet_group_id AND mog.cust_id = ?", params.ParentCustId).
				Select("promotion_outlet_attribute_group.*, mog.ot_grp_code AS outlet_group_code, mog.ot_grp_name AS outlet_group_name")
		}).
		Preload("AttributeClasses", func(db *gorm.DB) *gorm.DB {
			return db.Joins("LEFT JOIN mst.m_outlet_class moc ON moc.ot_class_id = promotion_outlet_attribute_class.outlet_class_id AND moc.cust_id = ?", params.ParentCustId).
				Select("promotion_outlet_attribute_class.*, moc.ot_class_code AS outlet_class_code, moc.ot_class_name AS outlet_class_name")
		}).
		Where("promo_id = ? AND cust_id = ?", params.PromoID, params.ParentCustId).
		Order("created_at ASC").
		Find(&outletCriteria).Error
	return outletCriteria, err
}

func (repository *RepositoryPromotionV2Impl) Update(c context.Context, promoID string, data *model.PromotionV2) error {
	err := repository.model(c).Where("promo_id = ? AND cust_id = ? AND distributor_cust_id = ?", promoID, data.CustID, data.DistributorCustID).Updates(data).Error
	if err != nil {
		log.Error("Error updating promotion:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) UpdateStatus(c context.Context, promoID string, status model.PromotionStatus) error {
	err := repository.model(c).Model(&model.PromotionV2{}).Where("promo_id = ?", promoID).Update("promo_status", status).Error
	if err != nil {
		log.Error("Error updating promotion status:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) DeleteSlabs(c context.Context, custID, promoID string) error {
	err := repository.model(c).Where("cust_id = ? AND promo_id = ?", custID, promoID).Delete(&model.PromotionV2Slabs{}).Error
	if err != nil {
		log.Error("Error deleting slabs:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) DeleteStratas(c context.Context, custID, promoID string) error {
	err := repository.model(c).Where("cust_id = ? AND promo_id = ?", custID, promoID).Delete(&model.PromotionV2Strata{}).Error
	if err != nil {
		log.Error("Error deleting stratas:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) DeleteProductCriteria(c context.Context, custIDs []string, promoID string) error {
	err := repository.model(c).Where("cust_id IN ? AND promo_id = ?", custIDs, promoID).Delete(&model.PromotionProductCriteria{}).Error
	if err != nil {
		log.Error("Error deleting product criteria:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) DeleteRewardProducts(c context.Context, custID, promoID string) error {
	err := repository.model(c).Where("cust_id = ? AND promo_id = ?", custID, promoID).Delete(&model.PromotionRewardProduct{}).Error
	if err != nil {
		log.Error("Error deleting reward products:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) DeleteCoverageDistributors(c context.Context, custID, promoID string) error {
	err := repository.model(c).Where("cust_id = ? AND promo_id = ?", custID, promoID).Delete(&model.PromotionCoverageDistributors{}).Error
	if err != nil {
		log.Error("Error deleting coverage distributors:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) DeleteOutletCriteria(c context.Context, custID, promoID string) error {
	// First delete related outlet criteria records
	err := repository.model(c).Where("cust_id = ? AND promo_id = ?", custID, promoID).Delete(&model.PromotionOutletCriteria{}).Error
	if err != nil {
		log.Error("Error deleting outlet criteria:", err)
		return err
	}
	return nil
}

func (repository *RepositoryPromotionV2Impl) ExistsPromo(custID, promoID string) (bool, error) {
	var count int64
	err := repository.WithContext(context.Background()).Model(&model.PromotionV2{}).
		Where("cust_id = ? AND promo_id = ?", custID, promoID).
		Count(&count).Error
	return count > 0, err
}

func (repository *RepositoryPromotionV2Impl) FindPromoIDsByBaseName(custID, basePromoID string) ([]string, error) {
	var promoIDs []string
	err := repository.WithContext(context.Background()).
		Model(&model.PromotionV2{}).
		Select("promo_id").
		Where("cust_id = ? AND promo_id LIKE ?", custID, basePromoID+"%%").
		Pluck("promo_id", &promoIDs).Error
	return promoIDs, err
}

func (repository *RepositoryPromotionV2Impl) FindOutletByID(outletID int64, custID string) (outlet model.OutletPromo, err error) {
	err = repository.Select(`
			mst.m_outlet.outlet_id, 
			mst.m_outlet.outlet_code, 
			mst.m_outlet.outlet_name, 
			mst.m_outlet.ot_grp_id, 
			mst.m_outlet.ot_class_id, 
			mst.m_outlet.ot_type_id 
		`).
		Where("mst.m_outlet.cust_id = ?", custID).
		Where("mst.m_outlet.outlet_id = ?", outletID).
		Take(&outlet).Error

	return outlet, err
}

func (repository *RepositoryPromotionV2Impl) FindSalesmanByID(salesmanID int64, custID string) (salesman model.SalesmanPromo, err error) {
	err = repository.Select(`
			mst.m_salesman.emp_id AS salesman_id, 
			ep.emp_code AS salesman_code,
			mst.m_salesman.sales_name AS salesman_name, 
			mst.m_salesman.sales_team_id,
			mst.m_salesman.wh_id
		`).
		Joins("LEFT JOIN mst.m_employee ep ON ep.cust_id = m_salesman.cust_id AND ep.emp_id = m_salesman.emp_id").
		Where("mst.m_salesman.cust_id = ?", custID).
		Where("mst.m_salesman.emp_id = ?", salesmanID).
		Where("mst.m_salesman.is_active = true").
		Take(&salesman).Error

	return salesman, err
}

func (repository *RepositoryPromotionV2Impl) FindWarehouseByID(whID int64, custID string) (wh model.WarehousePromo, err error) {
	err = repository.Select(`
			mst.m_warehouse.wh_id, 
			mst.m_warehouse.wh_code, 
			mst.m_warehouse.wh_name
		`).
		Where("mst.m_warehouse.cust_id = ?", custID).
		Where("mst.m_warehouse.wh_id = ?", whID).
		Where("mst.m_warehouse.is_active = true").
		Take(&wh).Error

	return wh, err
}

func (repository *RepositoryPromotionV2Impl) CloseExpiredPromotionStatuses(expiredBefore time.Time) (int64, error) {
	result := repository.Model(&model.PromotionV2{}).
		Where("promo.promotions.promo_status IN ?", []model.PromotionStatus{model.PromoStatusActive, model.PromoStatusInactive}).
		Where("promo.promotions.effective_to < ?", expiredBefore).
		Updates(map[string]interface{}{
			"promo_status": model.PromoStatusClosed,
			"updated_at":   time.Now().UTC(),
		})
	return result.RowsAffected, result.Error
}

func (repository *RepositoryPromotionV2Impl) FindActivePromotionsByOutletCriteria(req entity.ConsultPromoV2Req, outlet model.OutletPromo, salesman model.SalesmanPromo) (promotions []model.PromotionV2, err error) {
	query := repository.Select(`DISTINCT promo.promotions.*`).
		Joins("INNER JOIN promo.promotion_outlet_criteria poc ON poc.promo_id = promo.promotions.promo_id AND poc.cust_id = promo.promotions.cust_id").
		Where("promo.promotions.promo_status = ?", model.PromoStatusActive).
		Where("promo.promotions.effective_from <= ? AND promo.promotions.effective_to >= ?", req.OrderDate, req.OrderDate).
		Where("promo.promotions.cust_id = ? OR promo.promotions.distributor_cust_id = ?", req.ParentCustID, req.CustID).
		Where(`(
			-- Match by attribute selection
			(poc.selection_type = ? AND (
				EXISTS (SELECT 1 FROM promo.promotion_outlet_attribute_type poat WHERE poat.criteria_id = poc.id AND poat.outlet_type_id = ?) OR
				EXISTS (SELECT 1 FROM promo.promotion_outlet_attribute_group poag WHERE poag.criteria_id = poc.id AND poag.outlet_group_id = ?) OR
				EXISTS (SELECT 1 FROM promo.promotion_outlet_attribute_class poac WHERE poac.criteria_id = poc.id AND poac.outlet_class_id = ?) OR
				EXISTS (SELECT 1 FROM promo.promotion_outlet_attribute_sales_team poast WHERE poast.criteria_id = poc.id AND poast.sales_team_id = ?)
			)) OR
			-- Match by outlet selection
			(poc.selection_type = ? AND (
				EXISTS (SELECT 1 FROM promo.promotion_outlets_selected pos WHERE pos.criteria_id = poc.id AND pos.outlet_id = ?) OR
				EXISTS (SELECT 1 FROM promo.promotion_outlet_attribute_sales_team poast WHERE poast.criteria_id = poc.id AND poast.sales_team_id = ?)
			))
		)`, "by_attribute", outlet.OtTypeID, outlet.OtGrpID, outlet.OtClassID, salesman.SalesTeamId, "by_outlet", req.OutletID, salesman.SalesTeamId)

	err = query.Find(&promotions).Error
	return promotions, err
}

func (repository *RepositoryPromotionV2Impl) FindProductCriteriasByPromoIDs(promoIDs []string, custID string) (productCriterias []model.PromotionProductCriteria, err error) {
	if len(promoIDs) == 0 {
		return productCriterias, nil
	}
	err = repository.Select(`promo.promotion_product_criteria.*`).
		Where("promo.promotion_product_criteria.promo_id IN (?)", promoIDs).
		Where("promo.promotion_product_criteria.cust_id = ?", custID).
		Order("promo.promotion_product_criteria.promo_id DESC, promo.promotion_product_criteria.pro_id ASC").
		Find(&productCriterias).Error
	return productCriterias, err
}

func (repository *RepositoryPromotionV2Impl) FindSlabsByPromoIDs(promoIDs []string, custID string) (slabs []model.PromotionV2Slabs, err error) {
	if len(promoIDs) == 0 {
		return slabs, nil
	}
	err = repository.Select(`promo.promotion_slabs.*, promo.promotions.promo_desc, promo.promotions.slab_multiplied as is_multiplied`).
		Joins("INNER JOIN promo.promotions ON promo.promotions.promo_id = promo.promotion_slabs.promo_id AND promo.promotions.cust_id = promo.promotion_slabs.cust_id").
		Where("promo.promotion_slabs.promo_id IN (?)", promoIDs).
		Where("promo.promotion_slabs.cust_id = ?", custID).
		Order("promo.promotion_slabs.promo_id DESC, promo.promotion_slabs.ordinal ASC").
		Find(&slabs).Error
	return slabs, err
}

func (repository *RepositoryPromotionV2Impl) FindStratasByPromoIDs(promoIDs []string, custID string) (strata []model.PromotionV2Strata, err error) {
	if len(promoIDs) == 0 {
		return strata, nil
	}
	err = repository.Select(`promo.promotion_strata.*, promo.promotions.promo_desc, promo.promotions.slab_multiplied as is_multiplied`).
		Joins("INNER JOIN promo.promotions ON promo.promotions.promo_id = promo.promotion_strata.promo_id AND promo.promotions.cust_id = promo.promotion_strata.cust_id").
		Where("promo.promotion_strata.promo_id IN (?)", promoIDs).
		Where("promo.promotion_strata.cust_id = ?", custID).
		Order("promo.promotion_strata.promo_id DESC, promo.promotion_strata.ordinal ASC").
		Find(&strata).Error
	return strata, err
}

func (repository *RepositoryPromotionV2Impl) GetAllRewardProductFromStockV2(req entity.ConsultPromoV2Req, ctx model.RewardContext) (rewardProducts []model.PromotionRewardProduct, err error) {
	var qSelect string
	if ctx.RewardUom != nil {
		switch *ctx.RewardUom {
		case model.UomTypeSmallest:
			qSelect = "mp.unit_id1 as unit_id"
		case model.UomTypeMiddle:
			qSelect = "mp.unit_id2 as unit_id"
		default:
			qSelect = "mp.unit_id3 as unit_id"
		}
	} else {
		qSelect = "mp.unit_id3 as unit_id"
	}

	custIDs := []string{req.ParentCustID}
	if req.CustID != "" && req.CustID != req.ParentCustID {
		custIDs = append(custIDs, req.CustID)
	}

	err = repository.Select(`
        promo.promotion_reward_products.*,
        coalesce(s.qty_in, 0) as qty_in,
        coalesce(s.qty_out, 0) as qty_out,
        (coalesce(s.qty_in, 0) - coalesce(s.qty_out, 0)) as qty_stock,
		mp.conv_unit2 as conv_unit2,
		mp.conv_unit3 as conv_unit3,
        `+qSelect+`
    `).
		Joins("INNER JOIN mst.m_product mp ON mp.pro_id = promo.promotion_reward_products.pro_id AND mp.cust_id IN ?", custIDs).
		Joins(`LEFT JOIN (
            SELECT 
                inv.stock.pro_id,
                coalesce(sum(inv.stock.qty_in), 0) as qty_in, 
                coalesce(sum(inv.stock.qty_out), 0) as qty_out 
            FROM inv.stock
            WHERE inv.stock.wh_id = ? AND inv.stock.stock_date <= ? AND inv.stock.cust_id = ?
            GROUP BY inv.stock.pro_id
        ) s ON s.pro_id = promo.promotion_reward_products.pro_id`,
			req.WhID, req.OrderDate, req.CustID,
		).
		Where("promo.promotion_reward_products.promo_id = ?", ctx.PromoID).
		Where("promo.promotion_reward_products.cust_id = ?", req.ParentCustID).
		Where("(coalesce(s.qty_in, 0) - coalesce(s.qty_out, 0)) > 0").
		Order("promo.promotion_reward_products.ordinal ASC").
		Find(&rewardProducts).Error

	return rewardProducts, err
}

func (repository *RepositoryPromotionV2Impl) FindProductByID(productID int64) (product model.ProductRead, err error) {
	err = repository.Select(`
			mst.m_product.pro_id, 
			mst.m_product.pro_code, 
			mst.m_product.pro_name, 
			mst.m_product.conv_unit2, 
			mst.m_product.conv_unit3, 
			mst.m_product.principal_id 
		`).
		Where("mst.m_product.pro_id=?", productID).
		Take(&product).Error

	return product, err
}
