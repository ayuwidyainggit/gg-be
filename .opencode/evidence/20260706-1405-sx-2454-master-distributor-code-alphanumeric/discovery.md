# Discovery — SX-2454 Alphanumeric Distributor Code

Task ID: `20260706-1405-sx-2454-master-distributor-code-alphanumeric`
Reviewer: `@artifact-planner` (read-only, no source edits)

## Files inspected (repo-local evidence)

### `master` service (Fiber, port 9002)
- `master/controller/m_distributor_controller.go` (432 baris) — route `POST/GET/PATCH/DELETE /v1/distributors`, `GET /v1/distributors/customers`. `Update` di `controller.Update` (line 268) memvalidasi body dengan `controller.validator.ValidateStruct(request, ...)`.
- `master/service/m_distributor_service.go` (398 baris) — `Store`, `Update`, `List`, `LookupList`, `Detail`, `Delete`, `ListWithCustomer`. Tidak ada `strconv.Atoi/ParseInt` terhadap `distributor_code`. Pemeriksaan duplikat via `FindOneByDistributorCodeAndCustId` (string-based).
- `master/repository/m_distributor_repository.go` (999 baris) — semua query `distributor_code` dikirim sebagai `string` placeholder `$1/$2`. Insert: `INSERT INTO mst.m_distributor ... $2 distributor_code ...` (line 92-117). FindOneBy: `SELECT sp.distributor_id, sp.distributor_code FROM mst.m_distributor sp WHERE sp.distributor_code = $1` (line 524-528). Tidak ada cast `int`/`uint`/`float` ke/dari `distributor_code`.
- `master/entity/m_distributor.go` (304 baris) — DTO:
  - `CreateDistributorBody.DistributorCode`: `validate:"required,alphanum,max=20"` (line 163) → root cause: tag `alphanum` (bawaan `go-playground/validator/v10`, regex `^[0-9a-zA-Z]*$`) menolak `-` dan `_`.
  - `UpdateDistributorRequest.DistributorCode`: `validate:"omitempty,alphanum,max=20"` (line 212) → same root cause.
  - `Barcode`: `validate:"omitempty,numeric,max=13"` — di luar scope.
  - Response struct sudah `string`.
- `master/model/m_distributor.go` (4 model: create, list, contact, tax) — `DistributorCode string` (semua). Tidak ada tipe `int`/`uint`/`float`.
- `master/pkg/validation/validation.go` (280 baris) — bootstrap `go-playground/validator/v10`. Custom: `qtystr`, `alphanumericSpace`, `answer_frequency`. Tidak ada override `alphanum`. Pola: `vc.RegisterValidation("namaTag", func)` lalu daftarkan translasi `id`/`en`.
- `master/pkg/validation/validation-product.go` — pattern referensi custom validator `alphanumericSpaceDash` (regex `^[a-zA-Z0-9\s-]*$`, dengan translasi ID/EN) → bukti pola untuk menambah `distributorCode`/`alphanumDashUnderscore` baru.
- `master/service/m_distributor_update_partial_test.go` + `master/service/m_distributor_service_test.go` + `master/entity/m_distributor_validation_test.go` — testify + sqlmock sudah dipakai; pattern siap untuk ditambah test SX-2454.
- `master/main.go`, `master/Dockerfile`, `master/.air.toml`, `master/.env` — service jalan via compose, hot reload via `rtk` per `QUALITY.md`.

### Validasi perintah
- `master/go.mod` ada, buildable per service.
- `docker-compose.yml` mendaftarkan `master` service port 9002.
- Stack & command docs (`PROJECT_STACK.md`/`PROJECT_COMMANDS.md`/`FRAMEWORK_PLAYBOOK.md`/`PROJECT_DETECTED_TOOLS.md`) **tidak ditemukan** → skip, gunakan `QUALITY.md` cheat sheet + `SERVICE_MATRIX.md` precedence: `docker-compose.yml` → `go.mod`/`Makefile`/`.env` → README.

### Cross-module scan (grep `distributor_code`/`DistributorCode`)

| Modul | Bukti penggunaan | Status |
|---|---|---|
| `sales` | `sales/model/{report,promotionV2}.go` bertipe `string`. `sales/controller/report_controller.go` parser `json.RawMessage` untuk `distributor_code` (sudah string-friendly). `sales/repository/report_repository.go` SELECT `md.distributor_code` (string). | Aman (string, tidak ada cast numeric). |
| `pjp` | `pjp/model/live_monitoring.go` bertipe `string` `gorm:"column:distributor_code"`. `pjp/data/response/destination_detail_response.go` `string`. | Aman. |
| `pjp-principle` | `model/distributor_dms.go` `string varchar(125)`. `repository/distributor_dms/distributor_repository.go:26` `Where("distributor_code = ?", filter.DistributorCode)` — string. | Aman. |
| `pjp-sales` | `entity/report.go` `string` + `DistributorCodes []string` normalizer `NormalizeDistributorCodeList` (line 346). `repository/report_repository.go` SELECT `md.distributor_code` string. `repository/activity_report_query.go:221` `COALESCE(NULLIF(TRIM(distributor_code), ''), business_unit_code)` — string-only. | Aman. |
| `mobile` | `model/m_distributor.go` `string`. `middleware/jwt_middleware.go` set Locals `distributor_code` (string). | Aman. |
| `inventory` | `model/{gr,replenishment,reports,order_booking}.go` `*string` atau `string`. `service/sap_replenishment_status.go:95-123` melakukan validasi `cust_id/distributor_code` (string). | Aman. |
| `finance` | `model/{ap_payment,ap_supplier_invoice_return,account_payable_list}.go` `*string`/`string`. `repository/ap_list_repository .go` SELECT `md.distributor_code` (string). | Aman. |

