-- ============================================
-- Expense Type Migration - Add Fields
-- Created: 2025-01-XX
-- Description: Add is_active and source fields to acf.expense_type table
-- ============================================

BEGIN;

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS acf;

-- ============================================
-- Add is_active field (if not exists)
-- ============================================
-- Note: Field is_active mungkin sudah ada di tabel, jadi gunakan IF NOT EXISTS
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'acf' 
        AND table_name = 'expense_type' 
        AND column_name = 'is_active'
    ) THEN
        ALTER TABLE acf.expense_type
            ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true;
        
        -- Create index for is_active if not exists
        CREATE INDEX IF NOT EXISTS idx_expense_type_is_active ON acf.expense_type(is_active);
    END IF;
END$$;

-- ============================================
-- Add source field (if not exists)
-- ============================================
-- source: 1 = web, 2 = mobile
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = 'acf' 
        AND table_name = 'expense_type' 
        AND column_name = 'source'
    ) THEN
        ALTER TABLE acf.expense_type
            ADD COLUMN source INTEGER NOT NULL DEFAULT 1;
        
        -- Create index for source if not exists (sering digunakan untuk filtering)
        CREATE INDEX IF NOT EXISTS idx_expense_type_source ON acf.expense_type(source);
        
        -- Update existing data: set default source = 1 (web) untuk data yang sudah ada
        UPDATE acf.expense_type
        SET source = 1
        WHERE source IS NULL;
    END IF;
END$$;

-- ============================================
-- Add comments for documentation
-- ============================================
COMMENT ON COLUMN acf.expense_type.is_active IS 'Status aktif/inaktif expense type. True = Active, False = Inactive';
COMMENT ON COLUMN acf.expense_type.source IS 'Source of expense type: 1 = web, 2 = mobile';

COMMIT;
