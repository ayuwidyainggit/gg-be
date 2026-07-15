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
	RepositoryWhTrfImpl struct {
		*gorm.DB
	}
)
type WhTrfRepository interface {
	Store(c context.Context, data *model.WhTrf) error
	StoreDetail(c context.Context, data *model.WhTrfDet) error
	FindByNo(whTrfNo string, custId string) (whTrf model.WhTrfList, err error)
	FindWhTrfdetail(whTrfNo string, custId, parentCustId string, distributorID int64) (Details []model.WhTrfDetRead, err error)
	FindAllByCustId(dataFilter entity.WhQueryFilter, custId string) ([]model.WhTrfList, int64, int, error)
	Delete(c context.Context, custId string, smpIsNo string, deletedBy int64) error
	Update(c context.Context, whTrfNo string, custID string, data model.WhTrf) error
	DeleteDetailNotInIDs(c context.Context, whTrfNo string, custID string, IDs []int) error
	UpdateGrDetail(c context.Context, Details *model.WhTrfDet) error
	FindProductByListID(productIDs []int64) (products []model.Product, err error)
	FindWarehouse(dataFilter entity.StockTranferWarehouseQueryFilter, custId string) ([]model.WarehouseStockTransfer, int64, int, error)
}

func NewWhTrfRepo(db *gorm.DB) *RepositoryWhTrfImpl {
	return &RepositoryWhTrfImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryWhTrfImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryWhTrfImpl) Store(c context.Context, data *model.WhTrf) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryWhTrfImpl) StoreDetail(c context.Context, data *model.WhTrfDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryWhTrfImpl) FindByNo(whTrfNo string, custId string) (whTrf model.WhTrfList, err error) {
	err = repository.Select("wh_trf.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, whfrom.wh_code as wh_code_from, whfrom.wh_name as wh_name_from, whto.wh_code as wh_code_to, whto.wh_name as wh_name_to").
		Joins("left join sys.m_user us on us.user_id = wh_trf.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = wh_trf.closed_by").
		Joins("left join mst.m_warehouse whfrom on whfrom.wh_id = wh_trf.wh_id_from AND whfrom.cust_id = ?", custId).
		Joins("left join mst.m_warehouse whto on whto.wh_id = wh_trf.wh_id_to AND whto.cust_id = ?", custId).
		Where("wh_trf.wh_trf_no = ? AND wh_trf.cust_id=?", whTrfNo, custId).
		Take(&whTrf).Error
	return whTrf, err
}

