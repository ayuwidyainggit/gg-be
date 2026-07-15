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
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type (
	RepositoryDiscountImpl struct {
		*gorm.DB
	}
)
type DiscountRepository interface {
	Store(c context.Context, data *model.Discount) error
	// StoreStatusLog(c context.Context, data *model.DiscountStatusLog) error
	FindByDiscountID(params entity.DetailDiscountParams) (discounttion model.Discount, err error)
	FindAllDiscountPrincipalsByDiscountID(params entity.DetailDiscountParams) (discountCriterias []model.DiscountPrincipalDetail, err error)
	FindAllDiscountGroupsByDiscountID(params entity.DetailDiscountParams) (discountCriterias []model.DiscountGroupDetail, err error)
	FindAllDiscountCriteriasByDiscountID(params entity.DetailDiscountParams) (discountCriterias []model.DiscountCriteria, err error)

	FindAllByCustId(dataFilter entity.DiscountQueryFilter) ([]model.Discount, int64, int, error)
	FindDiscGrpId(DiscGrpId string) ([]model.OutletRead, error)

	Update(c context.Context, discountID string, data model.Discount) error
	Delete(c context.Context, custID, discountID string) error
	DeleteDiscountPrincipals(c context.Context, custID, discountID string) error
	DeleteDiscountGroups(c context.Context, custID, discountID string) error
	DeleteDiscountCriterias(c context.Context, custID, discountID string) error

	StoreDiscountPrincipal(c context.Context, data *model.DiscountPrincipal) error
	StoreDiscountGroup(c context.Context, data *model.DiscountGroup) error
	StoreDiscountCriteria(c context.Context, data *model.DiscountCriteria) error

	DeleteDiscountCriteriasNotInIDs(c context.Context, custID, discountID string, IDs []int64) error

	UpdateDiscountCriteria(c context.Context, discountCriteria *model.DiscountCriteria) error

	FindAllByCustIdAndDiscountID(request entity.PublishDiscountBody) ([]model.Discount, error)
	PublishDiscount(c context.Context, req entity.PublishDiscountBody) error

	FindOutletByID(outletID int, custId string, parentCustId string) (outlet model.OutletRead, err error)
	FindProductByID(productID int) (product model.ProductRead, err error)
	FindDiscountByProductAndOutlet(product model.ProductRead, outlet model.OutletRead, request entity.ConsultDiscountBody) (discount model.DiscountRead, err error)
	FindDiscountCriteriaBySubTotal(discountID string, subTotal int) (discountCriteria model.DiscountCriteria, err error)
}

func NewDiscountRepo(db *gorm.DB) *RepositoryDiscountImpl {
	return &RepositoryDiscountImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryDiscountImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryDiscountImpl) Store(c context.Context, data *model.Discount) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

// func (repository *RepositoryDiscountImpl) StoreStatusLog(c context.Context, data *model.DiscountStatusLog) error {
// 	err := repository.model(c).Create(data).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (repository *RepositoryDiscountImpl) FindByDiscountID(params entity.DetailDiscountParams) (discounttion model.Discount, err error) {
	err = repository.
		Select(`sls.discounts.*`).
		Where("sls.discounts.discount_id = ? AND sls.discounts.cust_id=?", params.DiscountID, params.CustID).
		Take(&discounttion).Error
	return discounttion, err
}

