package repository

import (
	"fmt"
	"sales/model"
	"strings"

	"gorm.io/gorm"
)

type activityReportPJPConfig struct {
	Schema             string
	FilterOutletByCust bool
	IsPrincipal        bool
}

type activityReportSQLParams struct {
	CustID                   string
	CustIDs                  []string
	SalesmanReferenceCustIDs []string
	AuthCustID               string
	ParentCustID             string
	FromDate                 string
	ToDate                   string
	SalesmanIDs              []int
	DistributorCodes         []string
	SkipSalesmanFilter       bool
	Limit                    int
	Offset                   int
	ForCount                 bool
}

func activityReportEffectiveCustIDs(p activityReportSQLParams) []string {
	if len(p.CustIDs) > 0 {
		return p.CustIDs
	}
	if id := strings.TrimSpace(p.CustID); id != "" {
		return []string{id}
	}
	return nil
}

func activityReportCustIDPredicate(column string, custIDs []string) (string, []interface{}) {
	if len(custIDs) == 0 {
		return "FALSE", nil
	}
	if len(custIDs) == 1 {
		return column + " = ?", []interface{}{custIDs[0]}
	}
	placeholders := make([]string, len(custIDs))
	args := make([]interface{}, len(custIDs))
	for i, id := range custIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	return column + " IN (" + strings.Join(placeholders, ",") + ")", args
}

func activityReportPrincipalScopeCustID(p activityReportSQLParams) string {
	if strings.TrimSpace(p.ParentCustID) != "" {
		return strings.TrimSpace(p.ParentCustID)
	}
	return strings.TrimSpace(p.CustID)
}

const activityReportPrincipalCustIDMaxLen = 6

func isActivityReportPrincipalCust(custID, parentCustID string) bool {
	custID = strings.TrimSpace(custID)
	parentCustID = strings.TrimSpace(parentCustID)
	if custID == "" {
		return false
	}
	if parentCustID != "" {
		return custID == parentCustID
	}
	return len(custID) <= activityReportPrincipalCustIDMaxLen
}

func splitActivityReportCustIDs(parentCustID string, custIDs []string) (principalIDs, distributorIDs []string) {
	for _, id := range custIDs {
		if isActivityReportPrincipalCust(id, parentCustID) {
			principalIDs = append(principalIDs, id)
			continue
		}
		distributorIDs = append(distributorIDs, id)
	}
	return principalIDs, distributorIDs
}

func hasActivityReportPrincipalPJP(db *gorm.DB, principalIDs []string) bool {
	if len(principalIDs) == 0 {
		return false
	}
	var count int64
	_ = db.Raw(`SELECT COUNT(*) FROM pjp_principles.permanent_journey_plans WHERE cust_id IN ?`, principalIDs).Scan(&count)
	return count > 0
}

func resolveActivityReportPJPConfig(db *gorm.DB, authCustID, parentCustID string, custIDs []string) activityReportPJPConfig {
	if isActivityReportPrincipalCust(authCustID, parentCustID) {
		principalIDs, _ := splitActivityReportCustIDs(parentCustID, custIDs)
		scopeIDs := principalIDs
		if len(scopeIDs) == 0 {
			scopeIDs = custIDs
		}
		if len(scopeIDs) == 0 {
			scopeIDs = []string{authCustID}
		}
		if hasActivityReportPrincipalPJP(db, scopeIDs) {
			return activityReportPJPConfig{Schema: "pjp_principles", FilterOutletByCust: false, IsPrincipal: true}
		}
		return activityReportPJPConfig{Schema: "pjp", FilterOutletByCust: false, IsPrincipal: false}
	}
	return activityReportPJPConfig{Schema: "pjp", FilterOutletByCust: true, IsPrincipal: false}
}

func appendActivityReportSalesmanFilter(p activityReportSQLParams) (string, []interface{}) {
	if p.SkipSalesmanFilter {
		return "", nil
	}
	return appendActivityReportSalesmanFilterForBranch(p, activityReportEffectiveCustIDs(p))
}

func appendActivityReportOrderSalesmanFilter(p activityReportSQLParams) (string, []interface{}) {
	if p.SkipSalesmanFilter {
		return "", nil
	}
	return appendActivityReportOrderSalesmanFilterForBranch(p, activityReportEffectiveCustIDs(p))
}

