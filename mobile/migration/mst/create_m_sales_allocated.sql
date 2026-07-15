BEGIN;

CREATE TABLE mst.m_sales_allocated (
	cust_id varchar(10) NOT NULL,
	sales_allocated_id serial4 NOT NULL,
	sales_target_id int4 NOT NULL,
	salesman_id int4 NOT NULL,
	sales_team_id int4 NULL,
	allocated int8 NOT NULL,
	is_active bool NULL DEFAULT true,
	created_by int4 NOT NULL,
	created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by int4 NULL,
	updated_at timestamptz(6) NULL,
	deleted_by int4 NULL,
	deleted_at timestamptz(6) NULL,
	is_del bool NULL DEFAULT false,
	CONSTRAINT m_sales_allocated_pkey PRIMARY KEY (sales_allocated_id)
);
CREATE INDEX idx_m_sales_allocated_cust_id ON mst.m_sales_allocated USING btree (cust_id);
CREATE INDEX idx_m_sales_allocated_is_active ON mst.m_sales_allocated USING btree (is_active);
CREATE INDEX idx_m_sales_allocated_is_del ON mst.m_sales_allocated USING btree (is_del);
CREATE INDEX idx_m_sales_allocated_sales_target_id ON mst.m_sales_allocated USING btree (sales_target_id);
CREATE INDEX idx_m_sales_allocated_sales_team_id ON mst.m_sales_allocated USING btree (sales_team_id);
CREATE INDEX idx_m_sales_allocated_salesman_id ON mst.m_sales_allocated USING btree (salesman_id);

CREATE TABLE mst.m_sales_target (
	cust_id varchar(10) NOT NULL,
	sales_target_id serial4 NOT NULL,
	sales_target_distributor_yearly_id int4 NOT NULL,
	sales_target_distributor_monthly_id int4 NOT NULL,
	"month" int4 NOT NULL,
	"year" int4 NOT NULL,
	allocated_total int8 NOT NULL,
	monthly_target int8 NOT NULL,
	remaining int8 NOT NULL,
	status int4 NOT NULL DEFAULT 1,
	created_by int4 NOT NULL,
	created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by int4 NULL,
	updated_at timestamptz(6) NULL,
	deleted_by int4 NULL,
	deleted_at timestamptz(6) NULL,
	is_del bool NULL DEFAULT false,
	CONSTRAINT m_sales_target_pkey PRIMARY KEY (sales_target_id)
);
CREATE INDEX idx_m_sales_target_cust_id ON mst.m_sales_target USING btree (cust_id);
CREATE INDEX idx_m_sales_target_is_del ON mst.m_sales_target USING btree (is_del);
CREATE INDEX idx_m_sales_target_month ON mst.m_sales_target USING btree (month);
CREATE INDEX idx_m_sales_target_status ON mst.m_sales_target USING btree (status);
CREATE INDEX idx_m_sales_target_year ON mst.m_sales_target USING btree (year);

ALTER TABLE mst.m_sales_allocated ADD CONSTRAINT fk_sales_target FOREIGN KEY (sales_target_id) REFERENCES mst.m_sales_target(sales_target_id);




COMMIT;
