package repository

import (
	"fmt"
	"strings"
)

func BuildInTransitStockSubqueries(isPrincipal bool, custId string) (string, string, string) {
	var inTransitStock1Subquery, inTransitStock2Subquery, inTransitStock3Subquery string

	if isPrincipal {
		// For principal: filter by rod.cust_id = ro.cust_id (already in JOIN condition)
		inTransitStock1Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod.order_booking_qty1, 0))
			FROM inv.replenishment_order_detail rod
			JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
			WHERE rod.pro_id = p.pro_id 
				AND rod.cust_id = ro.cust_id
				AND ro.status = 4
				AND rod.is_del = false
				AND ro.is_del = false
		), 0)`
		inTransitStock2Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod.order_booking_qty2, 0))
			FROM inv.replenishment_order_detail rod
			JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
			WHERE rod.pro_id = p.pro_id 
				AND rod.cust_id = ro.cust_id
				AND ro.status = 4
				AND rod.is_del = false
				AND ro.is_del = false
		), 0)`
		inTransitStock3Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod.order_booking_qty3, 0))
			FROM inv.replenishment_order_detail rod
			JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
			WHERE rod.pro_id = p.pro_id 
				AND rod.cust_id = ro.cust_id
				AND ro.status = 4
				AND rod.is_del = false
				AND ro.is_del = false
		), 0)`
	} else {
		// For non-principal: filter by custId from JWT token
		escapedCustId := strings.ReplaceAll(custId, "'", "''")
		inTransitStock1Subquery = fmt.Sprintf(`COALESCE((
			SELECT SUM(COALESCE(rod.order_booking_qty1, 0))
			FROM inv.replenishment_order_detail rod
			JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
			WHERE rod.pro_id = p.pro_id 
				AND rod.cust_id = '%s'
				AND ro.status = 4
				AND rod.is_del = false
				AND ro.is_del = false
		), 0)`, escapedCustId)
		inTransitStock2Subquery = fmt.Sprintf(`COALESCE((
			SELECT SUM(COALESCE(rod.order_booking_qty2, 0))
			FROM inv.replenishment_order_detail rod
			JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
			WHERE rod.pro_id = p.pro_id 
				AND rod.cust_id = '%s'
				AND ro.status = 4
				AND rod.is_del = false
				AND ro.is_del = false
		), 0)`, escapedCustId)
		inTransitStock3Subquery = fmt.Sprintf(`COALESCE((
			SELECT SUM(COALESCE(rod.order_booking_qty3, 0))
			FROM inv.replenishment_order_detail rod
			JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
			WHERE rod.pro_id = p.pro_id 
				AND rod.cust_id = '%s'
				AND ro.status = 4
				AND rod.is_del = false
				AND ro.is_del = false
		), 0)`, escapedCustId)
	}

	return inTransitStock1Subquery, inTransitStock2Subquery, inTransitStock3Subquery
}

func BuildWarehouseStockTotalSubquery(isPrincipal bool, custId string, whId int64) string {
	var whsTotalQtySubquery string
	warehouseFilter := ""
	if whId > 0 {
		warehouseFilter = fmt.Sprintf(" AND whs_sub.wh_id = %d", whId)
	}

	if isPrincipal {
		whsTotalQtySubquery = fmt.Sprintf(`COALESCE((
			SELECT SUM(COALESCE(whs_sub.qty, 0))
			FROM inv.warehouse_stock whs_sub
			WHERE whs_sub.pro_id = p.pro_id%s
		), 0)`, warehouseFilter)
	} else {
		escapedCustId := strings.ReplaceAll(custId, "'", "''")
		whsTotalQtySubquery = fmt.Sprintf(`COALESCE((
			SELECT SUM(COALESCE(whs_sub.qty, 0))
			FROM inv.warehouse_stock whs_sub
			WHERE whs_sub.pro_id = p.pro_id 
				AND whs_sub.cust_id = '%s'
				%s
		), 0)`, escapedCustId, warehouseFilter)
	}

	return whsTotalQtySubquery
}

func BuildQtyCalculationExpressions(whsTotalQtySubquery, tableAlias string, columnAliases ...string) (string, string, string) {
	a1, a2, a3 := "qty1", "qty2", "qty3"
	if len(columnAliases) >= 3 {
		a1, a2, a3 = columnAliases[0], columnAliases[1], columnAliases[2]
	}

	qty3Expression := fmt.Sprintf(`CASE 
		WHEN %s = 0 OR %s.conv_unit2 = 0 OR %s.conv_unit3 = 0 THEN 0
		ELSE FLOOR(%s / (%s.conv_unit2 * %s.conv_unit3))
	END AS %s`, whsTotalQtySubquery, tableAlias, tableAlias, whsTotalQtySubquery, tableAlias, tableAlias, a3)

	qty2Expression := fmt.Sprintf(`CASE 
		WHEN %s = 0 OR %s.conv_unit2 = 0 THEN 0
		ELSE FLOOR(
			(
				%s - 
				(FLOOR(%s / NULLIF(%s.conv_unit2 * %s.conv_unit3, 0)) * %s.conv_unit2 * %s.conv_unit3)
			) / NULLIF(%s.conv_unit2, 0)
		)
	END AS %s`, whsTotalQtySubquery, tableAlias, whsTotalQtySubquery, whsTotalQtySubquery, tableAlias, tableAlias, tableAlias, tableAlias, tableAlias, a2)

	qty1Expression := fmt.Sprintf(`CASE 
		WHEN %s = 0 THEN 0
		ELSE 
			%s -
			(FLOOR(%s / NULLIF(%s.conv_unit2 * %s.conv_unit3, 0)) * %s.conv_unit2 * %s.conv_unit3) -
			(FLOOR(
				(
					%s - 
					(FLOOR(%s / NULLIF(%s.conv_unit2 * %s.conv_unit3, 0)) * %s.conv_unit2 * %s.conv_unit3)
				) / NULLIF(%s.conv_unit2, 0)
			) * %s.conv_unit2)
	END AS %s`, whsTotalQtySubquery, whsTotalQtySubquery, whsTotalQtySubquery, tableAlias, tableAlias, tableAlias, tableAlias, whsTotalQtySubquery, whsTotalQtySubquery, tableAlias, tableAlias, tableAlias, tableAlias, tableAlias, tableAlias, a1)

	return qty1Expression, qty2Expression, qty3Expression
}
