package repository

import (
	"context"
	"errors"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"math"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryWhSoImpl struct {
		*gorm.DB
	}
)

func NewWhSoRepo(db *gorm.DB) *RepositoryWhSoImpl {
	return &RepositoryWhSoImpl{db}
}

type WhSoRepository interface {
	Store(c context.Context, data *model.WhSo) error
	StoreDetail(c context.Context, data *model.WhSoDet) error
	FindByNo(whSoNo string, custId string) (whAdj model.WhSoList, err error)
	FindDetail(whSoNo string, custId string) (Details []model.WhSoDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.WhSoList, int64, int, error)
	Delete(c context.Context, custId string, whSoNo string, deletedBy int64) error
	Update(c context.Context, whSoNo string, data model.WhSo) error
	DeleteDetailNotInIDs(c context.Context, whSoNo string, IDs []int64) error
	UpdateGrDetail(c context.Context, Details *model.WhSoDet) error
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryWhSoImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryWhSoImpl) Store(c context.Context, data *model.WhSo) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryWhSoImpl) StoreDetail(c context.Context, data *model.WhSoDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryWhSoImpl) FindByNo(whSoNo string, custId string) (whAdj model.WhSoList, err error) {
	err = repository.Select("inv.wh_so.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = inv.wh_so.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.wh_so.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.wh_so.wh_id AND wh.cust_id = ?", custId).
		Where("inv.wh_so.wh_so_no = ? AND inv.wh_so.cust_id=?", whSoNo, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryWhSoImpl) FindDetail(whSoNo string, custId string) (Details []model.WhSoDetRead, err error) {
	err = repository.Select("inv.wh_so_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = inv.wh_so_det.pro_id").
		Where("wh_so_no = ? AND inv.wh_so_det.cust_id=?", whSoNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryWhSoImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.WhSoList, int64, int, error) {
	var whSo []model.WhSoList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("wh_so_no")
	query := repository.Select("inv.wh_so.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = inv.wh_so.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.wh_so.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.wh_so.wh_id AND wh.cust_id = ?", custId)

	queryCount.Where("inv.wh_so.cust_id=?", custId)
	query.Where("inv.wh_so.cust_id=?", custId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.wh_so.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.wh_so.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("inv.wh_so.wh_so_no=?", dataFilter.Query)
		queryCount.Where("inv.wh_so.wh_so_no=?", dataFilter.Query)
	}

	if dataFilter.Sort != "" {

	} else {
		query.Order("wh_so_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&whSo).Error
	if err != nil {
		return whSo, total, 0, err
	}
	err = queryCount.Model(&whSo).Count(&total).Error
	if err != nil {
		return whSo, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return whSo, total, lastPage, nil
}
func (repository *RepositoryWhSoImpl) Delete(c context.Context, custId string, whSoNo string, deletedBy int64) error {
	var data model.WhSo
	result := repository.model(c).Model(&data).Where("wh_so_no=? AND cust_id = ? AND is_del= ? ", whSoNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryWhSoImpl) Update(c context.Context, whSoNo string, data model.WhSo) error {
	result := repository.model(c).Model(&data).Where("wh_so_no=?", whSoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryWhSoImpl) DeleteDetailNotInIDs(c context.Context, whSoNo string, IDs []int64) error {
	var Details model.WhSoDet
	err := repository.model(c).Where("wh_so_no=? AND wh_so_det_id not in (?) ", whSoNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryWhSoImpl) UpdateGrDetail(c context.Context, Details *model.WhSoDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
