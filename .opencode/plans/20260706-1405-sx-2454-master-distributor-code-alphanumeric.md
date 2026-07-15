# Plan — SX-2454: Alphanumeric Distributor Code (Master)

> Source of truth: file ini. Setelah disintesis, draft/evidence dapat dibersihkan kecuali yang masih operasional.
> Task ID: `20260706-1405-sx-2454-master-distributor-code-alphanumeric`
> Parent story: SX-799
> Target service: `master` (Fiber, port 9002, `mst.m_distributor.distributor_code`)
> Plan quality gate: `PASS_FOR_SLICE` (slice 1 di `master`; slice 2 parked)
> Active lane: `@artifact-planner` (read-only). Eksekusi ke lane berikutnya (`@orchestrator`/implementasi) — lihat §Active-lane reset.

## Goal

Menghapus semua constraint numerik terhadap `mst.m_distributor.distributor_code` di service `master` agar field tersebut dapat menyimpan nilai alfanumerik dengan karakter `-` dan `_` (panjang 1..20 char), tanpa mengubah DB schema, unique constraint, atau response shape, dan tanpa meregresi data existing yang numerik.

## Non-goals

- Migrasi data existing numerik → alfanumerik.
- Perubahan `Varchar(20)` ke panjang lain.
- Audit khusus perubahan `distributor_code` (pakai audit existing).
- Endpoint lookup khusus `GET /v1/distributors/by-code`.
- Modifikasi aturan validator `alphanum`/`numeric` untuk field lain (`barcode`, `zip_code`, dst.).
- Refactor handler/response/service yang tidak menyentuh validasi `distributor_code`.

## Scope

**Slice 1 (eksekusi):** service `master`.

- DTO `CreateDistributorBody.DistributorCode` (`master/entity/m_distributor.go:163`).
- DTO `UpdateDistributorRequest.DistributorCode` (`master/entity/m_distributor.go:212`).
- Validator custom baru di `master/pkg/validation/validation.go` (registrasi `alphanumDashUnderscore` + translasi `id`/`en`).
- Test table-driven di `master/entity/m_distributor_validation_test.go` (atau file baru) + 1 test service-level untuk PATCH alfanumerik (di `master/service/m_distributor_*_test.go`).

**Slice 2 (parkir — tidak dalam v1):** jika QA/exploration menemukan import/export Excel distributor atau modul lain dengan asumsi numeric, buat tiket turunan.

- Audit `master/service/...import*.go` & `...export*.go` khusus distributor — saat ini grep tidak menemukan handler tersebut di `master`. Konfirmasi lebih lanjut via eksplorasi runtime.
- Cross-module string-only verification (`pjp-sales`, `sales`, `inventory`, `finance`, `mobile`, `pjp-principle`) — hasil grep static sudah `string` semua; tetap diparkir untuk validasi runtime.

## Requirements

1. `distributor_code` menerima `A-Z`, `a-z`, `0-9`, `-`, `_`, panjang 1..20 char (regex `^[A-Za-z0-9_-]{1,20}$`).
2. `POST /v1/distributors` dengan `distributor_code` alfanumerik (contoh: `DIST-NEW1`) → 201 Created.
3. `PATCH /v1/distributors/:id` dengan `distributor_code` alfanumerik (contoh: `DIST128-AB`) → 200 OK, response `distributor_code` identik string.
4. `GET /v1/distributors/:id` → response `distributor_code` string utuh (termasuk leading zero & karakter `-`/`_`).
5. Distributor existing dengan code numerik (`162612`, dst.) tetap bisa di-GET dan di-PATCH tanpa error.
6. PATCH dengan `distributor_code` kosong → 400 dengan pesan kesalahan.
7. PATCH dengan `distributor_code` sepanjang 21 char → 400 dengan pesan kesalahan.
8. PATCH dengan `distributor_code` mengandung spasi, titik, atau karakter non-izinkan → 400.
9. PATCH dengan `distributor_code` duplikat (sudah dipakai distributor lain pada `cust_id` yang sama) → 400 dengan pesan "Distributor code already exists. Please use a different distributor code.".
10. Duplicate key DB (`pq.Error` kode `23505`) saat INSERT/UPDATE → tetap dipetakan ke pesan terstandar (perilaku `mapDistributorDuplicateError` yang sudah ada).
11. Tidak ada perubahan DB schema/column type/index/unique.
12. Tidak ada perubahan response shape API (field `distributor_code` tetap `string`).
13. Pesan kesalahan validator untuk `id` & `en` disesuaikan.

## Acceptance Criteria

1. `grep -RIn "distributor_code" --include="*.go" master | rg -iE "validate:\".*alphanum"` tidak lagi mengembalikan match untuk field `DistributorCode` (kecuali field lain di luar scope).
2. `cd master && rtk go test ./...` lulus.
3. `rtk go test ./entity -run "DistributorCode" -v` lulus untuk semua skenario pada TDD/Test Plan §TDD.
4. Verifikasi manual staging: PATCH distributor staging ke `DIST-<id>-ALPHA`, GET round-trip identik, log validasi tidak menolak.
5. Service `master` boot lokal tanpa error.
6. `cd master && rtk go mod tidy && rtk go build ./...` sukses.
7. PR menjelaskan: lokasi validator yang diganti, hasil grep pre/post, before/after curl PATCH, dan ringkasan hasil test.

## Existing Patterns / Reuse

