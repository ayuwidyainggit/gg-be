package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryGrBranchImpl struct {
		*gorm.DB
	}
)

type GrBranchRepository interface {
	Store(c context.Context, data *model.GrBranch) error
	FindByNo(grNo string, custId, parentCustId string) (gr model.GrBranchRead, err error)
	FindByGrBranchNo(grNo string, custId, parentCustId string) (gr model.GrBranchRead, err error)
	FindAllByCustId(dataFilter entity.GrBranchQueryFilter, custId, parentCustId string) ([]model.GrBranchList, int64, int, error)
	Update(c context.Context, grNo string, data *model.GrBranch) error
	Delete(c context.Context, custId string, grNo string, deletedBy int64) error
	FindGrBranchdetail(grBranchNo string, custId string) (grBranchDetails []model.GrBranchDetailList, err error)
	FindGrBranchdetailForUpdateWhStock(grNo string, custId string) (grDetails []model.GrBranchDetJoinGrBranchList, err error)
	UpdateGrBranchDetail(c context.Context, grDetails *model.GrBranchDet) error
	CreateGrBranchDetail(c context.Context, grDetails *model.GrBranchDetailCreate) (*model.GrBranchDetailCreate, error)
	DeleteGrBranchDetailNotInIDs(c context.Context, grNo string, IDs []int) error
	DeleteGrBranchDetailByGrBranchNo(c context.Context, grNo string, custId string) error
	DeleteStockNotInRefIds(c context.Context, custId, trNo string, newRefIds []int64) error
	DeleteStockInRefIds(c context.Context, custId, trNo string, oldRefIds []int64) error
	StoreWhStock(c context.Context, data *model.WhStock) error
	UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error
	StoreStock(c context.Context, data *model.Stock) error
	StoreProductCogs(c context.Context, data *model.ProductCogs) error
	FindQtyWhStock(custId string, proId, whId int64) (whStock model.WhStockList, err error)
	FindCogsProductDist(custId string, proId int64) (productDist model.ProductDist, err error)
	UpdateProductDist(c context.Context, custId string, proId int64, data model.ProductDist) error
	GetByInvoiceNo(invoiceNo string, custId, parentCustId string) (gr model.GrBranchList, err error)
	FindSupplierGrBranch(dataFilter entity.GrBranchSupplierQueryFilter, custId, parentCustId string) ([]model.GrBranchSupplier, int64, int, error)
	FindDistributorGrBranch(dataFilter entity.GrBranchDistributorQueryFilter, custId, parentCustId string) ([]model.GrBranchDistributor, int64, int, error)
	FindProductByListID(productIDs []int64) (products []model.Product, err error)
	FindWarehouseGrBranch(dataFilter entity.GrBranchWarehouseQueryFilter, custId, parentCustId string) ([]model.GrBranchWarehouse, int64, int, error)
	FindPrintWarehouseGrBranch(dataFilter entity.GrBranchPrintWarehouseQueryFilter, custId, parentCustId string) ([]model.GrBranchWarehouse, int64, int, error)
	FindGrBranchdetailWithDiscount(grBranchNo string, custId string) (grBranchDetails []model.GrBranchDetailList, err error)
	FindGrBranchOrderBookingDetails(orderBookingId int, custId string, parentCustId string) (orderBookingDetails []model.GrBranchOrderBookingDetail, err error)
	FindGrBranchOrderBookingList(dataFilter entity.GrBranchOrderBookingListQueryFilter, custId string, parentCustId string) ([]model.GrBranchOrderBooking, int64, int, error)
	FindGrBranchOrderBooking(OrderBookingId string, custId string, parentCustId string) (orderBooking model.GrBranchOrderBooking, err error)
	PrintGrBranch(c context.Context, custId string, grBranches []entity.GrBranchBulkPrintBody, printedBy int64) error
}

