BEGIN;

ALTER TABLE acf.deposit_detail
ADD COLUMN salesman_id INT;

COMMIT;