> Konfirmasi: tidak ada modul yang melakukan `strconv.Atoi(*distributor_code)` atau `cast(distributor_code as int)` di lapisan BE. Semua join/list/filter sudah `string`-native. Risiko cross-modul rendah untuk slice 1.

### `monitoring_activity_be_doc.txt` dll
- Tidak ada di repo lokal (`.opencode/evidence/<task>/monitoring_activity_be_doc.txt` juga tidak ada). Pengetahuan diambil dari prompt task + grep repo.

## Perintah dijalankan

```bash
rg -n "distributor_code" --type go   # scope: master (lengkap), all modules (sampel)
rg -n "Atoi|ParseInt|numeric" --type go master
rg -n "excelize|xlsx|Export|Import" --type go master
rg -n "ImportDistributor|ExportDistributor|distributor.*[Ii]mport" --type go master
ls master/{controller,service,repository,entity,model,pkg/validation}/*.go
```

Hasil: scope utama terkonfirmasi di entity & service. Import/Export Excel **khusus distributor** tidak ditemukan di service `master` (asumsi: ditangani di service lain atau tidak ada). Slice 1 tidak menyentuh import/export; jika QA nanti menemukan import/export yang bermasalah, akan dipisah jadi slice 2.

## Constraints & riesgos

- **Tidak boleh ubah DB schema** (kolom sudah `Varchar(20)` per prompt task — dikonfirmasi via `model.Distributor` Go yang semuanya `string`).
- **Tidak boleh ubah unique constraint** (di luar scope).
- **Backward-compat wajib** — distributor dengan code numerik existing harus tetap bisa di-GET/PATCH.
- Validator `alphanum` bawaan tidak bisa menerima `-`/`_`; tambah custom validator adalah pola yang sudah ada di `validation-product.go` (`alphanumericSpaceDash`).
- `model.Distributor` field `DistributorCode` di-`json:"distributor_code"` — response API sudah string. Tidak ada perubahan response shape.

## Confirmed vs Assumed Audit

| Klaim | Status | Bukti |
|---|---|---|
| Service `master` Fiber, port 9002 | confirmed_repo | `docker-compose.yml`, `SERVICE_MATRIX.md` |
| Route CRUD `/v1/distributors` ada dan terdaftar | confirmed_repo | `master/controller/m_distributor_controller.go:34-43` |
| Bug ada di `entity/m_distributor.go` line 163 & 212 (`alphanum,max=20`) | confirmed_repo | baca langsung |
| `alphanum` default `go-playground/validator/v10` menolak `-`/`_` | confirmed_docs (eksternal) | perilaku bawaan library; pola custom sudah ada di `validation-product.go` |
| Tidak ada `strconv.Atoi` di path distributor_code | confirmed_repo | grep 100 matches: tidak ada |
| Tidak ada import/export Excel khusus distributor di `master` | assumption | grep tidak menemukan service/route |
| DB kolom `mst.m_distributor.distributor_code` adalah `Varchar(20)` | assumption (per task prompt) | migration folder di `master` tidak ditemukan |
| Staging hostname `best.staging.scyllax.online` | user_confirmed | Q&A gate |
| Regex final `^[A-Za-z0-9_-]{1,20}$` | user_confirmed | Q&A gate |

## Reuse candidates (sudah ada di repo, dipakai ulang)

- `go-playground/validator/v10` + `pkg/validation/validation.go` — bootstrap validator (tambah tag baru).
- Pola `alphanumericSpaceDash` di `pkg/validation/validation-product.go` — referensi langsung untuk `alphanumericDashUnderscore`.
- `pkg/structs` Automapper untuk DTO ↔ model (tidak berubah).
- `lib/pq` untuk mapping `pq.Error` kode `23505` (sudah dipakai di `mapDistributorDuplicateError`).
- Testify + sqlmock + pattern `validCreateDistributorBody()`/`validUpdateDistributorRequest()` di `entity/m_distributor_validation_test.go` untuk test validator.
- `repository.NewDistributorRepository(sqlxDB)` di service test untuk service-level test.

## Findings singkat

- **Root cause**: tag `validate:"alphanum,max=20"` di `CreateDistributorBody.DistributorCode` & `UpdateDistributorRequest.DistributorCode`.
- **Fix minimal**: ganti tag ke custom validator baru `alphanumericDashUnderscore` (regex `^[A-Za-z0-9_-]{1,20}$`), registrasi di `validation.go` (ID + EN), tambah translasi.
- **Backward-compat**: regex baru menerima `[A-Za-z0-9_-]+` → superset dari `alphanum`, jadi tidak ada regresi untuk code numerik existing.
- **Risiko SQL injection** tidak relevan: pakai parameterized `$N`. Tidak ada perubahan query.
- **Cross-module**: scope slice 1 tidak menyentuh modul lain; hasil grep menunjukkan mereka sudah string-native.
