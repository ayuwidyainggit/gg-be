-- Migration 006: Create m_survey_distributor table
-- Schema: mst
-- Date: 2026-07-07
-- Ticket: SX-2445 / SX-2448 / SX-2452
-- Notes: Sprint 13 introduces explicit target Distributor mappings. cust_id
-- follows the existing tenant string convention used by m_survey_salesman.
-- The unique index is partial (is_del = false) so historical soft-deleted rows
-- do not block re-inserts.

CREATE TABLE IF NOT EXISTS mst.m_survey_distributor (
    m_survey_distributor_id SERIAL PRIMARY KEY,
    cust_id VARCHAR(10) NOT NULL,
    survey_id INT NOT NULL REFERENCES mst.m_survey(survey_id),
    distributor_id INT NOT NULL,
    is_del BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by BIGINT,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by BIGINT,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_m_survey_distributor_active
    ON mst.m_survey_distributor(survey_id, distributor_id)
    WHERE is_del = false;

CREATE INDEX IF NOT EXISTS idx_m_survey_distributor_survey_id
    ON mst.m_survey_distributor(survey_id);

CREATE INDEX IF NOT EXISTS idx_m_survey_distributor_cust_distributor
    ON mst.m_survey_distributor(cust_id, distributor_id);
