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
	RepositoryWhAdjImpl struct {
		*gorm.DB
	}
)

type WhAdjRepository interface {
	Store(c context.Context, data *model.WhAdj) error
	StoreDetail(c context.Context, data *model.WhAdjDet) error
	FindByNo(adjNo, custId, langId string) (whAdj model.WhAdjList, err error)
	FindDetail(adjNo, custId string) (Details []model.WhAdjDetRead, err error)
	FindAllByCustId(dataFilter entity.WhAdjQueryFilter, custId, langId string) ([]model.WhAdjList, int64, int, error)
	Delete(c context.Context, custId string, adjNo string, deletedBy int64) error
	Update(c context.Context, adjNo string, data model.WhAdj) error
	DeleteDetailNotInIDs(c context.Context, adjNo string, IDs []int) error
	UpdateGrDetail(c context.Context, Details *model.WhAdjDet) error
	FindWarehouseStockAdjusment(dataFilter entity.WhAdjWarehouseQueryFilter, custId string) ([]model.WarehouseAdjustment, int64, int, error)
	FindProductByListID(productIDs []int64) (products []model.Product, err error)
	UpdateStatus(c context.Context, adjNo string, custId string, status int) error
}

func NewWhAdjRepo(db *gorm.DB) *RepositoryWhAdjImpl {
	return &RepositoryWhAdjImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryWhAdjImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryWhAdjImpl) Store(c context.Context, data *model.WhAdj) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryWhAdjImpl) StoreDetail(c context.Context, data *model.WhAdjDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryWhAdjImpl) FindByNo(adjNo, custId, langId string) (whAdj model.WhAdjList, err error) {
	err = repository.Select("inv.wh_adj.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name, wh.stock_type, st.status_name AS item_cdn_name").
		Joins("left join sys.m_user us on us.user_id = inv.wh_adj.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.wh_adj.closed_by").
		Joins("left join mst.m_status st on st.status_id = 'item_cdn' AND st.status_value = inv.wh_adj.item_cdn AND st.lang_id = '"+langId+"'").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.wh_adj.wh_id AND wh.cust_id = ?", custId).
		Where("inv.wh_adj.adj_no = ? AND inv.wh_adj.cust_id=?", adjNo, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryWhAdjImpl) FindDetail(adjNo, custId string) (Details []model.WhAdjDetRead, err error) {
	err = repository.Select("inv.wh_adj_det.*, p.pro_code, p.pro_name, p.unit_id1, p.unit_id2, p.unit_id3, p.conv_unit2, p.conv_unit3").
		Joins("left join mst.m_product p on p.pro_id = inv.wh_adj_det.pro_id").
		Where("adj_no = ? AND inv.wh_adj_det.cust_id=?", adjNo, custId).Order("seq_no ASC").
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryWhAdjImpl) FindAllByCustId(dataFilter entity.WhAdjQueryFilter, custId, langId string) ([]model.WhAdjList, int64, int, error) {
	var whAdj []model.WhAdjList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("adj_no")
	query := repository.Select("inv.wh_adj.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name, st.status_name AS item_cdn_name").
		Joins("left join sys.m_user us on us.user_id = inv.wh_adj.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = inv.wh_adj.closed_by").
		Joins("left join mst.m_status st on st.status_id = 'item_cdn' AND st.status_value = inv.wh_adj.item_cdn AND st.lang_id = '"+langId+"'").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.wh_adj.wh_id AND wh.cust_id = ?", custId)

	queryCount.Where("inv.wh_adj.cust_id=?", custId)
	query.Where("inv.wh_adj.cust_id=?", custId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.wh_adj.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.wh_adj.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.AdjusmentNo != "" {
		query.Where("inv.wh_adj.adj_no=?", dataFilter.AdjusmentNo)
		queryCount.Where("inv.wh_adj.adj_no=?", dataFilter.AdjusmentNo)
	}

	if len(dataFilter.WhID) != 0 {
		query.Where("inv.wh_adj.wh_id in ?", dataFilter.WhID)
		queryCount.Where("inv.wh_adj.wh_id in ?", dataFilter.WhID)
	}

	if len(dataFilter.Status) != 0 {
		query.Where("inv.wh_adj.data_status in ?", dataFilter.Status)
		queryCount.Where("inv.wh_adj.data_status in ?", dataFilter.Status)
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
		query.Order("adj_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&whAdj).Error
	if err != nil {
		return whAdj, total, 0, err
	}
	err = queryCount.Model(&whAdj).Count(&total).Error
	if err != nil {
		return whAdj, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return whAdj, total, lastPage, nil

}
func (repository *RepositoryWhAdjImpl) Delete(c context.Context, custId string, adjNo string, deletedBy int64) error {
	var data model.WhAdj
	result := repository.model(c).Model(&data).Where("adj_no=? AND cust_id = ? AND is_del= ? ", adjNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryWhAdjImpl) Update(c context.Context, adjNo string, data model.WhAdj) error {
	result := repository.model(c).Model(&data).Where("adj_no=?", adjNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryWhAdjImpl) DeleteDetailNotInIDs(c context.Context, adjNo string, IDs []int) error {
	var Details model.WhAdjDet
	err := repository.model(c).Where("adj_no=? AND wh_adj_det_id not in (?) ", adjNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryWhAdjImpl) UpdateGrDetail(c context.Context, Details *model.WhAdjDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryWhAdjImpl) FindWarehouseStockAdjusment(dataFilter entity.WhAdjWarehouseQueryFilter, custId string) ([]model.WarehouseAdjustment, int64, int, error) {
	var adjusmentWarehouses []model.WarehouseAdjustment
	var total int64
	var limit int

	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("wh_id").Where("cust_id=? ", custId)

	query := repository.
		Select("wh.wh_id, wh.wh_code, wh_name").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh_adj.wh_id = wh.wh_id AND wh.cust_id = ?", custId)
	query.Where("wh_adj.cust_id = ?", custId)

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		query.Where("wh.adj_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
		queryCount.Where("wh.wh_adj_datecode between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit
	err := query.Group("wh.wh_id, wh.wh_code, wh_name").Limit(limit).Offset(offset).Find(&adjusmentWarehouses).Error
	if err != nil {
		return adjusmentWarehouses, total, 0, err
	}
	err = queryCount.Model(&adjusmentWarehouses).Group("wh_id").Count(&total).Error
	if err != nil {
		return adjusmentWarehouses, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return adjusmentWarehouses, total, lastPage, nil
}

func (repository *RepositoryWhAdjImpl) FindProductByListID(productIDs []int64) (products []model.Product, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryWhAdjImpl) UpdateStatus(c context.Context, adjNo string, custId string, status int) error {
	var data model.WhAdj
	result := repository.model(c).Model(&data).Where("adj_no=? AND cust_id = ?", adjNo, custId).
		Updates(map[string]interface{}{"data_status": status})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
