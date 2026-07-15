CREATE TABLE IF NOT EXISTS pjp_principles.routes (
	id serial4 NOT NULL,
	route_code int8 NOT NULL,
	route_name varchar(125) NOT NULL,
	pjp_id int4 NULL,
	cust_id varchar(125) DEFAULT NULL::character varying NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,

	CONSTRAINT routes_pkey PRIMARY KEY (id),
	CONSTRAINT routes_route_code_key UNIQUE (route_code)
);