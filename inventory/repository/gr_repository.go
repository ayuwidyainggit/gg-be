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
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryGrImpl struct {
		*gorm.DB
	}
)

type GrRepository interface {
	Store(c context.Context, data *model.Gr) error
	FindByNo(grNo string, custId, parentCustId string) (gr model.GrList, err error)
	FindAllByCustId(dataFilter entity.GrQueryFilter, custId, parentCustId string) ([]model.GrList, int64, int, error)
	Update(c context.Context, grNo string, data model.Gr) error
	Delete(c context.Context, custId string, grNo string, deletedBy int64) error
	FindGrdetail(grNo string, custId string) (grDetails []model.GrDetList, err error)
	FindGrdetailForUpdateWhStock(grNo string, custId string) (grDetails []model.GrDetJoinGrList, err error)
	UpdateGrDetail(c context.Context, grDetails *model.GrDet) error
	CreateGrDetail(c context.Context, grDetails *model.GrDetCreate) (*model.GrDetCreate, error)
	DeleteGrDetailNotInIDs(c context.Context, grNo string, IDs []int) error
	DeleteGrDetailByGrNo(c context.Context, grNo string) error
	DeleteStockNotInRefIds(c context.Context, custId, trNo string, newRefIds []int64) error
	DeleteStockInRefIds(c context.Context, custId, trNo string, oldRefIds []int64) error
	StoreWhStock(c context.Context, data *model.WhStock) error
	UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error
	StoreStock(c context.Context, data *model.Stock) error
	StoreProductCogs(c context.Context, data *model.ProductCogs) error
	FindQtyWhStock(custId string, proId, whId int64) (whStock model.WhStockList, err error)
	FindCogsProductDist(custId string, proId int64) (productDist model.ProductDist, err error)
	UpdateProductDist(c context.Context, custId string, proId int64, data model.ProductDist) error
	GetByInvoiceNo(invoiceNo string, custId, parentCustId string) (gr model.GrList, err error)
	FindSupplierGr(dataFilter entity.GrSupplierQueryFilter, custId, parentCustId string) ([]model.GrSupplier, int64, int, error)
	FindProductByListID(custID, parentCustId string, distributorID int64, productIDs []int64) (products []model.Product, err error)
	FindWarehouseGr(dataFilter entity.GrWarehouseQueryFilter, custId, parentCustId string) ([]model.GrWarehouse, int64, int, error)
	FindGrdetailWithDiscount(grNo string, custId string) (grDetails []model.GrDetList, err error)
	FindDistributorGr(custId string) (distrobutorGr []model.DistributorGr, err error)
	FindAllGrBranchByCustIdSupId(dataFilter entity.GrLookupQueryFilter, custId, parentCustId string) ([]model.GrBranchLookup, int64, int, error)
	FindAllGrByCustIdSupId(dataFilter entity.GrLookupQueryFilter, custId, parentCustId string) ([]model.GrLookup, int64, int, error)
	CheckReportInProgress(ctx context.Context, prefix string) (bool, error)
	FindGrDetailForDownload(grNo string, custId, parentCustId string) (gr model.GrList, grDetails []model.GrDetList, err error)
}

