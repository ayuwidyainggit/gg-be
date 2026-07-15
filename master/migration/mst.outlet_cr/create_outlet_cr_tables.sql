-- ============================================
-- Outlet Change Request Migration
-- Created: 2025-12-17
-- Description: Create tables for Outlet Change Request and Approval feature
-- ============================================

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS mst;

-- ============================================
-- 1. Create mst.outlet_cr table
-- ============================================
-- Tabel master untuk menyimpan semua perubahan mengenai outlet
-- yang dibedakan dari 2 sumber (source: mobile/web) serta approval dari setiap perubahan
CREATE TABLE IF NOT EXISTS mst.outlet_cr (
    cust_id varchar(10) NOT NULL,
    outlet_cr_id bigserial NOT NULL,
    outlet_id int8 NOT NULL,
    source int4 NOT NULL DEFAULT 1,
    status int4 NOT NULL DEFAULT 1,
    created_by int8 NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    approval_by int8 NULL,
    approval_at timestamptz(6) NULL,
    
    -- Primary Key
    CONSTRAINT outlet_cr_id_pkey PRIMARY KEY (outlet_cr_id),
    
    -- Foreign Keys (composite key karena mst.m_outlet menggunakan composite PK)
    CONSTRAINT fk_outlet_cr_outlet FOREIGN KEY (cust_id, outlet_id) 
        REFERENCES mst.m_outlet(cust_id, outlet_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_outlet_cr_status ON mst.outlet_cr(status);
CREATE INDEX IF NOT EXISTS idx_outlet_cr_outlet_id ON mst.outlet_cr(outlet_id);
CREATE INDEX IF NOT EXISTS idx_outlet_cr_cust_id ON mst.outlet_cr(cust_id);
CREATE INDEX IF NOT EXISTS idx_outlet_cr_source ON mst.outlet_cr(source);
CREATE INDEX IF NOT EXISTS idx_outlet_cr_created_at ON mst.outlet_cr(created_at);

-- Add comments
COMMENT ON TABLE mst.outlet_cr IS 'Tabel master yang menyimpan semua perubahan mengenai outlet yang dibedakan dari 2 sumber (source: mobile/web) serta approval dari setiap perubahan';
COMMENT ON COLUMN mst.outlet_cr.cust_id IS 'ID customer';
COMMENT ON COLUMN mst.outlet_cr.outlet_cr_id IS 'Primary Key';
COMMENT ON COLUMN mst.outlet_cr.outlet_id IS 'Foreign Key ke mst.m_outlet';
COMMENT ON COLUMN mst.outlet_cr.source IS '1: web, 2: mobile';
COMMENT ON COLUMN mst.outlet_cr.status IS '1: pending, 2: approve, 3: reject';
COMMENT ON COLUMN mst.outlet_cr.approval_by IS 'User ID yang melakukan approval';
COMMENT ON COLUMN mst.outlet_cr.approval_at IS 'Timestamp saat approval dilakukan';

-- ============================================
-- 2. Create mst.outlet_cr_det table
-- ============================================
-- Tabel untuk menyimpan field yang berubah
-- 1x perubahan untuk merubah lebih dari 1 field
CREATE TABLE IF NOT EXISTS mst.outlet_cr_det (
    outlet_cr_det_id bigserial NOT NULL,
    outlet_cr_id int8 NOT NULL,
    field_name varchar(30) NOT NULL,
    new_value varchar(225) NULL,
    old_value varchar(225) NULL,
    
    -- Primary Key
    CONSTRAINT outlet_cr_det_id_pkey PRIMARY KEY (outlet_cr_det_id),
    
    -- Foreign Keys
    CONSTRAINT fk_outlet_cr_det_cr FOREIGN KEY (outlet_cr_id) 
        REFERENCES mst.outlet_cr(outlet_cr_id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_outlet_cr_det_cr_id ON mst.outlet_cr_det(outlet_cr_id);
CREATE INDEX IF NOT EXISTS idx_outlet_cr_det_field_name ON mst.outlet_cr_det(field_name);
CREATE INDEX IF NOT EXISTS idx_outlet_cr_det_cr_field ON mst.outlet_cr_det(outlet_cr_id, field_name);

-- Add comments
COMMENT ON TABLE mst.outlet_cr_det IS 'Tabel untuk menyimpan field yang berubah. 1x perubahan untuk merubah lebih dari 1 field';
COMMENT ON COLUMN mst.outlet_cr_det.outlet_cr_det_id IS 'Primary Key';
COMMENT ON COLUMN mst.outlet_cr_det.outlet_cr_id IS 'Foreign Key ke mst.outlet_cr';
COMMENT ON COLUMN mst.outlet_cr_det.field_name IS 'Field name dari mst.m_outlet yang dilakukan perubahan';
COMMENT ON COLUMN mst.outlet_cr_det.new_value IS 'Value perubahan';
COMMENT ON COLUMN mst.outlet_cr_det.old_value IS 'Value sebelum perubahan';
