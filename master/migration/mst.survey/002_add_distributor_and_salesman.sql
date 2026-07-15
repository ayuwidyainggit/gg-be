ALTER TABLE mst.m_survey_area
    ADD COLUMN IF NOT EXISTS distributor_id INT;

UPDATE mst.m_survey_area
SET distributor_id = 0
WHERE distributor_id IS NULL;

ALTER TABLE mst.m_survey_area
    ALTER COLUMN distributor_id SET NOT NULL;

CREATE TABLE IF NOT EXISTS mst.m_survey_salesman (
    m_survey_salesman_id SERIAL PRIMARY KEY,
    cust_id VARCHAR(10) NOT NULL,
    survey_id INT NOT NULL REFERENCES mst.m_survey(survey_id),
    salesman_id INT NOT NULL,
    is_del BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_m_survey_area_distributor_id ON mst.m_survey_area(distributor_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_salesman_survey_id ON mst.m_survey_salesman(survey_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_salesman_cust_salesman ON mst.m_survey_salesman(cust_id, salesman_id);
