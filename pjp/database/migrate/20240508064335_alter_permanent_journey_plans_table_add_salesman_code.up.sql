ALTER TABLE pjp.permanent_journey_plans
ADD COLUMN IF NOT EXISTS salesman_code VARCHAR(125);