func appendActivityReportPaymentSalesmanFilter(p activityReportSQLParams) (string, []interface{}) {
	if p.SkipSalesmanFilter {
		return "", nil
	}
	return appendActivityReportPaymentSalesmanFilterForBranch(p, activityReportEffectiveCustIDs(p))
}

func activityReportSalesmanReferenceCustIDs(p activityReportSQLParams) []string {
	if len(p.SalesmanReferenceCustIDs) > 0 {
		return p.SalesmanReferenceCustIDs
	}
	return activityReportEffectiveCustIDs(p)
}

func appendActivityReportSalesmanFilterForBranch(p activityReportSQLParams, branchCustIDs []string) (string, []interface{}) {
	return buildActivityReportSalesmanEmpIDFilter("p.salesman_id", branchCustIDs, activityReportSalesmanReferenceCustIDs(p), p.SalesmanIDs)
}

func appendActivityReportOrderSalesmanFilterForBranch(p activityReportSQLParams, branchCustIDs []string) (string, []interface{}) {
	return buildActivityReportSalesmanEmpIDFilter("salesman_id", branchCustIDs, activityReportSalesmanReferenceCustIDs(p), p.SalesmanIDs)
}

func appendActivityReportPaymentSalesmanFilterForBranch(p activityReportSQLParams, branchCustIDs []string) (string, []interface{}) {
	return buildActivityReportSalesmanEmpIDFilter("o.salesman_id", branchCustIDs, activityReportSalesmanReferenceCustIDs(p), p.SalesmanIDs)
}

func buildActivityReportSalesmanEmpIDFilter(column string, branchCustIDs, referenceCustIDs []string, salesmanIDs []int) (string, []interface{}) {
	if len(salesmanIDs) == 0 {
		return "", nil
	}
	if len(branchCustIDs) == 0 {
		branchCustIDs = referenceCustIDs
	}
	if len(referenceCustIDs) == 0 {
		referenceCustIDs = branchCustIDs
	}

	salesmanPlaceholders := make([]string, len(salesmanIDs))
	salesmanArgs := make([]interface{}, len(salesmanIDs))
	for i, id := range salesmanIDs {
		salesmanPlaceholders[i] = "?"
		salesmanArgs[i] = id
	}

	if len(branchCustIDs) == 0 || len(referenceCustIDs) == 0 {
		return fmt.Sprintf(" AND %s IN (%s)", column, strings.Join(salesmanPlaceholders, ",")), salesmanArgs
	}

	branchIn, branchArgs := activityReportInClauseArgs(branchCustIDs)
	refIn, refArgs := activityReportInClauseArgs(referenceCustIDs)

	sql := fmt.Sprintf(` AND %s IN (
		SELECT e.emp_id FROM mst.m_employee e
		WHERE e.cust_id IN (%s)
		AND e.emp_code IN (
			SELECT DISTINCT e2.emp_code FROM mst.m_employee e2
			WHERE e2.emp_id IN (%s) AND e2.cust_id IN (%s)
		)
	)`, column, branchIn, strings.Join(salesmanPlaceholders, ","), refIn)

	args := make([]interface{}, 0, len(branchArgs)+len(salesmanArgs)+len(refArgs))
	args = append(args, branchArgs...)
	args = append(args, salesmanArgs...)
	args = append(args, refArgs...)
	return sql, args
}

func activityReportInClauseArgs(values []string) (string, []interface{}) {
	placeholders := make([]string, len(values))
	args := make([]interface{}, len(values))
	for i, value := range values {
		placeholders[i] = "?"
		args[i] = value
	}
	return strings.Join(placeholders, ","), args
}

const activityReportClockTimeWIBSelect = `
    CASE WHEN att.clock_in_time IS NOT NULL THEN TO_CHAR(att.clock_in_time, 'YYYY-MM-DD HH24:MI') ELSE '' END AS clock_in_time,
    CASE WHEN att.clock_out_time IS NOT NULL THEN TO_CHAR(att.clock_out_time, 'YYYY-MM-DD HH24:MI') ELSE '' END AS clock_out_time`

