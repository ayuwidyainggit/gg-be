-- Migration: Fix duplicate invoice_no per cust_id and add unique index
-- Purpose: Resolve duplicate data issue before creating unique index uq_order_cust_invoice_no
-- Date: 2026-03-02

BEGIN;

-- 1) Update duplicate rows by appending incremental suffix: -1, -2, ...
-- Keep the first row in each (cust_id, invoice_no) group unchanged.
WITH ranked AS (
    SELECT
        ctid,
        cust_id,
        invoice_no,
        ROW_NUMBER() OVER (
            PARTITION BY cust_id, invoice_no
            ORDER BY updated_at DESC NULLS LAST, created_at DESC NULLS LAST, ro_no DESC NULLS LAST, ctid DESC
        ) AS rn
    FROM sls."order"
    WHERE invoice_no IS NOT NULL
), to_update AS (
    SELECT
        ctid,
        invoice_no,
        (rn - 1) AS dup_seq
    FROM ranked
    WHERE rn > 1
)
UPDATE sls."order" o
SET invoice_no = CONCAT(to_update.invoice_no, '-', to_update.dup_seq::text)
FROM to_update
WHERE o.ctid = to_update.ctid;

-- 2) Safety check: fail transaction if duplicates still exist.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM sls."order"
        WHERE invoice_no IS NOT NULL
        GROUP BY cust_id, invoice_no
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'Duplicate (cust_id, invoice_no) still exists after normalization';
    END IF;
END $$;

-- 3) Create unique index after data is clean.
CREATE UNIQUE INDEX IF NOT EXISTS uq_order_cust_invoice_no
ON sls."order"(cust_id, invoice_no);

COMMIT;

-- Verification query (run manually if needed):
-- SELECT cust_id, invoice_no, COUNT(*)
-- FROM sls."order"
-- WHERE invoice_no IS NOT NULL
-- GROUP BY cust_id, invoice_no
-- HAVING COUNT(*) > 1;
