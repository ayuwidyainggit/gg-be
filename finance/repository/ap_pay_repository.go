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
	RepositoryApPayImpl struct {
		*gorm.DB
	}
)
type ApPayRepository interface {
	Store(c context.Context, data *model.ApPay) error
	StoreDetail(c context.Context, data *model.ApPayDet) error
	StoreApPaymethod(c context.Context, data *model.ApPayMethod) error

	FindByNo(apPayNo string, custId, parentCustId string) (ap model.ApPayList, err error)
	FindDetail(apPayNo string, custId string) (Details []model.ApPayDet, err error)
	FindApPayMethod(apPayNo string, custId string) (Details []model.ApPayMethod, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApPayList, int64, int, error)

	Delete(c context.Context, custId string, apPayNo string, deletedBy int64) error

	Update(c context.Context, apPayNo string, data model.ApPay) error
	DeleteDetailNotInIDs(c context.Context, apPayNo string, IDs []int64) error
	DeleteApPayMethodDetailNotInIDs(c context.Context, apPayNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.ApPayDet) error
	UpdateApPayMethodDetail(c context.Context, Details *model.ApPayMethod) error
}

func NewApPayRepo(db *gorm.DB) *RepositoryApPayImpl {
	return &RepositoryApPayImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryApPayImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryApPayImpl) Store(c context.Context, data *model.ApPay) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApPayImpl) StoreDetail(c context.Context, data *model.ApPayDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApPayImpl) StoreApPaymethod(c context.Context, data *model.ApPayMethod) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApPayImpl) FindByNo(apPayNo string, custId, parentCustId string) (ap model.ApPayList, err error) {
	err = repository.Select("acf.ap_pay.*,us.user_fullname AS updated_by_name,sup.sup_code,sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = acf.ap_pay.updated_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.ap_pay.sup_id AND sup.cust_id = ?", parentCustId).
		Where("acf.ap_pay.ap_pay_no = ? AND acf.ap_pay.cust_id=?", apPayNo, custId).
		Take(&ap).Error
	return ap, err
}

func (repository *RepositoryApPayImpl) FindDetail(apPayNo string, custId string) (Details []model.ApPayDet, err error) {
	err = repository.
		Where("ap_pay_no = ? AND cust_id=?", apPayNo, custId).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApPayImpl) FindApPayMethod(apPayNo string, custId string) (Details []model.ApPayMethod, err error) {
	err = repository.
		Where("ap_pay_no = ? AND cust_id=?", apPayNo, custId).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApPayImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApPayList, int64, int, error) {
	var ap []model.ApPayList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ap_pay_no")
	query := repository.Select("acf.ap_pay.*,us.user_fullname AS updated_by_name,sup.sup_code,sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = acf.ap_pay.updated_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.ap_pay.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.ap_pay.cust_id=?", dataFilter.CustId)
	query.Where("acf.ap_pay.cust_id=?", dataFilter.CustId)
	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.ap_pay.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.ap_pay.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}
	if dataFilter.Query != "" {
		queryCount.Where("acf.ap_pay.ap_pay_no=?", dataFilter.Query)
		query.Where("acf.ap_pay.ap_pay_no=?", dataFilter.Query)
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
		query.Order("ap_pay_no DESC")
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

func (repository *RepositoryApPayImpl) Delete(c context.Context, custId string, apPayNo string, deletedBy int64) error {
	var data model.ApPay
	result := repository.model(c).Model(&data).Where("ap_pay_no=? AND cust_id = ? AND is_del= ? ", apPayNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryApPayImpl) Update(c context.Context, apPayNo string, data model.ApPay) error {

	result := repository.model(c).Model(&data).Where("ap_pay_no=?", apPayNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryApPayImpl) DeleteDetailNotInIDs(c context.Context, apPayNo string, IDs []int64) error {
	var Details model.ApPayDet
	err := repository.model(c).Where("ap_pay_no=? AND ap_pay_det_id not in (?) ", apPayNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryApPayImpl) DeleteApPayMethodDetailNotInIDs(c context.Context, apPayNo string, IDs []int64) error {
	var Details model.ApPayMethod
	err := repository.model(c).Where("ap_pay_no=? AND ap_pay_method_id not in (?) ", apPayNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryApPayImpl) UpdateDetail(c context.Context, Details *model.ApPayDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryApPayImpl) UpdateApPayMethodDetail(c context.Context, Details *model.ApPayMethod) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