const activityReportVisitTimeWIBSelect = `
    CASE WHEN p.arrive_at IS NOT NULL THEN TO_CHAR(TO_TIMESTAMP(p.arrive_at / 1000.0) + INTERVAL '7 hour', 'YYYY-MM-DD HH24:MI:SS') ELSE '' END AS checkin_time,
    CASE WHEN p.leave_at IS NOT NULL THEN TO_CHAR(TO_TIMESTAMP(p.leave_at / 1000.0) + INTERVAL '7 hour', 'YYYY-MM-DD HH24:MI:SS') ELSE '' END AS checkout_time`

const activityReportRemarksSelect = `
    CASE
        WHEN EXISTS (
            SELECT 1
            FROM mobile.leave_request lr
            WHERE lr.emp_id = p.salesman_id
              AND lr.cust_id = p.cust_id
              AND LOWER(TRIM(lr.approval)) IN ('approve', 'approved')
              AND p.date::date BETWEEN lr.start_date AND lr.end_date
        ) THEN 'On Leave'
        ELSE '-'
    END AS remarks`

const activityReportSalesmanNameSelect = `COALESCE(NULLIF(TRIM(e.emp_name), ''), p.salesman_name) AS salesman_name`

const activityReportListOrderBy = ` ORDER BY salesman_name ASC, visit_date ASC`

const activityReportEffectiveDistributorCodeExpr = `COALESCE(NULLIF(TRIM(distributor_code), ''), business_unit_code)`

func buildActivitySalesReportSQLBase(cfg activityReportPJPConfig, p activityReportSQLParams) (string, []interface{}) {
	if cfg.IsPrincipal && cfg.Schema == "pjp_principles" {
		return buildActivitySalesReportPrincipalQuery(p)
	}
	return buildActivitySalesReportDistributorQuery(cfg, p)
}

func buildActivitySalesReportSQL(cfg activityReportPJPConfig, p activityReportSQLParams) (string, []interface{}) {
	base, args := buildActivitySalesReportSQLBase(cfg, p)
	return finalizeActivityReportSQL(base, p, args)
}

func buildActivitySalesReportPrincipalSQL(p activityReportSQLParams) (string, []interface{}) {
	base, args := buildActivitySalesReportPrincipalQuery(p)
	return finalizeActivityReportSQL(base, p, args)
}

