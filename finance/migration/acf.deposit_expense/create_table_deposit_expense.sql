BEGIN;

CREATE TABLE IF NOT EXISTS  acf.deposit_expense (
    deposit_expense_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    deposit_no VARCHAR(30) NOT NULL,
    expense_id BIGINT NOT NULL,
    payment_amount NUMERIC(20,4) NOT NULL,
    created_by INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by INT,
    updated_at TIMESTAMPTZ,
    deleted_by INT,
    deleted_at TIMESTAMPTZ,
    is_del BOOLEAN
);

COMMIT;