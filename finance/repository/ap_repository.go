package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryApImpl struct {
		*gorm.DB
	}
)
type ApRepository interface {
	Store(c context.Context, data *model.Ap) error
	StoreDetail(c context.Context, data *model.ApDet) error
	StoreQtyPromoDetail(c context.Context, data *model.ApQtyPromo) error
	StoreMoneyPromoDetail(c context.Context, data *model.ApMoneyPromo) error
	FindByNo(apNo string, custId string) (ap model.ApList, err error)
	FindDetail(apNo string, custId string) (Details []model.ApDetRead, err error)
	FindQtyPromoDetail(apNo string, custId string) (Details []model.ApQtyPromoRead, err error)
	FindMoneyPromoDetail(apNo string, custId string) (Details []model.ApMoneyPromoRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApList, int64, int, error)
	Delete(c context.Context, custId string, apNo string, deletedBy int64) error
	Update(c context.Context, apNo string, data model.Ap) error
	DeleteDetailNotInIDs(c context.Context, apNo string, IDs []int64) error
	DeleteQtyPromoDetailNotInIDs(c context.Context, apNo string, IDs []int64) error
	DeleteMoneyPromoDetailNotInIDs(c context.Context, apNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.ApDet) error
	UpdateQtyPromoDetail(c context.Context, Details *model.ApQtyPromo) error
	UpdateMoneyPromoDetail(c context.Context, Details *model.ApMoneyPromo) error
}

func NewApRepo(db *gorm.DB) *RepositoryApImpl {
	return &RepositoryApImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryApImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryApImpl) Store(c context.Context, data *model.Ap) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApImpl) StoreDetail(c context.Context, data *model.ApDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApImpl) StoreQtyPromoDetail(c context.Context, data *model.ApQtyPromo) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApImpl) StoreMoneyPromoDetail(c context.Context, data *model.ApMoneyPromo) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApImpl) FindByNo(apNo string, custId string) (ap model.ApList, err error) {
	err = repository.Select("acf.ap.*, us.user_fullname AS updated_by_name,s.sup_code,s.sup_name").
		Joins("left join sys.m_user us on us.user_id = acf.ap.updated_by").
		Joins("left join mst.m_supplier s on s.sup_id = acf.ap.sup_id AND s.cust_id = ?", custId).
		Where("ap_no = ? AND acf.ap.cust_id=?", apNo, custId).
		Take(&ap).Error
	return ap, err
}

func (repository *RepositoryApImpl) FindDetail(apNo string, custId string) (Details []model.ApDetRead, err error) {
	err = repository.Select(`acf.ap_det.*, 
			p.pro_code, p.pro_name`).
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = acf.ap_det.pro_id").
		Joins("LEFT JOIN smc.m_customer smc ON smc.cust_id = ?", custId).
		// Joins("LEFT JOIN mst.m_dist_price dp ON dp.pro_id = ap_det.pro_id AND dp.dist_price_group_id = smc.dist_price_grp_id AND dp.cust_id = smc.parent_cust_id").
		Where("ap_no = ? AND acf.ap_det.cust_id=?", apNo, custId).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApImpl) FindQtyPromoDetail(apNo string, custId string) (Details []model.ApQtyPromoRead, err error) {
	err = repository.Select(`acf.ap_qty_promo.*, 
			p.pro_code, p.pro_name`).
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = acf.ap_qty_promo.pro_id").
		Joins("LEFT JOIN smc.m_customer smc ON smc.cust_id = ?", custId).
		// Joins("LEFT JOIN mst.m_dist_price dp ON dp.pro_id = ap_qty_promo.pro_id AND dp.dist_price_group_id = smc.dist_price_grp_id AND dp.cust_id = smc.parent_cust_id").
		Where("ap_no = ? AND acf.ap_qty_promo.cust_id=?", apNo, custId).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApImpl) FindMoneyPromoDetail(apNo string, custId string) (Details []model.ApMoneyPromoRead, err error) {
	err = repository.Select("acf.ap_money_promo.*, p.pro_code, p.pro_name").
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = acf.ap_money_promo.pro_id").
		Joins("LEFT JOIN smc.m_customer smc ON smc.cust_id = ?", custId).
		// Joins("LEFT JOIN mst.m_dist_price dp ON dp.pro_id = ap_money_promo.pro_id AND dp.dist_price_group_id = smc.dist_price_grp_id AND dp.cust_id = smc.parent_cust_id").
		Where("ap_no = ? AND acf.ap_money_promo.cust_id=?", apNo, custId).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApList, int64, int, error) {
	var ap []model.ApList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ap_no")
	query := repository.Select("acf.ap.*, us.user_fullname AS updated_by_name,s.sup_code,s.sup_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.ap.updated_by").
		Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = acf.ap.sup_id AND s.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.ap.cust_id=?", dataFilter.CustId)
	query.Where("acf.ap.cust_id=?", dataFilter.CustId)
	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.ap.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.ap.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}
	if dataFilter.Query != "" {
		queryCount.Where("acf.ap.ap_no=?", dataFilter.Query)
		query.Where("acf.ap.ap_no=?", dataFilter.Query)
	}

	if dataFilter.TrCode != "" {
		queryCount.Where("acf.ap.tr_code=?", dataFilter.TrCode)
		query.Where("acf.ap.tr_code=?", dataFilter.TrCode)
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
		query.Order("ap_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ap).Error
	if err != nil {
		return ap, total, 0, err
	}
	err = queryCount.Model(&ap).Count(&total).Error
	if err != nil {
		return ap, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ap, total, lastPage, nil
}

func (repository *RepositoryApImpl) Delete(c context.Context, custId string, apNo string, deletedBy int64) error {
	var data model.Ap
	result := repository.model(c).Model(&data).Where("ap_no=? AND cust_id = ? AND is_del= ? ", apNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryApImpl) Update(c context.Context, apNo string, data model.Ap) error {

	result := repository.model(c).Model(&data).Where("ap_no=?", apNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryApImpl) DeleteDetailNotInIDs(c context.Context, apNo string, IDs []int64) error {
	var Details model.ApDet
	err := repository.model(c).Where("ap_no=? AND ap_det_id not in (?) ", apNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryApImpl) DeleteQtyPromoDetailNotInIDs(c context.Context, apNo string, IDs []int64) error {
	var Details model.ApQtyPromo
	err := repository.model(c).Where("ap_no=? AND ap_qty_promo_id not in (?) ", apNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryApImpl) DeleteMoneyPromoDetailNotInIDs(c context.Context, apNo string, IDs []int64) error {
	var Details model.ApMoneyPromo
	err := repository.model(c).Where("ap_no=? AND ap_money_promo_id not in (?) ", apNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryApImpl) UpdateDetail(c context.Context, Details *model.ApDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryApImpl) UpdateQtyPromoDetail(c context.Context, Details *model.ApQtyPromo) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryApImpl) UpdateMoneyPromoDetail(c context.Context, Details *model.ApMoneyPromo) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
