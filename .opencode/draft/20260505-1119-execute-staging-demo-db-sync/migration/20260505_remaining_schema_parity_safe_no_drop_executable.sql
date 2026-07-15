-- Safe executable remainder: functions, triggers, and non-unique indexes only.
-- Excludes FK/PK/UNIQUE/CHECK constraints due existing demo data/name risks; no table/index/constraint removal.

BEGIN;

-- Functions
CREATE OR REPLACE FUNCTION inv.generate_stock_opname_doc_no(p_cust_id character varying, p_schedule_date date)
 RETURNS character varying
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_sequence_name VARCHAR(100);
    v_sequence_identifier VARCHAR(50);
    v_date_str VARCHAR(6);
    v_next_val INTEGER;
    v_doc_no VARCHAR(20);
    v_sanitized_cust_id VARCHAR(50);
BEGIN
    -- Format date as YYMMDD
    v_date_str := TO_CHAR(p_schedule_date, 'YYMMDD');
    
    -- Sanitize cust_id to be safe for use in sequence name
    v_sanitized_cust_id := UPPER(REGEXP_REPLACE(p_cust_id, '[^a-zA-Z0-9_]', '_', 'g'));
    
    -- Create sequence identifier
    v_sequence_identifier := 'stock_opname_seq_' || v_date_str || '_' || v_sanitized_cust_id;
    v_sequence_name := 'inv.' || v_sequence_identifier;
    
    -- Create sequence if it doesn't exist (using dynamic SQL with DO block)
    EXECUTE format('
        DO $create_seq$
        DECLARE
            seq_exists BOOLEAN;
        BEGIN
            SELECT EXISTS (
                SELECT 1 FROM pg_sequences 
                WHERE schemaname = ''inv'' 
                AND sequencename = %L
            ) INTO seq_exists;
            
            IF NOT seq_exists THEN
                EXECUTE format(''CREATE SEQUENCE inv.%%I START 1'', %L);
            END IF;
        END $create_seq$;
    ', v_sequence_identifier, v_sequence_identifier);
    
    -- Get next value from sequence
    EXECUTE format('SELECT nextval(%L)', v_sequence_name) INTO v_next_val;
    
    -- Format doc_no: SO + YYMMDD + 3-digit sequence number
    v_doc_no := 'SO' || v_date_str || LPAD(v_next_val::TEXT, 3, '0');
    
    RETURN v_doc_no;
END;
$function$;

CREATE OR REPLACE FUNCTION mst.update_route_outlet_from_outlet()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    -- Gabungkan semua kondisi dalam satu IF
    IF NEW.outlet_name IS DISTINCT FROM OLD.outlet_name OR
       NEW.outlet_code IS DISTINCT FROM OLD.outlet_code OR
       NEW.address1    IS DISTINCT FROM OLD.address1 OR
       NEW.avg_sales_week IS DISTINCT FROM OLD.avg_sales_week OR
       NEW.latitude    IS DISTINCT FROM OLD.latitude OR
       NEW.longitude   IS DISTINCT FROM OLD.longitude THEN

        -- Update ke pjp.route_outlet berdasarkan outlet_id
        UPDATE pjp.route_outlet
        SET
            outlet_name     = NEW.outlet_name,
            outlet_code     = NEW.outlet_code,
            address1        = NEW.address1,
            avg_sales_week  = NEW.avg_sales_week,
            latitude        = NEW.latitude,
            longitude       = NEW.longitude
        WHERE outlet_id = NEW.outlet_id;
    END IF;

    RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION pjp.update_route_name_on_change()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
  -- Update route_outlet_history
  UPDATE pjp.route_outlet_history
  SET route_name = NEW.route_name
  WHERE route_code = NEW.route_code;

  -- Update route_outlet
  UPDATE pjp.route_outlet
  SET route_name = NEW.route_name
  WHERE route_code = NEW.route_code;

  RETURN NEW;
END;
$function$;

CREATE OR REPLACE FUNCTION promo.fn_check_promotion_slabs()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION promo.fn_limit_strata_to_five()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION public.generate_mongodb_id()
 RETURNS text
 LANGUAGE plpgsql
AS $function$
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
$function$;

CREATE OR REPLACE FUNCTION sls.generate_invoice_no(p_cust_id character varying, p_invoice_date date DEFAULT CURRENT_DATE)
 RETURNS character varying
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_seq INTEGER;
    v_existing_max_seq INTEGER;
BEGIN
    IF p_cust_id IS NULL OR LENGTH(TRIM(p_cust_id)) = 0 THEN
        RAISE EXCEPTION 'p_cust_id is required';
    END IF;

    SELECT COALESCE(MAX(SUBSTRING(o.invoice_no FROM 10 FOR 4)::INTEGER), 0)
    INTO v_existing_max_seq
    FROM sls."order" o
    WHERE o.cust_id = p_cust_id
      AND o.invoice_no IS NOT NULL
      AND o.invoice_no LIKE CONCAT('INV', TO_CHAR(p_invoice_date, 'YYMMDD'), '%');

    INSERT INTO sls.invoice_no_counter AS c (cust_id, seq_date, last_seq, updated_at)
    VALUES (p_cust_id, p_invoice_date, v_existing_max_seq + 1, NOW())
    ON CONFLICT (cust_id, seq_date)
    DO UPDATE
    SET last_seq = GREATEST(c.last_seq, v_existing_max_seq) + 1,
        updated_at = NOW()
    RETURNING last_seq INTO v_seq;

    RETURN CONCAT('INV', TO_CHAR(p_invoice_date, 'YYMMDD'), LPAD(v_seq::TEXT, 4, '0'));
END;
$function$;

-- Triggers
CREATE TRIGGER trg_update_route_name AFTER UPDATE OF route_name ON pjp.routes FOR EACH ROW EXECUTE FUNCTION pjp.update_route_name_on_change();
CREATE TRIGGER trg_check_promotion_slabs BEFORE INSERT OR UPDATE ON promo.promotion_slabs FOR EACH ROW EXECUTE FUNCTION promo.fn_check_promotion_slabs();
CREATE TRIGGER trg_limit_strata_to_five BEFORE INSERT ON promo.promotion_strata FOR EACH ROW EXECUTE FUNCTION promo.fn_limit_strata_to_five();

-- Non-unique indexes
CREATE INDEX IF NOT EXISTS idx_expense_source ON acf.expense USING btree (source);
CREATE INDEX IF NOT EXISTS idx_expense_det_collector_id ON acf.expense_det USING btree (collector_id);
CREATE INDEX IF NOT EXISTS idx_expense_det_expense_type_id ON acf.expense_det USING btree (expense_type_id);
CREATE INDEX IF NOT EXISTS idx_expense_type_cust_active_del ON acf.expense_type USING btree (cust_id, is_active, is_del);
CREATE INDEX IF NOT EXISTS idx_expense_type_cust_id ON acf.expense_type USING btree (cust_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_cust_date ON inv.replenishment_order USING btree (cust_id, date);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_distributor_id ON inv.replenishment_order USING btree (distributor_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_replenishment_no ON inv.replenishment_order USING btree (replenishment_no);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_ro_id ON inv.replenishment_order USING btree (replenishment_no);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_so_no ON inv.replenishment_order USING btree (so_no);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_sup_id ON inv.replenishment_order USING btree (sup_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_wh_id ON inv.replenishment_order USING btree (wh_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_detail_pro_id ON inv.replenishment_order_detail USING btree (pro_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_order_detail_replenishment_id ON inv.replenishment_order_detail USING btree (replenishment_id);
CREATE INDEX IF NOT EXISTS idx_replenishment_status_name ON inv.replenishment_status USING btree (status_name);
CREATE INDEX IF NOT EXISTS idx_replenishment_type_name ON inv.replenishment_type USING btree (replenishment_type_name);
CREATE INDEX IF NOT EXISTS idx_so_bulk_upload_doc_no ON inv.stock_opname_bulk_upload USING btree (doc_no);
CREATE INDEX IF NOT EXISTS idx_so_bulk_upload_items_upload_id ON inv.stock_opname_bulk_upload_items USING btree (upload_id);
CREATE INDEX IF NOT EXISTS idx_visits_cust_emp_outlet_created_visit ON mobile.visits USING btree (cust_id, emp_code, outlet_code, created_at DESC, visit_id DESC);
CREATE INDEX IF NOT EXISTS idx_m_distributor_parent_cust_id ON mst.m_distributor USING btree (parent_cust_id);
CREATE INDEX IF NOT EXISTS m_employee_emp_id_idx ON mst.m_employee USING btree (emp_id);
CREATE INDEX IF NOT EXISTS idx_m_outlet_code_created_at ON mst.m_outlet_code USING btree (created_at);
CREATE INDEX IF NOT EXISTS idx_m_outlet_code_cust_id ON mst.m_outlet_code USING btree (cust_id);
CREATE INDEX IF NOT EXISTS idx_m_outlet_code_status ON mst.m_outlet_code USING btree (status);
CREATE INDEX IF NOT EXISTS idx_m_salesman_cust_emp_lookup ON mst.m_salesman USING btree (cust_id, emp_id) WHERE ((is_del = false) AND (deleted_at IS NULL));
CREATE INDEX IF NOT EXISTS idx_m_survey_area_distributor_id ON mst.m_survey_area USING btree (distributor_id);
CREATE INDEX IF NOT EXISTS idx_arrival_report_created_at ON pjp.arrival_report USING btree (created_at);
CREATE INDEX IF NOT EXISTS idx_arrival_report_cust_id ON pjp.arrival_report USING btree (cust_id);
CREATE INDEX IF NOT EXISTS idx_arrival_report_outlet_id ON pjp.arrival_report USING btree (outlet_id);
CREATE INDEX IF NOT EXISTS idx_arrival_report_user_id ON pjp.arrival_report USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_ovl_date_pjp_outlet_extra ON pjp.outlet_visit_list USING btree (date, pjp_id, outlet_id, is_extra_call);
CREATE INDEX IF NOT EXISTS outlet_visit_list_date_idx ON pjp.outlet_visit_list USING btree (date);
CREATE INDEX IF NOT EXISTS outlet_visit_list_outlet_idx ON pjp.outlet_visit_list USING btree (outlet_id);
CREATE INDEX IF NOT EXISTS idx_roh_date_cust_pjp_outlet_extra_route ON pjp.route_outlet_history USING btree (date, cust_id, pjp_id, outlet_id, is_extra_call, route_code);

COMMIT;
