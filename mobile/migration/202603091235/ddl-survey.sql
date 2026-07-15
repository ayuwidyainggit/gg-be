CREATE SEQUENCE mst.survey_answer_id_seq
    START 1
    INCREMENT 1
    MINVALUE 1
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE mst.survey_answer_detail_id_seq
    START 1
    INCREMENT 1
    MINVALUE 1
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE mst.survey_answer_file_id_seq
    START 1
    INCREMENT 1
    MINVALUE 1
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE mst.survey_answer_option_id_seq
    START 1
    INCREMENT 1
    MINVALUE 1
    NO MAXVALUE
    CACHE 1;

CREATE TABLE mst.survey_answer (
    cust_id VARCHAR(10) NOT NULL,
    survey_answer_id INT8 PRIMARY KEY DEFAULT nextval('mst.survey_answer_id_seq'::regclass),
    survey_template_id INT8 NOT NULL, -- FK to mst.m_survey_template
    survey_id INT8 NOT NULL, -- FK to mst.m_survey
    emp_id INT8 NOT NULL, -- FK to mst.m_salesman
    outlet_id INT8 NOT NULL, -- FK to mst.m_outlet
    area_id INT8 NOT NULL, -- FK to mst.m_area
    answer_date TIMESTAMP DEFAULT CURRENT_DATE,
    status VARCHAR(20) DEFAULT 'Submitted', -- Enum: Draft, Submitted, Cancelled
    created_by INT8 NOT NULL,
    created_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    updated_by INT8 NOT NULL,
    updated_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    is_del BOOLEAN DEFAULT FALSE,
    deleted_by INT8 NULL,
    deleted_at TIMESTAMPTZ(6) NULL,
    
    -- Foreign Key Constraints
    CONSTRAINT fk_survey_template FOREIGN KEY (survey_template_id) REFERENCES mst.m_survey_template(survey_template_id),
    CONSTRAINT fk_survey FOREIGN KEY (survey_id) REFERENCES mst.m_survey(survey_id),
);

CREATE TABLE mst.survey_answer_detail (
    cust_id VARCHAR(10) NOT NULL,
    survey_answer_detail_id INT8 PRIMARY KEY DEFAULT nextval('mst.survey_answer_detail_id_seq'::regclass),
    survey_answer_id INT8 NOT NULL, -- FK to mst.survey_answer
    question_template_id INT8 NOT NULL, -- FK to mst.m_question_template
    input_type VARCHAR(225) NOT NULL, -- e.g., textfield, dropdown
    answer_type VARCHAR(20) NOT NULL, -- Enum: Single, Multiple, Free Text
    seq INT4 NOT NULL, -- Question sequence
    is_answered BOOLEAN DEFAULT FALSE,
    free_text_answer TEXT NULL,
    photo_path VARCHAR(255) NULL,
    created_by INT8 NOT NULL,
    created_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    updated_by INT8 NOT NULL,
    updated_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    is_del BOOLEAN DEFAULT FALSE,
    deleted_by INT8 NULL,
    deleted_at TIMESTAMPTZ(6) NULL,
    
    -- Foreign Key Constraints
    CONSTRAINT fk_survey_answer FOREIGN KEY (survey_answer_id) 
        REFERENCES mst.survey_answer(survey_answer_id),
    CONSTRAINT fk_question_template FOREIGN KEY (question_template_id) 
        REFERENCES mst.question_template(question_template_id)
);

CREATE TABLE mst.survey_answer_files (
    cust_id VARCHAR(10) NOT NULL,
    survey_answer_files INT8 PRIMARY KEY DEFAULT nextval('mst.survey_answer_file_id_seq'::regclass),
    survey_answer_detail_id INT8 NOT NULL, -- FK to mst.survey_answer_detail
    file_name VARCHAR(255) NOT NULL,
    file_data BYTEA,
    file_key VARCHAR(10) NOT NULL, -- Enum: 'image', 'video'
    media_category TEXT NOT NULL, -- PNG, JPG, JPEG, etc.
    file_size BIGINT,

    -- Foreign Key Constraints
    CONSTRAINT fk_survey_detail_file FOREIGN KEY (survey_answer_detail_id) 
        REFERENCES mst.survey_answer_detail(survey_answer_detail_id)
);

CREATE TABLE mst.survey_answer_option (
    cust_id VARCHAR(10) NOT NULL,
    survey_answer_option_id INT8 PRIMARY KEY DEFAULT nextval('mst.survey_answer_option_id_seq'::regclass),
    survey_answer_detail_id INT8 NOT NULL, -- FK to mst.survey_answer_detail
    q_option_template_id INT8 NOT NULL, -- FK to mst.m_q_option_template
    option_label VARCHAR(225),
    created_by INT8 NOT NULL,
    created_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    updated_by INT8 NOT NULL,
    updated_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    is_del BOOLEAN DEFAULT FALSE,
    deleted_by INT8 NULL,
    deleted_at TIMESTAMPTZ(6) NULL,

    -- Foreign Key Constraints
    CONSTRAINT fk_survey_detail_opt FOREIGN KEY (survey_answer_detail_id) 
        REFERENCES mst.survey_answer_detail(survey_answer_detail_id),
    CONSTRAINT fk_q_option_template FOREIGN KEY (q_option_template_id) 
        REFERENCES mst.m_q_option_template(q_option_template_id)
);
