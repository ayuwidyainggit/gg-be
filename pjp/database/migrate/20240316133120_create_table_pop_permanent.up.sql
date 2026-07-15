CREATE TABLE IF NOT EXISTS pjp.route_pop_permanent (
    id SERIAL PRIMARY KEY,
    year BIGINT,
    week BIGINT,
    date TIMESTAMP,
    day VARCHAR(125),
    route_code BIGINT,
    pjp_id BIGINT,
    pjp_code BIGINT,
    cust_id VARCHAR(125) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_route_pop_permanent_pjp FOREIGN KEY (pjp_id) REFERENCES pjp.permanent_journey_plans (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_route_pop_permanent_route FOREIGN KEY (route_code) REFERENCES pjp.routes (route_code) ON UPDATE CASCADE ON DELETE CASCADE
);