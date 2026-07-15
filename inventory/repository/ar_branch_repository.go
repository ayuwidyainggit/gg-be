package repository

import (
	"context"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"log"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryArBranchImpl struct {
		*gorm.DB
	}
)

type ArBranchRepository interface {
	// Store(c context.Context, data *model.ArBranch) error
	StoreArBranchPayment(c context.Context, data *model.ArBranchPaymentCreate) error
	FindByNo(grNo string, custId, parentCustId string) (gr model.ArBranchRead, err error)
	// FindByArBranchNo(grNo string, custId, parentCustId string) (gr model.ArBranchRead, err error)
	FindAllByCustId(dataFilter entity.ArBranchQueryFilter, custId, parentCustId string) ([]model.ArBranchList, int64, int, error)
	// Update(c context.Context, grNo string, data *model.ArBranch) error
	// Delete(c context.Context, custId string, grNo string, deletedBy int64) error
	FindArBranchdetail(grBranchNo string, custId string) (arBranchDetails []model.ArBranchDetailList, err error)
	// FindArBranchdetailForUpdateWhStock(grNo string, custId string) (grDetails []model.ArBranchDetJoinArBranchList, err error)
	// UpdateArBranchDetail(c context.Context, grDetails *model.ArBranchDet) error
	// CreateArBranchDetail(c context.Context, grDetails *model.ArBranchDetailCreate) (*model.ArBranchDetailCreate, error)
	// DeleteArBranchDetailNotInIDs(c context.Context, grNo string, IDs []int) error
	// DeleteArBranchDetailByArBranchNo(c context.Context, grNo string) error
	// DeleteStockNotInRefIds(c context.Context, custId, trNo string, newRefIds []int64) error
	// DeleteStockInRefIds(c context.Context, custId, trNo string, oldRefIds []int64) error
	// StoreWhStock(c context.Context, data *model.WhStock) error
	// UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error
	// StoreStock(c context.Context, data *model.Stock) error
	// StoreProductCogs(c context.Context, data *model.ProductCogs) error
	// FindQtyWhStock(custId string, proId, whId int64) (whStock model.WhStockList, err error)
	// FindCogsProductDist(custId string, proId int64) (productDist model.ProductDist, err error)
	// UpdateProductDist(c context.Context, custId string, proId int64, data model.ProductDist) error
	// GetByInvoiceNo(invoiceNo string, custId, parentCustId string) (gr model.ArBranchList, err error)
	// FindSupplierArBranch(dataFilter entity.ArBranchSupplierQueryFilter, custId, parentCustId string) ([]model.ArBranchSupplier, int64, int, error)
	// FindProductByListID(productIDs []int64) (products []model.Product, err error)
	// FindWarehouseArBranch(dataFilter entity.ArBranchWarehouseQueryFilter, custId, parentCustId string) ([]model.ArBranchWarehouse, int64, int, error)
	FindArBranchDetailWithDiscount(grBranchNo string, custId string) (arBranchDetails []model.ArBranchDetailList, err error)
	FindArBranchPayments(invoiceNoBranch string, custId string) (arBranchPayments []model.ArBranchPaymentList, err error)
	// FindArBranchOrderBookingDetails(orderBookingId int, custId string, parentCustId string) (orderBookingDetails []model.ArBranchOrderBookingDetail, err error)
	// FindArBranchOrderBookingList(dataFilter entity.ArBranchOrderBookingListQueryFilter, custId string, parentCustId string) ([]model.ArBranchOrderBooking, int64, int, error)
	// FindArBranchOrderBooking(OrderBookingId string, custId string, parentCustId string) (orderBooking model.ArBranchOrderBooking, err error)
	// PrintArBranch(c context.Context, custId string, arBranches []entity.ArBranchBulkPrintBody, printedBy int64) error
	FindDistributorsArBranch(dataFilter entity.ArBranchDistributorsFilterQueryFilter, custId, parentCustId string) ([]model.ArBranchDistributor, int64, int, error)
	FindSuppliersArBranch(dataFilter entity.ArBranchSuppliersFilterQueryFilter, custId, parentCustId string) ([]model.ArBranchSupplier, int64, int, error)
}

