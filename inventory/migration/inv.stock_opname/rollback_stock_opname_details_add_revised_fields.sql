-- Rollback: Remove revised fields from inv.stock_opname_details table
-- Description: Rollback migration for revised columns
-- Author: System
-- Date: 2026-01-06

-- Drop revised columns
ALTER TABLE inv.stock_opname_details
DROP COLUMN IF EXISTS qty_revised1,
DROP COLUMN IF EXISTS qty_revised2,
DROP COLUMN IF EXISTS qty_revised3,
DROP COLUMN IF EXISTS user_revised,
DROP COLUMN IF EXISTS revised_date;

