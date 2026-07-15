-- Migration: SX-828 - Add use_image column to mst.question_template
-- Purpose: Move use_image from template level to question level
-- Date: 2026-02-09

-- Up migration: Add use_image column to question_template
ALTER TABLE mst.question_template 
ADD COLUMN IF NOT EXISTS use_image BOOLEAN NOT NULL DEFAULT false;

-- Migrate existing data: Copy use_image value from template to all its questions
UPDATE mst.question_template qt
SET use_image = st.use_image
FROM mst.m_survey_template st
WHERE qt.survey_template_id = st.survey_template_id
AND qt.is_del = false;

-- Note: Column use_image in mst.m_survey_template is NOT removed for backward compatibility
-- It will be ignored by the application going forward
