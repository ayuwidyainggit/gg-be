-- ============================================
-- Stock Disposal Migration - Make File Fields Optional
-- Created: 2025-01-XX
-- Description: Alter stock_disposal_detail table to make file columns nullable
-- ============================================

-- Make file columns nullable in inv.stock_disposal_detail
ALTER TABLE inv.stock_disposal_detail 
  ALTER COLUMN file_name DROP NOT NULL,
  ALTER COLUMN file_type DROP NOT NULL,
  ALTER COLUMN media_category DROP NOT NULL,
  ALTER COLUMN file_base64 DROP NOT NULL,
  ALTER COLUMN file_size DROP NOT NULL;

-- Add comments for documentation
COMMENT ON COLUMN inv.stock_disposal_detail.file_name IS 'File name (optional - nullable)';
COMMENT ON COLUMN inv.stock_disposal_detail.file_type IS 'File type: jpg, jpeg, png, or MP4 (optional - nullable)';
COMMENT ON COLUMN inv.stock_disposal_detail.media_category IS 'Media type: image or video (optional - nullable)';
COMMENT ON COLUMN inv.stock_disposal_detail.file_base64 IS 'Base64 encoded file content (optional - nullable)';
COMMENT ON COLUMN inv.stock_disposal_detail.file_size IS 'File size in bytes (optional - nullable)';

