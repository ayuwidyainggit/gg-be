-- ============================================
-- Proforma Invoice Migration
-- Created: 2025-12-02
-- Description: Add proforma invoice fields to sls.order table
-- ============================================

BEGIN;

-- Add is_proforma_inv column
-- false/null: belum generate proforma invoice
-- true: sudah generate proforma invoice
ALTER TABLE IF EXISTS sls.order
    ADD COLUMN IF NOT EXISTS is_proforma_inv BOOLEAN;

-- Add generate_by column
-- ID user yang pertama kali melakukan generate proforma invoice
ALTER TABLE IF EXISTS sls.order
    ADD COLUMN IF NOT EXISTS generate_by BIGINT;

-- Add first_issue_date column
-- Waktu dan tanggal generate proforma invoice pertama kali
ALTER TABLE IF EXISTS sls.order
    ADD COLUMN IF NOT EXISTS first_issue_date TIMESTAMPTZ;

-- Add comments for documentation
COMMENT ON COLUMN sls.order.is_proforma_inv IS 'false/null: belum generate proforma invoice; true: sudah generate proforma invoice';
COMMENT ON COLUMN sls.order.generate_by IS 'ID user yang pertama kali melakukan generate proforma invoice';
COMMENT ON COLUMN sls.order.first_issue_date IS 'Waktu dan tanggal generate proforma invoice pertama kali';

COMMIT;