- `master/pkg/validation/validation.go` — bootstrap `go-playground/validator/v10`; pola `vc.RegisterValidation("alphanumericSpace", ...)` (line 56-57) + translasi `id`/`en` (line 84-94, 123-133) + `alphanumericSpaceDash` di `validation-product.go` (line 15-40). Reuse pola ini untuk `alphanumDashUnderscore`.
- `master/service/m_distributor_service.go` — `mapDistributorDuplicateError` untuk `pq.Error` `23505` (line 38-49) — tidak perlu diubah.
- `master/entity/m_distributor_validation_test.go` — test validator existing; tambah test baru di sini atau file baru.
- `master/service/m_distributor_service_test.go` + `m_distributor_update_partial_test.go` — pattern testify + sqlmock + `setupDistributorServiceTest` + `validUpdateRequest()`.
- `pkg/structs` Automapper (tidak berubah).
- `lib/pq` (sudah dipakai, tidak perlu ditambah).

## Source Anatomy

| Subsystem | File:line (master) | Catatan |
|---|---|---|
| Route CRUD | `controller/m_distributor_controller.go:34-43` | Fiber group `/v1/distributors` |
| Handler Update | `controller/m_distributor_controller.go:268-351` | `controller.validator.ValidateStruct(request, ...)` line 311 |
| Handler Create | `controller/m_distributor_controller.go:45-83` | validasi line 67 |
| Validator bootstrap | `pkg/validation/validation.go:51-162` | tempat daftarkan custom tag |
| Pola custom validator | `pkg/validation/validation-product.go:15-40` (`alphanumericSpaceDash`) | referensi langsung |
| Pola translasi ID/EN | `pkg/validation/validation.go:84-94, 123-133` | referensi |
| DTO Create | `entity/m_distributor.go:158-186` | line 163 = root cause Create |
| DTO Update | `entity/m_distributor.go:203-243` | line 212 = root cause Update |
| Service Update | `service/m_distributor_service.go:237-374` | tidak ada cast numeric |
| Service Store | `service/m_distributor_service.go:74-149` | duplikat check `FindOneByDistributorCodeAndCustId` (string) |
| Repository Insert | `repository/m_distributor_repository.go:90-123` | `$2 distributor_code` (string) |
| Repository Update | `repository/m_distributor_repository.go:594-739` | NamedExec + dynamic `sqlSetFields` (string) |
| Repository FindOneByCode | `repository/m_distributor_repository.go:522-535` | `WHERE sp.distributor_code = $1` (string) |
| Repository List filter | `repository/m_distributor_repository.go:189-193` | `ILIKE` substring — kompatibel string |

## Reference Map

| Fitur | Source basis | Reason |
|---|---|---|
| Custom validator dengan translasi | `pkg/validation/validation-product.go` (`alphanumericSpaceDash`) | repo-backed, pola sudah established |
| Test validator existing | `entity/m_distributor_validation_test.go` | repo-backed |
| Test service Update | `service/m_distributor_update_partial_test.go` | repo-backed |
| Regex `^[A-Za-z0-9_-]{1,20}$` | user_confirmed via Q&A gate + docs go-playground/validator | cukup |
| Backward-compat (superset dari `alphanum`) | first-principles (regex superset) | cukup |

## Constraints

- DB schema **tidak boleh** berubah.
- Unique constraint **tidak boleh** berubah.
- Response shape **tidak boleh** berubah (field `distributor_code` tetap `string`).
- Panjang maks tetap 20 char.
- Pola error response existing (HTTP code + `responsePayload` + `errors`) **harus** dipertahankan.
- Single service (`master`) untuk slice 1.
- Hot-reload service via `rtk` (`.air.toml`).

## Risks

- **R1 — Regex menolak karakter lain yang dipakai FE** (low): QA ticket secara eksplisit hanya menyebut `DIST-15676761A` & `DST128-AB`. Karakter lain (`.`, spasi, unicode) tetap ditolak sesuai `naming convention` singkat. *Mitigasi*: plan membuka "Add when" jika ada tiket lanjutan.
- **R2 — Backward-compat**: regex baru superset dari `alphanum`, jadi tidak ada regresi. *Mitigasi*: test eksplisit `162612` (numerik existing) lolos validasi.
- **R3 — Cross-module regression**: `sales`/`pjp`/`mobile`/`inventory`/`finance`/`pjp-sales`/`pjp-principle` sudah `string` di semua entity. *Mitigasi*: test service `master` + grep static; validasi runtime di slice 2.
- **R4 — Duplikat race condition** (existing): tidak berubah; perilaku `pq.Error 23505` mapping sudah ada.
- **R5 — Leading zero & zero-padded numeric**: regex menerima `0` di awal. *Mitigasi*: response sudah `string` di seluruh chain; tidak ada Excel numeric-format khusus distributor di `master`. *Residual*: jika ada modul FE yang me-`Number()` JSON, ticket terpisah.
- **R6 — JWT `distributor_code` claim** di `mobile/middleware/jwt_middleware.go:67`: claim diset dari struct upstream — tidak terkait validator SX-2454.

## Decisions / Assumptions

### Decisions
- **D1**: regex `^[A-Za-z0-9_-]{1,20}$` (user-confirmed).
- **D2**: nama tag baru `alphanumDashUnderscore`; didaftarkan di `pkg/validation/validation.go` saja.
- **D3**: slice 1 hanya `master`; slice 2 parked.
- **D4**: tidak ada perubahan response shape, DB schema, atau unique constraint.
- **D5**: pesan ID & EN mengikuti pola `alphanumericSpaceDash`.

