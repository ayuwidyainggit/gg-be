-- Migration: Create Survey tables
-- Schema: mst
-- Date: 2025-12-23

-- Create m_survey table
CREATE TABLE IF NOT EXISTS mst.m_survey (
    survey_id SERIAL PRIMARY KEY,
    cust_id VARCHAR(10) NOT NULL,
    survey_title VARCHAR(150) NOT NULL,
    answer_frequency VARCHAR(20) NOT NULL,
    response_type VARCHAR(20) NOT NULL,
    target_type VARCHAR(20),
    emp_id INT,
    efective_date_start DATE,
    efective_date_end DATE,
    status INT DEFAULT 1,
    is_del BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
);

-- Create m_survey_area table
CREATE TABLE IF NOT EXISTS mst.m_survey_area (
    survey_area_id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL REFERENCES mst.m_survey(survey_id),
    distributor_id INT NOT NULL,
    area_id INT NOT NULL,
    is_del BOOLEAN DEFAULT FALSE
);

-- Create m_survey_salesman table
CREATE TABLE IF NOT EXISTS mst.m_survey_salesman (
    m_survey_salesman_id SERIAL PRIMARY KEY,
    cust_id VARCHAR(10) NOT NULL,
    survey_id INT NOT NULL REFERENCES mst.m_survey(survey_id),
    salesman_id INT NOT NULL,
    is_del BOOLEAN DEFAULT FALSE
);

-- Create m_survey_outlet table
CREATE TABLE IF NOT EXISTS mst.m_survey_outlet (
    survey_outlet_id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL REFERENCES mst.m_survey(survey_id),
    outlet_id INT NOT NULL,
    is_del BOOLEAN DEFAULT FALSE
);

-- Create m_survey_detail table (survey to template mapping)
CREATE TABLE IF NOT EXISTS mst.m_survey_detail (
    survey_detail_id SERIAL PRIMARY KEY,
    survey_id INT NOT NULL REFERENCES mst.m_survey(survey_id),
    survey_template_id INT NOT NULL,
    is_del BOOLEAN DEFAULT FALSE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_m_survey_cust_id ON mst.m_survey(cust_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_area_survey_id ON mst.m_survey_area(survey_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_area_distributor_id ON mst.m_survey_area(distributor_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_salesman_survey_id ON mst.m_survey_salesman(survey_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_salesman_cust_salesman ON mst.m_survey_salesman(cust_id, salesman_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_outlet_survey_id ON mst.m_survey_outlet(survey_id);
CREATE INDEX IF NOT EXISTS idx_m_survey_detail_survey_id ON mst.m_survey_detail(survey_id);
