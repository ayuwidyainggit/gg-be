# A5 Runtime — Parent-Eligibility Fallback

Date: 2026-07-14
Claim level: `confirmed_runtime`

## Source guard

`master/repository/product_repository.go:143`:
```sql
mp.parent_pro_id <> 0
AND mp.is_product_mapping = true
AND md.allow_upload_secondary_sales = true
AND parent.pro_id IS NOT NULL
```

Without the guard, when parent row was missing or inactive, the `CASE` chose
`parent.pro_id` (NULL). Response scan into `ProductID int64` failed at runtime
with `gagal memecahkan kode: skema: kesalahan mengonversi nilai untuk pro_id`.
Source: production response on `C260020001`.

## SQLMock regression

`master/repository/product_report_repository_test.go:110-139`
(`TestProductReportRepository_NormalizationBranches`):

- `parent_missing_or_inactive` row: primary fields = `mp` (`pro_id=11`),
  `original_* = NULL`, `type='Product Mapping'`.
- `parent_eligible` row: primary fields = parent (`cust_id='C26002'`,
  `pro_id=7`), `original_*` set to source `mp` values.

Result: 10 repository report tests pass.

## Live local fixture (Docker + ggn_scyllax)

- `mst.m_product` rows with `is_product_mapping=true` and missing eligible
  parent: 0 in this dataset.
- `mst.m_product` rows with `is_product_mapping=true` and eligible parent
  (`C260020004`): 6 rows. Live curl pre-A5 normalized these to parent fields.
  Same SQL path now guarded; behaviour identical when parent exists.
- Local smoke cannot reproduce the production scan crash directly because
  no row in the dataset exercises `parent.pro_id IS NULL`. The production
  crash is closed at the source via the guard, not by live reproduction.

## Live production revalidation

Not performed. Bearer token from user prompt had `expires=1784091462`
(~2026-07-15) and was rejected as `401 Unauthorized` when re-sent. Production
token must be rotated by the user; new token required to call
`https://best.scyllax.online/master/v1/products/report?cust_id=C260020001&page=1&limit=20&sort_by=pro_name&sort_order=asc`
and confirm no conversion error.

## Residual limitation

- A3 fallback rule is closed by guard + SQLMock + local fixture unchanged.
- Production branch can only be re-validated by issuing a fresh bearer
  token and rerunning the original curl.
