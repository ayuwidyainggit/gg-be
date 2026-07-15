package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/structs"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	RepositoryStockImpl struct {
		*gorm.DB
	}
)

type invoiceReleaseKey struct {
	CustID   string
	TrNo     string
	RefDetID int64
}

const FLOAT_COMPARE_EPSILON = 1e-6

type StockRepository interface {
	StockUpdates(c context.Context, stockUpdates []*entity.StockUpdate) error
	SalesStockUpdates(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error
	InvoiceSalesStockUpdates(c context.Context, stockUpdates []*entity.InvoiceSalesStockUpdate) error
	CancelSalesStockUpdates(c context.Context, orderNo string, stockDate time.Time, rows []entity.CancelStockWrite) error
	GetCancelStockBasis(c context.Context, custID string, orderNo string) ([]entity.CancelStockBasis, error)
	UpdateOnCustomerOrder(c context.Context, custId string, whId int64, proId int64, delta float64) error
	GetCurrentStock(c context.Context, custId string, whId int64, proId int64) (float64, error)
}

func NewStockRepo(db *gorm.DB) *RepositoryStockImpl {
	return &RepositoryStockImpl{db}
}

func isZeroDelta(delta float64) bool {
	return math.Abs(delta) < FLOAT_COMPARE_EPSILON
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryStockImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func buildInvoiceReleaseKeys(stockUpdates []*entity.InvoiceSalesStockUpdate) ([]invoiceReleaseKey, string, error) {
	if len(stockUpdates) == 0 {
		return nil, "", nil
	}

	releaseKeys := make([]invoiceReleaseKey, 0, len(stockUpdates))
	custID := stockUpdates[0].CustID
	for _, su := range stockUpdates {
		if su.CustID != custID {
			return nil, "", errors.New("InvoiceSalesStockUpdates received mixed cust_id in a single call")
		}
		releaseKeys = append(releaseKeys, invoiceReleaseKey{
			CustID:   su.CustID,
			TrNo:     su.TrNo + "-" + "CO",
			RefDetID: su.RefDetId,
		})
	}

	return releaseKeys, custID, nil
}

func filterDuplicateInvoiceStockUpdates(stockUpdates []*entity.InvoiceSalesStockUpdate, releaseKeys []invoiceReleaseKey, existing map[invoiceReleaseKey]struct{}) ([]*entity.InvoiceSalesStockUpdate, []invoiceReleaseKey, int) {
	filteredUpdates := make([]*entity.InvoiceSalesStockUpdate, 0, len(stockUpdates))
	filteredKeys := make([]invoiceReleaseKey, 0, len(releaseKeys))
	skipped := 0
	for i, stockUpdate := range stockUpdates {
		key := releaseKeys[i]
		if _, dup := existing[key]; dup {
			skipped++
			continue
		}
		filteredUpdates = append(filteredUpdates, stockUpdate)
		filteredKeys = append(filteredKeys, key)
	}
	return filteredUpdates, filteredKeys, skipped
}

func (repository *RepositoryStockImpl) InvoiceSalesStockUpdates(c context.Context, stockUpdates []*entity.InvoiceSalesStockUpdate) error {
	if len(stockUpdates) == 0 {
		return nil
	}

	releaseKeys, custID, err := buildInvoiceReleaseKeys(stockUpdates)
	if err != nil {
		return err
	}

	existing := make(map[invoiceReleaseKey]struct{}, len(releaseKeys))
	releaseTrNo := make([]string, 0, len(releaseKeys))
	releaseRefDets := make([]int64, 0, len(releaseKeys))
	for _, k := range releaseKeys {
		releaseTrNo = append(releaseTrNo, k.TrNo)
		releaseRefDets = append(releaseRefDets, k.RefDetID)
	}

	type priorReleaseRow struct {
		TrNo     string
		RefDetID int64
	}
	var prior []priorReleaseRow
	err = repository.model(c).
		Table("inv.stock").
		Select("tr_no, ref_det_id").
		Where("cust_id = ? AND tr_code = 'CO' AND tr_no IN ? AND ref_det_id IN ? AND qty_out_order > 0", custID, releaseTrNo, releaseRefDets).
		Scan(&prior).Error
	if err != nil {
		return err
	}
	for _, row := range prior {
		existing[invoiceReleaseKey{CustID: custID, TrNo: row.TrNo, RefDetID: row.RefDetID}] = struct{}{}
	}

	filteredUpdates, filteredKeys, skipped := filterDuplicateInvoiceStockUpdates(stockUpdates, releaseKeys, existing)

	var whStocks []model.WarehouseStock
	var stocks []*model.Stock
	mergedWhStocks := make(map[string]*model.WarehouseStock)
	for i, stockUpdate := range filteredUpdates {
		key := filteredKeys[i]

		whKey := fmt.Sprintf("%v-%v", stockUpdate.ProID, stockUpdate.WhID)
		qty := 0 - stockUpdate.QtyOrderBefore
		if existingWh, ok := mergedWhStocks[whKey]; ok {
			existingWh.QtyOnOrder += qty
		} else {
			mergedWhStocks[whKey] = &model.WarehouseStock{
				CustID:     stockUpdate.CustID,
				WhID:       stockUpdate.WhID,
				ProID:      stockUpdate.ProID,
				QtyOnOrder: qty,
			}
		}

		stocks = append(stocks, &model.Stock{
			CustID:      stockUpdate.CustID,
			StockDate:   stockUpdate.StockDate,
			TrCode:      "CO",
			TrNo:        key.TrNo,
			WhID:        stockUpdate.WhID,
			ProID:       stockUpdate.ProID,
			ItemCdn:     1,
			UnitPrice:   stockUpdate.UnitPrice,
			RefDetId:    stockUpdate.RefDetId,
			QtyOutOrder: qty * -1,
		})
	}

	for _, whStock := range mergedWhStocks {
		whStocks = append(whStocks, *whStock)
	}
	if len(whStocks) == 0 && len(stocks) == 0 {
		return nil
	}
	log.Info("InvoiceSalesStockUpdates: incoming=", len(stockUpdates), " skipped=", skipped, " wh_deltas=", structs.StructToJson(whStocks))

	if len(whStocks) > 0 {
		if err := repository.UpsertWithExistingValueArr(c, whStocks); err != nil {
			return err
		}
	}
	if len(stocks) > 0 {
		if err := repository.StoreBulk(c, stocks); err != nil {
			return err
		}
	}

	return nil
}

func (repository *RepositoryStockImpl) SalesStockUpdates(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
	filteredUpdates := make([]*entity.SalesOrderStockUpdate, 0, len(stockUpdates))
	for _, stockUpdate := range stockUpdates {
		if stockUpdate == nil {
			continue
		}

		effectiveDelta := stockUpdate.QtyOrder
		if stockUpdate.QtyOrderBefore != nil {
			effectiveDelta = stockUpdate.QtyOrder - *stockUpdate.QtyOrderBefore
		}

		if isZeroDelta(effectiveDelta) {
			continue
		}

		filteredUpdates = append(filteredUpdates, stockUpdate)
	}

	if len(filteredUpdates) == 0 {
		return nil
	}

	var whStocks []model.WarehouseStock
	var stocks []*model.Stock
	mergedWhStocks := make(map[string]*model.WarehouseStock)

	for _, stockUpdate := range filteredUpdates {
		key := fmt.Sprintf("%v-%v", stockUpdate.ProID, stockUpdate.WhID)
		qty := stockUpdate.QtyOrder
		qtyOnStock := stockUpdate.QtyOrder * -1
		qtyOutStockWarehouse := stockUpdate.QtyOrder
		var qtyInStockWarehouse float64
		if stockUpdate.QtyOrderBefore != nil {
			if *stockUpdate.QtyOrderBefore > stockUpdate.QtyOrder {
				qty = (*stockUpdate.QtyOrderBefore - stockUpdate.QtyOrder) * -1
				qtyOnStock = *stockUpdate.QtyOrderBefore - stockUpdate.QtyOrder

				qtyOutStockWarehouse = 0
				qtyInStockWarehouse = (*stockUpdate.QtyOrderBefore - stockUpdate.QtyOrder)
			} else {
				qty = stockUpdate.QtyOrder - *stockUpdate.QtyOrderBefore
				qtyOnStock = (stockUpdate.QtyOrder - *stockUpdate.QtyOrderBefore) * -1

				qtyOutStockWarehouse = (stockUpdate.QtyOrder - *stockUpdate.QtyOrderBefore)
				qtyInStockWarehouse = 0
			}

		}

		if existing, ok := mergedWhStocks[key]; ok {
			existing.Qty += qty
			existing.QtyOnOrder += qty
		} else {
			whStock := &model.WarehouseStock{
				CustID:     stockUpdate.CustID,
				WhID:       stockUpdate.WhID,
				ProID:      stockUpdate.ProID,
				QtyOnOrder: qty,
				Qty:        qtyOnStock,
			}
			mergedWhStocks[key] = whStock
		}

		stockWarehouse := model.Stock{
			CustID:    stockUpdate.CustID,
			StockDate: stockUpdate.StockDate,
			TrCode:    stockUpdate.TrCode,
			TrNo:      stockUpdate.TrNo,
			WhID:      stockUpdate.WhID,
			ProID:     stockUpdate.ProID,
			ItemCdn:   1,
			UnitPrice: stockUpdate.UnitPrice,
			RefDetId:  stockUpdate.RefDetId,
			QtyIn:     qtyInStockWarehouse,
			QtyOut:    qtyOutStockWarehouse,
		}

		stocks = append(stocks, &stockWarehouse)

		stock := model.Stock{
			CustID:      stockUpdate.CustID,
			StockDate:   stockUpdate.StockDate,
			TrCode:      "CO",
			TrNo:        stockUpdate.TrNo + "-" + "CO",
			WhID:        stockUpdate.WhID,
			ProID:       stockUpdate.ProID,
			ItemCdn:     1,
			UnitPrice:   stockUpdate.UnitPrice,
			RefDetId:    stockUpdate.RefDetId,
			QtyInOrder:  qtyOutStockWarehouse,
			QtyOutOrder: qtyInStockWarehouse,
		}

		stocks = append(stocks, &stock)
	}

	for _, whStock := range mergedWhStocks {
		whStocks = append(whStocks, *whStock)
	}
	log.Info("mergeWhStocks:", structs.StructToJson(whStocks))

	err := repository.UpsertWithExistingValueArr(c, whStocks)
	if err != nil {
		return err
	}

	err = repository.StoreBulk(c, stocks)
	if err != nil {
		return err
	}

	return nil
}

type cancelStockBaseRow struct {
	CustID       string
	WhID         int64
	ProID        int64
	RefDetID     int64
	QtyOutSO     float64
	UnitPrice    float64
	SourceTrNo   string
	SourceTrCode string
}

func buildCancelStockMutations(orderNo string, stockDate time.Time, rows []cancelStockBaseRow) ([]*model.Stock, []model.WarehouseStock) {
	stocks := make([]*model.Stock, 0)
	whDeltas := make([]model.WarehouseStock, 0)

	for _, row := range rows {
		if row.QtyOutSO <= 0 {
			continue
		}

		stockRowB := &model.Stock{
			CustID:      row.CustID,
			StockDate:   stockDate,
			TrCode:      "CO",
			TrNo:        orderNo + "-CO",
			WhID:        row.WhID,
			ProID:       row.ProID,
			ItemCdn:     1,
			QtyIn:       0,
			QtyOut:      0,
			QtyInOrder:  0,
			QtyOutOrder: row.QtyOutSO,
			UnitPrice:   row.UnitPrice,
			Cogs:        0,
			RefDetId:    row.RefDetID,
		}

		stocks = append(stocks, stockRowB)
		whDeltas = append(whDeltas, model.WarehouseStock{
			CustID:     row.CustID,
			WhID:       row.WhID,
			ProID:      row.ProID,
			Qty:        row.QtyOutSO,
			QtyOnOrder: -row.QtyOutSO,
		})
	}

	return stocks, whDeltas
}

func (repository *RepositoryStockImpl) cancelStockBasisQuery(c context.Context, custID string, orderNo string) *gorm.DB {
	cancelTrNo := orderNo + "-CO"

	sourceAgg := repository.model(c).
		Table("inv.stock s").
		Select(`
			s.cust_id,
			s.wh_id,
			s.pro_id,
			s.ref_det_id,
			MAX(s.stock_date) AS stock_date,
			MAX(s.unit_price) AS unit_price,
			SUM(s.qty_out - s.qty_in) AS qty_out_so
		`).
		Where("s.cust_id = ? AND s.tr_no = ? AND s.tr_code = 'SO'", custID, orderNo).
		Group("s.cust_id, s.wh_id, s.pro_id, s.ref_det_id")

	cancelAgg := repository.model(c).
		Table("inv.stock c").
		Select(`
			c.cust_id,
			c.wh_id,
			c.pro_id,
			c.ref_det_id,
			SUM(c.qty_out_order) AS qty_out_order_cancel
		`).
		Where("c.cust_id = ? AND c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO%'))", custID, cancelTrNo).
		Group("c.cust_id, c.wh_id, c.pro_id, c.ref_det_id")

	activeDetailAgg := repository.model(c).
		Table("sls.order_detail od").
		Select(`
			od.cust_id,
			od.ro_no,
			od.pro_id,
			COUNT(*) AS active_detail_count
		`).
		Where(`od.cust_id = ? AND od.ro_no = ? AND od.item_type = 1
			AND (
				COALESCE(od.qty1_final, 0) > 0
				OR COALESCE(od.qty1, 0) > 0
				OR COALESCE(od.qty_po1, 0) > 0
				OR COALESCE(od.qty2_final, 0) > 0
				OR COALESCE(od.qty2, 0) > 0
				OR COALESCE(od.qty_po2, 0) > 0
				OR COALESCE(od.qty3_final, 0) > 0
				OR COALESCE(od.qty3, 0) > 0
				OR COALESCE(od.qty_po3, 0) > 0
			)`, custID, orderNo).
		Group("od.cust_id, od.ro_no, od.pro_id")

	return repository.model(c).
		Table("sls.order_detail od").
		Joins("JOIN sls.order o ON o.cust_id = od.cust_id AND o.ro_no = od.ro_no").
		Joins("LEFT JOIN mst.m_product mp ON mp.cust_id = od.cust_id AND mp.pro_id = od.pro_id").
		Joins(`
			LEFT JOIN (?) s
				ON s.cust_id = od.cust_id
				AND s.wh_id = o.wh_id
				AND s.pro_id = od.pro_id
				AND s.ref_det_id = od.order_detail_id
		`, sourceAgg).
		Joins(`
			LEFT JOIN (?) c
				ON c.cust_id = od.cust_id
				AND c.wh_id = o.wh_id
				AND c.pro_id = od.pro_id
				AND c.ref_det_id = od.order_detail_id
		`, cancelAgg).
		Joins(`
			LEFT JOIN (?) ad
				ON ad.cust_id = od.cust_id
				AND ad.ro_no = od.ro_no
				AND ad.pro_id = od.pro_id
		`, activeDetailAgg).
		Where(`od.cust_id = ? AND od.ro_no = ?
			AND (
				COALESCE(od.qty1_final, 0) > 0
				OR COALESCE(od.qty1, 0) > 0
				OR COALESCE(od.qty_po1, 0) > 0
				OR COALESCE(od.qty2_final, 0) > 0
				OR COALESCE(od.qty2, 0) > 0
				OR COALESCE(od.qty_po2, 0) > 0
				OR COALESCE(od.qty3_final, 0) > 0
				OR COALESCE(od.qty3, 0) > 0
				OR COALESCE(od.qty_po3, 0) > 0
			)`, custID, orderNo).
		Select(`
			od.cust_id,
			o.wh_id,
			od.pro_id,
			od.order_detail_id AS ref_det_id,
			COALESCE(s.stock_date, o.ro_date) AS stock_date,
			COALESCE(s.unit_price, COALESCE(od.sell_price_final1, COALESCE(od.sell_price1, COALESCE(od.sell_price_po1, 0)))) AS unit_price,
			COALESCE(od.qty1_final, COALESCE(od.qty1, COALESCE(od.qty_po1, 0))) AS qty_final,
			COALESCE(s.qty_out_so, 0) AS qty_out_so,
			COALESCE(c.qty_out_order_cancel, 0) AS qty_out_order_cancel,
			(
				CASE
					WHEN COALESCE(s.qty_out_so, 0) > 0 THEN COALESCE(s.qty_out_so, 0)
					ELSE GREATEST(
						COALESCE(od.qty1_final, COALESCE(od.qty1, COALESCE(od.qty_po1, 0)))
						* GREATEST(COALESCE(od.conv_unit2, COALESCE(mp.conv_unit2, 1)), 1)
						* GREATEST(COALESCE(od.conv_unit3, COALESCE(mp.conv_unit3, 1)), 1)
						+ COALESCE(od.qty2_final, COALESCE(od.qty2, COALESCE(od.qty_po2, 0)))
						* GREATEST(COALESCE(od.conv_unit3, COALESCE(mp.conv_unit3, 1)), 1)
						+ COALESCE(od.qty3_final, COALESCE(od.qty3, COALESCE(od.qty_po3, 0)))
					, 0)
				END
				- COALESCE(c.qty_out_order_cancel, 0)
			) AS qty_outstanding,
			GREATEST(
				CASE
					WHEN COALESCE(s.qty_out_so, 0) > 0 THEN COALESCE(s.qty_out_so, 0)
					ELSE GREATEST(
						COALESCE(od.qty1_final, COALESCE(od.qty1, COALESCE(od.qty_po1, 0)))
						* GREATEST(COALESCE(od.conv_unit2, COALESCE(mp.conv_unit2, 1)), 1)
						* GREATEST(COALESCE(od.conv_unit3, COALESCE(mp.conv_unit3, 1)), 1)
						+ COALESCE(od.qty2_final, COALESCE(od.qty2, COALESCE(od.qty_po2, 0)))
						* GREATEST(COALESCE(od.conv_unit3, COALESCE(mp.conv_unit3, 1)), 1)
						+ COALESCE(od.qty3_final, COALESCE(od.qty3, COALESCE(od.qty_po3, 0)))
					, 0)
				END
				- COALESCE(c.qty_out_order_cancel, 0)
			, 0) AS qty_out_smallest,
			(COALESCE(s.qty_out_so, 0) <= 0) AS is_missing_source,
			(COALESCE(ad.active_detail_count, 0) > 1) AS is_ambiguous
		`)
}

func (repository *RepositoryStockImpl) GetCancelStockBasis(c context.Context, custID string, orderNo string) ([]entity.CancelStockBasis, error) {
	var rows []entity.CancelStockBasis
	err := repository.cancelStockBasisQuery(c, custID, orderNo).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repository *RepositoryStockImpl) CancelSalesStockUpdates(c context.Context, orderNo string, stockDate time.Time, rows []entity.CancelStockWrite) error {
	if len(rows) == 0 {
		return nil
	}

	baseRows := make([]cancelStockBaseRow, 0, len(rows))
	for _, row := range rows {
		baseRows = append(baseRows, cancelStockBaseRow{
			CustID:       row.CustID,
			WhID:         row.WhID,
			ProID:        row.ProID,
			RefDetID:     row.RefDetID,
			QtyOutSO:     row.QtySmallest,
			UnitPrice:    row.UnitPrice,
			SourceTrNo:   row.RoNo,
			SourceTrCode: "SO",
		})
	}

	stocks, whDeltas := buildCancelStockMutations(orderNo, stockDate, baseRows)
	if len(stocks) == 0 && len(whDeltas) == 0 {
		return nil
	}

	if len(whDeltas) > 0 {
		if err := repository.UpsertWithExistingValueArr(c, whDeltas); err != nil {
			return err
		}
	}
	if len(stocks) > 0 {
		if err := repository.StoreBulk(c, stocks); err != nil {
			return err
		}
	}

	return nil
}

func (repository *RepositoryStockImpl) StockUpdates(c context.Context, stockUpdates []*entity.StockUpdate) error {
	var whStocks []model.WarehouseStock
	var stocks []*model.Stock

	mergedWhStocks := make(map[string]*model.WarehouseStock)

	for index, stockUpdate := range stockUpdates {
		key := fmt.Sprintf("%v-%v", stockUpdate.ProID, stockUpdate.WhID)
		var qty float64

		if stockUpdate.QtyOut > 0 && stockUpdate.QtyIn > 0 {
			return errors.New(fmt.Sprintf("qtyin and qtyout can't greater than 0 on same row for product ID %v at index %v", key, index))
		}

		if stockUpdate.QtyOut > 0 && stockUpdate.QtyIn == 0 {
			qty = stockUpdate.QtyOut * -1
		} else {
			qty = float64(stockUpdate.QtyIn)
		}

		if existing, ok := mergedWhStocks[key]; ok {
			existing.Qty += qty
		} else {
			whStock := &model.WarehouseStock{
				CustID: stockUpdate.CustID,
				WhID:   stockUpdate.WhID,
				ProID:  stockUpdate.ProID,
				Qty:    qty,
			}
			mergedWhStocks[key] = whStock
		}

		stock := model.Stock{
			CustID:    stockUpdate.CustID,
			StockDate: stockUpdate.StockDate,
			TrCode:    stockUpdate.TrCode,
			TrNo:      stockUpdate.TrNo,
			WhID:      stockUpdate.WhID,
			ProID:     stockUpdate.ProID,
			ItemCdn:   1,
			UnitPrice: stockUpdate.UnitPrice,
			RefDetId:  stockUpdate.RefDetId,
			QtyIn:     stockUpdate.QtyIn,
			QtyOut:    stockUpdate.QtyOut,
		}

		stocks = append(stocks, &stock)

	}

	for _, whStock := range mergedWhStocks {
		whStocks = append(whStocks, *whStock)
	}
	log.Info("mergeWhStocks:", structs.StructToJson(whStocks))

	err := repository.UpsertWithExistingValueArr(c, whStocks)
	if err != nil {
		return err
	}

	err = repository.StoreBulk(c, stocks)
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryStockImpl) UpsertWithExistingValueArr(c context.Context, data []model.WarehouseStock) error {
	err := repository.model(c).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "cust_id"}, {Name: "wh_id"}, {Name: "pro_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"qty":          gorm.Expr("inv.warehouse_stock.qty + EXCLUDED.qty"),
				"qty_on_order": gorm.Expr("inv.warehouse_stock.qty_on_order + EXCLUDED.qty_on_order"),
			}),
		},
	).Create(&data).Error
	if err != nil {
		log.Error("UpsertQty, error:", err.Error())
		return err
	}
	return nil
}

func (repository *RepositoryStockImpl) StoreBulk(c context.Context, data []*model.Stock) error {
	err := repository.model(c).Create(data).Error
	return err
}

func (repository *RepositoryStockImpl) UpdateOnCustomerOrder(c context.Context, custId string, whId int64, proId int64, delta float64) error {
	return repository.model(c).
		Table("inv.warehouse_stock").
		Where("cust_id = ? AND wh_id = ? AND pro_id = ?", custId, whId, proId).
		UpdateColumn("qty_on_order", gorm.Expr("qty_on_order + ?", delta)).
		Error
}

func (repository *RepositoryStockImpl) GetCurrentStock(c context.Context, custId string, whId int64, proId int64) (float64, error) {
	var qty float64
	err := repository.model(c).
		Table("inv.warehouse_stock").
		Select("qty").
		Where("cust_id = ? AND wh_id = ? AND pro_id = ?", custId, whId, proId).
		Scan(&qty).Error
	return qty, err
}