func (repository *RepositoryDiscountImpl) FindAllByCustId(dataFilter entity.DiscountQueryFilter) ([]model.Discount, int64, int, error) {
	var (
		discount []model.Discount
		total    int64
	)
	limit := 10
	if dataFilter.Limit != 0 {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("discount_id")
	query := repository.Select(
		`sls.discounts.*`)
	queryCount.Where("sls.discounts.cust_id=?", dataFilter.CustId)
	query.Where("sls.discounts.cust_id=?", dataFilter.CustId)

	if dataFilter.EffectiveFrom != nil && dataFilter.EffectiveTo != nil {
		query.Where(`sls.discounts.effective_from BETWEEN ? AND ? 
					OR sls.discounts.effective_to BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
		)
		queryCount.Where(`sls.discounts.effective_from BETWEEN ? AND ?
					OR sls.discounts.effective_to BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
			str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
		)
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.discounts.discount_id ILIKE ? OR sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query.Where("sls.discounts.discount_id ILIKE ? OR sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if dataFilter.DiscountID != "" {
		queryCount.Where("sls.discounts.discount_id ILIKE ?", "%"+dataFilter.DiscountID+"%")
		query.Where("sls.discounts.discount_id ILIKE ?", "%"+dataFilter.DiscountID+"%")
	}

	if dataFilter.DiscountDesc != "" {
		queryCount.Where("sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.DiscountDesc+"%")
		query.Where("sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.DiscountDesc+"%")
	}

	if len(dataFilter.DiscountStatusID) > 0 {
		queryCount.Where("sls.discounts.discount_status_id IN ?", dataFilter.DiscountStatusID)
		query.Where("sls.discounts.discount_status_id IN ?", dataFilter.DiscountStatusID)
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
	err := query.Limit(limit).Offset(offset).Find(&discount).Error
	if err != nil {
		return discount, total, 0, err
	}
	err = queryCount.Model(&discount).Count(&total).Error
	if err != nil {
		return discount, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return discount, total, lastPage, nil
}

func (repository *RepositoryDiscountImpl) Update(c context.Context, discountID string, data model.Discount) error {
	result := repository.model(c).Model(&data).Where("discount_id=?", discountID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.DiscountwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryDiscountImpl) Delete(c context.Context, custId, discountID string) error {
	var data model.Discount
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "discount_id": discountID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryDiscountImpl) DeleteDiscountPrincipals(c context.Context, custId, discountID string) error {
	var data model.DiscountPrincipal
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "discount_id": discountID})
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryDiscountImpl) DeleteDiscountGroups(c context.Context, custId, discountID string) error {
	var data model.DiscountGroup
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "discount_id": discountID})
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryDiscountImpl) DeleteDiscountCriterias(c context.Context, custId, discountID string) error {
	var data model.DiscountCriteria
	result := repository.model(c).
		Delete(&data, map[string]interface{}{"cust_id": custId, "discount_id": discountID})
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryDiscountImpl) StoreDiscountPrincipal(c context.Context, data *model.DiscountPrincipal) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryDiscountImpl) StoreDiscountGroup(c context.Context, data *model.DiscountGroup) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryDiscountImpl) StoreDiscountCriteria(c context.Context, data *model.DiscountCriteria) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryDiscountImpl) FindAllDiscountPrincipalsByDiscountID(params entity.DetailDiscountParams) (discountCriterias []model.DiscountPrincipalDetail, err error) {
	err = repository.
		Select(`discount_principals.*, princ.principal_code, princ.principal_name`).
		Joins("LEFT JOIN mst.m_principal princ ON princ.principal_id = discount_principals.principal_id AND princ.cust_id = ?", params.ParentCustId).
		Where("discount_principals.discount_id = ? AND discount_principals.cust_id = ?", params.DiscountID, params.CustID).
		Order("discount_principals.principal_id ASC").
		Find(&discountCriterias).Error
	return discountCriterias, err
}

func (repository *RepositoryDiscountImpl) FindAllDiscountGroupsByDiscountID(params entity.DetailDiscountParams) (discountGroups []model.DiscountGroupDetail, err error) {
	err = repository.
		Select(`discount_groups.*, disc_grp_code, disc_grp_name`).
		Joins("LEFT JOIN mst.m_disc_group mdisc ON mdisc.disc_grp_id = discount_groups.disc_grp_id AND mdisc.cust_id = ?", params.CustID).
		Where("discount_groups.discount_id = ? AND discount_groups.cust_id = ?", params.DiscountID, params.CustID).
		Order("discount_groups.disc_grp_id ASC").
		Find(&discountGroups).Error
	return discountGroups, err
}

