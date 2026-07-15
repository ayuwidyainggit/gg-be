-- Schema
CREATE SCHEMA IF NOT EXISTS promo;

-- ================
-- 1) ENUMS
-- ================
DO $$
BEGIN
  -- top-level
  CREATE TYPE promo.promotion_status AS ENUM ('draft','submit','approved','rejected','inactive','active','closed');
  CREATE TYPE promo.promotion_type   AS ENUM ('slab','strata');
  CREATE TYPE promo.creation_type    AS ENUM ('new','replacement');
  CREATE TYPE promo.claim_type       AS ENUM ('full','partial');
  CREATE TYPE promo.budget_ref_type  AS ENUM ('unlimited','limited');
  CREATE TYPE promo.control_level    AS ENUM ('region','area','distributor','salesman');
  CREATE TYPE promo.coverage_type    AS ENUM ('national','by_distributor');
  CREATE TYPE promo.outlet_sel_type  AS ENUM ('by_outlet','by_attribute');
  CREATE TYPE promo.rule_type        AS ENUM ('quantity','value');
  CREATE TYPE promo.uom_type         AS ENUM ('smallest','middle','largest');
  CREATE TYPE promo.reward_type      AS ENUM ('percentage','fixed_value','product');
  CREATE TYPE promo.reward_cap_type  AS ENUM ('amount','qty');
EXCEPTION WHEN duplicate_object THEN NULL;
END$$;

DROP TABLE IF EXISTS promo.promotions;
CREATE TABLE IF NOT EXISTS promo.promotions (
  "cust_id" varchar(10) COLLATE "pg_catalog"."default" NOT NULL,
  "promo_id" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "promo_desc" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "promo_type" promo.promotion_type NOT NULL, -- slab / strata
  "promo_creation_type" promo.creation_type NOT NULL, -- new / replacement
  "existing_promo_id" varchar(50) NULL REFERENCES promo.promotions("promo_id") ON UPDATE CASCADE ON DELETE RESTRICT,
  "promo_status" promo.promotion_status NOT NULL DEFAULT 'draft',
  "is_budget_reference" bool NOT NULL DEFAULT false,
  "budget_ref_type" promo.budget_ref_type DEFAULT NULL, -- unlimited / limited (enabled when is_budget_reference = true) 
  "budget_reference_id" int4,
  "budget_control_level" promo.control_level DEFAULT NULL,
  "budget_amount" numeric(20,4) DEFAULT 0, -- only for limited
  "execution_level" promo.control_level DEFAULT NULL,
  "effective_from" date NOT NULL,
  "effective_to" date NOT NULL,
   CHECK (effective_to >= effective_from),
  "is_claimable" bool NOT NULL DEFAULT false, -- Yes/No
  "claim_type" promo.claim_type DEFAULT NULL, -- full/partial (when claimable)
  "claim_start_after_days" int4, -- report claim start after X days from end date
  "claim_realization_pct" NUMERIC(5,2), -- 0..100 (when claim_type = partial)
  "max_total_reward_type" promo.reward_cap_type DEFAULT NULL, -- amount/qty (optional)
  "max_total_reward_value" numeric(20,4) DEFAULT 0, -- meaning depends on type
  "max_invoice_per_outlet" numeric(10,2) DEFAULT 0, -- optional

   -- Multipliers / global flags (applies to SLAB only)
  "slab_multiplied" BOOLEAN, -- if TRUE: single slab replicated; % reward disabled (business rule)

  -- Strata global (applies to STRATA only)
  "strata_sequential" BOOLEAN, -- sequential calc flag (Yes/No)

  -- Product criteria
  "minimum_sku" int4 NOT NULL DEFAULT 1 CHECK (minimum_sku >= 1),

  -- Coverage
  "coverage" promo.coverage_type NOT NULL DEFAULT 'national',
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "created_by" varchar(150) COLLATE "pg_catalog"."default",
  "updated_by" varchar(150) COLLATE "pg_catalog"."default",
  "max_discount_outlet_uom" int4 DEFAULT 1,
  "remarks" varchar(255) COLLATE "pg_catalog"."default",

  -- Derived / runtime fields (optional, useful for view/detail)
  "budget_realization" NUMERIC(20,4) NOT NULL DEFAULT 0,
  "remaining_budget" NUMERIC(20,4) GENERATED ALWAYS AS (CASE
                                WHEN is_budget_reference AND budget_ref_type = 'limited'
                                  THEN GREATEST(budget_amount - budget_realization, 0)
                                ELSE NULL
                               END) STORED,

  UNIQUE ("cust_id", "promo_id"),

  -- Validations
  CHECK (NOT is_budget_reference OR budget_ref_type IS NOT NULL),
  CHECK (NOT is_budget_reference OR budget_reference_id IS NOT NULL),
  CHECK (NOT (is_budget_reference AND budget_ref_type = 'limited') OR budget_amount >= 0),
  CHECK (NOT is_claimable OR claim_type IS NOT NULL),
  CHECK (claim_realization_pct IS NULL OR (claim_realization_pct >= 0 AND claim_realization_pct <= 100)),
  "distributor_cust_id" varchar(10) COLLATE "pg_catalog"."default",

	CONSTRAINT "promo_pkey" PRIMARY KEY ("promo_id")
);

