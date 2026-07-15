# Staging Evidence — SX-2314 Retry v2

Status: `SOFT_BLOCKED`

Staging DB credentials/connection details not available in current workspace environment. No read-only staging queries executed.

Required queries when staging access is available:
1. Legacy reversal row inventory:
```sql
SELECT tr_code, tr_no, COUNT(*)
FROM inv.stock
WHERE tr_no LIKE '%-CO'
GROUP BY tr_code, tr_no
ORDER BY 1;
```
2. Reward detail qty for QA SOs:
```sql
SELECT order_detail_id, item_type, pro_id, qty, qty_final, qty_po,
       qty1, qty2, qty3, qty1_final, qty2_final, qty3_final,
       qty_po1, qty_po2, qty_po3, conv_unit2, conv_unit3
FROM sls.order_detail
WHERE ro_no IN ('SO2606230004','SO2606240002','SO2606260002')
ORDER BY ro_no, order_detail_id;
```

Mitigation: implementation uses residual math with clamp-to-zero and legacy cancelAgg guard so staging data profile cannot cause over-reversal. Validate locally when staging access available before prod deploy.
