# A4 notes

Scope: A4 only.

Changed:
- `master/repository/product_repository.go:155-161`: `ReportList` now joins derived `md` relation, one row per `cust_id`, using `BOOL_OR(COALESCE(allow_upload_secondary_sales, false))`.
- `master/repository/product_report_repository_test.go:30-49, 76-92, 110-279`: updated SQLMock join expectations; added `TestProductReportRepository_DistributorAggregationPreservesCardinality` proving one count and one data row for aggregated distributor input.

Preserved:
- `filter.CustIDs` scope and SQL binding through `sqlx.In`/`Rebind`.
- Parent composite join: `parent.pro_id`, `parent.cust_id = LEFT(mp.cust_id, 6)`, active/non-deleted predicates.
- Mapping normalization fields and aggregated upload flag semantics.
- Controller → Service → Repository layering.
- No migrations, data writes, env, compose, dependency, or other-module edits.

Validation:
- `cd master && rtk go test ./repository -run 'Product.*Report' -v`: PASS, 9 tests.
- `cd master && rtk go test ./...`: PASS, 409 tests in 23 packages.
- `cd master && rtk go build ./...`: PASS.
- No write DB commands run.

Residual risks:
- Live Docker/DB cardinality proof not run in A4 lane; runtime evidence supplied before change remains source evidence for original duplication. SQLMock proves query shape and mapping scan cardinality only.
- Existing unrelated LSP diagnostics remain; package tests/build pass.
