CREATE TABLE IF NOT EXISTS tms.shipment_order_status (
    id SERIAL PRIMARY KEY,
    order_no VARCHAR(125) NULL,
    status_order VARCHAR(125) NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    updated_at timestamptz NULL
)