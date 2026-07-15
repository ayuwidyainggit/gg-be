CREATE TABLE mst.auto_replenishment_product (
    id                  BIGSERIAL PRIMARY KEY,

    cust_id             VARCHAR(10) NOT NULL,
    pro_id              INT NOT NULL,
    distributor_id      INT NOT NULL,

    limit_action        VARCHAR(20) NOT NULL
        CHECK (LOWER(limit_action) IN ('restricted', 'warning')),

    max_order_qty       INT NOT NULL,
    max_order_type      CHAR(1) NOT NULL
        CHECK (max_order_type IN ('S', 'M', 'L')),

    min_stock_qty       INT NOT NULL,
    min_stock_type      CHAR(1) NOT NULL
        CHECK (min_stock_type IN ('S', 'M', 'L')),

    safety_stock_qty    INT NOT NULL,
    safety_stock_type   CHAR(1) NOT NULL
        CHECK (safety_stock_type IN ('S', 'M', 'L')),

    min_order_qty       INT NOT NULL,
    min_order_type      CHAR(1) NOT NULL
        CHECK (min_order_type IN ('S', 'M', 'L')),

    is_active           BOOLEAN DEFAULT TRUE,

    created_by          INT NOT NULL,
    created_at          TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),

    updated_by          INT,
    updated_at          TIMESTAMPTZ(6),

    deleted_by          INT,
    deleted_at          TIMESTAMPTZ(6),

    is_del              BOOLEAN DEFAULT FALSE,

    -- Foreign Keys (assumed)
    CONSTRAINT fk_auto_replenishment_product_customer
        FOREIGN KEY (cust_id)
        REFERENCES smc.m_customer (cust_id),

    CONSTRAINT fk_auto_replenishment_product_product
        FOREIGN KEY (pro_id)
        REFERENCES mst.m_product (pro_id),

    CONSTRAINT fk_auto_replenishment_product_distributor
        FOREIGN KEY (distributor_id)
        REFERENCES mst.m_distributor (distributor_id)
);