func buildActivitySalesReportPrincipalQuery(p activityReportSQLParams) (string, []interface{}) {
	custIDs := activityReportEffectiveCustIDs(p)
	pjpCustWhere, pjpCustArgs := activityReportCustIDPredicate("p.cust_id", custIDs)
	empCustWhere, empCustArgs := activityReportCustIDPredicate("cust_id", custIDs)
	custDataWhere, custDataArgs := activityReportCustIDPredicate("mc.cust_id", custIDs)
	orderCustWhere, orderCustArgs := activityReportCustIDPredicate("cust_id", custIDs)
	returnCustWhere, returnCustArgs := activityReportCustIDPredicate("cust_id", custIDs)
	paymentCustWhere, paymentCustArgs := activityReportCustIDPredicate("d.cust_id", custIDs)

	salesmanFilterSQL, salesmanArgs := appendActivityReportSalesmanFilter(p)
	args := make([]interface{}, 0)
	args = append(args, pjpCustArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, salesmanArgs...)
	args = append(args, empCustArgs...)
	args = append(args, custDataArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, orderCustArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, returnCustArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, paymentCustArgs...)
	args = append(args, p.FromDate, p.ToDate)

	base := fmt.Sprintf(`
WITH pjp_data AS (
    SELECT
        v.id AS visit_id,
        v.pjp_code,
        v.outlet_code,
        v.outlet_id,
        v.date,
        v.arrive_at,
        v.leave_at,
        v.is_planned,
        v.skip_at,
        v.location_status,
        v.latitude AS v_lat,
        v.longitude AS v_lon,
        p.cust_id,
        p.salesman_id,
        p.salesman_name
    FROM pjp_principles.outlet_visit_list v
    JOIN pjp_principles.permanent_journey_plans p ON v.pjp_code = p.pjp_code
    WHERE %s
      AND v.date::date BETWEEN ?::date AND ?::date
      %s
),
emp_data AS (
    SELECT emp_id, emp_code, emp_name
    FROM mst.m_employee
    WHERE %s
),
	outlet_data AS (
    SELECT outlet_id, outlet_name, outlet_principal_code, latitude, longitude
    FROM mst.m_outlet
),
outlet_dist_data AS (
    SELECT
        o.outlet_id,
        COALESCE(md.distributor_code, '') AS distributor_code,
        COALESCE(md.distributor_name, '') AS distributor_name
    FROM mst.m_outlet o
    LEFT JOIN mst.m_distributor md ON md.cust_id = o.cust_id
),
cust_data AS (
    SELECT
        mc.cust_id,
        mc.cust_name,
        COALESCE(md.distributor_code, '') AS distributor_code
    FROM smc.m_customer mc
    LEFT JOIN mst.m_distributor md ON md.cust_id = mc.cust_id
    WHERE %s
),
dh_data AS (
    SELECT DISTINCT ON (pjp_code, date::date)
        pjp_code,
        date::date AS visit_date,
        destination_id
    FROM pjp_principles.destinations_history
    WHERE destination_type = 'distributor'
    ORDER BY pjp_code, date::date, id DESC
),
dh_all AS (
    SELECT DISTINCT ON (pjp_code, date::date, destination_id)
        pjp_code,
        date::date AS visit_date,
        destination_id,
        destination_type
    FROM pjp_principles.destinations_history
    ORDER BY pjp_code, date::date, destination_id, id DESC
),
dist_data AS (
    SELECT distributor_id, distributor_code, distributor_name, latitude, longitude
    FROM mst.m_distributor
),
attendance_data AS (
    SELECT
        emp_code,
        created_at::date AS date,
        MIN(CASE WHEN type = 1 THEN created_at END) AS clock_in_time,
        MAX(CASE WHEN type = 2 THEN created_at END) AS clock_out_time
    FROM mobile.attendances
    WHERE created_at::date BETWEEN ?::date AND ?::date
    GROUP BY emp_code, created_at::date
),
order_data AS (
    SELECT outlet_id, ro_date::date AS trx_date, SUM(total_final) AS sales_value
    FROM sls."order"
    WHERE %s
      AND data_status != 9
      AND ro_date::date BETWEEN ?::date AND ?::date
    GROUP BY outlet_id, ro_date::date
),
return_data AS (
    SELECT outlet_id, return_date::date AS trx_date, SUM(total) AS return_value
    FROM sls."return"
    WHERE %s
      AND data_status != 9
      AND return_date::date BETWEEN ?::date AND ?::date
    GROUP BY outlet_id, return_date::date
),
payment_data AS (
    SELECT
        o.outlet_id,
        d.deposit_date::date AS trx_date,
        SUM(dd.total_payment) AS payment_value
    FROM acf.deposit d
    LEFT JOIN acf.deposit_detail dd ON d.deposit_no = dd.deposit_no AND d.cust_id = dd.cust_id
    LEFT JOIN sls."order" o ON dd.invoice_no = o.invoice_no AND dd.cust_id = o.cust_id
    WHERE %s
      AND d.deposit_date::date BETWEEN ?::date AND ?::date
    GROUP BY o.outlet_id, d.deposit_date::date
)
SELECT
    '' AS business_unit_code,
    c.cust_name AS business_unit_name,
    p.pjp_code::text AS pjp_code,
    e.emp_code,
    `+activityReportSalesmanNameSelect+`,
    COALESCE(NULLIF(TRIM(d.distributor_code), ''), NULLIF(TRIM(o_dist.distributor_code), '')) AS distributor_code,
    COALESCE(NULLIF(TRIM(d.distributor_name), ''), NULLIF(TRIM(o_dist.distributor_name), '')) AS distributor_name,
    CASE WHEN dh_all.destination_type = 'distributor' THEN '' ELSE p.outlet_code END AS outlet_code,
    CASE WHEN dh_all.destination_type = 'distributor' THEN '' ELSE o.outlet_principal_code END AS outlet_principal_code,
    CASE WHEN dh_all.destination_type = 'distributor' THEN '' ELSE o.outlet_name END AS outlet_name,
    p.date AS visit_date,
`+activityReportClockTimeWIBSelect+`,
`+activityReportVisitTimeWIBSelect+`,
    CASE
        WHEN p.arrive_at IS NOT NULL AND p.leave_at IS NOT NULL AND (p.leave_at - p.arrive_at) > 0
        THEN ROUND((p.leave_at - p.arrive_at) / 60000.0)::bigint
        ELSE 0
    END AS duration_in_minutes,
    p.is_planned,
    p.skip_at,
    p.location_status,
    CASE WHEN p.is_planned THEN 'Planned' ELSE 'Unplanned' END AS pjp_status,
    CASE
        WHEN p.skip_at IS NOT NULL THEN 'Skipped'
        WHEN p.arrive_at IS NOT NULL THEN 'Visited'
        ELSE 'Not Visited'
    END AS visit_status,
    CASE WHEN p.is_planned = true AND p.arrive_at IS NOT NULL AND p.skip_at IS NULL THEN 'Yes' ELSE 'No' END AS compliance,
    COALESCE(od.sales_value, 0) AS sales_value,
    COALESCE(rd.return_value, 0) AS return_value,
    COALESCE(pd.payment_value, 0) AS payment_value,
    CASE
        WHEN dh_all.destination_type = 'distributor' AND d_visit.latitude IS NOT NULL AND d_visit.longitude IS NOT NULL
            THEN d_visit.latitude::text || ',' || d_visit.longitude::text
        WHEN (dh_all.destination_type IS DISTINCT FROM 'distributor') AND o.latitude IS NOT NULL AND o.longitude IS NOT NULL
            THEN o.latitude::text || ',' || o.longitude::text
        ELSE ''
    END AS location_master,
    CASE WHEN p.v_lat IS NOT NULL AND p.v_lon IS NOT NULL THEN p.v_lat::text || ',' || p.v_lon::text ELSE '' END AS location_actual,
    CASE
        WHEN p.location_status = 0 THEN 'Mismatch'
        WHEN p.location_status = 1 THEN 'Match'
        ELSE ''
    END AS geotag_status,
`+activityReportRemarksSelect+`
FROM pjp_data p
LEFT JOIN emp_data e ON p.salesman_id = e.emp_id
LEFT JOIN outlet_data o ON p.outlet_id = o.outlet_id
LEFT JOIN outlet_dist_data o_dist ON p.outlet_id = o_dist.outlet_id
LEFT JOIN cust_data c ON p.cust_id = c.cust_id
LEFT JOIN dh_data dh ON p.pjp_code = dh.pjp_code AND p.date::date = dh.visit_date
LEFT JOIN dh_all ON p.pjp_code = dh_all.pjp_code AND p.date::date = dh_all.visit_date AND p.outlet_id = dh_all.destination_id
LEFT JOIN dist_data d ON dh_all.destination_id = d.distributor_id
LEFT JOIN dist_data d_visit ON p.outlet_id = d_visit.distributor_id
LEFT JOIN attendance_data att ON e.emp_code = att.emp_code AND p.date::date = att.date
LEFT JOIN order_data od ON p.outlet_id = od.outlet_id AND p.date::date = od.trx_date
LEFT JOIN return_data rd ON p.outlet_id = rd.outlet_id AND p.date::date = rd.trx_date
LEFT JOIN payment_data pd ON p.outlet_id = pd.outlet_id AND p.date::date = pd.trx_date`,
		pjpCustWhere, salesmanFilterSQL, empCustWhere, custDataWhere, orderCustWhere, returnCustWhere, paymentCustWhere)

	return base, args
}

