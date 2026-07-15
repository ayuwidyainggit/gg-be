-- ============================================
-- Order Reward Reff ID Length Migration
-- Created: 2026-06-23
-- Description: Widen sls.order_reward.reff_id from varchar(20) to varchar(50)
--              to keep the promo reference width consistent with
--              promo.promotions.promo_id and sls.order_detail.promo_id.
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order_reward
    ALTER COLUMN reff_id TYPE VARCHAR(50);

COMMIT;
