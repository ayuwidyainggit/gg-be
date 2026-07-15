-- Migration: Enhance inv.stock_opname_details table
-- Description: Add stock_opname_det_id (PK auto increment), unit_id1-3, conv_unit1-3, qty1-3, qty_so1-3, stock_opname_date
-- Author: System
-- Date: 2026-01-06

-- Step 1: Add stock_opname_det_id as primary key with auto increment
-- First, check if column already exists and add if not
DO $$
BEGIN
    -- Add stock_opname_det_id column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'inv' 
        AND table_name = 'stock_opname_details' 
        AND column_name = 'stock_opname_det_id'
    ) THEN
        -- Drop existing primary key constraint if exists
        ALTER TABLE inv.stock_opname_details
        DROP CONSTRAINT IF EXISTS stock_opname_details_pkey;
        
        -- Add the column as BIGSERIAL (auto increment)
        ALTER TABLE inv.stock_opname_details
        ADD COLUMN stock_opname_det_id BIGSERIAL;
        
        -- Create new composite primary key including stock_opname_det_id
        ALTER TABLE inv.stock_opname_details
        ADD CONSTRAINT stock_opname_details_pkey PRIMARY KEY (cust_id, doc_no, pro_id, stock_opname_det_id);
        
        -- Add comment
        COMMENT ON COLUMN inv.stock_opname_details.stock_opname_det_id IS 'Primary key auto increment';
    END IF;
END $$;

-- Step 2: Add unit_id columns
ALTER TABLE inv.stock_opname_details
ADD COLUMN IF NOT EXISTS unit_id1 VARCHAR(5) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS unit_id2 VARCHAR(5) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS unit_id3 VARCHAR(5) NOT NULL DEFAULT '';

-- Step 3: Add conv_unit columns
ALTER TABLE inv.stock_opname_details
ADD COLUMN IF NOT EXISTS conv_unit1 FLOAT4 NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS conv_unit2 FLOAT4 NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS conv_unit3 FLOAT4 NOT NULL DEFAULT 0;

-- Step 4: Add qty columns
ALTER TABLE inv.stock_opname_details
ADD COLUMN IF NOT EXISTS qty1 FLOAT4 NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS qty2 FLOAT4 NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS qty3 FLOAT4 NOT NULL DEFAULT 0;

-- Step 5: Add qty_so columns (nullable with default 0)
ALTER TABLE inv.stock_opname_details
ADD COLUMN IF NOT EXISTS qty_so1 FLOAT4 NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS qty_so2 FLOAT4 NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS qty_so3 FLOAT4 NULL DEFAULT 0;

-- Step 6: Add stock_opname_date (nullable datetime)
ALTER TABLE inv.stock_opname_details
ADD COLUMN IF NOT EXISTS stock_opname_date TIMESTAMPTZ NULL;

-- Add comments
COMMENT ON COLUMN inv.stock_opname_details.unit_id1 IS 'Unit ID 1';
COMMENT ON COLUMN inv.stock_opname_details.unit_id2 IS 'Unit ID 2';
COMMENT ON COLUMN inv.stock_opname_details.unit_id3 IS 'Unit ID 3';
COMMENT ON COLUMN inv.stock_opname_details.conv_unit1 IS 'Conversion unit 1';
COMMENT ON COLUMN inv.stock_opname_details.conv_unit2 IS 'Conversion unit 2';
COMMENT ON COLUMN inv.stock_opname_details.conv_unit3 IS 'Conversion unit 3';
COMMENT ON COLUMN inv.stock_opname_details.qty1 IS 'Quantity 1';
COMMENT ON COLUMN inv.stock_opname_details.qty2 IS 'Quantity 2';
COMMENT ON COLUMN inv.stock_opname_details.qty3 IS 'Quantity 3';
COMMENT ON COLUMN inv.stock_opname_details.qty_so1 IS 'Stock opname quantity 1 - filled when assigned user has input stock opname quantity';
COMMENT ON COLUMN inv.stock_opname_details.qty_so2 IS 'Stock opname quantity 2 - filled when assigned user has input stock opname quantity';
COMMENT ON COLUMN inv.stock_opname_details.qty_so3 IS 'Stock opname quantity 3 - filled when assigned user has input stock opname quantity';
COMMENT ON COLUMN inv.stock_opname_details.stock_opname_date IS 'Date and time when stock opname is submitted';
