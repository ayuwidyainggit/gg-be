-- Fix inverted UOM hierarchy on product-mapping rows created before column mapping correction.
-- Master product convention: unit_id1 = smallest, unit_id2 = middle, unit_id3 = largest.
UPDATE mst.m_product
SET
    unit_id1 = unit_id3,
    unit_id3 = unit_id1,
    updated_at = CURRENT_TIMESTAMP
WHERE COALESCE(is_product_mapping, false) = true
    AND COALESCE(is_del, false) = false
    AND (
        COALESCE(unit_id1, '') <> ''
        OR COALESCE(unit_id3, '') <> ''
    );
