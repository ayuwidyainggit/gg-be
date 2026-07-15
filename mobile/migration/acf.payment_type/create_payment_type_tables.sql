-- ============================================
-- Payment Type Migration (PostgreSQL Style)
-- Created: 2026-02-18
-- Description: Create acf.payment_type table
-- ============================================

BEGIN;

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS acf;

-- ============================================
-- 1. Create acf.payment_type table
-- ============================================
CREATE TABLE IF NOT EXISTS acf.payment_type (
    payment_type_id serial NOT NULL,
    payment_type_code varchar(20) NOT NULL,
    payment_type_name varchar(50) NOT NULL,
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int8 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    CONSTRAINT payment_type_id_pkey PRIMARY KEY (payment_type_id)
);

-- ============================================
-- 2. Create Indexes
-- ============================================
CREATE INDEX idx_payment_type_code ON acf.payment_type USING btree (payment_type_code);

COMMIT;
