CREATE TABLE picklist.picklist (
    picklist_no VARCHAR(128) PRIMARY KEY,
    driver VARCHAR(128) NULL,
    helper VARCHAR(128) NULL,
    vehicle VARCHAR(128) NULL,
    order_no VARCHAR(128) NULL,
    updated_by VARCHAR(128) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);