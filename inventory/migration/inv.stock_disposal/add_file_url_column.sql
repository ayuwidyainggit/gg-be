-- ============================================
-- Stock Disposal Migration - Add file_url Column
-- Created: 2025-12-11
-- Description: Add file_url column to stock_disposal_detail table for storing file URLs from object storage
-- ============================================

-- Add file_url column to inv.stock_disposal_detail
ALTER TABLE inv.stock_disposal_detail 
  ADD COLUMN IF NOT EXISTS file_url varchar(500) NULL;

-- Add comment for documentation
COMMENT ON COLUMN inv.stock_disposal_detail.file_url IS 'File URL from object storage (optional - nullable, used for new uploads via multipart form data)';
