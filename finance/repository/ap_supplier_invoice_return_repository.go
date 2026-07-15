package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryApSupplierInvoiceReturnImpl struct {
		*gorm.DB
	}
)
type ApSupplierInvoiceReturnRepository interface {
	Store(c context.Context, data *model.ApSupplierInvoiceReturn) error
	StoreDetail(c context.Context, data *model.AccountPayableProduct) error
	StoreApProductPromo(c context.Context, data *model.AccountPayableProductPromo) error
	UpdateGr(c context.Context, GrNo string, data model.GrUpdate) error

	FindByID(accountPayableID uint, custId, parentCustId string) (ap model.ApSuppilerInvoiceReturnList, err error)
	FindProductList(accountPayableID uint, custId string) (Details []model.AccountPayableProductList, err error)
	FindAllByCustId(dataFilter entity.ApSupplierInoviceReturnQueryFilter) ([]model.ApSuppilerInvoiceReturnList, int64, int, error)

	Delete(c context.Context, custId string, accountPayableID uint, deletedBy int64) error

	Update(c context.Context, accountPayableID uint, data model.ApSupplierInvoiceReturnupdate) error
	FindProductByListID(productIDs []int64) (products []model.Product, err error)
	GetGrByInvoiceNo(invoiceNo string, custId, parentCustId string) (gr model.GrList, err error)
	GetGrbByInvoiceNo(invoiceNo string, custId, parentCustId string, custIdParam string) (gr model.GrbList, err error)
	GetGoodReceiptByProduct(c context.Context, GrNo string, proID int64) (grDet model.GrDetList, err error)
	FindGrdetail(grNo string, custId string) (grDetails []model.GrDetList, err error)
	FindGrbdetail(grNo string, custId string, custIdParam string) (grDetails []model.GrbDetList, err error)
	UpdateProductNormal(c context.Context, accountPayableDetailID int64, Details *model.AccountPayableProduct) error
	GetReturnSupplierByDocumentNo(documentNo, custId string) (supplierReturn model.SupplierReturnGet, err error)
	GetArPayByInvoiceNo(invoiceNo, custId string) (ApPay model.ApPayJoinDet, err error)
	DeleteDetail(c context.Context, custId string, accountPayableID uint, deletedBy int64) error
	FindByDocumentNoAndType(documentNo string, apType string, custId, parentCustId string) (ap model.ApSuppilerInvoiceReturnList, err error)
	FindByInvoiceNo(invoiceNo string, custId string) (ap model.ApSuppilerInvoiceReturnList, err error)
	GetWarehouseStockFromGr(custID, documentNo string, productID int64) (whStock model.WarehouseStock, err error)
	GetWarehouseStockFromGrb(custID, documentNo string, productID int64) (whStock model.WarehouseStockGrb, err error)
	FindDetBySupplierReturnNo(documentNo string, custId string) (Details []model.SupplierReturnDetGet, err error)
	GetWarehouseStockFromReturn(custID, documentNo string, productID int64) (whStock model.WarehouseStockFromReturn, err error)
	FindAllVatInByCustId(dataFilter entity.VatExtractQueryFilter) ([]model.ApSuppilerInvoiceReturnVatInList, int64, int, error)
	FindAllVatInReturnByCustId(dataFilter entity.VatExtractQueryFilter) ([]model.ApSuppilerInvoiceReturnVatInList, int64, int, error)
	FindByDocumentsNo(documentNo []string, custId, parentCustId string) (ap []model.ApSuppilerInvoiceReturnList, err error)
	// DeleteDetailNotInIDs(c context.Context, ApSupplierInvoiceReturnNo string, IDs []int64) error
	// DeleteApSupplierInvoiceReturnMethodDetailNotInIDs(c context.Context, ApSupplierInvoiceReturnNo string, IDs []int64) error
	// UpdateDetail(c context.Context, Details *model.ApSupplierInvoiceReturnDet) error
	// UpdateApSupplierInvoiceReturnMethodDetail(c context.Context, Details *model.ApSupplierInvoiceReturnMethod) error
}