### Assumptions (terbuka — track jika berubah)
- **A1** — `mst.m_distributor.distributor_code` adalah `Varchar(20)` (per prompt task). Tipe di `model.Distributor` Go adalah `string`. Tidak ada migration file di repo `master` (asumsi: schema sudah final).
- **A2** — Tidak ada service/route import-export Excel khusus distributor di `master` (grep static negatif). Validasi runtime ditunda ke slice 2.
- **A3** — `monitoring_activity_be_doc.txt` dan prompt-prompt eksternal Jira yang dirujuk user **tidak** ada di repo lokal; pengetahuan dari prompt task + grep lokal sudah cukup untuk slice 1.
- **A4** — Backend diuji via `go test` (testify + sqlmock) dan smoke test manual via Fiber endpoint; tidak ada integration test framework baru.

## Execution Source of Truth

Precedence untuk implementasi (kalau ada konflik, ikut urutan atas):

1. Instruksi eksplisit terbaru dari user (jawaban Q&A sudah termasuk: regex `^[A-Za-z0-9_-]{1,20}$`, staging `best.staging.scyllax.online`, scope `master` slice 1 dulu).
2. Non-negotiable Implementation Invariants (§di bawah).
3. Acceptance Criteria & Done Criteria.
4. Implementation Steps (§di bawah) & Execution-ready Worklist.
5. Rekomendasi/follow-up.

## Non-negotiable Implementation Invariants

1. **DB schema tidak berubah.** Tidak ada ALTER TABLE, migrasi, atau perubahan tipe kolom.
2. **Unique constraint DB tidak berubah.**
3. **Response shape tidak berubah.** Field `distributor_code` tetap `string` JSON.
4. **Backward-compat wajib.** Distributor existing numerik (`162612`, dst.) harus lulus validasi baru.
5. **Panjang maks 20 char** sesuai `Varchar(20)`.
6. **Whitespace handling**: payload `distributor_code: "  DIST001  "` tidak boleh lolos jadi beda key; trim sebelum validasi (default behavior JSON binding di Fiber sudah tidak auto-trim, jadi tambahkan `strings.TrimSpace` di handler Update, atau andalkan service `resolveUpdateCustID` upstream). *Implementasi*: gunakan `strings.TrimSpace` di helper validator (di `validation.go`) agar trimming konsisten di Create & Update.
7. **Pesan error response** ID/EN disesuaikan; pola payload (`responsePayload.Seterrors(errs)`) tidak berubah.
8. **Single PR** untuk slice 1, single commit (atau multi-commit logis dengan pesan jelas); tidak menggabungkan refactor lain.

## Do Not / Reject If

- **Reject** jika ada perubahan DB schema (`ALTER TABLE mst.m_distributor ...`).
- **Reject** jika ada perubahan `unique index/constraint` distributor_code.
- **Reject** jika response `distributor_code` di-cast ke `int`/`uint`/`float`/`json.Number`.
- **Reject** jika regex lebih longgar dari `^[A-Za-z0-9_-]{1,20}$` (mis. mengizinkan spasi/titik).
- **Reject** jika validasi `alphanum` dihapus total (tanpa custom validator) — risiko input liar.
- **Reject** jika PR mencampur refactor tidak terkait.
- **Reject** jika ada perubahan yang membuat `162612` (numerik existing) gagal validasi.
- **Reject** jika tidak ada test baru untuk regex.

## Diff Boundary

**Allowed**:
- `master/entity/m_distributor.go` — hanya 2 baris: ganti `validate:"required,alphanum,max=20"` dan `validate:"omitempty,alphanum,max=20"` untuk `DistributorCode`.
- `master/pkg/validation/validation.go` — tambah `vc.RegisterValidation("alphanumDashUnderscore", ...)` + translasi `id`/`en`. Tidak mengubah tag lain.
- `master/entity/m_distributor_validation_test.go` (atau `*_test.go` baru di package `entity`) — tambah test table-driven.
- `master/service/m_distributor_*_test.go` — tambah 1 test PATCH alfanumerik.

**Out of boundary** (harus di-revert atau justifikasi eksplisit):
- File di luar direktori `master/`.
- File `master/go.mod` / `go.sum` (kecuali dep baru yang benar-benar dibutuhkan — diharapkan 0).
- `master/docs/*` (kecuali bila ada OpenAPI yang harus diupdate; di repo ini tidak ada swagger yang eksplisit, jadi tidak wajib).
- Modul lain: `sales/`, `pjp*/`, `mobile/`, `inventory/`, `finance/`, `tms/`, `system/`, `cronjob/`.

**Generated-report exception**: artefak di `.opencode/evidence/<task-id>/` (log curl, log test) boleh ditambah selama eksekusi.

## TDD / Test Plan

**Pendekatan**: Red → Green → Refactor.

**TDD required**: ya, karena menyentuh validator (logika non-trivial).

**Existing test patterns**:
- `master/entity/m_distributor_validation_test.go` — pattern `validator := validation.NewValidator()` + `validator.Validator.Struct(payload)` + `t.Run` table-driven.
- `master/service/m_distributor_update_partial_test.go` — pattern `setupDistributorServiceTest` + sqlmock.

**First failing tests (Red)** — tambahkan di file test baru atau existing:
1. `TestCreateDistributorBodyValidation_DistributorCodeAcceptsAlphanumeric` — payload `DIST-NEW1` → saat ini gagal (`-` ditolak), setelah fix harus pass.
2. `TestUpdateDistributorRequestValidation_DistributorCodeAcceptsAlphanumeric` — payload `DIST128-AB` → saat ini gagal, setelah fix harus pass.
3. `TestUpdateDistributorRequestValidation_DistributorCodeRejectsSpace` — payload `DIST 001` → tetap gagal.
4. `TestUpdateDistributorRequestValidation_DistributorCodeRejectsDot` — payload `DIST.001` → tetap gagal.
5. `TestUpdateDistributorRequestValidation_DistributorCodeRejectsLength21` → tetap gagal.
6. `TestUpdateDistributorRequestValidation_DistributorCodeEmptyRejected` (Create only karena required) → gagal.
7. `TestUpdateDistributorRequestValidation_DistributorCodeBackwardCompatNumeric` — payload `162612` → pass (regresi test).

