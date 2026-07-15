-- Open API configuration tables for third-party integrations (Promo integration)
-- Schema: mst

CREATE TABLE IF NOT EXISTS mst.open_api_config (
    id BIGSERIAL PRIMARY KEY,
    system_integration VARCHAR(100) NOT NULL,
    client_id VARCHAR(255),
    client_secret TEXT,
    environment VARCHAR(20) NOT NULL,
    base_url VARCHAR(255),
    signature_algorithm VARCHAR(50),
    status CHAR(1) NOT NULL DEFAULT 'A',
    created_by VARCHAR(50) NOT NULL,
    created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(50),
    updated_date TIMESTAMP
);

COMMENT ON TABLE mst.open_api_config IS 'Master config for third-party Open API integrations';
COMMENT ON COLUMN mst.open_api_config.status IS 'A=Active, I=Inactive';

CREATE TABLE IF NOT EXISTS mst.open_api_config_ip (
    id BIGSERIAL PRIMARY KEY,
    open_api_config_id BIGINT NOT NULL REFERENCES mst.open_api_config(id) ON DELETE CASCADE,
    ip_address VARCHAR(50) NOT NULL,
    status CHAR(1) NOT NULL DEFAULT 'A'
);

COMMENT ON TABLE mst.open_api_config_ip IS 'IP whitelist for Open API access';
COMMENT ON COLUMN mst.open_api_config_ip.status IS 'A=Active, I=Inactive';

CREATE TABLE IF NOT EXISTS mst.open_api_config_customer (
    id BIGSERIAL PRIMARY KEY,
    open_api_config_id BIGINT NOT NULL REFERENCES mst.open_api_config(id) ON DELETE CASCADE,
    cust_id VARCHAR(50) NOT NULL,
    status CHAR(1) NOT NULL DEFAULT 'A',
    CONSTRAINT uq_open_api_config_customer UNIQUE (open_api_config_id, cust_id)
);

COMMENT ON TABLE mst.open_api_config_customer IS 'Customer whitelist per Open API integration';
COMMENT ON COLUMN mst.open_api_config_customer.cust_id IS 'FK reference: sys.m_customer.cust_id';

CREATE TABLE IF NOT EXISTS mst.open_api_endpoint (
    id BIGSERIAL PRIMARY KEY,
    open_api_config_id BIGINT NOT NULL REFERENCES mst.open_api_config(id) ON DELETE CASCADE,
    api_code VARCHAR(100) NOT NULL,
    api_name VARCHAR(255) NOT NULL,
    endpoint_url VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    api_type VARCHAR(20) NOT NULL,
    timeout_second INT DEFAULT 30,
    retry_count INT DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,
    CONSTRAINT uq_open_api_endpoint UNIQUE (open_api_config_id, api_code)
);

COMMENT ON TABLE mst.open_api_endpoint IS 'Registered endpoints per Open API integration';
COMMENT ON COLUMN mst.open_api_endpoint.api_type IS 'INBOUND or OUTBOUND';

CREATE INDEX IF NOT EXISTS idx_open_api_config_client_id ON mst.open_api_config(client_id) WHERE status = 'A';
CREATE INDEX IF NOT EXISTS idx_open_api_config_ip_config ON mst.open_api_config_ip(open_api_config_id) WHERE status = 'A';
CREATE INDEX IF NOT EXISTS idx_open_api_config_customer_cust ON mst.open_api_config_customer(cust_id) WHERE status = 'A';
