-- Add address1 column to sls.order table
-- Issue: SX-521
ALTER TABLE sls.order ADD COLUMN IF NOT EXISTS address1 VARCHAR(150);