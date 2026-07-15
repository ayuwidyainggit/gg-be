-- SX-2426 production data repair
-- Purpose:
--   Repair duplicated CO release rows that made inv.stock.qty_in_order - qty_out_order negative
--   for invoiced orders where release (`CO` qty_out_order) was written more than reserve (`CO` qty_in_order).
--
-- Safety model:
--   1) Dry-run the candidate set first.
--   2) Apply by INSERTING compensation rows only (append-only audit trail).
--   3) Recompute inv.warehouse_stock.qty_on_order from inv.stock after compensation.
--
-- Scope:
--   Default example is narrowed to the Jira sample products at wh_id=566.
--   Remove or widen the WHERE clause only after review.

-- ============================================================
-- 1) DRY RUN: inspect offending CO rows
-- ============================================================
WITH co_balance AS (
    SELECT
        s.cust_id,
        s.wh_id,
        s.pro_id,
        s.ref_det_id,
        s.tr_no,
        MIN(s.stock_date) AS stock_date,
        MAX(s.unit_price) AS unit_price,
        SUM(COALESCE(s.qty_in_order, 0)) AS sum_qty_in_order,
        SUM(COALESCE(s.qty_out_order, 0)) AS sum_qty_out_order,
        SUM(COALESCE(s.qty_out_order, 0)) - SUM(COALESCE(s.qty_in_order, 0)) AS over_release_qty
    FROM inv.stock s
    WHERE s.tr_code = 'CO'
      AND s.wh_id = 566
      AND s.pro_id IN (9634, 9636, 9637)
    GROUP BY s.cust_id, s.wh_id, s.pro_id, s.ref_det_id, s.tr_no
)
SELECT *
FROM co_balance
WHERE over_release_qty > 0
ORDER BY cust_id, wh_id, pro_id, tr_no, ref_det_id;

-- ============================================================
-- 2) APPLY: append compensation rows to inv.stock
-- ============================================================
-- Recommendation: run inside a transaction and review RETURNING output first.
BEGIN;

WITH co_balance AS (
    SELECT
        s.cust_id,
        s.wh_id,
        s.pro_id,
        s.ref_det_id,
        s.tr_no,
        MIN(s.stock_date) AS stock_date,
        MAX(s.unit_price) AS unit_price,
        SUM(COALESCE(s.qty_in_order, 0)) AS sum_qty_in_order,
        SUM(COALESCE(s.qty_out_order, 0)) AS sum_qty_out_order,
        SUM(COALESCE(s.qty_out_order, 0)) - SUM(COALESCE(s.qty_in_order, 0)) AS over_release_qty
    FROM inv.stock s
    WHERE s.tr_code = 'CO'
      AND s.wh_id = 566
      AND s.pro_id IN (9634, 9636, 9637)
    GROUP BY s.cust_id, s.wh_id, s.pro_id, s.ref_det_id, s.tr_no
), candidates AS (
    SELECT *
    FROM co_balance
    WHERE over_release_qty > 0
), inserted AS (
    INSERT INTO inv.stock (
        cust_id,
        stock_date,
        tr_code,
        tr_no,
        wh_id,
        pro_id,
        item_cdn,
        qty_in,
        qty_out,
        unit_price,
        cogs,
        ref_det_id,
        created_at,
        qty_in_order,
        qty_out_order
    )
    SELECT
        c.cust_id,
        c.stock_date,
        'CO' AS tr_code,
        c.tr_no,
        c.wh_id,
        c.pro_id,
        1 AS item_cdn,
        0 AS qty_in,
        0 AS qty_out,
        c.unit_price,
        0 AS cogs,
        c.ref_det_id,
        EXTRACT(EPOCH FROM NOW())::bigint AS created_at,
        c.over_release_qty AS qty_in_order,
        0 AS qty_out_order
    FROM candidates c
    RETURNING cust_id, wh_id, pro_id, ref_det_id, tr_no, qty_in_order
), wh_recalc AS (
    SELECT
        s.cust_id,
        s.wh_id,
        s.pro_id,
        COALESCE(SUM(s.qty_in_order), 0) - COALESCE(SUM(s.qty_out_order), 0) AS qty_on_order_new
    FROM inv.stock s
    WHERE s.wh_id = 566
      AND s.pro_id IN (9634, 9636, 9637)
    GROUP BY s.cust_id, s.wh_id, s.pro_id
)
UPDATE inv.warehouse_stock ws
SET qty_on_order = r.qty_on_order_new,
    updated_at = EXTRACT(EPOCH FROM NOW())::bigint
FROM wh_recalc r
WHERE ws.cust_id = r.cust_id
  AND ws.wh_id = r.wh_id
  AND ws.pro_id = r.pro_id;

-- Review before commit:
SELECT * FROM inserted ORDER BY cust_id, wh_id, pro_id, tr_no, ref_det_id;

SELECT
    ws.cust_id,
    ws.wh_id,
    ws.pro_id,
    ws.qty_on_order,
    COALESCE(SUM(s.qty_in_order), 0) - COALESCE(SUM(s.qty_out_order), 0) AS ledger_qty_on_order
FROM inv.warehouse_stock ws
LEFT JOIN inv.stock s
    ON s.cust_id = ws.cust_id
   AND s.wh_id = ws.wh_id
   AND s.pro_id = ws.pro_id
WHERE ws.wh_id = 566
  AND ws.pro_id IN (9634, 9636, 9637)
GROUP BY ws.cust_id, ws.wh_id, ws.pro_id, ws.qty_on_order
ORDER BY ws.cust_id, ws.pro_id;

COMMIT;

-- If verification fails, use ROLLBACK instead of COMMIT.
