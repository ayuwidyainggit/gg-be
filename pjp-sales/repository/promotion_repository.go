package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"sales/pkg/structs"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"gorm.io/gorm"
)

type (
	RepositoryPromotionImpl struct {
		*gorm.DB
	}
)
type PromotionRepository interface {
	Store(c context.Context, data *model.Promotion) error
	StoreStatusLog(c context.Context, data *model.PromoStatusLog) error
	FindByPromoID(params entity.DetailPromotionParams) (promotion model.Promotion, err error)
	FindAllPromoCriteriasByPromoID(params entity.DetailPromotionParams) (promoCriterias []model.PromoCriteria, err error)
	FindAllPromoAdditionalCriteriasByPromoID(params entity.DetailPromotionParams) (promoAddCriterias []model.PromoAdditionalCriteria, err error)
	FindAllRewardProductsByPromoID(params entity.DetailPromotionParams) (promoRewardProducts []model.PromoRewardProduct, err error)
	FindProductByProID(custID string, proID int64) (rewardProductDetail model.RewardProductDetail, err error)
	FindAllByCustId(dataFilter entity.PromotionQueryFilter) ([]model.Promotion, int64, int, error)
	Update(c context.Context, promoID string, data model.Promotion) error
	Delete(c context.Context, custID, promoID string) error
	DeletePromoCriterias(c context.Context, custID, promoID string) error
	DeletePromoAdditionalCriterias(c context.Context, custId, promoID string) error
	StorePromoCriteria(c context.Context, data *model.PromoCriteria) error
	StorePromoRewardProduct(c context.Context, data *model.PromoRewardProduct) error
	StorePromoAdditionalCriteria(c context.Context, data *model.PromoAdditionalCriteria) error
	DeletePromoCriteriasNotInIDs(c context.Context, custID, promoID string, IDs []int64) error
	DeletePromoAdditionalCriteriasNotInIDs(c context.Context, custID, promoID string, IDs []int64) error
	UpdatePromoCriteria(c context.Context, promoCriteria *model.PromoCriteria) error
	UpdatePromoAdditionalCriteria(c context.Context, promoAddCriteria *model.PromoAdditionalCriteria) error
	DeletePromoRewardProducts(c context.Context, custID, promoID string) error
	FindOneProductByProID(custID, parentCustID string, proID int64) (productAdditionalCriteria model.ProductAdditionalCriteria, err error)
	FindOneOutletClassByID(custID, parentCustID string, id int64) (outletClassAdditionalCriteria model.OutletClassAdditionalCriteria, err error)
	FindOneOutletTypeByID(custID, parentCustID string, id int64) (outletTypeAdditionalCriteria model.OutletTypeAdditionalCriteria, err error)
	FindOneOutletGroupByID(custID, parentCustID string, id int64) (outletGroupAdditionalCriteria model.OutletGroupAdditionalCriteria, err error)
	FindOneSalesTypeByID(custID, parentCustID string, id int64) (salesTypeAdditionalCriteria model.SalesTypeAdditionalCriteria, err error)
	FindOneSalesTeamByID(custID, parentCustID string, id int64) (salesTeamAdditionalCriteria model.SalesTeamAdditionalCriteria, err error)
	FindAllByCustIdAndPromoID(request entity.BulkUpdateStatusPromotionBody) ([]model.Promotion, error)
	BulkUpdateStatus(c context.Context, req entity.BulkUpdateStatusPromotionBody) error

	FindOutletByID(outletID int64, custId string, parentCustId string) (outlet model.OutletRead, err error)
	FindSalesmanByID(salesmanID int64, custId string, parentCustId string) (salesman model.SalesmanRead, err error)
	FindProductByID(productID int64) (product model.ProductRead, err error)
	FindProductByIDAndCustID(productID int64, custID string) (product model.ProductRead, err error)
	FindAllPromoAdditionalCriteriasByActivePromo(request entity.ConsultPromotionBody) (promoAdditionalCriteriasByActivePromo []model.PromoAdditionalCriteriaByActivePromo, err error)
	FindAllPromoCriteriasByPromoIDs(promoIDs []string) (promoCriterias []model.ConsultPromoCriteria, err error)
	FindPromoCriteriasBySlabIDs(slabIDs []string) (promoCriterias []model.ConsultPromoCriteria, err error)
	GetAllRewardProductFromStock(request entity.ConsultPromotionBody, promoCriteria model.ConsultPromoCriteria) (rewardProducts []model.PromoRewardProductRead, err error)
	FindPromoAdditionalCriteriasWithProductAttributeByPromoID(promoID string, custID string) (promoAdditionalCriteriasWithProductAttribute []model.PromoAdditionalCriteriaByActivePromo, err error)
	FindPromoCriteriaBySlabRule(orderReward model.FullPromoRewardRead, custID string, slabRule float64) (promoCriterias model.ConsultPromoCriteria, err error)
	FindProductAndPriceByID(productID, distributorID int64, transDate, custID, parentCustID string) (product model.Product, err error)
}

