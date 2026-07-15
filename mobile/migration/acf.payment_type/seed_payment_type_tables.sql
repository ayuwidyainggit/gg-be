-- ============================================
-- Payment Type Seeder (PostgreSQL Style)
-- Created: 2026-02-19
-- Description: Seed initial data for acf.payment_type table
-- ============================================

BEGIN;

INSERT INTO acf.payment_type (payment_type_id, payment_type_code, payment_type_name, created_by)
VALUES 
    (1, 'CASH', 'Cash', 1),
    (2, 'TRANSFER', 'Transfer', 1),
    (3, 'CHEQUE_GIRO', 'Cheque/Bilyet Giro', 1),
    (4, 'RETURN', 'Return', 1),
    (5, 'CREDIT_DEBIT', 'Credit/Debit', 1)
ON CONFLICT (payment_type_id) DO UPDATE SET
    payment_type_code = EXCLUDED.payment_type_code,
    payment_type_name = EXCLUDED.payment_type_name,
    updated_at = CURRENT_TIMESTAMP,
    updated_by = 1;

-- Reset sequence to the next value after 5
SELECT setval('acf.payment_type_payment_type_id_seq', 5, true);

COMMIT;