ALTER TABLE "promo"."promotions" 
  OWNER TO "postgres";

CREATE INDEX IF NOT EXISTS idx_promotions_customer ON promo.promotions(cust_id);
CREATE INDEX IF NOT EXISTS idx_promotions_status  ON promo.promotions(promo_status);
CREATE INDEX IF NOT EXISTS idx_promotions_effective
  ON promo.promotions(effective_from, effective_to);



-- ================
-- Prerequisites for MongoDB ObjectId Generation
-- ================
-- Enable pgcrypto extension for cryptographic functions (gen_random_bytes)
-- This extension provides secure random number generation capabilities
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create sequence for ObjectId counter component
-- This sequence provides incremental values for the 3-byte counter part
-- of MongoDB ObjectId to ensure uniqueness even at high generation rates
CREATE SEQUENCE IF NOT EXISTS mongodb_objectid_counter;

-- ================
-- MongoDB ObjectId Generator
-- ================
-- Generates a 24-character hexadecimal string that follows MongoDB ObjectId format
-- Used for creating unique identifiers compatible with MongoDB systems
-- Format: 4-byte timestamp + 5-byte random + 3-byte counter = 12 bytes (24 hex chars)
CREATE OR REPLACE FUNCTION public.generate_mongodb_id()
RETURNS text
LANGUAGE plpgsql
AS $$
DECLARE
    ts_int   integer;  -- Unix timestamp in seconds (4 bytes)
    rnd5     bytea;    -- 5 random bytes for uniqueness
    cnt      bigint;   -- Sequence counter from sequence
    cnt24    integer;  -- Counter wrapped to 24-bit value (0-16777215)
    cnt3     bytea := E'\\000\\000\\000'::bytea;  -- 3-byte buffer for counter
BEGIN
    -- 1) 4-byte timestamp (big-endian, network byte order)
    -- Represents seconds since Unix epoch, same as MongoDB ObjectId
    ts_int := EXTRACT(EPOCH FROM clock_timestamp())::integer;
    -- int4send(integer) returns 4-byte big-endian bytea
    -- (Exactly what MongoDB uses for the timestamp part)
    
    -- 2) 5 random bytes for uniqueness
    -- Ensures ObjectIds are unique even when generated at the same timestamp
    rnd5 := gen_random_bytes(5);

    -- 3) 3-byte counter (sequence, wrapped to 24 bits)
    -- Provides additional uniqueness for high-frequency ID generation
    cnt   := nextval('mongodb_objectid_counter');
    cnt24 := (cnt % 16777216)::integer;  -- 2^24 = 16,777,216 (max 24-bit value)

    -- Pack 3-byte counter: [high byte, mid byte, low byte]
    -- Extract each byte using bit shifting and masking
    cnt3 := set_byte(cnt3, 0, ((cnt24 >> 16) & 255)::int);  -- High byte (bits 23-16)
    cnt3 := set_byte(cnt3, 1, ((cnt24 >> 8)  & 255)::int);  -- Mid byte  (bits 15-8)
    cnt3 := set_byte(cnt3, 2, ( cnt24        & 255)::int);  -- Low byte  (bits 7-0)

    -- Concatenate all parts and encode as hexadecimal string
    -- Result: 24-character hex string (12 bytes total)
    RETURN encode(int4send(ts_int) || rnd5 || cnt3, 'hex');
