-- ============================================
-- Order Detail Discount & VAT PO Fields Migration
-- Created: 2026-01-22
-- Description: Add discount and VAT fields for Purchase Order tab to sls.order_detail table
-- ============================================

BEGIN;

-- Add disc_po and vat_value_po columns
ALTER TABLE IF EXISTS sls.order_detail
    ADD COLUMN IF NOT EXISTS disc_po NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS vat_value_po NUMERIC(20,4) DEFAULT 0;

-- Add comments for documentation
COMMENT ON COLUMN sls.order_detail.disc_po IS 'Discount value for Purchase Order tab';
COMMENT ON COLUMN sls.order_detail.vat_value_po IS 'VAT value for Purchase Order tab';

-- Backfill existing data: copy disc_value to disc_po, vat_value to vat_value_po
UPDATE sls.order_detail 
SET 
    disc_po = COALESCE(disc_value, 0),
    vat_value_po = COALESCE(vat_value, 0)
WHERE disc_po = 0 OR disc_po IS NULL;

COMMIT;
