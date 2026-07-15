-- Migration: sync demo schema toward staging (strict non-destructive/no-DML variant)
-- Generated from schema-only staging/demo comparison. No row payload is copied or modified.
-- Destructive DDL intentionally excluded; demo-only objects remain as residual drift.

BEGIN;

CREATE SCHEMA IF NOT EXISTS promo;

-- Missing sequences required before table defaults
CREATE SEQUENCE IF NOT EXISTS "inv"."replenishment_order_approval_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "mst"."m_survey_salesman_m_survey_salesman_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "mst"."product_ripening_history_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "mst"."survey_answer_detail_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "mst"."survey_answer_file_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "mst"."survey_answer_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "mst"."survey_answer_option_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."destinations_additional_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."destinations_history_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."destinations_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."outlet_visit_list_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."permanent_journey_plans_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."route_pop_dailies_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."route_pop_permanent_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "pjp_principles"."routes_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "promo"."mongodb_objectid_counter" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "public"."mongodb_objectid_counter" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "sys"."password_reset_requests_id_seq" START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

-- Required extension and shared sequence dependencies
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';

-- Name: mongodb_objectid_counter; Type: SEQUENCE; Schema: promo; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS promo.mongodb_objectid_counter
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;




-- Promo schema/types/tables/indexes/constraints excluding procedural bodies
--
-- PostgreSQL database dump
--


-- Dumped from database version 14.5 (Ubuntu 14.5-1.pgdg20.04+1)
-- Dumped by pg_dump version 18.3 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: promo; Type: SCHEMA; Schema: -; Owner: -
--

--
-- Name: budget_ref_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.budget_ref_type AS ENUM (
    'unlimited',
    'limited'
);


--
-- Name: claim_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.claim_type AS ENUM (
    'full',
    'partial'
);


--
-- Name: control_level; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.control_level AS ENUM (
    'region',
    'area',
    'distributor',
    'salesman'
);


--
-- Name: coverage_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.coverage_type AS ENUM (
    'national',
    'by_distributor'
);


--
-- Name: creation_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.creation_type AS ENUM (
    'new',
    'replacement'
);


--
-- Name: outlet_sel_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.outlet_sel_type AS ENUM (
    'by_outlet',
    'by_attribute'
);


--
-- Name: promotion_status; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.promotion_status AS ENUM (
    'draft',
    'submit',
    'approved',
    'rejected',
    'inactive',
    'active',
    'closed'
);


--
-- Name: promotion_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.promotion_type AS ENUM (
    'slab',
    'strata'
);


--
-- Name: reward_cap_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.reward_cap_type AS ENUM (
    'amount',
    'qty'
);


--
-- Name: reward_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.reward_type AS ENUM (
    'percentage',
    'fixed_value',
    'product'
);


--
-- Name: rule_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.rule_type AS ENUM (
    'quantity',
    'value'
);


--
-- Name: uom_type; Type: TYPE; Schema: promo; Owner: -
--

CREATE TYPE promo.uom_type AS ENUM (
    'smallest',
    'middle',
    'largest'
);






--
-- Name: mongodb_objectid_counter; Type: SEQUENCE; Schema: promo; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS promo.mongodb_objectid_counter
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: promotion_coverage_distributors; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_coverage_distributors (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    promo_id character varying(50) NOT NULL,
    distributor_id bigint NOT NULL
);


--
-- Name: promotion_outlet_attribute_class; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_class (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    criteria_id character varying(30) NOT NULL,
    outlet_class_id bigint NOT NULL
);


--
-- Name: promotion_outlet_attribute_group; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_group (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    criteria_id character varying(30) NOT NULL,
    outlet_group_id bigint NOT NULL
);


--
-- Name: promotion_outlet_attribute_sales_team; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_sales_team (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    criteria_id character varying(30) NOT NULL,
    sales_team_id bigint NOT NULL
);


--
-- Name: promotion_outlet_attribute_type; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_type (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    criteria_id character varying(30) NOT NULL,
    outlet_type_id bigint NOT NULL
);


