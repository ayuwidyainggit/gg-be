BEGIN;

ALTER TABLE acf.expense_type
    ADD COLUMN IF NOT EXISTS cust_id varchar(50);

WITH owner_candidates AS (
    SELECT
        et.expense_type_id,
        MIN(cust.parent_cust_id) AS parent_cust_id,
        COUNT(DISTINCT cust.parent_cust_id) AS owner_count
    FROM acf.expense_type et
    JOIN sys.m_user us ON us.user_id = et.created_by AND us.is_del = false
    JOIN smc.m_customer cust ON cust.cust_id = us.cust_id AND cust.is_del = false
    WHERE et.cust_id IS NULL
    GROUP BY et.expense_type_id
)
UPDATE acf.expense_type et
SET cust_id = owner_candidates.parent_cust_id
FROM owner_candidates
WHERE et.expense_type_id = owner_candidates.expense_type_id
  AND owner_candidates.owner_count = 1;

CREATE INDEX IF NOT EXISTS idx_expense_type_cust_id
    ON acf.expense_type(cust_id);

CREATE INDEX IF NOT EXISTS idx_expense_type_cust_active_del
    ON acf.expense_type(cust_id, is_active, is_del);

COMMIT;