func (repository *RepositoryDiscountImpl) FindAllDiscountCriteriasByDiscountID(params entity.DetailDiscountParams) (discountCriterias []model.DiscountCriteria, err error) {
	err = repository.
		Select(`discount_criterias.*`).
		Where("discount_criterias.discount_id = ? AND discount_criterias.cust_id = ?", params.DiscountID, params.CustID).
		Order("discount_criterias.slab_desc ASC").
		Find(&discountCriterias).Error
	return discountCriterias, err
}

func (repository *RepositoryDiscountImpl) FindProductByProID(custID string, proID int64) (rewardProductDetail model.RewardProductDetail, err error) {
	err = repository.
		Select(`m_product.*`).
		Where("m_product.cust_id = ? AND m_product.pro_id = ? AND m_product.is_active = ?", custID, proID, true).
		Take(&rewardProductDetail).Error
	return rewardProductDetail, err
}

func (repository *RepositoryDiscountImpl) DeleteDiscountCriteriasNotInIDs(c context.Context, custID, discountID string, IDs []int64) error {
	var discountCriterias model.DiscountCriteria
	err := repository.model(c).Where("cust_id = ? AND discount_id = ? AND slab_id NOT IN (?) ", custID, discountID, IDs).Delete(&discountCriterias).Error
	return err
}

