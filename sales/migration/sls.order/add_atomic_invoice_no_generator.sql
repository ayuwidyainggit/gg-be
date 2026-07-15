-- Migration: add atomic invoice number generator and enforce unique (cust_id, invoice_no)
-- Purpose: prevent race condition for invoice_no on concurrent invoice creation
-- Date: 2026-03-03

BEGIN;

-- NOTE:
-- Duplicate cleanup below sets duplicate invoice_no values to NULL (keeps the most recent row per group).
-- This is intentionally safe for index creation, but original duplicate values are not fully restorable in down migration.

-- 1) Prepare counter table for atomic daily sequence by cust_id.
CREATE TABLE IF NOT EXISTS sls.invoice_no_counter (
    cust_id VARCHAR(20) NOT NULL,
    seq_date DATE NOT NULL,
    last_seq INTEGER NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cust_id, seq_date)
);

-- 2) Create generator function.
CREATE OR REPLACE FUNCTION sls.generate_invoice_no(p_cust_id VARCHAR, p_invoice_date DATE DEFAULT CURRENT_DATE)
RETURNS VARCHAR
LANGUAGE plpgsql
AS $$
DECLARE
    v_seq INTEGER;
BEGIN
    IF p_cust_id IS NULL OR LENGTH(TRIM(p_cust_id)) = 0 THEN
        RAISE EXCEPTION 'p_cust_id is required';
    END IF;

    INSERT INTO sls.invoice_no_counter AS c (cust_id, seq_date, last_seq, updated_at)
    VALUES (p_cust_id, p_invoice_date, 1, NOW())
    ON CONFLICT (cust_id, seq_date)
    DO UPDATE
    SET last_seq = c.last_seq + 1,
        updated_at = NOW()
    RETURNING last_seq INTO v_seq;

    RETURN CONCAT('INV', TO_CHAR(p_invoice_date, 'YYMMDD'), LPAD(v_seq::TEXT, 4, '0'));
END;
$$;

-- 3) Cleanup existing duplicates safely (keep latest row, nullify others).
WITH ranked AS (
    SELECT
        ctid,
        ROW_NUMBER() OVER (
            PARTITION BY cust_id, invoice_no
            ORDER BY updated_at DESC NULLS LAST, created_at DESC NULLS LAST, ro_no DESC NULLS LAST, ctid DESC
        ) AS rn
    FROM sls."order"
    WHERE invoice_no IS NOT NULL
)
UPDATE sls."order" o
SET invoice_no = NULL
FROM ranked r
WHERE o.ctid = r.ctid
  AND r.rn > 1;

-- 4) Safety check before creating unique index.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM sls."order"
        WHERE invoice_no IS NOT NULL
        GROUP BY cust_id, invoice_no
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'Duplicate (cust_id, invoice_no) still exists after cleanup';
    END IF;
END $$;

-- 5) Enforce uniqueness for non-null invoice numbers only.
CREATE UNIQUE INDEX IF NOT EXISTS uq_order_cust_invoice_no
ON sls."order" (cust_id, invoice_no)
WHERE invoice_no IS NOT NULL;

COMMIT;
