-- ============================================
-- Payment Transaction Detail Migration Rollback (PostgreSQL Style)
-- Created: 2026-02-18
-- Description: Rollback script to remove acf.payment_trx_detail table
-- ============================================

BEGIN;

-- ============================================
-- 1. Drop Foreign Keys
-- ============================================
ALTER TABLE IF EXISTS acf.payment_trx_detail DROP CONSTRAINT IF EXISTS fk_payment_trx_det_header;
ALTER TABLE IF EXISTS acf.payment_trx_detail DROP CONSTRAINT IF EXISTS fk_payment_trx_det_cust;

-- ============================================
-- 2. Drop Indexes
-- ============================================
DROP INDEX IF EXISTS acf.idx_payment_trx_det_pay_type;
DROP INDEX IF EXISTS acf.idx_payment_trx_det_header;

-- ============================================
-- 3. Drop Table
-- ============================================
DROP TABLE IF EXISTS acf.payment_trx_detail;

COMMIT;
