ALTER TABLE pjp.permanent_journey_plans
ADD COLUMN approval_status VARCHAR(32) NOT NULL DEFAULT 'Draft';