# A4 Runtime Validation — SX-2513

Date: 2026-07-14
Claim level: `confirmed_runtime`

## Environment

- Docker service: `scylla-master`, exposed at `localhost:9002`.
- Database: local PostgreSQL `ggn_scyllax`.
- Request auth: `Cust_id` header through existing middleware fallback. No bearer token captured or recorded.
- No database writes performed.

## Cardinality regression: C22001

Database source count:

```sql
SELECT COUNT(*)
FROM mst.m_product
WHERE cust_id = 'C22001' AND is_del = false AND is_active = true;
```

Result: `3917`.

Request:

```text
POST /v1/products/report?cust_id[]=C22001&limit=100&sort_by=pro_id&sort_order=asc
```

Result: HTTP `200`; `paging.total_record=3917`; returned `100` rows; `100` distinct `pro_id`; `0` duplicate `pro_id` values.

This proves aggregate distributor relation prevents earlier `38x` multiplication from 38 `mst.m_distributor` rows for `C22001`.

## Mapping normalization: C260020004

Read-only DB query confirmed six active mapping products. All use `BOOL_OR(COALESCE(md.allow_upload_secondary_sales,false)) = true`; each eligible parent exists at `cust_id='C26002'` using composite `(parent_pro_id, LEFT(mp.cust_id,6))` identity.

Request:

```text
POST /v1/products/report?cust_id[]=C260020004&limit=100&sort_by=pro_id&sort_order=asc
```

Result: HTTP `200`; `paging.total_record=9`; six `type='Product Mapping'` rows. Each mapping row used parent principal identity in primary fields, with mapping source retained in `original_*` fields.

## Residual risk

A3 remains conditional: a live mapping row with unavailable eligible parent was not found in this fixture. Plan semantics prevent normalizing to stale parent; product owner decision remains needed only if such live data appears.
