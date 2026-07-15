CREATE TABLE picklist.order_picklist (
    order_no VARCHAR(128) PRIMARY KEY,
    invoice_no VARCHAR(128) NULL,
    outlet_name VARCHAR(128) NULL,
    salesman VARCHAR(128) NULL,
    invoice_date DATE,
    due_date DATE,
    total_price NUMERIC,
    ppn NUMERIC,
    discount NUMERIC,
    total_unpaid NUMERIC,
    payment_type VARCHAR(128) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);