**Green**: tambahkan `alphanumDashUnderscore` di `pkg/validation/validation.go`, ganti tag di entity.

**Refactor**: bersihkan duplikasi pesan translasi jika ada; konsistenkan error key.

**Edge cases** (sudah termasuk di Red):
- whitespace leading/trailing → trim dulu (pakai `strings.TrimSpace` di validator sebelum regex).
- Unicode karakter (emoji, huruf non-latin) → ditolak regex.

**Test commands** (jalankan dari direktori `master`):

```bash
rtk go mod download
rtk go mod tidy
rtk go test ./entity -run "Distributor" -v
rtk go test ./service -run "Distributor" -v
rtk go test ./...
```

## Implementation Steps

1. **A0** — Baca `master/entity/m_distributor.go` (sudah selesai sebagai discovery), verifikasi baris 163 & 212.
2. **A1** — Di `master/pkg/validation/validation.go`, tambah:
   - `vc.RegisterValidation("alphanumDashUnderscore", alphanumDashUnderscore)` di dalam `NewValidator()` (setelah `alphanumericSpace` line 57).
   - Translasi `id` di `case id.New().Locale():` (line 67-105) — daftarkan key `alphanumDashUnderscoreID`.
   - Translasi `en` di `case en.New().Locale():` (line 106-145) — daftarkan key `alphanumDashUnderscoreEN`.
   - Fungsi `func alphanumDashUnderscore(fl validator.FieldLevel) bool { s := strings.TrimSpace(fl.Field().String()); return regexp.MustCompile(`^[A-Za-z0-9_-]{1,20}$`).MatchString(s) }` di akhir file (sebelum `qtystr` atau setelah `alphanumericSpace` line 272-276, di-trim untuk konsistensi).
3. **A2** — Di `master/entity/m_distributor.go`:
   - Line 163: `validate:"required,alphanum,max=20"` → `validate:"required,alphanumDashUnderscore,max=20"`.
   - Line 212: `validate:"omitempty,alphanum,max=20"` → `validate:"omitempty,alphanumDashUnderscore,max=20"`.
4. **A3** — Tambah test di `master/entity/m_distributor_validation_test.go` (atau file baru `master/entity/m_distributor_code_validation_test.go` jika ingin pemisahan jelas):
   - Tujuh sub-test sesuai TDD plan.
   - Tetap import `master/pkg/validation` seperti test existing.
5. **A4** — Tambah 1 test service-level di `master/service/m_distributor_service_test.go` (atau file baru): `TestDistributorService_Update_AcceptsAlphanumericCode` — mock `FindOneByDistributorCodeAndCustId` return sql.ErrNoRows, `UPDATE mst.m_distributor` return rows affected 1, ExpectCommit. Payload `DistributorCode: serviceStrPtr("DIST-15676761A")`. Validator tidak diuji di service test (sudah di entity test), cukup bahwa service menerima string alfanumerik.
6. **A5** — Verifikasi `rtk go build ./...` dari `master/`.
7. **A6** — Verifikasi `rtk go test ./...` dari `master/`.
8. **A7** — Boot service `master` lokal via `rtk docker compose -f docker-compose.yml up -d master` (atau via `rtk go run main.go` di `master/`).
9. **A8** — Smoke test manual (sesuai §Validation Commands): curl PATCH ke distributor existing dengan payload alfanumerik.
10. **A9** — Verifikasi manual staging `best.staging.scyllax.online`:
    - Dapatkan distributor dengan code numerik; catat `distributor_id`.
    - PATCH ke `DIST-<id>-ALPHA`.
    - GET detail; response `distributor_code` identik string.
    - Catat timestamp & response.
11. **A10** — Tulis PR ringkas + checklist §Acceptance.

## Expected Files to Change

- `master/pkg/validation/validation.go` — tambah ~25-40 baris (registrasi, translasi ID/EN, fungsi validator).
- `master/entity/m_distributor.go` — 2 baris (ganti tag).
- `master/entity/m_distributor_validation_test.go` (atau file baru) — tambah 7 sub-test.
- `master/service/m_distributor_service_test.go` (atau file baru) — tambah 1 sub-test.

**Total diff**: < 80 baris netto.

## Agent / Tool Routing

| Tugas | Owner |
|---|---|
| Eksekusi perubahan kode Go | `@fixer` (atau `@backend` jika lebih cocok; `@fixer` default) |
| Quality gate | `@quality-gate` (review akhir; per AGENTS.md) |
| Penulisan PR | `@fixer` (atau user jika manual) |
| Verifikasi staging manual | user/QA (tim infra memberi akses) |

## Executor Handoff Prompt (copy-paste-ready)

