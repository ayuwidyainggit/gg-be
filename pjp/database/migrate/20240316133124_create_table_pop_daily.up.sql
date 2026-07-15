CREATE TABLE IF NOT EXISTS pjp.route_pop_daily (
    id SERIAL PRIMARY KEY,
    year BIGINT,
    week BIGINT,
    date TIMESTAMP,
    day VARCHAR(125),
    route_code BIGINT,
    pjp_id BIGINT,
    pjp_code BIGINT,
    parent_route BIGINT DEFAULT NULL,
    status VARCHAR(125) DEFAULT 'active',
    cust_id VARCHAR(125) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_route_pop_daily_pjp FOREIGN KEY (pjp_id) REFERENCES pjp.permanent_journey_plans (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_route_pop_daily_route FOREIGN KEY (route_code) REFERENCES pjp.routes (route_code) ON UPDATE CASCADE ON DELETE CASCADE,
    -- CONSTRAINT fk_route_pop_daily_route_outlet_additional FOREIGN KEY (parent_route) REFERENCES pjp.route_outlet_additional (route_code) ON UPDATE CASCADE ON DELETE CASCADE

);
