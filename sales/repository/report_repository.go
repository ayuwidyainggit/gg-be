package repository

import (
	"context"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	RepositoryReportImpl struct {
		*gorm.DB
	}
)
type ReportRepository interface {
	SecondarySales(dataFilter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReport, int64, int, error)
	StoreReportList(c context.Context, data *model.ReportList) error
	FindAllByCustID(dataFilter entity.ReportQueryFilter) (data []model.ReportList, total int64, lastPage int, err error)
	CountSecondarySalesReportByDate(dataFilter entity.SecondarySalesReportQueryFilter) int64
	UpdateReportByReportID(c context.Context, reportID string, data *model.ReportList) error
	ActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) ([]model.SalesActivityReportRow, int64, int, error)
	ActivitySalesReportList(dataFilter entity.ActivityReportQueryFilterList) ([]model.SalesActivityReportRow, int64, int, error)
	CountReportByDateAndReportName(custID string, exportDate, reportName string) int64
	SecondarySalesUnionPagination(dataFilter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, int64, int, error)
	SecondarySalesUnion(dataFilter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error)
	GetReportByReportID(reportID string) (data model.ReportList, err error)
	GetReportSecondarySalesReportOrder(custID string, date time.Time, limit, offset int) (data []model.SecondarySalesReportUnionReport, err error)
	GetReportSecondarySalesReportReturn(custID string, date time.Time, limit, offset int) (data []model.SecondarySalesReportUnionReturn, err error)
	SaveProductCategoriesDim(c context.Context, productCategories []model.DimProductCategory) (err error)
	SaveProductDim(c context.Context, products []model.DimProduct) (err error)
	SaveOutletsDim(c context.Context, outlets []model.DimOutlet) (err error)
	SaveSalemanDim(c context.Context, salesmans []model.DimSalesman) (err error)
	GetOrCreateBatchDimDate(ctx context.Context, dates []time.Time) (map[string]int64, error)
	SaveOrderfact(c context.Context, orders []model.FactOrder) (err error)
	SaveReturnfact(c context.Context, returns []model.FactReturn) (err error)
	ListCustIDReportSecondarySalesReportOrder(date time.Time) (data []model.SecondarySalesReportOrderCustID, err error)
	ListCustIDReportSecondarySalesReportReturn(date time.Time) (data []model.SecondarySalesReportReturnCustID, err error)
	ExistsCustomerInParentScope(custID string, parentCustID string) (bool, error)
	SecondarySalesReportSumReportByMonth(custIDs []string, filter entity.SecondarySalesReportDashboardSumPayload, year int) (data model.SumReportByMonthModel, err error)
	SecondarySalesReportGroupOutlet(custIDs []string, month int, year int) (data []model.SecondarySalesReportGroup, err error)
	SecondarySalesReportGroupSalesman(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error)
	SecondarySalesReportProductCategory(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error)
	SecondarySalesReportProduct(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error)
	SecondarySalesReportReturnSumReportByMonth(custIDs []string, month int, year int) (data model.SumReportReturnByMonthModel, err error)
	SecondarySalesReportTrendSales(custIDs []string, year int) (results []model.TrendSalesSecondarySalesModel, err error)
	SalesmanActivityReportSumByMonth(custIDs []string, month int, year int) (data model.SalesmanActivitySumByMonthModel, err error)
	SalesmanActivityReportTrendSales(custIDs []string, year int) (results []model.ActivityReportTrendSalesModel, err error)
	ActivityReportGeotag(parentCustID string, custIDs []string, year int, empID *int) (results []model.ActivityReportGeotagRow, err error)
	ActivitySalesmanReportGroupSalesman(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error)
	ActivitySalesmanReturnReportGroupSalesman(custIDs []string, month int, year int) (results []model.ReturnReportGroup, err error)
	ActivitySalesReportSalesmanList(dataFilter entity.ActivityReportSalesmanListQueryFilter) ([]model.SalesActivityReportSalesmanList, error)
	// DownloadSalesOrder report methods
	CountDownloadSalesOrderInProgress(custID string) (int64, error)
	CountDownloadSalesOrderByDate(custID, exportDate string) int64
}

func NewReportRepo(db *gorm.DB) *RepositoryReportImpl {
	return &RepositoryReportImpl{db}
}

