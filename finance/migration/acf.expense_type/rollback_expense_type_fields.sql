-- ============================================
-- Expense Type Migration Rollback
-- Created: 2025-01-XX
-- Description: Rollback script to remove is_active and source fields from acf.expense_type table
-- ============================================

BEGIN;

-- ============================================
-- Drop Indexes
-- ============================================
DROP INDEX IF EXISTS acf.idx_expense_type_source;
DROP INDEX IF EXISTS acf.idx_expense_type_is_active;

-- ============================================
-- Drop Columns
-- ============================================
ALTER TABLE IF EXISTS acf.expense_type
    DROP COLUMN IF EXISTS source,
    DROP COLUMN IF EXISTS is_active;

COMMIT;
