ALTER TABLE acf.payment_trx
ADD COLUMN IF NOT EXISTS invoice_no varchar(255) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS collection_no varchar(255) DEFAULT NULL;