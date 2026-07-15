-- ============================================
-- Rollback Proforma Invoice Migration
-- Created: 2025-12-02
-- Description: Remove proforma invoice fields from sls.order table
-- ============================================

BEGIN;

-- Remove first_issue_date column
ALTER TABLE IF EXISTS sls.order
    DROP COLUMN IF EXISTS first_issue_date;

-- Remove generate_by column
ALTER TABLE IF EXISTS sls.order
    DROP COLUMN IF EXISTS generate_by;

-- Remove is_proforma_inv column
ALTER TABLE IF EXISTS sls.order
    DROP COLUMN IF EXISTS is_proforma_inv;

COMMIT;
