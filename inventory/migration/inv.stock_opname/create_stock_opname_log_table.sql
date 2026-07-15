-- Migration: Create inv.stock_opname_log table
-- Description: Table for logging stock opname status changes
-- Author: System
-- Date: 2026-01-05

-- Create stock_opname_log table
CREATE TABLE IF NOT EXISTS inv.stock_opname_log (
    id_stock_opname_log SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    execution_time TIMESTAMPTZ NOT NULL,
    old_status INT NOT NULL,
    status INT NOT NULL,
    ref_id VARCHAR(50) NOT NULL,
    transaction_code VARCHAR(30) NOT NULL,
    ref_table_name VARCHAR(50) NOT NULL DEFAULT 'inv.stock_opname',
    triggered_by VARCHAR(20) NOT NULL DEFAULT 'MANUAL',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by INT8,
    cust_id VARCHAR(10) NOT NULL
);

-- Add indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_stock_opname_log_transaction_code ON inv.stock_opname_log(transaction_code);
CREATE INDEX IF NOT EXISTS idx_stock_opname_log_cust_id ON inv.stock_opname_log(cust_id);
CREATE INDEX IF NOT EXISTS idx_stock_opname_log_created_at ON inv.stock_opname_log(created_at);

-- Add comments
COMMENT ON TABLE inv.stock_opname_log IS 'Log table for stock opname status changes';
COMMENT ON COLUMN inv.stock_opname_log.id_stock_opname_log IS 'Primary key, auto-increment';
COMMENT ON COLUMN inv.stock_opname_log.title IS 'Action title: Assign, Need Review, Completed, Rejected';
COMMENT ON COLUMN inv.stock_opname_log.execution_time IS 'Time when the action was executed';
COMMENT ON COLUMN inv.stock_opname_log.old_status IS 'Previous status before update';
COMMENT ON COLUMN inv.stock_opname_log.status IS 'New status after update';
COMMENT ON COLUMN inv.stock_opname_log.ref_id IS 'Reference ID for the action (ObjectID)';
COMMENT ON COLUMN inv.stock_opname_log.transaction_code IS 'Document number (doc_no)';
COMMENT ON COLUMN inv.stock_opname_log.ref_table_name IS 'Reference table name';
COMMENT ON COLUMN inv.stock_opname_log.triggered_by IS 'Trigger source: MANUAL, SYSTEM';
COMMENT ON COLUMN inv.stock_opname_log.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN inv.stock_opname_log.created_by IS 'User ID who created the record';
COMMENT ON COLUMN inv.stock_opname_log.cust_id IS 'Customer ID';

