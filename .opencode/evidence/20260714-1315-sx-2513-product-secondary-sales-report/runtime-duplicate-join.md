# Runtime duplicate-join evidence — SX-2513 remediation

## Claim level

`confirmed_runtime` from live `ggn_scyllax` psql and endpoint output supplied in delegation payload.

## Observed facts

| Check | Result | Consequence |
|---|---:|---|
| Active `mst.m_product` rows for `cust_id='C22001'` | 3917 | Expected product-row cardinality before report pagination. |
| `mst.m_distributor` rows for `cust_id='C22001'` | 38 | `LEFT JOIN mst.m_distributor md ON md.cust_id=mp.cust_id` produces up to 38 copies per product row. |
| Endpoint `total_record` | 148846 | Equals `3917 × 38`; count is multiplied. |
| `pro_id=495` observed in first response page | 5 copies | Data query also multiplies rows; page can contain duplicate products. |

## Root cause

Original plan required direct `LEFT JOIN mst.m_distributor md ON md.cust_id=mp.cust_id` at plan line 231. `mst.m_distributor.cust_id` is not one-row-per-customer in live data. Joining it directly violates original Requirement 15: count/data must represent `mp` result rows before pagination.

## Chosen remediation semantics

Join one derived row per `cust_id`:

```sql
LEFT JOIN (
  SELECT
    cust_id,
    BOOL_OR(COALESCE(allow_upload_secondary_sales, false)) AS allow_upload_secondary_sales
  FROM mst.m_distributor
  GROUP BY cust_id
) md ON md.cust_id = mp.cust_id
```

Semantics: a customer allows secondary-sales upload when at least one of its distributor rows has `allow_upload_secondary_sales=true`. No matching distributor row behaves as false via `COALESCE(md.allow_upload_secondary_sales, false)` in mapping CASE predicates. This is smallest one-row-per-customer relation preserving feature input needed by report: mapping/upload boolean only.

## Explicit non-solutions

- Do not use `DISTINCT` in count/data queries. It masks join cardinality and can break deterministic paging.
- Do not select arbitrary distributor row with `MIN`, `MAX`, or `DISTINCT ON`; that changes flag truth semantics.
- Do not change request scope, parent composite join, mappings, migrations, data, compose/env, or modules.

## Required proof after remediation

1. SQLMock query assertion proves derived grouped `md` relation exists and direct raw `md.cust_id=mp.cust_id` join does not.
2. SQLMock cardinality regression: one product mapping input with 38 distributor records simulated by aggregate result yields one scanned report row and count 1.
3. Docker+DB read-only proof for `C22001`: `COUNT(*)` from report query base equals 3917; no duplicate `pro_id` group has `COUNT(*) > 1`.
4. Authenticated endpoint proof: `total_record=3917` for `cust_id[]=C22001`; first page has unique `pro_id` values for returned rows.

## Source strategy

Live runtime psql/curl evidence is primary. Repo query path was pre-inspected in `master/repository/product_repository.go:143-213`. External docs skipped: PostgreSQL `BOOL_OR` is existing database aggregate semantics; no dependency/API/version decision is introduced.
