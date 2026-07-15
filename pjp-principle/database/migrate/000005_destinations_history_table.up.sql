CREATE TABLE IF NOT EXISTS pjp_principles.destinations_history (
	id serial4 NOT NULL,
	route_code int8 NOT NULL,
	route_name varchar(125) NOT NULL,
	verified_date timestamp NULL,
    "date" timestamp NULL,
	week int4 NULL,
	"year" int4 NULL,
    index_day int4 NULL,
	start_week timestamp NULL,
	is_in_current_year bool NULL,
	is_additional bool DEFAULT false NULL,

	destination_id int8 NULL,
	destination_code varchar(125) NULL,
	destination_status varchar(125) NULL,
	destination_name varchar(125) NULL,
	destination_address varchar(125) NULL,
	destination_type varchar(125) NULL,
	longitude varchar(125) NULL,
	latitude varchar(125) NULL,

	pjp_id int8 NULL,
	pjp_code int8 NULL,
	old_pjp_id int8 NULL,
	old_pjp_code int8 NULL,
	old_route_code int8 NULL,
	old_route_name varchar(125) DEFAULT NULL::character varying NULL,

	photo varchar(125) DEFAULT NULL::character varying NULL,
	signature varchar(125) DEFAULT NULL::character varying NULL,
	avg_sales_week numeric(10, 2) DEFAULT 0 NULL,

	cust_id varchar(125) DEFAULT NULL::character varying NULL,

	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	
    CONSTRAINT destinations_history_pkey PRIMARY KEY (id)
);
