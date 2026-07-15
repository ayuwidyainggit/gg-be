ALTER TABLE tms.shipment_invoices
DROP CONSTRAINT IF EXISTS fk_shipment_invoices_shipment_no;

ALTER TABLE tms.shipments
DROP CONSTRAINT IF EXISTS unique_shipment_no;