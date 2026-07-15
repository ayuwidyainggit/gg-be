-- ============================================
-- Replenishment Order Migration
-- Created: 2025-11-27
-- Description: Create tables for Replenishment Order feature
-- ============================================

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS inv;

-- ============================================
-- Create replenishment_order table
-- ============================================
-- Note: ro_id is varchar (document number), not auto-increment
-- Format: ARO[YY][MM][DD][3-digit sequential]
CREATE TABLE IF NOT EXISTS inv.replenishment_order (
    cust_id varchar(10) NOT NULL,
    ro_id varchar(30) NOT NULL,
    date date NOT NULL,
    sup_id int8 NOT NULL,
    wh_id int8 NOT NULL,
    delivery_type varchar(25) NOT NULL,
    replenishment_type varchar(25) NOT NULL,
    so_start_date date NULL,
    so_end_date date NULL,
    delivery_date date NULL,
    note varchar(255) NOT NULL,
    status int4 NOT NULL DEFAULT 1,
    po_no varchar(20) NULL,
    so_no varchar(20) NULL,
    created_by int8 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int8 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    -- Primary Key
    CONSTRAINT replenishment_order_pkey PRIMARY KEY (cust_id, ro_id),
    
    -- Foreign Keys
    CONSTRAINT fk_replenishment_order_cust FOREIGN KEY (cust_id) 
        REFERENCES smc.m_customer(cust_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_replenishment_order_status FOREIGN KEY (status) 
        REFERENCES inv.replenishment_status(status_code) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_replenishment_order_delivery_type FOREIGN KEY (delivery_type) 
        REFERENCES inv.delivery_type(delivery_type_code) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_replenishment_order_replenishment_type FOREIGN KEY (replenishment_type) 
        REFERENCES inv.replenishment_type(replenishment_type_code) ON UPDATE CASCADE ON DELETE RESTRICT
    -- Note: Foreign keys to m_supplier and m_warehouse removed - these tables may use composite keys
    -- CONSTRAINT fk_replenishment_order_supplier FOREIGN KEY (sup_id) 
    --     REFERENCES mst.m_supplier(sup_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    -- CONSTRAINT fk_replenishment_order_warehouse FOREIGN KEY (wh_id) 
    --     REFERENCES mst.m_warehouse(wh_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_replenishment_order_ro_id ON inv.replenishment_order(ro_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_cust_date ON inv.replenishment_order(cust_id, date);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_wh_id ON inv.replenishment_order(wh_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_sup_id ON inv.replenishment_order(sup_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_status ON inv.replenishment_order(status);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_po_no ON inv.replenishment_order(po_no);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_so_no ON inv.replenishment_order(so_no);

-- ============================================
-- Create replenishment_order_detail table
-- ============================================
-- Note: ro_detail_id menggunakan BIGSERIAL untuk auto-increment
CREATE TABLE IF NOT EXISTS inv.replenishment_order_detail (
    cust_id varchar(10) NOT NULL,
    ro_detail_id bigserial NOT NULL,
    ro_id varchar(30) NOT NULL,
    pro_id int8 NOT NULL,
    order_booking_qty1 numeric(20,4) NOT NULL DEFAULT 0,
    order_booking_qty2 numeric(20,4) NOT NULL DEFAULT 0,
    order_booking_qty3 numeric(20,4) NOT NULL DEFAULT 0,
    purch_price1 numeric(20,4) NOT NULL DEFAULT 0,
    purch_price2 numeric(20,4) NOT NULL DEFAULT 0,
    purch_price3 numeric(20,4) NOT NULL DEFAULT 0,
    estimated_price numeric(20,4) NOT NULL DEFAULT 0,
    created_by int8 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int8 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false,
    
    -- Primary Key
    CONSTRAINT replenishment_order_detail_pkey PRIMARY KEY (cust_id, ro_detail_id),
    
    -- Foreign Keys
    CONSTRAINT fk_replenishment_order_detail_cust FOREIGN KEY (cust_id) 
        REFERENCES smc.m_customer(cust_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_replenishment_order_detail_header FOREIGN KEY (cust_id, ro_id) 
        REFERENCES inv.replenishment_order(cust_id, ro_id) ON UPDATE CASCADE ON DELETE CASCADE
    -- Note: Foreign key to m_product removed - m_product may use composite key
    -- CONSTRAINT fk_replenishment_order_detail_product FOREIGN KEY (pro_id) 
    --     REFERENCES mst.m_product(pro_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_replenishment_order_detail_ro_id ON inv.replenishment_order_detail(ro_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_detail_pro_id ON inv.replenishment_order_detail(pro_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_detail_cust_ro ON inv.replenishment_order_detail(cust_id, ro_id);

-- ============================================
-- Add comments for documentation
-- ============================================
COMMENT ON TABLE inv.replenishment_order IS 'Replenishment Order header table - records replenishment order transactions';
COMMENT ON TABLE inv.replenishment_order_detail IS 'Replenishment Order detail table - records products for each replenishment order';

COMMENT ON COLUMN inv.replenishment_order.ro_id IS 'Auto-generated document number format: ARO[YY][MM][DD][3-digit sequential]';
COMMENT ON COLUMN inv.replenishment_order.status IS 'Status code (1-7): Default 1 (Need Review)';
COMMENT ON COLUMN inv.replenishment_order.delivery_type IS 'Delivery type: Full, Partial';
COMMENT ON COLUMN inv.replenishment_order.replenishment_type IS 'Replenishment type: Replenishment, Replenishment Event';
COMMENT ON COLUMN inv.replenishment_order_detail.order_booking_qty1 IS 'Quantity 1 (Smallest unit)';
COMMENT ON COLUMN inv.replenishment_order_detail.order_booking_qty2 IS 'Quantity 2 (Middle unit)';
COMMENT ON COLUMN inv.replenishment_order_detail.order_booking_qty3 IS 'Quantity 3 (Largest unit)';
COMMENT ON COLUMN inv.replenishment_order_detail.purch_price1 IS 'Purchase price 1 (Smallest unit)';
COMMENT ON COLUMN inv.replenishment_order_detail.purch_price2 IS 'Purchase price 2 (Middle unit)';
COMMENT ON COLUMN inv.replenishment_order_detail.purch_price3 IS 'Purchase price 3 (Largest unit)';

