-- Step 1: add column (avoid IF NOT EXISTS + DEFAULT in one statement for CRDB compatibility)
ALTER TABLE mst.m_product
ADD COLUMN is_product_mapping BOOL;

-- Step 2: set default + backfill existing rows
ALTER TABLE mst.m_product
ALTER COLUMN is_product_mapping SET DEFAULT false;

UPDATE mst.m_product
SET is_product_mapping = false
WHERE is_product_mapping IS NULL;

-- Step 3: allow origin = 'product_mapping' for product mapping import
ALTER TABLE mst.m_product
DROP CONSTRAINT IF EXISTS m_product_origin_chk;

ALTER TABLE mst.m_product
ADD CONSTRAINT m_product_origin_chk
CHECK (
    origin = LOWER(origin)
    AND origin IN ('import', 'assignment', 'create', 'bulk', 'product_mapping')
);
