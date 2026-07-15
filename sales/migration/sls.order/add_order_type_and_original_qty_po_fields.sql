-- ============================================
-- SX-2184 order_type and original PO qty fields
-- Created: 2026-06-08
-- Description: add order_type to sls.order and original_qty_po* to sls.order_detail
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order
    ADD COLUMN IF NOT EXISTS order_type VARCHAR(2) NULL;

COMMENT ON COLUMN sls.order.order_type IS 'Order type snapshot for create order flow (O/C/SO)';

ALTER TABLE IF EXISTS sls.order_detail
    ADD COLUMN IF NOT EXISTS original_qty_po1 FLOAT4 NULL,
    ADD COLUMN IF NOT EXISTS original_qty_po2 FLOAT4 NULL,
    ADD COLUMN IF NOT EXISTS original_qty_po3 FLOAT4 NULL;

COMMENT ON COLUMN sls.order_detail.original_qty_po1 IS 'Original purchase order qty tier 1 from create request';
COMMENT ON COLUMN sls.order_detail.original_qty_po2 IS 'Original purchase order qty tier 2 from create request';
COMMENT ON COLUMN sls.order_detail.original_qty_po3 IS 'Original purchase order qty tier 3 from create request';

COMMIT;
