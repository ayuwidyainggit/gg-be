package repository

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"log"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositorySupplierReturnImpl struct {
		*gorm.DB
	}
)

type SupplierReturnRepository interface {
	Store(c context.Context, data *model.SupplierReturn) error
	StoreDetail(c context.Context, detail *model.SupplierReturnDet) (*model.SupplierReturnDet, error)
	GetGoodReceiptByProduct(c context.Context, GrNo string, proID int64) (grDet model.GrDetList, err error)
	FindAllByCustId(dataFilter entity.SupplierReturnQueryFilter) ([]model.SupplierReturnGet, int64, int, error)
	FindByNo(supplierReturnNo, custId, parentCustId string) (supplierReturn model.SupplierReturnGet, err error)
	FindDetBySupplierReturnNo(supplierReturnNo string, custId string) (Details []model.SupplierReturnDetGet, err error)
	FindSupplierReturn(dataFilter entity.ReturnSupplierQueryFilter, custId, parentCustId string) ([]model.ReturnSuppliers, int64, int, error)
	FindProductByListID(productIDs []int64) (products []model.Product, err error)
	GetRemainingQtyInvoice(invoiceNo string, custId string, proID int64) (remainingQty model.RemainingQty, err error)
	SaveRemainingQty(c context.Context, balances model.ProductInvoiceBalances)
	GetRemainingProductQtyByInvoiceNo(invoiceNo string) (remainingQtyProducts []model.RemainingQtyProduct, err error)
	SetInvoiceNoIsCanReturn(c context.Context, invoiceNo string, isCanReturn bool) (err error)
	UpdateStatus(c context.Context, supplierReturnNo string, custId string, status int) error
	GetApProductList(InvNo string, custId string, proID int64) (apDet model.AccountPayableProductList, err error)
	GetInvoiceFromAPReturn(custId string, invoiceNo string) (ap model.AccountPayable, err error)
}

func NewSupplierReturnRepo(db *gorm.DB) *RepositorySupplierReturnImpl {
	return &RepositorySupplierReturnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositorySupplierReturnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositorySupplierReturnImpl) Store(c context.Context, data *model.SupplierReturn) error {
	err := repository.model(c).Create(&data)
	if err != nil {
		return err.Error
	}
	return nil
}

func (repository *RepositorySupplierReturnImpl) StoreDetail(c context.Context, detail *model.SupplierReturnDet) (*model.SupplierReturnDet, error) {
	result := repository.model(c).Create(detail)
	if result.Error != nil {
		return detail, result.Error
	}
	if result.RowsAffected == 0 {
		return detail, errors.New("no rows affected")
	}

	return detail, nil
}

func (repository *RepositorySupplierReturnImpl) GetGoodReceiptByProduct(c context.Context, GrNo string, proID int64) (grDet model.GrDetList, err error) {
	err = repository.
		Select(`gr_det.*,
			pd.pro_code, pd.pro_name, pd.purch_price1, pd.purch_price2, pd.purch_price3, pd.conv_unit2, pd.conv_unit3`).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_det.pro_id").
		Where("gr_det.gr_no = ? AND gr_det.pro_id = ?", GrNo, proID).
		Take(&grDet).Error
	return grDet, err
}

