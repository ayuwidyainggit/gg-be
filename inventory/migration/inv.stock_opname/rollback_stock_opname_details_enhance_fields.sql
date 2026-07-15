-- Rollback Migration: Enhance inv.stock_opname_details table
-- Description: Remove added fields from stock_opname_details enhancement
-- Author: System
-- Date: 2026-01-06

-- Remove stock_opname_date
ALTER TABLE inv.stock_opname_details
DROP COLUMN IF EXISTS stock_opname_date;

-- Remove qty_so columns
ALTER TABLE inv.stock_opname_details
DROP COLUMN IF EXISTS qty_so1,
DROP COLUMN IF EXISTS qty_so2,
DROP COLUMN IF EXISTS qty_so3;

-- Remove qty columns
ALTER TABLE inv.stock_opname_details
DROP COLUMN IF EXISTS qty1,
DROP COLUMN IF EXISTS qty2,
DROP COLUMN IF EXISTS qty3;

-- Remove conv_unit columns
ALTER TABLE inv.stock_opname_details
DROP COLUMN IF EXISTS conv_unit1,
DROP COLUMN IF EXISTS conv_unit2,
DROP COLUMN IF EXISTS conv_unit3;

-- Remove unit_id columns
ALTER TABLE inv.stock_opname_details
DROP COLUMN IF EXISTS unit_id1,
DROP COLUMN IF EXISTS unit_id2,
DROP COLUMN IF EXISTS unit_id3;

-- Remove stock_opname_det_id and restore original primary key
DO $$
BEGIN
    -- Drop the new primary key constraint
    ALTER TABLE inv.stock_opname_details
    DROP CONSTRAINT IF EXISTS stock_opname_details_pkey;
    
    -- Restore original primary key (if it was cust_id, doc_no, pro_id)
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE table_schema = 'inv' 
        AND table_name = 'stock_opname_details' 
        AND constraint_name = 'stock_opname_details_pkey'
    ) THEN
        ALTER TABLE inv.stock_opname_details
        ADD CONSTRAINT stock_opname_details_pkey PRIMARY KEY (cust_id, doc_no, pro_id);
    END IF;
    
    -- Drop stock_opname_det_id column
    ALTER TABLE inv.stock_opname_details
    DROP COLUMN IF EXISTS stock_opname_det_id;
    
    -- Drop sequence if exists
    DROP SEQUENCE IF EXISTS inv.stock_opname_details_stock_opname_det_id_seq;
END $$;