END;
$$;

-- ================
-- SLAB (rules & rewards)
-- ================
-- For SLAB: arbitrary number of slabs; if slab_multiplied = true, app logic will replicate ranges/rewards.
DROP TABLE IF EXISTS promo.promotion_slabs CASCADE;
CREATE TABLE IF NOT EXISTS promo.promotion_slabs (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  promo_id       varchar(50) NOT NULL REFERENCES promo.promotions(promo_id) ON UPDATE CASCADE ON DELETE CASCADE,
  ordinal        INTEGER NOT NULL,                                   -- 1..N
  description    VARCHAR(50),
  rule_type      promo.rule_type NOT NULL,                           -- quantity/value
  rule_uom       promo.uom_type DEFAULT NULL,                        -- smallest/middle/largest
  range_from     NUMERIC(20,4),                                      -- may be NULL when multiplied
  range_to       NUMERIC(20,4) NOT NULL,
  CHECK (range_to > COALESCE(range_from, 0)),

  reward_type    promo.reward_type NOT NULL,                         -- percentage/fixed_value/product
  reward_value   NUMERIC(20,4),                                      -- % or fixed amount; NULL if product
  reward_uom     promo.uom_type DEFAULT NULL,                        -- smallest/middle/largest
  per_scope      VARCHAR(16),                                        -- 'per_product'|'per_order' if fixed_value chosen

  UNIQUE(promo_id, ordinal),

  -- Business constraints we can enforce structurally:
  CHECK (NOT (reward_type = 'percentage') OR (reward_value IS NOT NULL AND reward_value BETWEEN 1 AND 100)),
  CHECK (NOT (reward_type = 'fixed_value') OR (reward_value IS NOT NULL AND reward_value > 0)),
  CHECK (NOT (reward_type = 'product')    OR (reward_value IS NOT NULL AND reward_value > 0)),

  -- If multiplied at header, force range_from NULL and forbid percentage rewards (validated at insert/update via trigger if needed)
  -- (Use app/trigger to hard-enforce: slab_multiplied=true => range_from IS NULL AND reward_type!='percentage')
  -- Keep as comment because header flag lives in another table.
  CHECK (per_scope IS NULL OR per_scope IN ('per_product','per_order'))
);

CREATE INDEX IF NOT EXISTS idx_promo_slabs_promo ON promo.promotion_slabs(promo_id);


-- ================
-- Trigger validasi aturan SLAB
-- ================
-- Jika slab_multiplied = TRUE di header:
-- range_from harus NULL
-- reward_type != 'percentage'
-- Konsistensi dalam satu promo: semua slab harus punya rule_type yang sama dan reward_type yang sama
-- Reward meningkat seiring ordinal (lebih besar dari sebelumnya)
CREATE OR REPLACE FUNCTION promo.fn_check_promotion_slabs()
RETURNS TRIGGER AS $$
DECLARE
  v_promo_type promo.promotion_type;
	v_prev_range_from NUMERIC(20,4);
  v_multiplied boolean;
  v_any_rule   promo.rule_type;
  v_any_reward promo.reward_type;
  v_prev_value NUMERIC(20,4);
  v_prev_ord   int;
