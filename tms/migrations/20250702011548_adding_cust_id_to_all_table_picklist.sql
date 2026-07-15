ALTER TABLE picklist.order_picklist
ADD COLUMN cust_id VARCHAR(255);

ALTER TABLE picklist.order_product
ADD COLUMN cust_id VARCHAR(255);

ALTER TABLE picklist.picklist
ADD COLUMN cust_id VARCHAR(255);
