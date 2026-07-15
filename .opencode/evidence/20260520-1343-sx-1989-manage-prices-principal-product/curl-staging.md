# SX-1989 — Curl Staging Evidence

Date: 2026-05-20
Tester: Kiro (automated via OpenCode)
User: adminbm@gmail.com (distributor_id=102, cust_id=C260020001, parent_cust_id=C26002)

## Pre-fix (remote best.scyllax.online, kode lama)

### order_date=2026-05-16
```json
{
  "pro_code": "JY1-005",
  "purch_price1": 60000,
  "sell_price1": 1000000
}
```
Harga lama dari m_transaction_price start_date=2026-05-13 (parent scope, coverage=D, dist=102).

### order_date=2026-05-20
```json
{
  "pro_code": "JY1-005",
  "purch_price1": 15000,
  "sell_price1": 12000
}
```
Harga terbaru dari m_transaction_price start_date=2026-05-20 (parent scope, coverage=N).

## Post-fix (localhost:9002, kode baru dengan UNION ALL child+parent generic)

### order_date=2026-05-16 — FIXED ✅
```json
{
  "pro_code": "JY1-005",
  "purch_price1": 50000,
  "purch_price2": 500000,
  "purch_price3": 500000,
  "sell_price1": 150000,
  "sell_price2": 1500000,
  "sell_price3": 1500000
}
```
Harga dari child generic row: cust_id=C260020001, pro_id=8457, start_date=2026-05-16, coverage=N, scope_priority=0.
Child generic row menang atas parent old row (start_date=2026-05-13) karena start_date lebih baru.

## DB Audit

### Pre-backfill
- broken_child_count: 39 rows (distributor_id IS NOT NULL AND COALESCE(parent_pro_id,0)=0 AND is_del=false)
- Breakdown per tenant: C260020001=23, C260040001=9, C220010001=4, C220010002=3

### Backfill result (2026-05-20 08:47 WIB)
- FIXABLE: 3 rows (unambiguous 1:1 mapping)
- AMBIGUOUS: 0 rows
- Updated: pro_id=10776 (AF-007), 10777 (AF-008), 10778 (AF-009) → parent_pro_id set

### Post-backfill
- remaining_broken: 36 rows (sisa adalah produk tanpa child row atau tanpa parent match — masuk manual review queue)

## Root cause summary

Query `FindAllByDistributorLookupDistPrice` sebelumnya hanya membaca parent generic
`mst.m_transaction_price` (cust_id=parent, pro_id=pricing_lookup_pro_id).
Tidak membaca child generic rows (cust_id=child, pro_id=child_pro_id, outlet_id=0).

Akibatnya untuk order_date=2026-05-16:
- child generic row (start_date=2026-05-16, sell_price1=150000) diabaikan
- parent old row (start_date=2026-05-13, sell_price1=1000000) dipakai

Fix: UNION ALL child generic + parent generic, ORDER BY start_date DESC, scope_priority ASC.
Child (scope_priority=0) menang atas parent (scope_priority=1) jika start_date sama.
Latest start_date menang jika berbeda.

## Validation commands run

```bash
# Local tests
rtk go test ./repository -run TestFindAllByDistributorLookupDistPrice  # 3/3 pass
rtk go test ./...  # 267/267 pass

# DB simulation (staging)
-- query UNION ALL → effective_purch_price1=50000, effective_sell_price1=150000, scope_priority=0

# Curl local (post-fix)
curl http://localhost:9002/v1/products?mode=lookup_dist_price&...&order_date=2026-05-16&q=jersey+persinga
→ purch_price1=50000, sell_price1=150000 ✅
```