func (repository *RepositorySupplierReturnImpl) FindAllByCustId(dataFilter entity.SupplierReturnQueryFilter) ([]model.SupplierReturnGet, int64, int, error) {
	var supplierReturns []model.SupplierReturnGet
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("supplier_return_no")
	query := repository.Select("supplier_returns.*, ap.invoice_no, ap.invoice_date, ap.tax_invoice_date, ap.tax_invoice_no, ap.due_date ,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup.sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = supplier_returns.updated_by").
		Joins("left join acf.account_payable ap on supplier_returns.invoice_no = ap.invoice_no AND ap.ap_type = 'I' AND ap.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user us2 on us2.user_id = supplier_returns.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = supplier_returns.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_supplier sup on sup.sup_id = supplier_returns.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("supplier_returns.cust_id=?", dataFilter.CustId)
	query.Where("supplier_returns.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("supplier_returns.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("supplier_returns.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}
	if len(dataFilter.SupId) > 0 {
		query.Where("supplier_returns.sup_id in ?", dataFilter.SupId)
		queryCount.Where("supplier_returns.sup_id in ?", dataFilter.SupId)

	}

	if dataFilter.SupplierReturnNo != "" {
		query.Where("supplier_returns.supplier_return_no=?", dataFilter.SupplierReturnNo)
		queryCount.Where("supplier_returns.supplier_return_no=?", dataFilter.SupplierReturnNo)
	}

	if len(dataFilter.Status) != 0 {
		query.Where("inv.supplier_returns.data_status in ?", dataFilter.Status)
		queryCount.Where("inv.supplier_returns.data_status in ?", dataFilter.Status)
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
		query.Order("supplier_return_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&supplierReturns).Error
	if err != nil {
		return supplierReturns, total, 0, err
	}

	err = queryCount.Model(&supplierReturns).Count(&total).Error
	if err != nil {
		log.Println("queryCount, err:", err.Error())
		return supplierReturns, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))

	return supplierReturns, total, lastPage, nil

}

func (repository *RepositorySupplierReturnImpl) FindByNo(supplierReturnNo, custId, parentCustId string) (supplierReturn model.SupplierReturnGet, err error) {
	err = repository.
		Select("supplier_returns.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup.sup_name, wh.wh_code, wh.wh_name, ap.invoice_no, ap.invoice_date, ap.tax_invoice_date, ap.tax_invoice_no, ap.due_date").
		Joins("left join sys.m_user us on us.user_id = supplier_returns.updated_by").
		Joins("left join acf.account_payable ap on supplier_returns.invoice_no = ap.invoice_no").
		Joins("left join sys.m_user us2 on us2.user_id = supplier_returns.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = supplier_returns.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_supplier sup on sup.sup_id = supplier_returns.sup_id AND sup.cust_id = ?", parentCustId).
		Where("supplier_returns.supplier_return_no = ? AND supplier_returns.cust_id=?", supplierReturnNo, custId).
		Take(&supplierReturn).Error
	return supplierReturn, err
}

func (repository *RepositorySupplierReturnImpl) FindDetBySupplierReturnNo(supplierReturnNo string, custId string) (Details []model.SupplierReturnDetGet, err error) {
	err = repository.Select("supplier_return_details.supplier_return_det_id, supplier_return_details.cust_id, supplier_return_details.supplier_return_no, "+
		"supplier_return_details.seq_no, supplier_return_details.seq_no, supplier_return_details.qty, supplier_return_details.pro_id,p.pro_code, supplier_return_details.sub_total,supplier_return_details.discount, supplier_return_details.discount_value,supplier_return_details.total, "+
		"supplier_return_details.vat_value, supplier_return_details.vat_lg_value, supplier_return_details.vat_bg_value, "+
		"p.pro_name, p.unit_id1, p.unit_id2, p.unit_id3, p.conv_unit2, p.conv_unit3, gd.qty AS invoice_qty, COALESCE ( grb.qty, gd.qty ) AS remaining_qty, "+
		"gd.unit_price1, gd.unit_price2, gd.unit_price3,supplier_return_details.item_cdn,supplier_return_details.return_reason_id, retr.return_reason_name, COALESCE (whs.qty, 0) as wh_qty").
		Joins("left join inv.supplier_returns sr on sr.supplier_return_no = supplier_return_details.supplier_return_no AND sr.cust_id = ?", custId).
		Joins("left join mst.m_product p on p.pro_id = supplier_return_details.pro_id").
		Joins("left join acf.account_payable ap on ap.invoice_no = sr.invoice_no AND ap.ap_type = 'I'").
		Joins("LEFT JOIN acf.account_payable_detail gd ON gd.account_payable_id = ap.account_payable_id AND ap.invoice_no = sr.invoice_no AND gd.item_type = 1 AND gd.pro_id = supplier_return_details.pro_id ").
		Joins("left join mst.m_return_reason retr on retr.return_reason_id = supplier_return_details.return_reason_id ").
		Joins("LEFT JOIN inv.product_invoice_balances grb ON grb.invoice_no = sr.invoice_no AND grb.pro_id = supplier_return_details.pro_id").
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = supplier_return_details.pro_id AND whs.wh_id = sr.wh_id AND whs.cust_id = ?", custId).
		Where("supplier_return_details.supplier_return_no = ? AND supplier_return_details.cust_id=?", supplierReturnNo, custId).Order("seq_no ASC").
		Find(&Details).Error
	return Details, err
}

func (repository *RepositorySupplierReturnImpl) FindSupplierReturn(dataFilter entity.ReturnSupplierQueryFilter, custId, parentCustId string) ([]model.ReturnSuppliers, int64, int, error) {
	var returnSuppliers []model.ReturnSuppliers
	var total int64
	var limit int

	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("supplier_returns.sup_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = supplier_returns.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = supplier_returns.wh_id AND wh.cust_id = ?", custId)
	queryCount.Where("supplier_returns.cust_id=? ", custId)

	query := repository.
		Select("supplier_returns.sup_id, sup.sup_code, sup_name").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = supplier_returns.sup_id AND sup.cust_id = ?", parentCustId)
	query.Where("supplier_returns.cust_id = ?", custId)

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		query.Where("supplier_returns.supplier_return_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
		queryCount.Where("supplier_returns.supplier_return_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Group("supplier_returns.sup_id, sup.sup_code, sup_name").Find(&returnSuppliers).Error
	if err != nil {
		return returnSuppliers, total, 0, err
	}
	err = queryCount.Model(&returnSuppliers).Group("supplier_returns.sup_id").Count(&total).Error
	if err != nil {
		return returnSuppliers, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return returnSuppliers, total, lastPage, nil
}
func (repository *RepositorySupplierReturnImpl) FindProductByListID(productIDs []int64) (products []model.Product, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositorySupplierReturnImpl) GetRemainingQtyInvoice(invoiceNo string, custId string, proID int64) (remainingQty model.RemainingQty, err error) {
	err = repository.Select("coalesce(pib.qty, acf.account_payable_detail.qty ) AS remaining_qty ").
		Joins("LEFT JOIN acf.account_payable ap ON ap.account_payable_id = acf.account_payable_detail.account_payable_id").
		Joins("LEFT JOIN inv.product_invoice_balances pib ON pib.invoice_no = ap.invoice_no AND acf.account_payable_detail.pro_id = pib.pro_id AND pib.cust_id = acf.account_payable_detail.cust_id  ").
		Where("ap.cust_id = ? AND ap.invoice_no = ? AND acf.account_payable_detail.pro_id = ?  ", custId, invoiceNo, proID).
		Take(&remainingQty).Error
	return remainingQty, err
}

func (repository *RepositorySupplierReturnImpl) SaveRemainingQty(c context.Context, balances model.ProductInvoiceBalances) {
	repository.model(c).Save(&balances)
}

func (repository *RepositorySupplierReturnImpl) GetRemainingProductQtyByInvoiceNo(invoiceNo string) (remainingQtyProducts []model.RemainingQtyProduct, err error) {
	err = repository.Select("acf.account_payable_detail.pro_id, acf.account_payable_detail.qty AS qty_receipt,coalesce(grb.qty, acf.account_payable_detail.qty) as qty_remaining").
		Joins("LEFT JOIN acf.account_payable ap ON ap.account_payable_id = acf.account_payable_detail.account_payable_id").
		Joins("left join inv.product_invoice_balances pib on acf.account_payable_detail.invoice_no = pib.invoice_no AND acf.account_payable_detail.pro_id = pib.pro_id AND pib.cust_id = acf.account_payable_detail.cust_id ").
		Where("ap.invoice_no = ?  ", invoiceNo).
		Find(&remainingQtyProducts).Error

	return
}

func (repository *RepositorySupplierReturnImpl) GetInvoiceFromAPReturn(custId string, invoiceNo string) (ap model.AccountPayable, err error) {
	err = repository.Select("acf.account_payable.cust_id, acf.account_payable.invoice_no").
		Where("acf.account_payable.cust_id = ? AND acf.account_payable.invoice_no = ?  ", custId, invoiceNo).
		Take(&ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, errors.New("Invoice not found")
	}
	return
}

func (repository *RepositorySupplierReturnImpl) SetInvoiceNoIsCanReturn(c context.Context, invoiceNo string, isCanReturn bool) (err error) {
	var ap model.AccountPayable
	err = repository.model(c).Model(&ap).Where("invoice_no = ?", invoiceNo).Update("is_can_return", isCanReturn).Error
	return
}
func (repository *RepositorySupplierReturnImpl) UpdateStatus(c context.Context, supplierReturnNo string, custId string, status int) error {
	var data model.SupplierReturn
	result := repository.model(c).Model(&data).Where(" cust_id = ? AND supplier_return_no=?", custId, supplierReturnNo).
		Updates(map[string]interface{}{"data_status": status})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositorySupplierReturnImpl) GetApProductList(InvNo string, custId string, proID int64) (apDet model.AccountPayableProductList, err error) {
	err = repository.Select("acf.account_payable_detail.*, mp.pro_code, mp.pro_name, coalesce(pib.qty, acf.account_payable_detail.qty) AS qty_remaining, mp.conv_unit2, mp.conv_unit3, ap.invoice_no  ").
		Joins("left join mst.m_product mp on mp.pro_id = acf.account_payable_detail.pro_id  ").
		Joins("LEFT JOIN acf.account_payable ap on acf.account_payable_detail.account_payable_id = ap.account_payable_id").
		Joins("LEFT JOIN inv.product_invoice_balances pib ON pib.invoice_no = ap.invoice_no AND pib.pro_id = acf.account_payable_detail.pro_id").
		Where("ap.invoice_no = ? AND acf.account_payable_detail.cust_id=? AND acf.account_payable_detail.pro_id=? AND acf.account_payable_detail.item_type=1", InvNo, custId, proID).
		Take(&apDet).Error
	return apDet, err
}
