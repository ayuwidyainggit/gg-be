package repository

const activityReportGroupOrderDataCTE = `
order_data AS (
    SELECT
        o.salesman_id,
        o.cust_id,
        COALESCE(od.qty1_final, 0) AS qty1,
        COALESCE(od.qty2_final, 0) AS qty2,
        COALESCE(od.qty3_final, 0) AS qty3,
        COALESCE(od.sell_price_final1, 0) AS price1,
        COALESCE(od.sell_price_final2, 0) AS price2,
        COALESCE(od.sell_price_final3, 0) AS price3,
        COALESCE(od.disc_value_final, 0) AS special_discount,
        (
            COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) +
            COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) +
            COALESCE(od.promo_final5, 0)
        ) AS discount,
        COALESCE(od.vat_value_final, 0) AS ppn,
        1 AS multiplier
    FROM sls."order" o
    JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = o.cust_id
    WHERE o.cust_id IN ?
      AND o.data_status IN (6, 7)
      AND o.invoice_date >= ? AND o.invoice_date < ?
)`

const activityReportGroupReturnDataCTE = `
return_data AS (
    SELECT
        r.salesman_id,
        r.cust_id,
        COALESCE(rd.qty1, 0) AS qty1,
        COALESCE(rd.qty2, 0) AS qty2,
        COALESCE(rd.qty3, 0) AS qty3,
        COALESCE(rd.sell_price1, 0) AS price1,
        COALESCE(rd.sell_price2, 0) AS price2,
        COALESCE(rd.sell_price3, 0) AS price3,
        COALESCE(rd.promo_value, 0) AS special_discount,
        COALESCE(rd.disc_value, 0) AS discount,
        COALESCE(rd.vat_value, 0) AS ppn,
        -1 AS multiplier
    FROM sls.return_det rd
    JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
    JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
    WHERE rd.cust_id IN ?
      AND o.data_status IN (6, 7)
      AND o.invoice_date >= ? AND o.invoice_date < ?
)`

const activityReportGroupReturnOnlyDataCTE = `
return_data AS (
    SELECT
        r.salesman_id,
        r.cust_id,
        COALESCE(rd.qty1, 0) AS qty1,
        COALESCE(rd.qty2, 0) AS qty2,
        COALESCE(rd.qty3, 0) AS qty3,
        COALESCE(rd.sell_price1, 0) AS price1,
        COALESCE(rd.sell_price2, 0) AS price2,
        COALESCE(rd.sell_price3, 0) AS price3,
        COALESCE(rd.promo_value, 0) AS special_discount,
        COALESCE(rd.disc_value, 0) AS discount,
        COALESCE(rd.vat_value, 0) AS ppn
    FROM sls.return_det rd
    JOIN sls."return" r ON r.return_no = rd.return_no AND r.cust_id = rd.cust_id
    JOIN sls."order" o ON o.invoice_no = r.invoice_no AND o.cust_id = r.cust_id
    WHERE rd.cust_id IN ?
      AND o.data_status IN (6, 7)
      AND o.invoice_date >= ? AND o.invoice_date < ?
)`

const activityReportGroupNetIncVATExpr = `
(
    (t.qty1 * t.price1) +
    (t.qty2 * t.price2) +
    (t.qty3 * t.price3)
) - t.special_discount - t.discount + t.ppn`

const activityReportGroupSalesmanJoin = `
LEFT JOIN mst.m_salesman ms ON ms.emp_id = t.salesman_id AND ms.cust_id = t.cust_id AND ms.is_del = FALSE
LEFT JOIN mst.m_employee me ON me.emp_id = ms.emp_id AND me.cust_id = ms.cust_id`

func buildActivitySalesmanGroupSalesSQL() string {
	return `
WITH ` + activityReportGroupOrderDataCTE + `,
` + activityReportGroupReturnDataCTE + `,
trx AS (
    SELECT * FROM order_data
    UNION ALL
    SELECT * FROM return_data
)
SELECT
    t.salesman_id AS id,
    COALESCE(me.emp_code, '') AS code,
    COALESCE(ms.sales_name, '') AS name,
    SUM((` + activityReportGroupNetIncVATExpr + `) * t.multiplier) AS net_sales
FROM trx t
` + activityReportGroupSalesmanJoin + `
GROUP BY t.salesman_id, me.emp_code, ms.sales_name
ORDER BY net_sales DESC`
}

func buildActivitySalesmanGroupReturnSQL() string {
	return `
WITH ` + activityReportGroupReturnOnlyDataCTE + `
SELECT
    t.salesman_id AS id,
    COALESCE(me.emp_code, '') AS code,
    COALESCE(ms.sales_name, '') AS name,
    SUM(` + activityReportGroupNetIncVATExpr + `) AS net_sales
FROM return_data t
` + activityReportGroupSalesmanJoin + `
GROUP BY t.salesman_id, me.emp_code, ms.sales_name
ORDER BY net_sales DESC`
}
