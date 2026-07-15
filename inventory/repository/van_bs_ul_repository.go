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
	RepositoryVanBsUlImpl struct {
		*gorm.DB
	}
)

func NewVanBsUlRepo(db *gorm.DB) *RepositoryVanBsUlImpl {
	return &RepositoryVanBsUlImpl{db}
}

type VanBsUlRepository interface {
	Store(c context.Context, data *model.VanBsUl) error
	StoreDetail(c context.Context, data *model.VanBsUlDet) error
	FindByNo(gdsNo string, custId string) (whAdj model.VanBsUlList, err error)
	FindDetail(gdsNo string, custId string) (Details []model.VanBsUlDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanBsUlList, int64, int, error)
	Delete(c context.Context, custId string, vanBsUlNo string, deletedBy int64) error
	Update(c context.Context, vanBsUlNo string, data model.VanBsUl) error
	DeleteDetailNotInIDs(c context.Context, vanBsUlNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.VanBsUlDet) error
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryVanBsUlImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryVanBsUlImpl) Store(c context.Context, data *model.VanBsUl) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVanBsUlImpl) StoreDetail(c context.Context, data *model.VanBsUlDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryVanBsUlImpl) FindByNo(vanBsUlNo string, custId string) (whAdj model.VanBsUlList, err error) {
	err = repository.Select("inv.van_bs_ul.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name, sls.emp_id as salesman_code, sls.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_bs_ul.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_bs_ul.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.van_bs_ul.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_salesman sls on sls.emp_id = inv.van_bs_ul.salesman_id AND sls.cust_id = ?", custId).
		Where("inv.van_bs_ul.van_bs_ul_no = ? AND inv.van_bs_ul.cust_id=?", vanBsUlNo, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryVanBsUlImpl) FindDetail(vanBsUlNo string, custId string) (Details []model.VanBsUlDetRead, err error) {
	err = repository.Select("inv.van_bs_ul_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = inv.van_bs_ul_det.pro_id").
		Where("van_bs_ul_no = ? AND inv.van_bs_ul_det.cust_id=?", vanBsUlNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryVanBsUlImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanBsUlList, int64, int, error) {
	var vanBsUl []model.VanBsUlList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("van_bs_ul_no")
	query := repository.Select("inv.van_bs_ul.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name, sls.emp_id as salesman_code, sls.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_bs_ul.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_bs_ul.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.van_bs_ul.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_salesman sls on sls.emp_id = inv.van_bs_ul.salesman_id AND sls.cust_id = ?", dataFilter.CustId)

	queryCount.Where("inv.van_bs_ul.cust_id=?", dataFilter.CustId)
	query.Where("inv.van_bs_ul.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.van_bs_ul.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.van_bs_ul.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("inv.van_bs_ul.van_bs_ul_no=?", dataFilter.Query)
		queryCount.Where("inv.van_bs_ul.van_bs_ul_no=?", dataFilter.Query)
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
		query.Order("van_bs_ul_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&vanBsUl).Error
	if err != nil {
		return vanBsUl, total, 0, err
	}
	err = queryCount.Model(&vanBsUl).Count(&total).Error
	if err != nil {
		return vanBsUl, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return vanBsUl, total, lastPage, nil
}

func (repository *RepositoryVanBsUlImpl) Delete(c context.Context, custId string, vanBsUlNo string, deletedBy int64) error {
	var data model.VanBsUl
	result := repository.model(c).Model(&data).Where("van_bs_ul_no=? AND cust_id = ? AND is_del= ? ", vanBsUlNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryVanBsUlImpl) Update(c context.Context, vanBsUlNo string, data model.VanBsUl) error {
	result := repository.model(c).Model(&data).Where("van_bs_ul_no=?", vanBsUlNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryVanBsUlImpl) DeleteDetailNotInIDs(c context.Context, vanBsUlNo string, IDs []int64) error {
	var Details model.VanBsUlDet
	err := repository.model(c).Where("van_bs_ul_no=? AND van_bs_ul_det_id not in (?) ", vanBsUlNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryVanBsUlImpl) UpdateDetail(c context.Context, Details *model.VanBsUlDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