BEGIN
  -- ambil flag multiplied dari header (berdasarkan promo_id)
  SELECT promo_type, slab_multiplied
    INTO v_promo_type, v_multiplied
  FROM promo.promotions
  WHERE promo_id = NEW.promo_id;

  -- aturan multiplied
  IF COALESCE(v_multiplied, FALSE) THEN
    IF NEW.range_from IS NOT NULL THEN
      RAISE EXCEPTION 'SLAB: range_from must be NULL when slab_multiplied = true';
    END IF;
    IF NEW.reward_type = 'percentage' THEN
      RAISE EXCEPTION 'SLAB: percentage reward not allowed when slab_multiplied = true';
    END IF;
	ELSE 
		IF v_promo_type = 'slab' AND NEW.range_from IS NULL THEN
			RAISE EXCEPTION 'SLAB: range_from not allowed NULL when slab_multiplied = false';
		END IF;
		IF v_promo_type = 'slab' AND NEW.range_from = 0 THEN
			RAISE EXCEPTION 'SLAB: range_from must be > 0 when slab_multiplied = false';
		END IF;
  END IF;

  -- konsistensi rule_type
  SELECT s.rule_type INTO v_any_rule
  FROM promo.promotion_slabs s
  WHERE s.promo_id = NEW.promo_id
    AND s.id <> NEW.id
  LIMIT 1;

  IF v_any_rule IS NOT NULL AND v_any_rule <> NEW.rule_type THEN
    RAISE EXCEPTION 'SLAB: all slabs must use the same rule_type for a promotion';
  END IF;

  -- konsistensi reward_type
  SELECT s.reward_type INTO v_any_reward
  FROM promo.promotion_slabs s
  WHERE s.promo_id = NEW.promo_id
    AND s.id <> NEW.id
  LIMIT 1;

  IF v_any_reward IS NOT NULL AND v_any_reward <> NEW.reward_type THEN
    RAISE EXCEPTION 'SLAB: all slabs must use the same reward_type for a promotion';
  END IF;

  -- reward harus meningkat (bandingkan dgn ordinal sebelumnya)
  SELECT s.range_from, s.reward_value, s.ordinal
    INTO v_prev_range_from, v_prev_value, v_prev_ord
  FROM promo.promotion_slabs s
  WHERE s.promo_id = NEW.promo_id
    AND s.ordinal = (SELECT max(x.ordinal)
                     FROM promo.promotion_slabs x
                     WHERE x.promo_id = NEW.promo_id
                       AND x.ordinal < NEW.ordinal)
  LIMIT 1;

  IF v_prev_value IS NOT NULL AND NEW.reward_value IS NOT NULL
     AND NEW.reward_value <= v_prev_value THEN
    RAISE EXCEPTION 'SLAB: reward_value (ordinal %) must be > previous (ordinal %)', NEW.ordinal, v_prev_ord;
  END IF;
	
	IF v_prev_range_from IS NOT NULL AND NEW.range_from IS NOT NULL
     AND NEW.range_from < v_prev_range_from THEN
    RAISE EXCEPTION 'SLAB: range_from (ordinal %) must be >= previous (ordinal %)', NEW.ordinal, v_prev_ord;
  END IF;
	
	IF v_promo_type = 'slab' AND COALESCE(v_multiplied, FALSE) AND v_prev_value IS NOT NULL THEN
    RAISE EXCEPTION 'SLAB: slab item not allowed > 1 when slab_multiplied = true';
  END IF;

  RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_check_promotion_slabs ON promo.promotion_slabs;
CREATE TRIGGER trg_check_promotion_slabs
BEFORE INSERT OR UPDATE ON promo.promotion_slabs
FOR EACH ROW
EXECUTE FUNCTION promo.fn_check_promotion_slabs();


