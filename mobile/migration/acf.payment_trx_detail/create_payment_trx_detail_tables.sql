-- ============================================
-- Payment Transaction Detail Migration (PostgreSQL Style)
-- Created: 2026-02-18
-- Description: Create acf.payment_trx_detail table
-- ============================================

BEGIN;

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS acf;

-- ============================================
-- 1. Create acf.payment_trx_detail table
-- ============================================
CREATE TABLE IF NOT EXISTS acf.payment_trx_detail (
    cust_id varchar(10) NOT NULL,
    payment_trx_det_id bigserial NOT NULL,
    payment_trx_id int8 NOT NULL,
    cndn_no varchar(30) NULL,
    pay_type int2 NOT NULL, -- 1: cash, 2: transfer
    bank_transfer_no int4 NULL,
    cheque_giro_no int4 NULL,
    amount numeric(20, 4) NOT NULL DEFAULT 0,
    
    -- Audit Columns
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    CONSTRAINT payment_trx_det_id_pkey PRIMARY KEY (cust_id, payment_trx_det_id)
);

-- ============================================
-- 2. Create Indexes
-- ============================================
CREATE INDEX idx_payment_trx_det_header ON acf.payment_trx_detail USING btree (cust_id, payment_trx_id);
CREATE INDEX idx_payment_trx_det_pay_type ON acf.payment_trx_detail USING btree (pay_type);

-- ============================================
-- 3. Create Foreign Keys
-- ============================================
ALTER TABLE acf.payment_trx_detail ADD CONSTRAINT fk_payment_trx_det_cust FOREIGN KEY (cust_id) REFERENCES smc.m_customer(cust_id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE acf.payment_trx_detail ADD CONSTRAINT fk_payment_trx_det_header FOREIGN KEY (cust_id, payment_trx_id) REFERENCES acf.payment_trx(cust_id, payment_trx_id) ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;
