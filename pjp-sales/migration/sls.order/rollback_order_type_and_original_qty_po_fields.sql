-- ============================================
-- Rollback Order Type and Original PO Qty Migration
-- Created: 2026-06-04
-- Description: Remove order type from sls.order and original purchase order quantity fields from sls.order_detail
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order
    DROP CONSTRAINT IF EXISTS order_order_type_check;

ALTER TABLE IF EXISTS sls.order
    DROP COLUMN IF EXISTS order_type;

ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS original_qty_po1,
    DROP COLUMN IF EXISTS original_qty_po2,
    DROP COLUMN IF EXISTS original_qty_po3;

COMMIT;
