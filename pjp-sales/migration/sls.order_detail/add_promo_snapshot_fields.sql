-- ============================================
-- Order Detail Promo Snapshot Fields Migration
-- Created: 2026-03-11
-- Description: Add persisted promo snapshot fields to sls.order_detail table
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order_detail
    ADD COLUMN IF NOT EXISTS promo_so1 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_so2 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_so3 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_so4 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_so5 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_final1 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_final2 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_final3 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_final4 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_final5 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_po1 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_po2 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_po3 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_po4 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_po5 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS promo_remarks_so JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS promo_remarks_final JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS promo_remarks_po JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS is_product_promotion_so BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_product_promotion_final BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_product_promotion_po BOOLEAN DEFAULT FALSE;

COMMENT ON COLUMN sls.order_detail.promo_so1 IS 'Persisted promo snapshot amount 1 for Sales Order tab';
COMMENT ON COLUMN sls.order_detail.promo_so2 IS 'Persisted promo snapshot amount 2 for Sales Order tab';
COMMENT ON COLUMN sls.order_detail.promo_so3 IS 'Persisted promo snapshot amount 3 for Sales Order tab';
COMMENT ON COLUMN sls.order_detail.promo_so4 IS 'Persisted promo snapshot amount 4 for Sales Order tab';
COMMENT ON COLUMN sls.order_detail.promo_so5 IS 'Persisted promo snapshot amount 5 for Sales Order tab';
COMMENT ON COLUMN sls.order_detail.promo_final1 IS 'Persisted promo snapshot amount 1 for Final Order tab';
COMMENT ON COLUMN sls.order_detail.promo_final2 IS 'Persisted promo snapshot amount 2 for Final Order tab';
COMMENT ON COLUMN sls.order_detail.promo_final3 IS 'Persisted promo snapshot amount 3 for Final Order tab';
COMMENT ON COLUMN sls.order_detail.promo_final4 IS 'Persisted promo snapshot amount 4 for Final Order tab';
COMMENT ON COLUMN sls.order_detail.promo_final5 IS 'Persisted promo snapshot amount 5 for Final Order tab';
COMMENT ON COLUMN sls.order_detail.promo_po1 IS 'Persisted promo snapshot amount 1 for Purchase Order tab';
COMMENT ON COLUMN sls.order_detail.promo_po2 IS 'Persisted promo snapshot amount 2 for Purchase Order tab';
COMMENT ON COLUMN sls.order_detail.promo_po3 IS 'Persisted promo snapshot amount 3 for Purchase Order tab';
COMMENT ON COLUMN sls.order_detail.promo_po4 IS 'Persisted promo snapshot amount 4 for Purchase Order tab';
COMMENT ON COLUMN sls.order_detail.promo_po5 IS 'Persisted promo snapshot amount 5 for Purchase Order tab';
COMMENT ON COLUMN sls.order_detail.promo_remarks_so IS 'Persisted promo snapshot remarks for Sales Order tab';
COMMENT ON COLUMN sls.order_detail.promo_remarks_final IS 'Persisted promo snapshot remarks for Final Order tab';
COMMENT ON COLUMN sls.order_detail.promo_remarks_po IS 'Persisted promo snapshot remarks for Purchase Order tab';
COMMENT ON COLUMN sls.order_detail.is_product_promotion_so IS 'Persisted reward-product promo flag for Sales Order tab';
COMMENT ON COLUMN sls.order_detail.is_product_promotion_final IS 'Persisted reward-product promo flag for Final Order tab';
COMMENT ON COLUMN sls.order_detail.is_product_promotion_po IS 'Persisted reward-product promo flag for Purchase Order tab';

UPDATE sls.order_detail
SET
    promo_so1 = COALESCE(promo_so1, 0),
    promo_so2 = COALESCE(promo_so2, 0),
    promo_so3 = COALESCE(promo_so3, 0),
    promo_so4 = COALESCE(promo_so4, 0),
    promo_so5 = COALESCE(promo_so5, 0),
    promo_final1 = COALESCE(promo_final1, 0),
    promo_final2 = COALESCE(promo_final2, 0),
    promo_final3 = COALESCE(promo_final3, 0),
    promo_final4 = COALESCE(promo_final4, 0),
    promo_final5 = COALESCE(promo_final5, 0),
    promo_po1 = COALESCE(promo_po1, 0),
    promo_po2 = COALESCE(promo_po2, 0),
    promo_po3 = COALESCE(promo_po3, 0),
    promo_po4 = COALESCE(promo_po4, 0),
    promo_po5 = COALESCE(promo_po5, 0),
    promo_remarks_so = COALESCE(promo_remarks_so, '[]'::jsonb),
    promo_remarks_final = COALESCE(promo_remarks_final, '[]'::jsonb),
    promo_remarks_po = COALESCE(promo_remarks_po, '[]'::jsonb),
    is_product_promotion_so = COALESCE(is_product_promotion_so, FALSE),
    is_product_promotion_final = COALESCE(is_product_promotion_final, FALSE),
    is_product_promotion_po = COALESCE(is_product_promotion_po, FALSE)
WHERE promo_so1 IS NULL
   OR promo_so2 IS NULL
   OR promo_so3 IS NULL
   OR promo_so4 IS NULL
   OR promo_so5 IS NULL
   OR promo_final1 IS NULL
   OR promo_final2 IS NULL
   OR promo_final3 IS NULL
   OR promo_final4 IS NULL
   OR promo_final5 IS NULL
   OR promo_po1 IS NULL
   OR promo_po2 IS NULL
   OR promo_po3 IS NULL
   OR promo_po4 IS NULL
   OR promo_po5 IS NULL
   OR promo_remarks_so IS NULL
   OR promo_remarks_final IS NULL
   OR promo_remarks_po IS NULL
   OR is_product_promotion_so IS NULL
   OR is_product_promotion_final IS NULL
   OR is_product_promotion_po IS NULL;

COMMIT;
