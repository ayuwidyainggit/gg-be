-- Migration: Create inv.stock_opname_bulk_upload and inv.stock_opname_bulk_upload_items tables
-- Description: Tables for bulk upload stock opname functionality
-- Date: 2026-01-06

-- Create stock_opname_bulk_upload table
CREATE TABLE IF NOT EXISTS inv.stock_opname_bulk_upload (
    upload_id SERIAL PRIMARY KEY,
    doc_no VARCHAR(30) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    status INT NULL,
    total_row INT NULL,
    valid_row INT NULL,
    invalid_row INT NULL,
    uploaded_by VARCHAR(50) NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add comments for stock_opname_bulk_upload
COMMENT ON TABLE inv.stock_opname_bulk_upload IS 'Table for storing bulk upload stock opname files';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.upload_id IS 'Primary key, auto-increment';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.doc_no IS 'Document number reference';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.file_path IS 'Path to the uploaded file';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.status IS 'Upload status';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.total_row IS 'Total number of rows in the file';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.valid_row IS 'Number of valid rows';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.invalid_row IS 'Number of invalid rows';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.uploaded_by IS 'User who uploaded the file';
COMMENT ON COLUMN inv.stock_opname_bulk_upload.uploaded_at IS 'Timestamp when file was uploaded';

-- Create indexes for stock_opname_bulk_upload
CREATE INDEX IF NOT EXISTS idx_stock_opname_bulk_upload_doc_no ON inv.stock_opname_bulk_upload(doc_no);
CREATE INDEX IF NOT EXISTS idx_stock_opname_bulk_upload_uploaded_at ON inv.stock_opname_bulk_upload(uploaded_at);

-- Create stock_opname_bulk_upload_items table
CREATE TABLE IF NOT EXISTS inv.stock_opname_bulk_upload_items (
    stock_opname_bulk_upload_id SERIAL PRIMARY KEY,
    upload_id INT NOT NULL REFERENCES inv.stock_opname_bulk_upload(upload_id) ON DELETE CASCADE,
    product_id INT8 NOT NULL,
    qty_so1 FLOAT4 NOT NULL DEFAULT 0,
    qty_so2 FLOAT4 NOT NULL DEFAULT 0,
    qty_so3 FLOAT4 NOT NULL DEFAULT 0,
    qty_revised1 FLOAT4 NOT NULL DEFAULT 0,
    qty_revised2 FLOAT4 NOT NULL DEFAULT 0,
    qty_revised3 FLOAT4 NOT NULL DEFAULT 0,
    unit_id1 FLOAT4 NOT NULL DEFAULT 0,
    unit_id2 FLOAT4 NOT NULL DEFAULT 0,
    unit_id3 FLOAT4 NOT NULL DEFAULT 0
);

-- Add comments for stock_opname_bulk_upload_items
COMMENT ON TABLE inv.stock_opname_bulk_upload_items IS 'Table for storing bulk upload stock opname item details';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.stock_opname_bulk_upload_id IS 'Primary key, auto-increment';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.upload_id IS 'Foreign key to stock_opname_bulk_upload';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.product_id IS 'Product ID reference';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.qty_so1 IS 'Stock opname quantity 1';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.qty_so2 IS 'Stock opname quantity 2';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.qty_so3 IS 'Stock opname quantity 3';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.qty_revised1 IS 'Revised quantity 1';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.qty_revised2 IS 'Revised quantity 2';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.qty_revised3 IS 'Revised quantity 3';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.unit_id1 IS 'Unit ID 1';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.unit_id2 IS 'Unit ID 2';
COMMENT ON COLUMN inv.stock_opname_bulk_upload_items.unit_id3 IS 'Unit ID 3';

-- Create indexes for stock_opname_bulk_upload_items
CREATE INDEX IF NOT EXISTS idx_stock_opname_bulk_upload_items_upload_id ON inv.stock_opname_bulk_upload_items(upload_id);
CREATE INDEX IF NOT EXISTS idx_stock_opname_bulk_upload_items_product_id ON inv.stock_opname_bulk_upload_items(product_id);

