-- ============================================
-- Rollback SX-2184 order_type and original PO qty fields
-- Created: 2026-06-08
-- Description: remove order_type from sls.order and original_qty_po* from sls.order_detail
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS original_qty_po1,
    DROP COLUMN IF EXISTS original_qty_po2,
    DROP COLUMN IF EXISTS original_qty_po3;

ALTER TABLE IF EXISTS sls.order
    DROP COLUMN IF EXISTS order_type;

COMMIT;