func buildActivitySalesReportDistributorSQL(cfg activityReportPJPConfig, p activityReportSQLParams) (string, []interface{}) {
	base, args := buildActivitySalesReportDistributorQuery(cfg, p)
	return finalizeActivityReportSQL(base, p, args)
}

func buildActivitySalesReportDistributorQuery(cfg activityReportPJPConfig, p activityReportSQLParams) (string, []interface{}) {
	schema := cfg.Schema
	if schema != "pjp" && schema != "pjp_principles" {
		schema = "pjp"
	}

	custIDs := activityReportEffectiveCustIDs(p)
	pjpCustWhere, pjpCustArgs := activityReportCustIDPredicate("p.cust_id", custIDs)
	empCustWhere, empCustArgs := activityReportCustIDPredicate("cust_id", custIDs)
	custDataWhere, custDataArgs := activityReportCustIDPredicate("mc.cust_id", custIDs)
	orderCustWhere, orderCustArgs := activityReportCustIDPredicate("cust_id", custIDs)
	returnCustWhere, returnCustArgs := activityReportCustIDPredicate("cust_id", custIDs)
	paymentCustWhere, paymentCustArgs := activityReportCustIDPredicate("d.cust_id", custIDs)

	salesmanFilterSQL, salesmanArgs := appendActivityReportSalesmanFilter(p)
	orderSalesmanFilterSQL, orderSalesmanArgs := appendActivityReportOrderSalesmanFilter(p)
	paymentSalesmanFilterSQL, paymentSalesmanArgs := appendActivityReportPaymentSalesmanFilter(p)

	outletFilterSQL := ""
	var outletFilterArgs []interface{}
	if cfg.FilterOutletByCust {
		outletWhere, outletArgs := activityReportCustIDPredicate("cust_id", custIDs)
		outletFilterSQL = "WHERE " + outletWhere
		outletFilterArgs = outletArgs
	}

	args := make([]interface{}, 0)
	args = append(args, pjpCustArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, salesmanArgs...)
	args = append(args, outletFilterArgs...)
	args = append(args, empCustArgs...)
	args = append(args, custDataArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, orderCustArgs...)
	args = append(args, orderSalesmanArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, returnCustArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, paymentCustArgs...)
	args = append(args, p.FromDate, p.ToDate)
	args = append(args, paymentSalesmanArgs...)

	base := fmt.Sprintf(`
WITH pjp_data AS (
    SELECT
        v.id AS visit_id,
        v.pjp_code,
        v.outlet_code,
        v.outlet_id,
        v.date,
        v.arrive_at,
        v.leave_at,
        v.is_planned,
        v.skip_at,
        v.location_status,
        v.latitude AS v_lat,
        v.longitude AS v_lon,
        p.cust_id,
        p.salesman_id,
        p.salesman_name
    FROM %s.outlet_visit_list v
    JOIN %s.permanent_journey_plans p ON v.pjp_code = p.pjp_code
    WHERE %s
      AND v.date::date BETWEEN ?::date AND ?::date
      %s
),
emp_data AS (
    SELECT emp_id, emp_code, emp_name
    FROM mst.m_employee
    WHERE %s
),
outlet_data AS (
    SELECT outlet_id, outlet_name, outlet_principal_code, latitude, longitude
    FROM mst.m_outlet
    %s
),
cust_data AS (
    SELECT
        mc.cust_id,
        mc.cust_name,
        md.distributor_code,
        md.distributor_name
    FROM smc.m_customer mc
    JOIN mst.m_distributor md ON md.cust_id = mc.cust_id
    WHERE %s
),
attendance_data AS (
    SELECT
        emp_code,
        created_at::date AS date,
        MIN(CASE WHEN type = 1 THEN created_at END) AS clock_in_time,
        MAX(CASE WHEN type = 2 THEN created_at END) AS clock_out_time
    FROM mobile.attendances
    WHERE created_at::date BETWEEN ?::date AND ?::date
    GROUP BY emp_code, created_at::date
),
order_data AS (
    SELECT outlet_id, ro_date::date AS trx_date, SUM(total) AS sales_value
    FROM sls."order"
    WHERE %s
      %s
      AND data_status != 9
      AND ro_date::date BETWEEN ?::date AND ?::date
    GROUP BY outlet_id, ro_date::date
),
return_data AS (
    SELECT outlet_id, return_date::date AS trx_date, SUM(total) AS return_value
    FROM sls."return"
    WHERE %s
      AND data_status != 9
      AND return_date::date BETWEEN ?::date AND ?::date
    GROUP BY outlet_id, return_date::date
),
payment_data AS (
    SELECT
        o.outlet_id,
        o.salesman_id,
        d.deposit_date::date AS trx_date,
        SUM(dd.total_payment) AS payment_value
    FROM acf.deposit d
    LEFT JOIN acf.deposit_detail dd ON d.deposit_no = dd.deposit_no AND d.cust_id = dd.cust_id
    LEFT JOIN sls."order" o ON dd.invoice_no = o.invoice_no AND dd.cust_id = o.cust_id
    WHERE %s
      AND d.deposit_date::date BETWEEN ?::date AND ?::date
      %s
    GROUP BY o.outlet_id, d.deposit_date::date, o.salesman_id
)
SELECT
    c.distributor_code AS business_unit_code,
    c.cust_name AS business_unit_name,
    p.pjp_code::text AS pjp_code,
    e.emp_code,
    `+activityReportSalesmanNameSelect+`,
    '' AS distributor_code,
    '' AS distributor_name,
    p.outlet_code,
    o.outlet_principal_code,
    o.outlet_name,
    p.date AS visit_date,
`+activityReportClockTimeWIBSelect+`,
`+activityReportVisitTimeWIBSelect+`,
    CASE
        WHEN p.arrive_at IS NOT NULL AND p.leave_at IS NOT NULL THEN ROUND((p.leave_at - p.arrive_at) / 60000.0)::bigint
        ELSE 0
    END AS duration_in_minutes,
    p.is_planned,
    p.skip_at,
    p.location_status,
    CASE WHEN p.is_planned THEN 'Planned' ELSE 'Unplanned' END AS pjp_status,
    CASE
        WHEN p.skip_at IS NOT NULL THEN 'Skipped'
        WHEN p.arrive_at IS NOT NULL THEN 'Visited'
        ELSE 'Not Visited'
    END AS visit_status,
    CASE WHEN p.is_planned = true AND p.arrive_at IS NOT NULL AND p.skip_at IS NULL THEN 'Yes' ELSE 'No' END AS compliance,
    COALESCE(od.sales_value, 0) AS sales_value,
    COALESCE(rd.return_value, 0) AS return_value,
    COALESCE(pd.payment_value, 0) AS payment_value,
    CASE WHEN o.latitude IS NOT NULL AND o.longitude IS NOT NULL THEN o.latitude::text || ',' || o.longitude::text ELSE '' END AS location_master,
    CASE WHEN p.v_lat IS NOT NULL AND p.v_lon IS NOT NULL THEN p.v_lat::text || ',' || p.v_lon::text ELSE '' END AS location_actual,
    CASE
        WHEN p.location_status = 0 THEN 'Mismatch'
        WHEN p.location_status = 1 THEN 'Match'
        ELSE ''
    END AS geotag_status,
`+activityReportRemarksSelect+`
FROM pjp_data p
LEFT JOIN emp_data e ON p.salesman_id = e.emp_id
LEFT JOIN outlet_data o ON p.outlet_id = o.outlet_id
LEFT JOIN cust_data c ON p.cust_id = c.cust_id
LEFT JOIN attendance_data att ON e.emp_code = att.emp_code AND p.date::date = att.date
LEFT JOIN order_data od ON p.outlet_id = od.outlet_id AND p.date::date = od.trx_date
LEFT JOIN return_data rd ON p.outlet_id = rd.outlet_id AND p.date::date = rd.trx_date
LEFT JOIN payment_data pd ON p.outlet_id = pd.outlet_id AND p.date::date = pd.trx_date AND p.salesman_id = pd.salesman_id`,
		schema, schema, pjpCustWhere, salesmanFilterSQL, empCustWhere, outletFilterSQL, custDataWhere,
		orderCustWhere, orderSalesmanFilterSQL, returnCustWhere, paymentCustWhere, paymentSalesmanFilterSQL)

	return base, args
}

