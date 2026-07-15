# Discovery Evidence — SX-1917 Payment Deposit Report List

Task ID: `20260506-1345-sx-1917-payment-deposit-report-list`
Waktu lokal: 2026-05-06 13:45 Asia/Jakarta

## Pemeriksaan Awal

- Perintah wajib dari `AGENTS.md` dijalankan dari root repo:
  - `rtk docker compose -f docker-compose.yml ps`
- Hasil: Docker daemon tidak aktif: `Cannot connect to the Docker daemon at unix:///Users/ujang/.docker/run/docker.sock. Is the docker daemon running?`
- Implikasi: rencana validasi manual dengan service/DB lokal perlu mengaktifkan Docker terlebih dahulu atau menjalankan test yang tidak membutuhkan service Docker.

## File yang Diinspeksi

- `finance/main.go`
- `finance/controller/payment_deposit_report_controller.go`
- `finance/controller/payment_deposit_report_controller_test.go`
- `finance/controller/report_payment_deposit_controller.go`
- `finance/entity/payment_deposit_report.go`
- `finance/entity/report_payment_deposit.go`
- `finance/model/payment_deposit_report.go`
- `finance/repository/payment_deposit_report_repository.go`
- `finance/repository/payment_deposit_report_repository_test.go`
- `finance/repository/report_payment_deposit_repository.go`
- `finance/service/payment_deposit_report_service.go`
- `finance/service/report_payment_deposit_service.go` ditemukan sebagai implementasi lama/alternatif untuk endpoint `/v1/reports/payment-deposit`.

## Pola Project yang Ditemukan

- Module target adalah `finance/` dengan `go.mod` sendiri; tidak ada root module.
- Endpoint baru/spec sudah sebagian ada pada controller baru:
  - `PaymentDepositReportController.Route`: `app.Group("/finance/v1/reports/payment-deposit", middleware.JWTProtected())`
  - `GET ""` untuk list dan `GET "/download"` untuk download.
- Ada controller lama:
  - `ReportPaymentDepositController.Route`: `app.Group("/v1/reports", middleware.JWTProtected())`, endpoint `/payment-deposit` dan `/payment-deposit/download`.
  - `main.go` masih mendaftarkan keduanya.
- Response memakai `responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)` lalu `Setmsg`, `Setdata`.
- Auth tenant memakai locals Fiber:
  - `cust_id`
  - `parent_cust_id`
- Pagination khusus response report sudah ada di entity:
  - `PaymentDepositReportResponse{Items, Summary, Pagination}`
  - `PaymentDepositReportPagination{Page, Limit, TotalData, TotalPage}`
- Pattern test saat ini banyak memakai GORM `DryRun` untuk memeriksa SQL, bukan full integration DB.

## Reuse Candidates

- Reuse `PaymentDepositReportController`, `PaymentDepositReportService`, `PaymentDepositReportRepository`, `PaymentDepositReportQueryFilter`, `PaymentDepositReportResponse` sebagai target utama karena route sudah sesuai `/finance/v1/reports/payment-deposit`.
- Reuse helper controller yang sudah ada:
  - `normalizeDateInput` mendukung epoch dan `YYYY-MM-DD` via `str.UnixTimestampToUtcTime`.
  - `normalizeDepositNoQuery` untuk CSV `deposit_no`.
  - `validateAndNormalizeSort` perlu diperluas whitelist field.
- Reuse repository tests `payment_deposit_report_repository_test.go` sebagai basis Red tests untuk query generation.
- Reuse service response mapping di `payment_deposit_report_service.go`, tetapi perlu menambah `deposit_type` dan collector nullable untuk AP.

## Gap Implementasi Saat Ini

- Query repository saat ini hanya AR dari `acf.deposit`; belum ada AP dan belum ada `UNION ALL`.
- `PaymentDepositReportQueryFilter` belum punya `deposit_type` dan `emp_id` array; masih memakai single `salesman_id`.
- Controller masih mewajibkan `salesman_id`; spec SX-1917 mewajibkan `deposit_type` dan membuat `emp_id` conditional/opsional untuk AR.
- Repository `buildQuery` langsung join `acf.deposit_payment` dan `acf.deposit_expense`, sehingga berisiko row multiplication saat ada banyak payment dan banyak expense. Spec meminta pre-aggregation payment/expense.
- `expense_amount` saat ini `SUM(COALESCE(de.payment_amount, 0))` dari join langsung; perlu subquery aggregate `deposit_expense` per `deposit_no, cust_id`.
- `total_payment` saat ini memakai `d.total_payment`; spec meminta formula `cash + cheque + transfer + return + credit/debit - expense`.
- Field response belum punya `deposit_type`.
- Model row memakai `SalesmanID/SalesmanCode/SalesmanName`; spec menyebut `collector_*`. Bisa reuse nama internal atau rename untuk kejelasan, tetapi JSON harus `collector_*`.
- Sort whitelist saat ini hanya `created_date` dan `deposit_date`; spec perlu `deposit_date`, `deposit_no`, `deposit_type`, `collector_name`, `total_payment`.
- Default sort saat ini `created_date:desc`/`d.created_at DESC`; untuk union final perlu sort dari alias `t.deposit_date` dan default `deposit_date DESC`.
- AP soft delete wajib `app.deleted_by IS NULL`; AR perlu cek apakah kolom soft-delete sebenarnya `deleted_at` atau `deleted_by`. Implementasi saat ini menggunakan `d.deleted_at IS NULL`.
- AP join reference tidak menyertakan `cust_id` pada `account_payable_payment_options`; perlu dicek schema/model jika tersedia, tetapi minimal join by payment number sesuai spec dan filter parent `app.cust_id`.

## Constraints

- Semua query wajib filter `cust_id` dari token, bukan query param atau hardcoded sample.
- Dynamic SQL hanya boleh untuk fragment yang dikendalikan whitelist; value harus binding.
- Pagination dan count harus diterapkan setelah union final.
- `emp_id` hanya untuk AR dan tidak boleh menghilangkan AP ketika `deposit_type=AR,AP`.
- AP collector fields harus `NULL` sesuai table spec.
- Tidak boleh menambah debug prints `fmt.Println` / `log.Println` di feature code.

## Risiko

- Ada dua implementasi Payment Deposit Report (`payment_deposit_report_*` dan `report_payment_deposit_*`) yang bisa membingungkan route dan download behavior.
- Perubahan entity/controller dari `salesman_id` ke `emp_id` bisa mempengaruhi endpoint download jika masih memakai service yang sama.
- Jika download endpoint masih harus mendukung parameter lama, implementasi perlu backward compatibility atau pembaruan paralel.
- Test DryRun hanya memvalidasi SQL shape, bukan hasil aggregation aktual; perlu rencana integration/manual SQL validation bila DB tersedia.
- Docker daemon tidak aktif pada discovery sehingga service/DB local verification belum bisa dijalankan.

## Commands/Docs Checked

- `rtk docker compose -f docker-compose.yml ps` — gagal karena Docker daemon tidak aktif.
- Local discovery memakai `Glob`, `Grep`, dan `Read` pada module `finance`.
- Official docs/context7 tidak diperlukan untuk rencana ini karena perubahan memakai GORM/Go/Fiber pola lokal yang sudah ada dan tidak bergantung behavior library baru.
- GitHub/web/browser tidak diperlukan karena tidak ada upstream/reference eksternal yang mempengaruhi implementasi backend ini.