```
Task ID: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
Plan:    .opencode/plans/20260706-1405-sx-2454-master-distributor-code-alphanumeric.md
Caller:  orchestrator -> fixer
Scope:   Service `master` — ganti validator `distributor_code` dari `alphanum` ke custom `alphanumDashUnderscore` (regex ^[A-Za-z0-9_-]{1,20}$) di Create & Update DTO. Tambah test. Tidak ubah DB / response / unique constraint.

Must preserve:
- DB schema & unique constraint distributor_code (Varchar(20)).
- Response shape (distributor_code string).
- Backward-compat: code numerik existing (mis. 162612) tetap valid.
- Panjang maks 20 char.
- Pola error response (responsePayload + errors).
- Pola service-layer transaction (Controller -> Service -> Repository -> DB).

Do not touch:
- File di luar direktori master/.
- master/go.mod / go.sum (0 dep baru; cukup regex stdlib + validator existing).
- Tag validator field lain (barcode, zip_code, dst.).
- Modul lain (sales, pjp*, mobile, inventory, finance, tms, system, cronjob).
- Migration apa pun.

Validation:
- cd master && rtk go mod tidy && rtk go build ./...
- cd master && rtk go test ./entity -run "DistributorCode" -v
- cd master && rtk go test ./service -run "Distributor" -v
- cd master && rtk go test ./...
- Manual smoke: PATCH /v1/distributors/128 dengan {"distributor_code":"DIST128-AB"} -> 200, GET round-trip identik.

Evidence expected:
- Log rtk go test (all pass).
- Hasil curl PATCH pre-fix (rejected) vs post-fix (200) — simpan di .opencode/evidence/20260706-1405-.../curl-*.log
- Catatan verifikasi staging best.staging.scyllax.online — simpan di evidence/staging-verify.md.

Return:
- Daftar path:line yang berubah.
- Output rtk go test (lulus).
- Diff netto <= 80 baris.
- Status: PASS / FAIL / NEEDS_DECISION.
```

## Execution-ready Worklist / Handoff Contract

Catatan: payload handoff utama dipecah per task agar lolos `subagent-handoff-check.py` dan mengurangi context drift antar-worker.

### Worklist tasks (atomic, owner-explicit)

Setiap task di bawah punya YAML payload terstruktur untuk `@orchestrator`/`@quality-gate`.

#### Task A1 — Daftarkan custom validator `alphanumDashUnderscore`

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Tambah custom validator `alphanumDashUnderscore` (regex ^[A-Za-z0-9_-]{1,20}$, trim whitespace) dan translasi id/en di master/pkg/validation/validation.go.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya bila registrasi + 2 translasi (id/en) ada dan `rtk go build ./...` sukses.
  source_basis:
    - master/pkg/validation/validation.go
    - master/pkg/validation/validation-product.go (pola alphanumericSpaceDash)
  must_preserve:
    - Pola custom validator existing (alphanumericSpace, qtystr, answer_frequency)
    - Pola translasi id/en
    - 0 dep baru (stdlib regexp cukup)
  do_not_touch:
    - File di luar master/pkg/validation/
    - master/go.mod, master/go.sum
    - Field validator lain (alphanumericSpace, dll.)
  validation:
    - cd master && rtk go build ./...
  exit_criteria:
    - function alphanumDashUnderscore terdaftar via RegisterValidation
    - translasi id (alphanumDashUnderscoreID) dan en (alphanumDashUnderscoreEN) ada
    - build sukses
  evidence_required:
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/build-A1.log
  depends_on: none
  context_bundle:
    verified_by_planner:
      - Pola referensi: alphanumericSpaceDash di master/pkg/validation/validation-product.go line 15-40 (confirmed_repo)
      - Bootstrap validator: master/pkg/validation/validation.go line 51-162 (confirmed_repo)
      - Regex target: ^[A-Za-z0-9_-]{1,20}$ (user_confirmed)
    files_already_read:
      - master/pkg/validation/validation.go
      - master/pkg/validation/validation-product.go
    open_assumptions: []
```

#### Task A2 — Ganti tag DTO distributor_code

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Ganti validate tag DistributorCode di CreateDistributorBody (line 163) dan UpdateDistributorRequest (line 212) dari `alphanum` ke `alphanumDashUnderscore`.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya jika 2 baris berubah dan `rtk go build ./...` sukses.
  source_basis:
    - master/entity/m_distributor.go
  must_preserve:
    - max=20 tetap
    - required/omitempty pattern tetap
    - Field lain (barcode, zip_code) TIDAK berubah
    - Response shape JSON TIDAK berubah
  do_not_touch:
    - File di luar master/entity/m_distributor.go
    - master/go.mod, master/go.sum
    - Tag validator untuk field lain
  validation:
    - cd master && rtk go build ./...
  exit_criteria:
    - Line 163 berisi `alphanumDashUnderscore`
    - Line 212 berisi `alphanumDashUnderscore`
    - build sukses
  evidence_required:
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/build-A2.log
  depends_on: A1
  context_bundle:
    verified_by_planner:
      - Root cause: tag `alphanum` di line 163 & 212 (confirmed_repo, baca langsung)
      - Library go-playground/validator/v10 regex `^[0-9a-zA-Z]*$` (confirmed_docs, bawaan stabil)
    files_already_read:
      - master/entity/m_distributor.go
    open_assumptions:
      - A1: tipe kolom DB Varchar(20) (per task prompt)
```

#### Task A3 — Test entity validator (7 sub-test)

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Tambah 7 sub-test table-driven untuk DistributorCode validator di master/entity/m_distributor_validation_test.go (atau file baru *_test.go di package entity): alphanumeric PASS, backward-compat numeric PASS, spasi ditolak, titik ditolak, 21 char ditolak, kosong ditolak (Create), underscore PASS, dash PASS.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya jika semua 7 sub-test PASS.
  source_basis:
    - master/entity/m_distributor_validation_test.go (pola)
    - master/pkg/validation/validation.go
  must_preserve:
    - Pola test existing (validation.NewValidator() + Validator.Struct())
  do_not_touch:
    - File di luar master/entity/ kecuali _test.go baru
    - Test existing selain menambah; jangan hapus/ubah
  validation:
    - cd master && rtk go test ./entity -run "Distributor" -v
  exit_criteria:
    - 7 sub-test PASS
    - 0 test existing regresi
  evidence_required:
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/test-entity-A3.log
  depends_on: A2
  context_bundle:
    verified_by_planner:
      - Pattern test existing di master/entity/m_distributor_validation_test.go (confirmed_repo)
    files_already_read:
      - master/entity/m_distributor_validation_test.go
    open_assumptions: []