func (repository *RepositoryWhTrfImpl) FindWhTrfdetail(whTrfNo string, custId, parentCustId string, distributorID int64) (Details []model.WhTrfDetRead, err error) {

	err = repository.Select(`inv.wh_trf_det.*, p.pro_code, p.pro_name, p.unit_id1, p.unit_id2, p.unit_id3, p.conv_unit2, p.conv_unit3, p.vat, p.vat_bg, p.vat_lg_purch,
		p.purch_price1,
		p.purch_price2,
		p.purch_price3,
		p.sell_price1,
		p.sell_price2,
		p.sell_price3`).
		Joins("left join mst.m_product p on p.pro_id = inv.wh_trf_det.pro_id").
		Joins("LEFT JOIN smc.m_customer cus ON cus.cust_id = ?", custId).
		Where("wh_trf_no = ? AND inv.wh_trf_det.cust_id=?", whTrfNo, custId).Order("seq_no ASC").
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryWhTrfImpl) FindAllByCustId(dataFilter entity.WhQueryFilter, custId string) ([]model.WhTrfList, int64, int, error) {
	var whTrf []model.WhTrfList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("wh_trf_no")
	query := repository.Select("wh_trf.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, whfrom.wh_code as wh_code_from, whfrom.wh_name as wh_name_from, whto.wh_code as wh_code_to, whto.wh_name as wh_name_to").
		Joins("left join sys.m_user us on us.user_id = wh_trf.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = wh_trf.closed_by").
		Joins("left join mst.m_warehouse whfrom on whfrom.wh_id = wh_trf.wh_id_from AND whfrom.cust_id = ?", custId).
		Joins("left join mst.m_warehouse whto on whto.wh_id = wh_trf.wh_id_to AND whto.cust_id = ?", custId)

	queryCount.Where("wh_trf.cust_id=?", custId)
	query.Where("wh_trf.cust_id=?", custId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("wh_trf.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("wh_trf.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.StockTransferNo != "" {
		query.Where("wh_trf.wh_trf_no=?", dataFilter.StockTransferNo)
		queryCount.Where("wh_trf.wh_trf_no=?", dataFilter.StockTransferNo)
	}

	if len(dataFilter.WhIdFrom) > 0 {
		query.Where("wh_trf.wh_id_from in ?", dataFilter.WhIdFrom)
		queryCount.Where("wh_trf.wh_id_from in ?", dataFilter.WhIdFrom)
	}

	if len(dataFilter.WhIdTo) > 0 {
		query.Where("wh_trf.wh_id_to in ?", dataFilter.WhIdTo)
		queryCount.Where("wh_trf.wh_id_to in ?", dataFilter.WhIdTo)
	}

	if dataFilter.Sort != "" {

	} else {
		query.Order("wh_trf_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&whTrf).Error
	if err != nil {
		return whTrf, total, 0, err
	}
	err = queryCount.Model(&whTrf).Count(&total).Error
	if err != nil {
		return whTrf, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return whTrf, total, lastPage, nil

}
func (repository *RepositoryWhTrfImpl) Delete(c context.Context, custId string, whTrfNo string, deletedBy int64) error {
	var data model.WhTrf
	result := repository.model(c).Model(&data).Where("wh_trf_no=? AND cust_id = ? AND is_del= ? ", whTrfNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryWhTrfImpl) Update(c context.Context, whTrfNo string, custID string, data model.WhTrf) error {
	result := repository.model(c).Model(&data).Where("wh_trf_no=? AND cust_id = ?", whTrfNo, custID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryWhTrfImpl) DeleteDetailNotInIDs(c context.Context, whTrfNo string, custID string, IDs []int) error {
	var Details model.WhTrfDet
	err := repository.model(c).Where("wh_trf_no=? AND cust_id = ? AND wh_trf_det_id not in (?) ", whTrfNo, custID, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryWhTrfImpl) UpdateGrDetail(c context.Context, Details *model.WhTrfDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryWhTrfImpl) FindProductByListID(productIDs []int64) (products []model.Product, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryWhTrfImpl) FindWarehouse(dataFilter entity.StockTranferWarehouseQueryFilter, custId string) ([]model.WarehouseStockTransfer, int64, int, error) {
	var warehouses []model.WarehouseStockTransfer
	var total int64
	var limit int

	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("wh.wh_id").Where("wh_trf.cust_id=? ", custId)

	query := repository.
		Select("wh.wh_id, wh.wh_code, wh_name")

	if dataFilter.Method == "from" {
		query.Joins("LEFT JOIN mst.m_warehouse wh ON wh_trf.wh_id_from = wh.wh_id AND wh.cust_id = ?", custId)
		queryCount.Joins("LEFT JOIN mst.m_warehouse wh ON wh_trf.wh_id_from = wh.wh_id AND wh.cust_id = ?", custId)
	} else {
		query.Joins("LEFT JOIN mst.m_warehouse wh ON wh_trf.wh_id_to = wh.wh_id AND wh.cust_id = ?", custId)
		queryCount.Joins("LEFT JOIN mst.m_warehouse wh ON wh_trf.wh_id_to = wh.wh_id AND wh.cust_id = ?", custId)

	}

	query.Where("wh_trf.cust_id = ?", custId)

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		query.Where("wh_trf.trf_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
		queryCount.Where("wh_trf.trf_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit
	err := query.Group("wh.wh_id, wh.wh_code, wh_name").Limit(limit).Offset(offset).Find(&warehouses).Error
	if err != nil {
		return warehouses, total, 0, err
	}
	err = queryCount.Model(&warehouses).Group("wh.wh_id").Count(&total).Error
	if err != nil {
		return warehouses, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return warehouses, total, lastPage, nil
}
