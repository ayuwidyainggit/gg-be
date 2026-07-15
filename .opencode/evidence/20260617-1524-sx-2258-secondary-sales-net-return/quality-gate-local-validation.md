# SX-2258 quality gate after local validation

Status: `PASS_WITH_RISKS`

## Decision

- Local scope validated end-to-end.
- Prior "no DB/API evidence" risk is remediated for local `ggn_scyllax` + local Docker API.
- Remaining risk is environment data mismatch only, not code failure.
- Full cross-env `PASS` still needs staging/QA dataset proof because local snapshot does not contain the original Jira expected numbers.

## Local validation accepted by gate

- Direct SQL monthly summary for `C260020001`, June 2026:
  - `qty=265`
  - `qty_return=12`
  - `total_discount_promo=1403740`
  - `net_sales=6004889286`
- API/cURL principal child scope matched DB.
- API/cURL distributor own scope matched DB.
- Distributor parent scope denied with `403`.
- Optional filter DB/API matched:
  - `qty=60`
  - `total_discount_promo=0`
  - `qty_return=0`
- Trend June `total_discount_promo=1403740`.

## Remaining risk

- Exact Jira expected values are not reproducible on local `ggn_scyllax` snapshot:
  - expected `qty=134`
  - expected `total_discount_promo=1238740`
- Broad local DB search found no monthly row matching those expected values.
- Required for full cross-env `PASS`: rerun same sanitized SQL/API checks on staging/QA dataset or get QA confirmation of local dataset drift.

## Code status

No code remediation required from this gate.
