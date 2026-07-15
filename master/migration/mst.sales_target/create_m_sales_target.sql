-- Create table mst.m_sales_target
CREATE TABLE IF NOT EXISTS mst.m_sales_target (
    cust_id VARCHAR(10) NOT NULL,
    sales_target_id SERIAL PRIMARY KEY,
    sales_target_distributor_yearly_id INT NOT NULL,
    sales_target_distributor_monthly_id INT NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    allocated_total BIGINT NOT NULL,
    monthly_target BIGINT NOT NULL,
    remaining BIGINT NOT NULL,
    status INT NOT NULL DEFAULT 1,
    created_by INT NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INT,
    updated_at TIMESTAMPTZ(6),
    deleted_by INT,
    deleted_at TIMESTAMPTZ(6),
    is_del BOOLEAN DEFAULT FALSE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_m_sales_target_cust_id ON mst.m_sales_target(cust_id);
CREATE INDEX IF NOT EXISTS idx_m_sales_target_year ON mst.m_sales_target(year);
CREATE INDEX IF NOT EXISTS idx_m_sales_target_month ON mst.m_sales_target(month);
CREATE INDEX IF NOT EXISTS idx_m_sales_target_status ON mst.m_sales_target(status);
CREATE INDEX IF NOT EXISTS idx_m_sales_target_is_del ON mst.m_sales_target(is_del);

-- Add comments
COMMENT ON TABLE mst.m_sales_target IS 'Table to store sales target for salesman';
COMMENT ON COLUMN mst.m_sales_target.sales_target_id IS 'Primary key for sales target';
COMMENT ON COLUMN mst.m_sales_target.status IS '0: Draft, 1: Active, 2: Inactive';