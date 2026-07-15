CREATE TABLE mst.m_product_assignment (
    id BIGSERIAL PRIMARY KEY,
    cust_id VARCHAR(10) NOT NULL,
    action_date DATE NOT NULL,
    pro_id INT NOT NULL,
    distributor_id INT NOT NULL,
    assignment_type VARCHAR(20) NOT NULL
        CHECK (LOWER(assignment_type) IN ('assignment', 'remove_assignment')),
    created_by INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_m_product_assignment_customer
        FOREIGN KEY (cust_id)
        REFERENCES smc.m_customer (cust_id),
    CONSTRAINT fk_m_product_assignment_product
        FOREIGN KEY (pro_id)
        REFERENCES mst.m_product (pro_id),
    CONSTRAINT fk_m_product_assignment_distributor
        FOREIGN KEY (distributor_id)
        REFERENCES mst.m_distributor (distributor_id)
);

CREATE INDEX idx_m_product_assignment_cust_id
    ON mst.m_product_assignment (cust_id);

CREATE INDEX idx_m_product_assignment_pro_id
    ON mst.m_product_assignment (pro_id);

CREATE INDEX idx_m_product_assignment_distributor_id
    ON mst.m_product_assignment (distributor_id);

CREATE INDEX idx_m_product_assignment_created_at
    ON mst.m_product_assignment (created_at DESC);
