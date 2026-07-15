package repository

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryVanUlImpl struct {
		*gorm.DB
	}
)

func NewVanUlRepo(db *gorm.DB) *RepositoryVanUlImpl {
	return &RepositoryVanUlImpl{db}
}

type VanUlRepository interface {
	Store(c context.Context, data *model.VanUl) error
	StoreDetail(c context.Context, data *model.VanUlDet) error
	FindByNo(vanSoNo string, custId string) (vanSo model.VanUlRead, err error)
	FindDetail(vanSoNo string, custId string) (Details []model.VanUlDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanUlRead, int64, int, error)
	Delete(c context.Context, custId string, vanUlNo string, deletedBy int64) error
	Update(c context.Context, vanUlNo string, data model.VanUl) error
	DeleteDetailNotInIDs(c context.Context, vanUlNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.VanUlDet) error
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryVanUlImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryVanUlImpl) Store(c context.Context, data *model.VanUl) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVanUlImpl) StoreDetail(c context.Context, data *model.VanUlDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryVanUlImpl) FindByNo(vanSoNo string, custId string) (vanSo model.VanUlRead, err error) {
	err = repository.Select("inv.van_ul.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wr.wh_code, wr.wh_name,s.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_ul.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_ul.closed_by").
		Joins("left join mst.m_warehouse wr on wr.wh_id = inv.van_ul.wh_id AND wr.cust_id = ?", custId).
		Joins("left join mst.m_salesman s on s.emp_id = inv.van_ul.salesman_id AND s.cust_id = ?", custId).
		Where("inv.van_ul.van_ul_no = ? AND inv.van_ul.cust_id=?", vanSoNo, custId).
		Take(&vanSo).Error
	return vanSo, err
}

func (repository *RepositoryVanUlImpl) FindDetail(vanSoNo string, custId string) (Details []model.VanUlDetRead, err error) {
	err = repository.Select("inv.van_ul_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = inv.van_ul_det.pro_id").
		Where("van_ul_no = ? AND inv.van_ul_det.cust_id=?", vanSoNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryVanUlImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanUlRead, int64, int, error) {
	var vanSo []model.VanUlRead
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("van_ul_no")
	query := repository.Select("inv.van_ul.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wr.wh_code, wr.wh_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_ul.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_ul.closed_by").
		Joins("left join mst.m_warehouse wr on wr.wh_id = inv.van_ul.wh_id AND wr.cust_id = ?", dataFilter.CustId)

	queryCount.Where("inv.van_ul.cust_id=?", dataFilter.CustId)
	query.Where("inv.van_ul.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.van_ul.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.van_ul.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("inv.van_ul.van_ul_no=?", dataFilter.Query)
		queryCount.Where("inv.van_ul.van_ul_no=?", dataFilter.Query)
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
		query.Order("van_ul_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&vanSo).Error
	if err != nil {
		return vanSo, total, 0, err
	}
	err = queryCount.Model(&vanSo).Count(&total).Error
	if err != nil {
		return vanSo, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return vanSo, total, lastPage, nil
}
func (repository *RepositoryVanUlImpl) Delete(c context.Context, custId string, vanUlNo string, deletedBy int64) error {
	var data model.VanUl
	result := repository.model(c).Model(&data).Where("van_ul_no=? AND cust_id = ? AND is_del= ? ", vanUlNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryVanUlImpl) Update(c context.Context, vanUlNo string, data model.VanUl) error {
	result := repository.model(c).Model(&data).Where("van_ul_no=?", vanUlNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryVanUlImpl) DeleteDetailNotInIDs(c context.Context, vanUlNo string, IDs []int64) error {
	var Details model.VanUlDet
	err := repository.model(c).Where("van_ul_no=? AND van_ul_det_id not in (?) ", vanUlNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryVanUlImpl) UpdateDetail(c context.Context, Details *model.VanUlDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