func NewGrBranchRepo(db *gorm.DB) *RepositoryGrBranchImpl {
	return &RepositoryGrBranchImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryGrBranchImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryGrBranchImpl) FindByNo(grNo string, custId, parentCustId string) (gr model.GrBranchRead, err error) {
	err = repository.
		Select("gr_branch.*, us.user_fullname AS updated_by_name, us2.user_fullname AS printed_by_name, cust.cust_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name, ob.type_approval").
		Joins("left join smc.m_customer cust on cust.cust_id = gr_branch.cust_id").
		Joins("left join sys.m_user us on us.user_id = gr_branch.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gr_branch.printed_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr_branch.wh_id AND (wh.cust_id = ? OR wh.cust_id = ?)", custId, parentCustId).
		Joins("LEFT JOIN inv.order_booking ob ON ob.po_no = gr_branch.po_no AND ob.cust_id = ?", custId).
		Where("gr_branch.gr_branch_no = ? AND gr_branch.cust_id=?", grNo, custId).
		Take(&gr).Error
	return gr, err
}

func (repository *RepositoryGrBranchImpl) FindByGrBranchNo(grNo string, custId string, parentCustId string) (gr model.GrBranchRead, err error) {
	err = repository.
		Select("gr_branch.*").
		Where("gr_branch.gr_branch_no = ?", grNo).
		Where("gr_branch.cust_id = ?", custId).
		Take(&gr).Error
	return gr, err
}

func (repository *RepositoryGrBranchImpl) GetByInvoiceNo(invoiceNo string, custId, parentCustId string) (gr model.GrBranchList, err error) {
	err = repository.
		Select("gr_branch.*, us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = gr_branch.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gr_branch.closed_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr_branch.wh_id AND wh.cust_id = ?", custId).
		Where("gr_branch.invoice_no = ? AND gr_branch.cust_id=?", invoiceNo, custId).
		Take(&gr).Error
	return gr, err
}

