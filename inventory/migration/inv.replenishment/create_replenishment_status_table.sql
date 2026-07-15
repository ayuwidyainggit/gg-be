-- ============================================
-- Replenishment Status Migration
-- Created: 2025-11-27
-- Description: Create table for Replenishment Status master data
-- ============================================

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS inv;

-- ============================================
-- Create replenishment_status table
-- ============================================
-- Status code as primary key (1-7)
CREATE TABLE IF NOT EXISTS inv.replenishment_status (
    status_code int4 NOT NULL,
    status_name varchar(50) NOT NULL,
    description text NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz(6) NULL,
    
    -- Primary Key
    CONSTRAINT replenishment_status_pkey PRIMARY KEY (status_code)
);

-- Create index for status_name lookup (optional, for search)
CREATE INDEX IF NOT EXISTS idx_replenishment_status_name ON inv.replenishment_status(status_name);

-- ============================================
-- Insert Status Data
-- ============================================
-- Status Code 1: Need Review
INSERT INTO inv.replenishment_status (status_code, status_name, description, created_at, updated_at)
VALUES (
    1,
    'Need Review',
    'Replenishment Order has been submitted by the Distributor Admin for approval by a user at the Principal level',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (status_code) DO UPDATE
SET 
    status_name = EXCLUDED.status_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Status Code 2: Approved
INSERT INTO inv.replenishment_status (status_code, status_name, description, created_at, updated_at)
VALUES (
    2,
    'Approved',
    'Replenishment Order has been Approved by a User who has access to the Replenishment Order Approval menu',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (status_code) DO UPDATE
SET 
    status_name = EXCLUDED.status_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Status Code 3: Rejected
INSERT INTO inv.replenishment_status (status_code, status_name, description, created_at, updated_at)
VALUES (
    3,
    'Rejected',
    'Replenishment Order was rejected in the Replenishment Order Approval menu',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (status_code) DO UPDATE
SET 
    status_name = EXCLUDED.status_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Status Code 4: On Delivery
INSERT INTO inv.replenishment_status (status_code, status_name, description, created_at, updated_at)
VALUES (
    4,
    'On Delivery',
    'Status obtained when the ERP/related System posts PO No, SO No, DO No, Delivery date, Vehicle Number, Product code, QTY, Price then sends back the status to Scylla X. If replenishment is not integrated with another ERP/System, then the On Delivery status is obtained when the order has been approved and there is no API to another system',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (status_code) DO UPDATE
SET 
    status_name = EXCLUDED.status_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Status Code 5: Completed
INSERT INTO inv.replenishment_status (status_code, status_name, description, created_at, updated_at)
VALUES (
    5,
    'Completed',
    'PO No has been related to Goods Receipt document No',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (status_code) DO UPDATE
SET 
    status_name = EXCLUDED.status_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Status Code 6: Processed
INSERT INTO inv.replenishment_status (status_code, status_name, description, created_at, updated_at)
VALUES (
    6,
    'Processed',
    'Status obtained when the ERP/related System posts PO No, SO No, Product code, QTY, Price and sends back the status to Scylla X',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (status_code) DO UPDATE
SET 
    status_name = EXCLUDED.status_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Status Code 7: Cancelled
INSERT INTO inv.replenishment_status (status_code, status_name, description, created_at, updated_at)
VALUES (
    7,
    'Cancelled',
    'Status obtained when the ERP/related System posts rejection',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (status_code) DO UPDATE
SET 
    status_name = EXCLUDED.status_name,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- ============================================
-- Add comments for documentation
-- ============================================
COMMENT ON TABLE inv.replenishment_status IS 'Master data table for Replenishment Order status codes';
COMMENT ON COLUMN inv.replenishment_status.status_code IS 'Status code (1-7): Primary key for replenishment status';
COMMENT ON COLUMN inv.replenishment_status.status_name IS 'Status name: Need Review, Approved, Rejected, On Delivery, Completed, Processed, Cancelled';
COMMENT ON COLUMN inv.replenishment_status.description IS 'Detailed description of what each status means';

