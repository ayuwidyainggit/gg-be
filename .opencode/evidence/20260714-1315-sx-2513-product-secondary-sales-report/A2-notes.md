# A2 audit notes

## Verified implementation

- `master/repository/product_repository.go:145-214`
  - `ReportList` uses `sqlx.In` then `r.Rebind` for explicit caller `CustIDs` (`:160-173`, `:205-209`).
  - `q` is passed as two bound wildcard arguments (`:161-165`); no request value is interpolated into SQL syntax.
  - Count and data queries reuse the same `from`, `where`, and base argument set (`:155-175`, `:192-213`).
  - Composite parent eligibility join is `parent.pro_id = mp.parent_pro_id` and `parent.cust_id = LEFT(mp.cust_id, 6)` plus active/non-deleted predicates (`:157-159`).
  - `limit` and `offset` are bound (`:202-205`).
  - Closed in-process `sortColumns` map supplies output aliases and `sortOrder` normalizes to `ASC`/`DESC` (`:177-190`); no raw request string is appended.
- `master/service/product_service.go:104-106` delegates to repository. A1 temporary empty result no longer exists.
- `master/repository/product_report_repository_test.go:30-279` covers bound count/data arguments, composite parent join text, five output scan categories, nullable `original_*`, and `ExpectationsWereMet`.

## Validation

`confirmed_runtime` from `A2-repository.log`:

```text
rtk go build ./...                                  -> Success
rtk go test ./controller -run 'Product.*Report' -v  -> 5 passed in 1 packages
rtk go test ./repository -run 'Product.*Report' -v  -> 8 passed in 1 packages
```

## A2 remediation

- `master/controller/product_report_controller_test.go:19-23` updated `productServiceStub.ReportList` to return `data, total, lastPage, error`, with `0, 0` pagination defaults.
- `master/controller/product_controller.go:578-592` already consumed `total` and `lastPage` and populated `entity.Pagination`; no production change required.
- `rg -n "ReportList" master/` found no other old-signature `ProductService.ReportList` or `ProductRepository.ReportList` caller.

## Open assumption

A3 remains: an eligible mapping row whose parent is missing/inactive has no product-owner fallback decision (plan assumption A3). The current LEFT JOIN naturally returns null primary fields for the enabled-normalization CASE. Runtime smoke is not ready because compose has zero running services and no token has been verified.

## Scope

No migrations, JWT middleware, package manifests, compose/env, or other modules were changed by A2.
