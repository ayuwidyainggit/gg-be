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
	RepositoryVanSoImpl struct {
		*gorm.DB
	}
)

func NewVanSoRepo(db *gorm.DB) *RepositoryVanSoImpl {
	return &RepositoryVanSoImpl{db}
}

type VanSoRepository interface {
	Store(c context.Context, data *model.VanSo) error
	StoreDetail(c context.Context, data *model.VanSoDet) error
	FindByNo(vanSoNo string, custId string) (vanSo model.VanSoList, err error)
	FindDetail(vanSoNo string, custId string) (Details []model.VanSoDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanSoList, int64, int, error)
	Delete(c context.Context, custId string, vanSoNo string, deletedBy int64) error
	Update(c context.Context, whSoNo string, data model.VanSo) error
	DeleteDetailNotInIDs(c context.Context, whSoNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.VanSoDet) error
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryVanSoImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryVanSoImpl) Store(c context.Context, data *model.VanSo) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVanSoImpl) StoreDetail(c context.Context, data *model.VanSoDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryVanSoImpl) FindByNo(vanSoNo string, custId string) (vanSo model.VanSoList, err error) {
	err = repository.Select("inv.van_so.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, emp.emp_code, emp.emp_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_so.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_so.closed_by").
		Joins("left join mst.m_employee emp on emp.emp_id = inv.van_so.emp_id AND emp.cust_id = ?", custId).
		Where("inv.van_so.van_so_no = ? AND inv.van_so.cust_id=?", vanSoNo, custId).
		Take(&vanSo).Error
	return vanSo, err
}

func (repository *RepositoryVanSoImpl) FindDetail(vanSoNo string, custId string) (Details []model.VanSoDetRead, err error) {
	err = repository.Select("inv.van_so_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = inv.van_so_det.pro_id").
		Where("van_so_no = ? AND inv.van_so_det.cust_id=?", vanSoNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryVanSoImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.VanSoList, int64, int, error) {
	var vanSo []model.VanSoList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("van_so_no")
	query := repository.Select("inv.van_so.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, emp.emp_code, emp.emp_name").
		Joins("left join sys.m_user us on us.user_id = inv.van_so.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.van_so.closed_by").
		Joins("left join mst.m_employee emp on emp.emp_id = inv.van_so.emp_id AND emp.cust_id = ?", dataFilter.CustId)

	queryCount.Where("inv.van_so.cust_id=?", dataFilter.CustId)
	query.Where("inv.van_so.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.van_so.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.van_so.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("inv.van_so.van_so_no=?", dataFilter.Query)
		queryCount.Where("inv.van_so.van_so_no=?", dataFilter.Query)
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
		query.Order("van_so_no DESC")
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
func (repository *RepositoryVanSoImpl) Delete(c context.Context, custId string, vanSoNo string, deletedBy int64) error {
	var data model.VanSo
	result := repository.model(c).Model(&data).Where("van_so_no=? AND cust_id = ? AND is_del= ? ", vanSoNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryVanSoImpl) Update(c context.Context, whSoNo string, data model.VanSo) error {
	result := repository.model(c).Model(&data).Where("van_so_no=?", whSoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryVanSoImpl) DeleteDetailNotInIDs(c context.Context, whSoNo string, IDs []int64) error {
	var Details model.VanSoDet
	err := repository.model(c).Where("van_so_no=? AND van_so_det_id not in (?) ", whSoNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryVanSoImpl) UpdateDetail(c context.Context, Details *model.VanSoDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
