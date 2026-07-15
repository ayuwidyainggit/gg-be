-- ============================================
-- Expense Migration Rollback
-- Created: 2025-01-XX
-- Description: Rollback script to remove all expense-related tables and objects
-- ============================================

BEGIN;

-- ============================================
-- 1. Drop Foreign Key Constraints
-- ============================================
-- Drop foreign keys from child tables first
ALTER TABLE IF EXISTS acf.expense_file 
    DROP CONSTRAINT IF EXISTS fk_expense_file_header,
    DROP CONSTRAINT IF EXISTS fk_expense_file_cust;

ALTER TABLE IF EXISTS acf.expense_det 
    DROP CONSTRAINT IF EXISTS fk_expense_det_header,
    DROP CONSTRAINT IF EXISTS fk_expense_det_cust;

ALTER TABLE IF EXISTS acf.expense 
    DROP CONSTRAINT IF EXISTS fk_expense_type,
    DROP CONSTRAINT IF EXISTS fk_expense_cust;

-- ============================================
-- 2. Drop Indexes
-- ============================================
DROP INDEX IF EXISTS acf.idx_expense_file_cust_expense;
DROP INDEX IF EXISTS acf.idx_expense_file_expense_id;

DROP INDEX IF EXISTS acf.idx_expense_det_cust_expense;
DROP INDEX IF EXISTS acf.idx_expense_det_outlet_id;
DROP INDEX IF EXISTS acf.idx_expense_det_expense_id;

DROP INDEX IF EXISTS acf.idx_expense_cust_id;
DROP INDEX IF EXISTS acf.idx_expense_type_id;
DROP INDEX IF EXISTS acf.idx_expense_cust_date;

DROP INDEX IF EXISTS acf.idx_expense_type_is_active;
DROP INDEX IF EXISTS acf.idx_expense_type_code;

-- ============================================
-- 3. Drop Tables (Child tables first, then parent)
-- ============================================
-- Drop child tables first
DROP TABLE IF EXISTS acf.expense_file;
DROP TABLE IF EXISTS acf.expense_det;

-- Drop parent tables
DROP TABLE IF EXISTS acf.expense;
DROP TABLE IF EXISTS acf.expense_type;

-- ============================================
-- 4. Drop ENUM Type
-- ============================================
DROP TYPE IF EXISTS acf.media_category_type;

-- ============================================
-- 5. Drop Schema (Optional - only if no other tables exist)
-- ============================================
-- Note: Uncomment below if you want to drop the entire schema
-- WARNING: This will drop ALL tables in acf schema, not just expense tables
-- DROP SCHEMA IF EXISTS acf CASCADE;

COMMIT;
