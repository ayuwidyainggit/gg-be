ALTER TABLE IF EXISTS tms.shipments
ADD CONSTRAINT IF EXISTS unique_shipment_no UNIQUE (shipment_no);

ALTER TABLE IF EXISTS tms.shipment_invoices
ADD CONSTRAINT IF EXISTS fk_shipment_invoices_shipment_no
FOREIGN KEY (shipment_no)
REFERENCES tms.shipments (shipment_no)
ON UPDATE CASCADE
ON DELETE CASCADE;
