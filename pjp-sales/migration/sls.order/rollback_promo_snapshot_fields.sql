-- ============================================
-- Rollback Order Promo Snapshot Header Fields Migration
-- Created: 2026-03-11
-- Description: Remove persisted promo snapshot remarks fields from sls.order table
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order
    DROP COLUMN IF EXISTS promo_remarks_po,
    DROP COLUMN IF EXISTS promo_remarks_final,
    DROP COLUMN IF EXISTS promo_remarks_so;

COMMIT;
