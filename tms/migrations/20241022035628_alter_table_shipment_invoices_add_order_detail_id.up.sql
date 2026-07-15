ALTER TABLE tms.shipment_invoices
ADD COLUMN IF NOT EXISTS order_detail_id BIGINT DEFAULT NULL;