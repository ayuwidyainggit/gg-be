-- ============================================
-- Rollback Add opr_type Column Migration
-- Created: 2026-03-31
-- Description: Remove opr_type column from sls.order table
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order
    DROP COLUMN IF EXISTS opr_type;

COMMIT;
