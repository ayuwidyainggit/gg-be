CREATE TABLE IF NOT EXISTS pjp_principles.route_pop_permanent (
	id serial4 NOT NULL,
	"year" int8 NULL,
	week int8 NULL,
	"date" timestamp NULL,
	"day" varchar(125) NULL,
	route_code int8 NULL,
	pjp_id int8 NULL,
	pjp_code int8 NULL,
	cust_id varchar(125) DEFAULT NULL::character varying NULL,

	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,

	CONSTRAINT route_pop_permanent_pkey PRIMARY KEY (id),
    CONSTRAINT fk_route_pop_permanent_pjp FOREIGN KEY (pjp_id) REFERENCES pjp_principles.permanent_journey_plans(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_route_pop_permanent_route FOREIGN KEY (route_code) REFERENCES pjp_principles.routes(route_code) ON DELETE CASCADE ON UPDATE CASCADE
);