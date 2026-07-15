-- ============================================
-- Rollback Order Reward Reff ID Length Migration
-- Created: 2026-06-23
-- Description: Revert sls.order_reward.reff_id back to varchar(20).
--              Only safe if all existing reff_id values are <= 20 chars.
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order_reward
    ALTER COLUMN reff_id TYPE VARCHAR(20);

COMMIT;
