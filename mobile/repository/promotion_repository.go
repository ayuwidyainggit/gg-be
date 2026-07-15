package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/str"
	"mobile/pkg/structs"
	"strconv"
	"strings"
	"time"

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
	FindAllPromoAdditionalCriteriasByActivePromo(request entity.ConsultPromotionBody) (promoAdditionalCriteriasByActivePromo []model.PromoAdditionalCriteriaByActivePromo, err error)
	FindAllPromoCriteriasByPromoIDs(promoIDs []string) (promoCriterias []model.ConsultPromoCriteria, err error)
	GetAllRewardProductFromStock(request entity.ConsultPromotionBody, promoCriteria model.ConsultPromoCriteria) (rewardProducts []model.PromoRewardProductRead, err error)

	// Mobile Promotion List
	FindAllActivePromotions(dataFilter entity.PromotionMobileListQueryFilter) ([]model.PromotionMobileList, int64, int, error)
	FindPromotedProductsByPromoIDs(promoIDs []string, parentCustID string) ([]model.PromotedProductRead, error)
	FindRewardProductsByPromoIDs(promoIDs []string, parentCustID string) ([]model.RewardProductRead, error)
	FindAdditionalCriteriasByPromoIDs(promoIDs []string) ([]model.AdditionalCriteriaRead, error)
	FindPromoCriteriaMinMaxByPromoIDs(promoIDs []string) ([]model.PromoCriteriaMinMax, error)

	// Mobile Promotion Detail
	FindPromotionDetailByPromoID(promoID, parentCustID string) (model.PromotionMobileDetail, error)
	FindPromotedProductsByPromoID(promoID, parentCustID string) ([]model.PromotedProductDetailRead, error)
	FindPromotionCriteriaByPromoID(promoID string) ([]model.PromotionCriteriaDetailRead, error)
	FindPromotionRewardByPromoID(promoID, parentCustID string) ([]model.PromotionRewardDetailRead, error)
	FindPromotionSlabByPromoID(promoID string) ([]model.PromotionSlabRead, error)
	FindPromotionStrataByPromoID(promoID string) ([]model.PromotionStrataRead, error)
	FindOutletTypesByPromoID(promoID, parentCustID string) ([]model.OutletTypeDetailRead, error)
	FindOutletGroupsByPromoID(promoID, parentCustID string) ([]model.OutletGroupDetailRead, error)
	FindOutletClassesByPromoID(promoID, parentCustID string) ([]model.OutletClassDetailRead, error)

	// Promotion Outlet List
	FindOutletsByTypeGroupClass(dataFilter entity.PromotionOutletListQueryFilter, otTypeID int64, custId, parentCustId string) ([]model.PromotionOutletList, int64, int, error)
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

