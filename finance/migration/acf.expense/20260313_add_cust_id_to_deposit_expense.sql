BEGIN;

ALTER TABLE acf.deposit_expense
ADD COLUMN cust_id VARCHAR(10) NULL;

ALTER TABLE acf.deposit_expense
ADD CONSTRAINT fk_deposit_expense_cust
FOREIGN KEY (cust_id)
REFERENCES smc.m_customer (cust_id)
ON UPDATE CASCADE
ON DELETE RESTRICT;

COMMIT;