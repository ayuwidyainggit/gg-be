-- Rollback migration: remove atomic invoice generator and unique index
-- Date: 2026-03-03
-- NOTE:
-- This rollback removes schema objects only.
-- Data changed by the up migration duplicate cleanup (invoice_no set to NULL) cannot be automatically restored.

BEGIN;

DROP INDEX IF EXISTS uq_order_cust_invoice_no;
DROP FUNCTION IF EXISTS sls.generate_invoice_no(VARCHAR, DATE);
DROP TABLE IF EXISTS sls.invoice_no_counter;

COMMIT;
