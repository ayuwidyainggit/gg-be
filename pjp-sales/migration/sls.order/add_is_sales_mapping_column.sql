BEGIN;

ALTER TABLE IF EXISTS sls.order
ADD COLUMN IF NOT EXISTS is_sales_mapping BOOLEAN DEFAULT FALSE;

UPDATE sls.order
SET is_sales_mapping = FALSE
WHERE is_sales_mapping IS NULL;

ALTER TABLE IF EXISTS sls.order
ALTER COLUMN is_sales_mapping SET DEFAULT FALSE;

COMMIT;
