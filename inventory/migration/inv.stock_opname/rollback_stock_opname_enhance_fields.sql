-- Rollback Migration: Enhance inv.stock_opname table
-- Description: Remove added/updated fields from stock_opname enhancement
-- Author: System
-- Date: 2026-01-06

-- Note: We only remove columns that were added, not modify existing ones back
-- If columns already existed, they will remain

-- Remove is_revised
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS is_revised;

-- Remove emp_id (only if it was added by this migration)
-- Note: Check if this was already in the table before removing
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS emp_id;

-- Remove input_by (only if it was added by this migration)
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS input_by;

-- Remove sbrand1_id
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS sbrand1_id;

-- Remove brand_id
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS brand_id;

-- Remove pl_lane
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS pl_lane;

-- Remove principal_id
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS principal_id;

-- Note: product_hierarchy type change cannot be easily rolled back
-- If needed, manually convert back to integer type

-- Remove division_id
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS division_id;

-- Remove stock_type (only if it was added by this migration)
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS stock_type;

-- Remove is_process (only if it was added by this migration)
ALTER TABLE inv.stock_opname
DROP COLUMN IF EXISTS is_process;
