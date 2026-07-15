-- ============================================
-- Leave Request Migration
-- Description: Create table for leave request feature
-- ============================================

BEGIN;

CREATE SCHEMA IF NOT EXISTS mobile;

CREATE TABLE IF NOT EXISTS mobile.leave_request (
    leave_id     bigserial PRIMARY KEY,
    cust_id      varchar(10) NOT NULL,
    emp_id       int4 NOT NULL,
    start_date   date NOT NULL,
    end_date     date NOT NULL,
    reason       varchar(500) NOT NULL,
    file_url     varchar(500) NOT NULL,
    file_name    varchar(255) NOT NULL,
    approval     varchar(20) NOT NULL DEFAULT 'Approved',
    created_by   int4 NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved_by  int4 NULL,
    approved_at  timestamptz NULL,
    canceled_by  int4 NULL,
    canceled_at  timestamptz NULL
);

CREATE INDEX IF NOT EXISTS idx_leave_request_cust_emp_dates
    ON mobile.leave_request (cust_id, emp_id, start_date, end_date);

CREATE INDEX IF NOT EXISTS idx_leave_request_cust_emp_created
    ON mobile.leave_request (cust_id, emp_id, created_at DESC);

COMMIT;
