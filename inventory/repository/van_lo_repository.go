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
	RepositoryVanLoImpl struct {
		*gorm.DB
	}
)

func NewVanLoRepo(db *gorm.DB) *RepositoryVanLoImpl {
	return &RepositoryVanLoImpl{db}
}

type VanLoRepository interface {
	Store(c context.Context, data *model.VanLo) error
	StoreDetail(c context.Context, data *model.VanLoDet) error
	FindByNo(vanLoNo string, custId string) (vanLo model.VanLoRead, err error)
	FindDetail(vanSoNo string, custId string) (Details []model.VanLoDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanLoRead, int64, int, error)
	Delete(c context.Context, custId string, vanLoNo string, deletedBy int64) error
	Update(c context.Context, vanLoNo string, data model.VanLo) error
	DeleteDetailNotInIDs(c context.Context, vanLoNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.VanLoDet) error
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryVanLoImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryVanLoImpl) Store(c context.Context, data *model.VanLo) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVanLoImpl) StoreDetail(c context.Context, data *model.VanLoDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryVanLoImpl) FindByNo(vanLoNo string, custId string) (vanLo model.VanLoRead, err error) {
	err = repository.Select("inv.van_lo.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, s.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_lo.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_lo.closed_by").
		Joins("left join mst.m_salesman s on s.emp_id = inv.van_lo.salesman_id AND s.cust_id = ?", custId).
		Where("inv.van_lo.van_lo_no = ? AND inv.van_lo.cust_id=?", vanLoNo, custId).
		Take(&vanLo).Error
	return vanLo, err
}
func (repository *RepositoryVanLoImpl) FindDetail(vanLoNo string, custId string) (Details []model.VanLoDetRead, err error) {
	err = repository.Select("inv.van_lo_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = inv.van_lo_det.pro_id").
		Where("inv.van_lo_det.van_lo_no = ? AND inv.van_lo_det.cust_id=?", vanLoNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryVanLoImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanLoRead, int64, int, error) {
	var vanSo []model.VanLoRead
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("van_lo_no")
	query := repository.Select("inv.van_lo.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name,s.sales_name as salesman_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_lo.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_lo.closed_by").
		Joins("left join mst.m_salesman s on s.emp_id = inv.van_lo.salesman_id AND s.cust_id = ?", dataFilter.CustId)

	queryCount.Where("inv.van_lo.cust_id=?", dataFilter.CustId)
	query.Where("inv.van_lo.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.van_lo.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.van_lo.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("inv.van_lo.van_lo_no=?", dataFilter.Query)
		queryCount.Where("inv.van_lo.van_lo_no=?", dataFilter.Query)
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
		query.Order("van_lo_no DESC")
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
func (repository *RepositoryVanLoImpl) Delete(c context.Context, custId string, vanLoNo string, deletedBy int64) error {
	var data model.VanLo
	result := repository.model(c).Model(&data).Where("van_lo_no=? AND cust_id = ? AND is_del= ? ", vanLoNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryVanLoImpl) Update(c context.Context, vanLoNo string, data model.VanLo) error {
	result := repository.model(c).Model(&data).Where("van_lo_no=?", vanLoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryVanLoImpl) DeleteDetailNotInIDs(c context.Context, vanLoNo string, IDs []int64) error {
	var Details model.VanLoDet
	err := repository.model(c).Where("van_lo_no=? AND van_lo_det_id not in (?) ", vanLoNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryVanLoImpl) UpdateDetail(c context.Context, Details *model.VanLoDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
