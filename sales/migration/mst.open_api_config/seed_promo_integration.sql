-- Seed: Promo integration Open API (whitelist cust_id C26009)
-- Run after create_open_api_tables.sql
-- Replace client_id / client_secret before production use.

INSERT INTO mst.open_api_config (
    system_integration,
    client_id,
    client_secret,
    environment,
    signature_algorithm,
    status,
    created_by
)
SELECT
    'Promo Integration',
    'promo-integration-client',
    'change-me-in-production',
    'DEV',
    'HMAC_SHA256',
    'A',
    'system'
WHERE NOT EXISTS (
    SELECT 1 FROM mst.open_api_config WHERE system_integration = 'Promo Integration'
);

-- IP whitelist not used by application logic (table retained for future use).

INSERT INTO mst.open_api_config_customer (open_api_config_id, cust_id, status)
SELECT c.id, 'C26009', 'A'
FROM mst.open_api_config c
WHERE c.system_integration = 'Promo Integration'
  AND NOT EXISTS (
      SELECT 1 FROM mst.open_api_config_customer cc
      WHERE cc.open_api_config_id = c.id AND cc.cust_id = 'C26009'
  );

INSERT INTO mst.open_api_endpoint (
    open_api_config_id,
    api_code,
    api_name,
    endpoint_url,
    method,
    api_type,
    is_active,
    description
)
SELECT
    c.id,
    'CREATE_PROMOTION_V2',
    'Create Promotion V2',
    '/open-api/v1/promotions',
    'POST',
    'INBOUND',
    TRUE,
    'Third-party create promotion (same logic as POST /v2/promotions)'
FROM mst.open_api_config c
WHERE c.system_integration = 'Promo Integration'
  AND NOT EXISTS (
      SELECT 1 FROM mst.open_api_endpoint e
      WHERE e.open_api_config_id = c.id AND e.api_code = 'CREATE_PROMOTION_V2'
  );
