-- Adds columns written by SAP callback: sap_qty3, sap_purch_price3
-- Applied to: inv.replenishment_order_detail (Scylla inventory service)

ALTER TABLE inv.replenishment_order_detail
	ADD COLUMN IF NOT EXISTS sap_qty3 numeric(20, 4) NULL,
	ADD COLUMN IF NOT EXISTS sap_purch_price3 numeric(20, 4) NULL;

COMMENT ON COLUMN inv.replenishment_order_detail.sap_qty3 IS 'Large qty approved / reported by SAP (unit 3)';
COMMENT ON COLUMN inv.replenishment_order_detail.sap_purch_price3 IS 'Large purchase price from SAP (unit 3)';