```

#### Task A4 — Test service Update dengan alphanumeric

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Tambah 1 test service-level `TestDistributorService_Update_AcceptsAlphanumericCode` (atau nama setara) di master/service/m_distributor_service_test.go (atau _test.go baru). Mock FindOneByDistributorCodeAndCustId return sql.ErrNoRows, mock UPDATE return rows affected 1, ExpectCommit. Payload DistributorCode = `DIST-15676761A`.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya jika 1 test baru PASS dan test service existing tidak regresi.
  source_basis:
    - master/service/m_distributor_update_partial_test.go (pola)
    - master/service/m_distributor_service_test.go (helper)
  must_preserve:
    - setupDistributorServiceTest helper
    - serviceInt64Ptr/serviceStrPtr helper
  do_not_touch:
    - File di luar master/service/ kecuali _test.go baru
    - Test existing
  validation:
    - cd master && rtk go test ./service -run "Distributor" -v
  exit_criteria:
    - 1 test baru PASS
    - 0 test existing regresi
  evidence_required:
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/test-service-A4.log
  depends_on: A2
  context_bundle:
    verified_by_planner:
      - Pola testify + sqlmock + setupDistributorServiceTest (confirmed_repo)
      - Service Update tidak cast numeric (confirmed_repo, baca 398 baris)
    files_already_read:
      - master/service/m_distributor_service_test.go
      - master/service/m_distributor_update_partial_test.go
    open_assumptions: []
```

#### Task A5 — Full build + test master

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Jalankan rtk go test ./... dan rtk go build ./... di master/. Pastikan 0 failure, 0 error.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya jika kedua perintah exit code 0.
  source_basis:
    - .opencode/docs/QUALITY.md
  must_preserve: []
  do_not_touch: []
  validation:
    - cd master && rtk go test ./...
    - cd master && rtk go build ./...
  exit_criteria:
    - exit code 0 untuk kedua perintah
    - tidak ada test existing regresi
  evidence_required:
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/test-full-A5.log
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/build-A5.log
  depends_on: A3, A4
  context_bundle:
    verified_by_planner:
      - Cheat sheet per-service di .opencode/docs/QUALITY.md line 18-32 (confirmed_repo)
    files_already_read:
      - .opencode/docs/QUALITY.md
    open_assumptions:
      - A4: 0 dep baru; cukup stdlib regexp + validator existing
```

#### Task A6 — Smoke test manual local

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Boot service master lokal (rtk docker compose ... up -d master atau rtk go run main.go), lakukan curl PATCH ke distributor id=128 (atau distributor uji) dengan payload alfanumerik DIST128-AB. Verifikasi response 200 dan distributor_code string utuh di GET. Catat di evidence.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya jika PATCH 200 + GET round-trip identik. Token TIDAK boleh di-hardcode; ambil dari env lokal.
  source_basis:
    - master/.env (env lokal)
    - master/controller/m_distributor_controller.go (route)
  must_preserve:
    - Data existing (setelah PATCH test, kembalikan ke nilai awal bila distributor uji dipakai)
    - Pola error response jika gagal
  do_not_touch:
    - DB selain 1 row distributor uji (rollback setelah test)
    - File source apa pun
  validation:
    - curl -X PATCH $BASE_URL/master/v1/distributors/128 -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"distributor_code":"DIST128-AB","distributor_name":"..."}' -w "\nHTTP %{http_code}\n"
    - curl -X GET $BASE_URL/master/v1/distributors/128 -H "Authorization: Bearer $TOKEN" -w "\nHTTP %{http_code}\n"
  exit_criteria:
    - PATCH HTTP 200, response distributor_code="DIST128-AB"
    - GET HTTP 200, response distributor_code="DIST128-AB" (round-trip identik)
  evidence_required:
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/curl-pre-fix.log
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/curl-post-fix.log
  depends_on: A5
  context_bundle:
    verified_by_planner:
      - Route PATCH: master/controller/m_distributor_controller.go line 42 (confirmed_repo)
      - Validasi: line 311 controller.validator.ValidateStruct(request, ...) (confirmed_repo)
    files_already_read:
      - master/controller/m_distributor_controller.go
    open_assumptions: []
```

#### Task A7 — Verifikasi manual staging

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Ulangi smoke PATCH + GET di environment staging best.staging.scyllax.online. Catat timestamp, distributor id, request, response di evidence.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya jika PATCH 200 + GET round-trip identik; rollback distributor staging ke nilai awal.
  source_basis:
    - .opencode/docs/ARCHITECTURE.md (default staging lokal)
  must_preserve:
    - Data staging (rollback nilai distributor_code ke kondisi awal)
  do_not_touch:
    - DB staging selain 1 row distributor uji
    - Service production
  validation:
    - PATCH $STAGING_URL/master/v1/distributors/$ID_STAGING dengan distributor_code=DIST-<id>-ALPHA
    - GET $STAGING_URL/master/v1/distributors/$ID_STAGING
  exit_criteria:
    - PATCH 200, response distributor_code=DIST-<id>-ALPHA
    - GET 200, response distributor_code=DIST-<id>-ALPHA
    - Catatan timestamp + curl output di evidence
  evidence_required:
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/staging-verify.md
  depends_on: A6
  context_bundle:
    verified_by_planner:
      - Staging hostname: best.staging.scyllax.online (user_confirmed)
    files_already_read: []
    open_assumptions: []
