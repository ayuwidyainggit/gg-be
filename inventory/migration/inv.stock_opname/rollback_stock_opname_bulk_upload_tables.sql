-- Rollback: Drop inv.stock_opname_bulk_upload and inv.stock_opname_bulk_upload_items tables
-- Description: Rollback migration for bulk upload tables
-- Date: 2026-01-06

-- Drop indexes for stock_opname_bulk_upload_items
DROP INDEX IF EXISTS inv.idx_stock_opname_bulk_upload_items_upload_id;
DROP INDEX IF EXISTS inv.idx_stock_opname_bulk_upload_items_product_id;

-- Drop stock_opname_bulk_upload_items table first (has FK dependency)
DROP TABLE IF EXISTS inv.stock_opname_bulk_upload_items;

-- Drop indexes for stock_opname_bulk_upload
DROP INDEX IF EXISTS inv.idx_stock_opname_bulk_upload_doc_no;
DROP INDEX IF EXISTS inv.idx_stock_opname_bulk_upload_uploaded_at;

-- Drop stock_opname_bulk_upload table
DROP TABLE IF EXISTS inv.stock_opname_bulk_upload;

