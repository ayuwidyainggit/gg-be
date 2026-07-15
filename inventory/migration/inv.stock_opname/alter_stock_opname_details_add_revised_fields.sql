-- Migration: Add revised fields to inv.stock_opname_details table
-- Description: Add qty_revised1, qty_revised2, qty_revised3, user_revised, revised_date columns
-- Author: System
-- Date: 2026-01-06

-- Add revised quantity columns
ALTER TABLE inv.stock_opname_details
ADD COLUMN IF NOT EXISTS qty_revised1 FLOAT4 NULL,
ADD COLUMN IF NOT EXISTS qty_revised2 FLOAT4 NULL,
ADD COLUMN IF NOT EXISTS qty_revised3 FLOAT4 NULL,
ADD COLUMN IF NOT EXISTS user_revised INT8 NULL,
ADD COLUMN IF NOT EXISTS revised_date TIMESTAMPTZ NULL;

-- Add comments
COMMENT ON COLUMN inv.stock_opname_details.qty_revised1 IS 'Revised quantity 1';
COMMENT ON COLUMN inv.stock_opname_details.qty_revised2 IS 'Revised quantity 2';
COMMENT ON COLUMN inv.stock_opname_details.qty_revised3 IS 'Revised quantity 3';
COMMENT ON COLUMN inv.stock_opname_details.user_revised IS 'User ID who performed the revision';
COMMENT ON COLUMN inv.stock_opname_details.revised_date IS 'Date and time when revision was performed';

