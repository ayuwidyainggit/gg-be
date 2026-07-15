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
	RepositoryItemStChImpl struct {
		*gorm.DB
	}
)

func NewItemStChRepo(db *gorm.DB) *RepositoryItemStChImpl {
	return &RepositoryItemStChImpl{db}
}

type ItemStChRepository interface {
	Store(c context.Context, data *model.ItemStCh) error
	StoreDetail(c context.Context, data *model.ItemStChDet) error
	FindByNo(isc_no string, custId string) (itemStCh model.ItemStChList, err error)
	FindSmpIssuedetail(isc_no string, custId string) (Details []model.ItemStChDetResponse, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ItemStChList, int64, int, error)
	Delete(c context.Context, custId string, isc_no string, deletedBy int64) error
	Update(c context.Context, isc_no string, data model.ItemStCh) error
	DeleteDetailNotInIDs(c context.Context, isc_no string, IDs []int) error
	UpdateGrDetail(c context.Context, Details *model.ItemStChDet) error
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryItemStChImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryItemStChImpl) Store(c context.Context, data *model.ItemStCh) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryItemStChImpl) StoreDetail(c context.Context, data *model.ItemStChDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryItemStChImpl) FindByNo(isc_no string, custId string) (itemStCh model.ItemStChList, err error) {
	err = repository.
		Select("inv.item_st_ch.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = inv.item_st_ch.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.item_st_ch.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.item_st_ch.wh_id AND wh.cust_id = ?", custId).
		Where("inv.item_st_ch.isc_no = ? AND inv.item_st_ch.cust_id=?", isc_no, custId).
		Take(&itemStCh).Error
	return itemStCh, err
}

func (repository *RepositoryItemStChImpl) FindSmpIssuedetail(isc_no string, custId string) (Details []model.ItemStChDetResponse, err error) {
	err = repository.Select("item_st_ch_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = item_st_ch_det.pro_id").
		Where("inv.item_st_ch_det.isc_no = ? AND inv.item_st_ch_det.cust_id=?", isc_no, custId).Order("inv.item_st_ch_det.seq_no ASC").
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryItemStChImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ItemStChList, int64, int, error) {
	var itemStCh []model.ItemStChList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("isc_no")
	query := repository.Select("inv.item_st_ch.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name").
		Joins("LEFT JOIN sys.m_user us on us.user_id = inv.item_st_ch.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.item_st_ch.closed_by").
		Joins("LEFT JOIN mst.m_warehouse wh on wh.wh_id = inv.item_st_ch.wh_id AND wh.cust_id = ?", custId)

	queryCount.Where("inv.item_st_ch.cust_id=?", custId)
	query.Where("inv.item_st_ch.cust_id=?", custId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.item_st_ch.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.item_st_ch.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("inv.item_st_ch.isc_no=?", dataFilter.Query)
		queryCount.Where("inv.item_st_ch.isc_no=?", dataFilter.Query)
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
		query.Order("isc_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&itemStCh).Error
	if err != nil {
		return itemStCh, total, 0, err
	}
	err = queryCount.Model(&itemStCh).Count(&total).Error
	if err != nil {
		return itemStCh, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return itemStCh, total, lastPage, nil
}

func (repository *RepositoryItemStChImpl) Delete(c context.Context, custId string, isc_no string, deletedBy int64) error {
	var data model.ItemStCh
	result := repository.model(c).Model(&data).Where("isc_no=? AND cust_id = ? AND is_del= ? ", isc_no, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryItemStChImpl) Update(c context.Context, isc_no string, data model.ItemStCh) error {
	result := repository.model(c).Model(&data).Where("isc_no=?", isc_no).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryItemStChImpl) DeleteDetailNotInIDs(c context.Context, isc_no string, IDs []int) error {
	var Details model.ItemStChDet
	err := repository.model(c).Where("isc_no=? AND isc_det_id not in (?) ", isc_no, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryItemStChImpl) UpdateGrDetail(c context.Context, Details *model.ItemStChDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
