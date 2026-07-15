package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"sales/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBuildCancelStockMutations_SingleSKU(t *testing.T) {
	stockDate := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)
	baseRows := []cancelStockBaseRow{
		{
			CustID:       "C220010001",
			WhID:         63,
			ProID:        474,
			RefDetID:     4949,
			QtyOutSO:     2400,
			UnitPrice:    6432,
			SourceTrNo:   "SO2602090002",
			SourceTrCode: "SO",
		},
	}

	stocks, whDeltas := buildCancelStockMutations("SO2602090002", stockDate, baseRows)
	if len(stocks) != 1 {
		t.Fatalf("expected 1 stock row (single CO reversal), got %d", len(stocks))
	}
	if len(whDeltas) != 1 {
		t.Fatalf("expected 1 warehouse delta row, got %d", len(whDeltas))
	}

	rowB := stocks[0]
	if rowB.TrNo != "SO2602090002-CO" || rowB.TrCode != "CO" {
		t.Fatalf("unexpected reversal row identity: tr_no=%s tr_code=%s", rowB.TrNo, rowB.TrCode)
	}
	if rowB.QtyIn != 0 || rowB.QtyOut != 0 || rowB.QtyInOrder != 0 || rowB.QtyOutOrder != 2400 {
		t.Fatalf("unexpected reversal row quantities: %+v", rowB)
	}

	delta := whDeltas[0]
	if delta.Qty != 2400 || delta.QtyOnOrder != -2400 {
		t.Fatalf("unexpected warehouse delta: %+v", delta)
	}
}

func TestSalesStockUpdates_AllZeroDelta_ShouldReturnWithoutDBWrites(t *testing.T) {
	repo := &RepositoryStockImpl{}
	stockDate := time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)
	before := 10.0

	updates := []*entity.SalesOrderStockUpdate{
		{
			CustID:         "C220010001",
			WhID:           63,
			ProID:          474,
			StockDate:      stockDate,
			TrCode:         "SO",
			TrNo:           "SO2603260001",
			RefDetId:       1001,
			QtyOrder:       10,
			QtyOrderBefore: &before,
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected zero-delta updates to return before DB writes, panic=%v", r)
		}
	}()

	if err := repo.SalesStockUpdates(context.Background(), updates); err != nil {
		t.Fatalf("expected nil error for zero-delta updates, got %v", err)
	}
}

func TestSalesStockUpdates_RealDelta_ShouldStillReachWritePath(t *testing.T) {
	repo := &RepositoryStockImpl{}
	stockDate := time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)
	before := 10.0

	updates := []*entity.SalesOrderStockUpdate{
		{
			CustID:         "C220010001",
			WhID:           63,
			ProID:          474,
			StockDate:      stockDate,
			TrCode:         "SO",
			TrNo:           "SO2603260002",
			RefDetId:       1001,
			QtyOrder:       12,
			QtyOrderBefore: &before,
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected real delta to continue into write path")
		}
	}()

	_ = repo.SalesStockUpdates(context.Background(), updates)
}

func TestBuildInvoiceReleaseKeys_MixedCustIDRejected(t *testing.T) {
	updates := []*entity.InvoiceSalesStockUpdate{
		{CustID: "C1", TrNo: "SO1", RefDetId: 10},
		{CustID: "C2", TrNo: "SO1", RefDetId: 11},
	}

	_, _, err := buildInvoiceReleaseKeys(updates)
	if err == nil {
		t.Fatalf("expected mixed cust_id to be rejected")
	}
}

