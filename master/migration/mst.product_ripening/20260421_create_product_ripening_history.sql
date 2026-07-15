/*
History for product ripening is intentionally disabled for now.

CREATE TABLE IF NOT EXISTS mst.product_ripening_history (
    id BIGSERIAL PRIMARY KEY,
    cust_id VARCHAR(30) NOT NULL,
    distributor_id BIGINT NULL,
    per_year INTEGER NULL,
    per_id INTEGER NULL,
    week_id INTEGER NULL,
    week_start DATE NULL,
    week_end DATE NULL,
    source_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    file_url TEXT NULL,
    file_name VARCHAR(255) NULL,
    total_row INTEGER NOT NULL DEFAULT 0,
    success_row INTEGER NOT NULL DEFAULT 0,
    failed_row INTEGER NOT NULL DEFAULT 0,
    error_summary TEXT NULL,
    processed_by BIGINT NOT NULL,
    processed_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_ripening_history_customer FOREIGN KEY (cust_id) REFERENCES smc.m_customer (cust_id),
    CONSTRAINT fk_product_ripening_history_distributor FOREIGN KEY (distributor_id) REFERENCES mst.m_distributor (distributor_id),
    CONSTRAINT fk_product_ripening_history_non_negative CHECK (
        total_row >= 0 AND success_row >= 0 AND failed_row >= 0
    )
);

CREATE INDEX IF NOT EXISTS idx_product_ripening_history_cust_processed
    ON mst.product_ripening_history (cust_id, processed_at DESC);

CREATE INDEX IF NOT EXISTS idx_product_ripening_history_plan
    ON mst.product_ripening_history (cust_id, distributor_id, per_year, per_id, week_id);

CREATE INDEX IF NOT EXISTS idx_product_ripening_history_source_status
    ON mst.product_ripening_history (cust_id, source_type, status);

COMMENT ON TABLE mst.product_ripening_history IS 'Audit trail for product ripening imports and manual edits.';
*/
