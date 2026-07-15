-- ============================================
-- Payment Transaction Migration (PostgreSQL Style)
-- Created: 2026-02-18
-- Description: Create acf.payment_trx table with all required columns
-- ============================================

BEGIN;

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS acf;

-- ============================================
-- 1. Create acf.payment_trx table
-- ============================================
CREATE TABLE IF NOT EXISTS acf.payment_trx (
    cust_id varchar(10) NOT NULL,
    payment_trx_id bigserial NOT NULL,
    outlet_id int4 NOT NULL,
    emp_id int4 NOT NULL,
    po_number varchar(50) NOT NULL,
    document_no varchar(50) NOT NULL,
    trx_source varchar(1) NOT NULL, -- 'C' for Canvas, 'O' for Taking Order
    trx_ref_no int8 NULL,
    total_transaction numeric(20, 4) NOT NULL DEFAULT 0,
    payment_amount numeric(20, 4) NOT NULL DEFAULT 0,
    remaining_amount numeric(20, 4) NOT NULL DEFAULT 0,
    date date NOT NULL DEFAULT CURRENT_DATE,
    notes varchar(250) NULL,
    files jsonb NULL,
    
    -- Audit Columns
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    CONSTRAINT payment_trx_id_pkey PRIMARY KEY (cust_id, payment_trx_id)
);

-- ============================================
-- 2. Create Indexes
-- ============================================
CREATE INDEX idx_payment_trx_cust_po ON acf.payment_trx USING btree (cust_id, po_number);
CREATE INDEX idx_payment_trx_ref_no ON acf.payment_trx USING btree (trx_ref_no);
CREATE INDEX idx_payment_trx_doc_no ON acf.payment_trx USING btree (document_no);
CREATE INDEX idx_payment_trx_outlet_id ON acf.payment_trx USING btree (outlet_id);
CREATE INDEX idx_payment_trx_emp_id ON acf.payment_trx USING btree (emp_id);
CREATE INDEX idx_payment_trx_cust_id ON acf.payment_trx USING btree (cust_id);

-- ============================================
-- 3. Create Foreign Keys
-- ============================================
ALTER TABLE acf.payment_trx ADD CONSTRAINT fk_payment_trx_cust FOREIGN KEY (cust_id) REFERENCES smc.m_customer(cust_id) ON DELETE RESTRICT ON UPDATE CASCADE;

COMMIT;
