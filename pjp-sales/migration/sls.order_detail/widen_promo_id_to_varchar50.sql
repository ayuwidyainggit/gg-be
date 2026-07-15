-- ============================================
-- Order Detail Promo ID Length Migration
-- Created: 2026-06-23
-- Description: Widen sls.order_detail.promo_id from varchar(20) to varchar(50)
--              to match Promo V2 master schema (promo.promotions.promo_id varchar(50)).
--              Without this, reward detail rows generated from Promo V2 product
--              rewards fail insert with "value too long for type varying(20)".
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order_detail
    ALTER COLUMN promo_id TYPE VARCHAR(50);

COMMIT;
