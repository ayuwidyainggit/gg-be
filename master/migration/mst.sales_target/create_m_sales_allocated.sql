-- Create table mst.m_sales_allocated
CREATE TABLE IF NOT EXISTS mst.m_sales_allocated (
    cust_id VARCHAR(10) NOT NULL,
    sales_allocated_id SERIAL PRIMARY KEY,
    sales_target_id INT NOT NULL,
    salesman_id INT NOT NULL,
    sales_team_id INT,
    allocated BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INT NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INT,
    updated_at TIMESTAMPTZ(6),
    deleted_by INT,
    deleted_at TIMESTAMPTZ(6),
    is_del BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_sales_target FOREIGN KEY (sales_target_id) REFERENCES mst.m_sales_target(sales_target_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_m_sales_allocated_cust_id ON mst.m_sales_allocated(cust_id);
CREATE INDEX IF NOT EXISTS idx_m_sales_allocated_sales_target_id ON mst.m_sales_allocated(sales_target_id);
CREATE INDEX IF NOT EXISTS idx_m_sales_allocated_salesman_id ON mst.m_sales_allocated(salesman_id);
CREATE INDEX IF NOT EXISTS idx_m_sales_allocated_sales_team_id ON mst.m_sales_allocated(sales_team_id);
CREATE INDEX IF NOT EXISTS idx_m_sales_allocated_is_active ON mst.m_sales_allocated(is_active);
CREATE INDEX IF NOT EXISTS idx_m_sales_allocated_is_del ON mst.m_sales_allocated(is_del);

-- Add comments
COMMENT ON TABLE mst.m_sales_allocated IS 'Table to store allocated sales target per salesman';
COMMENT ON COLUMN mst.m_sales_allocated.sales_allocated_id IS 'Primary key for sales allocated';
COMMENT ON COLUMN mst.m_sales_allocated.salesman_id IS 'Foreign key reference to mst.m_salesman (emp_id)';
COMMENT ON COLUMN mst.m_sales_allocated.allocated IS 'Target amount allocated to salesman';