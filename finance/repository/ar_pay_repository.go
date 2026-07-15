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
	RepositoryArPayImpl struct {
		*gorm.DB
	}
)

type ArPayRepository interface {
	Store(c context.Context, data *model.ArPay) error
	StoreDetail(c context.Context, data *model.ArPayDet) error
	FindByNo(arPayNo string, custId string) (whAdj model.ArPayList, err error)
	FindDetail(arPayNo string, custId string) (Details []model.ArPayDet, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ArPayList, int64, int, error)
	Delete(c context.Context, custId string, arPayNo string, deletedBy int64) error
	Update(c context.Context, arPayNo string, data model.ArPay) error
	DeleteDetailNotInIDs(c context.Context, arPayNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.ArPayDet) error
}

func NewArPayRepo(db *gorm.DB) *RepositoryArPayImpl {
	return &RepositoryArPayImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryArPayImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryArPayImpl) Store(c context.Context, data *model.ArPay) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryArPayImpl) StoreDetail(c context.Context, data *model.ArPayDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryArPayImpl) FindByNo(arPayNo string, custId string) (whAdj model.ArPayList, err error) {
	err = repository.Select("acf.ar_pay.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = acf.ar_pay.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.ar_pay.salesman_id AND sls.cust_id = ?", custId).
		Where("acf.ar_pay.ar_pay_no = ? AND acf.ar_pay.cust_id=?", arPayNo, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryArPayImpl) FindDetail(arPayNo string, custId string) (Details []model.ArPayDet, err error) {
	err = repository.
		Where("ar_pay_no = ? AND cust_id=?", arPayNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryArPayImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ArPayList, int64, int, error) {
	var arPay []model.ArPayList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ar_pay_no")
	query := repository.Select("acf.ar_pay.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = acf.ar_pay.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.ar_pay.salesman_id AND sls.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.ar_pay.cust_id=?", dataFilter.CustId)
	query.Where("acf.ar_pay.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.ar_pay.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.ar_pay.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.ar_pay.ar_pay_no=?", dataFilter.Query)
		query.Where("acf.ar_pay.ar_pay_no=?", dataFilter.Query)
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
		query.Order("ar_pay_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&arPay).Error
	if err != nil {
		return arPay, total, 0, err
	}
	err = queryCount.Model(&arPay).Count(&total).Error
	if err != nil {
		return arPay, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return arPay, total, lastPage, nil
}
func (repository *RepositoryArPayImpl) Delete(c context.Context, custId string, arPayNo string, deletedBy int64) error {
	var data model.ArPay
	result := repository.model(c).Model(&data).Where("ar_pay_no=? AND cust_id = ? AND is_del= ? ", arPayNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryArPayImpl) Update(c context.Context, arPayNo string, data model.ArPay) error {

	result := repository.model(c).Model(&data).Where("ar_pay_no=?", arPayNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (repository *RepositoryArPayImpl) DeleteDetailNotInIDs(c context.Context, arPayNo string, IDs []int64) error {
	var Details model.ArPayDet
	err := repository.model(c).Where("ar_pay_no=? AND ar_pay_det_id not in (?) ", arPayNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryArPayImpl) UpdateDetail(c context.Context, Details *model.ArPayDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