--
-- Name: promotion_outlet_criteria; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_outlet_criteria (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    promo_id character varying(50) NOT NULL,
    selection_type promo.outlet_sel_type DEFAULT 'by_attribute'::promo.outlet_sel_type NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: promotion_outlets_selected; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_outlets_selected (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    criteria_id character varying(30) NOT NULL,
    outlet_id bigint NOT NULL
);


--
-- Name: promotion_product_criteria; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_product_criteria (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    promo_id character varying(50) NOT NULL,
    pro_id bigint NOT NULL,
    mandatory boolean DEFAULT false NOT NULL,
    min_buy_type promo.rule_type,
    min_buy_qty numeric(20,4),
    min_buy_value numeric(20,4),
    min_buy_uom promo.uom_type,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT promotion_product_criteria_check CHECK ((((min_buy_type IS NULL) AND (min_buy_qty IS NULL) AND (min_buy_value IS NULL)) OR ((min_buy_type = 'quantity'::promo.rule_type) AND (min_buy_qty IS NOT NULL) AND (min_buy_qty >= (0)::numeric) AND (min_buy_value IS NULL)) OR ((min_buy_type = 'value'::promo.rule_type) AND (min_buy_value IS NOT NULL) AND (min_buy_value >= (0)::numeric) AND (min_buy_qty IS NULL))))
);


--
-- Name: promotion_reward_products; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_reward_products (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    promo_id character varying(50) NOT NULL,
    pro_id bigint NOT NULL,
    ordinal integer NOT NULL
);


--
-- Name: promotion_slabs; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_slabs (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    promo_id character varying(50) NOT NULL,
    ordinal integer NOT NULL,
    description character varying(50),
    rule_type promo.rule_type NOT NULL,
    rule_uom promo.uom_type,
    range_from numeric(20,4),
    range_to numeric(20,4) NOT NULL,
    reward_type promo.reward_type NOT NULL,
    reward_value numeric(20,4),
    reward_uom promo.uom_type,
    per_scope character varying(16),
    CONSTRAINT promotion_slabs_check CHECK ((range_to > COALESCE(range_from, (0)::numeric))),
    CONSTRAINT promotion_slabs_check1 CHECK (((NOT (reward_type = 'percentage'::promo.reward_type)) OR ((reward_value IS NOT NULL) AND ((reward_value >= (1)::numeric) AND (reward_value <= (100)::numeric))))),
    CONSTRAINT promotion_slabs_check2 CHECK (((NOT (reward_type = 'fixed_value'::promo.reward_type)) OR ((reward_value IS NOT NULL) AND (reward_value > (0)::numeric)))),
    CONSTRAINT promotion_slabs_check3 CHECK (((NOT (reward_type = 'product'::promo.reward_type)) OR ((reward_value IS NOT NULL) AND (reward_value > (0)::numeric)))),
    CONSTRAINT promotion_slabs_per_scope_check CHECK (((per_scope IS NULL) OR ((per_scope)::text = ANY ((ARRAY['per_product'::character varying, 'per_order'::character varying])::text[]))))
);


--
-- Name: promotion_strata; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotion_strata (
    cust_id character varying(10) NOT NULL,
    id character varying(30) NOT NULL,
    promo_id character varying(50) NOT NULL,
    ordinal integer NOT NULL,
    description character varying(50),
    rule_type promo.rule_type NOT NULL,
    rule_uom promo.uom_type,
    range_from numeric(20,4) NOT NULL,
    range_to numeric(20,4) NOT NULL,
    reward_type promo.reward_type NOT NULL,
    reward_value numeric(20,4),
    reward_uom promo.uom_type,
    per_scope character varying(16),
    claimable boolean DEFAULT false NOT NULL,
    claim_realization_pct numeric(5,2),
    CONSTRAINT promotion_strata_check CHECK ((range_to > range_from)),
    CONSTRAINT promotion_strata_check1 CHECK (((NOT (reward_type = 'percentage'::promo.reward_type)) OR ((reward_value IS NOT NULL) AND ((reward_value >= (0)::numeric) AND (reward_value <= (100)::numeric))))),
    CONSTRAINT promotion_strata_check2 CHECK (((NOT (reward_type = 'fixed_value'::promo.reward_type)) OR ((reward_value IS NOT NULL) AND (reward_value >= (0)::numeric)))),
    CONSTRAINT promotion_strata_check3 CHECK (((NOT (reward_type = 'product'::promo.reward_type)) OR (reward_value IS NOT NULL))),
    CONSTRAINT promotion_strata_claim_realization_pct_check CHECK (((claim_realization_pct IS NULL) OR ((claim_realization_pct >= (0)::numeric) AND (claim_realization_pct <= (100)::numeric)))),
    CONSTRAINT promotion_strata_ordinal_check CHECK (((ordinal >= 1) AND (ordinal <= 5))),
    CONSTRAINT promotion_strata_per_scope_check CHECK (((per_scope IS NULL) OR ((per_scope)::text = ANY ((ARRAY['per_product'::character varying, 'per_order'::character varying])::text[]))))
);


--
-- Name: promotions; Type: TABLE; Schema: promo; Owner: -
--

CREATE TABLE IF NOT EXISTS promo.promotions (
    cust_id character varying(10) NOT NULL,
    promo_id character varying(50) NOT NULL,
    promo_desc character varying(100) NOT NULL,
    promo_type promo.promotion_type NOT NULL,
    promo_creation_type promo.creation_type NOT NULL,
    existing_promo_id character varying(50),
    promo_status promo.promotion_status DEFAULT 'draft'::promo.promotion_status NOT NULL,
    is_budget_reference boolean DEFAULT false NOT NULL,
    budget_ref_type promo.budget_ref_type,
    budget_reference_id integer,
    budget_control_level promo.control_level,
    budget_amount numeric(20,4) DEFAULT 0,
    execution_level promo.control_level,
    effective_from date NOT NULL,
    effective_to date NOT NULL,
    is_claimable boolean DEFAULT false NOT NULL,
    claim_type promo.claim_type,
    claim_start_after_days integer,
    claim_realization_pct numeric(5,2),
    max_total_reward_type promo.reward_cap_type,
    max_total_reward_value numeric(20,4) DEFAULT 0,
    max_invoice_per_outlet numeric(10,2) DEFAULT 0,
    slab_multiplied boolean,
    strata_sequential boolean,
    minimum_sku integer DEFAULT 1 NOT NULL,
    coverage promo.coverage_type DEFAULT 'national'::promo.coverage_type NOT NULL,
    created_at timestamp(6) with time zone,
    updated_at timestamp(6) with time zone,
    created_by character varying(150),
    updated_by character varying(150),
    max_discount_outlet_uom integer DEFAULT 1,
    remarks character varying(255),
    budget_realization numeric(20,4) DEFAULT 0 NOT NULL,
    remaining_budget numeric(20,4) GENERATED ALWAYS AS (
CASE
    WHEN (is_budget_reference AND (budget_ref_type = 'limited'::promo.budget_ref_type)) THEN GREATEST((budget_amount - budget_realization), (0)::numeric)
    ELSE NULL::numeric
END) STORED,
    distributor_cust_id character varying(10),
    CONSTRAINT promotions_check CHECK ((effective_to >= effective_from)),
    CONSTRAINT promotions_check1 CHECK (((NOT is_budget_reference) OR (budget_ref_type IS NOT NULL))),
    CONSTRAINT promotions_check2 CHECK (((NOT is_budget_reference) OR (budget_reference_id IS NOT NULL))),
    CONSTRAINT promotions_check3 CHECK (((NOT (is_budget_reference AND (budget_ref_type = 'limited'::promo.budget_ref_type))) OR (budget_amount >= (0)::numeric))),
    CONSTRAINT promotions_check4 CHECK (((NOT is_claimable) OR (claim_type IS NOT NULL))),
    CONSTRAINT promotions_claim_realization_pct_check CHECK (((claim_realization_pct IS NULL) OR ((claim_realization_pct >= (0)::numeric) AND (claim_realization_pct <= (100)::numeric)))),
    CONSTRAINT promotions_minimum_sku_check CHECK ((minimum_sku >= 1))
);


--
-- Name: promotions promo_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotions
    ADD CONSTRAINT promo_pkey PRIMARY KEY (promo_id);


--
-- Name: promotion_coverage_distributors promotion_coverage_distributors_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_coverage_distributors
    ADD CONSTRAINT promotion_coverage_distributors_pkey PRIMARY KEY (id);


--
-- Name: promotion_coverage_distributors promotion_coverage_distributors_promo_id_distributor_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_coverage_distributors
    ADD CONSTRAINT promotion_coverage_distributors_promo_id_distributor_id_key UNIQUE (promo_id, distributor_id);


--
-- Name: promotion_outlet_attribute_class promotion_outlet_attribute_clas_criteria_id_outlet_class_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_class
    ADD CONSTRAINT promotion_outlet_attribute_clas_criteria_id_outlet_class_id_key UNIQUE (criteria_id, outlet_class_id);


--
-- Name: promotion_outlet_attribute_class promotion_outlet_attribute_class_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_class
    ADD CONSTRAINT promotion_outlet_attribute_class_pkey PRIMARY KEY (id);


--
-- Name: promotion_outlet_attribute_group promotion_outlet_attribute_grou_criteria_id_outlet_group_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_group
    ADD CONSTRAINT promotion_outlet_attribute_grou_criteria_id_outlet_group_id_key UNIQUE (criteria_id, outlet_group_id);


--
-- Name: promotion_outlet_attribute_group promotion_outlet_attribute_group_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_group
    ADD CONSTRAINT promotion_outlet_attribute_group_pkey PRIMARY KEY (id);


--
-- Name: promotion_outlet_attribute_sales_team promotion_outlet_attribute_sales__criteria_id_sales_team_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_sales_team
    ADD CONSTRAINT promotion_outlet_attribute_sales__criteria_id_sales_team_id_key UNIQUE (criteria_id, sales_team_id);


--
-- Name: promotion_outlet_attribute_sales_team promotion_outlet_attribute_sales_team_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_sales_team
    ADD CONSTRAINT promotion_outlet_attribute_sales_team_pkey PRIMARY KEY (id);


--
-- Name: promotion_outlet_attribute_type promotion_outlet_attribute_type_criteria_id_outlet_type_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_type
    ADD CONSTRAINT promotion_outlet_attribute_type_criteria_id_outlet_type_id_key UNIQUE (criteria_id, outlet_type_id);


--
-- Name: promotion_outlet_attribute_type promotion_outlet_attribute_type_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_type
    ADD CONSTRAINT promotion_outlet_attribute_type_pkey PRIMARY KEY (id);


--
-- Name: promotion_outlet_criteria promotion_outlet_criteria_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_criteria
    ADD CONSTRAINT promotion_outlet_criteria_pkey PRIMARY KEY (id);


--
-- Name: promotion_outlets_selected promotion_outlets_selected_criteria_id_outlet_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlets_selected
    ADD CONSTRAINT promotion_outlets_selected_criteria_id_outlet_id_key UNIQUE (criteria_id, outlet_id);


--
-- Name: promotion_outlets_selected promotion_outlets_selected_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlets_selected
    ADD CONSTRAINT promotion_outlets_selected_pkey PRIMARY KEY (id);


--
-- Name: promotion_product_criteria promotion_product_criteria_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_product_criteria
    ADD CONSTRAINT promotion_product_criteria_pkey PRIMARY KEY (id);


--
-- Name: promotion_product_criteria promotion_product_criteria_promo_id_pro_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_product_criteria
    ADD CONSTRAINT promotion_product_criteria_promo_id_pro_id_key UNIQUE (promo_id, pro_id);


--
-- Name: promotion_reward_products promotion_reward_products_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_reward_products
    ADD CONSTRAINT promotion_reward_products_pkey PRIMARY KEY (id);


--
-- Name: promotion_reward_products promotion_reward_products_promo_id_ordinal_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_reward_products
    ADD CONSTRAINT promotion_reward_products_promo_id_ordinal_key UNIQUE (promo_id, ordinal);


--
-- Name: promotion_reward_products promotion_reward_products_promo_id_pro_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_reward_products
    ADD CONSTRAINT promotion_reward_products_promo_id_pro_id_key UNIQUE (promo_id, pro_id);


--
-- Name: promotion_slabs promotion_slabs_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_slabs
    ADD CONSTRAINT promotion_slabs_pkey PRIMARY KEY (id);


--
-- Name: promotion_slabs promotion_slabs_promo_id_ordinal_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_slabs
    ADD CONSTRAINT promotion_slabs_promo_id_ordinal_key UNIQUE (promo_id, ordinal);


--
-- Name: promotion_strata promotion_strata_pkey; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_strata
    ADD CONSTRAINT promotion_strata_pkey PRIMARY KEY (id);


--
-- Name: promotion_strata promotion_strata_promo_id_ordinal_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_strata
    ADD CONSTRAINT promotion_strata_promo_id_ordinal_key UNIQUE (promo_id, ordinal);


--
-- Name: promotions promotions_cust_id_promo_id_key; Type: CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotions
    ADD CONSTRAINT promotions_cust_id_promo_id_key UNIQUE (cust_id, promo_id);


--
-- Name: idx_promo_cov_dist_promo; Type: INDEX; Schema: promo; Owner: -
--

CREATE INDEX idx_promo_cov_dist_promo ON promo.promotion_coverage_distributors USING btree (promo_id);


--
-- Name: idx_promo_outlet_sel_crit; Type: INDEX; Schema: promo; Owner: -
--

CREATE INDEX idx_promo_outlet_sel_crit ON promo.promotion_outlets_selected USING btree (criteria_id);


--
-- Name: idx_promo_reward_products_promo; Type: INDEX; Schema: promo; Owner: -
--

CREATE INDEX idx_promo_reward_products_promo ON promo.promotion_reward_products USING btree (promo_id);


--
-- Name: idx_promo_strata_promo; Type: INDEX; Schema: promo; Owner: -
--

CREATE INDEX idx_promo_strata_promo ON promo.promotion_strata USING btree (promo_id);


--
-- Name: idx_promotions_customer; Type: INDEX; Schema: promo; Owner: -
--

CREATE INDEX idx_promotions_customer ON promo.promotions USING btree (cust_id);


--
-- Name: idx_promotions_effective; Type: INDEX; Schema: promo; Owner: -
--

CREATE INDEX idx_promotions_effective ON promo.promotions USING btree (effective_from, effective_to);


--
-- Name: idx_promotions_status; Type: INDEX; Schema: promo; Owner: -
--

CREATE INDEX idx_promotions_status ON promo.promotions USING btree (promo_status);






--
-- Name: promotion_coverage_distributors promotion_coverage_distributors_promo_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_coverage_distributors
    ADD CONSTRAINT promotion_coverage_distributors_promo_id_fkey FOREIGN KEY (promo_id) REFERENCES promo.promotions(promo_id);


--
-- Name: promotion_outlet_attribute_class promotion_outlet_attribute_class_criteria_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_class
    ADD CONSTRAINT promotion_outlet_attribute_class_criteria_id_fkey FOREIGN KEY (criteria_id) REFERENCES promo.promotion_outlet_criteria(id);


--
-- Name: promotion_outlet_attribute_group promotion_outlet_attribute_group_criteria_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_group
    ADD CONSTRAINT promotion_outlet_attribute_group_criteria_id_fkey FOREIGN KEY (criteria_id) REFERENCES promo.promotion_outlet_criteria(id);


--
-- Name: promotion_outlet_attribute_sales_team promotion_outlet_attribute_sales_team_criteria_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_sales_team
    ADD CONSTRAINT promotion_outlet_attribute_sales_team_criteria_id_fkey FOREIGN KEY (criteria_id) REFERENCES promo.promotion_outlet_criteria(id);


--
-- Name: promotion_outlet_attribute_type promotion_outlet_attribute_type_criteria_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_attribute_type
    ADD CONSTRAINT promotion_outlet_attribute_type_criteria_id_fkey FOREIGN KEY (criteria_id) REFERENCES promo.promotion_outlet_criteria(id);


--
-- Name: promotion_outlet_criteria promotion_outlet_criteria_promo_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlet_criteria
    ADD CONSTRAINT promotion_outlet_criteria_promo_id_fkey FOREIGN KEY (promo_id) REFERENCES promo.promotions(promo_id);


--
-- Name: promotion_outlets_selected promotion_outlets_selected_criteria_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_outlets_selected
    ADD CONSTRAINT promotion_outlets_selected_criteria_id_fkey FOREIGN KEY (criteria_id) REFERENCES promo.promotion_outlet_criteria(id);


--
-- Name: promotion_product_criteria promotion_product_criteria_promo_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_product_criteria
    ADD CONSTRAINT promotion_product_criteria_promo_id_fkey FOREIGN KEY (promo_id) REFERENCES promo.promotions(promo_id);


--
-- Name: promotion_reward_products promotion_reward_products_promo_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_reward_products
    ADD CONSTRAINT promotion_reward_products_promo_id_fkey FOREIGN KEY (promo_id) REFERENCES promo.promotions(promo_id);


--
-- Name: promotion_slabs promotion_slabs_promo_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_slabs
    ADD CONSTRAINT promotion_slabs_promo_id_fkey FOREIGN KEY (promo_id) REFERENCES promo.promotions(promo_id);


--
-- Name: promotion_strata promotion_strata_promo_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotion_strata
    ADD CONSTRAINT promotion_strata_promo_id_fkey FOREIGN KEY (promo_id) REFERENCES promo.promotions(promo_id);


--
-- Name: promotions promotions_existing_promo_id_fkey; Type: FK CONSTRAINT; Schema: promo; Owner: -
--

ALTER TABLE ONLY promo.promotions
    ADD CONSTRAINT promotions_existing_promo_id_fkey FOREIGN KEY (existing_promo_id) REFERENCES promo.promotions(promo_id);


--
-- PostgreSQL database dump complete
--



-- Missing non-promo tables excluding procedural bodies
--
-- PostgreSQL database dump
--


-- Dumped from database version 14.5 (Ubuntu 14.5-1.pgdg20.04+1)
-- Dumped by pg_dump version 18.3 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: delivery_type; Type: TABLE; Schema: inv; Owner: -
--

CREATE TABLE IF NOT EXISTS inv.delivery_type (
    delivery_type_code character varying(25) NOT NULL,
    delivery_type_name character varying(100) NOT NULL,
    description text,
    created_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp(6) with time zone
);


--
-- Name: replenishment_order_approval; Type: TABLE; Schema: inv; Owner: -
--

CREATE TABLE IF NOT EXISTS inv.replenishment_order_approval (
    id bigint NOT NULL,
    cust_id character varying(10) NOT NULL,
    replenishment_order_id bigint NOT NULL,
    level smallint NOT NULL,
    sequence smallint NOT NULL,
    pic bigint NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    remarks character varying(225),
    CONSTRAINT ck_replenishment_order_approval_status CHECK ((status = ANY (ARRAY[1, 2, 3])))
);


--
-- Name: TABLE replenishment_order_approval; Type: COMMENT; Schema: inv; Owner: -
--

COMMENT ON TABLE inv.replenishment_order_approval IS 'Queue approval untuk transaksi replenishment order';


--
-- Name: replenishment_order_approval_id_seq; Type: SEQUENCE; Schema: inv; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS inv.replenishment_order_approval_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: replenishment_order_approval_id_seq; Type: SEQUENCE OWNED BY; Schema: inv; Owner: -
--

ALTER SEQUENCE inv.replenishment_order_approval_id_seq OWNED BY inv.replenishment_order_approval.id;


--
-- Name: distributor_replenishment_approval; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.distributor_replenishment_approval (
    cust_id character varying(10) NOT NULL,
    id integer NOT NULL,
    dist_replenishment_setup_id integer NOT NULL,
    level smallint NOT NULL,
    sequence smallint NOT NULL,
    business_unit integer NOT NULL,
    pic integer NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_by integer NOT NULL,
    created_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_by bigint,
    updated_at timestamp(6) with time zone,
    deleted_by integer,
    deleted_at timestamp(6) with time zone,
    is_del boolean DEFAULT false
);


--
-- Name: TABLE distributor_replenishment_approval; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON TABLE mst.distributor_replenishment_approval IS 'Approval chain for distributor replenishment setup (level, sequence, PIC).';


--
-- Name: COLUMN distributor_replenishment_approval.dist_replenishment_setup_id; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON COLUMN mst.distributor_replenishment_approval.dist_replenishment_setup_id IS 'FK to mst.distributor_replenishment_setup.id';


--
-- Name: COLUMN distributor_replenishment_approval.business_unit; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON COLUMN mst.distributor_replenishment_approval.business_unit IS 'Business unit id (see data dictionary; spec ref may vary).';


--
-- Name: COLUMN distributor_replenishment_approval.pic; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON COLUMN mst.distributor_replenishment_approval.pic IS 'Person in charge (approver); logical ref mst.m_employee (no DB FK — enforce in app or add PK on employee)';


--
-- Name: distributor_replenishment_approval_id_seq; Type: SEQUENCE; Schema: mst; Owner: -
--

ALTER TABLE mst.distributor_replenishment_approval ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME mst.distributor_replenishment_approval_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: distributor_replenishment_setup; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.distributor_replenishment_setup (
    cust_id character varying(10) NOT NULL,
    id integer NOT NULL,
    sup_id integer NOT NULL,
    distributor_id integer NOT NULL,
    distributor_type character varying(20) NOT NULL,
    wh_limit_action character varying(20),
    wh_capacity integer,
    wh_volume integer,
    credit_limit_action integer NOT NULL,
    plafon_credit integer,
    lead_time_days integer NOT NULL,
    is_approval_required boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_by integer NOT NULL,
    created_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_by bigint,
    updated_at timestamp(6) with time zone,
    deleted_by integer,
    deleted_at timestamp(6) with time zone,
    is_del boolean DEFAULT false,
    CONSTRAINT distributor_replenishment_setup_credit_limit_check CHECK ((credit_limit_action = ANY (ARRAY[1, 2]))),
    CONSTRAINT distributor_replenishment_setup_distributor_type_check CHECK (((distributor_type)::text = ANY (ARRAY[('FMCG'::character varying)::text, ('Fresh'::character varying)::text]))),
    CONSTRAINT distributor_replenishment_setup_wh_limit_check CHECK (((wh_limit_action IS NULL) OR ((wh_limit_action)::text = ANY ((ARRAY['Restricted'::character varying, 'Unrestricted'::character varying])::text[]))))
);


--
-- Name: TABLE distributor_replenishment_setup; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON TABLE mst.distributor_replenishment_setup IS 'Distributor replenishment / auto push allocation parameters per supplier & distributor.';


--
-- Name: COLUMN distributor_replenishment_setup.cust_id; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON COLUMN mst.distributor_replenishment_setup.cust_id IS 'Customer id; ref smc.m_customer';


--
-- Name: COLUMN distributor_replenishment_setup.sup_id; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON COLUMN mst.distributor_replenishment_setup.sup_id IS 'Supplier id; logical ref mst.m_supplier (no DB FK — enforce in app or add PK on supplier)';


--
-- Name: COLUMN distributor_replenishment_setup.distributor_id; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON COLUMN mst.distributor_replenishment_setup.distributor_id IS 'Distributor id; ref mst.m_distributor.distributor_id';


--
-- Name: COLUMN distributor_replenishment_setup.credit_limit_action; Type: COMMENT; Schema: mst; Owner: -
--

COMMENT ON COLUMN mst.distributor_replenishment_setup.credit_limit_action IS '1 = Restricted, 2 = Unrestricted';


--
-- Name: distributor_replenishment_setup_id_seq; Type: SEQUENCE; Schema: mst; Owner: -
--

ALTER TABLE mst.distributor_replenishment_setup ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME mst.distributor_replenishment_setup_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: m_outlet_backup_file_url_sht; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.m_outlet_backup_file_url_sht (
    cust_id character varying(10),
    outlet_id bigint,
    outlet_code character varying(30),
    barcode character varying(50),
    outlet_name character varying(150),
    outlet_status smallint,
    address1 character varying(150),
    address2 character varying(150),
    city character varying(100),
    zip_code character varying(6),
    phone_no character varying(20),
    wa_no character varying(20),
    fax_no character varying(20),
    email character varying(20),
    disc_grp_id bigint,
    ot_loc_id bigint,
    ot_grp_id bigint,
    price_grp_id bigint,
    district_id bigint,
    beat_id bigint,
    sbeat_id bigint,
    ot_class_id bigint,
    industry_id bigint,
    market_id bigint,
    top integer,
    payment_type smallint,
    is_contra_bon boolean,
    plu_grp_id bigint,
    conv_grp_id bigint,
    disc_inv_id bigint,
    agent_from character varying(50),
    credit_limit_type smallint,
    credit_limit numeric(20,4),
    sales_inv_limit_type smallint,
    sales_inv_limit smallint,
    avg_sales_week numeric(10,2),
    avg_sales_month numeric(10,2),
    first_trans_date date,
    last_trans_date date,
    first_week_no smallint,
    ot_start_date date,
    ot_reg_date date,
    building_own smallint,
    dob date,
    ar_status smallint,
    ar_total numeric(20,4),
    closed_date date,
    is_emb_bail boolean,
    tax_name character varying(150),
    tax_addr1 character varying(150),
    tax_addr2 character varying(150),
    tax_city character varying(100),
    tax_no character varying(30),
    owner_name character varying(150),
    owner_addr1 character varying(150),
    owner_addr2 character varying(150),
    owner_city character varying(100),
    owner_phone_no character varying(20),
    owner_id_no character varying(50),
    delv_addr1 character varying(150),
    delv_addr2 character varying(150),
    delv_city character varying(100),
    inv_addr1 character varying(150),
    inv_addr2 character varying(150),
    inv_city character varying(100),
    is_active boolean,
    created_by bigint,
    created_at timestamp(6) with time zone,
    updated_by bigint,
    updated_at timestamp(6) with time zone,
    is_del boolean,
    deleted_by bigint,
    deleted_at timestamp(6) with time zone,
    latitude character varying(50),
    longitude character varying(50),
    image_url character varying(255),
    ot_type_id bigint,
    is_obs boolean,
    obs bigint,
    outlet_ward_id character varying(10),
    is_wa_no boolean,
    delv_ward_id character varying(10),
    delv_zip_code character varying(6),
    delv_is_same_addr boolean,
    inv_ward_id character varying(10),
    inv_zip_code character varying(6),
    inv_is_same_addr boolean,
    verification_status smallint,
    verified_at timestamp(6) with time zone,
    verified_by bigint,
    tax_invoice_form smallint,
    obs_type bigint,
    credit_limit_action bigint,
    sales_inv_limit_action bigint,
    obs_limit_action bigint,
    outlet_establishment_date date,
    delv_city2 character varying(100),
    delv_latitude character varying(50),
    delv_longitude character varying(50),
    delv_latitude2 character varying(50),
    delv_longitude2 character varying(50),
    delv_ward_id2 character varying(10),
    delv_zip_code2 character varying(6),
    file_url text
);


--
-- Name: m_outlet_code_seq; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.m_outlet_code_seq (
    outlet_code_id uuid NOT NULL,
    last_sequence_no integer DEFAULT 0 NOT NULL,
    updated_by character varying(50),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: m_outlet_principal_code_seq; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.m_outlet_principal_code_seq (
    prefix character varying(64) NOT NULL,
    last_sequence_no integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: m_survey_salesman; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.m_survey_salesman (
    m_survey_salesman_id integer NOT NULL,
    cust_id character varying(10) NOT NULL,
    survey_id integer NOT NULL,
    salesman_id integer NOT NULL,
    is_del boolean DEFAULT false
);


--
-- Name: m_survey_salesman_m_survey_salesman_id_seq; Type: SEQUENCE; Schema: mst; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS mst.m_survey_salesman_m_survey_salesman_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: m_survey_salesman_m_survey_salesman_id_seq; Type: SEQUENCE OWNED BY; Schema: mst; Owner: -
--

ALTER SEQUENCE mst.m_survey_salesman_m_survey_salesman_id_seq OWNED BY mst.m_survey_salesman.m_survey_salesman_id;


--
-- Name: product_ripening_history; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.product_ripening_history (
    id bigint NOT NULL,
    cust_id character varying(30) NOT NULL,
    distributor_id bigint,
    per_year integer,
    per_id integer,
    week_id integer,
    week_start date,
    week_end date,
    source_type character varying(20) NOT NULL,
    status character varying(20) NOT NULL,
    file_url text,
    file_name character varying(255),
    total_row integer DEFAULT 0 NOT NULL,
    success_row integer DEFAULT 0 NOT NULL,
    failed_row integer DEFAULT 0 NOT NULL,
    error_summary text,
    processed_by bigint NOT NULL,
    processed_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT fk_product_ripening_history_non_negative CHECK (((total_row >= 0) AND (success_row >= 0) AND (failed_row >= 0)))
);


--
-- Name: product_ripening_history_id_seq; Type: SEQUENCE; Schema: mst; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS mst.product_ripening_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_ripening_history_id_seq; Type: SEQUENCE OWNED BY; Schema: mst; Owner: -
--

ALTER SEQUENCE mst.product_ripening_history_id_seq OWNED BY mst.product_ripening_history.id;


--
-- Name: survey_answer; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.survey_answer (
    cust_id character varying(10) NOT NULL,
    survey_answer_id bigint DEFAULT nextval('mst.survey_answer_id_seq'::regclass) NOT NULL,
    survey_template_id bigint NOT NULL,
    survey_id bigint NOT NULL,
    emp_id bigint NOT NULL,
    outlet_id bigint NOT NULL,
    area_id bigint,
    answer_date timestamp without time zone DEFAULT CURRENT_DATE,
    status character varying(20) DEFAULT 'Submitted'::character varying,
    created_by bigint NOT NULL,
    created_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_by bigint NOT NULL,
    updated_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP,
    is_del boolean DEFAULT false,
    deleted_by bigint,
    deleted_at timestamp(6) with time zone
);


--
-- Name: survey_answer_detail; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.survey_answer_detail (
    cust_id character varying(10) NOT NULL,
    survey_answer_detail_id bigint DEFAULT nextval('mst.survey_answer_detail_id_seq'::regclass) NOT NULL,
    survey_answer_id bigint NOT NULL,
    question_template_id bigint NOT NULL,
    input_type character varying(225) NOT NULL,
    answer_type character varying(20) NOT NULL,
    seq integer NOT NULL,
    is_answered boolean DEFAULT false,
    free_text_answer text,
    photo_path character varying(255),
    created_by bigint NOT NULL,
    created_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_by bigint NOT NULL,
    updated_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP,
    is_del boolean DEFAULT false,
    deleted_by bigint,
    deleted_at timestamp(6) with time zone
);


--
-- Name: survey_answer_files; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.survey_answer_files (
    cust_id character varying(10) NOT NULL,
    survey_answer_files bigint DEFAULT nextval('mst.survey_answer_file_id_seq'::regclass) NOT NULL,
    survey_answer_detail_id bigint NOT NULL,
    file_name character varying(255) NOT NULL,
    file_data bytea,
    file_key character varying(10) NOT NULL,
    media_category text NOT NULL,
    file_size bigint
);


--
-- Name: survey_answer_option; Type: TABLE; Schema: mst; Owner: -
--

CREATE TABLE IF NOT EXISTS mst.survey_answer_option (
    cust_id character varying(10) NOT NULL,
    survey_answer_option_id bigint DEFAULT nextval('mst.survey_answer_option_id_seq'::regclass) NOT NULL,
    survey_answer_detail_id bigint NOT NULL,
    q_option_template_id bigint NOT NULL,
    option_label character varying(225),
    created_by bigint NOT NULL,
    created_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_by bigint NOT NULL,
    updated_at timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP,
    is_del boolean DEFAULT false,
    deleted_by bigint,
    deleted_at timestamp(6) with time zone
);


--
-- Name: destinations; Type: TABLE; Schema: pjp_principles; Owner: -
--

CREATE TABLE IF NOT EXISTS pjp_principles.destinations (
    id integer NOT NULL,
    route_code bigint NOT NULL,
    route_name character varying(125) NOT NULL,
    status character varying(125) DEFAULT 'pending'::character varying,
    verified_date timestamp without time zone,
    destination_id bigint,
    destination_code character varying(125),
    destination_status character varying(125),
    destination_name character varying(125),
    destination_address character varying(125),
    destination_type character varying(125),
    longitude character varying(125),
    latitude character varying(125),
    pjp_id bigint,
    pjp_code bigint,
    old_pjp_id bigint,
    old_pjp_code bigint,
    old_route_code bigint,
    old_route_name character varying(125) DEFAULT NULL::character varying,
    photo character varying(125) DEFAULT NULL::character varying,
    signature character varying(125) DEFAULT NULL::character varying,
    avg_sales_week numeric(10,2) DEFAULT 0,
    cust_id character varying(125) DEFAULT NULL::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: destinations_additional; Type: TABLE; Schema: pjp_principles; Owner: -
--

CREATE TABLE IF NOT EXISTS pjp_principles.destinations_additional (
    id integer NOT NULL,
    route_code bigint NOT NULL,
    route_name character varying(125) NOT NULL,
    status character varying(125) DEFAULT 'additional'::character varying,
    verified_date timestamp without time zone,
    date timestamp without time zone,
    destination_id bigint,
    destination_code character varying(125),
    destination_status character varying(125),
    destination_name character varying(125),
    destination_address character varying(125),
    destination_type character varying(125),
    longitude character varying(125),
    latitude character varying(125),
    pjp_id bigint,
    pjp_code bigint,
    old_pjp_id bigint,
    old_pjp_code bigint,
    old_route_code bigint,
    old_route_name character varying(125) DEFAULT NULL::character varying,
    photo character varying(125) DEFAULT NULL::character varying,
    signature character varying(125) DEFAULT NULL::character varying,
    avg_sales_week numeric(10,2) DEFAULT 0,
    cust_id character varying(125) DEFAULT NULL::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_planned boolean DEFAULT false
);


--
-- Name: destinations_additional_id_seq; Type: SEQUENCE; Schema: pjp_principles; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS pjp_principles.destinations_additional_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: destinations_additional_id_seq; Type: SEQUENCE OWNED BY; Schema: pjp_principles; Owner: -
--

ALTER SEQUENCE pjp_principles.destinations_additional_id_seq OWNED BY pjp_principles.destinations_additional.id;


--
-- Name: destinations_history; Type: TABLE; Schema: pjp_principles; Owner: -
--

CREATE TABLE IF NOT EXISTS pjp_principles.destinations_history (
    id integer NOT NULL,
    route_code bigint NOT NULL,
    route_name character varying(125) NOT NULL,
    verified_date timestamp without time zone,
    date timestamp without time zone,
    week integer,
    year integer,
    index_day integer,
    start_week timestamp without time zone,
    is_in_current_year boolean,
    is_additional boolean DEFAULT false,
    destination_id bigint,
    destination_code character varying(125),
    destination_status character varying(125),
    destination_name character varying(125),
    destination_address character varying(125),
    destination_type character varying(125),
    longitude character varying(125),
    latitude character varying(125),
    pjp_id bigint,
    pjp_code bigint,
    old_pjp_id bigint,
    old_pjp_code bigint,
    old_route_code bigint,
    old_route_name character varying(125) DEFAULT NULL::character varying,
    photo character varying(125) DEFAULT NULL::character varying,
    signature character varying(125) DEFAULT NULL::character varying,
    avg_sales_week numeric(10,2) DEFAULT 0,
    cust_id character varying(125) DEFAULT NULL::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_extra_call boolean DEFAULT false
);


--
-- Name: destinations_history_id_seq; Type: SEQUENCE; Schema: pjp_principles; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS pjp_principles.destinations_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: destinations_history_id_seq; Type: SEQUENCE OWNED BY; Schema: pjp_principles; Owner: -
--

ALTER SEQUENCE pjp_principles.destinations_history_id_seq OWNED BY pjp_principles.destinations_history.id;


--
-- Name: destinations_id_seq; Type: SEQUENCE; Schema: pjp_principles; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS pjp_principles.destinations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: destinations_id_seq; Type: SEQUENCE OWNED BY; Schema: pjp_principles; Owner: -
--

ALTER SEQUENCE pjp_principles.destinations_id_seq OWNED BY pjp_principles.destinations.id;


--
-- Name: route_pop_dailies; Type: TABLE; Schema: pjp_principles; Owner: -
--

CREATE TABLE IF NOT EXISTS pjp_principles.route_pop_dailies (
    id integer NOT NULL,
    year bigint,
    week bigint,
    date timestamp without time zone,
    day character varying(125),
    route_code bigint,
    pjp_id bigint,
    pjp_code bigint,
    parent_route bigint,
    status character varying(125) DEFAULT 'active'::character varying,
    cust_id character varying(125) DEFAULT NULL::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: route_pop_dailies_id_seq; Type: SEQUENCE; Schema: pjp_principles; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS pjp_principles.route_pop_dailies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: route_pop_dailies_id_seq; Type: SEQUENCE OWNED BY; Schema: pjp_principles; Owner: -
--

ALTER SEQUENCE pjp_principles.route_pop_dailies_id_seq OWNED BY pjp_principles.route_pop_dailies.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: pjp_principles; Owner: -
--

CREATE TABLE IF NOT EXISTS pjp_principles.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: invoice_no_counter; Type: TABLE; Schema: sls; Owner: -
--

CREATE TABLE IF NOT EXISTS sls.invoice_no_counter (
    cust_id character varying(20) NOT NULL,
    seq_date date NOT NULL,
    last_seq integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: password_reset_requests; Type: TABLE; Schema: sys; Owner: -
--

CREATE TABLE IF NOT EXISTS sys.password_reset_requests (
    id integer NOT NULL,
    user_id bigint NOT NULL,
    email character varying(255) NOT NULL,
    otp_code character varying(255) NOT NULL,
    otp_expired_at timestamp without time zone NOT NULL,
    otp_attempt_count integer DEFAULT 0 NOT NULL,
    otp_max_attempt integer DEFAULT 3 NOT NULL,
    resend_count integer DEFAULT 0 NOT NULL,
    resend_max integer DEFAULT 3 NOT NULL,
    resend_cooldown_until timestamp without time zone,
    request_id character varying(100) NOT NULL,
    reset_token character varying(255),
    reset_token_expired_at timestamp without time zone,
    status character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: password_reset_requests_id_seq; Type: SEQUENCE; Schema: sys; Owner: -
--

CREATE SEQUENCE IF NOT EXISTS sys.password_reset_requests_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: password_reset_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: sys; Owner: -
--

ALTER SEQUENCE sys.password_reset_requests_id_seq OWNED BY sys.password_reset_requests.id;


--
-- Name: replenishment_order_approval id; Type: DEFAULT; Schema: inv; Owner: -
--

ALTER TABLE ONLY inv.replenishment_order_approval ALTER COLUMN id SET DEFAULT nextval('inv.replenishment_order_approval_id_seq'::regclass);


--
-- Name: m_survey_salesman m_survey_salesman_id; Type: DEFAULT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.m_survey_salesman ALTER COLUMN m_survey_salesman_id SET DEFAULT nextval('mst.m_survey_salesman_m_survey_salesman_id_seq'::regclass);


--
-- Name: product_ripening_history id; Type: DEFAULT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.product_ripening_history ALTER COLUMN id SET DEFAULT nextval('mst.product_ripening_history_id_seq'::regclass);


--
-- Name: destinations id; Type: DEFAULT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations ALTER COLUMN id SET DEFAULT nextval('pjp_principles.destinations_id_seq'::regclass);


--
-- Name: destinations_additional id; Type: DEFAULT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations_additional ALTER COLUMN id SET DEFAULT nextval('pjp_principles.destinations_additional_id_seq'::regclass);


--
-- Name: destinations_history id; Type: DEFAULT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations_history ALTER COLUMN id SET DEFAULT nextval('pjp_principles.destinations_history_id_seq'::regclass);


--
-- Name: route_pop_dailies id; Type: DEFAULT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.route_pop_dailies ALTER COLUMN id SET DEFAULT nextval('pjp_principles.route_pop_dailies_id_seq'::regclass);


--
-- Name: password_reset_requests id; Type: DEFAULT; Schema: sys; Owner: -
--

ALTER TABLE ONLY sys.password_reset_requests ALTER COLUMN id SET DEFAULT nextval('sys.password_reset_requests_id_seq'::regclass);


--
-- Name: delivery_type delivery_type_pkey; Type: CONSTRAINT; Schema: inv; Owner: -
--

ALTER TABLE ONLY inv.delivery_type
    ADD CONSTRAINT delivery_type_pkey PRIMARY KEY (delivery_type_code);


--
-- Name: replenishment_order_approval replenishment_order_approval_pkey; Type: CONSTRAINT; Schema: inv; Owner: -
--

ALTER TABLE ONLY inv.replenishment_order_approval
    ADD CONSTRAINT replenishment_order_approval_pkey PRIMARY KEY (id);


--
-- Name: distributor_replenishment_approval distributor_replenishment_approval_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.distributor_replenishment_approval
    ADD CONSTRAINT distributor_replenishment_approval_pkey PRIMARY KEY (id);


--
-- Name: distributor_replenishment_setup distributor_replenishment_setup_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.distributor_replenishment_setup
    ADD CONSTRAINT distributor_replenishment_setup_pkey PRIMARY KEY (id);


--
-- Name: m_outlet_code_seq m_outlet_code_seq_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.m_outlet_code_seq
    ADD CONSTRAINT m_outlet_code_seq_pkey PRIMARY KEY (outlet_code_id);


--
-- Name: m_outlet_principal_code_seq m_outlet_principal_code_seq_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.m_outlet_principal_code_seq
    ADD CONSTRAINT m_outlet_principal_code_seq_pkey PRIMARY KEY (prefix);


--
-- Name: m_survey_salesman m_survey_salesman_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.m_survey_salesman
    ADD CONSTRAINT m_survey_salesman_pkey PRIMARY KEY (m_survey_salesman_id);


--
-- Name: product_ripening_history product_ripening_history_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.product_ripening_history
    ADD CONSTRAINT product_ripening_history_pkey PRIMARY KEY (id);


--
-- Name: survey_answer_detail survey_answer_detail_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_detail
    ADD CONSTRAINT survey_answer_detail_pkey PRIMARY KEY (survey_answer_detail_id);


--
-- Name: survey_answer_files survey_answer_files_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_files
    ADD CONSTRAINT survey_answer_files_pkey PRIMARY KEY (survey_answer_files);


--
-- Name: survey_answer_option survey_answer_option_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_option
    ADD CONSTRAINT survey_answer_option_pkey PRIMARY KEY (survey_answer_option_id);


--
-- Name: survey_answer survey_answer_pkey; Type: CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer
    ADD CONSTRAINT survey_answer_pkey PRIMARY KEY (survey_answer_id);


--
-- Name: destinations_additional destinations_additional_pkey; Type: CONSTRAINT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations_additional
    ADD CONSTRAINT destinations_additional_pkey PRIMARY KEY (id);


--
-- Name: destinations_history destinations_history_pkey; Type: CONSTRAINT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations_history
    ADD CONSTRAINT destinations_history_pkey PRIMARY KEY (id);


--
-- Name: destinations destinations_pkey; Type: CONSTRAINT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations
    ADD CONSTRAINT destinations_pkey PRIMARY KEY (id);


--
-- Name: invoice_no_counter invoice_no_counter_pkey; Type: CONSTRAINT; Schema: sls; Owner: -
--

ALTER TABLE ONLY sls.invoice_no_counter
    ADD CONSTRAINT invoice_no_counter_pkey PRIMARY KEY (cust_id, seq_date);


--
-- Name: password_reset_requests password_reset_requests_pkey; Type: CONSTRAINT; Schema: sys; Owner: -
--

ALTER TABLE ONLY sys.password_reset_requests
    ADD CONSTRAINT password_reset_requests_pkey PRIMARY KEY (id);


--
-- Name: password_reset_requests password_reset_requests_request_id_key; Type: CONSTRAINT; Schema: sys; Owner: -
--

ALTER TABLE ONLY sys.password_reset_requests
    ADD CONSTRAINT password_reset_requests_request_id_key UNIQUE (request_id);


--
-- Name: idx_delivery_type_name; Type: INDEX; Schema: inv; Owner: -
--

CREATE INDEX idx_delivery_type_name ON inv.delivery_type USING btree (delivery_type_name);


--
-- Name: idx_roa_created_at; Type: INDEX; Schema: inv; Owner: -
--

CREATE INDEX idx_roa_created_at ON inv.replenishment_order_approval USING btree (created_at);


--
-- Name: idx_roa_cust_repl; Type: INDEX; Schema: inv; Owner: -
--

CREATE INDEX idx_roa_cust_repl ON inv.replenishment_order_approval USING btree (cust_id, replenishment_order_id);


--
-- Name: idx_roa_pic_status; Type: INDEX; Schema: inv; Owner: -
--

CREATE INDEX idx_roa_pic_status ON inv.replenishment_order_approval USING btree (pic, status);


--
-- Name: idx_distributor_replenishment_approval_cust_id; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_approval_cust_id ON mst.distributor_replenishment_approval USING btree (cust_id);


--
-- Name: idx_distributor_replenishment_approval_is_del; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_approval_is_del ON mst.distributor_replenishment_approval USING btree (is_del);


--
-- Name: idx_distributor_replenishment_approval_pic; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_approval_pic ON mst.distributor_replenishment_approval USING btree (pic);


--
-- Name: idx_distributor_replenishment_approval_setup_id; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_approval_setup_id ON mst.distributor_replenishment_approval USING btree (dist_replenishment_setup_id);


--
-- Name: idx_distributor_replenishment_setup_cust_id; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_setup_cust_id ON mst.distributor_replenishment_setup USING btree (cust_id);


--
-- Name: idx_distributor_replenishment_setup_distributor_id; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_setup_distributor_id ON mst.distributor_replenishment_setup USING btree (distributor_id);


--
-- Name: idx_distributor_replenishment_setup_is_del; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_setup_is_del ON mst.distributor_replenishment_setup USING btree (is_del);


--
-- Name: idx_distributor_replenishment_setup_sup_id; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_distributor_replenishment_setup_sup_id ON mst.distributor_replenishment_setup USING btree (sup_id);


--
-- Name: idx_m_survey_salesman_cust_salesman; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_m_survey_salesman_cust_salesman ON mst.m_survey_salesman USING btree (cust_id, salesman_id);


--
-- Name: idx_m_survey_salesman_cust_survey; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_m_survey_salesman_cust_survey ON mst.m_survey_salesman USING btree (cust_id, survey_id);


--
-- Name: idx_m_survey_salesman_survey_id; Type: INDEX; Schema: mst; Owner: -
--

CREATE INDEX idx_m_survey_salesman_survey_id ON mst.m_survey_salesman USING btree (survey_id);


--
-- Name: replenishment_order_approval fk_roa_customer; Type: FK CONSTRAINT; Schema: inv; Owner: -
--

ALTER TABLE ONLY inv.replenishment_order_approval
    ADD CONSTRAINT fk_roa_customer FOREIGN KEY (cust_id) REFERENCES smc.m_customer(cust_id);


--
-- Name: distributor_replenishment_approval fk_distributor_replenishment_approval_cust; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.distributor_replenishment_approval
    ADD CONSTRAINT fk_distributor_replenishment_approval_cust FOREIGN KEY (cust_id) REFERENCES smc.m_customer(cust_id);


--
-- Name: distributor_replenishment_approval fk_distributor_replenishment_approval_setup; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.distributor_replenishment_approval
    ADD CONSTRAINT fk_distributor_replenishment_approval_setup FOREIGN KEY (dist_replenishment_setup_id) REFERENCES mst.distributor_replenishment_setup(id);


--
-- Name: distributor_replenishment_setup fk_distributor_replenishment_setup_cust; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.distributor_replenishment_setup
    ADD CONSTRAINT fk_distributor_replenishment_setup_cust FOREIGN KEY (cust_id) REFERENCES smc.m_customer(cust_id);


--
-- Name: distributor_replenishment_setup fk_distributor_replenishment_setup_distributor; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.distributor_replenishment_setup
    ADD CONSTRAINT fk_distributor_replenishment_setup_distributor FOREIGN KEY (distributor_id) REFERENCES mst.m_distributor(distributor_id);


--
-- Name: product_ripening_history fk_product_ripening_history_customer; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.product_ripening_history
    ADD CONSTRAINT fk_product_ripening_history_customer FOREIGN KEY (cust_id) REFERENCES smc.m_customer(cust_id);


--
-- Name: product_ripening_history fk_product_ripening_history_distributor; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.product_ripening_history
    ADD CONSTRAINT fk_product_ripening_history_distributor FOREIGN KEY (distributor_id) REFERENCES mst.m_distributor(distributor_id);


--
-- Name: survey_answer_option fk_q_option_template; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_option
    ADD CONSTRAINT fk_q_option_template FOREIGN KEY (q_option_template_id) REFERENCES mst.m_q_option_template(q_option_template_id);


--
-- Name: survey_answer_detail fk_question_template; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_detail
    ADD CONSTRAINT fk_question_template FOREIGN KEY (question_template_id) REFERENCES mst.question_template(question_template_id);


--
-- Name: survey_answer fk_survey; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer
    ADD CONSTRAINT fk_survey FOREIGN KEY (survey_id) REFERENCES mst.m_survey(survey_id);


--
-- Name: survey_answer_detail fk_survey_answer; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_detail
    ADD CONSTRAINT fk_survey_answer FOREIGN KEY (survey_answer_id) REFERENCES mst.survey_answer(survey_answer_id);


--
-- Name: survey_answer_files fk_survey_detail_file; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_files
    ADD CONSTRAINT fk_survey_detail_file FOREIGN KEY (survey_answer_detail_id) REFERENCES mst.survey_answer_detail(survey_answer_detail_id);


--
-- Name: survey_answer_option fk_survey_detail_opt; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer_option
    ADD CONSTRAINT fk_survey_detail_opt FOREIGN KEY (survey_answer_detail_id) REFERENCES mst.survey_answer_detail(survey_answer_detail_id);


--
-- Name: survey_answer fk_survey_template; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.survey_answer
    ADD CONSTRAINT fk_survey_template FOREIGN KEY (survey_template_id) REFERENCES mst.m_survey_template(survey_template_id);


--
-- Name: m_survey_salesman m_survey_salesman_survey_id_fkey; Type: FK CONSTRAINT; Schema: mst; Owner: -
--

ALTER TABLE ONLY mst.m_survey_salesman
    ADD CONSTRAINT m_survey_salesman_survey_id_fkey FOREIGN KEY (survey_id) REFERENCES mst.m_survey(survey_id);


--
-- Name: destinations_additional fk_destinations_additional_pjp_principles; Type: FK CONSTRAINT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations_additional
    ADD CONSTRAINT fk_destinations_additional_pjp_principles FOREIGN KEY (pjp_id) REFERENCES pjp_principles.permanent_journey_plans(id);


--
-- Name: destinations_additional fk_destinations_additional_routes; Type: FK CONSTRAINT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations_additional
    ADD CONSTRAINT fk_destinations_additional_routes FOREIGN KEY (route_code) REFERENCES pjp_principles.routes(route_code);


--
-- Name: destinations fk_destinations_pjp_principles; Type: FK CONSTRAINT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations
    ADD CONSTRAINT fk_destinations_pjp_principles FOREIGN KEY (pjp_id) REFERENCES pjp_principles.permanent_journey_plans(id);


--
-- Name: destinations fk_destinations_routes; Type: FK CONSTRAINT; Schema: pjp_principles; Owner: -
--

ALTER TABLE ONLY pjp_principles.destinations
    ADD CONSTRAINT fk_destinations_routes FOREIGN KEY (route_code) REFERENCES pjp_principles.routes(route_code);


--
-- PostgreSQL database dump complete
--



-- Missing columns on existing demo tables
-- Missing columns on demo existing tables (generated from staging catalog)
-- Review NOT NULL additions carefully; defaults from staging are included when present.

ALTER TABLE "acf"."bank_transfer" ADD COLUMN IF NOT EXISTS "account_name" character varying(255) DEFAULT ''::character varying;
ALTER TABLE "acf"."collection_det" ADD COLUMN IF NOT EXISTS "salesman_id" bigint;
ALTER TABLE "acf"."deposit_payment" ADD COLUMN IF NOT EXISTS "deleted_at" timestamp with time zone;
ALTER TABLE "acf"."expense" ADD COLUMN IF NOT EXISTS "doc_no" character varying(50);
ALTER TABLE "acf"."expense" ADD COLUMN IF NOT EXISTS "source" integer;
ALTER TABLE "acf"."expense" ADD COLUMN IF NOT EXISTS "balance" numeric(20,4) DEFAULT 0;
ALTER TABLE "acf"."expense_det" ADD COLUMN IF NOT EXISTS "collector_id" bigint;
ALTER TABLE "acf"."expense_det" ADD COLUMN IF NOT EXISTS "expense_type_id" integer;
ALTER TABLE "acf"."expense_det" ADD COLUMN IF NOT EXISTS "amount" numeric(20,4) DEFAULT 0;
ALTER TABLE "acf"."expense_det" ADD COLUMN IF NOT EXISTS "notes" character varying(100);
ALTER TABLE "acf"."expense_type" ADD COLUMN IF NOT EXISTS "cust_id" character varying(10);
ALTER TABLE "acf"."expense_type" ADD COLUMN IF NOT EXISTS "note" character varying(100);
ALTER TABLE "inv"."replenishment_order" ADD COLUMN IF NOT EXISTS "distributor_id" bigint;
ALTER TABLE "inv"."replenishment_order_detail" ADD COLUMN IF NOT EXISTS "return_reason_id" bigint;
ALTER TABLE "mst"."m_dist_price" ADD COLUMN IF NOT EXISTS "status" smallint DEFAULT 1;
ALTER TABLE "mst"."m_distributor" ADD COLUMN IF NOT EXISTS "allow_add_product" boolean DEFAULT false;
ALTER TABLE "mst"."m_distributor" ADD COLUMN IF NOT EXISTS "allow_edit_product" boolean DEFAULT false;
ALTER TABLE "mst"."m_distributor" ADD COLUMN IF NOT EXISTS "allow_manage_pricing" boolean DEFAULT false;
ALTER TABLE "mst"."m_distributor" ADD COLUMN IF NOT EXISTS "allow_upload_secondary_sales" boolean DEFAULT false;
ALTER TABLE "mst"."m_distributor" ADD COLUMN IF NOT EXISTS "parent_cust_id" character varying(10);
ALTER TABLE "mst"."m_outlet" ADD COLUMN IF NOT EXISTS "credit_limit_type_name" character varying(100);
ALTER TABLE "mst"."m_outlet" ADD COLUMN IF NOT EXISTS "credit_limit_action_name" character varying(100);
ALTER TABLE "mst"."m_outlet" ADD COLUMN IF NOT EXISTS "sales_inv_limit_type_name" character varying(100);
ALTER TABLE "mst"."m_outlet" ADD COLUMN IF NOT EXISTS "sales_inv_limit_action_name" character varying(100);
ALTER TABLE "mst"."m_outlet" ADD COLUMN IF NOT EXISTS "source" integer;
ALTER TABLE "mst"."m_price" ADD COLUMN IF NOT EXISTS "created_by_id" bigint;
ALTER TABLE "mst"."m_price" ADD COLUMN IF NOT EXISTS "updated_by_id" bigint;
ALTER TABLE "mst"."m_product" ADD COLUMN IF NOT EXISTS "distributor_id" bigint;
ALTER TABLE "mst"."m_product" ADD COLUMN IF NOT EXISTS "level" integer DEFAULT 0;
ALTER TABLE "mst"."m_product" ADD COLUMN IF NOT EXISTS "description" character varying(225);
ALTER TABLE "mst"."m_product" ADD COLUMN IF NOT EXISTS "referral_code" character varying(20);
ALTER TABLE "mst"."m_survey_area" ADD COLUMN IF NOT EXISTS "distributor_id" integer;
ALTER TABLE "mst"."question_template" ADD COLUMN IF NOT EXISTS "use_image" boolean DEFAULT false;
ALTER TABLE "mst"."question_template" ADD COLUMN IF NOT EXISTS "input_type" character varying(20);
ALTER TABLE "pjp"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "is_extra_call" boolean DEFAULT false;
ALTER TABLE "pjp"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "distance_meter" integer;
ALTER TABLE "pjp"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "allowed_radius" integer DEFAULT 100;
ALTER TABLE "pjp"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "location_status" smallint;
ALTER TABLE "pjp"."route_outlet_history" ADD COLUMN IF NOT EXISTS "is_extra_call" boolean DEFAULT false;
ALTER TABLE "pjp"."routes" ADD COLUMN IF NOT EXISTS "sequence" integer DEFAULT 0;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "is_extra_call" boolean DEFAULT false;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "photo_path" character varying(500);
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "latitude" character varying(50);
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "longitude" character varying(50);
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "is_update_location" boolean DEFAULT false;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "folder" character varying(255);
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "file_base64" text;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "file_name" character varying(255);
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "file_type" character varying(50);
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "media_category" character varying(20);
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "file_url" text;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "file_size" integer;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "distance_meter" integer;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "allowed_radius" integer DEFAULT 100;
ALTER TABLE "pjp_principles"."outlet_visit_list" ADD COLUMN IF NOT EXISTS "location_status" integer;
ALTER TABLE "pjp_principles"."routes" ADD COLUMN IF NOT EXISTS "pjp_id" integer;
ALTER TABLE "pjp_principles"."routes" ADD COLUMN IF NOT EXISTS "sequence" integer DEFAULT 0;
ALTER TABLE "report"."list" ADD COLUMN IF NOT EXISTS "file_base64" text;
ALTER TABLE "report"."list" ADD COLUMN IF NOT EXISTS "updated_at" timestamp with time zone;
ALTER TABLE "sls"."order" ADD COLUMN IF NOT EXISTS "is_proforma_inv" boolean;
ALTER TABLE "sls"."order" ADD COLUMN IF NOT EXISTS "generate_by" bigint;
ALTER TABLE "sls"."order" ADD COLUMN IF NOT EXISTS "first_issue_date" timestamp with time zone;
ALTER TABLE "sls"."order" ADD COLUMN IF NOT EXISTS "promo_remarks_so" jsonb DEFAULT '[]'::jsonb;
ALTER TABLE "sls"."order" ADD COLUMN IF NOT EXISTS "promo_remarks_final" jsonb DEFAULT '[]'::jsonb;
ALTER TABLE "sls"."order" ADD COLUMN IF NOT EXISTS "promo_remarks_po" jsonb DEFAULT '[]'::jsonb;
ALTER TABLE "sls"."order" ADD COLUMN IF NOT EXISTS "opr_type" character(1);
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_so1" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_so2" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_so3" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_so4" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_so5" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_final1" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_final2" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_final3" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_final4" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_final5" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_po1" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_po2" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_po3" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_po4" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_po5" numeric(20,4) DEFAULT 0;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_remarks_so" jsonb DEFAULT '[]'::jsonb;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_remarks_final" jsonb DEFAULT '[]'::jsonb;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "promo_remarks_po" jsonb DEFAULT '[]'::jsonb;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "is_product_promotion_so" boolean DEFAULT false;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "is_product_promotion_final" boolean DEFAULT false;
ALTER TABLE "sls"."order_detail" ADD COLUMN IF NOT EXISTS "is_product_promotion_po" boolean DEFAULT false;

COMMIT;
