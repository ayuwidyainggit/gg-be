# DB Validation — Task 120/121 Secondary Sales `cust_id` Filters

Task id: `20260520-1851-secondary-sales-cust-id-filters`
Tanggal: `2026-05-20`
Target DB: `scylla_citus_dev` (PostgreSQL 14.5, host `103.28.219.73:25431`, user `postgres`, `sslmode=disable`).
Mode: read-only `SELECT` only. Password tidak ditulis di artefak; gunakan `PGPASSWORD` env atau secret out-of-band.

## Connectivity

```text
SELECT current_database(), current_user, version();

scylla_citus_dev|postgres|PostgreSQL 14.5 (Ubuntu 14.5-1.pgdg20.04+1) on x86_64-pc-linux-gnu, ...
```

## Schema availability

Semua tabel relevan untuk Export dan Trend Sales hadir di DB dev:

```text
mst.m_distributor
mst.m_outlet
mst.m_product
mst.m_salesman
mst.m_supplier
report.dim_dates
report.dim_outlets
report.dim_product_categories
report.dim_products
report.dim_salesmans
report.fact_orders
report.fact_returns
report.list
sls.order
sls.order_detail
sls.return
sls.return_det
smc.m_customer
```

## Kolom kunci

`smc.m_customer`:

```text
cust_id, parent_cust_id, is_del, is_active
```

`report.fact_orders`:

```text
cust_id, date_id, salesman_id, outlet_id, pro_id,
gross_sale, discount, special_discount, net_sales_exclude_ppn,
qty, extracted_at
```

`report.list`:

```text
cust_id varchar
report_id varchar
report_name varchar
start_date date
end_date date
file_status smallint
file_url varchar
created_by varchar
created_at timestamptz
file_base64 text
updated_at timestamptz
```

Implikasi: `report.list.cust_id` ada → keputusan user "owner = auth user" tetap konsisten dengan kolom existing.

## Scope rule (Principal `C26002`)

```sql
SELECT cust_id, parent_cust_id, is_del, is_active
FROM smc.m_customer
WHERE parent_cust_id = 'C26002'
ORDER BY cust_id
LIMIT 20;

C26002      |C26002|f|t
C260020001  |C26002|f|t
C260020002  |C26002|f|t
C260020003  |C26002|f|t
```

Implikasi:

- Principal user dengan auth `C26002` (`auth == parent`) memang punya child `C260020001`, `C260020002`, `C260020003`.
- Distributor user `C260020001` (`auth != parent`) tidak punya child di kolom `parent_cust_id`. Konsisten dengan rule plan: distributor hanya boleh request `cust_id`-nya sendiri.
- `ExistsCustomerInParentScope` di repo (filter `is_del=false AND is_active=true`) cocok dengan data: semua 4 baris `is_del=f, is_active=t`.

## Trend Sales SQL (cust child `C260020001`, year 2026)

```sql
SELECT m.month,
       COALESCE(SUM(fo.gross_sale), 0)::bigint                          AS total_gross_sale,
       COALESCE(SUM(fo.discount + fo.special_discount), 0)::bigint      AS total_discount_promo,
       COALESCE(SUM(fo.net_sales_exclude_ppn), 0)::bigint               AS net_sales
FROM (SELECT generate_series(1,12) AS month) m
LEFT JOIN report.dim_dates dt
       ON dt.month = m.month AND dt.year = 2026
LEFT JOIN report.fact_orders fo
       ON fo.date_id = dt.id AND fo.cust_id = 'C260020001'
GROUP BY m.month
ORDER BY m.month;
```

Output:

```text
1 |0
2 |0
3 |0
4 |469500000|2132000|477368000
5 |0
...
12|0
```

Verifikasi: hasil 12 bulan dengan zero-fill, dan baris bulan 4 sama dengan response example pada `docs/Secondary Sales Report_BE.md` (`total_gross_sale=469500000`, `total_discount_promo=2132000`, `net_sales=477368000`). Logika repo `SecondarySalesReportTrendSales` selaras dengan SQL ini.

## fact_orders coverage per cust (filter `C26002%`)

```sql
SELECT cust_id, COUNT(*) AS rows,
       SUM(gross_sale)::bigint, SUM(net_sales_exclude_ppn)::bigint
FROM report.fact_orders
WHERE cust_id LIKE 'C26002%'
GROUP BY cust_id
ORDER BY cust_id;

C260020001|37|469500000|477368000
```

Implikasi:

- Hanya `C260020001` yang punya data di `fact_orders` saat ini.
- Saat principal `C26002` mengirim `cust_id=C260020001` lewat Trend Sales/SUM Date, response harus mengandung angka di atas.
- Saat principal `C26002` mengirim `cust_id=C260020002` atau `C260020003`, response harus 12 baris zero-fill (cust valid scope, tetapi belum ada data).
- Saat distributor `C260020001` tidak mengirim `cust_id`, fallback ke auth → harus dapat data sama dengan principal+`C260020001`.
- Saat distributor `C260020001` mengirim `cust_id=C260020002` → harus 403, scope check di service tolak request.

## Catatan keamanan

- Password DB tidak ditulis ke file ini, ke command history checked-in, atau ke test fixture. Gunakan `PGPASSWORD` env / secret manager saat menjalankan ulang query.
- Semua command di sesi ini read-only `SELECT`. Tidak ada `INSERT/UPDATE/DELETE/DDL` dieksekusi.
- DB ini adalah `scylla_citus_dev`, bukan staging/prod. Jangan ulang validasi langsung ke staging/prod tanpa otorisasi terpisah.
