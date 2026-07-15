ALTER TABLE inv.replenishment_order_detail
	DROP COLUMN IF EXISTS sap_qty3,
	DROP COLUMN IF EXISTS sap_purch_price3;
