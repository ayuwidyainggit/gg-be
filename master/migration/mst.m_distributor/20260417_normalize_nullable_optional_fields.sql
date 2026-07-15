-- Normalize optional distributor fields to database NULL
-- This keeps legacy rows consistent with PATCH semantics for optional fields.

UPDATE mst.m_distributor
SET
    barcode = NULLIF(BTRIM(barcode), ''),
    province_id = NULLIF(BTRIM(province_id), ''),
    regency_id = NULLIF(BTRIM(regency_id), ''),
    sub_district_id = NULLIF(BTRIM(sub_district_id), ''),
    ward_id = NULLIF(BTRIM(ward_id), ''),
    zip_code = NULLIF(BTRIM(zip_code), ''),
    phone = NULLIF(BTRIM(phone), ''),
    fax_number = NULLIF(BTRIM(fax_number), ''),
    parent_cust_id = NULLIF(BTRIM(parent_cust_id), '')
WHERE
    barcode IS NOT NULL
    OR province_id IS NOT NULL
    OR regency_id IS NOT NULL
    OR sub_district_id IS NOT NULL
    OR ward_id IS NOT NULL
    OR zip_code IS NOT NULL
    OR phone IS NOT NULL
    OR fax_number IS NOT NULL
    OR parent_cust_id IS NOT NULL;

WITH parent_candidates AS (
    SELECT
        d.distributor_id,
        CASE
            WHEN LENGTH(d.cust_id) > 6 THEN LEFT(d.cust_id, 6)
            ELSE d.cust_id
        END AS parent_candidate
    FROM mst.m_distributor d
    WHERE d.parent_cust_id IS NULL
)
UPDATE mst.m_distributor d
SET parent_cust_id = pc.parent_candidate
FROM parent_candidates pc
JOIN smc.m_customer mc ON mc.cust_id = pc.parent_candidate
WHERE d.distributor_id = pc.distributor_id;