-- ================
-- STRATA (max 5 strata)
-- ================
DROP TABLE IF EXISTS promo.promotion_strata CASCADE;
CREATE TABLE IF NOT EXISTS promo.promotion_strata (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  promo_id       varchar(50) NOT NULL REFERENCES promo.promotions(promo_id) ON UPDATE CASCADE ON DELETE CASCADE,
  ordinal        INTEGER NOT NULL CHECK (ordinal BETWEEN 1 AND 5),
  description    VARCHAR(50),
  rule_type      promo.rule_type NOT NULL,            -- quantity/value
  rule_uom       promo.uom_type,             -- smallest/middle/largest
  range_from     NUMERIC(20,4) NOT NULL,
  range_to       NUMERIC(20,4) NOT NULL,
  CHECK (range_to > range_from),

  reward_type    promo.reward_type NOT NULL,          -- percentage/fixed_value/product
  reward_value   NUMERIC(20,4),                       -- % or fixed amount; NULL if product
  reward_uom     promo.uom_type,                      -- smallest/middle/largest
  per_scope      VARCHAR(16),                         -- 'per_product'|'per_order' if fixed_value chosen
  claimable      BOOLEAN NOT NULL DEFAULT FALSE,      -- claimable per strata
  claim_realization_pct NUMERIC(5,2),                 -- only if claimable && header claim_type=partial

  UNIQUE(promo_id, ordinal),

  CHECK (NOT (reward_type = 'percentage') OR (reward_value IS NOT NULL AND reward_value BETWEEN 0 AND 100)),
  CHECK (NOT (reward_type = 'fixed_value') OR (reward_value IS NOT NULL AND reward_value >= 0)),
  CHECK (NOT (reward_type = 'product')    OR reward_value IS NOT NULL),
  CHECK (claim_realization_pct IS NULL OR (claim_realization_pct >= 0 AND claim_realization_pct <= 100)),
  CHECK (per_scope IS NULL OR per_scope IN ('per_product','per_order'))
);
CREATE INDEX IF NOT EXISTS idx_promo_strata_promo ON promo.promotion_strata(promo_id);

-- Optional: trigger to cap total strata per promo to 5 (cannot be done with a plain CHECK)
-- You can add a BEFORE INSERT trigger to count rows per promotion_id.

-- Here’s a safe, concurrency-aware BEFORE INSERT trigger to cap strata at max 5 per promo_id. 
-- It locks the parent promo row so concurrent inserts on the same promo can’t slip past the limit.

-- Function: limit strata rows to <= 5 per promo_id
CREATE OR REPLACE FUNCTION promo.fn_limit_strata_to_five()
RETURNS TRIGGER AS $$
DECLARE
  v_cnt integer;
BEGIN
  -- Serialize writes per promo: lock the parent row
  PERFORM 1
  FROM promo.promotions
  WHERE promo_id = NEW.promo_id
  FOR UPDATE;

  -- Count existing strata for this promo
  SELECT COUNT(*) INTO v_cnt
  FROM promo.promotion_strata
  WHERE promo_id = NEW.promo_id;

  IF v_cnt >= 5 THEN
    RAISE EXCEPTION 'STRATA: maximum 5 strata allowed per promo_id (%).', NEW.promo_id;
  END IF;

  RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Trigger: run before each insert
DROP TRIGGER IF EXISTS trg_limit_strata_to_five ON promo.promotion_strata;
CREATE TRIGGER trg_limit_strata_to_five
BEFORE INSERT ON promo.promotion_strata
FOR EACH ROW
EXECUTE FUNCTION promo.fn_limit_strata_to_five();