func (repository *RepositoryReportImpl) ExistsCustomerInParentScope(custID string, parentCustID string) (bool, error) {
	var count int64
	if err := repository.Model(&model.SmcMCustomer{}).
		Where("cust_id = ? AND parent_cust_id = ? AND is_del = false AND is_active = true", custID, parentCustID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func secondarySalesProductSelectWithFallback(
	parentProductAlias, customerProductAlias, parentSupplierAlias, customerSupplierAlias string,
	fallbackUnit1Expr, fallbackUnit2Expr, fallbackUnit3Expr, fallbackConv2Expr, fallbackConv3Expr string,
) string {
	return fmt.Sprintf(`
		COALESCE(%s.sup_code, %s.sup_code) AS sup_code,
		COALESCE(%s.sup_name, %s.sup_name) AS sup_name,
		COALESCE(%s.pro_code, %s.pro_code) AS pro_code,
		COALESCE(%s.pro_name, %s.pro_name) AS pro_name,
		COALESCE(%s.unit_id1, %s.unit_id1, %s) AS unit_id1,
		COALESCE(%s.unit_id2, %s.unit_id2, %s) AS unit_id2,
		COALESCE(%s.unit_id3, %s.unit_id3, %s) AS unit_id3,
		COALESCE(%s.conv_unit2, %s.conv_unit2, %s, 0) AS conv_unit2,
		COALESCE(%s.conv_unit3, %s.conv_unit3, %s, 0) AS conv_unit3,
	`,
		customerSupplierAlias, parentSupplierAlias,
		customerSupplierAlias, parentSupplierAlias,
		customerProductAlias, parentProductAlias,
		customerProductAlias, parentProductAlias,
		customerProductAlias, parentProductAlias, fallbackUnit1Expr,
		customerProductAlias, parentProductAlias, fallbackUnit2Expr,
		customerProductAlias, parentProductAlias, fallbackUnit3Expr,
		customerProductAlias, parentProductAlias, fallbackConv2Expr,
		customerProductAlias, parentProductAlias, fallbackConv3Expr,
	)
}

func secondarySalesProductSelect(parentProductAlias, customerProductAlias, parentSupplierAlias, customerSupplierAlias, detailAlias string) string {
	return secondarySalesProductSelectWithFallback(
		parentProductAlias,
		customerProductAlias,
		parentSupplierAlias,
		customerSupplierAlias,
		fmt.Sprintf("%s.unit_id1", detailAlias),
		fmt.Sprintf("%s.unit_id2", detailAlias),
		fmt.Sprintf("%s.unit_id3", detailAlias),
		fmt.Sprintf("%s.conv_unit2", detailAlias),
		fmt.Sprintf("%s.conv_unit3", detailAlias),
	)
}

func secondarySalesProductJoins(detailAlias, rowCustExpr, parentCustExpr string) string {
	return fmt.Sprintf(`
		LEFT JOIN mst.m_product cp ON cp.pro_id = %s AND cp.cust_id = %s
		LEFT JOIN LATERAL (
			SELECT pp.*
			FROM mst.m_product pp
			WHERE pp.cust_id = %s
				AND (pp.pro_id = %s OR (cp.pro_code IS NOT NULL AND pp.pro_code = cp.pro_code))
			ORDER BY CASE WHEN pp.pro_id = %s THEN 0 ELSE 1 END, pp.pro_id ASC
			LIMIT 1
		) pp ON TRUE
		LEFT JOIN mst.m_supplier psup ON psup.sup_id = pp.sup_id
		LEFT JOIN mst.m_supplier csup ON csup.sup_id = cp.sup_id
	`, detailAlias, rowCustExpr, parentCustExpr, detailAlias, detailAlias)
}

func buildSecondarySalesUnionQuery(dataFilter entity.SecondarySalesReportQueryFilter, withPagination bool) (string, []interface{}, int) {
	limit := dataFilter.Limit
	if limit == 0 {
		limit = 10
	}

	offset := 0
	if dataFilter.Page > 0 {
		offset = (dataFilter.Page - 1) * limit
	}

	custIDs := dataFilter.CustIDs
	if len(custIDs) == 0 && dataFilter.CustID != "" {
		custIDs = []string{dataFilter.CustID}
	}

	whereOrder := `od.cust_id IN ? AND o.data_status IN (6,7)`
	whereReturn := `rd.cust_id IN ? AND o.data_status IN (6,7)`
	outerWhere := ""

	paramsOrder := []interface{}{custIDs}
	paramsReturn := []interface{}{custIDs}
	paramsOuter := []interface{}{}

	if dataFilter.From != nil && dataFilter.To != nil {
		whereOrder += " AND o.invoice_date BETWEEN ? AND ?"
		whereReturn += " AND o.invoice_date BETWEEN ? AND ?"

		fromT := str.UnixTimestampToUtcTime(*dataFilter.From)
		toT := str.UnixTimestampToUtcTime(*dataFilter.To)

		paramsOrder = append(paramsOrder, fromT, toT)
		paramsReturn = append(paramsReturn, fromT, toT)
	}

	if len(dataFilter.DistributorIDs) > 0 {
		outerWhere += " AND md.distributor_id IN ?"
		paramsOuter = append(paramsOuter, dataFilter.DistributorIDs)
	}

	if len(dataFilter.SalesmanIDs) > 0 {
		whereOrder += " AND o.salesman_id IN ?"
		whereReturn += " AND r.salesman_id IN ?"

		paramsOrder = append(paramsOrder, dataFilter.SalesmanIDs)
		paramsReturn = append(paramsReturn, dataFilter.SalesmanIDs)
	}

	if len(dataFilter.OutletIDs) > 0 {
		whereOrder += " AND o.outlet_id IN ?"
		whereReturn += " AND r.outlet_id IN ?"

		paramsOrder = append(paramsOrder, dataFilter.OutletIDs)
		paramsReturn = append(paramsReturn, dataFilter.OutletIDs)
	}

	if len(dataFilter.ProIDs) > 0 {
		whereOrder += " AND od.pro_id IN ?"
		whereReturn += " AND rd.product_id IN ?"

		paramsOrder = append(paramsOrder, dataFilter.ProIDs)
		paramsReturn = append(paramsReturn, dataFilter.ProIDs)
	}

	whereOrder += " AND o.invoice_no IS NOT NULL"
	whereReturn += " AND r.data_status = 6"

	productSelect := secondarySalesProductSelectWithFallback(
		"pp", "cp", "psup", "csup",
		"t.t_unit_id1", "t.t_unit_id2", "t.t_unit_id3",
		"t.t_conv_unit2", "t.t_conv_unit3",
	)
	productJoins := secondarySalesProductJoins("t.product_id", "t.cust_id", "?")

	sql := fmt.Sprintf(`
WITH trans AS (
	SELECT
		od.cust_id,
		CASE WHEN od.item_type = 1 THEN 'ORDER' ELSE 'PROMO ORDER' END AS trx_type,
		o.invoice_no,
		o.invoice_date,
		o.invoice_no AS document_no,
		o.ro_date AS document_date,
		o.outlet_id,
		o.salesman_id,
		od.pro_id AS product_id,
		od.unit_id1 AS t_unit_id1,
		od.unit_id2 AS t_unit_id2,
		od.unit_id3 AS t_unit_id3,
		od.conv_unit2 AS t_conv_unit2,
		od.conv_unit3 AS t_conv_unit3,
		COALESCE(od.qty1_final, 0) AS qty1,
		COALESCE(od.qty2_final, 0) AS qty2,
		COALESCE(od.qty3_final, 0) AS qty3,
		((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
		 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
		 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0))) AS gross_sales,
		COALESCE(od.promo_value_final, 0) AS special_discount,
		COALESCE(od.disc_value_final, 0) AS discount,
		CASE WHEN od.item_type = 1 THEN
			((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
			 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
			 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0)))
			- COALESCE(od.promo_value_final, 0)
			- COALESCE(od.disc_value_final, 0)
		ELSE 0 END AS net_sales_exc_ppn,
		CASE WHEN od.item_type = 1 THEN od.vat_value_final ELSE 0 END AS ppn,
		CASE WHEN od.item_type = 1 THEN
			((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
			 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
			 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0)))
			- COALESCE(od.promo_value_final, 0)
			- COALESCE(od.disc_value_final, 0)
			+ COALESCE(od.vat_value_final, 0)
		ELSE 0 END AS net_sales_inc_ppn,
		od.sell_price1,
		od.sell_price2,
		od.sell_price3
	FROM sls.order_detail od
	JOIN sls."order" o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id
	WHERE %s

	UNION ALL

	SELECT
		rd.cust_id,
		'RETURN' AS trx_type,
		o.invoice_no,
		o.invoice_date,
		r.return_no AS document_no,
		r.return_date AS document_date,
		r.outlet_id,
		r.salesman_id,
		rd.product_id,
		rd.unit_id1 AS t_unit_id1,
		rd.unit_id2 AS t_unit_id2,
		rd.unit_id3 AS t_unit_id3,
		rd.conv_unit2 AS t_conv_unit2,
		rd.conv_unit3 AS t_conv_unit3,
		CASE WHEN rd.qty1 > 0 THEN rd.qty1 * -1 ELSE rd.qty1 END AS qty1,
		CASE WHEN rd.qty2 > 0 THEN rd.qty2 * -1 ELSE rd.qty2 END AS qty2,
		CASE WHEN rd.qty3 > 0 THEN rd.qty3 * -1 ELSE rd.qty3 END AS qty3,
		((COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)) +
		 (COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)) +
		 (COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0))) * -1 AS gross_sales,
		COALESCE(rd.promo_value, 0) AS special_discount,
		COALESCE(rd.disc_value, 0) AS discount,
		(COALESCE(rd.total, 0) - COALESCE(rd.vat_value, 0)) * -1 AS net_sales_exc_ppn,
		COALESCE(rd.vat_value, 0) AS ppn,
		COALESCE(rd.total, 0) * -1 AS net_sales_inc_ppn,
		rd.sell_price1,
		rd.sell_price2,
		rd.sell_price3
	FROM sls.return_det rd
	JOIN sls."return" r ON rd.return_no = r.return_no AND rd.cust_id = r.cust_id
	JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
	WHERE %s
)
SELECT
	md.distributor_code,
	md.distributor_name,
	t.trx_type,
	t.invoice_no,
	t.invoice_date,
	t.document_no,
	t.document_date,
	t.outlet_id,
	mo.outlet_code,
	mo.outlet_name,
	t.salesman_id,
	e.emp_code,
	e.emp_name,
	t.product_id,
	%s
	t.qty1,
	t.qty2,
	t.qty3,
	t.gross_sales,
	t.special_discount,
	t.discount,
	t.net_sales_exc_ppn,
	t.ppn,
	t.net_sales_inc_ppn,
	t.sell_price1,
	t.sell_price2,
	t.sell_price3
FROM trans t
LEFT JOIN mst.m_outlet mo ON mo.outlet_id = t.outlet_id AND mo.cust_id = t.cust_id
LEFT JOIN mst.m_salesman s ON s.emp_id = t.salesman_id AND s.cust_id = t.cust_id AND s.is_del = FALSE
LEFT JOIN mst.m_employee e ON e.emp_id = s.emp_id AND e.cust_id = s.cust_id
%s
LEFT JOIN smc.m_customer mc ON mc.cust_id = t.cust_id
LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.cust_id = mc.cust_id
WHERE 1=1%s
ORDER BY t.invoice_no ASC, t.trx_type, t.document_no ASC`, whereOrder, whereReturn, productSelect, productJoins, outerWhere)

	allParams := []interface{}{}
	allParams = append(allParams, paramsOrder...)
	allParams = append(allParams, paramsReturn...)
	allParams = append(allParams, dataFilter.ParentCustID)
	allParams = append(allParams, paramsOuter...)

	if withPagination {
		sql += " LIMIT ? OFFSET ?"
		allParams = append(allParams, limit, offset)
	}

	return sql, allParams, limit
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryReportImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryReportImpl) SecondarySales(dataFilter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReport, int64, int, error) {
	var invoice []model.SecondarySalesReport
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	// queryCount := repository.Select("o.invoice_no")
	query := repository.Select(`
			md.distributor_code, md.distributor_name,
			o.invoice_no, o.invoice_date, 
			o.outlet_id, mo.outlet_code, mo.outlet_name, 
			o.salesman_id, e.emp_code, emp_name,
			sls.order_detail.pro_id, 
			sup.sup_code, sup.sup_name,
			p.pro_code, p.pro_name,
			p.unit_id1, p.unit_id2, p.unit_id3,
			p.conv_unit2, p.conv_unit3,
			sls.order_detail.qty1_final, sls.order_detail.qty2_final, sls.order_detail.qty3_final,
			((COALESCE(sls.order_detail.qty1_final,0)*COALESCE(sls.order_detail.sell_price1,0)) + (COALESCE(sls.order_detail.qty2_final,0)*COALESCE(sls.order_detail.sell_price2,0)) + (COALESCE(sls.order_detail.qty3_final,0)*COALESCE(sls.order_detail.sell_price3,0))) AS gross_sales,
			COALESCE(sls.order_detail.promo_value_final,0) AS special_discount,
			COALESCE(sls.order_detail.disc_value_final,0) AS discount,
			((COALESCE(sls.order_detail.qty1_final,0)*COALESCE(sls.order_detail.sell_price1,0)) + (COALESCE(sls.order_detail.qty2_final,0)*COALESCE(sls.order_detail.sell_price2,0)) + (COALESCE(sls.order_detail.qty3_final,0)*COALESCE(sls.order_detail.sell_price3,0))) - COALESCE(sls.order_detail.promo_value_final,0) - COALESCE(sls.order_detail.disc_value_final,0) AS net_sales_exc_ppn,
			COALESCE(sls.order_detail.vat_value_final,0) AS ppn,
			((COALESCE(sls.order_detail.qty1_final,0)*COALESCE(sls.order_detail.sell_price1,0)) + (COALESCE(sls.order_detail.qty2_final,0)*COALESCE(sls.order_detail.sell_price2,0)) + (COALESCE(sls.order_detail.qty3_final,0)*COALESCE(sls.order_detail.sell_price3,0))) - COALESCE(sls.order_detail.promo_value_final,0) - COALESCE(sls.order_detail.disc_value_final,0) + COALESCE(sls.order_detail.vat_value_final,0) AS net_sales_inc_ppn,
			sum(rd.qty1) AS qty1_return,
			sum(rd.qty2) AS qty2_return,
			sum(rd.qty3) AS qty3_return,
			sls.order_detail.sell_price1,
			sls.order_detail.sell_price2,
			sls.order_detail.sell_price3
		`).
		Joins("LEFT JOIN sls.order o ON o.ro_no = sls.order_detail.ro_no AND sls.order_detail.cust_id = o.cust_id").
		Joins("LEFT JOIN mst.m_outlet mo ON mo.outlet_id = o.outlet_id AND mo.cust_id = o.cust_id").
		Joins("LEFT JOIN mst.m_salesman s ON s.emp_id = o.salesman_id AND s.cust_id = o.cust_id AND s.is_del = false").
		Joins("LEFT JOIN mst.m_employee e ON e.emp_id = s.emp_id AND e.cust_id = s.cust_id").
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = sls.order_detail.pro_id AND p.cust_id = ?", dataFilter.ParentCustID).
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = p.sup_id AND sup.cust_id = ?", dataFilter.ParentCustID).
		Joins("LEFT JOIN sls.return r ON r.invoice_no = o.invoice_no AND r.data_status = 6 AND r.cust_id = o.cust_id").
		Joins("LEFT JOIN sls.return_det rd ON rd.return_no = r.return_no AND sls.order_detail.pro_id = rd.product_id AND rd.cust_id = r.cust_id").
		Joins("LEFT JOIN smc.m_customer mc ON mc.cust_id = o.cust_id").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.cust_id = mc.cust_id")

	groupBy := `md.distributor_code, md.distributor_name, 
		o.invoice_no,
		o.invoice_date,
		o.outlet_id,
		mo.outlet_code,
		mo.outlet_name,
		o.salesman_id,
		e.emp_code,
		e.emp_name,
		sls.order_detail.pro_id,
		sup.sup_code, 
		sup.sup_name,
		p.pro_code,
		p.pro_name,
		p.unit_id1,
		p.unit_id2,
		p.unit_id3,
		p.conv_unit2,
		p.conv_unit3,
		sls.order_detail.qty1_final,
		sls.order_detail.qty2_final,
		sls.order_detail.qty3_final,
		sls.order_detail.sell_price1,
		sls.order_detail.sell_price2,
		sls.order_detail.sell_price3,
		sls.order_detail.promo_value_final,
		sls.order_detail.disc_value_final,
		sls.order_detail.vat_value_final`

	custIDs := dataFilter.CustIDs
	if len(custIDs) == 0 && dataFilter.CustID != "" {
		custIDs = []string{dataFilter.CustID}
	}
	query.Where("sls.order_detail.cust_id IN ? AND o.data_status IN (6,7) AND sls.order_detail.item_type = 1", custIDs).Group(groupBy)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("o.invoice_date BETWEEN ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.OutletIDs) > 0 {
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletIDs)
	}

	if len(dataFilter.ProIDs) > 0 {
		query.Where("sls.order.pro_id in ?", dataFilter.ProIDs)
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
		query.Order("o.invoice_no ASC")
	}

	err := query.Find(&invoice).Error
	if err != nil {
		return invoice, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return invoice, total, lastPage, nil
}

func (repository *RepositoryReportImpl) SecondarySalesUnionPagination(
	dataFilter entity.SecondarySalesReportQueryFilter,
) ([]model.SecondarySalesReportUnion, int64, int, error) {
	var results []model.SecondarySalesReportUnion
	var total int64
	sql, allParams, limit := buildSecondarySalesUnionQuery(dataFilter, true)
	err := repository.Raw(sql, allParams...).Scan(&results).Error
	if err != nil {
		return results, 0, 0, err
	}

	countSQL, countParams, _ := buildSecondarySalesUnionQuery(dataFilter, false)
	countSQL = fmt.Sprintf(`SELECT COUNT(*) FROM (%s) AS combined_count`, countSQL)
	err = repository.Raw(countSQL, countParams...).Scan(&total).Error
	if err != nil {
		return results, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return results, total, lastPage, nil
}

func (repository *RepositoryReportImpl) SecondarySalesUnion(dataFilter entity.SecondarySalesReportQueryFilter) ([]model.SecondarySalesReportUnion, error) {
	var results []model.SecondarySalesReportUnion
	sql, allParams, _ := buildSecondarySalesUnionQuery(dataFilter, false)
	err := repository.Raw(sql, allParams...).Scan(&results).Error
	if err != nil {
		return results, err
	}

	return results, nil
}

func (repository *RepositoryReportImpl) StoreReportList(c context.Context, data *model.ReportList) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryReportImpl) FindAllByCustID(dataFilter entity.ReportQueryFilter) ([]model.ReportList, int64, int, error) {
	var (
		reports []model.ReportList
		total   int64
	)
	limit := 10
	if dataFilter.Limit != 0 {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("report_id")
	query := repository.Select(
		`report.list.*`) // .

	queryCount.Where("report.list.cust_id=?", dataFilter.CustId)
	query.Where("report.list.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where(`report.list.created_at BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To),
		)
		queryCount.Where(`report.list.created_at BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To),
		)
	}

	if dataFilter.Query != "" {
		queryCount.Where("report.list.report_id ILIKE ? OR report.list.report_name ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query.Where("report.list.report_id ILIKE ? OR report.list.report_name ILIKE ?", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if dataFilter.FileStatus != nil {
		// log.Debugf("FileStatus: %v", dataFilter.FileStatus)
		if len(dataFilter.FileStatus) > 0 {
			if dataFilter.FileStatus[0] != 0 {
				queryCount.Where("report.list.file_status IN ?", dataFilter.FileStatus)
				query.Where("report.list.file_status IN ?", dataFilter.FileStatus)
			}
		}
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
		query.Order("created_at DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&reports).Error
	if err != nil {
		return reports, total, 0, err
	}
	err = queryCount.Model(&reports).Count(&total).Error
	if err != nil {
		return reports, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return reports, total, lastPage, nil
}

func (repository *RepositoryReportImpl) CountSecondarySalesReportByDate(dataFilter entity.SecondarySalesReportQueryFilter) int64 {
	var (
		report model.ReportList
		total  int64
	)

	queryCount := repository.Select("report_id")
	reportNamePrefix := entity.REPORT_NAME_SECONDARY_SALES + "-" + dataFilter.ExportDate + "-%"
	queryCount.Where("report.list.cust_id = ? AND report.list.report_name LIKE ?",
		dataFilter.CustID, reportNamePrefix)

	err := queryCount.Model(&report).Count(&total).Error
	if err != nil {
		return total
	}

	total += 1 // Add 1 to the total count to account for the current report being created
	return total
}

func (repository *RepositoryReportImpl) UpdateReportByReportID(c context.Context, reportID string, data *model.ReportList) error {
	result := repository.model(c).Model(data).Where("report_id=?", reportID).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (repository *RepositoryReportImpl) ActivitySalesReportSalesmanList(dataFilter entity.ActivityReportSalesmanListQueryFilter) ([]model.SalesActivityReportSalesmanList, error) {
	var results []model.SalesActivityReportSalesmanList

	likeName := "%"
	if dataFilter.SalesmanName != "" {
		likeName = "%" + strings.ToLower(dataFilter.SalesmanName) + "%"
	}

	query := `select * from(
    SELECT 
        o.salesman_id,
        emp.emp_code AS salesman_code,
        emp.emp_name AS salesman_name
    FROM sls."order" o
    LEFT JOIN mst.m_employee emp 
        ON emp.emp_id = o.salesman_id 
       AND emp.cust_id = ?
    WHERE 
        o.cust_id = ?
        AND o.ro_date BETWEEN ? AND ?
        AND LOWER(emp.emp_name) LIKE ?

    UNION

    SELECT 
        r.salesman_id,
        emp.emp_code AS salesman_code,
        emp.emp_name AS salesman_name
    FROM sls."return" r
    LEFT JOIN mst.m_employee emp 
        ON emp.emp_id = r.salesman_id 
       AND emp.cust_id = ?
    WHERE 
        r.cust_id = ?
        AND r.return_date BETWEEN ? AND ?
        AND LOWER(emp.emp_name) LIKE ?
	) A
	 group by salesman_id, salesman_code, salesman_name
    `

	err := repository.Raw(query,
		dataFilter.CustID, dataFilter.CustID, dataFilter.FromDate, dataFilter.ToDate, likeName, // bagian order
		dataFilter.CustID, dataFilter.CustID, dataFilter.FromDate, dataFilter.ToDate, likeName, // bagian return
	).Find(&results).Error
	if err != nil {
		return results, err
	}
	return results, nil
}

func activityReportAuthAndCustIDsFromExportFilter(dataFilter entity.ActivityReportQueryFilter) (string, []string) {
	authCustID := strings.TrimSpace(dataFilter.AuthCustID)
	if authCustID == "" {
		authCustID = strings.TrimSpace(dataFilter.CustID)
	}
	custIDs := dataFilter.CustIDs
	if len(custIDs) == 0 && authCustID != "" {
		custIDs = []string{authCustID}
	}
	return authCustID, custIDs
}

func activityReportAuthAndCustIDsFromListFilter(dataFilter entity.ActivityReportQueryFilterList) (string, []string) {
	authCustID := strings.TrimSpace(dataFilter.AuthCustID)
	if authCustID == "" {
		authCustID = strings.TrimSpace(dataFilter.CustID)
	}
	custIDs := dataFilter.CustIDs
	if len(custIDs) == 0 && authCustID != "" {
		custIDs = []string{authCustID}
	}
	return authCustID, custIDs
}

func (repository *RepositoryReportImpl) ActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) ([]model.SalesActivityReportRow, int64, int, error) {
	authCustID, custIDs := activityReportAuthAndCustIDsFromExportFilter(dataFilter)
	rows, err := repository.queryActivitySalesReportRows(
		authCustID,
		custIDs,
		dataFilter.ParentCustID,
		dataFilter.FromDate,
		dataFilter.ToDate,
		dataFilter.SalesmanIDs,
		dataFilter.DistributorCodes,
		0,
		0,
	)
	if err != nil {
		return nil, 0, 0, err
	}
	return rows, int64(len(rows)), 1, nil
}

func (repository *RepositoryReportImpl) ActivitySalesReportList(dataFilter entity.ActivityReportQueryFilterList) ([]model.SalesActivityReportRow, int64, int, error) {
	limit := dataFilter.Limit
	if limit == 0 {
		limit = 10
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	authCustID, custIDs := activityReportAuthAndCustIDsFromListFilter(dataFilter)

	total, err := repository.countActivitySalesReportRows(
		authCustID,
		custIDs,
		dataFilter.ParentCustID,
		dataFilter.FromDate,
		dataFilter.ToDate,
		dataFilter.SalesmanIDs,
		dataFilter.DistributorCodes,
	)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, err := repository.queryActivitySalesReportRows(
		authCustID,
		custIDs,
		dataFilter.ParentCustID,
		dataFilter.FromDate,
		dataFilter.ToDate,
		dataFilter.SalesmanIDs,
		dataFilter.DistributorCodes,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return rows, total, lastPage, nil
}

func (repository *RepositoryReportImpl) CountReportByDateAndReportName(custID string, exportDate, reportName string) int64 {
	var (
		report model.ReportList
		total  int64
	)

	queryCount := repository.Select("report_id")
	// Convert ExportDate (string) to time.Time range: 00:00:00 - 23:59:59
	startDate, _ := time.Parse("020106", exportDate)
	endDate := startDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	queryCount.Where("report.list.cust_id = ? AND report.list.created_at BETWEEN ? AND ? AND report.list.report_name ILIKE ?",
		custID, startDate, endDate, "%"+reportName+"%")

	err := queryCount.Model(&report).Count(&total).Error
	if err != nil {
		return total
	}

	total += 1 // Add 1 to the total count to account for the current report being created
	return total
}

func (repository *RepositoryReportImpl) GetReportByReportID(reportID string) (data model.ReportList, err error) {
	err = repository.
		Select("*").
		Where("report_id = ?", reportID).
		Take(&data).Error

	return data, err
}

func (repository *RepositoryReportImpl) ListCustIDReportSecondarySalesReportOrder(date time.Time) (data []model.SecondarySalesReportOrderCustID, err error) {
	err = repository.Select("cust_id").Where("invoice_date = ?", str.FormatTimeToDateString(date)).Group("cust_id").Find(&data).Error
	return
}

func (repository *RepositoryReportImpl) ListCustIDReportSecondarySalesReportReturn(date time.Time) (data []model.SecondarySalesReportReturnCustID, err error) {
	err = repository.Select("cust_id").Where("closed_at::date = ?", str.FormatTimeToDateString(date)).Group("cust_id").Find(&data).Error
	return
}

func buildReportSecondarySalesReportOrderQuery(custID string, date time.Time, limit, offset int) (string, []interface{}) {
	sql := fmt.Sprintf(`
		SELECT
			o.cust_id,
			o.ro_no,
			o.ro_date,
			md.distributor_code,
			md.distributor_name,
			CASE WHEN od.item_type = 1 THEN 'ORDER' ELSE 'PROMO ORDER' END AS trx_type,
			o.invoice_no,
			o.invoice_date,
			o.outlet_id,
			mo.outlet_code,
			mo.outlet_name,
			o.salesman_id,
			e.emp_code,
			e.emp_name,
			od.pro_id AS product_id,
			%s
			od.qty1_final AS qty1,
			od.qty2_final AS qty2,
			od.qty3_final AS qty3,
			od.qty_final AS qty,
			od.item_type,
			((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) + (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) + (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0))) AS gross_sales,
			COALESCE(od.promo_value_final, 0) AS special_discount,
			COALESCE(od.disc_value_final, 0) AS discount,
			CASE WHEN od.item_type = 1 THEN
				((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
				 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
				 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0)))
				- COALESCE(od.promo_value_final, 0)
				- COALESCE(od.disc_value_final, 0)
			ELSE 0 END AS net_sales_exc_ppn,
			CASE WHEN od.item_type = 1 THEN od.vat_value_final ELSE 0 END AS ppn,
			CASE WHEN od.item_type = 1 THEN
				((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
				 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
				 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0)))
				- COALESCE(od.promo_value_final, 0)
				- COALESCE(od.disc_value_final, 0)
				+ COALESCE(od.vat_value_final, 0)
			ELSE 0 END AS net_sales_inc_ppn,
			od.sell_price1,
			od.sell_price2,
			od.sell_price3,
			prdcat.pcat_id,
			prdcat.pcat_code,
			prdcat.pcat_name
		FROM sls.order_detail od
		LEFT JOIN sls.order o ON o.ro_no = od.ro_no AND od.cust_id = o.cust_id
		LEFT JOIN mst.m_outlet mo ON mo.outlet_id = o.outlet_id AND mo.cust_id = o.cust_id
		LEFT JOIN mst.m_salesman s ON s.emp_id = o.salesman_id AND s.cust_id = o.cust_id AND s.is_del = FALSE
		LEFT JOIN mst.m_employee e ON e.emp_id = s.emp_id AND e.cust_id = s.cust_id
		%s
		LEFT JOIN smc.m_customer mc ON mc.cust_id = o.cust_id
		LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.cust_id = mc.cust_id
		LEFT JOIN mst.m_product_cat prdcat ON prdcat.pcat_id = COALESCE(pp.pcat_id, cp.pcat_id)
		WHERE o.cust_id = ? AND o.invoice_date = ?
		ORDER BY o.ro_no
		LIMIT ? OFFSET ?
	`, secondarySalesProductSelect("pp", "cp", "psup", "csup", "od"), secondarySalesProductJoins("od.pro_id", "od.cust_id", "?"))

	params := []interface{}{
		custID,
		custID,
		str.FormatTimeToDateString(date),
		limit,
		offset,
	}

	return sql, params
}

func (repository *RepositoryReportImpl) GetReportSecondarySalesReportOrder(custID string, date time.Time, limit, offset int) (data []model.SecondarySalesReportUnionReport, err error) {
	var results []model.SecondarySalesReportUnionReport

	sql, params := buildReportSecondarySalesReportOrderQuery(custID, date, limit, offset)

	err = repository.Raw(sql, params...).Scan(&results).Error
	if err != nil {
		return results, err
	}

	return results, nil
}

func (repository *RepositoryReportImpl) GetReportSecondarySalesReportReturn(custID string, date time.Time, limit, offset int) (data []model.SecondarySalesReportUnionReturn, err error) {
	var results []model.SecondarySalesReportUnionReturn

	sql, params := buildReportSecondarySalesReportReturnQuery(custID, date, limit, offset)

	err = repository.Raw(sql, params...).Scan(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}

func buildReportSecondarySalesReportReturnQuery(custID string, date time.Time, limit, offset int) (string, []interface{}) {
	sql := fmt.Sprintf(`
		SELECT
		r.cust_id,
        md.distributor_code,
        md.distributor_name,
        'RETURN' AS trx_type,
        o.invoice_no,
        o.invoice_date,
		r.return_no,
		r.return_date,
        o.outlet_id,
        mo.outlet_code,
        mo.outlet_name,
        o.salesman_id,
        e.emp_code,
        e.emp_name,
        rd.product_id,
		%s
        CASE WHEN rd.qty1 > 0 THEN rd.qty1 ELSE rd.qty1 END AS qty1,
        CASE WHEN rd.qty2 > 0 THEN rd.qty2 ELSE rd.qty2 END AS qty2,
        CASE WHEN rd.qty3 > 0 THEN rd.qty3  ELSE rd.qty3 END AS qty3,
        ((COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)) + (COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)) + (COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0))) AS gross_sales,
        COALESCE(rd.promo_value, 0) AS special_discount,
        COALESCE(rd.disc_value, 0) AS discount,
        (COALESCE(rd.total, 0) - COALESCE(rd.vat_value, 0)) AS net_sales_exc_ppn,
        COALESCE(rd.vat_value, 0) AS ppn,
        COALESCE(rd.total, 0) AS net_sales_inc_ppn,
        rd.sell_price1,
        rd.sell_price2,
        rd.sell_price3,
		prdcat.pcat_id,
			prdcat.pcat_code,
			prdcat.pcat_name
    FROM sls.return_det rd
     LEFT JOIN sls."return" r ON rd.return_no = r.return_no AND rd.cust_id = r.cust_id
    LEFT JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
    LEFT JOIN mst.m_outlet mo ON mo.outlet_id = r.outlet_id AND mo.cust_id = r.cust_id
    LEFT JOIN mst.m_salesman s ON s.emp_id = r.salesman_id AND s.cust_id = r.cust_id AND s.is_del = FALSE
    LEFT JOIN mst.m_employee e ON e.emp_id = s.emp_id AND e.cust_id = s.cust_id
		%s
    LEFT JOIN smc.m_customer mc ON mc.cust_id = r.cust_id
    LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.cust_id = mc.cust_id
	LEFT JOIN mst.m_product_cat prdcat ON prdcat.pcat_id = COALESCE(pp.pcat_id, cp.pcat_id)
		WHERE r.cust_id = ? AND r.closed_at::date = ?
		ORDER BY rd.return_no
		LIMIT ? OFFSET ?
	`, secondarySalesProductSelect("pp", "cp", "psup", "csup", "rd"), secondarySalesProductJoins("rd.product_id", "rd.cust_id", "?"))

	params := []interface{}{
		custID,
		custID,
		str.FormatTimeToDateString(date),
		limit,
		offset,
	}

	return sql, params
}

func (repository *RepositoryReportImpl) SaveProductCategoriesDim(c context.Context, productCategories []model.DimProductCategory) (err error) {
	err = repository.model(c).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"code", "name"}),
		},
	).Create(&productCategories).Error
	if err != nil {
		log.Error("SaveProductCategoriesDim Model, error:", err.Error())
		return err
	}
	return
}

func (repository *RepositoryReportImpl) SaveProductDim(c context.Context, products []model.DimProduct) (err error) {
	err = repository.model(c).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"category_id", "code", "name", "unit_id1", "unit_id2", "unit_id3", "conv_unit2", "conv_unit3"}),
		},
	).Create(&products).Error
	if err != nil {
		log.Error("SaveProductDim Model, error:", err.Error())
		return err
	}
	return
}

func (repository *RepositoryReportImpl) SaveOutletsDim(c context.Context, outlets []model.DimOutlet) (err error) {
	err = repository.model(c).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"code", "name"}),
		},
	).Create(&outlets).Error
	if err != nil {
		log.Error("SaveOutletsDim Model, error:", err.Error())
		return err
	}
	return
}

func (repository *RepositoryReportImpl) SaveSalemanDim(c context.Context, salesmans []model.DimSalesman) (err error) {
	err = repository.model(c).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"code", "name"}),
		},
	).Create(&salesmans).Error
	if err != nil {
		log.Error("SaveSalemanDim Model, error:", err.Error())
		return err
	}
	return
}

func (repository *RepositoryReportImpl) SaveOrderfact(c context.Context, orders []model.FactOrder) (err error) {
	err = repository.model(c).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "cust_id"}, {Name: "ro_no"}, {Name: "pro_id"}, {Name: "item_type"}},
			DoUpdates: clause.AssignmentColumns([]string{"invoice_no",
				"date_id",
				"salesman_id",
				"outlet_id",
				"qty",
				"qty1",
				"qty2",
				"qty3",
				"gross_sale",
				"special_discount",
				"discount",
				"net_sales_exclude_ppn",
				"ppn",
				"net_sales_include_ppn",
				"sell_price1",
				"sell_price2",
				"sell_price3",
				"extracted_at"}),
		},
	).Create(&orders).Error
	if err != nil {
		log.Error("SaveSalemanDim Model, error:", err.Error())
		return err
	}

	return
}

func (repository *RepositoryReportImpl) SaveReturnfact(c context.Context, returns []model.FactReturn) (err error) {
	err = repository.model(c).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "cust_id"}, {Name: "return_no"}, {Name: "pro_id"}, {Name: "item_type"}},
			DoUpdates: clause.AssignmentColumns([]string{"cust_id", "invoice_no",
				"date_id",
				"salesman_id",
				"outlet_id",
				"qty",
				"qty1",
				"qty2",
				"qty3",
				"gross_sale",
				"special_discount",
				"discount",
				"net_sales_exclude_ppn",
				"ppn",
				"net_sales_include_ppn",
				"sell_price1",
				"sell_price2",
				"sell_price3",
				"extracted_at"}),
		},
	).Create(&returns).Error
	if err != nil {
		log.Error("SaveSalemanDim Model, error:", err.Error())
		return err
	}

	return
}

func (repository *RepositoryReportImpl) GetOrCreateBatchDimDate(ctx context.Context, dates []time.Time) (map[string]int64, error) {

	resultMap := make(map[string]int64)
	if len(dates) == 0 {
		return resultMap, nil
	}

	// 1️⃣ Buat map unik berdasarkan day-month-year
	uniqueDates := make(map[string]model.DimDate)
	for _, t := range dates {
		key := fmt.Sprintf("%d-%d-%d", t.Day(), int(t.Month()), t.Year())
		if _, exists := uniqueDates[key]; !exists {
			uniqueDates[key] = model.DimDate{
				Day:   t.Day(),
				Month: int(t.Month()),
				Year:  t.Year(),
			}
		}
	}

	// 2️⃣ Ambil data yang sudah ada
	var existing []model.DimDate
	err := repository.model(ctx).
		Where("(day, month, year) IN ?", getDateTuples(uniqueDates)).
		Find(&existing).Error
	if err != nil {
		return nil, err
	}

	// Isi map untuk yang sudah ada
	for _, d := range existing {
		key := fmt.Sprintf("%d-%d-%d", d.Day, d.Month, d.Year)
		resultMap[key] = *d.ID
		delete(uniqueDates, key) // hapus dari kandidat insert
	}

	// 3️⃣ Insert sisanya (yang belum ada)
	var newDates []model.DimDate
	for _, d := range uniqueDates {
		newDates = append(newDates, d)
	}

	if len(newDates) > 0 {
		if err := repository.model(ctx).CreateInBatches(&newDates, 1000).Error; err != nil {
			return nil, err
		}
		for _, d := range newDates {
			key := fmt.Sprintf("%d-%d-%d", d.Day, d.Month, d.Year)
			resultMap[key] = *d.ID
		}
	}

	return resultMap, nil
}

// Helper untuk bikin tuple WHERE (day, month, year)
func getDateTuples(unique map[string]model.DimDate) [][]int {
	tuples := make([][]int, 0, len(unique))
	for _, d := range unique {
		tuples = append(tuples, []int{d.Day, d.Month, d.Year})
	}
	return tuples
}
func buildSecondarySalesSummaryDateRange(filter entity.SecondarySalesReportDashboardSumPayload, year int) (time.Time, time.Time) {
	if filter.From != nil && filter.To != nil {
		return str.UnixTimestampToUtcTime(*filter.From), str.UnixTimestampToUtcTime(*filter.To)
	}

	dateFrom := time.Date(year, time.Month(filter.Month), 1, 0, 0, 0, 0, time.UTC)
	return dateFrom, dateFrom.AddDate(0, 1, 0)
}

func appendSecondarySalesSummaryOrderFilters(query *strings.Builder, params *[]interface{}, filter entity.SecondarySalesReportDashboardSumPayload) {
	if len(filter.OutletIDs) > 0 {
		query.WriteString(" AND o.outlet_id IN ?")
		*params = append(*params, filter.OutletIDs)
	}
	if len(filter.SalesmanIDs) > 0 {
		query.WriteString(" AND o.salesman_id IN ?")
		*params = append(*params, filter.SalesmanIDs)
	}
	if len(filter.ProIDs) > 0 {
		query.WriteString(" AND od.pro_id IN ?")
		*params = append(*params, filter.ProIDs)
	}
}

func appendSecondarySalesSummaryReturnFilters(query *strings.Builder, params *[]interface{}, filter entity.SecondarySalesReportDashboardSumPayload) {
	if len(filter.OutletIDs) > 0 {
		query.WriteString(" AND r.outlet_id IN ?")
		*params = append(*params, filter.OutletIDs)
	}
	if len(filter.SalesmanIDs) > 0 {
		query.WriteString(" AND r.salesman_id IN ?")
		*params = append(*params, filter.SalesmanIDs)
	}
	if len(filter.ProIDs) > 0 {
		query.WriteString(" AND rd.product_id IN ?")
		*params = append(*params, filter.ProIDs)
	}
}

func buildSecondarySalesReportSummarySQL(filter entity.SecondarySalesReportDashboardSumPayload) (string, []interface{}) {
	year := time.Now().Year()
	if filter.Year != nil {
		year = *filter.Year
	}
	dateFrom, dateTo := buildSecondarySalesSummaryDateRange(filter, year)

	params := []interface{}{filter.CustIDs, dateFrom, dateTo}
	var query strings.Builder
	query.WriteString(`
		WITH order_summary AS (
			SELECT
				COALESCE(SUM(
					(COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
					(COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
					(COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0))
				), 0) AS gross_sale,
				COALESCE(SUM(
					COALESCE(od.disc_value_final, 0) +
					COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) +
					COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) +
					COALESCE(od.promo_final5, 0)
				), 0) AS discount_promo,
				COALESCE(SUM(CASE WHEN od.item_type = 1 THEN COALESCE(od.vat_value_final, 0) ELSE 0 END), 0) AS ppn,
				COALESCE(SUM(CASE WHEN od.item_type = 1 THEN
					((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
					 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
					 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0))) -
					COALESCE(od.promo_value_final, 0) -
					COALESCE(od.disc_value_final, 0)
				ELSE 0 END), 0) AS net_sales_exc_ppn,
				COALESCE(SUM(CASE WHEN od.item_type = 1 THEN
					((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
					 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
					 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0))) -
					COALESCE(od.promo_value_final, 0) -
					COALESCE(od.disc_value_final, 0) +
					COALESCE(od.vat_value_final, 0)
				ELSE 0 END), 0) AS net_sales_inc_ppn,
				COUNT(DISTINCT o.salesman_id) AS total_salesman,
				COUNT(DISTINCT o.outlet_id) AS total_outlet,
				COUNT(DISTINCT od.pro_id) AS total_product,
				COALESCE(SUM(
					(COALESCE(od.qty3_final, 0) * COALESCE(od.conv_unit2, 1) * COALESCE(od.conv_unit3, 1)) +
					(COALESCE(od.qty2_final, 0) * COALESCE(od.conv_unit2, 1)) +
					COALESCE(od.qty1_final, 0)
				), 0)::bigint AS qty,
				MAX(o.updated_at) AS last_update
			FROM sls."order" o
			JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id
			WHERE o.cust_id IN ? AND o.data_status IN (6,7) AND o.invoice_date >= ? AND o.invoice_date < ?`)
	appendSecondarySalesSummaryOrderFilters(&query, &params, filter)
	query.WriteString(`
		), return_summary AS (
			SELECT
				COALESCE(SUM(
					(COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)) +
					(COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)) +
					(COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0))
				), 0) AS gross_sale,
				COALESCE(SUM(COALESCE(rd.disc_value, 0) + COALESCE(rd.promo_value, 0)), 0) AS discount_promo,
				COALESCE(SUM(COALESCE(rd.vat_value, 0)), 0) AS ppn,
				COALESCE(SUM(COALESCE(rd.total, 0) - COALESCE(rd.vat_value, 0)), 0) AS net_sales_exc_ppn,
				COALESCE(SUM(COALESCE(rd.total, 0)), 0) AS net_sales_inc_ppn,
				COALESCE(SUM(
					(COALESCE(rd.qty3, 0) * COALESCE(rd.conv_unit2, 1) * COALESCE(rd.conv_unit3, 1)) +
					(COALESCE(rd.qty2, 0) * COALESCE(rd.conv_unit2, 1)) +
					COALESCE(rd.qty1, 0)
				), 0)::bigint AS qty_return,
				MAX(r.updated_at) AS last_update
			FROM sls.return_det rd
			JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
			JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
			WHERE rd.cust_id IN ? AND o.data_status IN (6,7) AND o.invoice_date >= ? AND o.invoice_date < ?`)
	params = append(params, filter.CustIDs, dateFrom, dateTo)
	appendSecondarySalesSummaryReturnFilters(&query, &params, filter)
	query.WriteString(`
		)
		SELECT
			(os.gross_sale - rs.gross_sale) AS total_gross_sale,
			COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo,
			(os.ppn - rs.ppn) AS total_ppn,
			(os.net_sales_exc_ppn - rs.net_sales_exc_ppn) AS net_sales_exc_ppn,
			(os.net_sales_inc_ppn - rs.net_sales_inc_ppn) AS net_sales,
			os.total_salesman,
			os.total_outlet,
			os.total_product,
			COALESCE(os.qty, 0) - COALESCE(rs.qty_return, 0) AS qty,
			rs.qty_return AS qty_return,
			COALESCE(ROUND(((rs.net_sales_inc_ppn / NULLIF(os.net_sales_inc_ppn, 0)) * 100)::numeric, 2), 0) AS return_rate,
			rs.net_sales_inc_ppn AS net_sales_return,
			CASE
				WHEN os.last_update IS NULL THEN rs.last_update
				WHEN rs.last_update IS NULL THEN os.last_update
				ELSE GREATEST(os.last_update, rs.last_update)
			END AS last_update
		FROM order_summary os
		CROSS JOIN return_summary rs
	`)

	return query.String(), params
}

func (repository *RepositoryReportImpl) SecondarySalesReportSumReportByMonth(custIDs []string, filter entity.SecondarySalesReportDashboardSumPayload, year int) (data model.SumReportByMonthModel, err error) {
	filter.CustIDs = custIDs
	if filter.Year == nil {
		filter.Year = &year
	}

	query, params := buildSecondarySalesReportSummarySQL(filter)
	err = repository.Raw(query, params...).Take(&data).Error

	return
}

func (repository *RepositoryReportImpl) SecondarySalesReportReturnSumReportByMonth(custIDs []string, month int, year int) (data model.SumReportReturnByMonthModel, err error) {
	err = repository.Select(`
		COALESCE(SUM ( qty ), 0) AS qty,
		COALESCE(SUM( report.fact_returns.net_sales_exclude_ppn), 0) AS net_sales,
		max(extracted_at) AS last_update
		`).
		Joins("JOIN report.dim_dates dt on report.fact_returns.date_id = dt.id").
		Where("report.fact_returns.cust_id IN ? AND dt.month = ? AND dt.\"year\" = ?", custIDs, month, year).Take(&data).Error

	return
}

func buildSecondarySalesReportGroupQuery(groupBy string) string {
	orderNetSales := `
		(
			((COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price1, 0)) +
			 (COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price2, 0)) +
			 (COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price3, 0))) -
			COALESCE(od.promo_value_final, 0) -
			COALESCE(od.disc_value_final, 0) +
			COALESCE(od.vat_value_final, 0)
		) AS net_sales`
	returnNetSales := `
		(
			((COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)) +
			 (COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)) +
			 (COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0))) -
			COALESCE(rd.promo_value, 0) -
			COALESCE(rd.disc_value, 0) +
			COALESCE(rd.vat_value, 0)
		) * -1 AS net_sales`

	orderSelect := map[string]string{
		"outlet": fmt.Sprintf(`
			o.outlet_id AS id,
			COALESCE(mo.outlet_code, '') AS code,
			COALESCE(mo.outlet_name, '') AS name,
			%s`, orderNetSales),
		"salesman": fmt.Sprintf(`
			o.salesman_id AS id,
			COALESCE(e.emp_code, '') AS code,
			COALESCE(e.emp_name, '') AS name,
			%s`, orderNetSales),
		"product_category": fmt.Sprintf(`
			COALESCE(NULLIF(mpc.pcat_id, 0), 0) AS id,
			COALESCE(NULLIF(mpc.pcat_code, ''), '') AS code,
			COALESCE(NULLIF(mpc.pcat_name, ''), '') AS name,
			%s`, orderNetSales),
		"product": fmt.Sprintf(`
			COALESCE(NULLIF(mp.pro_id, 0), od.pro_id) AS id,
			COALESCE(NULLIF(mp.pro_code, ''), '') AS code,
			COALESCE(NULLIF(mp.pro_name, ''), '') AS name,
			%s`, orderNetSales),
	}
	orderJoin := map[string]string{
		"outlet": `
			LEFT JOIN mst.m_outlet mo ON mo.outlet_id = o.outlet_id AND mo.cust_id = o.cust_id`,
		"salesman": `
			LEFT JOIN mst.m_salesman s ON s.emp_id = o.salesman_id AND s.cust_id = o.cust_id AND s.is_del = FALSE
			LEFT JOIN mst.m_employee e ON e.emp_id = s.emp_id AND e.cust_id = s.cust_id`,
		"product_category": `
			LEFT JOIN mst.m_product mp ON mp.pro_id = od.pro_id AND mp.cust_id = od.cust_id
			LEFT JOIN mst.m_product_cat mpc ON mpc.pcat_id = mp.pcat_id`,
		"product": `
			LEFT JOIN mst.m_product mp ON mp.pro_id = od.pro_id AND mp.cust_id = od.cust_id`,
	}
	returnSelect := map[string]string{
		"outlet": fmt.Sprintf(`
			r.outlet_id AS id,
			COALESCE(mo.outlet_code, '') AS code,
			COALESCE(mo.outlet_name, '') AS name,
			%s`, returnNetSales),
		"salesman": fmt.Sprintf(`
			r.salesman_id AS id,
			COALESCE(e.emp_code, '') AS code,
			COALESCE(e.emp_name, '') AS name,
			%s`, returnNetSales),
		"product_category": fmt.Sprintf(`
			COALESCE(NULLIF(mpc.pcat_id, 0), 0) AS id,
			COALESCE(NULLIF(mpc.pcat_code, ''), '') AS code,
			COALESCE(NULLIF(mpc.pcat_name, ''), '') AS name,
			%s`, returnNetSales),
		"product": fmt.Sprintf(`
			COALESCE(NULLIF(mp.pro_id, 0), rd.product_id) AS id,
			COALESCE(NULLIF(mp.pro_code, ''), '') AS code,
			COALESCE(NULLIF(mp.pro_name, ''), '') AS name,
			%s`, returnNetSales),
	}
	returnJoin := map[string]string{
		"outlet": `
			LEFT JOIN mst.m_outlet mo ON mo.outlet_id = r.outlet_id AND mo.cust_id = rd.cust_id`,
		"salesman": `
			LEFT JOIN mst.m_salesman s ON s.emp_id = r.salesman_id AND s.cust_id = rd.cust_id AND s.is_del = FALSE
			LEFT JOIN mst.m_employee e ON e.emp_id = s.emp_id AND e.cust_id = s.cust_id`,
		"product_category": `
			LEFT JOIN mst.m_product mp ON mp.pro_id = rd.product_id AND mp.cust_id = rd.cust_id
			LEFT JOIN mst.m_product_cat mpc ON mpc.pcat_id = mp.pcat_id`,
		"product": `
			LEFT JOIN mst.m_product mp ON mp.pro_id = rd.product_id AND mp.cust_id = rd.cust_id`,
	}

	return fmt.Sprintf(`
		WITH grouped AS (
			SELECT %s
			FROM sls."order" o
			JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id
			%s
			WHERE o.cust_id IN ? AND o.data_status IN (6,7) AND o.invoice_date >= ? AND o.invoice_date < ?

			UNION ALL

			SELECT %s
			FROM sls.return_det rd
			JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
			JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
			%s
			WHERE rd.cust_id IN ? AND o.data_status IN (6,7) AND o.invoice_date >= ? AND o.invoice_date < ?
		)
		SELECT id, code, name, COALESCE(SUM(net_sales), 0) AS net_sales
		FROM grouped
		GROUP BY id, code, name
		ORDER BY net_sales DESC
	`, orderSelect[groupBy], orderJoin[groupBy], returnSelect[groupBy], returnJoin[groupBy])
}

func (repository *RepositoryReportImpl) SecondarySalesReportGroupOutlet(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error) {
	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0)

	query := buildSecondarySalesReportGroupQuery("outlet")
	err = repository.Raw(query, custIDs, dateFrom, dateTo, custIDs, dateFrom, dateTo).Find(&results).Error

	return
}

func (repository *RepositoryReportImpl) SecondarySalesReportGroupSalesman(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error) {
	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0)

	query := buildSecondarySalesReportGroupQuery("salesman")
	err = repository.Raw(query, custIDs, dateFrom, dateTo, custIDs, dateFrom, dateTo).Find(&results).Error

	return
}

func (repository *RepositoryReportImpl) SecondarySalesReportProductCategory(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error) {
	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0)

	query := buildSecondarySalesReportGroupQuery("product_category")
	err = repository.Raw(query, custIDs, dateFrom, dateTo, custIDs, dateFrom, dateTo).Find(&results).Error

	return
}

func (repository *RepositoryReportImpl) SecondarySalesReportProduct(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error) {
	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0)

	query := buildSecondarySalesReportGroupQuery("product")
	err = repository.Raw(query, custIDs, dateFrom, dateTo, custIDs, dateFrom, dateTo).Find(&results).Error

	return
}

func buildSecondarySalesReportTrendSalesSQL() string {
	return `
		WITH months AS (
			SELECT 1 AS month UNION ALL
			SELECT 2 UNION ALL
			SELECT 3 UNION ALL
			SELECT 4 UNION ALL
			SELECT 5 UNION ALL
			SELECT 6 UNION ALL
			SELECT 7 UNION ALL
			SELECT 8 UNION ALL
			SELECT 9 UNION ALL
			SELECT 10 UNION ALL
			SELECT 11 UNION ALL
			SELECT 12
		),
		order_summary AS (
			SELECT
				EXTRACT(MONTH FROM o.invoice_date)::INTEGER AS month,
				SUM(
					(COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price_final1, 0)) +
					(COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price_final2, 0)) +
					(COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price_final3, 0))
				) AS gross_sales,
				SUM(
					COALESCE(od.disc_value_final, 0) +
					COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) +
					COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) +
					COALESCE(od.promo_final5, 0)
				) AS discount_promo,
				SUM(COALESCE(od.vat_value_final, 0)) AS ppn
			FROM sls."order" o
			JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id
			WHERE o.cust_id IN ? AND o.data_status IN (6, 7)
				AND o.invoice_date >= ? AND o.invoice_date < ?
			GROUP BY EXTRACT(MONTH FROM o.invoice_date)::INTEGER
		),
		return_summary AS (
			SELECT
				EXTRACT(MONTH FROM r.return_date)::INTEGER AS month,
				SUM(
					(COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)) +
					(COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)) +
					(COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0))
				) AS gross_sales,
				SUM(
					COALESCE(rd.disc_value, 0) +
					COALESCE(rd.promo_value, 0)
				) AS discount_promo,
				SUM(COALESCE(rd.vat_value, 0)) AS ppn
			FROM sls.return_det rd
			JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
			JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
			WHERE rd.cust_id IN ? AND o.data_status IN (6, 7)
				AND r.return_date >= ? AND r.return_date < ?
			GROUP BY EXTRACT(MONTH FROM r.return_date)::INTEGER
		)
		SELECT
			m.month AS month,
			COALESCE(os.gross_sales, 0) - COALESCE(rs.gross_sales, 0) AS total_gross_sale,
			COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo,
			((COALESCE(os.gross_sales, 0) - COALESCE(os.discount_promo, 0)) -
			 (COALESCE(rs.gross_sales, 0) - COALESCE(rs.discount_promo, 0))) +
			(COALESCE(os.ppn, 0) - COALESCE(rs.ppn, 0)) AS net_sales
		FROM months m
		LEFT JOIN order_summary os ON os.month = m.month
		LEFT JOIN return_summary rs ON rs.month = m.month
		ORDER BY m.month
	`
}

func (repository *RepositoryReportImpl) SecondarySalesReportTrendSales(custIDs []string, year int) (results []model.TrendSalesSecondarySalesModel, err error) {
	dateFrom := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(1, 0, 0)

	err = repository.
		Raw(buildSecondarySalesReportTrendSalesSQL(), custIDs, dateFrom, dateTo, custIDs, dateFrom, dateTo).
		Find(&results).Error

	return
}

func buildActivityReportTrendSalesSQL() string {
	return `
		WITH months AS (
			SELECT 1 AS month_num UNION ALL
			SELECT 2 UNION ALL
			SELECT 3 UNION ALL
			SELECT 4 UNION ALL
			SELECT 5 UNION ALL
			SELECT 6 UNION ALL
			SELECT 7 UNION ALL
			SELECT 8 UNION ALL
			SELECT 9 UNION ALL
			SELECT 10 UNION ALL
			SELECT 11 UNION ALL
			SELECT 12
		),
		order_data AS (
			SELECT
				EXTRACT(MONTH FROM o.invoice_date)::INTEGER AS month_num,
				SUM(
					(
						(COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price_final1, 0)) +
						(COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price_final2, 0)) +
						(COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price_final3, 0))
					) -
					COALESCE(od.disc_value_final, 0) -
					(
						COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) +
						COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) +
						COALESCE(od.promo_final5, 0)
					) +
					COALESCE(od.vat_value_final, 0)
				) AS net_sales_order
			FROM sls."order" o
			JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id
			WHERE o.cust_id IN ? AND o.data_status IN (6, 7)
				AND o.invoice_date >= ? AND o.invoice_date < ?
			GROUP BY EXTRACT(MONTH FROM o.invoice_date)::INTEGER
		),
		return_data AS (
			SELECT
				EXTRACT(MONTH FROM r.return_date)::INTEGER AS month_num,
				SUM(
					(
						(COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)) +
						(COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)) +
						(COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0))
					) -
					COALESCE(rd.promo_value, 0) -
					COALESCE(rd.disc_value, 0) +
					COALESCE(rd.vat_value, 0)
				) AS net_sales_return
			FROM sls.return_det rd
			JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
			JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
			WHERE rd.cust_id IN ? AND o.data_status IN (6, 7)
				AND r.return_date >= ? AND r.return_date < ?
			GROUP BY EXTRACT(MONTH FROM r.return_date)::INTEGER
		)
		SELECT
			m.month_num AS month,
			COALESCE(o.net_sales_order, 0) AS total_invoice,
			COALESCE(r.net_sales_return, 0) AS total_return,
			COALESCE(o.net_sales_order, 0) - COALESCE(r.net_sales_return, 0) AS net_sales
		FROM months m
		LEFT JOIN order_data o ON o.month_num = m.month_num
		LEFT JOIN return_data r ON r.month_num = m.month_num
		ORDER BY m.month_num
	`
}

func (repository *RepositoryReportImpl) SalesmanActivityReportTrendSales(custIDs []string, year int) (results []model.ActivityReportTrendSalesModel, err error) {
	dateFrom := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(1, 0, 0)

	err = repository.
		Raw(buildActivityReportTrendSalesSQL(), custIDs, dateFrom, dateTo, custIDs, dateFrom, dateTo).
		Find(&results).Error

	return
}

func buildSalesmanActivitySumByMonthSQL() string {
	return `
		WITH order_summary AS (
			SELECT
				SUM(
					(COALESCE(od.qty1_final, 0) * COALESCE(od.sell_price_final1, 0)) +
					(COALESCE(od.qty2_final, 0) * COALESCE(od.sell_price_final2, 0)) +
					(COALESCE(od.qty3_final, 0) * COALESCE(od.sell_price_final3, 0))
				) AS gross_sales,
				SUM(
					COALESCE(od.disc_value_final, 0) +
					COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) +
					COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) +
					COALESCE(od.promo_final5, 0)
				) AS discount_promo,
				SUM(COALESCE(od.vat_value_final, 0)) AS ppn,
				MAX(o.updated_at) AS last_update
			FROM sls."order" o
			JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id
			WHERE o.cust_id IN ? AND o.data_status IN (6, 7)
				AND o.invoice_date >= ? AND o.invoice_date < ?
		),
		return_summary AS (
			SELECT
				SUM(
					(COALESCE(rd.qty1, 0) * COALESCE(rd.sell_price1, 0)) +
					(COALESCE(rd.qty2, 0) * COALESCE(rd.sell_price2, 0)) +
					(COALESCE(rd.qty3, 0) * COALESCE(rd.sell_price3, 0))
				) AS gross_sales,
				SUM(
					COALESCE(rd.disc_value, 0) +
					COALESCE(rd.promo_value, 0)
				) AS discount_promo,
				SUM(COALESCE(rd.vat_value, 0)) AS ppn,
				MAX(r.updated_at) AS last_update
			FROM sls.return_det rd
			JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
			JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
			WHERE rd.cust_id IN ? AND o.data_status IN (6, 7)
				AND o.invoice_date >= ? AND o.invoice_date < ?
		),
		salesman_summary AS (
			SELECT COUNT(DISTINCT trx.salesman_id) AS total_salesman
			FROM (
				SELECT o.salesman_id
				FROM sls."order" o
				JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id
				WHERE o.cust_id IN ? AND o.data_status IN (6, 7)
					AND o.invoice_date >= ? AND o.invoice_date < ?
				UNION
				SELECT r.salesman_id
				FROM sls.return_det rd
				JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
				JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
				WHERE rd.cust_id IN ? AND o.data_status IN (6, 7)
					AND o.invoice_date >= ? AND o.invoice_date < ?
			) trx
			WHERE trx.salesman_id IS NOT NULL
		)
		SELECT
			((COALESCE(os.gross_sales, 0) - COALESCE(os.discount_promo, 0)) -
			 (COALESCE(rs.gross_sales, 0) - COALESCE(rs.discount_promo, 0))) +
			(COALESCE(os.ppn, 0) - COALESCE(rs.ppn, 0)) AS total_sales,
			((COALESCE(rs.gross_sales, 0) - COALESCE(rs.discount_promo, 0))) +
			COALESCE(rs.ppn, 0) AS total_return,
			COALESCE(ss.total_salesman, 0) AS total_salesman,
			CASE
				WHEN os.last_update IS NULL THEN rs.last_update
				WHEN rs.last_update IS NULL THEN os.last_update
				ELSE GREATEST(os.last_update, rs.last_update)
			END AS last_update
		FROM order_summary os
		CROSS JOIN return_summary rs
		CROSS JOIN salesman_summary ss
	`
}

func (repository *RepositoryReportImpl) SalesmanActivityReportSumByMonth(custIDs []string, month int, year int) (data model.SalesmanActivitySumByMonthModel, err error) {
	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0)

	err = repository.Raw(
		buildSalesmanActivitySumByMonthSQL(),
		custIDs, dateFrom, dateTo,
		custIDs, dateFrom, dateTo,
		custIDs, dateFrom, dateTo,
		custIDs, dateFrom, dateTo,
	).Take(&data).Error

	return
}

func (repository *RepositoryReportImpl) ActivitySalesmanReportGroupSalesman(custIDs []string, month int, year int) (results []model.SecondarySalesReportGroup, err error) {
	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0)

	err = repository.Raw(
		buildActivitySalesmanGroupSalesSQL(),
		custIDs, dateFrom, dateTo,
		custIDs, dateFrom, dateTo,
	).Find(&results).Error

	return
}

func (repository *RepositoryReportImpl) ActivitySalesmanReturnReportGroupSalesman(custIDs []string, month int, year int) (results []model.ReturnReportGroup, err error) {
	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	dateTo := dateFrom.AddDate(0, 1, 0)

	err = repository.Raw(
		buildActivitySalesmanGroupReturnSQL(),
		custIDs, dateFrom, dateTo,
	).Find(&results).Error

	return
}

// CountDownloadSalesOrderInProgress checks if there's an in-progress DownloadSalesOrder report
func (repository *RepositoryReportImpl) CountDownloadSalesOrderInProgress(custID string) (int64, error) {
	var count int64
	err := repository.Model(&model.ReportList{}).
		Where("cust_id = ? AND report_name LIKE ? AND file_status = ?",
			custID, "DownloadSalesOrder%", entity.FILE_STATUS_PROCESSING).
		Count(&count).Error
	return count, err
}

// CountDownloadSalesOrderByDate counts DownloadSalesOrder reports for running number generation
func (repository *RepositoryReportImpl) CountDownloadSalesOrderByDate(custID, exportDate string) int64 {
	var count int64
	// Parse exportDate (ddmmyy format) to get today's range
	loc, _ := time.LoadLocation("Asia/Jakarta")
	startOfDay := time.Now().In(loc).Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	repository.Model(&model.ReportList{}).
		Where("cust_id = ? AND created_at >= ? AND created_at < ? AND report_name LIKE ?",
			custID, startOfDay, endOfDay, "DownloadSalesOrder%").
		Count(&count)
	return count + 1
}
