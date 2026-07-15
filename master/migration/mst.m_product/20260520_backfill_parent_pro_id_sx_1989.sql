-- SX-1989: Backfill parent_pro_id for legacy distributor child rows
-- that have a clear 1:1 unambiguous mapping to a principal parent product
-- based on matching pro_code within the same parent_cust_id scope.
--
-- Safe conditions:
--   - child.distributor_id IS NOT NULL (distributor-owned row)
--   - COALESCE(child.parent_pro_id, 0) = 0 (not yet linked)
--   - child.is_del = false
--   - exactly 1 parent candidate per child (HAVING COUNT = 1)
--   - parent.distributor_id IS NULL (principal-owned row)
--   - parent.is_del = false
--
-- Rows with ambiguous mapping (multiple parent candidates) are skipped
-- and must be resolved manually.
--
-- Idempotent: re-running only touches rows still with parent_pro_id = 0/NULL.

WITH parent_lookup AS (
    SELECT
        parent.pro_code,
        parent.pro_id   AS parent_pro_id,
        parent.cust_id  AS parent_cust_id
    FROM mst.m_product parent
    WHERE parent.distributor_id IS NULL
      AND parent.is_del = false
),
child_candidates AS (
    SELECT
        child.pro_id,
        child.cust_id,
        child.pro_code,
        parent_lookup.parent_pro_id,
        COUNT(DISTINCT parent_lookup.parent_pro_id) AS parent_count
    FROM mst.m_product child
    JOIN mst.m_distributor d
        ON d.cust_id        = child.cust_id
       AND d.distributor_id = child.distributor_id
    JOIN parent_lookup
        ON parent_lookup.parent_cust_id = d.parent_cust_id
       AND parent_lookup.pro_code       = child.pro_code
    WHERE child.distributor_id IS NOT NULL
      AND COALESCE(child.parent_pro_id, 0) = 0
      AND child.is_del = false
    GROUP BY child.pro_id, child.cust_id, child.pro_code, parent_lookup.parent_pro_id
    HAVING COUNT(DISTINCT parent_lookup.parent_pro_id) = 1
)
UPDATE mst.m_product child
   SET parent_pro_id = child_candidates.parent_pro_id,
       updated_at    = CURRENT_TIMESTAMP
  FROM child_candidates
 WHERE child.pro_id  = child_candidates.pro_id
   AND child.cust_id = child_candidates.cust_id;
