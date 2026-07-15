-- ============================================
-- Goods Receipt Migration - Add good_receipt_type Column
-- Created: 2025-12-02
-- Description: Add optional good_receipt_type column to inv.gr table
-- ============================================

-- Add good_receipt_type column to inv.gr table
-- This column is nullable and optional - if not provided in payload, it will be NULL
ALTER TABLE inv.gr 
  ADD COLUMN IF NOT EXISTS good_receipt_type VARCHAR(50) NULL;

-- Add comment for documentation
COMMENT ON COLUMN inv.gr.good_receipt_type IS 'Type of goods receipt (optional - nullable). If not provided in request payload, this field will be NULL in database.';
