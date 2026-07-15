-- ============================================
-- Expense Migration
-- Created: 2025-01-XX
-- Description: Create tables for Expense feature
-- ============================================

BEGIN;

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS acf;

-- ============================================
-- 1. Create ENUM for media_category
-- ============================================
DO $$
BEGIN
  CREATE TYPE acf.media_category_type AS ENUM ('image', 'video');
EXCEPTION WHEN duplicate_object THEN NULL;
END$$;

-- ============================================
-- 2. Create acf.expense_type table (Parameter)
-- ============================================
-- Note: expense_type_id menggunakan SERIAL untuk auto-increment
-- GORM akan handle auto-increment dengan tag autoIncrement
-- Tabel ini global/tidak per-customer, jadi tidak ada cust_id
CREATE TABLE IF NOT EXISTS acf.expense_type (
    expense_type_id serial NOT NULL,
    expense_type_code varchar(10) NULL,
    expense_type_name varchar(100) NULL,
    is_active bool NOT NULL DEFAULT true,
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    -- Primary Key
    CONSTRAINT expense_type_id_pkey PRIMARY KEY (expense_type_id)
    -- Note: Foreign key to sys.m_user removed - m_user uses composite PK (cust_id, user_id)
    -- created_by stores only user_id, so FK constraint cannot be created
);

-- Create index for expense_type_code lookup
CREATE INDEX IF NOT EXISTS idx_expense_type_code ON acf.expense_type(expense_type_code);
CREATE INDEX IF NOT EXISTS idx_expense_type_is_active ON acf.expense_type(is_active);

-- ============================================
-- 3. Create acf.expense table (Header)
-- ============================================
-- Note: expense_id menggunakan BIGSERIAL untuk auto-increment
-- GORM akan handle auto-increment dengan tag autoIncrement
CREATE TABLE IF NOT EXISTS acf.expense (
    cust_id varchar(10) NOT NULL,
    expense_id bigserial NOT NULL,
    expense_type_id int4 NOT NULL,
    date date NOT NULL,
    amount numeric(20,4) NOT NULL DEFAULT 0,
    note varchar(100) NULL,
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    -- Primary Key
    CONSTRAINT expense_id_pkey PRIMARY KEY (cust_id, expense_id),
    
    -- Foreign Keys
    CONSTRAINT fk_expense_cust FOREIGN KEY (cust_id) 
        REFERENCES smc.m_customer(cust_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_expense_type FOREIGN KEY (expense_type_id) 
        REFERENCES acf.expense_type(expense_type_id) ON UPDATE CASCADE ON DELETE RESTRICT
    -- Note: Foreign key to sys.m_user removed - m_user uses composite PK (cust_id, user_id)
    -- created_by stores only user_id, so FK constraint cannot be created
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_expense_cust_date ON acf.expense(cust_id, date);
CREATE INDEX IF NOT EXISTS idx_expense_type_id ON acf.expense(expense_type_id);
CREATE INDEX IF NOT EXISTS idx_expense_cust_id ON acf.expense(cust_id);

-- ============================================
-- 4. Create acf.expense_det table (Detail Outlet)
-- ============================================
-- Note: expense_det_id menggunakan BIGSERIAL untuk auto-increment
-- GORM akan handle auto-increment dengan tag autoIncrement
CREATE TABLE IF NOT EXISTS acf.expense_det (
    cust_id varchar(10) NOT NULL,
    expense_det_id bigserial NOT NULL,
    expense_id int8 NOT NULL,
    outlet_id int4 NOT NULL,
    
    -- Primary Key
    CONSTRAINT expense_det_id_pkey PRIMARY KEY (cust_id, expense_det_id),
    
    -- Foreign Keys
    CONSTRAINT fk_expense_det_cust FOREIGN KEY (cust_id) 
        REFERENCES smc.m_customer(cust_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_expense_det_header FOREIGN KEY (cust_id, expense_id) 
        REFERENCES acf.expense(cust_id, expense_id) ON UPDATE CASCADE ON DELETE CASCADE
    -- Note: Foreign key to mst.m_outlet removed - m_outlet may use composite key
    -- CONSTRAINT fk_expense_det_outlet FOREIGN KEY (outlet_id) 
    --     REFERENCES mst.m_outlet(outlet_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_expense_det_expense_id ON acf.expense_det(expense_id);
CREATE INDEX IF NOT EXISTS idx_expense_det_outlet_id ON acf.expense_det(outlet_id);
CREATE INDEX IF NOT EXISTS idx_expense_det_cust_expense ON acf.expense_det(cust_id, expense_id);

-- ============================================
-- 5. Create acf.expense_file table (Lampiran File)
-- ============================================
-- Note: expense_file_id menggunakan BIGSERIAL untuk auto-increment
-- GORM akan handle auto-increment dengan tag autoIncrement
CREATE TABLE IF NOT EXISTS acf.expense_file (
    cust_id varchar(10) NOT NULL,
    expense_file_id bigserial NOT NULL,
    expense_id int8 NOT NULL,
    file_name varchar(255) NOT NULL,
    file_url varchar(500) NOT NULL,
    file_key acf.media_category_type NOT NULL,
    media_category text NULL,
    file_size bigint NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Primary Key
    CONSTRAINT expense_file_id_pkey PRIMARY KEY (cust_id, expense_file_id),
    
    -- Foreign Keys
    CONSTRAINT fk_expense_file_cust FOREIGN KEY (cust_id) 
        REFERENCES smc.m_customer(cust_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_expense_file_header FOREIGN KEY (cust_id, expense_id) 
        REFERENCES acf.expense(cust_id, expense_id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_expense_file_expense_id ON acf.expense_file(expense_id);
CREATE INDEX IF NOT EXISTS idx_expense_file_cust_expense ON acf.expense_file(cust_id, expense_id);

-- ============================================
-- 6. Add comments for documentation
-- ============================================
COMMENT ON TABLE acf.expense_type IS 'Expense Type master data table - lookup parameter for expense types';
COMMENT ON TABLE acf.expense IS 'Expense header table - records expense transactions';
COMMENT ON TABLE acf.expense_det IS 'Expense detail table - records outlets associated with each expense';
COMMENT ON TABLE acf.expense_file IS 'Expense file attachment table - records file attachments (image/video) for each expense';

COMMENT ON COLUMN acf.expense_type.expense_type_code IS 'Code for expense type';
COMMENT ON COLUMN acf.expense_type.expense_type_name IS 'Name/description of expense type';
COMMENT ON COLUMN acf.expense.date IS 'Date of expense transaction';
COMMENT ON COLUMN acf.expense.amount IS 'Amount of expense (numeric with 4 decimal precision)';
COMMENT ON COLUMN acf.expense.note IS 'Additional notes for the expense';
COMMENT ON COLUMN acf.expense_det.outlet_id IS 'Reference to mst.m_outlet.outlet_id';
COMMENT ON COLUMN acf.expense_file.file_key IS 'Media type: image or video (ENUM)';
COMMENT ON COLUMN acf.expense_file.file_url IS 'URL/path to the file';
COMMENT ON COLUMN acf.expense_file.media_category IS 'Additional media category information';
COMMENT ON COLUMN acf.expense_file.file_size IS 'File size in bytes';

COMMIT;
