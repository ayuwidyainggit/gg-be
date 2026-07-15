-- ============================================
-- Rollback Order Detail PO Fields Migration
-- Created: 2025-12-04
-- Description: Remove purchase order and final price fields from sls.order_detail table
-- ============================================

BEGIN;

-- Remove sell_price_final columns
ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS sell_price_final3,
    DROP COLUMN IF EXISTS sell_price_final2,
    DROP COLUMN IF EXISTS sell_price_final1;

-- Remove sell_price_po columns
ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS sell_price_po3,
    DROP COLUMN IF EXISTS sell_price_po2,
    DROP COLUMN IF EXISTS sell_price_po1;

-- Remove qty_po columns
ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS qty_po3,
    DROP COLUMN IF EXISTS qty_po2,
    DROP COLUMN IF EXISTS qty_po1;

COMMIT;
