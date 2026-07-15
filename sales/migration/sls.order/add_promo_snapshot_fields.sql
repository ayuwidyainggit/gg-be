-- ============================================
-- Order Promo Snapshot Header Fields Migration
-- Created: 2026-03-11
-- Description: Add persisted promo snapshot remarks fields to sls.order table
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order
    ADD COLUMN IF NOT EXISTS promo_remarks_so JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS promo_remarks_final JSONB DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS promo_remarks_po JSONB DEFAULT '[]'::jsonb;

COMMENT ON COLUMN sls.order.promo_remarks_so IS 'Persisted promo snapshot remarks for Sales Order tab';
COMMENT ON COLUMN sls.order.promo_remarks_final IS 'Persisted promo snapshot remarks for Final Order tab';
COMMENT ON COLUMN sls.order.promo_remarks_po IS 'Persisted promo snapshot remarks for Purchase Order tab';

UPDATE sls.order
SET
    promo_remarks_so = COALESCE(promo_remarks_so, '[]'::jsonb),
    promo_remarks_final = COALESCE(promo_remarks_final, '[]'::jsonb),
    promo_remarks_po = COALESCE(promo_remarks_po, '[]'::jsonb)
WHERE promo_remarks_so IS NULL
   OR promo_remarks_final IS NULL
   OR promo_remarks_po IS NULL;

COMMIT;
