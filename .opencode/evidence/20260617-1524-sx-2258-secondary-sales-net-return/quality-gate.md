# SX-2258 quality gate

Status: `PASS_WITH_RISKS`

## Scope reviewed

- `sales/entity/report.go`
- `sales/service/report_service.go`
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`
- `sales/service/report_service_test.go`

## Pass basis

- Summary `qty` is now net sold qty:
  - `COALESCE(os.qty, 0) - COALESCE(rs.qty_return, 0) AS qty`
- Summary `total_discount_promo` is now net:
  - `COALESCE(os.discount_promo, 0) - COALESCE(rs.discount_promo, 0) AS total_discount_promo`
- Trend `total_discount_promo` is now net.
- Order discount/promo uses docs/Jira formula:
  - `COALESCE(od.disc_value_final, 0) + COALESCE(od.promo_final1, 0) + COALESCE(od.promo_final2, 0) + COALESCE(od.promo_final3, 0) + COALESCE(od.promo_final4, 0) + COALESCE(od.promo_final5, 0)`
- Return discount/promo uses:
  - `COALESCE(rd.disc_value, 0) + COALESCE(rd.promo_value, 0)`
- Return invoice join preserved.
- Optional summary filters use correct aliases and parameterized SQL.
- PPN, return value, return rate, and `qty_return` semantics preserved.
- Tests cover old formula rejection, optional filters, arithmetic regression, and trend formula.
- `pjp-sales` not changed.
- No secrets or `.env` touched.

## Validation evidence

- `rtk docker compose -f docker-compose.yml ps` ran; current env showed empty compose service table.
- `rtk go test ./repository -run 'SecondarySalesReport(SumReportByMonth|TrendSales)|SX2258'` passed.
- `rtk go test ./service -run SecondarySalesReportSumReportByMonth` passed.
- `rtk go test ./...` passed: `265 passed in 22 packages`.

## Residual risk

- Exact QA dataset proof is missing because safe internal DB/API access is not available in this session.
- Required before full `PASS`: run safe staging/local API or SQL check for Jira dataset and confirm:
  - `qty = 134`
  - `total_discount_promo = 1238740`

## Remediation status

No code remediation required. Only runtime proof gap remains.
