-- ============================================
-- Rollback Order Detail Promo ID Length Migration
-- Created: 2026-06-23
-- Description: Revert sls.order_detail.promo_id back to varchar(20).
--              Only safe if all existing promo_id values are <= 20 chars.
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order_detail
    ALTER COLUMN promo_id TYPE VARCHAR(20);

COMMIT;
