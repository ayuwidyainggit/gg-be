-- Migration: Create Survey Template tables
-- Schema: mst
-- Date: 2025-12-23

-- Create m_survey_template table
CREATE TABLE IF NOT EXISTS mst.m_survey_template (
    survey_template_id SERIAL PRIMARY KEY,
    cust_id VARCHAR(10) NOT NULL,
    template_code VARCHAR(10) NOT NULL,
    template_title VARCHAR(150) NOT NULL,
    question_total INT DEFAULT 0,
    use_image BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    is_del BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
);

-- Create question_template table
CREATE TABLE IF NOT EXISTS mst.question_template (
    question_template_id SERIAL PRIMARY KEY,
    survey_template_id INT NOT NULL REFERENCES mst.m_survey_template(survey_template_id),
    question VARCHAR(225) NOT NULL,
    input_type VARCHAR(20) NOT NULL CHECK (input_type IN ('textfield', 'dropdown', 'radiobutton', 'toggle', 'checkbox')),
    answer_type VARCHAR(20) NOT NULL CHECK (answer_type IN ('Single', 'Multiple', 'Free Text')),
    seq INT DEFAULT 0,
    is_del BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
);

-- Create m_q_option_template table
CREATE TABLE IF NOT EXISTS mst.m_q_option_template (
    q_option_template_id SERIAL PRIMARY KEY,
    question_template_id INT NOT NULL REFERENCES mst.question_template(question_template_id),
    option VARCHAR(225) NOT NULL,
    seq INT DEFAULT 0,
    is_del BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_m_survey_template_cust_id ON mst.m_survey_template(cust_id);
CREATE INDEX IF NOT EXISTS idx_question_template_survey_id ON mst.question_template(survey_template_id);
CREATE INDEX IF NOT EXISTS idx_m_q_option_template_question_id ON mst.m_q_option_template(question_template_id);
