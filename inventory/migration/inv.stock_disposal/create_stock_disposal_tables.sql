-- ============================================
-- Stock Disposal Migration
-- Created: 2025-11-17
-- Description: Create tables for Stock Disposal feature
-- ============================================

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS inv;

-- ============================================
-- 1. Create ENUM for media_category
-- ============================================
DO $$
BEGIN
  CREATE TYPE inv.media_category_type AS ENUM ('image', 'video');
EXCEPTION WHEN duplicate_object THEN NULL;
END$$;

-- ============================================
-- 2. Create stock_disposal table
-- ============================================
-- Note: sd_id menggunakan BIGSERIAL untuk auto-increment
-- GORM akan handle auto-increment dengan tag autoIncrement
CREATE TABLE IF NOT EXISTS inv.stock_disposal (
    cust_id varchar(10) NOT NULL,
    sd_id bigserial NOT NULL,
    tr_code varchar(5) NOT NULL,
    disposal_date date NOT NULL,
    sd_number varchar(30) NOT NULL,
    sup_id int8 NOT NULL,
    wh_id int8 NOT NULL,
    stock_type varchar(3) NOT NULL,
    gr_no varchar(30) NULL,
    note varchar(100) NOT NULL,
    sub_total numeric(20,4) NOT NULL DEFAULT 0,
    vat_value numeric(20,4) NOT NULL DEFAULT 0,
    total numeric(20,4) NOT NULL DEFAULT 0,
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    -- Primary Key
    CONSTRAINT sd_id_pkey PRIMARY KEY (cust_id, sd_id),
    
    -- Foreign Keys
    CONSTRAINT fk_stock_disposal_cust FOREIGN KEY (cust_id) 
        REFERENCES smc.m_customer(cust_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_stock_disposal_trans FOREIGN KEY (tr_code) 
        REFERENCES sys.m_trans(tr_code) ON UPDATE CASCADE ON DELETE RESTRICT
    -- Note: Foreign keys to m_supplier and m_warehouse removed - these tables may use composite keys
    -- CONSTRAINT fk_stock_disposal_supplier FOREIGN KEY (sup_id) 
    --     REFERENCES mst.m_supplier(sup_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    -- CONSTRAINT fk_stock_disposal_warehouse FOREIGN KEY (wh_id) 
    --     REFERENCES mst.m_warehouse(wh_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- Create index for sd_number lookup
CREATE INDEX IF NOT EXISTS idx_stock_disposal_sd_number ON inv.stock_disposal(sd_number);
CREATE INDEX IF NOT EXISTS idx_stock_disposal_cust_date ON inv.stock_disposal(cust_id, disposal_date);
CREATE INDEX IF NOT EXISTS idx_stock_disposal_wh_id ON inv.stock_disposal(wh_id);
CREATE INDEX IF NOT EXISTS idx_stock_disposal_sup_id ON inv.stock_disposal(sup_id);

-- ============================================
-- 3. Create stock_disposal_detail table
-- ============================================
-- Note: sd_detail_id menggunakan BIGSERIAL untuk auto-increment
-- GORM akan handle auto-increment dengan tag autoIncrement
CREATE TABLE IF NOT EXISTS inv.stock_disposal_detail (
    cust_id varchar(10) NOT NULL,
    sd_detail_id bigserial NOT NULL,
    sd_id int8 NOT NULL,
    pro_id int8 NOT NULL,
    file_name varchar(255) NOT NULL,
    file_type varchar(50) NOT NULL,
    media_category inv.media_category_type NOT NULL,
    file_base64 text NOT NULL,
    file_size bigint NOT NULL,
    unit_id1 varchar(5) NOT NULL,
    unit_id2 varchar(5) NOT NULL,
    unit_id3 varchar(5) NOT NULL,
    qty1 numeric(20,4) NOT NULL DEFAULT 0,
    qty2 numeric(20,4) NOT NULL DEFAULT 0,
    qty3 numeric(20,4) NOT NULL DEFAULT 0,
    purch_price1 numeric(20,4) NOT NULL DEFAULT 0,
    purch_price2 numeric(20,4) NOT NULL DEFAULT 0,
    purch_price3 numeric(20,4) NOT NULL DEFAULT 0,
    gross_price numeric(20,4) NOT NULL DEFAULT 0,
    vat numeric NOT NULL DEFAULT 0,
    vat_value numeric(20,4) NOT NULL DEFAULT 0,
    sub_total numeric(20,4) NOT NULL DEFAULT 0,
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    -- Primary Key
    CONSTRAINT sd_detail_id_pkey PRIMARY KEY (cust_id, sd_detail_id),
    
    -- Foreign Keys
    CONSTRAINT fk_stock_disposal_detail_cust FOREIGN KEY (cust_id) 
        REFERENCES smc.m_customer(cust_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_stock_disposal_detail_header FOREIGN KEY (cust_id, sd_id) 
        REFERENCES inv.stock_disposal(cust_id, sd_id) ON UPDATE CASCADE ON DELETE CASCADE
    -- Note: Foreign key to m_product removed - m_product may use composite key
    -- CONSTRAINT fk_stock_disposal_detail_product FOREIGN KEY (pro_id) 
    --     REFERENCES mst.m_product(pro_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_stock_disposal_detail_sd_id ON inv.stock_disposal_detail(sd_id);
CREATE INDEX IF NOT EXISTS idx_stock_disposal_detail_pro_id ON inv.stock_disposal_detail(pro_id);
CREATE INDEX IF NOT EXISTS idx_stock_disposal_detail_cust_sd ON inv.stock_disposal_detail(cust_id, sd_id);

-- ============================================
-- 4. Notes on Auto-Increment
-- ============================================
-- sd_id dan sd_detail_id menggunakan BIGSERIAL untuk auto-increment
-- GORM akan handle auto-increment dengan tag `autoIncrement` di model
-- Composite primary key (cust_id, sd_id) tetap bekerja dengan BIGSERIAL
-- 
-- Alternative: Jika perlu sequence manual per cust_id, bisa menggunakan:
-- CREATE SEQUENCE IF NOT EXISTS inv.stock_disposal_sd_id_seq;
-- ALTER TABLE inv.stock_disposal ALTER COLUMN sd_id SET DEFAULT nextval('inv.stock_disposal_sd_id_seq');
--
-- Tapi untuk consistency dengan existing tables (seperti inv.gr), 
-- menggunakan BIGSERIAL sudah cukup karena GORM handle auto-increment

-- ============================================
-- 5. Add comments for documentation
-- ============================================
COMMENT ON TABLE inv.stock_disposal IS 'Stock Disposal header table - records disposal transactions';
COMMENT ON TABLE inv.stock_disposal_detail IS 'Stock Disposal detail table - records products and files for each disposal';

COMMENT ON COLUMN inv.stock_disposal.stock_type IS 'Stock type: G (Good), E (Expired), BS (Bad Stock)';
COMMENT ON COLUMN inv.stock_disposal.sd_number IS 'Auto-generated document number format: SD[YY][MM][DD][3-digit sequential]';
COMMENT ON COLUMN inv.stock_disposal_detail.media_category IS 'Media type: image or video';
COMMENT ON COLUMN inv.stock_disposal_detail.file_base64 IS 'Base64 encoded file content';

