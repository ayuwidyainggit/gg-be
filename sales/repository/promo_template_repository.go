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
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type (
	RepositoryPromoTemplateImpl struct {
		*gorm.DB
	}
)
type PromoTemplateRepository interface {
	Store(c context.Context, data *model.PromoTemplate) (promoTemplateID string, err error)
	FindByPromoTemplateID(params entity.DetailPromoTemplateParams) (promotion model.PromoTemplate, err error)
	FindAllPromoCriteriasByPromoTemplateID(params entity.DetailPromoTemplateParams) (promoCriterias []model.PromoTemplateCriteria, err error)
	FindAllPromoAdditionalCriteriasByPromoTemplateID(params entity.DetailPromoTemplateParams) (promoAddCriterias []model.PromoTemplateAdditionalCriteria, err error)
	FindAllRewardProductsByPromoTemplateID(params entity.DetailPromoTemplateParams) (promoRewardProducts []model.PromoTemplateRewardProduct, err error)
	FindProductByProID(custID string, proID int64) (rewardProductDetail model.PromoTemplateRewardProductDetail, err error)
	FindAllByCustId(dataFilter entity.PromoTemplateQueryFilter) ([]model.PromoTemplate, int64, int, error)
	Update(c context.Context, promoTemplateID string, data model.PromoTemplate) error
	Delete(c context.Context, custID, promoTemplateID string) error
	DeletePromoCriterias(c context.Context, custID, promoTemplateID string) error
	DeletePromoAdditionalCriterias(c context.Context, custId, promoTemplateID string) error
	DeletePromoTemplateRewardProducts(c context.Context, custId, promoTemplateID string) error
	StorePromoCriteria(c context.Context, data *model.PromoTemplateCriteria) error
	StorePromoRewardProduct(c context.Context, data *model.PromoTemplateRewardProduct) error
	StorePromoAdditionalCriteria(c context.Context, data *model.PromoTemplateAdditionalCriteria) error
	DeletePromoCriteriasNotInIDs(c context.Context, custID, promoTemplateID string, IDs []int64) error
	DeletePromoAdditionalCriteriasNotInIDs(c context.Context, custID, promoTemplateID string, IDs []int64) error
	UpdatePromoCriteria(c context.Context, promoTemplateCriteria *model.PromoTemplateCriteria) error
	UpdatePromoAdditionalCriteria(c context.Context, promoAddCriteria *model.PromoTemplateAdditionalCriteria) error
	FindOneProductByProID(custID, parentCustID string, proID int64) (productAdditionalCriteria model.ProductAdditionalCriteria, err error)
	FindOneOutletClassByID(custID, parentCustID string, id int64) (outletClassAdditionalCriteria model.OutletClassAdditionalCriteria, err error)
	FindOneOutletTypeByID(custID, parentCustID string, id int64) (outletTypeAdditionalCriteria model.OutletTypeAdditionalCriteria, err error)
	FindOneOutletGroupByID(custID, parentCustID string, id int64) (outletGroupAdditionalCriteria model.OutletGroupAdditionalCriteria, err error)
	FindOneSalesTypeByID(custID, parentCustID string, id int64) (salesTypeAdditionalCriteria model.SalesTypeAdditionalCriteria, err error)
	FindOneSalesTeamByID(custID, parentCustID string, id int64) (salesTeamAdditionalCriteria model.SalesTeamAdditionalCriteria, err error)
}

func NewPromoTemplateRepo(db *gorm.DB) *RepositoryPromoTemplateImpl {
	return &RepositoryPromoTemplateImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryPromoTemplateImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryPromoTemplateImpl) Store(c context.Context, data *model.PromoTemplate) (promoTemplateID string, err error) {
	err = repository.model(c).Create(data).Error
	if err != nil {
		return data.PromoTemplateID, err
	}
	return data.PromoTemplateID, err
}

func (repository *RepositoryPromoTemplateImpl) FindByPromoTemplateID(params entity.DetailPromoTemplateParams) (promotion model.PromoTemplate, err error) {
	err = repository.
		Select(`sls.promo_templates.*`).
		Where("sls.promo_templates.promo_template_id = ? AND sls.promo_templates.cust_id=?", params.PromoTemplateID, params.CustID).
		Take(&promotion).Error
	return promotion, err
}

