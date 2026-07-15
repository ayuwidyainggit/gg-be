-- Rollback: Remove product snapshot fields from inv.stock_opname_details table
-- Description: Rollback migration for pro_status_before, is_active_before columns
-- Author: System
-- Date: 2026-02-07

-- Drop product snapshot columns
ALTER TABLE inv.stock_opname_details
DROP COLUMN IF EXISTS pro_status_before,
DROP COLUMN IF EXISTS is_active_before;
