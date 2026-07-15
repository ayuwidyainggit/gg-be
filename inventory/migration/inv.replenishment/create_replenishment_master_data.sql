-- ============================================
-- Replenishment Master Data Migration
-- Created: 2025-11-27
-- Description: Create master data tables for Replenishment Type and Delivery Type
-- ============================================

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS inv;

-- ============================================
-- Create replenishment_type table
-- ============================================
CREATE TABLE IF NOT EXISTS inv.replenishment_type (
    replenishment_type_code varchar(25) NOT NULL,
    replenishment_type_name varchar(100) NOT NULL,
    description text NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz(6) NULL,
    
    -- Primary Key
    CONSTRAINT replenishment_type_pkey PRIMARY KEY (replenishment_type_code)
);

-- Create index for replenishment_type_name lookup
CREATE INDEX IF NOT EXISTS idx_replenishment_type_name ON inv.replenishment_type(replenishment_type_name);

-- ============================================
-- Insert Replenishment Type Data
-- ============================================
-- Replenishment Type: Replenishment
INSERT INTO inv.replenishment_type (replenishment_type_code, replenishment_type_name, description, created_at, updated_at)
VALUES (
    'Replenishment',
    'Replenishment',
    'Standard replenishment order type',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (replenishment_type_code) DO UPDATE
SET 
    replenishment_type_name = EXCLUDED.replenishment_type_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Replenishment Type: Replenishment Event
INSERT INTO inv.replenishment_type (replenishment_type_code, replenishment_type_name, description, created_at, updated_at)
VALUES (
    'Replenishment Event',
    'Replenishment Event',
    'Event-based replenishment order type',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (replenishment_type_code) DO UPDATE
SET 
    replenishment_type_name = EXCLUDED.replenishment_type_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- ============================================
-- Create delivery_type table
-- ============================================
CREATE TABLE IF NOT EXISTS inv.delivery_type (
    delivery_type_code varchar(25) NOT NULL,
    delivery_type_name varchar(100) NOT NULL,
    description text NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz(6) NULL,
    
    -- Primary Key
    CONSTRAINT delivery_type_pkey PRIMARY KEY (delivery_type_code)
);

-- Create index for delivery_type_name lookup
CREATE INDEX IF NOT EXISTS idx_delivery_type_name ON inv.delivery_type(delivery_type_name);

-- ============================================
-- Insert Delivery Type Data
-- ============================================
-- Delivery Type: Full
INSERT INTO inv.delivery_type (delivery_type_code, delivery_type_name, description, created_at, updated_at)
VALUES (
    'Full',
    'Full',
    'Full delivery type - complete order delivery',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (delivery_type_code) DO UPDATE
SET 
    delivery_type_name = EXCLUDED.delivery_type_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Delivery Type: Partial
INSERT INTO inv.delivery_type (delivery_type_code, delivery_type_name, description, created_at, updated_at)
VALUES (
    'Partial',
    'Partial',
    'Partial delivery type - partial order delivery',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (delivery_type_code) DO UPDATE
SET 
    delivery_type_name = EXCLUDED.delivery_type_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- ============================================
-- Add comments for documentation
-- ============================================
COMMENT ON TABLE inv.replenishment_type IS 'Master data table for Replenishment Type lookup';
COMMENT ON TABLE inv.delivery_type IS 'Master data table for Delivery Type lookup';

COMMENT ON COLUMN inv.replenishment_type.replenishment_type_code IS 'Replenishment type code: Replenishment, Replenishment Event';
COMMENT ON COLUMN inv.delivery_type.delivery_type_code IS 'Delivery type code: Full, Partial';

