-- ============================================
-- Expense Type Migration - Drop Source Field
-- Created: 2025-01-XX
-- Description: Remove source field and index from acf.expense_type table
-- ============================================

BEGIN;

-- Drop index
DROP INDEX IF EXISTS acf.idx_expense_type_source;

-- Drop column
ALTER TABLE IF EXISTS acf.expense_type
    DROP COLUMN IF EXISTS source;

COMMIT;
