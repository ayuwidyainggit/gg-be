DROP INDEX IF EXISTS mst.idx_m_work_day_working_day_calendar_date;
DROP INDEX IF EXISTS mst.idx_m_week_working_day_calendar_calendar_week;
DROP INDEX IF EXISTS mst.idx_m_week_working_day_calendar_week;
DROP INDEX IF EXISTS mst.ux_working_day_calendar_holiday_scope_date;
DROP INDEX IF EXISTS mst.idx_working_day_calendar_holiday_calendar_date;
DROP INDEX IF EXISTS mst.idx_working_day_calendar_latest;
DROP INDEX IF EXISTS mst.idx_working_day_calendar_cust_date_range;

DO $$
DECLARE
    constraint_record RECORD;
BEGIN
    FOR constraint_record IN
        SELECT format('%I.%I', n.nspname, c.relname) AS table_name, con.conname
        FROM pg_constraint con
        INNER JOIN pg_class c ON c.oid = con.conrelid
        INNER JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = 'mst'
            AND (
                (
                    (c.relname = 'm_work_day' OR c.relname LIKE 'm\_work\_day\_%' ESCAPE '\')
                    AND (
                        con.conname LIKE 'fk_m_work_day_working_day_calendar%'
                        OR con.conname LIKE 'chk_m_work_day_holiday_source%'
                    )
                )
                OR (
                    (c.relname = 'm_week' OR c.relname LIKE 'm\_week\_%' ESCAPE '\')
                    AND (
                        con.conname LIKE 'fk_m_week_working_day_calendar%'
                        OR con.conname LIKE 'chk_m_week_calendar_week_no%'
                    )
                )
            )
    LOOP
        EXECUTE format(
            'ALTER TABLE %s DROP CONSTRAINT IF EXISTS %I',
            constraint_record.table_name,
            constraint_record.conname
        );
    END LOOP;
END
$$;

ALTER TABLE mst.m_work_day
    DROP COLUMN IF EXISTS holiday_note,
    DROP COLUMN IF EXISTS holiday_source,
    DROP COLUMN IF EXISTS working_day_calendar_id;

ALTER TABLE mst.m_week
    DROP COLUMN IF EXISTS calendar_week_no,
    DROP COLUMN IF EXISTS working_day_calendar_id;

DROP TABLE IF EXISTS mst.working_day_calendar_holiday;
DROP TABLE IF EXISTS mst.working_day_calendar;
