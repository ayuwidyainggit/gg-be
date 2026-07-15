BEGIN;

ALTER TABLE acf.expense
    ADD COLUMN collector_id INT;

COMMIT;