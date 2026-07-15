-- ============================================
-- Order Detail Enhancement: disc_value_po & vat_po Migration
-- Created: 2026-01-26
-- Description: Add disc_value_po and vat_po columns to sls.order_detail table
--              Required for Enhance Sales Order feature (Taking Order & Canvas)
-- Reference: Enhance Sales Order - BE.md
-- ============================================

BEGIN;

-- Add disc_value_po and vat_po columns
ALTER TABLE IF EXISTS sls.order_detail
    ADD COLUMN IF NOT EXISTS disc_value_po NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS vat_po NUMERIC(20,4) DEFAULT 0;

-- Add comments for documentation
COMMENT ON COLUMN sls.order_detail.disc_value_po IS 'Discount value for Purchase Order - mapped from disc_value in request';
COMMENT ON COLUMN sls.order_detail.vat_po IS 'VAT percentage for Purchase Order - mapped from vat in request';

-- Backfill existing data: copy disc_value to disc_value_po, vat to vat_po
UPDATE sls.order_detail 
SET 
    disc_value_po = COALESCE(disc_value, 0),
    vat_po = COALESCE(vat, 0)
WHERE disc_value_po = 0 OR disc_value_po IS NULL
   OR vat_po = 0 OR vat_po IS NULL;

COMMIT;