func TestFilterDuplicateInvoiceStockUpdates_SkipsExistingCORelease(t *testing.T) {
	updates := []*entity.InvoiceSalesStockUpdate{
		{CustID: "C220010001", TrNo: "SO2606190001", RefDetId: 23780},
		{CustID: "C220010001", TrNo: "SO2606190001", RefDetId: 23781},
	}
	keys, _, err := buildInvoiceReleaseKeys(updates)
	if err != nil {
		t.Fatalf("buildInvoiceReleaseKeys error: %v", err)
	}

	existing := map[invoiceReleaseKey]struct{}{
		{CustID: "C220010001", TrNo: "SO2606190001-CO", RefDetID: 23780}: {},
	}

	filteredUpdates, filteredKeys, skipped := filterDuplicateInvoiceStockUpdates(updates, keys, existing)
	if skipped != 1 {
		t.Fatalf("expected 1 skipped duplicate, got %d", skipped)
	}
	if len(filteredUpdates) != 1 || len(filteredKeys) != 1 {
		t.Fatalf("expected only one surviving invoice release row, got updates=%d keys=%d", len(filteredUpdates), len(filteredKeys))
	}
	if filteredUpdates[0].RefDetId != 23781 {
		t.Fatalf("expected ref_det_id 23781 to survive, got %d", filteredUpdates[0].RefDetId)
	}
	if filteredKeys[0].TrNo != "SO2606190001-CO" {
		t.Fatalf("expected normalized CO tr_no, got %s", filteredKeys[0].TrNo)
	}
}

func TestBuildCancelStockMutations_MultiSKUAndIdempotent(t *testing.T) {
	stockDate := time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC)
	baseRows := []cancelStockBaseRow{
		{CustID: "C220010001", WhID: 63, ProID: 474, RefDetID: 1001, QtyOutSO: 2400, UnitPrice: 6432, SourceTrNo: "SO2602090002", SourceTrCode: "SO"},
		{CustID: "C220010001", WhID: 63, ProID: 475, RefDetID: 1002, QtyOutSO: 1200, UnitPrice: 5000, SourceTrNo: "SO2602090002", SourceTrCode: "SO"},
	}

	stocks, whDeltas := buildCancelStockMutations("SO2602090002", stockDate, baseRows)

	if len(stocks) != 2 {
		t.Fatalf("expected 2 stock rows (one CO reversal per SKU), got %d", len(stocks))
	}
	if len(whDeltas) != 2 {
		t.Fatalf("expected 2 warehouse delta rows for 2 active ref_det_id, got %d", len(whDeltas))
	}
}

func TestBuildCancelStockMutations_PartialReverseExistingRefStillBuildsOutstanding(t *testing.T) {
	stockDate := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	baseRows := []cancelStockBaseRow{
		{
			CustID:       "C220010001",
			WhID:         63,
			ProID:        474,
			RefDetID:     1001,
			QtyOutSO:     12,
			UnitPrice:    6432,
			SourceTrNo:   "SO2603180011",
			SourceTrCode: "SO",
		},
	}

	stocks, whDeltas := buildCancelStockMutations("SO2603180011", stockDate, baseRows)
	if len(stocks) != 1 {
		t.Fatalf("expected outstanding reservation to still be reversed on final cancel with 1 CO row, got %d stock rows", len(stocks))
	}
	if len(whDeltas) != 1 {
		t.Fatalf("expected 1 warehouse delta for remaining outstanding reservation, got %d", len(whDeltas))
	}
	if whDeltas[0].QtyOnOrder != -12 {
		t.Fatalf("expected qty_on_order reverse to clear outstanding 12, got %+v", whDeltas[0].QtyOnOrder)
	}
}