func NewArBranchRepo(db *gorm.DB) *RepositoryArBranchImpl {
	return &RepositoryArBranchImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryArBranchImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryArBranchImpl) FindByNo(grNo string, custId, parentCustId string) (gr model.ArBranchRead, err error) {
	err = repository.
		Select(`
			gr_branch.*, 
			us.user_fullname AS updated_by_name, 
			us2.user_fullname AS printed_by_name, 
			sup.sup_code, 
			sup_name, 
			wh.wh_code, 
			wh.wh_name, 
			ob.type_approval,
			gr_branch.total as invoice_amount,
			cust.cust_name, 
			parent_cust.cust_id as parent_cust_id, 
			parent_cust.cust_name as parent_cust_name,
			coalesce(paid_invoices.paid_amount, 0) as paid_amount,
			CASE
				WHEN (gr_branch.total - coalesce(paid_invoices.paid_amount, 0)) < 0 THEN 0
				ELSE (gr_branch.total - coalesce(paid_invoices.paid_amount, 0))
			END 
			as remaining_amount,
			coalesce(paid_invoices.deposit_amount, 0) as deposit_amount
		`).
		Joins("left join sys.m_user us on us.user_id = gr_branch.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gr_branch.printed_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr_branch.wh_id AND (wh.cust_id = ? OR wh.cust_id = ?)", custId, parentCustId).
		Joins("LEFT JOIN smc.m_customer cust on cust.cust_id = gr_branch.cust_id AND cust.parent_cust_id = ?", parentCustId).
		Joins("LEFT JOIN smc.m_customer parent_cust on parent_cust.cust_id = cust.parent_cust_id").
		Joins("LEFT JOIN inv.order_booking ob ON ob.po_no = gr_branch.po_no AND ob.cust_id = gr_branch.cust_id", custId).
		Joins(`
		left join (
			select inv.gr_branch_payment.invoice_no_branch, 
				inv.gr_branch_payment.cust_id,
				sum(
					case 
						when inv.gr_branch_payment.verification_status = `+strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED)+` then coalesce(inv.gr_branch_payment.payment_amount, 0) + coalesce(inv.gr_branch_payment.discount, 0) + coalesce(inv.gr_branch_payment.payment_balance, 0)
						else 0
					end
				) as paid_amount,
				(coalesce(sum(inv.gr_branch_payment.payment_amount), 0) + coalesce(sum(inv.gr_branch_payment.discount), 0) + coalesce(sum(inv.gr_branch_payment.payment_balance), 0)) as deposit_amount
			from inv.gr_branch_payment
			left join inv.gr_branch on inv.gr_branch_payment.invoice_no_branch = inv.gr_branch.invoice_no_branch AND gr_branch.cust_id = inv.gr_branch_payment.cust_id
			where inv.gr_branch_payment.verification_status in (`+strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED)+`, `+strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_NEED_REVIEW)+`) 
			group by inv.gr_branch_payment.invoice_no_branch, inv.gr_branch_payment.cust_id
		) paid_invoices on paid_invoices.invoice_no_branch = inv.gr_branch.invoice_no_branch AND gr_branch.cust_id = paid_invoices.cust_id
	 	`).
		Where("gr_branch.gr_branch_no = ? AND gr_branch.cust_id=?", grNo, custId).
		Take(&gr).Error
	return gr, err
}

