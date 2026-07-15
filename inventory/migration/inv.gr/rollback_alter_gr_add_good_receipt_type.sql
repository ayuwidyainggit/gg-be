-- ============================================
-- Goods Receipt Migration Rollback - Remove good_receipt_type Column
-- Created: 2025-12-02
-- Description: Remove good_receipt_type column from inv.gr table
-- ============================================

-- Remove good_receipt_type column from inv.gr table
ALTER TABLE inv.gr 
  DROP COLUMN IF EXISTS good_receipt_type;
