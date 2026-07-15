-- ============================================
-- Order Detail PO Fields Migration
-- Created: 2025-12-04
-- Description: Add purchase order and final price fields to sls.order_detail table
-- ============================================

BEGIN;

-- Add qty_po columns (float4, default 0)
-- Kuantitas produk purchase order
ALTER TABLE IF EXISTS sls.order_detail
    ADD COLUMN IF NOT EXISTS qty_po1 FLOAT4 DEFAULT 0,
    ADD COLUMN IF NOT EXISTS qty_po2 FLOAT4 DEFAULT 0,
    ADD COLUMN IF NOT EXISTS qty_po3 FLOAT4 DEFAULT 0;

-- Add sell_price_po columns (numeric(20,4), default 0)
-- Harga Beli purchase order
ALTER TABLE IF EXISTS sls.order_detail
    ADD COLUMN IF NOT EXISTS sell_price_po1 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sell_price_po2 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sell_price_po3 NUMERIC(20,4) DEFAULT 0;

-- Add sell_price_final columns (numeric(20,4), default 0)
-- Harga Jual final order
ALTER TABLE IF EXISTS sls.order_detail
    ADD COLUMN IF NOT EXISTS sell_price_final1 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sell_price_final2 NUMERIC(20,4) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sell_price_final3 NUMERIC(20,4) DEFAULT 0;

-- Add comments for documentation
COMMENT ON COLUMN sls.order_detail.qty_po1 IS 'Kuantitas produk purchase order';
COMMENT ON COLUMN sls.order_detail.qty_po2 IS 'Kuantitas produk purchase order';
COMMENT ON COLUMN sls.order_detail.qty_po3 IS 'Kuantitas produk purchase order';
COMMENT ON COLUMN sls.order_detail.sell_price_po1 IS 'Harga Beli purchase order';
COMMENT ON COLUMN sls.order_detail.sell_price_po2 IS 'Harga Beli purchase order';
COMMENT ON COLUMN sls.order_detail.sell_price_po3 IS 'Harga Beli purchase order';
COMMENT ON COLUMN sls.order_detail.sell_price_final1 IS 'Harga Jual final order';
COMMENT ON COLUMN sls.order_detail.sell_price_final2 IS 'Harga Jual final order';
COMMENT ON COLUMN sls.order_detail.sell_price_final3 IS 'Harga Jual final order';

COMMIT;