/*
	func (repository *RepositoryArBranchImpl) FindByArBranchNo(grNo string, custId string, parentCustId string) (gr model.ArBranchRead, err error) {
		err = repository.
			Select("gr_branch.*").
			Where("gr_branch.gr_branch_no = ?", grNo).
			Where("gr_branch.cust_id = ?", custId).
			Take(&gr).Error
		return gr, err
	}

	func (repository *RepositoryArBranchImpl) GetByInvoiceNo(invoiceNo string, custId, parentCustId string) (gr model.ArBranchList, err error) {
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
*/
func (repository *RepositoryArBranchImpl) StoreArBranchPayment(c context.Context, data *model.ArBranchPaymentCreate) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryArBranchImpl) FindAllByCustId(dataFilter entity.ArBranchQueryFilter, custId, parentCustId string) ([]model.ArBranchList, int64, int, error) {
	var grs []model.ArBranchList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryPaidInvoices := `left join (
			select inv.gr_branch_payment.invoice_no_branch, 
			inv.gr_branch_payment.cust_id,
			(coalesce(sum(inv.gr_branch_payment.payment_amount), 0) + coalesce(sum(inv.gr_branch_payment.discount), 0) + coalesce(sum(inv.gr_branch_payment.payment_balance), 0)) as paid_amount
		from inv.gr_branch_payment
		left join inv.gr_branch on inv.gr_branch_payment.invoice_no_branch = inv.gr_branch.invoice_no_branch AND gr_branch.cust_id = inv.gr_branch_payment.cust_id
		where inv.gr_branch_payment.verification_status = ` + strconv.Itoa(entity.AR_BRANCH_VERIFICATION_STATUS_APPROVED) + `
		group by inv.gr_branch_payment.invoice_no_branch, inv.gr_branch_payment.cust_id
	) paid_invoices on paid_invoices.invoice_no_branch = inv.gr_branch.invoice_no_branch AND gr_branch.cust_id = paid_invoices.cust_id`

	queryCount := repository.Select("gr_branch_no").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = gr_branch.updated_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = gr_branch.printed_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr_branch.wh_id AND (wh.cust_id = gr_branch.cust_id OR wh.cust_id = ?)", parentCustId).
		Joins("LEFT JOIN smc.m_customer cust on cust.cust_id = gr_branch.cust_id AND cust.parent_cust_id = ?", parentCustId).
		Joins("LEFT JOIN smc.m_customer parent_cust on parent_cust.cust_id = cust.parent_cust_id").
		Joins(queryPaidInvoices).
		Where("gr_branch.data_status = ? AND gr_branch.is_print IS TRUE", entity.GR_BRANCH_COMPLETED)

	query := repository.
		Select(`
			gr_branch.*, 
			us.user_fullname AS updated_by_name, 
			us2.user_fullname AS printed_by_name, 
			sup.sup_code, 
			sup_name, 
			wh.wh_code, 
			wh.wh_name, 
			cust.cust_name,
			parent_cust.cust_id as parent_cust_id, 
			parent_cust.cust_name as parent_cust_name,                                                                                                                                                                                
			gr_branch.total as invoice_amount,
			coalesce(paid_invoices.paid_amount, 0) as paid_amount,
			CASE
				WHEN (gr_branch.total - coalesce(paid_invoices.paid_amount, 0)) < 0 THEN 0
				ELSE (gr_branch.total - coalesce(paid_invoices.paid_amount, 0))
			END 
			as remaining_amount
		`).
		Joins("LEFT JOIN sys.m_user us ON us.user_id = gr_branch.updated_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = gr_branch.printed_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN smc.m_customer cust on cust.cust_id = gr_branch.cust_id AND cust.parent_cust_id = ?", parentCustId).
		Joins("LEFT JOIN smc.m_customer parent_cust on parent_cust.cust_id = cust.parent_cust_id").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr_branch.wh_id AND (wh.cust_id = gr_branch.cust_id OR wh.cust_id = ?)", parentCustId).
		Joins(queryPaidInvoices).
		Where("gr_branch.data_status = ? AND gr_branch.is_print IS TRUE", entity.GR_BRANCH_COMPLETED)

	if custId == parentCustId {
		queryCount.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
		query.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
	} else {
		queryCount.Where("gr_branch.cust_id=?", custId)
		query.Where("gr_branch.cust_id = ?", custId)
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("gr_branch.invoice_date_branch between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("gr_branch.invoice_date_branch between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.CustId) > 0 {
		query.Where("gr_branch.cust_id in ?", dataFilter.CustId)
		queryCount.Where("gr_branch.cust_id in ?", dataFilter.CustId)

	}

	if len(dataFilter.SupID) > 0 {
		query.Where("gr_branch.sup_id in ?", dataFilter.SupID)
		queryCount.Where("gr_branch.sup_id in ?", dataFilter.SupID)
	}

	// if len(dataFilter.DataStatus) > 0 {
	// 	query.Where("gr_branch.data_status in ?", dataFilter.DataStatus)
	// 	queryCount.Where("gr_branch.data_status in ?", dataFilter.DataStatus)

	// }

	if dataFilter.Query != "" {
		query.Where(`
			gr_branch.invoice_no_branch ILIKE '%` + dataFilter.Query + `%' 	
		`)
		queryCount.Where(`
			gr_branch.invoice_no_branch ILIKE '%` + dataFilter.Query + `%'
			
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
		query.Order("gr_branch.invoice_no_branch DESC")
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

/*
	func (repository *RepositoryArBranchImpl) Update(c context.Context, arBranchNo string, data *model.ArBranch) error {
		log.Println("data update:", data)
		result := repository.model(c).Model(data).Where("gr_branch_no=?", arBranchNo).Updates(data)
		if result.Error != nil {
			return result.Error
		}
		// if result.RowsAffected == 0 {
		// 	return errors.New("no rows affected")
		// }
		return nil
	}

	func (repository *RepositoryArBranchImpl) Delete(c context.Context, custId string, grNo string, deletedBy int64) error {
		var data model.ArBranch
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
*/
func (repository *RepositoryArBranchImpl) FindArBranchdetail(arBranchNo string, custId string) (arBranchDetails []model.ArBranchDetailList, err error) {
	err = repository.
		Select(`gr_branch_det.*,
			pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3, pd.conv_unit2, pd.conv_unit3,
			mapd.disc_p, coalesce(grb.qty, gr_branch_det.qty) AS qty_remaining, COALESCE (whs.qty, 0) as wh_qty `).
		Joins("LEFT JOIN inv.gr_branch on gr_branch.gr_branch_no = gr_branch_det.gr_branch_no AND gr_branch.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_branch_det.pro_id").
		Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_branch_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN inv.good_receipt_balances grb ON grb.gr_branch_no = gr_branch_det.gr_branch_no AND grb.pro_id = gr_branch_det.pro_id", arBranchNo).
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_branch_det.pro_id AND whs.wh_id = gr_branch.wh_id").
		Where("gr_branch_det.gr_branch_no = ? AND gr_branch_det.cust_id = ?", arBranchNo, custId).Order("gr_branch_det.seq_no ASC").
		Find(&arBranchDetails).Error
	return arBranchDetails, err
}

func (repository *RepositoryArBranchImpl) FindArBranchDetailWithDiscount(grBranchNo string, custId string) (arBranchDetails []model.ArBranchDetailList, err error) {
	err = repository.
		Select(`
			gr_branch_det.*,
			ob_det.qty1_alloc,
			ob_det.qty2_alloc,
			ob_det.qty3_alloc,
			pd.pro_code, 
			pd.pro_name, 
			COALESCE(gr_branch_det.conv_unit2, pd.conv_unit2) as conv_unit2, 
			COALESCE(gr_branch_det.conv_unit3, pd.conv_unit3) as conv_unit3
		`).
		Joins("LEFT JOIN inv.gr_branch on gr_branch.gr_branch_no = gr_branch_det.gr_branch_no AND gr_branch.cust_id = ?", custId).
		Joins("LEFT JOIN inv.order_booking ob on ob.po_no = gr_branch.po_no AND ob.cust_id = ?", custId).
		Joins("LEFT JOIN inv.order_booking_detail ob_det on ob_det.order_booking_id = ob.order_booking_id AND ob_det.pro_id = inv.gr_branch_det.pro_id").
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_branch_det.pro_id").
		// Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_branch_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		// Joins("LEFT JOIN acf.account_payable_discounts apd ON apd.pro_id = gr_branch_det.pro_id AND apd.cust_id = ? AND apd.deleted_at IS NULL", custId).
		// Joins("LEFT JOIN inv.good_receipt_balances grb ON grb.gr_branch_no = gr_branch_det.gr_branch_no AND grb.pro_id = gr_branch_det.pro_id", grNo).
		// Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_branch_det.pro_id AND whs.wh_id = gr_branch.wh_id").
		Where("gr_branch_det.gr_branch_no = ? AND gr_branch_det.cust_id = ?", grBranchNo, custId).
		Order("gr_branch_det.seq_no ASC").
		Find(&arBranchDetails).Error
	return arBranchDetails, err
}

func (repository *RepositoryArBranchImpl) FindArBranchPayments(invoiceNoBranch string, custId string) (arBranchPayments []model.ArBranchPaymentList, err error) {
	err = repository.
		Select(`
			gr_branch_payment.*,
			us.user_fullname as verified_by_name
		`).
		Joins("LEFT JOIN sys.m_user us ON us.user_id = gr_branch_payment.verified_by").
		Joins("LEFT JOIN inv.gr_branch on gr_branch.invoice_no_branch = gr_branch_payment.invoice_no_branch AND inv.gr_branch.cust_id = ?", custId).
		Where("gr_branch_payment.cust_id = ? AND gr_branch_payment.invoice_no_branch = ?", custId, invoiceNoBranch).
		Order("gr_branch_payment.gr_branch_payment_id ASC").
		Find(&arBranchPayments).Error
	return arBranchPayments, err
}

/*
func (repository *RepositoryArBranchImpl) FindArBranchdetailForUpdateWhStock(grNo string, custId string) (grDetails []model.ArBranchDetJoinArBranchList, err error) {
	err = repository.
		Select(`gr_branch_det.gr_branch_det_id, gr_branch.wh_id, gr_branch_det.pro_id, gr_branch_det.qty`).
		Joins("LEFT JOIN inv.gr_branch gr ON gr_branch.gr_branch_no = gr_branch_det.gr_branch_no AND gr_branch.cust_id = ?", custId).
		Where("gr_branch_det.gr_branch_no = ? AND gr_branch_det.cust_id = ?", grNo, custId).Order("gr_branch_det.seq_no ASC").
		Find(&grDetails).Error
	return grDetails, err
}

func (repository *RepositoryArBranchImpl) CreateArBranchDetail(c context.Context, grDetails *model.ArBranchDetailCreate) (*model.ArBranchDetailCreate, error) {

	// var grDetail model.ArBranchDetailCreate
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

	// var arBranchDetail model.ArBranchDets
	result := repository.model(c).Exec(query,
		grDetails.CustID, grDetails.ArBranchNo, grDetails.SeqNo, grDetails.ProID, grDetails.ItemType, grDetails.Qty, grDetails.UnitPrice1, grDetails.UnitPrice2, grDetails.UnitPrice3,
		grDetails.Vat, grDetails.UnitId1, grDetails.UnitId2, grDetails.UnitId3, grDetails.ConvUnit2, grDetails.ConvUnit3, grDetails.QtyShip1, grDetails.QtyShip2, grDetails.QtyShip3, grDetails.QtyShip,
		grDetails.QtyReceived1, grDetails.QtyReceived2, grDetails.QtyReceived3, grDetails.QtyReceived, grDetails.VatValue, grDetails.Amount).Take(&grDetails) //.Scan(&arBranchDetail)

	if result.Error != nil {
		log.Println("CreateArBranchDetail, result.Error:", structs.StructToJson(result.Error))
		return grDetails, result.Error
	}

	// log.Println("ArBranchDetId arBranchDetail :", grDetail.ArBranchDetId)
	// log.Println("ArBranchDetId grDetails :", grDetails.ArBranchDetId)
	// grDetails.ArBranchDetId = arBranchDetail.ArBranchDetId

		// result := repository.model(c).Create(grDetails)
		// if result.Error != nil {
		// 	log.Println("CreateGrDetail, result.Error:", structs.StructToJson(result.Error))
		// 	return grDetails, result.Error
		// }
		// if result.RowsAffected == 0 {
		// 	return grDetails, errors.New("no rows affected")
		// }

	// grDet.ArBranchDetId = grDetails.ArBranchDetId
	return grDetails, nil
}

func (repository *RepositoryArBranchImpl) UpdateArBranchDetail(c context.Context, grDetails *model.ArBranchDet) error {
	result := repository.model(c).Updates(&grDetails)
	if result.Error != nil {
		log.Println("UpdateArBranchDetail, result:", structs.StructToJson(result.Error))
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryArBranchImpl) DeleteArBranchDetailNotInIDs(c context.Context, grNo string, IDs []int) error {
	var grDetails model.ArBranchDet
	err := repository.model(c).Where("gr_branch_no=? AND gr_branch_det_id not in (?) ", grNo, IDs).Delete(&grDetails).Error
	return err
}

func (repository *RepositoryArBranchImpl) DeleteArBranchDetailByArBranchNo(c context.Context, grNo string) error {
	var grDetails model.ArBranchDet
	err := repository.model(c).Where("gr_branch_no = ?", grNo).Delete(&grDetails).Error
	return err
}

func (repository *RepositoryArBranchImpl) FindQtyWhStock(custId string, proId, whId int64) (whs model.WhStockList, err error) {
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

func (repository *RepositoryArBranchImpl) FindCogsProductDist(custId string, proId int64) (productDist model.ProductDist, err error) {
	err = repository.
		Select("cogs").
		Where("cust_id = ? AND pro_id = ?", custId, proId).
		Take(&productDist).Error
	return productDist, err
}

func (repository *RepositoryArBranchImpl) UpdateProductDist(c context.Context, custId string, proId int64, data model.ProductDist) error {
	result := repository.model(c).Model(&data).Where("cust_id = ? AND pro_id = ?", custId, proId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryArBranchImpl) StoreWhStock(c context.Context, data *model.WhStock) error {
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

func (repository *RepositoryArBranchImpl) UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error {
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

func (repository *RepositoryArBranchImpl) StoreStock(c context.Context, data *model.Stock) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryArBranchImpl) StoreProductCogs(c context.Context, data *model.ProductCogs) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryArBranchImpl) DeleteStockNotInRefIds(c context.Context, custId, trNo string, newRefIds []int64) error {
	var stock model.Stock
	err := repository.model(c).Where("cust_id = ? AND tr_no = ? AND ref_det_id NOT IN (?) ", custId, trNo, newRefIds).Delete(&stock).Error
	return err
}

func (repository *RepositoryArBranchImpl) DeleteStockInRefIds(c context.Context, custId, trNo string, oldRefIds []int64) error {
	var stock model.Stock
	err := repository.model(c).Where("cust_id = ? AND tr_no = ? AND ref_det_id IN (?) ", custId, trNo, oldRefIds).Delete(&stock).Error
	return err
}

func (repository *RepositoryArBranchImpl) FindSupplierArBranch(dataFilter entity.ArBranchSupplierQueryFilter, custId, parentCustId string) ([]model.ArBranchSupplier, int64, int, error) {
	var grSuppliers []model.ArBranchSupplier
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

func (repository *RepositoryArBranchImpl) FindProductByListID(productIDs []int64) (products []model.Product, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryArBranchImpl) FindWarehouseArBranch(dataFilter entity.ArBranchWarehouseQueryFilter, custId, parentCustId string) ([]model.ArBranchWarehouse, int64, int, error) {
	var arBranchWarehouses []model.ArBranchWarehouse
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

	err := query.Find(&arBranchWarehouses).Error
	if err != nil {
		return arBranchWarehouses, total, 0, err
	}
	err = queryCount.Model(&arBranchWarehouses).Count(&total).Error
	if err != nil {
		return arBranchWarehouses, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return arBranchWarehouses, total, lastPage, nil
}

func (repository *RepositoryArBranchImpl) FindArBranchOrderBookingDetails(orderBookingId int, custId string, parentCustId string) (orderBookingDetails []model.ArBranchOrderBookingDetail, err error) {
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
		Where("order_booking_detail.order_booking_id = ? AND order_booking_detail.cust_id = ?", orderBookingId, custId).
		Order("order_booking_detail.order_booking_id ASC").
		Find(&orderBookingDetails).Error
	return orderBookingDetails, err
}

func (repository *RepositoryArBranchImpl) FindArBranchOrderBookingList(dataFilter entity.ArBranchOrderBookingListQueryFilter, custId string, parentCustId string) ([]model.ArBranchOrderBooking, int64, int, error) {
	var orderBookings []model.ArBranchOrderBooking
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
			WHERE gr_branch.cust_id='` + custId + `'
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
			WHERE gr_branch.cust_id='` + custId + `'
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

func (repository *RepositoryArBranchImpl) FindArBranchOrderBooking(orderBookingId string, custId string, parentCustId string) (orderBooking model.ArBranchOrderBooking, err error) {
	err = repository.
		Select("order_booking.order_booking_id, order_booking.po_no, order_booking.type_approval, order_booking.so_po as so_no, order_booking.delivery_fee, sup.sup_id, sup.sup_code, sup.sup_name").
		Joins("INNER JOIN mst.m_supplier sup ON sup.sup_id = order_booking.sup_id AND sup.cust_id = ?", parentCustId).
		Where("order_booking.po_no = ? AND order_booking.cust_id = ? AND order_booking.status_order_booking=2", orderBookingId, custId).
		Take(&orderBooking).Error
	return orderBooking, err
}

func (repository *RepositoryArBranchImpl) PrintArBranch(c context.Context, custId string, arBranches []entity.ArBranchBulkPrintBody, printedBy int64) error {
	var data model.ArBranchPrint

	data.IsPrint = true
	data.PrintedBy = printedBy
	data.PrintedAt = time.Now()

	for _, arBranch := range arBranches {
		fmt.Println("ArBranchRepository Print")
		fmt.Println("ArBranchRepository ArBranchNo : ", arBranch.ArBranchNo)
		fmt.Println("ArBranchRepository CustId : ", arBranch.CustId)
		fmt.Println("ArBranchRepository PrintedBy : ", printedBy)
		result := repository.model(c).Model(&data).Where("gr_branch_no = ? AND cust_id = ? AND is_print= ?", arBranch.ArBranchNo, arBranch.CustId, false).
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
*/

func (repository *RepositoryArBranchImpl) FindDistributorsArBranch(dataFilter entity.ArBranchDistributorsFilterQueryFilter, custId, parentCustId string) ([]model.ArBranchDistributor, int64, int, error) {
	var distributors []model.ArBranchDistributor
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_branch.cust_id").
		Joins("LEFT JOIN smc.m_customer dist ON dist.cust_id = gr_branch.cust_id AND dist.parent_cust_id = ?", parentCustId).
		Where("gr_branch.data_status = ? AND gr_branch.is_print IS TRUE", entity.GR_BRANCH_COMPLETED)

	query := repository.Select("gr_branch.cust_id, dist.cust_name").
		Joins("LEFT JOIN smc.m_customer dist ON dist.cust_id = gr_branch.cust_id AND dist.parent_cust_id = ?", parentCustId).
		Where("gr_branch.data_status = ? AND gr_branch.is_print IS TRUE", entity.GR_BRANCH_COMPLETED)

	if custId == parentCustId {
		queryCount.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
		query.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
	} else {
		queryCount.Where("gr_branch.cust_id=?", custId)
		query.Where("gr_branch.cust_id = ?", custId)
	}

	if dataFilter.Query != "" {
		query.Where(`(
			gr_branch.cust_id ILIKE '%` + dataFilter.Query + `%' OR
			dist.cust_name ILIKE '%` + dataFilter.Query + `%'
		)`)
		queryCount.Where(`(
			gr_branch.cust_id ILIKE '%` + dataFilter.Query + `%' OR
			dist.cust_name ILIKE '%` + dataFilter.Query + `%'
		)`)
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

	err := query.Limit(limit).Offset(offset).Group("gr_branch.cust_id, dist.cust_name").Find(&distributors).Error
	if err != nil {
		return distributors, total, 0, err
	}
	err = queryCount.Model(&distributors).Distinct("gr_branch.cust_id").Count(&total).Error
	if err != nil {
		return distributors, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return distributors, total, lastPage, nil
}

func (repository *RepositoryArBranchImpl) FindSuppliersArBranch(dataFilter entity.ArBranchSuppliersFilterQueryFilter, custId, parentCustId string) ([]model.ArBranchSupplier, int64, int, error) {
	var suppliers []model.ArBranchSupplier
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_branch.sup_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Where("gr_branch.data_status = ? AND gr_branch.is_print IS TRUE", entity.GR_BRANCH_COMPLETED)

	query := repository.Select("gr_branch.sup_id, sup.sup_code, sup.sup_name").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Where("gr_branch.data_status = ? AND gr_branch.is_print IS TRUE", entity.GR_BRANCH_COMPLETED)

	if custId == parentCustId {
		queryCount.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
		query.Where("gr_branch.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = ?)", parentCustId)
	} else {
		queryCount.Where("gr_branch.cust_id=?", custId)
		query.Where("gr_branch.cust_id = ?", custId)
	}

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

	err := query.Limit(limit).Offset(offset).Group("gr_branch.sup_id, sup.sup_code, sup.sup_name").Find(&suppliers).Error
	if err != nil {
		return suppliers, total, 0, err
	}
	err = queryCount.Model(&suppliers).Distinct("gr_branch.sup_id").Count(&total).Error
	if err != nil {
		return suppliers, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return suppliers, total, lastPage, nil
}
