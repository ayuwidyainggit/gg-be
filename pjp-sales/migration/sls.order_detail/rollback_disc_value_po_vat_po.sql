-- Rollback: Remove disc_value_po and vat_po from sls.order_detail

BEGIN;

ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS vat_po,
    DROP COLUMN IF EXISTS disc_value_po;

COMMIT;