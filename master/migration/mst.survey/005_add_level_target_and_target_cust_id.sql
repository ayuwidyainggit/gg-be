-- Migration 005: Add level_target to m_survey and target_cust_id to m_survey_area
-- Schema: mst
-- Date: 2026-07-07
-- Ticket: SX-2445 / SX-2448 / SX-2452
-- Notes: Both columns are idempotent (IF NOT EXISTS) and nullable to remain
-- backward compatible with rows written by releases prior to Sprint 13.

ALTER TABLE mst.m_survey
    ADD COLUMN IF NOT EXISTS level_target VARCHAR(20);

ALTER TABLE mst.m_survey_area
    ADD COLUMN IF NOT EXISTS target_cust_id VARCHAR(10);

-- Backfill safety: leave existing rows NULL. Downstream code treats NULL as
-- "legacy payload" so we do not invent a value that the original payload did
-- not provide.

CREATE INDEX IF NOT EXISTS idx_m_survey_level_target ON mst.m_survey(level_target);
CREATE INDEX IF NOT EXISTS idx_m_survey_area_target_cust_id ON mst.m_survey_area(target_cust_id);