func finalizeActivityReportSQL(base string, p activityReportSQLParams, args []interface{}) (string, []interface{}) {
	if len(p.DistributorCodes) > 0 {
		whereSQL, whereArgs := activityReportCustIDPredicate(activityReportEffectiveDistributorCodeExpr, p.DistributorCodes)
		base = fmt.Sprintf("SELECT * FROM (%s) AS activity_report_src WHERE %s", base, whereSQL)
		args = append(args, whereArgs...)
	}

	if p.ForCount {
		return "SELECT COUNT(*) FROM (" + base + ") AS activity_report_cnt", args
	}
	sql := base + activityReportListOrderBy
	if p.Limit > 0 {
		sql += " LIMIT ?"
		args = append(args, p.Limit)
	}
	if p.Offset > 0 {
		sql += " OFFSET ?"
		args = append(args, p.Offset)
	}
	return sql, args
}

func appendActivityReportQueryArgs(allArgs []interface{}, args []interface{}) []interface{} {
	return append(allArgs, args...)
}

func buildActivitySalesReportCombinedSQL(
	db *gorm.DB,
	authCustID string,
	custIDs []string,
	parentCustID string,
	fromDate, toDate string,
	salesmanIDs []int,
	distributorCodes []string,
	limit, offset int,
	forCount bool,
) (string, []interface{}) {
	baseParams := activityReportSQLParams{
		AuthCustID:               authCustID,
		ParentCustID:             parentCustID,
		FromDate:                 fromDate,
		ToDate:                   toDate,
		SalesmanIDs:              salesmanIDs,
		DistributorCodes:         distributorCodes,
		SalesmanReferenceCustIDs: custIDs,
	}

	principalIDs, distributorIDs := splitActivityReportCustIDs(parentCustID, custIDs)
	if len(principalIDs) == 0 && len(distributorIDs) == 0 {
		effectiveCustIDs := activityReportEffectiveCustIDs(activityReportSQLParams{CustIDs: custIDs, CustID: authCustID})
		principalIDs, distributorIDs = splitActivityReportCustIDs(parentCustID, effectiveCustIDs)
	}

	var subqueries []string
	var allArgs []interface{}

	if len(principalIDs) > 0 && hasActivityReportPrincipalPJP(db, principalIDs) {
		p := baseParams
		p.CustIDs = principalIDs
		if len(principalIDs) == 1 {
			p.CustID = principalIDs[0]
		}
		sql, args := buildActivitySalesReportPrincipalQuery(
			p,
		)
		subqueries = append(subqueries, sql)
		allArgs = appendActivityReportQueryArgs(allArgs, args)
	}

	if len(distributorIDs) > 0 {
		p := baseParams
		p.CustIDs = distributorIDs
		if len(distributorIDs) == 1 {
			p.CustID = distributorIDs[0]
		}
		// Principal salesman IDs do not map to distributor visits; when both BU types
		// are selected, distributor rows are scoped by cust_id + date only.
		if len(principalIDs) > 0 && len(p.SalesmanIDs) > 0 {
			p.SkipSalesmanFilter = true
		}
		sql, args := buildActivitySalesReportDistributorQuery(
			activityReportPJPConfig{Schema: "pjp", FilterOutletByCust: true, IsPrincipal: false},
			p,
		)
		subqueries = append(subqueries, sql)
		allArgs = appendActivityReportQueryArgs(allArgs, args)
	}

	if len(subqueries) == 0 {
		cfg := resolveActivityReportPJPConfig(db, authCustID, parentCustID, custIDs)
		p := baseParams
		p.CustIDs = custIDs
		p.SalesmanReferenceCustIDs = custIDs
		p.Limit = limit
		p.Offset = offset
		p.ForCount = forCount
		return buildActivitySalesReportSQL(cfg, p)
	}

	if len(subqueries) == 1 {
		p := baseParams
		p.Limit = limit
		p.Offset = offset
		p.ForCount = forCount
		return finalizeActivityReportSQL(subqueries[0], p, allArgs)
	}

	unionBase := "(" + subqueries[0] + ")"
	for i := 1; i < len(subqueries); i++ {
		unionBase += " UNION ALL (" + subqueries[i] + ")"
	}

	if forCount {
		p := baseParams
		p.ForCount = true
		return finalizeActivityReportSQL(unionBase, p, allArgs)
	}

	p := baseParams
	p.Limit = limit
	p.Offset = offset
	return finalizeActivityReportSQL(unionBase, p, allArgs)
}

