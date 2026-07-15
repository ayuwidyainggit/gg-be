CREATE TABLE IF NOT EXISTS pjp_principles.permanent_journey_plans (
	id serial4 NOT NULL,
	pjp_code int8 NOT NULL,
	operation_type varchar(125) NOT NULL,
	team_salesman varchar(125) NULL,
	salesman_id int8 NULL,
	salesman_name varchar(125) NULL,
	warehouse_id int8 NULL,
	warehouse_name varchar(125) NULL,
	pjp_mode varchar(125) DEFAULT 'manual'::character varying NULL,
	status varchar(125) DEFAULT 'pending'::character varying NULL,
	salesman_code varchar(125) NULL,
	approval_status varchar(32) DEFAULT 'Draft'::character varying NOT NULL,
	cust_id varchar(125) DEFAULT NULL::character varying NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,

	CONSTRAINT permanent_journey_plans_pjp_code_cust_id_key UNIQUE (pjp_code, cust_id),
	CONSTRAINT permanent_journey_plans_pkey PRIMARY KEY (id)
);