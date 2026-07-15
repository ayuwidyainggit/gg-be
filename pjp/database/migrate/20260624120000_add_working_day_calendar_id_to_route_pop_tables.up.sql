ALTER TABLE pjp.route_pop_permanent
ADD COLUMN IF NOT EXISTS working_day_calendar_id BIGINT;

ALTER TABLE pjp.route_pop_daily
ADD COLUMN IF NOT EXISTS working_day_calendar_id BIGINT;
