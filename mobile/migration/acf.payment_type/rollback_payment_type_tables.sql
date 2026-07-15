-- ============================================
-- Payment Type Migration Rollback (PostgreSQL Style)
-- Created: 2026-02-18
-- Description: Rollback script to remove acf.payment_type table
-- ============================================

BEGIN;

-- ============================================
-- 1. Drop Indexes
-- ============================================
DROP INDEX IF EXISTS acf.idx_payment_type_code;

-- ============================================
-- 2. Drop Table
-- ============================================
DROP TABLE IF EXISTS acf.payment_type;

COMMIT;
