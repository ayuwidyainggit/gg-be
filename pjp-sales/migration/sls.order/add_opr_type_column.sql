-- ============================================
-- Add opr_type Column Migration
-- Created: 2026-03-31
-- Description: Add opr_type column to sls.order table
-- ============================================

BEGIN;

ALTER TABLE IF EXISTS sls.order
    ADD COLUMN IF NOT EXISTS opr_type CHAR(1);

COMMENT ON COLUMN sls.order.opr_type IS 'Operation type code snapshot (aligned with mst.m_salesman.opr_type / mst.m_operation_type.operation_type_code)';

-- Backfill from salesman master data when available
UPDATE sls.order o
SET opr_type = ms.opr_type
FROM mst.m_salesman ms
WHERE o.cust_id = ms.cust_id
  AND o.salesman_id = ms.emp_id
  AND o.opr_type IS NULL
  AND ms.opr_type IS NOT NULL;

COMMIT;