func NewGrRepo(db *gorm.DB) *RepositoryGrImpl {
	return &RepositoryGrImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryGrImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryGrImpl) FindByNo(grNo string, custId, parentCustId string) (gr model.GrList, err error) {
	err = repository.
		Select("gr.*, us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = gr.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gr.closed_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId).
		Where("gr.gr_no = ? AND gr.cust_id=?", grNo, custId).
		Take(&gr).Error
	return gr, err
}

func (repository *RepositoryGrImpl) GetByInvoiceNo(invoiceNo string, custId, parentCustId string) (gr model.GrList, err error) {
	err = repository.
		Select("gr.*, us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = gr.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gr.closed_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId).
		Where("gr.invoice_no = ? AND gr.cust_id=?", invoiceNo, custId).
		Take(&gr).Error
	return gr, err
}

func (repository *RepositoryGrImpl) Store(c context.Context, data *model.Gr) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGrImpl) FindAllByCustId(dataFilter entity.GrQueryFilter, custId, parentCustId string) ([]model.GrList, int64, int, error) {
	var grs []model.GrList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_no").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = gr.updated_by").
		Joins("LEFT JOIN sys.m_user us2 on us2.user_id = gr.closed_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId)
	queryCount.Where("gr.cust_id=?", custId)

	query := repository.
		Select("gr.*, us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = gr.updated_by").
		Joins("LEFT JOIN sys.m_user us2 on us2.user_id = gr.closed_by").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId)

	query.Where("gr.cust_id = ?", custId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("gr.gr_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("gr.gr_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.SupId) > 0 {
		query.Where("gr.sup_id in ?", dataFilter.SupId)
		queryCount.Where("gr.sup_id in ?", dataFilter.SupId)

	}

	if dataFilter.IsAp != nil {
		if *dataFilter.IsAp == 1 {
			query.Where("gr.ap_no IS NULL")
			queryCount.Where("gr.ap_no IS NULL")
		}
		if *dataFilter.IsAp == 2 {
			query.Where("gr.ap_no IS NOT NULL")
			queryCount.Where("gr.ap_no IS NOT NULL")
		}
	}

	if dataFilter.DataStatus != nil {
		if *dataFilter.DataStatus == 1 {
			query.Where("gr.data_status = 1")
			queryCount.Where("gr.data_status = 1")
		}
		if *dataFilter.DataStatus == 2 {
			query.Where("gr.data_status = 2")
			queryCount.Where("gr.data_status = 2")
		}
	}

	// fmt.Println("=====>", dataFilter.Query)
	if dataFilter.InvoiceNo != "" {
		query.Where("gr.invoice_no = ?", dataFilter.InvoiceNo)
		queryCount.Where("gr.invoice_no = ?", dataFilter.InvoiceNo)
	}

	if dataFilter.ExcludeEmptyInvoice {
		query.Where("gr.invoice_no <> '' AND is_can_return = true")
		queryCount.Where("gr.invoice_no <> '' AND is_can_return = true")
	}

	if dataFilter.WhID != nil {
		query.Where("gr.wh_id = ?", *dataFilter.WhID)
		queryCount.Where("gr.wh_id = ?", *dataFilter.WhID)
	}

	if dataFilter.GrType != nil {
		query.Where("gr.wh_id = ?", dataFilter.GrType)
		queryCount.Where("gr.wh_id = ?", dataFilter.GrType)
	}

	if dataFilter.GrNo != "" {
		query.Where(`
			gr.gr_no ILIKE '%` + dataFilter.GrNo + `%' 	
		`)
		queryCount.Where(`
			gr.gr_no ILIKE '%` + dataFilter.GrNo + `%'
			
		`)
	}

	if dataFilter.Query != "" {
		query.Where(`(
			gr.gr_no ILIKE '%` + dataFilter.Query + `%' OR
			sup.sup_name ILIKE '%` + dataFilter.Query + `%' OR
			gr.delivery_no ILIKE '%` + dataFilter.Query + `%' OR
			gr.invoice_no ILIKE '%` + dataFilter.Query + `%'
		)`)
		queryCount.Where(`(
			gr.gr_no ILIKE '%` + dataFilter.Query + `%' OR
			sup.sup_name ILIKE '%` + dataFilter.Query + `%' OR
			gr.delivery_no ILIKE '%` + dataFilter.Query + `%' OR
			gr.invoice_no ILIKE '%` + dataFilter.Query + `%'
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
		query.Order("gr.gr_no DESC")
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

func (repository *RepositoryGrImpl) Update(c context.Context, grNo string, data model.Gr) error {
	result := repository.model(c).Model(&data).Where("gr_no=?", grNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryGrImpl) Delete(c context.Context, custId string, grNo string, deletedBy int64) error {
	var data model.Gr
	result := repository.model(c).Model(&data).Where("gr_no=? AND cust_id = ? AND is_del= ? ", grNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryGrImpl) FindGrdetail(grNo string, custId string) (grDetails []model.GrDetList, err error) {
	err = repository.
		Select(`gr_det.*,
			pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3, pd.conv_unit2, pd.conv_unit3,
			mapd.disc_p, coalesce(grb.qty, gr_det.qty) AS qty_remaining, COALESCE (whs.qty, 0) as wh_qty `).
		Joins("LEFT JOIN inv.gr on gr.gr_no = gr_det.gr_no AND gr.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_det.pro_id").
		Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN inv.good_receipt_balances grb ON grb.gr_no = gr_det.gr_no AND grb.pro_id = gr_det.pro_id", grNo).
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_det.pro_id AND whs.wh_id = gr.wh_id AND whs.cust_id = ? ", custId).
		Where("gr_det.gr_no = ? AND gr_det.cust_id = ?", grNo, custId).Order("gr_det.seq_no ASC").
		Find(&grDetails).Error
	return grDetails, err
}

func (repository *RepositoryGrImpl) FindGrdetailWithDiscount(grNo string, custId string) (grDetails []model.GrDetList, err error) {
	err = repository.
		Select(`gr_det.*,
			pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3, pd.conv_unit2, pd.conv_unit3,
			mapd.disc_p, coalesce(grb.qty, gr_det.qty) AS qty_remaining, COALESCE (whs.qty, 0) as wh_qty, apd.discount `).
		Joins("LEFT JOIN inv.gr on gr.gr_no = gr_det.gr_no AND gr.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_det.pro_id").
		Joins("LEFT JOIN acf.m_ap_disc mapd ON mapd.pro_id = gr_det.pro_id AND mapd.cust_id = ? AND mapd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN acf.account_payable_discounts apd ON apd.pro_id = gr_det.pro_id AND apd.cust_id = ? AND apd.deleted_at IS NULL", custId).
		Joins("LEFT JOIN inv.good_receipt_balances grb ON grb.gr_no = gr_det.gr_no AND grb.pro_id = gr_det.pro_id", grNo).
		Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = gr_det.pro_id AND whs.wh_id = gr.wh_id AND whs.cust_id = ? ", custId).
		Where("gr_det.gr_no = ? AND gr_det.cust_id = ?", grNo, custId).Order("gr_det.seq_no ASC").
		Find(&grDetails).Error
	return grDetails, err
}

func (repository *RepositoryGrImpl) FindGrdetailForUpdateWhStock(grNo string, custId string) (grDetails []model.GrDetJoinGrList, err error) {
	err = repository.
		Select(`gr_det.gr_det_id, gr.wh_id, gr_det.pro_id, gr_det.qty`).
		Joins("LEFT JOIN inv.gr gr ON gr.gr_no = gr_det.gr_no AND gr.cust_id = ?", custId).
		Where("gr_det.gr_no = ? AND gr_det.cust_id = ?", grNo, custId).Order("gr_det.seq_no ASC").
		Find(&grDetails).Error
	return grDetails, err
}

func (repository *RepositoryGrImpl) CreateGrDetail(c context.Context, grDetails *model.GrDetCreate) (*model.GrDetCreate, error) {
	result := repository.model(c).Create(grDetails)
	if result.Error != nil {
		log.Println("CreateGrDetail, result.Error:", structs.StructToJson(result.Error))
		return grDetails, result.Error
	}
	if result.RowsAffected == 0 {
		return grDetails, errors.New("no rows affected")
	}

	return grDetails, nil
}

func (repository *RepositoryGrImpl) UpdateGrDetail(c context.Context, grDetails *model.GrDet) error {
	result := repository.model(c).Updates(&grDetails)
	if result.Error != nil {
		log.Println("UpdateGrDetail, result:", structs.StructToJson(result.Error))
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryGrImpl) DeleteGrDetailNotInIDs(c context.Context, grNo string, IDs []int) error {
	var grDetails model.GrDet
	err := repository.model(c).Where("gr_no=? AND gr_det_id not in (?) ", grNo, IDs).Delete(&grDetails).Error
	return err
}

func (repository *RepositoryGrImpl) DeleteGrDetailByGrNo(c context.Context, grNo string) error {
	var grDetails model.GrDet
	err := repository.model(c).Where("gr_no = ?", grNo).Delete(&grDetails).Error
	return err
}

func (repository *RepositoryGrImpl) FindQtyWhStock(custId string, proId, whId int64) (whs model.WhStockList, err error) {
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

func (repository *RepositoryGrImpl) FindCogsProductDist(custId string, proId int64) (productDist model.ProductDist, err error) {
	err = repository.
		Select("cogs").
		Where("cust_id = ? AND pro_id = ?", custId, proId).
		Take(&productDist).Error
	return productDist, err
}

func (repository *RepositoryGrImpl) UpdateProductDist(c context.Context, custId string, proId int64, data model.ProductDist) error {
	result := repository.model(c).Model(&data).Where("cust_id = ? AND pro_id = ?", custId, proId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryGrImpl) StoreWhStock(c context.Context, data *model.WhStock) error {
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

func (repository *RepositoryGrImpl) UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error {
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

func (repository *RepositoryGrImpl) StoreStock(c context.Context, data *model.Stock) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGrImpl) StoreProductCogs(c context.Context, data *model.ProductCogs) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGrImpl) DeleteStockNotInRefIds(c context.Context, custId, trNo string, newRefIds []int64) error {
	var stock model.Stock
	err := repository.model(c).Where("cust_id = ? AND tr_no = ? AND ref_det_id NOT IN (?) ", custId, trNo, newRefIds).Delete(&stock).Error
	return err
}

func (repository *RepositoryGrImpl) DeleteStockInRefIds(c context.Context, custId, trNo string, oldRefIds []int64) error {
	var stock model.Stock
	err := repository.model(c).Where("cust_id = ? AND tr_no = ? AND ref_det_id IN (?) ", custId, trNo, oldRefIds).Delete(&stock).Error
	return err
}

func (repository *RepositoryGrImpl) FindSupplierGr(dataFilter entity.GrSupplierQueryFilter, custId, parentCustId string) ([]model.GrSupplier, int64, int, error) {
	var grSuppliers []model.GrSupplier
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr.sup_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId)
	queryCount.Where("gr.cust_id=? ", custId)

	query := repository.
		Select("gr.sup_id, sup.sup_code, sup_name").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId)
	query.Where("gr.cust_id = ?", custId)

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
		query.Where("gr.gr_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
		queryCount.Where("gr.gr_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
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
		query.Order("gr.sup_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Group("gr.sup_id, sup.sup_code, sup.sup_name").Find(&grSuppliers).Error
	if err != nil {
		return grSuppliers, total, 0, err
	}
	err = queryCount.Model(&grSuppliers).Distinct("gr.sup_id").Count(&total).Error
	if err != nil {
		return grSuppliers, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return grSuppliers, total, lastPage, nil
}

func (repository *RepositoryGrImpl) FindProductByListID(custID, parentCustId string, distributorID int64, productIDs []int64) (products []model.Product, err error) {
	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	timeNowJkt := time.Now().In(asiaJkt)
	currentDate := time.Date(timeNowJkt.Year(), timeNowJkt.Month(), timeNowJkt.Day(), 0, 0, 0, 0, timeNowJkt.Location())
	transDate := currentDate.Format("2006-01-02")

	query := repository.
		Select(`
	p.cust_id, p.pro_id,
        p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5,
        p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5,
       	CASE WHEN mtp_mg_pr.purch_price1 IS NULL THEN p.purch_price1 ELSE mtp_mg_pr.purch_price1 END AS purch_price1,
		CASE WHEN mtp_mg_pr.purch_price2 is NULL THEN p.purch_price2 ELSE mtp_mg_pr.purch_price2 END AS purch_price2,
		CASE WHEN mtp_mg_pr.purch_price3 is NULL THEN p.purch_price3 ELSE mtp_mg_pr.purch_price3 END AS purch_price3,
		p.purch_price4, 
		p.purch_price5,
		CASE WHEN (mtp.sell_price1=0 or mtp.sell_price1 is null) THEN (CASE WHEN (mtp_mg_pr_sell.sell_price1=0 or mtp_mg_pr_sell.sell_price1 is null) THEN p.sell_price1 ELSE mtp_mg_pr_sell.sell_price1 END) ELSE mtp.sell_price1 END AS sell_price1,
		CASE WHEN (mtp.sell_price2=0 or mtp.sell_price2 is null) THEN (CASE WHEN (mtp_mg_pr_sell.sell_price2=0 or mtp_mg_pr_sell.sell_price2 is null) THEN p.sell_price2 ELSE mtp_mg_pr_sell.sell_price2 END) ELSE mtp.sell_price2 END AS sell_price2,
		CASE WHEN (mtp.sell_price3=0 or mtp.sell_price3 is null) THEN (CASE WHEN (mtp_mg_pr_sell.sell_price3=0 or mtp_mg_pr_sell.sell_price3 is null) THEN p.sell_price3 ELSE mtp_mg_pr_sell.sell_price3 END) ELSE mtp.sell_price3 END AS sell_price3
        `).
		Table("mst.m_product p").
		Joins("LEFT JOIN mst.m_distributor dist ON dist.cust_id = ? AND dist.distributor_id = ?", parentCustId, distributorID).
		Joins("LEFT JOIN smc.m_customer cus ON cus.cust_id = ?", custID).
		Joins("LEFT JOIN mst.m_transaction_price mtp ON mtp.pro_id = p.pro_id AND mtp.cust_id = ? AND mtp.outlet_id = ? AND (? BETWEEN mtp.start_date AND mtp.end_date) ", custID, 0, transDate).
		Joins("LEFT JOIN LATERAL (SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3 FROM mst.m_transaction_price mtp_mg_pr WHERE mtp_mg_pr.cust_id = ? AND mtp_mg_pr.pro_id = p.pro_id AND mtp_mg_pr.start_date <= ? AND (mtp_mg_pr.distributor_id = (CASE WHEN mtp_mg_pr.coverage='N' THEN 0 ELSE dist.distributor_id END) OR mtp_mg_pr.price_group_reff = dist.dist_price_grp_id)ORDER BY mtp_mg_pr.start_date DESC LIMIT 1) mtp_mg_pr ON true", parentCustId, transDate).
		Joins("LEFT JOIN LATERAL (SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3 FROM mst.m_transaction_price mtp_mg_pr_sell WHERE mtp_mg_pr_sell.cust_id = ? AND mtp_mg_pr_sell.pro_id = p.pro_id and mtp_mg_pr_sell.source = 10 AND mtp_mg_pr_sell.start_date <= ? AND (mtp_mg_pr_sell.distributor_id = (CASE WHEN mtp_mg_pr_sell.coverage = 'N' THEN 0 ELSE dist.distributor_id END) OR mtp_mg_pr_sell.price_group_reff = dist.dist_price_grp_id) ORDER BY mtp_mg_pr_sell.start_date DESC LIMIT 1) mtp_mg_pr_sell ON true", parentCustId, transDate).
		Where("p.pro_id in ?", productIDs)

	err = query.Scan(&products).Error

	return products, err
}

func (repository *RepositoryGrImpl) FindWarehouseGr(dataFilter entity.GrWarehouseQueryFilter, custId, parentCustId string) ([]model.GrWarehouse, int64, int, error) {
	var grWarehouses []model.GrWarehouse
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr.wh_id").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId)
	queryCount.Where("gr.cust_id=? ", custId)

	query := repository.
		Select("gr.wh_id, wh.wh_code, wh.wh_name").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId)
	query.Where("gr.cust_id = ?", custId)

	if len(dataFilter.SupID) > 0 {
		queryCount.Where("inv.gr.sup_id IN ?", dataFilter.SupID)
		query.Where("inv.gr.sup_id IN ?", dataFilter.SupID)
	}

	if dataFilter.Query != "" {
		query.Where(`(
			wh.wh_code ILIKE '%` + dataFilter.Query + `%' OR
			wh.wh_name ILIKE '%` + dataFilter.Query + `%'
		)`)
		queryCount.Where(`(
			wh.wh_code ILIKE '%` + dataFilter.Query + `%' OR
			wh.wh_name ILIKE '%` + dataFilter.Query + `%'
		)`)
	}

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		query.Where("gr.gr_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
		queryCount.Where("gr.gr_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.StartDate), str.UnixTimestampToUtcTime(*dataFilter.EndDate))
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
		query.Order("gr.wh_id ASC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Group("gr.wh_id, wh.wh_code, wh.wh_name").Find(&grWarehouses).Error
	if err != nil {
		return grWarehouses, total, 0, err
	}
	err = queryCount.Model(&grWarehouses).Distinct("gr.wh_id").Count(&total).Error
	if err != nil {
		return grWarehouses, total, 0, err
	}
	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return grWarehouses, total, lastPage, nil
}

func (repository *RepositoryGrImpl) FindDistributorGr(custId string) (distrobutorGr []model.DistributorGr, err error) {
	query := repository.Select("smc.cust_id, smc.distributor_id, distributor_code, distributor_name").
		Joins("LEFT JOIN smc.m_customer smc ON smc.distributor_id = mst.m_distributor.distributor_id")

	if len(custId) < 10 {
		query = query.Where("smc.cust_id LIKE ?", custId+"%")
	} else {
		query = query.Where("smc.cust_id = ?", custId)
	}

	err = query.Find(&distrobutorGr).Error
	return distrobutorGr, err
}

func (repository *RepositoryGrImpl) FindAllGrByCustIdSupId(dataFilter entity.GrLookupQueryFilter, custId, parentCustId string) ([]model.GrLookup, int64, int, error) {
	var grs []model.GrLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_no").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = inv.gr.sup_id AND sup.cust_id = ?", parentCustId)
	queryCount.Where("inv.gr.cust_id=?", custId)

	query := repository.
		Select("inv.gr.*, sup.sup_code, sup_name").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = inv.gr.sup_id AND sup.cust_id = ?", parentCustId)

	query.Where("inv.gr.cust_id = ?", custId)

	query.Where("(inv.gr.invoice_no = ? OR inv.gr.invoice_no IS NULL)", "")
	queryCount.Where("(inv.gr.invoice_no = ? OR inv.gr.invoice_no IS NULL)", "")

	if len(dataFilter.SupId) > 0 {
		query.Where("inv.gr.sup_id in ?", dataFilter.SupId)
		queryCount.Where("inv.gr.sup_id in ?", dataFilter.SupId)
	}

	if dataFilter.GrNo != "" {
		query.Where(`
			inv.gr.gr_no ILIKE '%` + dataFilter.GrNo + `%' 	
		`)
		queryCount.Where(`
			inv.gr.gr_no ILIKE '%` + dataFilter.GrNo + `%'
			
		`)
	}

	if dataFilter.Query != "" {
		query.Where(`(
			inv.gr.gr_no ILIKE '%` + dataFilter.Query + `%' OR
			sup.sup_name ILIKE '%` + dataFilter.Query + `%'
		)`)
		queryCount.Where(`(
			inv.gr.gr_no ILIKE '%` + dataFilter.Query + `%' OR
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
		query.Order("inv.gr.gr_no DESC")
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
	return grs, total, lastPage, nil
}

func (repository *RepositoryGrImpl) FindAllGrBranchByCustIdSupId(dataFilter entity.GrLookupQueryFilter, custId, parentCustId string) ([]model.GrBranchLookup, int64, int, error) {
	var grs []model.GrBranchLookup
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gr_branch_no").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = inv.gr_branch.sup_id AND sup.cust_id = ?", parentCustId)
	queryCount.Where("inv.gr_branch.cust_id=?", custId)

	query := repository.
		Select("inv.gr_branch.*, sup.sup_code, sup_name").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = inv.gr_branch.sup_id AND sup.cust_id = ?", parentCustId)

	query.Where("inv.gr_branch.cust_id = ?", custId)

	query.Where("(inv.gr_branch.invoice_no = ? OR inv.gr_branch.invoice_no IS NULL)", "")
	queryCount.Where("(inv.gr_branch.invoice_no = ? OR inv.gr_branch.invoice_no IS NULL)", "")

	if len(dataFilter.SupId) > 0 {
		query.Where("inv.gr_branch.sup_id in ?", dataFilter.SupId)
		queryCount.Where("inv.gr_branch.sup_id in ?", dataFilter.SupId)
	}

	if dataFilter.GrNo != "" {
		query.Where(`
			inv.gr_branch.gr_branch_no ILIKE '%` + dataFilter.GrNo + `%' 	
		`)
		queryCount.Where(`
			inv.gr_branch.gr_branch_no ILIKE '%` + dataFilter.GrNo + `%'
			
		`)
	}

	if dataFilter.Query != "" {
		query.Where(`(
			inv.gr_branch.gr_branch_no ILIKE '%` + dataFilter.Query + `%' OR
			sup.sup_name ILIKE '%` + dataFilter.Query + `%'
		)`)
		queryCount.Where(`(
			inv.gr_branch.gr_branch_no ILIKE '%` + dataFilter.Query + `%' OR
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
		query.Order("inv.gr_branch.gr_branch_no DESC")
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
	return grs, total, lastPage, nil
}

func (repository *RepositoryGrImpl) CheckReportInProgress(ctx context.Context, prefix string) (bool, error) {
	var count int64
	err := repository.model(ctx).
		Table("report.list").
		Where("report_name LIKE ?", prefix+"%").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *RepositoryGrImpl) FindGrDetailForDownload(grNo string, custId, parentCustId string) (gr model.GrList, grDetails []model.GrDetList, err error) {
	// Get GR header
	err = repository.
		Select("gr.*, us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = gr.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gr.closed_by").
		Joins("left join mst.m_supplier sup on sup.sup_id = gr.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = gr.wh_id AND wh.cust_id = ?", custId).
		Where("gr.gr_no = ? AND gr.cust_id=?", grNo, custId).
		Take(&gr).Error
	if err != nil {
		return gr, grDetails, err
	}

	// Get GR details with product info
	err = repository.
		Select(`gr_det.*,
			pd.pro_code, pd.pro_name, pd.unit_id1, pd.unit_id2, pd.unit_id3, pd.conv_unit2, pd.conv_unit3`).
		Joins("LEFT JOIN mst.m_product pd ON pd.pro_id = gr_det.pro_id").
		Where("gr_det.gr_no = ? AND gr_det.cust_id = ?", grNo, custId).
		Order("gr_det.seq_no ASC").
		Find(&grDetails).Error

	return gr, grDetails, err
}
