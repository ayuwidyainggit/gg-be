# Local Smoke + DB Validation SX-2079

Date: 2026-05-26
Target DB: `ggn_scyllax`
Target service: local `master` on `http://127.0.0.1:9002`

## Runtime checks

```bash
rtk docker compose -f docker-compose.yml config --services
# rabbitmq redis system cronjob finance master pjp sales tms inventory mobile

rtk docker compose -f docker-compose.yml up -d master
# Container scylla-master Running

curl -sS "http://127.0.0.1:9002/ping"
# It works

PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -d ggn_scyllax -c "SELECT current_database(), current_user;"
# ggn_scyllax | postgres
```

## Test user + fixture strategy

- Principal user used for smoke: `user_id=1`, `user_name=princ@idetama.id`, `cust_id=C22001`, `emp_id=278`.
- JWT dibuat lokal dengan secret lokal `master/.env` dan claim principal (`distributor_id=0`). Secret tidak ditulis ulang di evidence ini.
- Kondisi awal DB untuk employee `C22001/278`:
  - `region_scope=ALL`
  - `area_scope=ALL`
  - `distributor_scope=ALL`
  - region mapping count `0`
  - area mapping count `0`
  - distributor mapping count `0`
- Untuk validasi specific scope dipakai fixture sementara di DB lokal:
  - update scope employee `278`
  - insert mapping rows temporary
  - hit API
  - compare hasil API vs query SQL langsung
  - cleanup via trap sesudah test

## Baseline ALL/ALL/ALL assertions

SQL baseline:
- active regions by `cust_id='C22001'` = `66,67,68,70,71,72,73,74,75,76,77,78,79,80`
- active areas by `cust_id='C22001'` = `70,82,83,84,85,86,89`
- active distributors by `parent_cust_id='C22001'` = set `42,44,45,67,68,69,70,86,107,108,114,115,116,117,118` plus current local extras exactly matched by API set comparison

API baseline matched DB set comparison:
- `GET /v1/regions?q=&page=1&limit=9999&sort=region_id:asc&is_active=1`
- `GET /v1/areas?q=&page=1&limit=9999&sort=area_id:asc&is_active=1`
- `GET /v1/business-unit?q=&page=1&limit=9999&sort=area_id:asc&is_active=1`

Observed API IDs:
- `baseline_regions=66,67,68,70,71,72,73,74,75,76,77,78,79,80`
- `baseline_areas=70,82,83,84,85,86,89`
- `baseline_business_units=70,42,44,45,116,67,68,115,108,112,118,69,114,86,117,107`

Note: business-unit compare dipakai **set comparison**, bukan stable order compare, karena query hanya sort by `area_id` and intra-area order tidak selalu deterministik.

## Specific scope smoke assertions

### 1. Region specific

Temporary DB fixture:
- `region_scope='SELECTED'`
- region mappings: `66`, `67`
- `area_scope='ALL'`
- `distributor_scope='ALL'`

Expected vs actual:
- region API returned `66,67`
- area API fallback by mapped region returned `70,82,83,85`
- business-unit API fallback by mapped region returned distributor set from `region_id IN (66,67)`

Observed:
- `region_specific=66,67`
- `area_fallback_from_region=70,82,83,85`
- `business_units_from_region=42,44,45,63,66,116,67,68,101,99,100,109,110,115,114,119`

### 2. Area specific + explicit region intersection

Temporary DB fixture:
- `region_scope='ALL'`
- `area_scope='SELECTED'`
- area mappings: `70`, `85`
- `distributor_scope='ALL'`

Expected vs actual:
- area API returned only mapped areas `70,85`
- business-unit API with explicit `region_id=67` returned intersection of mapped areas and explicit region filter

Observed:
- `area_specific=70,85`
- `business_units_area_intersection=42,44,45,66`

### 3. Distributor specific + explicit region/area intersection

Temporary DB fixture:
- `region_scope='ALL'`
- `area_scope='ALL'`
- `distributor_scope='SELECTED'`
- distributor mappings: `67`, `69`

Expected vs actual:
- business-unit API returned only mapped distributors
- business-unit API with explicit `region_id=68&area_id=85` returned intersection only `69`

Observed:
- `business_units_specific=67,69`
- `business_units_specific_intersection=69`

## Smoke command outcome

End-to-end smoke script result:

```text
SMOKE_OK
baseline_regions=66,67,68,70,71,72,73,74,75,76,77,78,79,80
baseline_areas=70,82,83,84,85,86,89
baseline_business_units=70,42,44,45,116,67,68,115,108,112,118,69,114,86,117,107
region_specific=66,67
area_fallback_from_region=70,82,83,85
business_units_from_region=42,44,45,63,66,116,67,68,101,99,100,109,110,115,114,119
area_specific=70,85
business_units_area_intersection=42,44,45,66
business_units_specific=67,69
business_units_specific_intersection=69
```

## Cleanup verification

Sesudah trap cleanup:

```sql
SELECT region_scope, area_scope, distributor_scope
FROM mst.m_employee
WHERE cust_id='C22001' AND emp_id=278;
-- ALL | ALL | ALL

SELECT 'region', COUNT(*) FROM mst.m_employee_region_mapping WHERE cust_id='C22001' AND emp_id=278 AND is_del=false
UNION ALL
SELECT 'area', COUNT(*) FROM mst.m_employee_area_mapping WHERE cust_id='C22001' AND emp_id=278 AND is_del=false
UNION ALL
SELECT 'distributor', COUNT(*) FROM mst.m_employee_distributor_mapping WHERE cust_id='C22001' AND emp_id=278 AND is_del=false;
-- all counts = 0
```

## Conclusion

- Local API smoke against local DB `ggn_scyllax` passed.
- Baseline ALL behavior matched DB.
- Temporary specific-scope fixtures validated principal scope enforcement for:
  - region specific
  - area specific
  - distributor specific
  - explicit query intersection safety
- Local DB restored to original tested state for employee `278`.