func (repository *RepositoryDiscountImpl) UpdateDiscountCriteria(c context.Context, discountCriteria *model.DiscountCriteria) error {
	result := repository.model(c).Updates(&discountCriteria)
	if result.Error != nil {
		log.Error("UpdateDiscountCriteria, result:", structs.StructToJson(result.Error))
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryDiscountImpl) FindAllByCustIdAndDiscountID(request entity.PublishDiscountBody) ([]model.Discount, error) {
	var (
		discount []model.Discount
	)

	query := repository.Select(`sls.discounts.discount_id, sls.discounts.discount_status_id`)
	query.Where("discount_status_id = 1 AND publish_status_id = 1 AND discounts.cust_id = ? AND discount_id IN ?", request.CustID, request.DiscountID)

	err := query.Find(&discount).Error
	if err != nil {
		return discount, err
	}

	return discount, nil
}

func (repository *RepositoryDiscountImpl) PublishDiscount(c context.Context, request entity.PublishDiscountBody) error {
	discountStatusID := model.Discount{
		DiscountStatusID: 2,
		PublishStatusID:  2,
		UpdatedBy:        request.UpdatedBy,
	}
	discountModel := model.Discount{}
	result := repository.model(c).Model(&discountModel).
		Where(`cust_id = ? AND discount_id IN ?`, request.CustID, request.DiscountID).
		Updates(discountStatusID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != int64(len(request.DiscountID)) {
		return errors.New("rows affected is different with the request")
	}
	return nil
}

func (repository *RepositoryDiscountImpl) FindOutletByID(outletID int, custId string, parentCustId string) (outlet model.OutletRead, err error) {
	err = repository.Select(`
			mst.m_outlet.outlet_id, 
			mst.m_outlet.outlet_code, 
			mst.m_outlet.outlet_name, 
			mst.m_outlet.disc_grp_id 
		`).
		Where("mst.m_outlet.outlet_id=?", outletID).
		Take(&outlet).Error

	return outlet, err
}

func (repository *RepositoryDiscountImpl) FindProductByID(productID int) (product model.ProductRead, err error) {
	err = repository.Select(`
			mst.m_product.pro_id, 
			mst.m_product.pro_code, 
			mst.m_product.pro_name, 
			mst.m_product.principal_id 
		`).
		Where("mst.m_product.pro_id=?", productID).
		Take(&product).Error

	return product, err
}

func (repository *RepositoryDiscountImpl) FindDiscountByProductAndOutlet(product model.ProductRead, outlet model.OutletRead, request entity.ConsultDiscountBody) (discount model.DiscountRead, err error) {
	err = repository.Select(`
			sls.discounts.discount_id, 
			sls.discounts.discount_desc
		`).
		// Joins("inner join").
		Where("sls.discounts.discount_status_id = 2").
		Where("sls.discounts.publish_status_id = 2").
		Where("sls.discounts.effective_from <= ?", request.OrderDate).
		Where("sls.discounts.effective_to >= ?", request.OrderDate).
		Where("sls.discounts.discount_id IN (SELECT discount_id from sls.discount_principals WHERE principal_id = ?)", product.PrincipalId).
		Where("sls.discounts.discount_id IN (SELECT discount_id from sls.discount_groups WHERE disc_grp_id = ?)", outlet.DiscGrpId).
		Take(&discount).Error

	return discount, err
}

func (repository *RepositoryDiscountImpl) FindDiscountCriteriaBySubTotal(discountID string, subTotal int) (discountCriteria model.DiscountCriteria, err error) {
	err = repository.Select(`
			sls.discount_criterias.*
		`).
		Where("sls.discount_criterias.discount_id = ?", discountID).
		Where("sls.discount_criterias.slab_rule_from <= ?", subTotal).
		Where("sls.discount_criterias.slab_rule_to >= ?", subTotal).
		Take(&discountCriteria).Error

	return discountCriteria, err
}

func (repository *RepositoryDiscountImpl) FindDiscGrpId(DiscGrpId string) ([]model.OutletRead, error) {
	var (
		discount []model.OutletRead
		// total    int64
	)
	// limit := 10
	// if dataFilter.Limit != 0 {
	// 	limit = dataFilter.Limit
	// }

	// queryCount := repository.Select("disc_grp_id")
	query := repository.Select(`mst.m_outlet.*`)
	query.Where("mst.m_outlet.disc_grp_id=?", DiscGrpId)
	// query.Where("sls.discounts.cust_id=?", dataFilter.CustId)

	// if dataFilter.EffectiveFrom != nil && dataFilter.EffectiveTo != nil {
	// 	query.Where(`sls.discounts.effective_from BETWEEN ? AND ?
	// 				OR sls.discounts.effective_to BETWEEN ? AND ?`,
	// 		str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
	// 		str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
	// 	)
	// 	queryCount.Where(`sls.discounts.effective_from BETWEEN ? AND ?
	// 				OR sls.discounts.effective_to BETWEEN ? AND ?`,
	// 		str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
	// 		str.UnixTimestampToUtcTime(*dataFilter.EffectiveFrom), str.UnixTimestampToUtcTime(*dataFilter.EffectiveTo),
	// 	)
	// }

	// if dataFilter.Query != "" {
	// 	queryCount.Where("sls.discounts.discount_id ILIKE ? OR sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	// 	query.Where("sls.discounts.discount_id ILIKE ? OR sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	// }

	// if dataFilter.DiscountID != "" {
	// 	queryCount.Where("sls.discounts.discount_id ILIKE ?", "%"+dataFilter.DiscountID+"%")
	// 	query.Where("sls.discounts.discount_id ILIKE ?", "%"+dataFilter.DiscountID+"%")
	// }

	// if dataFilter.DiscountDesc != "" {
	// 	queryCount.Where("sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.DiscountDesc+"%")
	// 	query.Where("sls.discounts.discount_desc ILIKE ?", "%"+dataFilter.DiscountDesc+"%")
	// }

	// if len(dataFilter.DiscountStatusID) > 0 {
	// 	queryCount.Where("sls.discounts.discount_status_id IN ?", dataFilter.DiscountStatusID)
	// 	query.Where("sls.discounts.discount_status_id IN ?", dataFilter.DiscountStatusID)
	// }

	// sortBy := ``
	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	query.Order(sortBy)
	// } else {
	// 	query.Order("created_at DESC")
	// }

	// page := dataFilter.Page
	// if page-1 < 1 {
	// 	page = 1
	// }
	// offset := (page - 1) * dataFilter.Limit
	err := query.Find(&discount).Error
	if err != nil {
		return discount, err
	}
	// err = queryCount.Model(&discount).Count(&total).Error
	// if err != nil {
	// 	return discount, total, 0, err
	// }
	// lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return discount, nil

}
