ALTER TABLE pjp.route_pop_daily
DROP COLUMN IF EXISTS working_day_calendar_id;

ALTER TABLE pjp.route_pop_permanent
DROP COLUMN IF EXISTS working_day_calendar_id;
