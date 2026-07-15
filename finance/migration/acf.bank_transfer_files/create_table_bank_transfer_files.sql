-- ============================================
-- Create Table acf.bank_transfer_files
-- Created: 2026-01-30
-- Description: Table to store proof of payment files for bank transfers
-- ============================================

BEGIN;

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS acf;

CREATE TABLE IF NOT EXISTS acf.bank_transfer_files (
    bank_transfer_file_id SERIAL PRIMARY KEY,
    cust_id VARCHAR(10) NOT NULL,
    bank_transfer_no VARCHAR(30) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_key TEXT NOT NULL,
    media_category TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add comments for documentation
COMMENT ON TABLE acf.bank_transfer_files IS 'Stores file attachments for bank transfers (e.g. proof of payment)';
COMMENT ON COLUMN acf.bank_transfer_files.bank_transfer_file_id IS 'Primary key';
COMMENT ON COLUMN acf.bank_transfer_files.cust_id IS 'Customer ID';
COMMENT ON COLUMN acf.bank_transfer_files.bank_transfer_no IS 'Relation to bank transfer document number';
COMMENT ON COLUMN acf.bank_transfer_files.file_name IS 'Original file name';
COMMENT ON COLUMN acf.bank_transfer_files.file_url IS 'Public/Private URL to access the file';
COMMENT ON COLUMN acf.bank_transfer_files.file_key IS 'Object storage key (S3/MinIO)';
COMMENT ON COLUMN acf.bank_transfer_files.media_category IS 'Category of the media (e.g. proof_of_payment)';
COMMENT ON COLUMN acf.bank_transfer_files.file_size IS 'File size in bytes';
COMMENT ON COLUMN acf.bank_transfer_files.created_at IS 'Timestamp when the file was uploaded';

-- Create index for faster lookup by bank_transfer_no
CREATE INDEX IF NOT EXISTS idx_bank_transfer_files_no ON acf.bank_transfer_files(bank_transfer_no);
CREATE INDEX IF NOT EXISTS idx_bank_transfer_files_cust ON acf.bank_transfer_files(cust_id);

COMMIT;
