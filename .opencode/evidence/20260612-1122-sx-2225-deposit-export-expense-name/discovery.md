# Discovery — SX-2225 Expense Name kosong di export Payment Deposit Report (mobile)

Task ID: `20260612-1122-sx-2225-deposit-export-expense-name`
Tanggal discovery: 2026-06-12 (Asia/Jakarta)

## Ringkasan masalah
Expense yang dibuat dari aplikasi mobile tidak menampilkan Expense Type/Expense name di file Excel export Payment Deposit Report. Expense yang dibuat dari web tampil benar. Keduanya memakai query export yang sama, jadi perbedaan ada pada **data yang direferensikan**, bukan dua jalur export berbeda.

## File yang diinspeksi
- `finance/repository/payment_deposit_report_repository.go` (query export AR + AP, `buildDownloadARQuery`).
- `finance/service/payment_deposit_report_service.go` (`generateExcel`, mapping kolom Excel).
- `finance/model/payment_deposit_report.go` (`PaymentDepositReportDownloadRow.ExpenseName *string`).
- `finance/model/expense_type.go` (`acf.expense_type` punya kolom `cust_id`).
- `finance/repository/expense_repository.go` (web expense-type lookup di-scope `parent_cust_id`).
- `mobile/service/expense.go` (`Create`, jalur create expense dari mobile).
- `mobile/model/expense.go` (model `Expense`, `ExpenseType`; komentar "Global table, tidak ada cust_id").
- `mobile/repository/expense_repository.go` (`FindAllExpenseTypeLookup`, `FindExpenseTypeById` — tidak ada filter `cust_id`).

## Temuan kunci (code-confirmed)

### 1. Export AR-expense join di-scope ke `parentCustId`
`finance/repository/payment_deposit_report_repository.go` baris ~403-419:
```text
COALESCE(etr.expense_type_name, '') AS expense_name
...
LEFT JOIN acf.expense ex ON ex.expense_id = de.expense_id AND ex.cust_id = d.cust_id AND ex.deleted_at IS NULL
LEFT JOIN acf.expense_type etr ON etr.expense_type_id = ex.expense_type_id AND etr.cust_id = ?
```
Argumen `?` untuk `etr.cust_id` di-bind dengan `parentCustId` (baris ~410: `args = append(args, parentCustId, custId, startDate, endDate)`).

Karena expense row sendiri tetap muncul (amount tampil, hanya nama kosong), join `de` dan `ex` sukses. Yang gagal hanya join `etr` → kondisi `etr.cust_id = parentCustId` tidak match → `COALESCE(..., '')` mengembalikan string kosong.

### 2. Mobile memilih expense_type tanpa scope cust_id; web di-scope parent_cust_id
- Web: `finance/repository/expense_repository.go` `buildExpenseTypeListBaseQuery` memfilter `acf.expense_type.cust_id = parentCustId`; `FindByCodeAndName`/`FindById` juga selalu filter `cust_id`. Web create (`finance/service/expense_service.go`) menyimpan `cust_id`. → Expense type yang dipilih web selalu ber-`cust_id = parentCustId`, sehingga join export match.
- Mobile: `mobile/repository/expense_repository.go` `FindAllExpenseTypeLookup` dan `FindExpenseTypeById` hanya filter `is_del = false AND is_active = true` (TIDAK ada `cust_id`). Komentar model menyatakan tabel dianggap global. → Mobile bisa memilih `expense_type_id` yang `cust_id`-nya bukan `parentCustId` (mis. `custId` distributor atau tenant lain), sehingga join export yang di-scope `parentCustId` tidak match → nama kosong.

### 3. Tidak ada kolom snapshot nama di acf.expense
`mobile/model/expense.go` dan `finance/model/expense.go`: `acf.expense` menyimpan `expense_type_id` (FK) saja, tidak ada kolom `expense_type_name` denormalized. Jadi ini **bug read-side** (resolusi nama via join), bukan write-side. Backfill kemungkinan TIDAK diperlukan jika fix read-side berhasil meresolusi via FK.

### 4. `expense_type_id` adalah PK auto-increment global
`acf.expense_type.expense_type_id` adalah `primaryKey;autoIncrement` (unik global lintas tenant). Maka join hanya pada `expense_type_id` bersifat deterministik; kondisi `etr.cust_id = parentCustId` adalah penyebab kebocoran kosong.

## Hal yang BELUM bisa dikonfirmasi (butuh DB)
- Nilai `cust_id` aktual baris `acf.expense_type` untuk expense `E20260611001` vs nilai `parentCustId`/`custId` yang dipakai report. (Compose tidak running saat discovery; tidak mengakses credentials.)
- Format tampilan kanonik yang benar: query saat ini HANYA select `expense_type_name`, BUKAN `{code} - {name}`. Tiket menyebut ekspektasi `000 - Uang Parkir` (`code - name`). Perlu cek baris web yang "benar" di export untuk memastikan apakah format target adalah `name` saja atau `code - name`.

## Preseden pola di repo
Query yang sama (`buildDownloadARQuery`, join `mst.m_outlet`) sudah memakai pola scoping permisif: `mo.cust_id = dp2.cust_id OR mo.cust_id = ?` (parentCustId). Fix expense_type yang konsisten dapat meniru pola OR ini, atau join hanya pada PK global.

## Test pattern yang tersedia untuk reuse
- `finance/service/payment_deposit_report_service_test.go`: sudah ada test `generateExcel` (`TestPaymentDepositReportService_GenerateExcel*`) yang memetakan `PaymentDepositReportDownloadRow` → cell Excel. Bisa ditambah kasus `ExpenseName` populated dari mobile fixture.
- `finance/repository/expense_repository_test.go`: pola dry-run GORM (`newExpenseTypeDryRunDB`, assert substring SQL). Untuk raw SQL builder `buildDownloadARQuery`, test bisa memanggil method dan assert string SQL pada kondisi join expense_type.

## Validasi yang relevan
- Dari direktori `finance`: `rtk go mod download && rtk go mod tidy`, `rtk go test ./...`, dan targeted `rtk go test ./service -run TestPaymentDepositReport` / `rtk go test ./repository -run TestPaymentDepositReport`.