func (repository *RepositoryPromotionImpl) GetAllRewardProductFromStock(request entity.ConsultPromotionBody, promoCriteria model.ConsultPromoCriteria) (rewardProducts []model.PromoRewardProductRead, err error) {

	var qSelect string
	// var qSlabValue string
	switch promoCriteria.SlabRewardUom {
	case 3:
		qSelect = "mp.unit_id3 as unit_id"
		// qSlabValue = "(pc.slab_reward * mp.conv_unit3 * mp.conv_unit2)"
	case 2:
		qSelect = "mp.unit_id3 as unit_id"
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

// FindAllActivePromotions finds all active promotions (promo_status = 'active') for mobile
func (repository *RepositoryPromotionImpl) FindAllActivePromotions(dataFilter entity.PromotionMobileListQueryFilter) ([]model.PromotionMobileList, int64, int, error) {
	var (
		promotions []model.PromotionMobileList
		total      int64
		sortBy     string
	)

	limit := 20
	if dataFilter.Limit != 0 {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("promo_id")
	query := repository.Select(`
		promo.promotions.cust_id,
		promo.promotions.promo_id,
		promo.promotions.promo_desc,
		promo.promotions.effective_from,
		promo.promotions.effective_to,
		promo.promotions.slab_multiplied,
		promo.promotions.max_invoice_per_outlet
	`)

	// Filter by cust_id (use ParentCustId for master data like promotions)
	queryCount.Where("promo.promotions.cust_id = ?", dataFilter.ParentCustId)
	query.Where("promo.promotions.cust_id = ?", dataFilter.ParentCustId)

	// Filter by promo_status = 'active'
	queryCount.Where("promo.promotions.promo_status = ?", "active")
	query.Where("promo.promotions.promo_status = ?", "active")

	// Filter by promo_desc if provided
	if dataFilter.PromoDesc != "" {
		queryCount.Where("promo.promotions.promo_desc ILIKE ?", "%"+dataFilter.PromoDesc+"%")
		query.Where("promo.promotions.promo_desc ILIKE ?", "%"+dataFilter.PromoDesc+"%")
	}

	// Sorting
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
		query.Order("promo_desc ASC")
	}

	// Pagination
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&promotions).Error
	if err != nil {
		return promotions, total, 0, err
	}

	err = queryCount.Model(&model.PromotionMobileList{}).Count(&total).Error
	if err != nil {
		return promotions, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return promotions, total, lastPage, nil
}

// FindPromotedProductsByPromoIDs finds promoted products by promo IDs (attribute = 'PRO')
func (repository *RepositoryPromotionImpl) FindPromotedProductsByPromoIDs(promoIDs []string, parentCustID string) ([]model.PromotedProductRead, error) {
	var promotedProducts []model.PromotedProductRead

	if len(promoIDs) == 0 {
		return promotedProducts, nil
	}

	err := repository.Select(`
		promo.promo_additional_criterias.promo_id,
		promo.promo_additional_criterias.reference_id AS pro_id,
		mp.pro_code,
		mp.pro_name,
		promo.promo_additional_criterias.is_mandatory,
		promo.promo_additional_criterias.min_buy_value
	`).
		Joins("INNER JOIN mst.m_product mp ON mp.pro_id = promo.promo_additional_criterias.reference_id AND mp.cust_id = ?", parentCustID).
		Where("promo.promo_additional_criterias.promo_id IN ('"+strings.Join(promoIDs, "', '")+"')").
		Where("promo.promo_additional_criterias.attribute = ?", "PRO").
		Order("promo.promo_additional_criterias.promo_id, promo.promo_additional_criterias.promo_add_criteria_id ASC").
		Find(&promotedProducts).Error

	return promotedProducts, err
}

// FindRewardProductsByPromoIDs finds reward products by promo IDs
func (repository *RepositoryPromotionImpl) FindRewardProductsByPromoIDs(promoIDs []string, parentCustID string) ([]model.RewardProductRead, error) {
	var rewardProducts []model.RewardProductRead

	if len(promoIDs) == 0 {
		return rewardProducts, nil
	}

	err := repository.Select(`
		promo.promo_reward_products.promo_id,
		promo.promo_reward_products.pro_id,
		mp.pro_code,
		mp.pro_name
	`).
		Joins("INNER JOIN mst.m_product mp ON mp.pro_id = promo.promo_reward_products.pro_id AND mp.cust_id = ?", parentCustID).
		Where("promo.promo_reward_products.promo_id IN ('" + strings.Join(promoIDs, "', '") + "')").
		Order("promo.promo_reward_products.promo_id, promo.promo_reward_products.promo_reward_id ASC").
		Find(&rewardProducts).Error

	return rewardProducts, err
}

// FindAdditionalCriteriasByPromoIDs finds additional criterias by promo IDs (excluding PRO attribute)
func (repository *RepositoryPromotionImpl) FindAdditionalCriteriasByPromoIDs(promoIDs []string) ([]model.AdditionalCriteriaRead, error) {
	var additionalCriterias []model.AdditionalCriteriaRead

	if len(promoIDs) == 0 {
		return additionalCriterias, nil
	}

	// Map attribute codes to names
	attributeNameCase := `
		CASE 
			WHEN promo.promo_additional_criterias.attribute = 'OCL' THEN 'Outlet Class'
			WHEN promo.promo_additional_criterias.attribute = 'OTG' THEN 'Outlet Group'
			WHEN promo.promo_additional_criterias.attribute = 'OTY' THEN 'Outlet Type'
			WHEN promo.promo_additional_criterias.attribute = 'STE' THEN 'Sales Team'
			WHEN promo.promo_additional_criterias.attribute = 'STY' THEN 'Sales Type'
			ELSE promo.promo_additional_criterias.attribute
		END AS attribute_name
	`

	err := repository.Select(`
		promo.promo_additional_criterias.promo_id,
		promo.promo_additional_criterias.promo_add_criteria_id,
		promo.promo_additional_criterias.attribute,
		`+attributeNameCase).
		Where("promo.promo_additional_criterias.promo_id IN ('"+strings.Join(promoIDs, "', '")+"')").
		Where("promo.promo_additional_criterias.attribute != ?", "PRO").
		Order("promo.promo_additional_criterias.promo_id, promo.promo_additional_criterias.promo_add_criteria_id ASC").
		Find(&additionalCriterias).Error

	return additionalCriterias, err
}

// FindPromoCriteriaMinMaxByPromoIDs finds min and max purchase from promo criterias
func (repository *RepositoryPromotionImpl) FindPromoCriteriaMinMaxByPromoIDs(promoIDs []string) ([]model.PromoCriteriaMinMax, error) {
	var criteriaMinMax []model.PromoCriteriaMinMax

	if len(promoIDs) == 0 {
		return criteriaMinMax, nil
	}

	err := repository.Select(`
		promo.promo_criterias.promo_id,
		MIN(promo.promo_criterias.slab_rule_from) AS min_purchase,
		MAX(promo.promo_criterias.slab_rule_to) AS max_purchase
	`).
		Where("promo.promo_criterias.promo_id IN ('" + strings.Join(promoIDs, "', '") + "')").
		Group("promo.promo_criterias.promo_id").
		Find(&criteriaMinMax).Error

	return criteriaMinMax, err
}

// FindPromotionDetailByPromoID finds promotion detail by promo_id
func (repository *RepositoryPromotionImpl) FindPromotionDetailByPromoID(promoID, parentCustID string) (model.PromotionMobileDetail, error) {
	var promotion model.PromotionMobileDetail

	err := repository.Select(`
		promo.promotions.promo_id,
		promo.promotions.promo_desc,
		promo.promotions.effective_from,
		promo.promotions.effective_to,
		promo.promotions.slab_multiplied,
		promo.promotions.max_invoice_per_outlet,
		promo.promotions.promo_type,
		promo.promotions.max_total_reward_type,
		promo.promotions.max_total_reward_value,
		COALESCE(promo.promotion_slabs.rule_uom::text, promo.promotion_strata.rule_uom::text, NULL) AS rule_uom
	`).
		Table("promo.promotions").
		Joins("LEFT JOIN promo.promotion_slabs ON promo.promotion_slabs.promo_id = promo.promotions.promo_id AND promo.promotions.promo_type = 'slab'").
		Joins("LEFT JOIN promo.promotion_strata ON promo.promotion_strata.promo_id = promo.promotions.promo_id AND promo.promotions.promo_type = 'strata'").
		Where("promo.promotions.promo_id = ? AND promo.promotions.cust_id = ?", promoID, parentCustID).
		First(&promotion).Error

	return promotion, err
}

// FindPromotedProductsByPromoID finds promoted products by promo_id
func (repository *RepositoryPromotionImpl) FindPromotedProductsByPromoID(promoID, parentCustID string) ([]model.PromotedProductDetailRead, error) {
	var products []model.PromotedProductDetailRead

	err := repository.Select(`
		ppc.promo_id,
		ppc.pro_id,
		mp.pro_code,
		mp.pro_name,
		ppc.mandatory,
		ppc.min_buy_type,
		ppc.min_buy_qty,
		ppc.min_buy_value,
		ppc.min_buy_uom
	`).
		Table("promo.promotion_product_criteria AS ppc").
		Joins("INNER JOIN mst.m_product mp ON mp.pro_id = ppc.pro_id AND mp.cust_id = ?", parentCustID).
		Where("ppc.promo_id = ?", promoID).
		Find(&products).Error

	return products, err
}

// FindPromotionCriteriaByPromoID finds promotion criteria by promo_id
func (repository *RepositoryPromotionImpl) FindPromotionCriteriaByPromoID(promoID string) ([]model.PromotionCriteriaDetailRead, error) {
	var criteria []model.PromotionCriteriaDetailRead

	err := repository.Select(`
		ppc.promo_id,
		ppc.pro_id,
		mp.pro_code,
		mp.pro_name,
		COUNT(ppc.pro_id) AS count_promo,
		MIN(CASE 
			WHEN ps.range_from IS NOT NULL THEN ps.range_from::numeric
			WHEN pst.range_from IS NOT NULL THEN pst.range_from::numeric
			ELSE 0
		END) AS min_purchase,
		MAX(CASE 
			WHEN ps.range_to IS NOT NULL THEN ps.range_to::numeric
			WHEN pst.range_to IS NOT NULL THEN pst.range_to::numeric
			ELSE 0
		END) AS max_purchase,
		COALESCE(ps.rule_uom::text, pst.rule_uom::text, NULL) AS uom
	`).
		Table("promo.promotion_product_criteria AS ppc").
		Joins("INNER JOIN promo.promotions p ON p.promo_id = ppc.promo_id").
		Joins("LEFT JOIN mst.m_product mp ON mp.pro_id = ppc.pro_id").
		Joins("LEFT JOIN promo.promotion_slabs ps ON ps.promo_id = p.promo_id AND p.promo_type = 'slab'").
		Joins("LEFT JOIN promo.promotion_strata pst ON pst.promo_id = p.promo_id AND p.promo_type = 'strata'").
		Where("ppc.promo_id = ?", promoID).
		Group("ppc.promo_id, ppc.pro_id, mp.pro_code, mp.pro_name, ps.rule_uom::text, pst.rule_uom::text").
		Find(&criteria).Error

	return criteria, err
}

// FindPromotionRewardByPromoID finds promotion reward products by promo_id
func (repository *RepositoryPromotionImpl) FindPromotionRewardByPromoID(promoID, parentCustID string) ([]model.PromotionRewardDetailRead, error) {
	var rewards []model.PromotionRewardDetailRead

	err := repository.Select(`
		prp.promo_id,
		prp.id,
		prp.pro_id,
		mp.pro_code,
		mp.pro_name
	`).
		Table("promo.promotion_reward_products AS prp").
		Joins("INNER JOIN mst.m_product mp ON mp.pro_id = prp.pro_id AND mp.cust_id = ?", parentCustID).
		Where("prp.promo_id = ?", promoID).
		Find(&rewards).Error

	return rewards, err
}

// FindPromotionSlabByPromoID finds promotion slab by promo_id
func (repository *RepositoryPromotionImpl) FindPromotionSlabByPromoID(promoID string) ([]model.PromotionSlabRead, error) {
	var slabs []model.PromotionSlabRead

	err := repository.Select(`
		promo_id,
		id::text,
		description,
		ordinal,
		rule_type,
		range_from,
		range_to,
		rule_uom,
		reward_type,
		reward_uom,
		reward_value
	`).
		Table("promo.promotion_slabs").
		Where("promo_id = ?", promoID).
		Order("ordinal ASC").
		Find(&slabs).Error

	return slabs, err
}

// FindPromotionStrataByPromoID finds promotion strata by promo_id
func (repository *RepositoryPromotionImpl) FindPromotionStrataByPromoID(promoID string) ([]model.PromotionStrataRead, error) {
	var stratas []model.PromotionStrataRead

	err := repository.Select(`
		promo_id,
		id::text,
		description,
		ordinal,
		rule_type,
		range_from,
		range_to,
		rule_uom,
		reward_type,
		reward_uom,
		reward_value
	`).
		Table("promo.promotion_strata").
		Where("promo_id = ?", promoID).
		Order("ordinal ASC").
		Find(&stratas).Error

	return stratas, err
}

// FindOutletTypesByPromoID finds outlet types by promo_id
func (repository *RepositoryPromotionImpl) FindOutletTypesByPromoID(promoID, parentCustID string) ([]model.OutletTypeDetailRead, error) {
	var outletTypes []model.OutletTypeDetailRead

	err := repository.Select(`
		p.promo_id,
		mot.ot_type_id,
		mot.ot_type_code,
		mot.ot_type_name
	`).
		Table("promo.promotions AS p").
		Joins("LEFT JOIN promo.promotion_outlet_criteria poc ON p.promo_id = poc.promo_id").
		Joins("LEFT JOIN promo.promotion_outlet_attribute_type poat ON poat.criteria_id = poc.id").
		Joins("LEFT JOIN mst.m_outlet_type mot ON mot.ot_type_id = poat.outlet_type_id AND mot.cust_id = ?", parentCustID).
		Where("p.promo_id = ? AND mot.ot_type_id IS NOT NULL", promoID).
		Find(&outletTypes).Error

	return outletTypes, err
}

// FindOutletGroupsByPromoID finds outlet groups by promo_id
func (repository *RepositoryPromotionImpl) FindOutletGroupsByPromoID(promoID, parentCustID string) ([]model.OutletGroupDetailRead, error) {
	var outletGroups []model.OutletGroupDetailRead

	err := repository.Select(`
		p.promo_id,
		mog.ot_grp_id,
		mog.ot_grp_code,
		mog.ot_grp_name
	`).
		Table("promo.promotions AS p").
		Joins("LEFT JOIN promo.promotion_outlet_criteria poc ON p.promo_id = poc.promo_id").
		Joins("LEFT JOIN promo.promotion_outlet_attribute_group poag ON poag.criteria_id = poc.id").
		Joins("LEFT JOIN mst.m_outlet_group mog ON mog.ot_grp_id = poag.outlet_group_id AND mog.cust_id = ?", parentCustID).
		Where("p.promo_id = ? AND mog.ot_grp_id IS NOT NULL", promoID).
		Find(&outletGroups).Error

	return outletGroups, err
}

// FindOutletClassesByPromoID finds outlet classes by promo_id
func (repository *RepositoryPromotionImpl) FindOutletClassesByPromoID(promoID, parentCustID string) ([]model.OutletClassDetailRead, error) {
	var outletClasses []model.OutletClassDetailRead

	err := repository.Select(`
		p.promo_id,
		moc.ot_class_id,
		moc.ot_class_code,
		moc.ot_class_name
	`).
		Table("promo.promotions AS p").
		Joins("LEFT JOIN promo.promotion_outlet_criteria poc ON p.promo_id = poc.promo_id").
		Joins("LEFT JOIN promo.promotion_outlet_attribute_class poac ON poac.criteria_id = poc.id").
		Joins("LEFT JOIN mst.m_outlet_class moc ON moc.ot_class_id = poac.outlet_class_id AND moc.cust_id = ?", parentCustID).
		Where("p.promo_id = ? AND moc.ot_class_id IS NOT NULL", promoID).
		Find(&outletClasses).Error

	return outletClasses, err
}

func (repository *RepositoryPromotionImpl) FindOutletsByTypeGroupClass(dataFilter entity.PromotionOutletListQueryFilter, otTypeID int64, custId, parentCustId string) (outlets []model.PromotionOutletList, total int64, lastPage int, err error) {
	page := dataFilter.Page
	if page == 0 {
		page = 1
	}
	limit := dataFilter.Limit
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)
	todayStartUnix := todayStart.Unix()
	todayEndUnix := todayEnd.Unix()

	baseQuery := repository.
		Table("mst.m_outlet AS o").
		Joins("LEFT JOIN pjp.outlet_visit_list AS ovl ON ovl.outlet_id = o.outlet_id AND ovl.date = ?", todayStart).
		Where("o.cust_id = ? AND o.is_del = false AND o.ot_type_id = ?", custId, otTypeID)

	if dataFilter.OtGrpID != nil {
		baseQuery = baseQuery.Where("o.ot_grp_id = ?", *dataFilter.OtGrpID)
	}
	if dataFilter.OtClassID != nil {
		baseQuery = baseQuery.Where("o.ot_class_id = ?", *dataFilter.OtClassID)
	}

	// Count total records
	var count int64
	if err := baseQuery.Count(&count).Error; err != nil {
		return outlets, total, lastPage, err
	}
	total = count
	lastPage = int(math.Ceil(float64(total) / float64(limit)))

	query := baseQuery.
		Select(`
			o.ot_type_id,
			o.outlet_code,
			o.outlet_name,
			COALESCE(o.address1, '') AS address1,
			CASE 
				WHEN ovl.arrive_at IS NOT NULL 
					AND ovl.arrive_at >= ? 
					AND ovl.arrive_at <= ? 
				THEN true 
				ELSE false 
			END AS today_visit
		`, todayStartUnix, todayEndUnix).
		Order("o.outlet_name ASC").
		Offset(offset).
		Limit(limit)

	err = query.Scan(&outlets).Error

	if err != nil {
		return outlets, total, lastPage, err
	}

	return outlets, total, lastPage, nil
}
