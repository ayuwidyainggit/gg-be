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
	RepositoryRoImpl struct {
		*gorm.DB
	}
)
type RoRepository interface {
	Store(c context.Context, data *model.Ro) error
	StoreDetail(c context.Context, data *model.RoDet) error
	FindByNo(RoNo string, custId string) (realOrder model.RoList, err error)
	FindDetail(RoNo string, custId string) (details []model.RoDetRead, err error)
	FindAllByCustId(dataFilter entity.RoQueryFilter) ([]model.RoList, int64, int, error)

	Update(c context.Context, RoNo string, data model.Ro) error
	DeleteDetailNotInIDs(c context.Context, RoNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.RoDet) error
	Delete(c context.Context, custId string, RoNo string, deletedBy int64) error
}

func NewRoRepo(db *gorm.DB) *RepositoryRoImpl {
	return &RepositoryRoImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryRoImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryRoImpl) Store(c context.Context, data *model.Ro) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryRoImpl) StoreDetail(c context.Context, data *model.RoDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryRoImpl) FindByNo(roNo string, custId string) (realOrder model.RoList, err error) {
	err = repository.
		Select(`ro.*, 
			us.user_fullname AS updated_by_name,
			ot.outlet_code, ot.outlet_name,
			sls.sales_name,
			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.ro.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.ro.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.ro.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.ro.outlet_id AND ot.cust_id = ?", custId).
		Where("ro.ro_no = ? AND ro.cust_id=?", roNo, custId).
		Take(&realOrder).Error
	return realOrder, err
}

func (repository *RepositoryRoImpl) FindDetail(roNo string, custId string) (details []model.RoDetRead, err error) {
	err = repository.Select("sls.ro_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = sls.ro_det.pro_id").
		Where("ro_no = ? AND sls.ro_det.cust_id=?", roNo, custId).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryRoImpl) FindAllByCustId(dataFilter entity.RoQueryFilter) ([]model.RoList, int64, int, error) {
	var ro []model.RoList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ro_no")
	query := repository.Select(
		`ro.*, 
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, ot.outlet_name, 
			sls.sales_name,
			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.ro.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.ro.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.ro.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.ro.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.ro.cust_id=?", dataFilter.CustId)
	query.Where("sls.ro.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.ro.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.ro.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.ro.ro_no=?", dataFilter.Query)
		query.Where("sls.ro.ro_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.ro.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.ro.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.ro.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.ro.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.ro.data_status in ?", dataFilter.Status)
		query.Where("sls.ro.data_status in ?", dataFilter.Status)
	}

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
		query.Order("ro_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryRoImpl) Update(c context.Context, RoNo string, data model.Ro) error {
	result := repository.model(c).Model(&data).Where("ro_no=?", RoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryRoImpl) DeleteDetailNotInIDs(c context.Context, RoNo string, IDs []int64) error {
	var Details model.RoDet
	err := repository.model(c).Where("ro_no=? AND ro_det_id not in (?) ", RoNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryRoImpl) UpdateDetail(c context.Context, Details *model.RoDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryRoImpl) Delete(c context.Context, custId string, RoNo string, deletedBy int64) error {
	var data model.Ro
	result := repository.model(c).Model(&data).Where("ro_no=? AND cust_id = ? AND is_del= ? ", RoNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
