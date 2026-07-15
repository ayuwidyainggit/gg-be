# SX-2258 final quality gate after staging data sync

Status: `PASS`

## Decision

- Synced local data reproduced exact SX-2258 expected values.
- DB direct SQL and local API/cURL match.
- No remaining code/data validation blocker.

## Pass basis

- Local DB `ggn_scyllax` expected dataset found:
  - `cust_id=C260020001`
  - `invoice_date=2026-06-03`
- Direct SQL:
  - `qty=134`
  - `total_discount_promo=1238740`
  - `qty_return=12`
  - `net_sales=5720067386`
- Principal cURL:
  - status `200`
  - `qty=134`
  - `total_discount_promo=1238740`
  - `qty_return=12`
  - `net_sales=5720067386`
- Distributor cURL:
  - status `200`
  - `qty=134`
  - `total_discount_promo=1238740`
  - `qty_return=12`
  - `net_sales=5720067386`
- Distributor parent scope remains denied:
  - status `403`
  - `cust_id is outside authorized scope`
- Tokens/passwords not stored in evidence.

## Remediation

None.
