-- mst.m_sales_target_distributor_yearly definition

CREATE TABLE IF NOT EXISTS mst.m_sales_target_distributor_yearly (
    cust_id varchar(10) NOT NULL,
    sales_target_distributor_yearly_id serial4 NOT NULL,
    area_id int4 NOT NULL,
    region_id int4 NOT NULL,
    distributor_id int4 NOT NULL,
    year int4 NOT NULL,
    yearly_target int4 NOT NULL,
    status int4 NOT NULL,
    is_active bool DEFAULT true,
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool DEFAULT false,
    CONSTRAINT m_sales_target_distributor_yearly_pkey PRIMARY KEY (sales_target_distributor_yearly_id)
);

-- mst.m_sales_target_distributor_monthly definition

CREATE TABLE IF NOT EXISTS mst.m_sales_target_distributor_monthly (
    cust_id varchar(10) NOT NULL,
    sales_target_distributor_monthly_id serial4 NOT NULL,
    sales_target_distributor_yearly_id int4 NOT NULL,
    month int2 NOT NULL,
    monthly_target int4 NOT NULL,
    is_active bool DEFAULT true,
    created_by int4 NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int4 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool DEFAULT false,
    CONSTRAINT m_sales_target_distributor_monthly_pkey PRIMARY KEY (sales_target_distributor_monthly_id),
    CONSTRAINT fk_sales_target_distributor_yearly FOREIGN KEY (sales_target_distributor_yearly_id) REFERENCES mst.m_sales_target_distributor_yearly(sales_target_distributor_yearly_id)
);