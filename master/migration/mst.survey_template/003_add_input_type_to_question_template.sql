-- Migration: SX-XXX - Add input_type column to mst.question_template
-- Purpose: Persist survey question input type based on Template Survei docs
-- Date: 2026-03-24

ALTER TABLE mst.question_template
ADD COLUMN IF NOT EXISTS input_type VARCHAR(20);

UPDATE mst.question_template
SET input_type = 'textfield'
WHERE input_type IS NULL;

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'question_template_input_type_check'
			AND conrelid = 'mst.question_template'::regclass
	) THEN
		ALTER TABLE mst.question_template
		ADD CONSTRAINT question_template_input_type_check
		CHECK (input_type IN ('textfield', 'dropdown', 'radiobutton', 'toggle', 'checkbox'));
	END IF;
END $$;

ALTER TABLE mst.question_template
ALTER COLUMN input_type SET NOT NULL;