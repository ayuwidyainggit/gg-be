-- Adds is_addition_from flag for SAP replenishment export.
-- true  = created by system
-- false = created by user

ALTER TABLE inv.replenishment_order
	ADD COLUMN IF NOT EXISTS is_addition_from bool NOT NULL DEFAULT true;

COMMENT ON COLUMN inv.replenishment_order.is_addition_from IS 'true = created by system, false = created by user (SAP export)';