func TestGetCancelStockBasisQuery_UsesOutstandingReservationFormula(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres dbname=postgres sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{DisableAutomaticPing: true, DryRun: true})
	if err != nil {
		t.Fatalf("failed to open dry-run gorm db: %v", err)
	}

	repo := NewStockRepo(db)
	stmt := repo.cancelStockBasisQuery(context.Background(), "C220010001", "SO2603180011").Find(&[]struct{}{}).Statement
	sql := stmt.SQL.String()

	if !strings.Contains(sql, "SUM(s.qty_out - s.qty_in) AS qty_out_so") {
		t.Fatalf("expected source subquery to aggregate outstanding reservation, sql=%s", sql)
	}
	if !strings.Contains(sql, "SUM(c.qty_out_order) AS qty_out_order_cancel") {
		t.Fatalf("expected cancel subquery to aggregate previous reverse rows, sql=%s", sql)
	}
	if !strings.Contains(sql, "c.tr_code = 'CO'") {
		t.Fatalf("expected cancel subquery to detect prior CO reversal rows, sql=%s", sql)
	}
	if !strings.Contains(sql, "FROM sls.order_detail od") {
		t.Fatalf("expected basis query to be rooted on final active order detail rows, sql=%s", sql)
	}
	if !strings.Contains(sql, "JOIN sls.order o ON o.cust_id = od.cust_id AND o.ro_no = od.ro_no") {
		t.Fatalf("expected basis query to resolve warehouse from order header, sql=%s", sql)
	}
	if !strings.Contains(sql, "od.order_detail_id AS ref_det_id") {
		t.Fatalf("expected basis query to reconcile by order_detail_id/ref_det_id, sql=%s", sql)
	}
	if !strings.Contains(sql, "(COALESCE(s.qty_out_so, 0) <= 0) AS is_missing_source") {
		t.Fatalf("expected basis query to expose missing source signal, sql=%s", sql)
	}
	if !strings.Contains(sql, "LEFT JOIN") || !strings.Contains(sql, ") ad") {
		t.Fatalf("expected basis query to join active detail aggregate as ad, sql=%s", sql)
	}
	if !strings.Contains(sql, "AND ad.ro_no = od.ro_no") || !strings.Contains(sql, "AND ad.pro_id = od.pro_id") {
		t.Fatalf("expected active detail aggregate to join by ro_no and pro_id, sql=%s", sql)
	}
	if !strings.Contains(sql, "(COALESCE(ad.active_detail_count, 0) > 1) AS is_ambiguous") {
		t.Fatalf("expected basis query to expose ambiguous split signal, sql=%s", sql)
	}
	// SX-2314: cancelAgg must include legacy rows (tr_code='SO' tr_no LIKE '%-CO%')
	if !strings.Contains(sql, "c.tr_code = 'CO'") {
		t.Fatalf("expected cancelAgg to match tr_code=CO companion reversal rows, sql=%s", sql)
	}
	if !strings.Contains(sql, "c.tr_code = 'SO'") {
		t.Fatalf("expected cancelAgg to also match legacy tr_code=SO tr_no suffix rows, sql=%s", sql)
	}
	if !strings.Contains(sql, "c.tr_no LIKE '%-CO%'") {
		t.Fatalf("expected cancelAgg to filter legacy rows by tr_no suffix, sql=%s", sql)
	}
}

func TestGetCancelStockBasisQuery_IncludesRewardLine(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres dbname=postgres sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{DisableAutomaticPing: true, DryRun: true})
	if err != nil {
		t.Fatalf("failed to open dry-run gorm db: %v", err)
	}

	repo := NewStockRepo(db)
	stmt := repo.cancelStockBasisQuery(context.Background(), "C220010001", "SO2603180011").Find(&[]struct{}{}).Statement
	sql := stmt.SQL.String()

	if !strings.Contains(sql, "AS qty_outstanding") {
		t.Fatalf("expected main SELECT to define qty_outstanding, sql=%s", sql)
	}
	if !strings.Contains(sql, "active_detail_count") {
		t.Fatalf("expected activeDetailAgg subquery, sql=%s", sql)
	}
	count := strings.Count(sql, "od.item_type = 1")
	if count > 1 {
		t.Fatalf("expected main SELECT to not filter item_type=1 (only activeDetailAgg should), got %d occurrences of 'od.item_type = 1', sql=%s", count, sql)
	}
}

func TestGetCancelStockBasisQuery_LegacySORowsIncludedInCancelAgg(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres dbname=postgres sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{DisableAutomaticPing: true, DryRun: true})
	if err != nil {
		t.Fatalf("failed to open dry-run gorm db: %v", err)
	}

	repo := NewStockRepo(db)
	stmt := repo.cancelStockBasisQuery(context.Background(), "C220010001", "SO2603180011").Find(&[]struct{}{}).Statement
	sql := stmt.SQL.String()

	// cancelAgg must match both new tr_code='CO' reversal rows AND legacy tr_code='SO' rows
	// where tr_no LIKE '%-CO%'. The combined predicate is the new contract.
	if !strings.Contains(sql, "c.tr_no = ? AND (c.tr_code = 'CO' OR (c.tr_code = 'SO' AND c.tr_no LIKE '%-CO%'))") {
		// Fallback: accept the OR-equivalent form that some formatters produce
		if !strings.Contains(sql, "c.tr_code = 'CO' OR (c.tr_code = 'SO'") {
			t.Fatalf("expected cancelAgg to match legacy tr_code=SO tr_no suffix rows via OR, sql=%s", sql)
		}
	}
}