func NewPromotionRepo(db *gorm.DB) *RepositoryPromotionImpl {
	return &RepositoryPromotionImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryPromotionImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryPromotionImpl) Store(c context.Context, data *model.Promotion) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromotionImpl) StoreStatusLog(c context.Context, data *model.PromoStatusLog) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromotionImpl) FindByPromoID(params entity.DetailPromotionParams) (promotion model.Promotion, err error) {
	err = repository.
		Select(`sls.promotions.*`).
		Where("sls.promotions.promo_id = ? AND sls.promotions.cust_id=?", params.PromoID, params.CustID).
		Take(&promotion).Error
	return promotion, err
}

func (repository *RepositoryPromotionImpl) FindAllByCustId(dataFilter entity.PromotionQueryFilter) ([]model.Promotion, int64, int, error) {
	var (
		promo []model.Promotion
		total int64
	)
	limit := 10
	if dataFilter.Limit != 0 {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("promo_id")
	query := repository.Select(
		`sls.promotions.*`) // .
	// Joins("left join mst.m_outlet ot on ot.outlet_id = sls.promotions.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.promotions.cust_id=?", dataFilter.CustId)
	query.Where("sls.promotions.cust_id=?", dataFilter.CustId)

	if dataFilter.EffectiveFrom != nil && dataFilter.EffectiveTo != nil {
		query.Where(`sls.promotions.effective_from BETWEEN ? AND ? 
					OR sls.promotions.effective_to BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
		)
		queryCount.Where(`sls.promotions.effective_from BETWEEN ? AND ? 
						OR sls.promotions.effective_to BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
		)
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.promotions.promo_id ILIKE ? OR sls.promotions.promo_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query.Where("sls.promotions.promo_id ILIKE ? OR sls.promotions.promo_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if dataFilter.PromoID != "" {
		queryCount.Where("sls.promotions.promo_id ILIKE ?", "%"+dataFilter.PromoID+"%")
		query.Where("sls.promotions.promo_id ILIKE ?", "%"+dataFilter.PromoID+"%")
	}

	if dataFilter.PromoDesc != "" {
		queryCount.Where("sls.promotions.promo_desc ILIKE ?", "%"+dataFilter.PromoDesc+"%")
		query.Where("sls.promotions.promo_desc ILIKE ?", "%"+dataFilter.PromoDesc+"%")
	}

	if len(dataFilter.PromoStatusID) > 0 {
		queryCount.Where("sls.promotions.promo_status_id IN ?", dataFilter.PromoStatusID)
		query.Where("sls.promotions.promo_status_id IN ?", dataFilter.PromoStatusID)
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
		query.Order("created_at DESC")
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

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return promo, total, lastPage, nil
}

func (repository *RepositoryPromotionImpl) Update(c context.Context, promoID string, data model.Promotion) error {
	result := repository.model(c).Model(&data).Where("promo_id=?", promoID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.PromotionwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryPromotionImpl) Delete(c context.Context, custId, promoID string) error {
	var data model.Promotion
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "promo_id": promoID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryPromotionImpl) DeletePromoCriterias(c context.Context, custId, promoID string) error {
	var data model.PromoCriteria
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "promo_id": promoID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryPromotionImpl) DeletePromoAdditionalCriterias(c context.Context, custId, promoID string) error {
	var data model.PromoAdditionalCriteria
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "promo_id": promoID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryPromotionImpl) StorePromoCriteria(c context.Context, data *model.PromoCriteria) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromotionImpl) StorePromoAdditionalCriteria(c context.Context, data *model.PromoAdditionalCriteria) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromotionImpl) FindAllPromoCriteriasByPromoID(params entity.DetailPromotionParams) (promoCriterias []model.PromoCriteria, err error) {
	err = repository.
		Select(`promo_criterias.*`).
		// Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Where("promo_criterias.promo_id = ? AND promo_criterias.cust_id = ?", params.PromoID, params.CustID).
		Order("promo_criterias.slab_id ASC").
		Find(&promoCriterias).Error
	return promoCriterias, err
}

func (repository *RepositoryPromotionImpl) FindAllPromoAdditionalCriteriasByPromoID(params entity.DetailPromotionParams) (promoAddCriterias []model.PromoAdditionalCriteria, err error) {
	err = repository.
		Select(`promo_additional_criterias.*`).
		Where("promo_additional_criterias.promo_id = ? AND promo_additional_criterias.cust_id = ?", params.PromoID, params.CustID).
		Order("promo_additional_criterias.promo_add_criteria_id ASC").
		Find(&promoAddCriterias).Error
	return promoAddCriterias, err
}

func (repository *RepositoryPromotionImpl) FindAllRewardProductsByPromoID(params entity.DetailPromotionParams) (promoRewardProducts []model.PromoRewardProduct, err error) {
	err = repository.
		Select(`promo_reward_products.*`).
		Where("promo_reward_products.promo_id = ? AND promo_reward_products.cust_id = ?", params.PromoID, params.CustID).
		Order("promo_reward_products.promo_reward_id ASC").
		Find(&promoRewardProducts).Error
	return promoRewardProducts, err
}

func (repository *RepositoryPromotionImpl) FindProductByProID(custID string, proID int64) (rewardProductDetail model.RewardProductDetail, err error) {
	err = repository.
		Select(`m_product.*`).
		Where("m_product.cust_id = ? AND m_product.pro_id = ? AND m_product.is_active = ?", custID, proID, true).
		Take(&rewardProductDetail).Error
	return rewardProductDetail, err
}

func (repository *RepositoryPromotionImpl) DeletePromoCriteriasNotInIDs(c context.Context, custID, promoID string, IDs []int64) error {
	var promoCriterias model.PromoCriteria
	err := repository.model(c).Where("cust_id = ? AND promo_id = ? AND slab_id NOT IN (?) ", custID, promoID, IDs).Delete(&promoCriterias).Error
	return err
}

func (repository *RepositoryPromotionImpl) DeletePromoAdditionalCriteriasNotInIDs(c context.Context, custID, promoID string, IDs []int64) error {
	var promoAddCriterias model.PromoAdditionalCriteria
	err := repository.model(c).Where("cust_id = ? AND promo_id = ? AND promo_add_criteria_id NOT IN (?) ", custID, promoID, IDs).Delete(&promoAddCriterias).Error
	return err
}

func (repository *RepositoryPromotionImpl) UpdatePromoCriteria(c context.Context, promoCriteria *model.PromoCriteria) error {
	result := repository.model(c).Updates(&promoCriteria)
	if result.Error != nil {
		log.Error("UpdatePromoCriteria, result:", structs.StructToJson(result.Error))
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryPromotionImpl) UpdatePromoAdditionalCriteria(c context.Context, promoAddCriteria *model.PromoAdditionalCriteria) error {
	result := repository.model(c).Updates(&promoAddCriteria)
	if result.Error != nil {
		log.Error("UpdatePromoAdditionalCriteria, result:", structs.StructToJson(result.Error))
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryPromotionImpl) StorePromoRewardProduct(c context.Context, data *model.PromoRewardProduct) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromotionImpl) DeletePromoRewardProducts(c context.Context, custID, promoID string) error {
	var rewardProduct model.PromoRewardProduct
	err := repository.model(c).Where("cust_id = ? AND promo_id = ?", custID, promoID).Delete(&rewardProduct).Error
	return err
}

func (repository *RepositoryPromotionImpl) FindOneProductByProID(custID, parentCustID string, proID int64) (productAddCriteria model.ProductAdditionalCriteria, err error) {
	err = repository.
		Select(`m_product.pro_id AS reference_id, m_product.pro_code AS reference_code, m_product.pro_name AS reference_name`).
		Where("m_product.cust_id = ? AND m_product.pro_id = ?", parentCustID, proID).
		Take(&productAddCriteria).Error
	return productAddCriteria, err
}

func (repository *RepositoryPromotionImpl) FindOneOutletClassByID(custID, parentCustID string, outletClassID int64) (outletClassAddCriteria model.OutletClassAdditionalCriteria, err error) {
	err = repository.
		Select(`m_outlet_class.ot_class_id AS reference_id, m_outlet_class.ot_class_code AS reference_code, m_outlet_class.ot_class_name AS reference_name`).
		Where("m_outlet_class.cust_id = ? AND m_outlet_class.ot_class_id = ?", parentCustID, outletClassID).
		Take(&outletClassAddCriteria).Error
	return outletClassAddCriteria, err
}

func (repository *RepositoryPromotionImpl) FindOneOutletTypeByID(custID, parentCustID string, outletTypeID int64) (outletTypeAddCriteria model.OutletTypeAdditionalCriteria, err error) {
	err = repository.
		Select(`m_outlet_type.ot_type_id AS reference_id, m_outlet_type.ot_type_code AS reference_code, m_outlet_type.ot_type_name AS reference_name`).
		Where("m_outlet_type.cust_id = ? AND m_outlet_type.ot_type_id = ?", parentCustID, outletTypeID).
		Take(&outletTypeAddCriteria).Error
	return outletTypeAddCriteria, err
}

func (repository *RepositoryPromotionImpl) FindOneOutletGroupByID(custID, parentCustID string, outletGroupID int64) (outletGroupAddCriteria model.OutletGroupAdditionalCriteria, err error) {
	err = repository.
		Select(`m_outlet_group.ot_grp_id AS reference_id, m_outlet_group.ot_grp_code AS reference_code, m_outlet_group.ot_grp_name AS reference_name`).
		Where("m_outlet_group.cust_id = ? AND m_outlet_group.ot_grp_id = ?", parentCustID, outletGroupID).
		Take(&outletGroupAddCriteria).Error
	return outletGroupAddCriteria, err
}

func (repository *RepositoryPromotionImpl) FindOneSalesTypeByID(custID, parentCustID string, salesTypeID int64) (salesTypeAddCriteria model.SalesTypeAdditionalCriteria, err error) {
	err = repository.
		Select(`m_sales_type.sales_type_id AS reference_id, m_sales_type.sales_type_code AS reference_code, m_sales_type.sales_type_name AS reference_name`).
		Where("m_sales_type.cust_id = ? AND m_sales_type.sales_type_id = ?", custID, salesTypeID).
		Take(&salesTypeAddCriteria).Error
	return salesTypeAddCriteria, err
}

func (repository *RepositoryPromotionImpl) FindOneSalesTeamByID(custID, parentCustID string, salesTeamID int64) (salesTeamAddCriteria model.SalesTeamAdditionalCriteria, err error) {
	err = repository.
		Select(`m_sales_team.sales_team_id AS reference_id, m_sales_team.sales_team_code AS reference_code, m_sales_team.sales_team_name AS reference_name`).
		Where("m_sales_team.cust_id = ? AND m_sales_team.sales_team_id = ?", custID, salesTeamID).
		Take(&salesTeamAddCriteria).Error
	return salesTeamAddCriteria, err
}

func (repository *RepositoryPromotionImpl) FindAllByCustIdAndPromoID(request entity.BulkUpdateStatusPromotionBody) ([]model.Promotion, error) {
	var (
		promo []model.Promotion
	)

	query := repository.Select(`sls.promotions.promo_id, sls.promotions.promo_status_id`)
	query.Where("sls.promotions.cust_id = ? AND promo_id IN ?", request.CustID, request.PromoID)

	err := query.Find(&promo).Error
	if err != nil {
		return promo, err
	}

	return promo, nil
}

func (repository *RepositoryPromotionImpl) BulkUpdateStatus(c context.Context, request entity.BulkUpdateStatusPromotionBody) error {
	promoStatusID := model.Promotion{
		PromoStatusID: request.PromoStatusID,
		Remarks:       request.Remarks,
		UpdatedBy:     request.UpdatedBy,
	}
	promoModel := model.Promotion{}
	result := repository.model(c).Model(&promoModel).
		Where(`cust_id = ? AND promo_id IN ?`, request.CustID, request.PromoID).
		Updates(promoStatusID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != int64(len(request.PromoID)) {
		return errors.New("rows affected is different with the request")
	}
	return nil
}

func (repository *RepositoryPromotionImpl) FindOutletByID(outletID int64, custId string, parentCustId string) (outlet model.OutletRead, err error) {
	err = repository.Select(`
			mst.m_outlet.outlet_id, 
			mst.m_outlet.outlet_code, 
			mst.m_outlet.outlet_name, 
			mst.m_outlet.ot_grp_id, 
			mst.m_outlet.ot_class_id, 
			mst.m_outlet.ot_type_id 
		`).
		Where("mst.m_outlet.outlet_id=?", outletID).
		Take(&outlet).Error

	return outlet, err
}

func (repository *RepositoryPromotionImpl) FindSalesmanByID(salesmanId int64, custId string, parentCustId string) (salesman model.SalesmanRead, err error) {
	err = repository.Select(`
			mst.m_salesman.emp_id as salesman_id, 
			mst.m_salesman.sales_name as salesman_name, 
			mst.m_salesman.sales_team_id,
			mst.m_salesman.wh_id
		`).
		Where("mst.m_salesman.emp_id=?", salesmanId).
		Take(&salesman).Error

	return salesman, err
}

func (repository *RepositoryPromotionImpl) FindProductByID(productID int64) (product model.ProductRead, err error) {
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

func (repository *RepositoryPromotionImpl) FindAllPromoAdditionalCriteriasByActivePromo(request entity.ConsultPromotionBody) (promoAdditionalCriteriasByActivePromo []model.PromoAdditionalCriteriaByActivePromo, err error) {
	err = repository.Select(`
			promotions.promo_desc,
			promotions.is_multiplied,
			sls.promo_additional_criterias.*
		`).
		// Joins("inner join sls.promotions promotions ON promotions.promo_id = sls.promo_additional_criterias.promo_id AND promotions.cust_id = ? AND promotions.promo_status_id = 6", request.CustID).
		Joins("inner join sls.promotions promotions ON promotions.promo_id = sls.promo_additional_criterias.promo_id AND promotions.cust_id = ? AND promotions.promo_status_id = 6 AND promotions.effective_from <= ? AND promotions.effective_to >= ?", request.CustID, request.OrderDate, request.OrderDate).
		// Joins("inner join sls.promotions promotions ON promotions.promo_id = sls.promo_additional_criterias.promo_id AND promotions.cust_id = ? AND promotions.promo_status_id IN (5, 6, 7) AND promotions.effective_from <= ? AND promotions.effective_to >= ?", request.CustID, request.OrderDate, request.OrderDate).
		Where("sls.promo_additional_criterias.attribute IN ('OCL', 'OTG', 'OTY', 'STE', 'PRO')").
		Order("sls.promo_additional_criterias.promo_id DESC, sls.promo_additional_criterias.attribute ASC, sls.promo_additional_criterias.is_mandatory DESC, sls.promo_additional_criterias.reference_id ASC").
		Find(&promoAdditionalCriteriasByActivePromo).Error

	return promoAdditionalCriteriasByActivePromo, err
}

func (repository *RepositoryPromotionImpl) FindPromoAdditionalCriteriasWithProductAttributeByPromoID(promoID string, custID string) (promoAdditionalCriteriasWithProductAttribute []model.PromoAdditionalCriteriaByActivePromo, err error) {
	err = repository.Select(`
			promotions.promo_desc,
			promotions.is_multiplied,
			sls.promo_additional_criterias.*
		`).
		Joins("inner join sls.promotions promotions ON promotions.promo_id = sls.promo_additional_criterias.promo_id AND promotions.cust_id = ? AND promotions.promo_id = ?", custID, promoID).
		Where("sls.promo_additional_criterias.attribute = 'PRO'").
		Order("sls.promo_additional_criterias.is_mandatory DESC, sls.promo_additional_criterias.reference_id ASC").
		Find(&promoAdditionalCriteriasWithProductAttribute).Error

	return promoAdditionalCriteriasWithProductAttribute, err
}

func (repository *RepositoryPromotionImpl) FindAllPromoCriteriasByPromoIDs(promoIDs []string) (promoCriterias []model.ConsultPromoCriteria, err error) {
	if len(promoIDs) == 0 {
		return promoCriterias, nil
	}

	err = repository.Select(`
			sls.promo_criterias.*,
			sls.promo_criterias.slab_rule_type AS slab_rule,
			promotions.promo_desc,
			promotions.is_multiplied
		`).
		Joins("inner join sls.promotions promotions ON promotions.promo_id = sls.promo_criterias.promo_id").
		Where("sls.promo_criterias.promo_id IN ('" + strings.Join(promoIDs, "', '") + "')").
		Order("sls.promo_criterias.promo_id DESC, sls.promo_criterias.slab_id ASC").
		Find(&promoCriterias).Error

	return promoCriterias, err
}

func (repository *RepositoryPromotionImpl) FindPromoCriteriasBySlabIDs(slabIDs []string) (promoCriterias []model.ConsultPromoCriteria, err error) {
	if len(slabIDs) == 0 {
		return promoCriterias, nil
	}

	err = repository.Select(`
			sls.promo_criterias.*,
			sls.promo_criterias.slab_rule_type AS slab_rule,
			promotions.promo_desc,
			promotions.is_multiplied
		`).
		Joins("inner join sls.promotions promotions ON promotions.promo_id = sls.promo_criterias.promo_id").
		Where("sls.promo_criterias.slab_id IN (" + strings.Join(slabIDs, ", ") + ")").
		Order("sls.promo_criterias.slab_reward_uom ASC, sls.promo_criterias.slab_reward ASC").
		Find(&promoCriterias).Error

	return promoCriterias, err
}

func (repository *RepositoryPromotionImpl) FindPromoCriteriaBySlabRule(orderReward model.FullPromoRewardRead, custID string, slabRule float64) (promoCriteria model.ConsultPromoCriteria, err error) {
	// if len(slabIDs) == 0 {
	// 	return promoCriterias, nil
	// }
	var query string
	if orderReward.IsMultiplied {
		query = "sls.promo_criterias.slab_rule_to <= " + strconv.FormatInt(int64(slabRule), 10)
	} else {
		query = "sls.promo_criterias.slab_rule_from <= " + strconv.FormatInt(int64(slabRule), 10) + "AND sls.promo_criterias.slab_rule_to >= " + strconv.FormatInt(int64(slabRule), 10)
	}

	err = repository.Select(`
			sls.promo_criterias.*,
			sls.promo_criterias.slab_rule_type AS slab_rule,
			promotions.promo_desc,
			promotions.is_multiplied
		`).
		Joins("inner join sls.promotions promotions ON promotions.promo_id = sls.promo_criterias.promo_id").
		Where("sls.promo_criterias.promo_id = ?", orderReward.PromoID).
		Where("sls.promo_criterias.cust_id = ?", custID).
		Where(query).
		Take(&promoCriteria).Error

	return promoCriteria, err
}

func (repository *RepositoryPromotionImpl) GetAllRewardProductFromStock(request entity.ConsultPromotionBody, promoCriteria model.ConsultPromoCriteria) (rewardProducts []model.PromoRewardProductRead, err error) {

	var qSelect string
	// var qSlabValue string
	switch promoCriteria.SlabRewardUom {
	case 3:
		qSelect = "mp.unit_id3 as unit_id"
		// qSlabValue = "(pc.slab_reward * mp.conv_unit3 * mp.conv_unit2)"
	case 2:
		qSelect = "mp.unit_id2 as unit_id"
		// qSlabValue = "(pc.slab_reward * mp.conv_unit2)"
	default:
		qSelect = "mp.unit_id1 as unit_id"
		// qSlabValue = "pc.slab_reward"
	}

	// if promoCriteria.IsMultiplied {
	// 	qSlabValue = qSlabValue + " * " + strconv.FormatInt(promoCriteria.SlabRule/int64(promoCriteria.SlabRuleTo), 10)
	// }

	err = repository.Select(`
			sls.promo_reward_products.*,
			pc.slab_reward_type,
			pc.slab_reward,
			pc.slab_reward_uom,	
			coalesce(s.qty_in, 0) as qty_in,
			coalesce(s.qty_out, 0) as qty_out,
			(coalesce(s.qty_in, 0) - coalesce(s.qty_out, 0)) as qty_stock,
			`+qSelect+`
		`).
		Joins("inner join sls.promo_criterias pc on pc.promo_id = sls.promo_reward_products.promo_id and pc.slab_id = ? and pc.cust_id = ?", promoCriteria.SlabID, request.CustID).
		Joins("inner join mst.m_product mp on mp.pro_id = sls.promo_reward_products.pro_id and mp.cust_id = ?", request.ParentCustID).
		Joins(`left join (
			select 
				inv.stock.pro_id,
				coalesce(sum(inv.stock.qty_in), 0) as qty_in, 
				coalesce(sum(inv.stock.qty_out), 0) as qty_out 
			from inv.stock
			where inv.stock.wh_id = `+strconv.Itoa(request.WhId)+` and inv.stock.stock_date <= '`+request.OrderDate+`' and inv.stock.cust_id = '`+request.CustID+`'
			group by inv.stock.pro_id
		) s on s.pro_id = sls.promo_reward_products.pro_id`).
		Where("sls.promo_reward_products.promo_id = ?", promoCriteria.PromoID).
		Where("(coalesce(s.qty_in, 0) - coalesce(s.qty_out, 0)) > 0").
		// Where("(coalesce(s.qty_in, 0) - coalesce(s.qty_out, 0)) > " + qSlabValue).
		Order("sls.promo_reward_products.promo_reward_id ASC").
		Find(&rewardProducts).Error

	return rewardProducts, err
}

func (repository *RepositoryPromotionImpl) FindProductByIDAndCustID(productID int64, custID string) (product model.ProductRead, err error) {
	err = repository.Select(`
			mst.m_product.pro_id, 
			mst.m_product.pro_code, 
			mst.m_product.pro_name, 
			mst.m_product.conv_unit2, 
			mst.m_product.conv_unit3
		`).
		Where("mst.m_product.cust_id=? AND mst.m_product.pro_id=?", custID, productID).
		Take(&product).Error

	return product, err
}

func (repository *RepositoryPromotionImpl) FindProductAndPriceByID(productID, distributorID int64, transDate, custID, parentCustID string) (model.Product, error) {
	query := `SELECT 
				p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name,
				p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
				p.pcat_id, pc.pcat_code, pc.pcat_name,
				br.brand_id, br.brand_code, br.brand_name, 
				br.pl_id, pl.pl_code, pl.pl_name,
				p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
				p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
				p.flavor_id, fv.flavor_code, fv.flavor_name, 
				p.ptype_id, pt.ptype_code, pt.ptype_name,
				p.psize_id, ps.psize_code, ps.psize_name,
				p.sup_id, su.sup_code, su.sup_name,
				p.principal_id, pr.principal_code, pr.principal_name,
				p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
				p.is_main_pro, p.sort_no, p.item_no, p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
				un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3,
				unc1.unit_id_coretax AS unit_id_coretax1, unc2.unit_id_coretax AS unit_id_coretax2, unc3.unit_id_coretax AS unit_id_coretax3,
				unc1.unit_name_coretax AS unit_name_coretax1, unc2.unit_name_coretax AS unit_name_coretax2, unc3.unit_name_coretax AS unit_name_coretax3,
				p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
				p.is_batch, p.is_exp_date, 
				p.weight,p.length, p.width, p.height, p.volume,
				
				CASE WHEN mtp_mg_pr.purch_price1 IS NULL THEN p.purch_price1 
					ELSE mtp_mg_pr.purch_price1 END AS purch_price1,
					
				CASE WHEN mtp_mg_pr.purch_price2 IS NULL THEN p.purch_price2
					ELSE mtp_mg_pr.purch_price2 END AS purch_price2,
					
				CASE WHEN mtp_mg_pr.purch_price3 IS NULL THEN p.purch_price3 
					ELSE mtp_mg_pr.purch_price3 END AS purch_price3,

				p.purch_price4, p.purch_price5,

				CASE WHEN (mtp.sell_price1=0 or mtp.sell_price1 is null) THEN 
				(CASE WHEN (mtp_mg_pr_sell.sell_price1=0 or mtp_mg_pr_sell.sell_price1 is null) THEN p.sell_price1 
					ELSE mtp_mg_pr_sell.sell_price1 END)
				ELSE mtp.sell_price1 END AS sell_price1,

				CASE WHEN (mtp.sell_price2=0 or mtp.sell_price2 is null) THEN 
				(CASE WHEN (mtp_mg_pr_sell.sell_price2=0 or mtp_mg_pr_sell.sell_price2 is null) THEN p.sell_price2 
					ELSE mtp_mg_pr_sell.sell_price2 END)
				ELSE mtp.sell_price2 END AS sell_price2,

				CASE WHEN (mtp.sell_price3=0 or mtp.sell_price3 is null) THEN 
				(CASE WHEN (mtp_mg_pr_sell.sell_price3=0 or mtp_mg_pr_sell.sell_price3 is null) THEN p.sell_price3 
					ELSE mtp_mg_pr_sell.sell_price3 END)
				ELSE mtp.sell_price3 END AS sell_price3,
				p.sell_price4, p.sell_price5, 

				p.weight1, p.length1, p.width1, p.height1, p.volume1, 
				p.weight2, p.length2, p.width2, p.height2, p.volume2, 
				p.weight3, p.length3, p.width3, p.height3, p.volume3, 
				p.weight4, p.length4, p.width4, p.height4, p.volume4, 
				p.weight5, p.length5, p.width5, p.height5, p.volume5,  
				p.parent_pro_id, 
				p.excise_rate, p.excise_tax, 
				p.is_active, p.created_by, p.created_at, 
				p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
				p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
				saf_unit.unit_name AS saf_stock_unit_name,
				min_unit.unit_name AS min_stock_unit_name,
				u.user_fullname AS updated_by_name,
				parent.pro_code AS parent_pro_code, parent.pro_name AS parent_pro_name,p.pro_code_coretax,pct.pro_name_coretax
			FROM mst.m_product p
			LEFT JOIN sys.m_user u ON u.user_id = p.updated_by
			LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + parentCustID + `' 
			LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + parentCustID + `' 
			LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + parentCustID + `' 
			LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc1 ON un1.unit_id_coretax = unc1.unit_id_coretax AND un1.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc2 ON un2.unit_id_coretax = unc2.unit_id_coretax AND un2.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc3 ON un3.unit_id_coretax = unc3.unit_id_coretax AND un3.cust_id = '` + parentCustID + `'
			LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id
			LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id
			LEFT JOIN mst.m_product parent ON parent.pro_id = p.parent_pro_id
			LEFT JOIN mst.m_product_coretax pct ON pct.pro_code_coretax = p.pro_code_coretax
			LEFT JOIN mst.m_distributor dist ON dist.cust_id = '` + parentCustID + `' AND dist.distributor_id = ` + fmt.Sprintf("%d", distributorID) + `
			LEFT JOIN mst.m_transaction_price mtp ON mtp.pro_id = p.pro_id 
					AND mtp.cust_id = '` + custID + `' AND mtp.outlet_id = 0
					AND ('` + transDate + `' BETWEEN mtp.start_date AND mtp.end_date) 
			LEFT JOIN LATERAL (
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3
				FROM mst.m_transaction_price mtp_mg_pr
				WHERE mtp_mg_pr.cust_id = '` + parentCustID + `' 
					AND mtp_mg_pr.pro_id = p.pro_id
					AND mtp_mg_pr.start_date <= '` + transDate + `' 
					AND (mtp_mg_pr.distributor_id = (CASE WHEN mtp_mg_pr.coverage='N' THEN 0 ELSE dist.distributor_id END)
								OR mtp_mg_pr.price_group_reff = dist.dist_price_grp_id)	
				ORDER BY mtp_mg_pr.start_date DESC LIMIT 1
			) mtp_mg_pr ON true
			LEFT JOIN LATERAL (
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3 
				FROM mst.m_transaction_price mtp_mg_pr_sell
				WHERE mtp_mg_pr_sell.cust_id = '` + parentCustID + `' 
					AND mtp_mg_pr_sell.pro_id = p.pro_id
					and mtp_mg_pr_sell.source = 10
					AND mtp_mg_pr_sell.start_date <= '` + transDate + `' 
					AND (mtp_mg_pr_sell.distributor_id = (CASE WHEN mtp_mg_pr_sell.coverage = 'N' THEN 0 ELSE dist.distributor_id END)
								OR mtp_mg_pr_sell.price_group_reff = dist.dist_price_grp_id)	
				ORDER BY mtp_mg_pr_sell.start_date DESC LIMIT 1
			) mtp_mg_pr_sell ON true
			WHERE p.pro_id = $1 
			AND p.cust_id IN ($2, $3)`

	var product model.Product
	err := repository.Raw(query, productID, parentCustID, custID).Scan(&product).Error
	return product, err
}
