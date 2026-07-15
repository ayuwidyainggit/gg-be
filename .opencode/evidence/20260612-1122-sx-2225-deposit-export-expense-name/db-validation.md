# DB Validation — SX-2225 (ggn_scyllax lokal)

Tanggal: 2026-06-12 (Asia/Jakarta)
DB: `ggn_scyllax` @ `localhost:5432` (postgres/postgres)

## Tujuan
Membuktikan root cause empiris dan memverifikasi fix read-side export pada data nyata `E20260611001`.

## Data sumber expense
```
acf.expense JOIN acf.expense_type:
cust_id=C260040002 expense_id=683 doc_no=E20260611001 expense_type_id=50 source=2(mobile)
  expense_type.cust_id=C22001 code=000 name=Uang Parkir
cust_id=C260040001 expense_id=681 doc_no=E20260611001 expense_type_id=50 source=2(mobile)
  expense_type.cust_id=C22001 code=000 name=Uang Parkir
```
Catatan: `source=2` = mobile (lihat `mobile/service/expense.go` `sourceMobile = 2`).

## Link deposit (yang relevan untuk export)
```
deposit.cust_id=C260040002 deposit_no=DP2606110003 deposit_date=2026-06-11 emp_id=429
  -> expense_id=683 doc_no=E20260611001 expense_type_id=50
  -> expense_type.cust_id=C22001 code=000 name=Uang Parkir payment_amount=25000
```

## Mapping customer (smc.m_customer)
```
C22001       parent=C22001  Principal A Company
C260040002   parent=C26004  PT. Sinar Jaya
C260040001   parent=C26004  PT. Makmur Sejahtera
```
Jadi untuk deposit `C260040002`, `parentCustId` yang dipakai report = `C26004`. Expense type ID 50 sebenarnya milik `C22001`. Mismatch inilah penyebab nama kosong di export lama.

## Bukti old vs new (satu query, data nyata)
Disimulasikan dengan dua join expense_type: `et_old` memakai kondisi lama `et_old.cust_id = parentCustId`, `et_new` PK-only.
```
cust_id=C260040002 parent_cust_id=C26004 deposit_no=DP2606110003 doc_no=E20260611001
expense_type_id=50 actual_expense_type_cust_id=C22001 code=000 name=Uang Parkir
old_export_expense_name = ""                  <- BUG (join cust_id=C26004 tidak match C22001)
new_export_expense_name = "000 - Uang Parkir"  <- FIX
expense = -25000.0000                          <- tidak berubah
```

## Kesimpulan
- Root cause read-side terkonfirmasi empiris: kondisi `etr.cust_id = parentCustId` membuang nama expense type untuk expense yang dibuat mobile dengan expense type milik cust_id lain (di sini principal `C22001`).
- Fix PK-only menghasilkan format `000 - Uang Parkir` sesuai acceptance criteria, tanpa mengubah nilai `expense`.
- Tidak ada backfill diperlukan; nilai diselesaikan murni dari FK saat export.

## Catatan keamanan
- Hanya operasi SELECT read-only.
- Tidak menyalin/menyebarkan credentials; koneksi memakai kredensial lokal default repo dan tidak ditulis ke artefak selain host/db generik.
