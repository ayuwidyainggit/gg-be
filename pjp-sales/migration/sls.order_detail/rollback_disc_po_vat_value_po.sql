-- Rollback: Remove disc_po and vat_value_po from sls.order_detail

BEGIN;

ALTER TABLE IF EXISTS sls.order_detail
    DROP COLUMN IF EXISTS vat_value_po,
    DROP COLUMN IF EXISTS disc_po;

COMMIT;
