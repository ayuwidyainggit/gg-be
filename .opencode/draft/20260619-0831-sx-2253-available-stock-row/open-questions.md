# Open Questions SX-2253

- Apakah `RefreshOrderDetailStock` pernah menulis `qty*_stok` yang salah saat row lain dengan `pro_id` sama diproses (mis. cancellation, mobile sync, atau path lain)? Trace `rg` di T1: tidak ada call site `RefreshOrderDetailStock` di `sales/`. Repository method hanya ada di interface/impl, tidak dipanggil dari service. **Resolved T1: not used in sales/ service.**
- Apakah ada path upsert/persist lain (selain `RefreshOrderDetailStock`) yang menulis `qty1_stok/qty2_stok/qty3_stok`? Trace `rg` di T1: satu-satunya writer di sales adalah `sales/service/order_service.go:5555-5565` di `UpdateEnhance`-style path, memakai `currentStock` saja tanpa row qty. **Resolved T1: single writer in service, but DetailV2 recomputes from scratch so display is row-keyed.**
- Apakah FE menghitung `available_stock` dari jumlah `qty*_stok` lintas row, atau murni per row? FE comment di Jira bilang tidak ada kalkulasi FE, hanya menampilkan API. **Out of scope for BE worklist.**

## Open question raised by T2-T4 results

- Premise plan: "rows that share the same `pro_id` end up with the same displayed `available_stock`". Not reproducible in `DetailV2` mapping using T2-T4 test inputs (pro_id=8435, wh=100, rows with qty1=10 and qty1=5). Current `DetailV2` returns different `Qty1Stok/Qty2Stok/Qty3Stok` per row even before any code change. Need a concrete failing production repro (call path, input SO/pro_id/qty, expected vs actual) that exhibits the bug for `@oracle` to plan a precise fix.
