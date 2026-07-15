package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryTlsImpl struct {
		*gorm.DB
	}
)
type TlsRepository interface {
	Store(c context.Context, data *model.Tls) error
	StoreDetail(c context.Context, data *model.TlsDet) error
	FindByNo(TlsId int, custId string) (whAdj model.TlsList, err error)
	FindDetail(TlsId int, custId string) (Details []model.TlsDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.TlsList, int64, int, error)

	Update(c context.Context, TlsId int, data model.Tls) error
	DeleteDetailNotInIDs(c context.Context, TlsId int, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.TlsDet) error
	Delete(c context.Context, custId string, TlsId int, deletedBy int64) error
}

func NewTlsRepo(db *gorm.DB) *RepositoryTlsImpl {
	return &RepositoryTlsImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryTlsImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryTlsImpl) Store(c context.Context, data *model.Tls) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryTlsImpl) StoreDetail(c context.Context, data *model.TlsDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryTlsImpl) FindByNo(TlsId int, custId string) (whAdj model.TlsList, err error) {
	err = repository.Select("tls.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = sls.tls.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = tls.salesman_id AND sls.cust_id = ?", custId).
		Where("sls.tls.cust_id=? AND sls.tls.tls_id = ?", custId, TlsId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryTlsImpl) FindDetail(TlsId int, custId string) (Details []model.TlsDetRead, err error) {
	err = repository.Select("sls.tls_det.*, o.outlet_code, o.outlet_name").
		Joins("left join mst.m_outlet o on o.outlet_id = sls.tls_det.outlet_id AND o.cust_id = ?", custId).
		Where("sls.tls_det.cust_id=? AND tls_id = ?", custId, TlsId).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryTlsImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.TlsList, int64, int, error) {
	var tls []model.TlsList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("tls_id")
	query := repository.Select("tls.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = sls.tls.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = tls.salesman_id AND sls.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.tls.cust_id=?", dataFilter.CustId)
	query.Where("sls.tls.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.tls.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.tls.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.tls.tls_id=?", dataFilter.Query)
		query.Where("sls.tls.tls_id=?", dataFilter.Query)
	}

	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			// sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 			query.Order(colSort[0] + " " + colSort[1])
	// 		}
	// 	}
	// } else {
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
		query.Order("tls_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&tls).Error
	if err != nil {
		return tls, total, 0, err
	}
	err = queryCount.Model(&tls).Count(&total).Error
	if err != nil {
		return tls, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return tls, total, lastPage, nil
}

func (repository *RepositoryTlsImpl) Update(c context.Context, TlsId int, data model.Tls) error {
	result := repository.model(c).Model(&data).Where("tls_id=?", TlsId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryTlsImpl) DeleteDetailNotInIDs(c context.Context, TlsId int, IDs []int64) error {
	var Details model.TlsDet
	err := repository.model(c).Where("tls_id=? AND tls_det_id not in (?) ", TlsId, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryTlsImpl) UpdateDetail(c context.Context, Details *model.TlsDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryTlsImpl) Delete(c context.Context, custId string, TlsId int, deletedBy int64) error {
	var data model.Tls
	result := repository.model(c).Model(&data).Where("cust_id = ? AND tls_id=? AND is_del= ? ", custId, TlsId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
