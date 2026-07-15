-- Migration: Add Master Employee territory mapping
-- Schema: mst
-- Date: 2026-05-20
--
-- Scope columns are dynamic:
-- - ALL means the employee is allowed to access all current and future records.
-- - SELECTED means the employee is allowed only records listed in the mapping tables.

ALTER TABLE mst.m_employee
    ADD COLUMN IF NOT EXISTS region_scope varchar(10) NOT NULL DEFAULT 'ALL',
    ADD COLUMN IF NOT EXISTS area_scope varchar(10) NOT NULL DEFAULT 'ALL',
    ADD COLUMN IF NOT EXISTS distributor_scope varchar(10) NOT NULL DEFAULT 'ALL';

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'm_employee_region_scope_check'
        AND conrelid = 'mst.m_employee'::regclass
    ) THEN
        ALTER TABLE mst.m_employee
            ADD CONSTRAINT m_employee_region_scope_check
            CHECK (region_scope IN ('ALL', 'SELECTED'));
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'm_employee_area_scope_check'
        AND conrelid = 'mst.m_employee'::regclass
    ) THEN
        ALTER TABLE mst.m_employee
            ADD CONSTRAINT m_employee_area_scope_check
            CHECK (area_scope IN ('ALL', 'SELECTED'));
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'm_employee_distributor_scope_check'
        AND conrelid = 'mst.m_employee'::regclass
    ) THEN
        ALTER TABLE mst.m_employee
            ADD CONSTRAINT m_employee_distributor_scope_check
            CHECK (distributor_scope IN ('ALL', 'SELECTED'));
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS mst.m_employee_region_mapping (
    employee_region_mapping_id bigserial PRIMARY KEY,
    cust_id varchar(10) NOT NULL,
    emp_id int4 NOT NULL,
    region_id int4 NOT NULL,
    created_by int8 NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int8 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_m_employee_region_mapping_active
    ON mst.m_employee_region_mapping (cust_id, emp_id, region_id)
    WHERE is_del = false;

CREATE INDEX IF NOT EXISTS idx_m_employee_region_mapping_emp
    ON mst.m_employee_region_mapping (cust_id, emp_id)
    WHERE is_del = false;

CREATE INDEX IF NOT EXISTS idx_m_employee_region_mapping_region
    ON mst.m_employee_region_mapping (cust_id, region_id)
    WHERE is_del = false;

CREATE TABLE IF NOT EXISTS mst.m_employee_area_mapping (
    employee_area_mapping_id bigserial PRIMARY KEY,
    cust_id varchar(10) NOT NULL,
    emp_id int4 NOT NULL,
    area_id int4 NOT NULL,
    created_by int8 NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int8 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_m_employee_area_mapping_active
    ON mst.m_employee_area_mapping (cust_id, emp_id, area_id)
    WHERE is_del = false;

CREATE INDEX IF NOT EXISTS idx_m_employee_area_mapping_emp
    ON mst.m_employee_area_mapping (cust_id, emp_id)
    WHERE is_del = false;

CREATE INDEX IF NOT EXISTS idx_m_employee_area_mapping_area
    ON mst.m_employee_area_mapping (cust_id, area_id)
    WHERE is_del = false;

CREATE TABLE IF NOT EXISTS mst.m_employee_distributor_mapping (
    employee_distributor_mapping_id bigserial PRIMARY KEY,
    cust_id varchar(10) NOT NULL,
    emp_id int4 NOT NULL,
    distributor_id int4 NOT NULL,
    created_by int8 NULL,
    created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by int8 NULL,
    updated_at timestamptz(6) NULL,
    deleted_by int8 NULL,
    deleted_at timestamptz(6) NULL,
    is_del bool NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_m_employee_distributor_mapping_active
    ON mst.m_employee_distributor_mapping (cust_id, emp_id, distributor_id)
    WHERE is_del = false;

CREATE INDEX IF NOT EXISTS idx_m_employee_distributor_mapping_emp
    ON mst.m_employee_distributor_mapping (cust_id, emp_id)
    WHERE is_del = false;

CREATE INDEX IF NOT EXISTS idx_m_employee_distributor_mapping_distributor
    ON mst.m_employee_distributor_mapping (cust_id, distributor_id)
    WHERE is_del = false;