func (repository *RepositoryGrBranchImpl) Store(c context.Context, data *model.GrBranch) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGrBranchImpl) FindAllByCustId(dataFilter entity.GrBranchQueryFilter, custId, parentCustId string) ([]model.GrBranchList, int64, int, error) {
	var grs []model.GrBranchList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_branch_no").
		Joins("LEFT JOIN smc.m_customer cust ON cust.cust_id = gr_branch.cust_id").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = gr_branch.updated_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = gr_branch.printed_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr_branch.wh_id AND (wh.cust_id = gr_branch.cust_id OR wh.cust_id = ?)", parentCustId)

	query := repository.
		Select("gr_branch.*, us.user_fullname AS updated_by_name, cust.cust_name, us2.user_fullname AS printed_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("LEFT JOIN smc.m_customer cust ON cust.cust_id = gr_branch.cust_id").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = gr_branch.updated_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = gr_branch.printed_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr_branch.wh_id AND (wh.cust_id = gr_branch.cust_id OR wh.cust_id = ?)", parentCustId)

	if custId == parentCustId {
		queryCount.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
		query.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
	} else {
		queryCount.Where("gr_branch.cust_id=?", custId)
		query.Where("gr_branch.cust_id = ?", custId)
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.SupId) > 0 {
		query.Where("gr_branch.sup_id in ?", dataFilter.SupId)
		queryCount.Where("gr_branch.sup_id in ?", dataFilter.SupId)

	}

	if len(dataFilter.DistributorId) > 0 {
		query.Where("gr_branch.cust_id in ?", dataFilter.DistributorId)
		queryCount.Where("gr_branch.cust_id in ?", dataFilter.DistributorId)

	}

	if len(dataFilter.WhID) > 0 {
		query.Where("gr_branch.wh_id in ?", dataFilter.WhID)
		queryCount.Where("gr_branch.wh_id in ?", dataFilter.WhID)

	}

	if len(dataFilter.DataStatus) > 0 {
		query.Where("gr_branch.data_status in ?", dataFilter.DataStatus)
		queryCount.Where("gr_branch.data_status in ?", dataFilter.DataStatus)

	}

	if dataFilter.Query != "" {
		query.Where(`
			gr_branch.gr_branch_no ILIKE '%` + dataFilter.Query + `%' 	
		`)
		queryCount.Where(`
			gr_branch.gr_branch_no ILIKE '%` + dataFilter.Query + `%'
			
		`)
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
		query.Order("gr_branch.gr_branch_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&grs).Error
	if err != nil {
		return grs, total, 0, err
	}
	err = queryCount.Model(&grs).Count(&total).Error
	if err != nil {
		return grs, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	log.Println("repo list - total:", total)
	log.Println("repo list - lastPage:", lastPage)
	return grs, total, lastPage, nil

}

func (repository *RepositoryGrBranchImpl) Update(c context.Context, grBranchNo string, data *model.GrBranch) error {
	log.Println("data update:", data)
	result := repository.model(c).Model(data).Where("gr_branch_no=?", grBranchNo).Where("cust_id=?", data.CustID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryGrBranchImpl) Delete(c context.Context, custId string, grNo string, deletedBy int64) error {
	var data model.GrBranch
	result := repository.model(c).Model(&data).Where("gr_branch_no=? AND cust_id = ? AND is_del= ? ", grNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryGrBranchImpl) FindGrBranchdetail(grBranchNo string, custId string) (grBranchDetails []model.GrBranchDetailList, err error) {
	err = repository.
		Select(`gr_branch_det.*,
			pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3, pd.conv_unit2, pd.conv_unit3,
			mapd.disc_p, coalesce(grb.qty, gr_branch_det.qty) AS qty_remaining, COALESCE (whs.qty, 0) as wh_qty `).
		Joins("LEFT JOIN inv.gr_branch on gr_branch.gr_branch_no = gr_branch_det.gr_branch_no AND gr_branch.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_branch_det.pro_id").
		Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_branch_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN inv.good_receipt_balances grb ON grb.gr_branch_no = gr_branch_det.gr_branch_no AND grb.pro_id = gr_branch_det.pro_id", grBranchNo).
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_branch_det.pro_id AND whs.wh_id = gr_branch.wh_id").
		Where("gr_branch_det.gr_branch_no = ? AND gr_branch_det.cust_id = ?", grBranchNo, custId).Order("gr_branch_det.seq_no ASC").
		Find(&grBranchDetails).Error
	return grBranchDetails, err
}

func (repository *RepositoryGrBranchImpl) FindGrBranchdetailWithDiscount(grBranchNo string, custId string) (grBranchDetails []model.GrBranchDetailList, err error) {
	err = repository.
		Select(`inv.gr_branch_det.*,
			ob_det.qty1_alloc,
			ob_det.qty2_alloc,
			ob_det.qty3_alloc,
			pd.pro_code, pd.pro_name, COALESCE(inv.gr_branch_det.conv_unit2, pd.conv_unit2) as conv_unit2, COALESCE(inv.gr_branch_det.conv_unit3, pd.conv_unit3) as conv_unit3
			`).
		Joins("LEFT JOIN inv.gr_branch gr on gr.gr_branch_no = inv.gr_branch_det.gr_branch_no AND gr.cust_id = ?", custId).
		Joins("LEFT JOIN inv.order_booking ob on ob.po_no = gr.po_no AND ob.cust_id = ?", custId).
		Joins("LEFT JOIN inv.order_booking_detail ob_det on ob_det.order_booking_id = ob.order_booking_id AND ob_det.pro_id = inv.gr_branch_det.pro_id").
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = inv.gr_branch_det.pro_id").
		// Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_branch_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		// Joins("LEFT JOIN acf.account_payable_discounts apd ON apd.pro_id = gr_branch_det.pro_id AND apd.cust_id = ? AND apd.deleted_at IS NULL", custId).
		// Joins("LEFT JOIN inv.good_receipt_balances grb ON grb.gr_branch_no = gr_branch_det.gr_branch_no AND grb.pro_id = gr_branch_det.pro_id", grNo).
		// Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_branch_det.pro_id AND whs.wh_id = gr_branch.wh_id").
		Where("inv.gr_branch_det.gr_branch_no = ? AND inv.gr_branch_det.cust_id = ?", grBranchNo, custId).
		Order("inv.gr_branch_det.seq_no ASC").
		Find(&grBranchDetails).Error
	return grBranchDetails, err
}

func (repository *RepositoryGrBranchImpl) FindGrBranchdetailForUpdateWhStock(grNo string, custId string) (grDetails []model.GrBranchDetJoinGrBranchList, err error) {
	err = repository.
		Select(`gr_branch_det.gr_branch_det_id, gr_branch.wh_id, gr_branch_det.pro_id, gr_branch_det.qty`).
		Joins("LEFT JOIN inv.gr_branch gr ON gr_branch.gr_branch_no = gr_branch_det.gr_branch_no AND gr_branch.cust_id = ?", custId).
		Where("gr_branch_det.gr_branch_no = ? AND gr_branch_det.cust_id = ?", grNo, custId).Order("gr_branch_det.seq_no ASC").
		Find(&grDetails).Error
	return grDetails, err
}

func (repository *RepositoryGrBranchImpl) CreateGrBranchDetail(c context.Context, grDetails *model.GrBranchDetailCreate) (*model.GrBranchDetailCreate, error) {

	// var grDetail model.GrBranchDetailCreate
	// if err := structs.Automapper(*grDetails, &grDetail); err != nil {
	// 	return err
	// }

	query :=
		`INSERT INTO inv.gr_branch_det(
			cust_id, gr_branch_no, seq_no, pro_id, item_type, qty, unit_price1, unit_price2, unit_price3,
			vat, unit_id1, unit_id2, unit_id3, conv_unit2, conv_unit3, qty_ship1, qty_ship2, qty_ship3, qty_ship,
			qty_received1, qty_received2, qty_received3, qty_received, vat_value, amount)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23, $24, $25
		) RETURNING gr_branch_det_id;`

	// var grBranchDetail model.GrBranchDets
	result := repository.model(c).Exec(query,
		grDetails.CustID, grDetails.GrBranchNo, grDetails.SeqNo, grDetails.ProID, grDetails.ItemType, grDetails.Qty, grDetails.UnitPrice1, grDetails.UnitPrice2, grDetails.UnitPrice3,
		grDetails.Vat, grDetails.UnitId1, grDetails.UnitId2, grDetails.UnitId3, grDetails.ConvUnit2, grDetails.ConvUnit3, grDetails.QtyShip1, grDetails.QtyShip2, grDetails.QtyShip3, grDetails.QtyShip,
		grDetails.QtyReceived1, grDetails.QtyReceived2, grDetails.QtyReceived3, grDetails.QtyReceived, grDetails.VatValue, grDetails.Amount).Take(&grDetails) //.Scan(&grBranchDetail)

	if result.Error != nil {
		log.Println("CreateGrBranchDetail, result.Error:", structs.StructToJson(result.Error))
		return grDetails, result.Error
	}

	// log.Println("GrBranchDetId grBranchDetail :", grDetail.GrBranchDetId)
	// log.Println("GrBranchDetId grDetails :", grDetails.GrBranchDetId)
	// grDetails.GrBranchDetId = grBranchDetail.GrBranchDetId
	/*
		result := repository.model(c).Create(grDetails)
		if result.Error != nil {
			log.Println("CreateGrDetail, result.Error:", structs.StructToJson(result.Error))
			return grDetails, result.Error
		}
		if result.RowsAffected == 0 {
			return grDetails, errors.New("no rows affected")
		}
	*/
	// grDet.GrBranchDetId = grDetails.GrBranchDetId
	return grDetails, nil
}

func (repository *RepositoryGrBranchImpl) UpdateGrBranchDetail(c context.Context, grDetails *model.GrBranchDet) error {
	result := repository.model(c).Updates(&grDetails)
	if result.Error != nil {
		log.Println("UpdateGrBranchDetail, result:", structs.StructToJson(result.Error))
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryGrBranchImpl) DeleteGrBranchDetailNotInIDs(c context.Context, grNo string, IDs []int) error {
	var grDetails model.GrBranchDet
	err := repository.model(c).Where("gr_branch_no=? AND gr_branch_det_id not in (?) ", grNo, IDs).Delete(&grDetails).Error
	return err
}

func (repository *RepositoryGrBranchImpl) DeleteGrBranchDetailByGrBranchNo(c context.Context, grNo string, custId string) error {
	var grDetails model.GrBranchDet
	err := repository.model(c).Where("gr_branch_no = ?", grNo).Delete(&grDetails).Error
	return err
}

func (repository *RepositoryGrBranchImpl) FindQtyWhStock(custId string, proId, whId int64) (whs model.WhStockList, err error) {
	err = repository.
		Select("qty").
		Where("cust_id = ? AND pro_id = ? AND wh_id = ? ", custId, proId, whId).
		Take(&whs).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		qtyNotFound := float64(0)
		whs.Qty = &qtyNotFound
		err = nil
	}
	return whs, err
}

func (repository *RepositoryGrBranchImpl) FindCogsProductDist(custId string, proId int64) (productDist model.ProductDist, err error) {
	err = repository.
		Select("cogs").
		Where("cust_id = ? AND pro_id = ?", custId, proId).
		Take(&productDist).Error
	return productDist, err
}

func (repository *RepositoryGrBranchImpl) UpdateProductDist(c context.Context, custId string, proId int64, data model.ProductDist) error {
	result := repository.model(c).Model(&data).Where("cust_id = ? AND pro_id = ?", custId, proId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryGrBranchImpl) StoreWhStock(c context.Context, data *model.WhStock) error {
	err := repository.model(c).Exec(
		`INSERT INTO inv.wh_stock (
			cust_id, wh_id, pro_id, qty
		) VALUES (
			@cust_id, @wh_id, @pro_id, @qty 
		) ON CONFLICT ON CONSTRAINT wh_stock_pkey 
		DO UPDATE SET qty = inv.wh_stock.qty + EXCLUDED.qty;`,
		sql.Named("cust_id", data.CustID),
		sql.Named("wh_id", data.WhID),
		sql.Named("pro_id", data.ProID),
		sql.Named("qty", data.Qty)).Error
	if err != nil {
		log.Println("StoreWhStock, error:", err.Error())
		return err
	}
	return nil
}

func (repository *RepositoryGrBranchImpl) UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error {
	err := repository.model(c).Exec(
		`UPDATE inv.wh_stock 
		SET qty = qty-@qty 
		WHERE cust_id = @cust_id AND pro_id = @pro_id AND wh_id = @wh_id;`,
		sql.Named("cust_id", custId),
		sql.Named("wh_id", whId),
		sql.Named("pro_id", proId),
		sql.Named("qty", qty)).Error
	if err != nil {
		log.Println("UpdateOldStock, error:", err.Error())
		return err
	}
	return nil
}

func (repository *RepositoryGrBranchImpl) StoreStock(c context.Context, data *model.Stock) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGrBranchImpl) StoreProductCogs(c context.Context, data *model.ProductCogs) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGrBranchImpl) DeleteStockNotInRefIds(c context.Context, custId, trNo string, newRefIds []int64) error {
	var stock model.Stock
	err := repository.model(c).Where("cust_id = ? AND tr_no = ? AND ref_det_id NOT IN (?) ", custId, trNo, newRefIds).Delete(&stock).Error
	return err
}

func (repository *RepositoryGrBranchImpl) DeleteStockInRefIds(c context.Context, custId, trNo string, oldRefIds []int64) error {
	var stock model.Stock
	err := repository.model(c).Where("cust_id = ? AND tr_no = ? AND ref_det_id IN (?) ", custId, trNo, oldRefIds).Delete(&stock).Error
	return err
}

func (repository *RepositoryGrBranchImpl) FindSupplierGrBranch(dataFilter entity.GrBranchSupplierQueryFilter, custId, parentCustId string) ([]model.GrBranchSupplier, int64, int, error) {
	var grSuppliers []model.GrBranchSupplier
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_branch.sup_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId)
	queryCount.Where("gr_branch.cust_id=? ", custId)

	query := repository.
		Select("gr_branch.sup_id, sup.sup_code, sup_name").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId)
	query.Where("gr_branch.cust_id = ?", custId)

	if dataFilter.Query != "" {
		query.Where(`(
			sup.sup_code ILIKE '%` + dataFilter.Query + `%' OR
			sup.sup_name ILIKE '%` + dataFilter.Query + `%'
		)`)
		queryCount.Where(`(
			sup.sup_code ILIKE '%` + dataFilter.Query + `%' OR
			sup.sup_name ILIKE '%` + dataFilter.Query + `%'
		)`)
	}

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		query.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
		queryCount.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
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
		query.Order("gr_branch.sup_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Group("gr_branch.sup_id, sup.sup_code, sup.sup_name").Find(&grSuppliers).Error
	if err != nil {
		return grSuppliers, total, 0, err
	}
	err = queryCount.Model(&grSuppliers).Distinct("gr_branch.sup_id").Count(&total).Error
	if err != nil {
		return grSuppliers, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return grSuppliers, total, lastPage, nil
}

func (repository *RepositoryGrBranchImpl) FindDistributorGrBranch(dataFilter entity.GrBranchDistributorQueryFilter, custId, parentCustId string) ([]model.GrBranchDistributor, int64, int, error) {
	var grSuppliers []model.GrBranchDistributor
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_branch.cust_id").
		Joins("LEFT JOIN smc.m_customer cust ON cust.cust_id = gr_branch.cust_id AND cust.parent_cust_id = ?", parentCustId)
	// queryCount.Where("gr_branch.cust_id=? ", custId)

	query := repository.
		Select("gr_branch.cust_id, cust.cust_name").
		Joins("LEFT JOIN smc.m_customer cust ON cust.cust_id = gr_branch.cust_id AND cust.parent_cust_id = ?", parentCustId)
	// query.Where("gr_branch.cust_id = ?", custId)

	if custId == parentCustId {
		queryCount.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ? AND cust_id <> parent_cust_id)", parentCustId)
		query.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ? AND cust_id <> parent_cust_id)", parentCustId)
	} else {
		queryCount.Where("gr_branch.cust_id=?", custId)
		query.Where("gr_branch.cust_id = ?", custId)
	}

	if dataFilter.Query != "" {
		query.Where(`(
			cust.cust_id ILIKE '%` + dataFilter.Query + `%' OR
			cust.cust_name ILIKE '%` + dataFilter.Query + `%'
		)`)
		queryCount.Where(`(
			cust.cust_id ILIKE '%` + dataFilter.Query + `%' OR
			cust.cust_name ILIKE '%` + dataFilter.Query + `%'
		)`)
	}

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		query.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
		queryCount.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
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
		query.Order("gr_branch.cust_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Group("gr_branch.cust_id, cust.cust_name").Find(&grSuppliers).Error
	if err != nil {
		return grSuppliers, total, 0, err
	}
	err = queryCount.Model(&grSuppliers).Distinct("gr_branch.cust_id").Count(&total).Error
	if err != nil {
		return grSuppliers, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return grSuppliers, total, lastPage, nil
}

func (repository *RepositoryGrBranchImpl) FindProductByListID(productIDs []int64) (products []model.Product, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryGrBranchImpl) FindWarehouseGrBranch(dataFilter entity.GrBranchWarehouseQueryFilter, custId, parentCustId string) ([]model.GrBranchWarehouse, int64, int, error) {
	var grBranchWarehouses []model.GrBranchWarehouse
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	if dataFilter.TypeApproval == 2 {
		custId = parentCustId
	}

	queryCount := repository.Select("wh_id")
	queryCount.Where("cust_id=? ", custId).Where("wh_name NOT ILIKE '%Kanvas%' AND wh_name NOT ILIKE '%Canvas%'")

	query := repository.Select("wh_id, wh_code, wh_name").
		Where("wh_name NOT ILIKE '%Kanvas%' AND wh_name NOT ILIKE '%Canvas%'")
	query.Where("cust_id = ?", custId)

	// if len(dataFilter.SupID) > 0 {
	// 	queryCount.Where("inv.gr_branch.sup_id IN ?", dataFilter.SupID)
	// 	query.Where("inv.gr_branch.sup_id IN ?", dataFilter.SupID)
	// }

	// if dataFilter.Query != "" {
	// 	query.Where(`(
	// 		wh.wh_code ILIKE '%` + dataFilter.Query + `%' OR
	// 		wh.wh_name ILIKE '%` + dataFilter.Query + `%'
	// 	)`)
	// 	queryCount.Where(`(
	// 		wh.wh_code ILIKE '%` + dataFilter.Query + `%' OR
	// 		wh.wh_name ILIKE '%` + dataFilter.Query + `%'
	// 	)`)
	// }

	// if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
	// 	query.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
	// 	queryCount.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
	// }

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
		query.Order("wh_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	// offset := (page - 1) * dataFilter.Limit

	err := query.Find(&grBranchWarehouses).Error
	if err != nil {
		return grBranchWarehouses, total, 0, err
	}
	err = queryCount.Model(&grBranchWarehouses).Count(&total).Error
	if err != nil {
		return grBranchWarehouses, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return grBranchWarehouses, total, lastPage, nil
}

func (repository *RepositoryGrBranchImpl) FindGrBranchOrderBookingDetails(orderBookingId int, custId string, parentCustId string) (orderBookingDetails []model.GrBranchOrderBookingDetail, err error) {
	err = repository.
		Select(`order_booking_detail.*,
			COALESCE(order_booking_detail.qty1, 0) as qty1,
			COALESCE(order_booking_detail.qty2, 0) as qty2,
			COALESCE(order_booking_detail.qty3, 0) as qty3,
			COALESCE(order_booking_detail.qty1_alloc, 0) as qty1_alloc,
			COALESCE(order_booking_detail.qty2_alloc, 0) as qty2_alloc,
			COALESCE(order_booking_detail.qty3_alloc, 0) as qty3_alloc,
			pd.pro_code, pd.pro_name`).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = order_booking_detail.pro_id").
		Where("order_booking_detail.order_booking_id = ? AND order_booking_detail.cust_id IN (?, ?)", orderBookingId, custId, parentCustId).
		Order("order_booking_detail.order_booking_id ASC").
		Find(&orderBookingDetails).Error
	return orderBookingDetails, err
}

func (repository *RepositoryGrBranchImpl) FindGrBranchOrderBookingList(dataFilter entity.GrBranchOrderBookingListQueryFilter, custId string, parentCustId string) ([]model.GrBranchOrderBooking, int64, int, error) {
	var orderBookings []model.GrBranchOrderBooking
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.
		Select("order_booking.order_booking_id").
		Joins("INNER JOIN mst.m_supplier sup ON sup.sup_id = order_booking.sup_id AND sup.cust_id = ?", parentCustId).
		Where("order_booking.cust_id=?", custId).
		Where("order_booking.status_order_booking=2").
		Where(`order_booking.po_no NOT IN (
			SELECT gr_branch.po_no 
			FROM inv.gr_branch
			WHERE gr_branch.cust_id='` + custId + `' AND gr_branch.data_status <> ` + strconv.Itoa(entity.GR_BRANCH_REJECTED) + ` 
		)`).
		Where("order_booking.po_no IS NOT NULL")

	query := repository.
		Select("order_booking.order_booking_id, order_booking.po_no, order_booking.type_approval, order_booking.so_po as so_no, sup.sup_id, sup.sup_code, sup.sup_name").
		Joins("INNER JOIN mst.m_supplier sup ON sup.sup_id = order_booking.sup_id AND sup.cust_id = ?", parentCustId).
		// Where("order_booking.cust_id = ? AND order_booking.status_order_booking=2", custId).
		Where("order_booking.cust_id=?", custId).
		Where("order_booking.status_order_booking=2").
		Where(`order_booking.po_no NOT IN (
			SELECT gr_branch.po_no 
			FROM inv.gr_branch
			WHERE gr_branch.cust_id='` + custId + `' AND gr_branch.data_status <> ` + strconv.Itoa(entity.GR_BRANCH_REJECTED) + `
		)`).
		Where("order_booking.po_no IS NOT NULL")

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
		query.Order("order_booking.order_booking_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&orderBookings).Error
	if err != nil {
		return orderBookings, total, 0, err
	}
	err = queryCount.Model(&orderBookings).Count(&total).Error
	if err != nil {
		return orderBookings, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return orderBookings, total, lastPage, nil
}

func (repository *RepositoryGrBranchImpl) FindGrBranchOrderBooking(orderBookingId string, custId string, parentCustId string) (orderBooking model.GrBranchOrderBooking, err error) {
	err = repository.
		Select("order_booking.order_booking_id, order_booking.po_no, order_booking.type_approval, order_booking.so_po as so_no, order_booking.delivery_fee, sup.sup_id, sup.sup_code, sup.sup_name").
		Joins("INNER JOIN mst.m_supplier sup ON sup.sup_id = order_booking.sup_id AND sup.cust_id = ?", parentCustId).
		Where("order_booking.po_no = ? AND order_booking.cust_id = ? AND order_booking.status_order_booking=2", orderBookingId, custId).
		Take(&orderBooking).Error
	return orderBooking, err
}

func (repository *RepositoryGrBranchImpl) PrintGrBranch(c context.Context, custId string, grBranches []entity.GrBranchBulkPrintBody, printedBy int64) error {
	var data model.GrBranchPrint

	data.IsPrint = true
	data.PrintedBy = printedBy
	data.PrintedAt = time.Now()
	data.DataStatus = entity.GR_BRANCH_COMPLETED

	for _, grBranch := range grBranches {
		fmt.Println("GrBranchRepository Print")
		fmt.Println("GrBranchRepository GrBranchNo : ", grBranch.GrBranchNo)
		fmt.Println("GrBranchRepository CustId : ", grBranch.CustId)
		fmt.Println("GrBranchRepository PrintedBy : ", printedBy)
		result := repository.model(c).Model(&data).Where("gr_branch_no = ? AND cust_id = ? AND is_print= ?", grBranch.GrBranchNo, grBranch.CustId, false).
			Updates(data)
		if result.Error != nil {
			return result.Error
		}
		// if result.RowsAffected == 0 {
		// 	return errors.New("no rows affected")
		// }
	}
	return nil
}

func (repository *RepositoryGrBranchImpl) FindPrintWarehouseGrBranch(dataFilter entity.GrBranchPrintWarehouseQueryFilter, custId, parentCustId string) ([]model.GrBranchWarehouse, int64, int, error) {
	var grBranchWarehouses []model.GrBranchWarehouse
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	custId = parentCustId
	if dataFilter.TypeApproval == 2 {
		custId = dataFilter.CustId
	}

	queryCount := repository.Select("wh_id")
	queryCount.Where("cust_id=? ", custId).Where("wh_name NOT ILIKE '%Kanvas%' AND wh_name NOT ILIKE '%Canvas%'")

	query := repository.Select("wh_id, wh_code, wh_name").
		Where("wh_name NOT ILIKE '%Kanvas%' AND wh_name NOT ILIKE '%Canvas%'")
	query.Where("cust_id = ?", custId)

	// if len(dataFilter.SupID) > 0 {
	// 	queryCount.Where("inv.gr_branch.sup_id IN ?", dataFilter.SupID)
	// 	query.Where("inv.gr_branch.sup_id IN ?", dataFilter.SupID)
	// }

	// if dataFilter.Query != "" {
	// 	query.Where(`(
	// 		wh.wh_code ILIKE '%` + dataFilter.Query + `%' OR
	// 		wh.wh_name ILIKE '%` + dataFilter.Query + `%'
	// 	)`)
	// 	queryCount.Where(`(
	// 		wh.wh_code ILIKE '%` + dataFilter.Query + `%' OR
	// 		wh.wh_name ILIKE '%` + dataFilter.Query + `%'
	// 	)`)
	// }

	// if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
	// 	query.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
	// 	queryCount.Where("gr_branch.gr_branch_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
	// }

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
		query.Order("wh_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	// offset := (page - 1) * dataFilter.Limit

	err := query.Find(&grBranchWarehouses).Error
	if err != nil {
		return grBranchWarehouses, total, 0, err
	}
	err = queryCount.Model(&grBranchWarehouses).Count(&total).Error
	if err != nil {
		return grBranchWarehouses, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return grBranchWarehouses, total, lastPage, nil
}
