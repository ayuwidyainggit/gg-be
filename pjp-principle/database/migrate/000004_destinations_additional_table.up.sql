CREATE TABLE IF NOT EXISTS pjp_principles.destinations_additional (
	id serial4 NOT NULL,
	route_code int8 NOT NULL,
	route_name varchar(125) NOT NULL,
	status varchar(125) DEFAULT 'additional'::character varying NULL,
	verified_date timestamp NULL,
    "date" timestamp NULL,

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
	
    CONSTRAINT destinations_additional_pkey PRIMARY KEY (id),
    CONSTRAINT fk_destinations_additional_pjp_principles FOREIGN KEY (pjp_id) REFERENCES pjp_principles.permanent_journey_plans(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_destinations_additional_routes FOREIGN KEY (route_code) REFERENCES pjp_principles.routes(route_code) ON DELETE CASCADE ON UPDATE CASCADE
);