```

#### Task A8 — PR ringkas

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: fixer
  scope: Tulis PR body dengan: scope, daftar path:line, sebelum/sesudah curl, ringkasan test, link evidence.
  claim_level: scoped
  claim_scope: Boleh klaim DONE hanya jika PR body memenuhi checklist §Acceptance dan path:line dilampirkan.
  source_basis:
    - .opencode/plans/20260706-1405-sx-2454-master-distributor-code-alphanumeric.md
  must_preserve:
    - 1 commit (atau multi-commit logis dengan pesan jelas)
  do_not_touch: []
  validation:
    - PR body review
  exit_criteria:
    - PR URL tersedia
    - Body berisi path:line, ringkasan test, sebelum/sesudah curl
  evidence_required:
    - PR link dicatat di final summary
  depends_on: A5, A6, A7
  context_bundle:
    verified_by_planner:
      - Scope file: 4 file di master/ (confirmed via plan)
    files_already_read:
      - .opencode/plans/20260706-1405-sx-2454-master-distributor-code-alphanumeric.md
    open_assumptions: []
```

#### Task Q1 — Quality gate

```yaml
handoff:
  task_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  plan_id: 20260706-1405-sx-2454-master-distributor-code-alphanumeric
  caller: orchestrator
  callee: quality-gate
  scope: Review akhir: security, regression, evidence, checklist §Acceptance.
  claim_level: scoped
  claim_scope: APPROVE atau REQUEST_CHANGES; tidak mengeksekusi perubahan kode.
  source_basis:
    - .opencode/plans/20260706-1405-sx-2454-master-distributor-code-alphanumeric.md
    - .opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/
  must_preserve:
    - Read-only posture
  do_not_touch:
    - File source/eksisting di master/
  validation:
    - Cek evidence path (curl-*.log, test-*.log, staging-verify.md)
    - Cek diff netto <= 80 baris
    - Cek tidak ada perubahan DB / response shape
    - Cek backward-compat numeric existing
  exit_criteria:
    - APPROVE dengan checklist §Acceptance lengkap; atau
    - REQUEST_CHANGES dengan daftar blocker spesifik
  evidence_required:
    - Quality gate notes
  depends_on: A8
  context_bundle:
    verified_by_planner:
      - Acceptance criteria 1-7 di plan
    files_already_read: []
    open_assumptions: []
```

#### Daftar task atomic (1-liner untuk progress tracker)

1. **A1** | `@fixer` | Tambah `alphanumDashUnderscore` + translasi `id`/`en` di `master/pkg/validation/validation.go`
2. **A2** | `@fixer` | Ganti tag `alphanum` → `alphanumDashUnderscore` di `master/entity/m_distributor.go` (2 baris)
3. **A3** | `@fixer` | Tambah 7 sub-test validator di `master/entity/m_distributor_validation_test.go`
4. **A4** | `@fixer` | Tambah 1 test service-level PATCH alfanumerik di `master/service/m_distributor_service_test.go`
5. **A5** | `@fixer` | Full test suite + build sukses di `master/`
6. **A6** | `@fixer` | Smoke test manual local (PATCH alfanumerik → 200)
7. **A7** | `@fixer` (atau user/QA) | Verifikasi manual staging `best.staging.scyllax.online`
8. **A8** | `@fixer` | Tulis PR ringkas + checklist §Acceptance
9. **Q1** | `@quality-gate` | Review akhir

### Execution ownership table

| Subsystem | Implementation owner | Review owner |
|---|---|---|
| Custom validator + translasi (`pkg/validation/validation.go`) | `@fixer` | `@quality-gate` |
| DTO tag update (`entity/m_distributor.go`) | `@fixer` | `@quality-gate` |
| Entity test | `@fixer` | `@quality-gate` |
| Service test | `@fixer` | `@quality-gate` |
| Smoke manual + staging verify | `@fixer` (+ user/QA untuk akses) | `@quality-gate` |
| PR & handoff | `@fixer` | `@quality-gate` |
| Cross-module deep audit (slice 2) | parked (`@explorer`/`@librarian` saat dipromosikan) | parked |

### Subagent Context Bundle (per-worker-task)

Untuk A1..A5: bundle di atas (sudah include verified facts, files_already_read, open_assumptions).

## Validation Commands

```bash
# dari direktori master/
cd master
rtk go mod download
rtk go mod tidy
rtk go build ./...

# Validator test
rtk go test ./entity -run "Distributor" -v

# Service test
rtk go test ./service -run "Distributor" -v

# Full test
rtk go test ./...

# (Opsional) Compose-level smoke
rtk docker compose -f docker-compose.yml ps
rtk docker compose -f docker-compose.yml up -d master

# Smoke test manual
BASE_URL="http://localhost:9002"   # sesuaikan dengan .env master
TOKEN="$TOKEN_STAGING"             # gunakan env, JANGAN hardcode

# Pre-fix expectation: 400 (alphanumeric)
# Post-fix expectation: 200 OK
curl -X PATCH "$BASE_URL/master/v1/distributors/128" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"distributor_code":"DIST128-AB","distributor_name":"Dist Sapi Madura"}' \
  -w "\nHTTP %{http_code}\n"

curl -X GET "$BASE_URL/master/v1/distributors/128" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nHTTP %{http_code}\n"
```

**Expected output post-fix**:
- PATCH → 200, body berisi `"distributor_code":"DIST128-AB"`.
- GET → 200, body `"distributor_code":"DIST128-AB"` (string utuh, termasuk `-`).