-- ================
-- PRODUCT CRITERIA (step 3)
-- ================
DROP TABLE IF EXISTS promo.promotion_product_criteria CASCADE;
CREATE TABLE IF NOT EXISTS promo.promotion_product_criteria (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  promo_id       varchar(50) NOT NULL REFERENCES promo.promotions(promo_id) ON UPDATE CASCADE ON DELETE CASCADE,
  pro_id         bigint NOT NULL,                     -- FK to master product
  mandatory      BOOLEAN NOT NULL DEFAULT FALSE,
  min_buy_type   promo.rule_type,                     -- 'quantity' or 'value' (nullable if no min)
  min_buy_qty    NUMERIC(20,4),
  min_buy_value  NUMERIC(20,4),
  min_buy_uom     promo.uom_type DEFAULT NULL,         -- smallest/middle/largest
  UNIQUE(promo_id, pro_id),
  CHECK ( (min_buy_type IS NULL AND min_buy_qty IS NULL AND min_buy_value IS NULL)
       OR (min_buy_type='quantity' AND min_buy_qty  IS NOT NULL AND min_buy_qty  >= 0 AND min_buy_value IS NULL)
       OR (min_buy_type='value'    AND min_buy_value IS NOT NULL AND min_buy_value>= 0 AND min_buy_qty  IS NULL)
  ),
  -- filters shown in UI (principal/category/brand) are selection helpers; persisted items are in detail table below
  created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- ================
-- REWARD PRODUCT SETUP (step 4 when reward_type='product')
--    If a promo has product rewards, list them here (priority by ordinal).
-- ================
DROP TABLE IF EXISTS promo.promotion_reward_products;
CREATE TABLE IF NOT EXISTS promo.promotion_reward_products (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  promo_id       varchar(50) NOT NULL REFERENCES promo.promotions(promo_id) ON UPDATE CASCADE ON DELETE CASCADE,
  pro_id         bigint NOT NULL,
  ordinal        INTEGER NOT NULL,
  UNIQUE (promo_id, pro_id),
  UNIQUE (promo_id, ordinal)
);
CREATE INDEX IF NOT EXISTS idx_promo_reward_products_promo ON promo.promotion_reward_products(promo_id);

-- ================
-- COVERAGE (national or by distributor)
-- ================
DROP TABLE IF EXISTS promo.promotion_coverage_distributors;
CREATE TABLE IF NOT EXISTS promo.promotion_coverage_distributors (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  promo_id       varchar(50) NOT NULL REFERENCES promo.promotions(promo_id) ON UPDATE CASCADE ON DELETE CASCADE,
  distributor_id bigint NOT NULL,               -- FK to master distributor
  UNIQUE (promo_id, distributor_id)
);
CREATE INDEX IF NOT EXISTS idx_promo_cov_dist_promo ON promo.promotion_coverage_distributors(promo_id);

-- ================
-- 8) OUTLET CRITERIA (step 4/5 depending on reward type)
--    Either explicitly selected outlets OR attribute-based inclusion.
-- ================
DROP TABLE IF EXISTS promo.promotion_outlet_criteria;
CREATE TABLE IF NOT EXISTS promo.promotion_outlet_criteria (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  promo_id       varchar(50) NOT NULL REFERENCES promo.promotions(promo_id) ON UPDATE CASCADE ON DELETE CASCADE,
  selection_type   promo.outlet_sel_type NOT NULL DEFAULT 'by_attribute',  -- by_outlet / by_attribute

  -- When by_attribute, we use link tables below. When by_outlet, use *_selected table.
  created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Explicitly selected outlets (when selection_type='by_outlet')
CREATE TABLE IF NOT EXISTS promo.promotion_outlets_selected (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  criteria_id    varchar(30) NOT NULL REFERENCES promo.promotion_outlet_criteria(id) ON UPDATE CASCADE ON DELETE CASCADE,
  outlet_id      bigint NOT NULL,  -- FK to master outlet
  UNIQUE (criteria_id, outlet_id)
);
CREATE INDEX IF NOT EXISTS idx_promo_outlet_sel_crit ON promo.promotion_outlets_selected(criteria_id);

-- Attribute-based filters (many-to-many)
CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_class (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  criteria_id    varchar(30) NOT NULL REFERENCES promo.promotion_outlet_criteria(id) ON UPDATE CASCADE ON DELETE CASCADE,
  outlet_class_id  bigint NOT NULL,
  UNIQUE (criteria_id, outlet_class_id)
);
CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_group (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  criteria_id    varchar(30) NOT NULL REFERENCES promo.promotion_outlet_criteria(id) ON UPDATE CASCADE ON DELETE CASCADE,
  outlet_group_id  bigint NOT NULL,
  UNIQUE (criteria_id, outlet_group_id)
);
CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_type (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  criteria_id    varchar(30) NOT NULL REFERENCES promo.promotion_outlet_criteria(id) ON UPDATE CASCADE ON DELETE CASCADE,
  outlet_type_id   bigint NOT NULL,
  UNIQUE (criteria_id, outlet_type_id)
);
CREATE TABLE IF NOT EXISTS promo.promotion_outlet_attribute_sales_team (
  cust_id        varchar(10) NOT NULL,
  id             varchar(30) PRIMARY KEY DEFAULT public.generate_mongodb_id(),
  criteria_id    varchar(30) NOT NULL REFERENCES promo.promotion_outlet_criteria(id) ON UPDATE CASCADE ON DELETE CASCADE,
  sales_team_id    bigint NOT NULL,
  UNIQUE (criteria_id, sales_team_id)
);
