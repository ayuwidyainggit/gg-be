-- Add base64 payload storage for report.list download history.
ALTER TABLE report.list
ADD COLUMN IF NOT EXISTS file_base64 TEXT;
