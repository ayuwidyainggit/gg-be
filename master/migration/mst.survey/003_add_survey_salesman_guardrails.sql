-- Guardrail migration for survey salesman mapping.
-- Foreign key to mst.m_salesman cannot be added because remote table mst.m_salesman
-- uses composite primary key (cust_id, emp_id, created_at) and emp_id is not unique.
-- Validation is enforced in backend service layer.

CREATE INDEX IF NOT EXISTS idx_m_survey_salesman_cust_survey
    ON mst.m_survey_salesman(cust_id, survey_id);

CREATE INDEX IF NOT EXISTS idx_m_salesman_cust_emp_lookup
    ON mst.m_salesman(cust_id, emp_id)
    WHERE is_del = false AND deleted_at IS NULL;
