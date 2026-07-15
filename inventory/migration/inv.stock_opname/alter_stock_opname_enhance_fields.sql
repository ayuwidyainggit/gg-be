-- Migration: Enhance inv.stock_opname table
-- Description: Ensure all required fields exist with correct data types
-- Author: System
-- Date: 2026-01-06

-- Step 1: Add/Update is_process (bool, NOT NULL, default false)
ALTER TABLE inv.stock_opname
ADD COLUMN IF NOT EXISTS is_process BOOLEAN NOT NULL DEFAULT false;

-- Step 2: Add/Update stock_type (varchar(3), NOT NULL)
ALTER TABLE inv.stock_opname
ADD COLUMN IF NOT EXISTS stock_type VARCHAR(3) NOT NULL DEFAULT 'G';

-- Step 3: Add/Update division_id (int8, nullable)
ALTER TABLE inv.stock_opname
ADD COLUMN IF NOT EXISTS division_id INT8 NULL;

-- Step 4: Update product_hierarchy to VARCHAR(50) if it's currently INT
-- First check current type and alter if needed
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'inv' 
        AND table_name = 'stock_opname' 
        AND column_name = 'product_hierarchy'
        AND data_type = 'integer'
    ) THEN
        -- Convert integer to varchar(50)
        ALTER TABLE inv.stock_opname
        ALTER COLUMN product_hierarchy TYPE VARCHAR(50) USING product_hierarchy::VARCHAR;
        
        -- Make it NOT NULL if not already
        ALTER TABLE inv.stock_opname
        ALTER COLUMN product_hierarchy SET NOT NULL;
    ELSIF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'inv' 
        AND table_name = 'stock_opname' 
        AND column_name = 'product_hierarchy'
    ) THEN
        -- Add column if it doesn't exist
        ALTER TABLE inv.stock_opname
        ADD COLUMN product_hierarchy VARCHAR(50) NOT NULL;
    END IF;
END $$;

-- Step 5: Add/Update principal_id, pl_lane, brand_id, sbrand1_id
-- Note: User mentioned Array<int>, but PostgreSQL doesn't have native array for integers in this context
-- We'll keep them as INT8 (single value) or use JSON/JSONB if array is needed
-- For now, keeping as INT8 as per existing model
ALTER TABLE inv.stock_opname
ADD COLUMN IF NOT EXISTS principal_id INT8 NULL,
ADD COLUMN IF NOT EXISTS pl_lane INT8 NULL,
ADD COLUMN IF NOT EXISTS brand_id INT8 NULL,
ADD COLUMN IF NOT EXISTS sbrand1_id INT8 NULL;

-- Step 6: Add/Update input_by (varchar(50), NOT NULL)
ALTER TABLE inv.stock_opname
ADD COLUMN IF NOT EXISTS input_by VARCHAR(50) NOT NULL DEFAULT 'Web';

-- Step 7: Add/Update emp_id (int8, NOT NULL)
ALTER TABLE inv.stock_opname
ADD COLUMN IF NOT EXISTS emp_id INT8 NOT NULL DEFAULT 0;

-- Step 8: Add/Update is_revised (bool, nullable, default false)
ALTER TABLE inv.stock_opname
ADD COLUMN IF NOT EXISTS is_revised BOOLEAN NULL DEFAULT false;

-- Add comments
COMMENT ON COLUMN inv.stock_opname.is_process IS 'False: manual process has not been performed, True: manual process has been performed';
COMMENT ON COLUMN inv.stock_opname.stock_type IS 'Stock type (G/E/BS)';
COMMENT ON COLUMN inv.stock_opname.division_id IS 'Division ID';
COMMENT ON COLUMN inv.stock_opname.product_hierarchy IS 'Product hierarchy';
COMMENT ON COLUMN inv.stock_opname.principal_id IS 'Principal ID - relation with mst.m_product.principal_id';
COMMENT ON COLUMN inv.stock_opname.pl_lane IS 'Product Line Lane - relation with mst.m_brand.pl_id';
COMMENT ON COLUMN inv.stock_opname.brand_id IS 'Brand ID - relation with mst.sub_brand1.brand_id';
COMMENT ON COLUMN inv.stock_opname.sbrand1_id IS 'Sub Brand 1 ID - relation with mst.m_product.sbrand1_id';
COMMENT ON COLUMN inv.stock_opname.input_by IS 'Input by (Mobile/Manual/Web)';
COMMENT ON COLUMN inv.stock_opname.emp_id IS 'Employee ID - assigned to';
COMMENT ON COLUMN inv.stock_opname.is_revised IS 'Will be true when stock opname ON GOING status has been edited';