func (repository *RepositoryPromoTemplateImpl) FindAllByCustId(dataFilter entity.PromoTemplateQueryFilter) ([]model.PromoTemplate, int64, int, error) {
	var (
		promo []model.PromoTemplate
		total int64
	)
	limit := 10
	if dataFilter.Limit != 0 {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("promo_template_id")
	query := repository.Select(
		`sls.promo_templates.*`) // .
	// Joins("left join mst.m_outlet ot on ot.outlet_id = sls.promo_templates.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.promo_templates.cust_id=?", dataFilter.CustId)
	query.Where("sls.promo_templates.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.promo_templates.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.promo_templates.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.promo_templates.promo_template_id ILIKE ? OR sls.promo_templates.promo_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query.Where("sls.promo_templates.promo_template_id ILIKE ? OR sls.promo_templates.promo_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if dataFilter.PromoTemplateID != "" {
		queryCount.Where("sls.promo_templates.promo_template_id ILIKE ?", "%"+dataFilter.PromoTemplateID+"%")
		query.Where("sls.promo_templates.promo_template_id ILIKE ?", "%"+dataFilter.PromoTemplateID+"%")
	}

	if dataFilter.PromoDesc != "" {
		queryCount.Where("sls.promo_templates.promo_desc ILIKE ?", "%"+dataFilter.PromoDesc+"%")
		query.Where("sls.promo_templates.promo_desc ILIKE ?", "%"+dataFilter.PromoDesc+"%")
	}

	if len(dataFilter.PromoTemplateStatusID) > 0 {
		queryCount.Where("sls.promo_templates.promo_template_status_id IN ?", dataFilter.PromoTemplateStatusID)
		query.Where("sls.promo_templates.promo_template_status_id IN ?", dataFilter.PromoTemplateStatusID)
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

func (repository *RepositoryPromoTemplateImpl) Update(c context.Context, RoNo string, data model.PromoTemplate) error {
	result := repository.model(c).Model(&data).Where("promo_template_id=?", RoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.PromoTemplatewsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryPromoTemplateImpl) Delete(c context.Context, custId, promoID string) error {
	var data model.PromoTemplate
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "promo_template_id": promoID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryPromoTemplateImpl) DeletePromoCriterias(c context.Context, custId, promoTemplateID string) error {
	var data model.PromoTemplateCriteria
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "promo_template_id": promoTemplateID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryPromoTemplateImpl) DeletePromoAdditionalCriterias(c context.Context, custId, promoTemplateID string) error {
	var data model.PromoTemplateAdditionalCriteria
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "promo_template_id": promoTemplateID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryPromoTemplateImpl) DeletePromoTemplateRewardProducts(c context.Context, custId, promoTemplateID string) error {
	var data model.PromoTemplateRewardProduct
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "promo_template_id": promoTemplateID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryPromoTemplateImpl) StorePromoCriteria(c context.Context, data *model.PromoTemplateCriteria) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromoTemplateImpl) StorePromoAdditionalCriteria(c context.Context, data *model.PromoTemplateAdditionalCriteria) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromoTemplateImpl) FindAllPromoCriteriasByPromoTemplateID(params entity.DetailPromoTemplateParams) (promoCriterias []model.PromoTemplateCriteria, err error) {
	err = repository.
		Select(`promo_template_criterias.*`).
		// Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Where("promo_template_criterias.promo_template_id = ? AND promo_template_criterias.cust_id = ?", params.PromoTemplateID, params.CustID).
		Order("promo_template_criterias.promo_template_slab_id ASC").
		Find(&promoCriterias).Error
	return promoCriterias, err
}

func (repository *RepositoryPromoTemplateImpl) FindAllPromoAdditionalCriteriasByPromoTemplateID(params entity.DetailPromoTemplateParams) (promoAddCriterias []model.PromoTemplateAdditionalCriteria, err error) {
	err = repository.
		Select(`promo_template_additional_criterias.*`).
		Where("promo_template_additional_criterias.promo_template_id = ? AND promo_template_additional_criterias.cust_id = ?", params.PromoTemplateID, params.CustID).
		Order("promo_template_additional_criterias.promo_template_add_criteria_id ASC").
		Find(&promoAddCriterias).Error
	return promoAddCriterias, err
}

func (repository *RepositoryPromoTemplateImpl) FindAllRewardProductsByPromoTemplateID(params entity.DetailPromoTemplateParams) (promoRewardProducts []model.PromoTemplateRewardProduct, err error) {
	err = repository.
		Select(`promo_template_reward_products.*`).
		Where("promo_template_reward_products.promo_template_id = ? AND promo_template_reward_products.cust_id = ?", params.PromoTemplateID, params.CustID).
		Order("promo_template_reward_products.promo_template_reward_id ASC").
		Find(&promoRewardProducts).Error
	return promoRewardProducts, err
}

func (repository *RepositoryPromoTemplateImpl) FindProductByProID(custID string, proID int64) (rewardProductDetail model.PromoTemplateRewardProductDetail, err error) {
	err = repository.
		Select(`m_product.*`).
		Where("m_product.cust_id = ? AND m_product.pro_id = ? AND m_product.is_active = ?", custID, proID, true).
		Take(&rewardProductDetail).Error
	return rewardProductDetail, err
}

func (repository *RepositoryPromoTemplateImpl) DeletePromoCriteriasNotInIDs(c context.Context, custID, promoTemplateID string, IDs []int64) error {
	var promoCriterias model.PromoTemplateCriteria
	err := repository.model(c).Where("cust_id = ? AND promo_template_id = ? AND promo_template_slab_id NOT IN (?) ", custID, promoTemplateID, IDs).Delete(&promoCriterias).Error
	return err
}

func (repository *RepositoryPromoTemplateImpl) DeletePromoAdditionalCriteriasNotInIDs(c context.Context, custID, promoID string, IDs []int64) error {
	var promoAddCriterias model.PromoTemplateAdditionalCriteria
	err := repository.model(c).Where("cust_id = ? AND promo_template_id = ? AND promo_template_add_criteria_id NOT IN (?) ", custID, promoID, IDs).Delete(&promoAddCriterias).Error
	return err
}

func (repository *RepositoryPromoTemplateImpl) UpdatePromoCriteria(c context.Context, promoCriteria *model.PromoTemplateCriteria) error {
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

func (repository *RepositoryPromoTemplateImpl) UpdatePromoAdditionalCriteria(c context.Context, promoAddCriteria *model.PromoTemplateAdditionalCriteria) error {
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

func (repository *RepositoryPromoTemplateImpl) StorePromoRewardProduct(c context.Context, data *model.PromoTemplateRewardProduct) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryPromoTemplateImpl) FindOneProductByProID(custID, parentCustID string, proID int64) (productAddCriteria model.ProductAdditionalCriteria, err error) {
	err = repository.
		Select(`m_product.pro_id AS reference_id, m_product.pro_code AS reference_code, m_product.pro_name AS reference_name`).
		Where("m_product.cust_id = ? AND m_product.pro_id = ?", parentCustID, proID).
		Take(&productAddCriteria).Error
	return productAddCriteria, err
}

func (repository *RepositoryPromoTemplateImpl) FindOneOutletClassByID(custID, parentCustID string, outletClassID int64) (outletClassAddCriteria model.OutletClassAdditionalCriteria, err error) {
	err = repository.
		Select(`m_outlet_class.ot_class_id AS reference_id, m_outlet_class.ot_class_code AS reference_code, m_outlet_class.ot_class_name AS reference_name`).
		Where("m_outlet_class.cust_id = ? AND m_outlet_class.ot_class_id = ?", parentCustID, outletClassID).
		Take(&outletClassAddCriteria).Error
	return outletClassAddCriteria, err
}

func (repository *RepositoryPromoTemplateImpl) FindOneOutletTypeByID(custID, parentCustID string, outletTypeID int64) (outletTypeAddCriteria model.OutletTypeAdditionalCriteria, err error) {
	err = repository.
		Select(`m_outlet_type.ot_type_id AS reference_id, m_outlet_type.ot_type_code AS reference_code, m_outlet_type.ot_type_name AS reference_name`).
		Where("m_outlet_type.cust_id = ? AND m_outlet_type.ot_type_id = ?", parentCustID, outletTypeID).
		Take(&outletTypeAddCriteria).Error
	return outletTypeAddCriteria, err
}

func (repository *RepositoryPromoTemplateImpl) FindOneOutletGroupByID(custID, parentCustID string, outletGroupID int64) (outletGroupAddCriteria model.OutletGroupAdditionalCriteria, err error) {
	err = repository.
		Select(`m_outlet_group.ot_grp_id AS reference_id, m_outlet_grp.ot_group_code AS reference_code, m_outlet_group.ot_grp_name AS reference_name`).
		Where("m_outlet_group.cust_id = ? AND m_outlet_group.ot_grp_id = ?", parentCustID, outletGroupID).
		Take(&outletGroupAddCriteria).Error
	return outletGroupAddCriteria, err
}

func (repository *RepositoryPromoTemplateImpl) FindOneSalesTypeByID(custID, parentCustID string, salesTypeID int64) (salesTypeAddCriteria model.SalesTypeAdditionalCriteria, err error) {
	err = repository.
		Select(`m_sales_type.sales_type_id AS reference_id, m_sales_type.sales_type_code AS reference_code, m_sales_type.sales_type_name AS reference_name`).
		Where("m_sales_type.cust_id = ? AND m_sales_type.sales_type_id = ?", custID, salesTypeID).
		Take(&salesTypeAddCriteria).Error
	return salesTypeAddCriteria, err
}

func (repository *RepositoryPromoTemplateImpl) FindOneSalesTeamByID(custID, parentCustID string, salesTeamID int64) (salesTeamAddCriteria model.SalesTeamAdditionalCriteria, err error) {
	err = repository.
		Select(`m_sales_team.sales_team_id AS reference_id, m_sales_team.sales_team_code AS reference_code, m_sales_team.sales_team_name AS reference_name`).
		Where("m_sales_team.cust_id = ? AND m_sales_team.sales_team_id = ?", custID, salesTeamID).
		Take(&salesTeamAddCriteria).Error
	return salesTeamAddCriteria, err
}