**Staging (best.staging.scyllax.online)**:
```bash
BASE_URL="https://best.staging.scyllax.online"
# ulangi PATCH + GET dengan distributor staging
```

## Evidence Requirements

Untuk klaim `PASS`/`PASS_FOR_SLICE`, simpan di `.opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/`:

- `discovery.md` — sudah ada (sudah dirangkum dari prompt).
- `test-output.log` — output `rtk go test ./...` (semua PASS).
- `curl-pre-fix.log` — output curl PATCH sebelum fix (400 expected).
- `curl-post-fix.log` — output curl PATCH setelah fix (200 expected, body berisi `DIST128-AB`).
- `staging-verify.md` — catatan verifikasi staging + timestamp.
- `index.json` — manifest (lihat §Evidence Index Schema di bawah).
- `grep-post-fix.txt` — `rg -n "validate:\".*alphanum.*DistributorCode" master` kosong (atau eksplisit mengapa ada match residual).

### Evidence Index Schema (`.opencode/evidence/<task-id>/index.json`)

```json
{
  "task_id": "20260706-1405-sx-2454-master-distributor-code-alphanumeric",
  "artifacts": [
    {"path": "discovery.md", "kind": "discovery", "summary": "..."},
    {"path": "test-output.log", "kind": "test_log", "command": "cd master && rtk go test ./...", "exit_code": 0},
    {"path": "curl-pre-fix.log", "kind": "smoke_log", "command": "curl PATCH pre-fix", "http_code": 400},
    {"path": "curl-post-fix.log", "kind": "smoke_log", "command": "curl PATCH post-fix", "http_code": 200},
    {"path": "staging-verify.md", "kind": "runtime_log", "target": "best.staging.scyllax.online"}
  ],
  "readiness": "PASS_FOR_SLICE"
}
```

## Done Criteria

1. Semua file di §Expected Files to Change ter-update.
2. `rtk go test ./...` di `master/` lulus 100% (tidak ada test existing yang regresi).
3. 7 test validator + 1 test service PATCH alfanumerik tertulis & lulus.
4. `grep -RIn "DistributorCode" --include="*.go" master | rg "validate:\".*\\balphanum\\b"` tidak lagi match (atau hanya match di luar scope — jelaskan jika ada).
5. Smoke test manual: PATCH + GET alfanumerik round-trip identik.
6. Verifikasi staging: PATCH + GET `best.staging.scyllax.online` round-trip identik, dicatat dengan timestamp.
7. PR deskripsi menjelaskan path:line, ringkasan test, dan sebelum/sesudah curl.
8. `.opencode/evidence/<task-id>/index.json` terisi dan siap dirujuk oleh quality-gate.

## Final Planning Summary

### Artifacts consulted/created
- Dibaca: `master/{entity,service,repository,controller}/m_distributor*.go`, `master/pkg/validation/*.go`, `.opencode/docs/{index,ARCHITECTURE,AGENT_ROUTING,QUALITY,SERVICE_MATRIX}.md`.
- Dilewati (tidak ada di repo): `PROJECT_STACK.md`, `PROJECT_COMMANDS.md`, `FRAMEWORK_PLAYBOOK.md`, `PROJECT_DETECTED_TOOLS.md`, `monitoring_activity_be_doc.txt`, dan prompt eksternal Jira.
- Diciptakan: `.opencode/plans/20260706-1405-sx-2454-master-distributor-code-alphanumeric.md` (plan utama), `.opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/discovery.md` (bukti), `.opencode/draft/20260706-1405-sx-2454-master-distributor-code-alphanumeric/decisions.md` (keputusan).

### Key decisions
- Regex: `^[A-Za-z0-9_-]{1,20}$` (user-confirmed, super-set dari `alphanum`).
- Tag baru: `alphanumDashUnderscore` di `pkg/validation/validation.go`.
- Scope: slice 1 = `master`; slice 2 parked.
- Staging source of truth: `best.staging.scyllax.online`.

### Assumptions
- A1: kolom DB `Varchar(20)` (per task prompt; model Go `string`).
- A2: tidak ada import/export Excel khusus distributor di `master` (grep negatif).
- A3: cross-module sudah `string` (grep positif; runtime check parked).
- A4: 0 dep baru.

### Open questions
- (None blocking). Jika QA menemukan import/export atau join yang masih numeric, naikkan ke slice 2 (tiket turunan).

### Readiness
- **Plan quality gate**: `PASS_FOR_SLICE` (slice 1 di `master` aman; slice 2 parked).
- **Eksekusi**: boleh dilakukan oleh `@fixer` (atau `@backend`) dengan mengikuti Executor Handoff Prompt. Active lane berikutnya: `@orchestrator` (lihat §Active-lane reset).

### Active-lane reset
- Lane saat ini: `@artifact-planner` (read-only, hanya tulis artefak di `.opencode/`).
- Setelah plan disetujui, lane eksekusi adalah `@orchestrator` → `@fixer` (atau domain equivalent).
- Setiap lane eksekusi WAJIB me-refresh permission/konteksnya sendiri (lihat `AGENT_ROUTING.md`); batasan read-only planner TIDAK otomatis terbawa.

### Cleanup performed
- File ini adalah primary plan; siap dirujuk oleh quality-gate.
- `.opencode/draft/20260706-1405-sx-2454-master-distributor-code-alphanumeric/decisions.md` — boleh dihapus setelah eksekusi selesai (sudah terkonsolidasi di sini).
- `.opencode/evidence/20260706-1405-sx-2454-master-distributor-code-alphanumeric/discovery.md` — boleh dihapus setelah eksekusi selesai (sudah dirangkum di sini), kecuali quality-gate ingin merujuknya.