func NewApSupplierInvoiceReturnRepo(db *gorm.DB) *RepositoryApSupplierInvoiceReturnImpl {
	return &RepositoryApSupplierInvoiceReturnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryApSupplierInvoiceReturnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) Store(c context.Context, data *model.ApSupplierInvoiceReturn) error {
	query := repository.model(c).Create(data)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) StoreDetail(c context.Context, data *model.AccountPayableProduct) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) StoreApProductPromo(c context.Context, data *model.AccountPayableProductPromo) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) UpdateGr(c context.Context, GrNo string, data model.GrUpdate) error {

	result := repository.model(c).Model(&data).Where("gr_no=?", GrNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindByID(accountPayableID uint, custId, parentCustId string) (ap model.ApSuppilerInvoiceReturnList, err error) {
	err = repository.Select(`acf.account_payable.*,us.user_fullname AS created_by_name, us.user_fullname AS updated_by_name,
			sup.sup_code,sup.sup_name,sup.sup_code,
			sc.distributor_id as distributor_id,md.distributor_code as distributor_code, md.distributor_name as distributor,
			CASE 
				WHEN LEFT(acf.account_payable.document_no, 3) = 'GRB' THEN gb.po_no
				ELSE gr.po_no
			END AS po_no_doc
			`).
		Joins("left join sys.m_user us on us.user_id = acf.account_payable.updated_by").
		Joins("left join smc.m_customer sc ON sc.cust_id = acf.account_payable.cust_id").
		Joins("left join mst.m_distributor md ON md.distributor_id = sc.distributor_id").
		Joins("left join inv.gr_branch gb ON gb.gr_branch_no = acf.account_payable.document_no").
		Joins("left join inv.gr gr ON gr.gr_no = acf.account_payable.document_no").
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.account_payable.sup_id AND sup.cust_id = ?", parentCustId).
		Where("acf.account_payable.account_payable_id = ?", accountPayableID).
		Take(&ap).Error
	return ap, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindByDocumentNoAndType(documentNo string, apType string, custId, parentCustId string) (ap model.ApSuppilerInvoiceReturnList, err error) {

	err = repository.Select("acf.account_payable.*,us.user_fullname AS created_by_name, us.user_fullname AS updated_by_name,sup.sup_code,sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = acf.account_payable.updated_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.account_payable.sup_id AND sup.cust_id = ?", parentCustId).
		Where("acf.account_payable.cust_id=? AND acf.account_payable.document_no = ? AND ap_type = ?", custId, documentNo, apType).
		Take(&ap).Error
	return ap, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindByInvoiceNo(invoiceNo string, custId string) (ap model.ApSuppilerInvoiceReturnList, err error) {
	err = repository.Select("acf.account_payable.*").
		Where("acf.account_payable.cust_id = ? AND acf.account_payable.invoice_no = ?", custId, invoiceNo).
		Where("acf.account_payable.deleted_at IS NULL").
		Take(&ap).Error
	return ap, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindByDocumentsNo(documentNo []string, custId, parentCustId string) (ap []model.ApSuppilerInvoiceReturnList, err error) {

	err = repository.Select("acf.account_payable.*,us.user_fullname AS created_by_name, us.user_fullname AS updated_by_name,sup.sup_code,sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = acf.account_payable.updated_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.account_payable.sup_id AND sup.cust_id = ?", parentCustId).
		Where("acf.account_payable.cust_id=? AND acf.account_payable.document_no = ? ", custId, documentNo).
		Find(&ap).Error
	return ap, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindProductList(accountPayableID uint, custId string) (Details []model.AccountPayableProductList, err error) {
	err = repository.Select("acf.account_payable_detail.*, mp.pro_code, mp.pro_name, mp.conv_unit2, mp.conv_unit3, mp.unit_id1, mp.unit_id2, mp.unit_id3,coalesce(pib.qty, acf.account_payable_detail.qty) AS qty_remaining ").
		Joins("left join mst.m_product mp on mp.pro_id = acf.account_payable_detail.pro_id").
		Joins("join acf.account_payable ap on acf.account_payable_detail.account_payable_id = ap.account_payable_id AND ap.cust_id = ?", custId).
		Joins("LEFT JOIN inv.product_invoice_balances pib ON pib.invoice_no = ap.invoice_no AND pib.pro_id = acf.account_payable_detail.pro_id").
		Where("acf.account_payable_detail.account_payable_id = ? AND acf.account_payable_detail.cust_id=?", accountPayableID, custId).
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindAllByCustId(dataFilter entity.ApSupplierInoviceReturnQueryFilter) ([]model.ApSuppilerInvoiceReturnList, int64, int, error) {
	var ap []model.ApSuppilerInvoiceReturnList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("invoice_no")
	query := repository.Select(`acf.account_payable.*,us.user_fullname AS updated_by_name, us.user_fullname AS created_by_name ,
			sup.sup_code,sup.sup_name,sup.sup_code,
			sc.distributor_id as distributor_id,md.distributor_code as distributor_code, md.distributor_name as distributor
			`).
		Joins("left join sys.m_user us on us.user_id = acf.account_payable.updated_by").
		Joins("left join smc.m_customer sc ON sc.cust_id = acf.account_payable.cust_id").
		Joins("left join mst.m_distributor md ON md.distributor_id = sc.distributor_id").
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.account_payable.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.account_payable.cust_id=?", dataFilter.CustId)
	query.Where("acf.account_payable.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.account_payable.account_payable_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.account_payable.account_payable_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.account_payable.invoice_no=?", dataFilter.Query)
		query.Where("acf.account_payable.invoice_no=?", dataFilter.Query)
	}

	if dataFilter.DocumentNo != "" {
		queryCount.Where("acf.account_payable.document_no=?", dataFilter.DocumentNo)
		query.Where("acf.account_payable.document_no=?", dataFilter.DocumentNo)
	}

	if dataFilter.Type != "" {
		queryCount.Where("acf.account_payable.ap_type=?", dataFilter.Type)
		query.Where("acf.account_payable.ap_type=?", dataFilter.Type)
	}

	if dataFilter.SuppId != 0 {
		queryCount.Where("acf.account_payable.sup_id=?", dataFilter.SuppId)
		query.Where("acf.account_payable.sup_id=?", dataFilter.SuppId)
	}

	if dataFilter.InvoiceNo != "" {
		queryCount.Where("acf.account_payable.invoice_no=?", dataFilter.InvoiceNo)
		query.Where("acf.account_payable.invoice_no=?", dataFilter.InvoiceNo)
	}

	if dataFilter.ExcludeEmptyInvoice {
		query.Where("acf.account_payable.invoice_no <> '' AND acf.account_payable.is_can_return = true")
		queryCount.Where("acf.account_payable.invoice_no <> '' AND acf.account_payable.is_can_return = true")
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
		query.Order("invoice_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ap).Error
	if err != nil {
		return ap, total, 0, err
	}
	err = queryCount.Model(&ap).Count(&total).Error
	if err != nil {
		return ap, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ap, total, lastPage, nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) Delete(c context.Context, custId string, accountPayableID uint, deletedBy int64) error {
	var data model.ApSupplierInvoiceReturn
	result := repository.model(c).Model(&data).Where("account_payable_id=? AND cust_id = ? AND is_del= ? ", accountPayableID, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) DeleteDetail(c context.Context, custId string, accountPayableID uint, deletedBy int64) error {
	var data model.AccountPayableProduct
	result := repository.model(c).Where("account_payable_id=? AND cust_id = ?", accountPayableID, custId).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) Update(c context.Context, accountPayableID uint, data model.ApSupplierInvoiceReturnupdate) error {

	result := repository.model(c).Model(&data).Where("account_payable_id = ?", accountPayableID).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindProductByListID(productIDs []int64) (products []model.Product, err error) {
	err = repository.
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) GetGrByInvoiceNo(grNo string, custId, parentCustId string) (gr model.GrList, err error) {
	err = repository.
		Select("gr.*, us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = gr.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gr.closed_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId).
		Where("gr.gr_no = ? AND gr.cust_id=?", grNo, custId).
		Take(&gr).Error
	if err == gorm.ErrRecordNotFound {
		return gr, errors.New("Good receipt not found")
	}
	return gr, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) GetGrbByInvoiceNo(grNo string, custId, parentCustId string, custIdParam string) (grb model.GrbList, err error) {
	err = repository.
		Select("gr_branch.*, us.user_fullname AS updated_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = gr_branch.updated_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr_branch.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr_branch.wh_id AND wh.cust_id = ?", custId).
		Where("gr_branch.gr_branch_no = ? AND gr_branch.cust_id=?", grNo, custIdParam).
		Take(&grb).Error
	if err == gorm.ErrRecordNotFound {
		return grb, errors.New("Good receipt branch not found")
	}
	return grb, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) GetGoodReceiptByProduct(c context.Context, GrNo string, proID int64) (grDet model.GrDetList, err error) {
	err = repository.
		Select(`gr_det.*,
			pd.pro_code, pd.pro_name, pd.purch_price1, pd.purch_price2, pd.purch_price3, pd.conv_unit2, pd.conv_unit3`).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_det.pro_id").
		Where("gr_det.gr_no = ? AND gr_det.pro_id = ?", GrNo, proID).
		Take(&grDet).Error
	return grDet, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindGrdetail(grNo string, custId string) (grDetails []model.GrDetList, err error) {
	err = repository.
		Select(`gr_det.*,
			pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3, pd.conv_unit2, pd.conv_unit3,
			mapd.disc_p, coalesce(grb.qty, gr_det.qty) AS qty_remaining, COALESCE (whs.qty, 0) as wh_qty, apd.discount `).
		Joins("LEFT JOIN inv.gr on gr.gr_no = gr_det.gr_no AND gr.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_det.pro_id").
		Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN acf.account_payable_discounts apd ON apd.pro_id = gr_det.pro_id AND apd.cust_id = ? AND apd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN inv.good_receipt_balances grb ON grb.gr_no = gr_det.gr_no AND grb.pro_id = gr_det.pro_id", grNo).
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_det.pro_id AND whs.wh_id = gr.wh_id").
		Where("gr_det.gr_no = ? AND gr_det.cust_id = ?", grNo, custId).Order("gr_det.seq_no ASC").
		Find(&grDetails).Error
	return grDetails, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindGrbdetail(grNo string, custId string, custIdParam string) (grDetails []model.GrbDetList, err error) {
	err = repository.
		Select(`gr_branch_det.*,
			pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3, pd.conv_unit2, pd.conv_unit3,
			mapd.disc_p, COALESCE (whs.qty, 0) as wh_qty, apd.discount `).
		Joins("LEFT JOIN inv.gr_branch on gr_branch.gr_branch_no = gr_branch_det.gr_branch_no AND gr_branch.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_branch_det.pro_id").
		Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_branch_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN acf.account_payable_discounts apd ON apd.pro_id = gr_branch_det.pro_id AND apd.cust_id = ? AND apd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_branch_det.pro_id AND whs.wh_id = gr_branch.wh_id").
		Where("gr_branch_det.gr_branch_no = ? AND gr_branch_det.cust_id = ?", grNo, custIdParam).Order("gr_branch_det.seq_no ASC").
		Find(&grDetails).Error
	return grDetails, err
}

// func (repository *RepositoryApSupplierInvoiceReturnImpl) DeleteDetailNotInIDs(c context.Context, ApSupplierInvoiceReturnNo string, IDs []int64) error {
// 	var Details model.ApSupplierInvoiceReturnDet
// 	err := repository.model(c).Where("ap_pay_no=? AND ap_pay_det_id not in (?) ", ApSupplierInvoiceReturnNo, IDs).Delete(&Details).Error
// 	return err
// }

// func (repository *RepositoryApSupplierInvoiceReturnImpl) DeleteApSupplierInvoiceReturnMethodDetailNotInIDs(c context.Context, ApSupplierInvoiceReturnNo string, IDs []int64) error {
// 	var Details model.ApSupplierInvoiceReturnMethod
// 	err := repository.model(c).Where("ap_pay_no=? AND ap_pay_method_id not in (?) ", ApSupplierInvoiceReturnNo, IDs).Delete(&Details).Error
// 	return err
// }

func (repository *RepositoryApSupplierInvoiceReturnImpl) UpdateProductNormal(c context.Context, accountPayableDetailID int64, Details *model.AccountPayableProduct) error {
	result := repository.model(c).Where("account_payable_detail_id = ?", accountPayableDetailID).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryApSupplierInvoiceReturnImpl) GetReturnSupplierByDocumentNo(invoiceNo, custId string) (supplierReturn model.SupplierReturnGet, err error) {
	err = repository.Select("*").Where("supplier_returns.supplier_return_no = ? AND supplier_returns.cust_id=?", invoiceNo, custId).Take(&supplierReturn).Error
	if err == gorm.ErrRecordNotFound {
		return supplierReturn, errors.New("return supplier not found")
	}
	return supplierReturn, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) GetArPayByInvoiceNo(invoiceNo, custId string) (ApPay model.ApPayJoinDet, err error) {
	err = repository.Select("app.*").
		Joins("LEFT join account_payable_payment_detail appd on app.account_payable_payment_no = appd.account_payable_payment_no").
		Where("supplier_returns.invoice_no = ? AND supplier_returns.cust_id=?", invoiceNo, custId).Take(&ApPay).Error
	return ApPay, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) GetWarehouseStockFromGr(custID, documentNo string, productID int64) (whStock model.WarehouseStock, err error) {
	err = repository.
		Select("ws.pro_id, ws.qty").
		Joins("left join inv.warehouse_stock ws on ws.wh_id = inv.gr.wh_id").
		Where("ws.cust_id = ? AND inv.gr.gr_no = ? AND ws.pro_id = ?", custID, documentNo, productID).Take(&whStock).Error
	return whStock, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) GetWarehouseStockFromGrb(custID, documentNo string, productID int64) (whStock model.WarehouseStockGrb, err error) {
	err = repository.
		Select("ws.pro_id, ws.qty").
		Joins("left join inv.warehouse_stock ws on ws.wh_id = inv.gr_branch.wh_id").
		Where("ws.cust_id = ? AND inv.gr_branch.gr_branch_no = ? AND ws.pro_id = ?", custID, documentNo, productID).Take(&whStock).Error
	return whStock, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) GetWarehouseStockFromReturn(custID, documentNo string, productID int64) (whStock model.WarehouseStockFromReturn, err error) {
	err = repository.
		Select("ws.pro_id, ws.qty").
		Joins("left join inv.warehouse_stock ws on ws.wh_id = inv.supplier_returns.wh_id").
		Where("ws.cust_id = ? AND inv.supplier_returns.supplier_return_no = ? AND ws.pro_id = ?", custID, documentNo, productID).Take(&whStock).Error
	return whStock, err
}

//	func (repository *RepositoryApSupplierInvoiceReturnImpl) UpdateApSupplierInvoiceReturnMethodDetail(c context.Context, Details *model.ApSupplierInvoiceReturnMethod) error {
//		result := repository.model(c).Updates(&Details)
//		if result.Error != nil {
//			return result.Error
//		}
//		// if result.RowsAffected == 0 {
//		// 	return errors.New("no rows affected")
//		// }
//		return nil
//	}
func (repository *RepositoryApSupplierInvoiceReturnImpl) FindDetBySupplierReturnNo(documentNo string, custId string) (Details []model.SupplierReturnDetGet, err error) {
	err = repository.Select("supplier_return_details.supplier_return_det_id, supplier_return_details.cust_id, supplier_return_details.supplier_return_no, "+
		"supplier_return_details.seq_no, supplier_return_details.seq_no, supplier_return_details.qty, supplier_return_details.pro_id,p.pro_code, supplier_return_details.sub_total,supplier_return_details.discount, supplier_return_details.discount_value,supplier_return_details.total, "+
		"supplier_return_details.vat_value, supplier_return_details.vat_lg_value, supplier_return_details.vat_bg_value, "+
		"p.pro_name, p.unit_id1, p.unit_id2, p.unit_id3, p.conv_unit2, p.conv_unit3, gd.qty AS invoice_qty, COALESCE ( grb.qty, gd.qty ) AS remaining_qty, "+
		"gd.unit_price1, gd.unit_price2, gd.unit_price3,supplier_return_details.item_cdn,supplier_return_details.return_reason_id, retr.return_reason_name, COALESCE (whs.qty, 0) as wh_qty").
		Joins("left join inv.supplier_returns sr on sr.supplier_return_no = supplier_return_details.supplier_return_no AND sr.cust_id = ?", custId).
		Joins("left join mst.m_product p on p.pro_id = supplier_return_details.pro_id").
		Joins("left join acf.account_payable ap on ap.invoice_no = sr.invoice_no").
		Joins("left join acf.account_payable_detail gd on gd.pro_id = supplier_return_details.pro_id AND ap.invoice_no = sr.invoice_no AND gd.item_type = 1 ").
		Joins("left join mst.m_return_reason retr on retr.return_reason_id = supplier_return_details.return_reason_id ").
		Joins("LEFT JOIN inv.product_invoice_balances grb ON grb.invoice_no = sr.invoice_no AND grb.pro_id = supplier_return_details.pro_id").
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = supplier_return_details.pro_id AND whs.wh_id = sr.wh_id").
		Where("supplier_return_details.supplier_return_no = ? AND supplier_return_details.cust_id=?", documentNo, custId).Order("seq_no ASC").
		Find(&Details).Error
	return Details, err
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindAllVatInByCustId(dataFilter entity.VatExtractQueryFilter) ([]model.ApSuppilerInvoiceReturnVatInList, int64, int, error) {
	var ap []model.ApSuppilerInvoiceReturnVatInList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}
	subquery := repository.Table("acf.account_payable").
		Select(`
			DISTINCT ON (acf.account_payable.account_payable_id)
			acf.account_payable.*,
			us.user_fullname AS updated_by_name,
			us.user_fullname AS created_by_name,
			CASE WHEN ved.vat_extract_id IS NULL THEN 'not extracted' ELSE 'extracted' END AS extract_status,
			ve.created_at AS extracted_at,
			sup.sup_code,
			sup.sup_name,
			sup.address1 AS address,
			sup.tax_no AS npwp
		`).
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.account_payable.updated_by").
		Joins("LEFT JOIN acf.vat_extract_details ved ON ved.reference_id = acf.account_payable.account_payable_id").
		Joins("LEFT JOIN acf.vat_extracts ve ON ve.vat_extract_id = ved.vat_extract_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = acf.account_payable.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)
	if dataFilter.From != nil && dataFilter.To != nil {
		subquery.Where("acf.account_payable.tax_invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.InvoiceType != "" {
		subquery.Where("acf.account_payable.ap_type=?", dataFilter.InvoiceType)
	}

	if dataFilter.ExtractionStatus != "" {
		if dataFilter.ExtractionStatus == "E" {
			subquery.Where("ved.vat_extract_id is not null")
		} else if dataFilter.ExtractionStatus == "NE" {
			subquery.Where("ved.vat_extract_id is null")
		}
	}

	if len(dataFilter.TransactionID) > 0 {
		subquery.Where("acf.account_payable.account_payable_id in ?", dataFilter.TransactionID)
	}

	subquery.Where("acf.account_payable.cust_id=?", dataFilter.CustId)
	subquery.Order("acf.account_payable.account_payable_id, ve.created_at DESC")

	if dataFilter.From != nil && dataFilter.To != nil {
		subquery.Where("acf.account_payable.account_payable_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}
	query := repository.Table("(?) AS subquery", subquery)
	queryCount := repository.Table("(?) AS subquery", subquery)
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
		query.Order("invoice_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ap).Error
	if err != nil {
		return ap, total, 0, err
	}
	err = queryCount.Model(&ap).Count(&total).Error
	if err != nil {
		return ap, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ap, total, lastPage, nil
}

func (repository *RepositoryApSupplierInvoiceReturnImpl) FindAllVatInReturnByCustId(dataFilter entity.VatExtractQueryFilter) ([]model.ApSuppilerInvoiceReturnVatInList, int64, int, error) {
	var ap []model.ApSuppilerInvoiceReturnVatInList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}
	subquery := repository.Table("acf.account_payable").
		Select(`
			DISTINCT ON (acf.account_payable.account_payable_id)
			acf.account_payable.*,
			us.user_fullname AS updated_by_name,
			us.user_fullname AS created_by_name,
			CASE WHEN ved.vat_extract_id IS NULL THEN 'not extracted' ELSE 'extracted' END AS extract_status,
			ve.created_at AS extracted_at,
			sup.sup_code,
			sup.sup_name,
			sup.address1 AS address,
			sup.tax_no AS npwp
		`).
		Joins("LEFT JOIN sys.m_user us ON us.user_id = acf.account_payable.updated_by").
		Joins("LEFT JOIN acf.vat_extract_details ved ON ved.reference_id = acf.account_payable.account_payable_id").
		Joins("LEFT JOIN acf.vat_extracts ve ON ve.vat_extract_id = ved.vat_extract_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = acf.account_payable.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)
	if dataFilter.From != nil && dataFilter.To != nil {
		subquery.Where("acf.account_payable.tax_invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.InvoiceType != "" {
		subquery.Where("acf.account_payable.ap_type=?", dataFilter.InvoiceType)
	}

	if dataFilter.ExtractionStatus != "" {
		if dataFilter.ExtractionStatus == "E" {
			subquery.Where("ved.vat_extract_id is not null")
		} else if dataFilter.ExtractionStatus == "NE" {
			subquery.Where("ved.vat_extract_id is null")
		}
	}

	if len(dataFilter.TransactionID) > 0 {
		subquery.Where("acf.account_payable.account_payable_id in ?", dataFilter.TransactionID)
	}

	subquery.Where("acf.account_payable.cust_id=?", dataFilter.CustId)
	subquery.Order("acf.account_payable.account_payable_id, ve.created_at DESC")

	if dataFilter.From != nil && dataFilter.To != nil {
		subquery.Where("acf.account_payable.account_payable_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}
	query := repository.Table("(?) AS subquery", subquery)
	queryCount := repository.Table("(?) AS subquery", subquery)
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
		query.Order("invoice_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ap).Error
	if err != nil {
		return ap, total, 0, err
	}
	err = queryCount.Model(&ap).Count(&total).Error
	if err != nil {
		return ap, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ap, total, lastPage, nil
}
