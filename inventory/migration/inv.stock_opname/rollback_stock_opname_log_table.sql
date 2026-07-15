-- Rollback: Drop inv.stock_opname_log table
-- Description: Rollback migration for stock opname log table
-- Author: System
-- Date: 2026-01-05

-- Drop indexes
DROP INDEX IF EXISTS inv.idx_stock_opname_log_transaction_code;
DROP INDEX IF EXISTS inv.idx_stock_opname_log_cust_id;
DROP INDEX IF EXISTS inv.idx_stock_opname_log_created_at;

-- Drop table
DROP TABLE IF EXISTS inv.stock_opname_log;