func TestGetCancelStockBasisQuery_POQtyFallback(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres dbname=postgres sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{DisableAutomaticPing: true, DryRun: true})
	if err != nil {
		t.Fatalf("failed to open dry-run gorm db: %v", err)
	}

	repo := NewStockRepo(db)
	stmt := repo.cancelStockBasisQuery(context.Background(), "C220010001", "SO2603180011").Find(&[]struct{}{}).Statement
	sql := stmt.SQL.String()

	if !strings.Contains(sql, "COALESCE(od.qty1_final, 0) > 0") {
		t.Fatalf("expected qty1_final in active detail filter, sql=%s", sql)
	}
	if !strings.Contains(sql, "COALESCE(od.qty1, 0) > 0") {
		t.Fatalf("expected qty1 fallback in active detail filter, sql=%s", sql)
	}
	if !strings.Contains(sql, "COALESCE(od.qty_po1, 0) > 0") {
		t.Fatalf("expected qty_po1 fallback in active detail filter, sql=%s", sql)
	}
	if !strings.Contains(sql, "COALESCE(od.qty1_final, COALESCE(od.qty1, COALESCE(od.qty_po1, 0))) AS qty_final") {
		t.Fatalf("expected qty priority select for qty1/final/po, sql=%s", sql)
	}
	if !strings.Contains(sql, "COALESCE(s.unit_price, COALESCE(od.sell_price_final1, COALESCE(od.sell_price1, COALESCE(od.sell_price_po1, 0)))) AS unit_price") {
		t.Fatalf("expected sell price fallback to po price, sql=%s", sql)
	}
	// Detail qty smallest-unit fallback expression (final → sales → PO with conv units)
	if !strings.Contains(sql, "qty1_final") || !strings.Contains(sql, "COALESCE(od.qty1,") || !strings.Contains(sql, "COALESCE(od.qty_po1,") {
		t.Fatalf("expected detail qty final/sales/po fallback with conv units, sql=%s", sql)
	}
	if !strings.Contains(sql, "conv_unit2") || !strings.Contains(sql, "conv_unit3") {
		t.Fatalf("expected conv_unit2/conv_unit3 in smallest-unit fallback expression, sql=%s", sql)
	}
	if !strings.Contains(sql, "GREATEST") {
		t.Fatalf("expected GREATEST for smallest-unit fallback qty, sql=%s", sql)
	}
}

func TestCancelStockBasisFallback_DetailOnlyPOCase(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres dbname=postgres sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{DisableAutomaticPing: true, DryRun: true})
	if err != nil {
		t.Fatalf("failed to open dry-run gorm db: %v", err)
	}

	repo := NewStockRepo(db)
	stmt := repo.cancelStockBasisQuery(context.Background(), "C220010001", "SO2603180011").Find(&[]struct{}{}).Statement
	sql := stmt.SQL.String()

	// Verify the COALESCE chain: qty1_final → qty1 → qty_po1
	if !strings.Contains(sql, "COALESCE(od.qty1_final, COALESCE(od.qty1, COALESCE(od.qty_po1, 0))) AS qty_final") {
		t.Fatalf("expected qty1 priority COALESCE chain with qty_po1 fallback, sql=%s", sql)
	}
	// Verify smallest-unit conversion uses conv units
	if !strings.Contains(sql, "conv_unit2") || !strings.Contains(sql, "conv_unit3") {
		t.Fatalf("expected conv_unit2/conv_unit3 for smallest-unit conversion, sql=%s", sql)
	}
	// Verify the detail qty fallback is used as source qty when ledger missing
	if !strings.Contains(sql, "GREATEST") || !strings.Contains(sql, "qty_final") {
		t.Fatalf("expected GREATEST wrapper around detail qty fallback, sql=%s", sql)
	}
}