func (repository *RepositoryReportImpl) queryActivitySalesReportRows(
	authCustID string,
	custIDs []string,
	parentCustID string,
	fromDate, toDate string,
	salesmanIDs []int,
	distributorCodes []string,
	limit, offset int,
) ([]model.SalesActivityReportRow, error) {
	sql, args := buildActivitySalesReportCombinedSQL(
		repository.DB,
		authCustID,
		custIDs,
		parentCustID,
		fromDate,
		toDate,
		salesmanIDs,
		distributorCodes,
		limit,
		offset,
		false,
	)
	var rows []model.SalesActivityReportRow
	err := repository.Raw(sql, args...).Scan(&rows).Error
	return rows, err
}

func (repository *RepositoryReportImpl) countActivitySalesReportRows(
	authCustID string,
	custIDs []string,
	parentCustID string,
	fromDate, toDate string,
	salesmanIDs []int,
	distributorCodes []string,
) (int64, error) {
	sql, args := buildActivitySalesReportCombinedSQL(
		repository.DB,
		authCustID,
		custIDs,
		parentCustID,
		fromDate,
		toDate,
		salesmanIDs,
		distributorCodes,
		0,
		0,
		true,
	)
	var total int64
	err := repository.Raw(sql, args...).Scan(&total).Error
	return total, err
}
