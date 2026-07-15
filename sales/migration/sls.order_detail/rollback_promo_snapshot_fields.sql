-- ============================================
-- Rollback Order Detail Promo Snapshot Fields Migration
-- Created: 2026-03-11
-- Description: Remove persisted promo snapshot fields from sls.order_detail table
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS is_product_promotion_po,
    DROP COLUMN IF EXISTS is_product_promotion_final,
    DROP COLUMN IF EXISTS is_product_promotion_so,
    DROP COLUMN IF EXISTS promo_remarks_po,
    DROP COLUMN IF EXISTS promo_remarks_final,
    DROP COLUMN IF EXISTS promo_remarks_so,
    DROP COLUMN IF EXISTS promo_po5,
    DROP COLUMN IF EXISTS promo_po4,
    DROP COLUMN IF EXISTS promo_po3,
    DROP COLUMN IF EXISTS promo_po2,
    DROP COLUMN IF EXISTS promo_po1,
    DROP COLUMN IF EXISTS promo_final5,
    DROP COLUMN IF EXISTS promo_final4,
    DROP COLUMN IF EXISTS promo_final3,
    DROP COLUMN IF EXISTS promo_final2,
    DROP COLUMN IF EXISTS promo_final1,
    DROP COLUMN IF EXISTS promo_so5,
    DROP COLUMN IF EXISTS promo_so4,
    DROP COLUMN IF EXISTS promo_so3,
    DROP COLUMN IF EXISTS promo_so2,
    DROP COLUMN IF EXISTS promo_so1;

COMMIT;
