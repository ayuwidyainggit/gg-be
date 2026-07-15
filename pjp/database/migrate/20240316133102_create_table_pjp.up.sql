CREATE TABLE IF NOT EXISTS pjp.permanent_journey_plans (
    id SERIAL PRIMARY KEY,
    pjp_code BIGINT UNIQUE NOT NULL,
    operation_type VARCHAR(125) NOT NULL,
    team_salesman VARCHAR(125),
    salesman_id BIGINT,
    salesman_name VARCHAR(125),
    warehouse_id BIGINT,
    warehouse_name VARCHAR(125),
    pjp_mode VARCHAR(125) DEFAULT 'manual',
    status VARCHAR(125) DEFAULT 'pending',
    cust_id VARCHAR(125) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);