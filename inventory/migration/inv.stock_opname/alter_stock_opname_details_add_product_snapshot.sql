-- Migration: Add product status snapshot fields to inv.stock_opname_details table
-- Description: Add pro_status_before, is_active_before for restoring m_product status on completed/rejected
-- Author: System
-- Date: 2026-02-07

-- Add product snapshot columns
ALTER TABLE inv.stock_opname_details
ADD COLUMN IF NOT EXISTS pro_status_before INT4 NULL,
ADD COLUMN IF NOT EXISTS is_active_before BOOLEAN NULL;

-- Add comments
COMMENT ON COLUMN inv.stock_opname_details.pro_status_before IS 'Product pro_status value before stock opname (snapshot on create)';
COMMENT ON COLUMN inv.stock_opname_details.is_active_before IS 'Product is_active value before stock opname (snapshot on create)';
