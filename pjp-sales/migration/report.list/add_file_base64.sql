-- Migration: Add file_base64 column to report.list table
-- Created: 2024-12-18
-- Purpose: Store Excel file content as base64 string for Download Sales Order feature

ALTER TABLE report.list ADD COLUMN IF NOT EXISTS file_base64 TEXT;

COMMENT ON COLUMN report.list.file_base64 IS 'Base64-encoded Excel file content for async download';
