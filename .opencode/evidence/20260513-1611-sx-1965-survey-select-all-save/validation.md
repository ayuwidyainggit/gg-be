# Validation — SX-1965 Survey Select All Save

## Pre-check runtime

Command:

```bash
rtk docker compose -f docker-compose.yml ps
```

Hasil ringkas:

- `scylla-master` status `Up`
- `scylla-system`, `scylla-sales`, `scylla-pjp`, dan `scylla-redis` juga `Up`

Command:

```bash
curl http://localhost:9002/ping
```

Hasil:

```text
It works
```

## Test targeted

Command:

```bash
rtk go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'
```

Hasil:

```text
Go test: 22 passed in 1 packages
```

Command:

```bash
rtk go test ./controller -run 'TestSurveyController_Create'
```

Hasil:

```text
Go test: 3 passed in 1 packages
```

## Full regression

Command:

```bash
rtk go test ./...
```

Hasil:

```text
Go test: 257 passed in 23 packages
```

## Runtime scope evidence

### Query DB master data payload QA

Hasil query menunjukkan:

| emp_id | sales_name | cust_id | parent_cust_id | distributor_id | sales_team_code | sales_team_name | is_active |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 415 | Jaka | C260020001 | C26002 | 102 | 01 | MIX | true |
| 421 | Piere Njangka | C260020001 | C26002 | 102 | 02 | GT | true |
| 435 | Erling Braut Caraka | C26002 | C26002 | null | 77 | Tim Yuhu | false |
| 450 | Bagus Prima | C26002 | C26002 | null | 77 | Tim Yuhu | true |
| 458 | Richard | C260020001 | C26002 | 102 | 02 | GT | true |
| 459 | Rizal | C260020001 | C26002 | 102 | 02 | GT | true |
| 466 | Subiwo | C260020001 | C26002 | 102 | 02 | GT | true |

### Percobaan hit endpoint salesman scope FE

Percobaan ke endpoint salesman dengan scope FE dari environment ini belum memberi evidence final karena token runtime yang tersedia tidak lolos autentikasi di host yang dicoba:

1. Ke host remote `https://best.scyllax.online/...` → `Unauthorized`
2. Ke local `http://localhost:9002/v1/salesman?...` dengan token yang sama → `Unauthorized`

Kesimpulan evidence saat ini:

- Konsistensi scope code-level sudah diikat lewat reuse builder `buildSalesmanCustScopeCondition(...)` dan validasi service yang kini menerima principal + child distributor saat `distributor_id` mengandung `0`.
- Evidence runtime langsung terhadap endpoint salesman masih perlu dilengkapi bila token dev yang valid tersedia di sesi retest manual.

### Evidence deterministik setara scope endpoint salesman dari DB

Sebagai pengganti sementara runtime hit yang tertahan autentikasi, dijalankan query DB yang meniru scope endpoint salesman untuk filter FE:

- principal scope: `cust_id = 'C26002'`
- child distributor scope: distributor `102,103,119` di bawah `parent_cust_id = 'C26002'`
- sales team code: `82,81,80,78,77,66,65`
- emp_id subset: payload QA `450,435,415,421,458,459,466`

Hasil:

| emp_id | sales_name | cust_id | is_active | sales_team_code | sales_team_name |
| --- | --- | --- | --- | --- | --- |
| 435 | Erling Braut Caraka | C26002 | false | 77 | Tim Yuhu |
| 450 | Bagus Prima | C26002 | true | 77 | Tim Yuhu |

Interpretasi:

- Untuk scope FE yang diberikan user, payload QA **tidak sepenuhnya** termasuk dalam himpunan salesman yang lolos scope list salesman berbasis query setara ini.
- Dari tujuh `emp_id` payload QA, yang cocok dengan scope FE + sales team FE hanya dua salesman principal-owned di team `77`: `435` dan `450`.
- Ini menguatkan bahwa ada mismatch nyata antara payload create QA dan source scope salesman FE yang diklaim, sehingga retest manual perlu memastikan apakah:
  1. payload QA memang diambil dari FE response yang sama, atau
  2. ada filter FE lain/behavior FE tambahan yang belum tercermin di query ini.

## Perubahan yang tervalidasi

1. Create/update survey menerima validasi scope principal + child distributor saat sentinel `0` hadir.
2. `distributor_id=0` tidak ikut lookup distributor positif.
3. Invalid salesman kini dikembalikan dengan payload actionable:
   - `invalid_emp_id`
   - `invalid_salesman`
4. Error write DB dibungkus dengan context lebih jelas dan ada logging sementara terstruktur di service.
5. Full test suite modul `master` tetap hijau.

## Retest local dengan user Princessa

### Login local

Endpoint:

```bash
POST http://localhost:9001/v1/users/login
```

Payload:

```json
{
  "email": "princessa@gmail.com",
  "password": "Admin123"
}
```

Hasil:

- HTTP `200`
- login berhasil untuk:
  - `user_id = 140`
  - `cust_id = C26002`
  - `parent_cust_id = C26002`
  - `distributor_id = 0`

### Retest endpoint salesman di local

Endpoint:

```bash
GET http://localhost:9002/v1/salesman?page=1&sort=sales_name:asc&sales_team_id=82,81,80,78,77,66,65&distributor_id=0,102,103,119&q=&limit=9999
```

Hasil:

- HTTP `200`
- `total_record = 8`
- Payload mencakup seluruh emp_id QA:
  - `450`
  - `435`
  - `415`
  - `421`
  - `458`
  - `459`
  - `466`
- Endpoint local juga mengembalikan `emp_id 435 / Erling Braut Caraka` dengan `is_active = false`.

Interpretasi:

- Source-of-truth runtime local untuk picker FE memang **mengembalikan** semua emp_id QA, termasuk salesman inactive `435`.
- Karena itu, create survey local yang selaras dengan picker FE memang seharusnya menerima payload tersebut.

### Retest create survey di local

Endpoint:

```bash
POST http://localhost:9002/v1/survey
```

Payload:

```json
{
  "survey_title": "Testing May 12 SX1965 local",
  "efective_date_start": "2026-05-12",
  "efective_date_end": "2026-05-13",
  "answer_frequency": "Multiple",
  "response_type": "Optional",
  "target_type": "Specific",
  "distributor_id": [0, 102, 103, 119],
  "area_id": [91, 88],
  "outlet_id": [],
  "survey_template_id": 53,
  "emp_id": [450, 435, 415, 421, 458, 459, 466]
}
```

Hasil:

- HTTP `201`
- response message: `Survey has been successfully created`

### Validasi database sesudah create

Survey header:

- `survey_id = 124`
- `cust_id = C26002`
- `survey_title = Testing May 12 SX1965 local`
- `emp_id = 450`
- `created_by = 140`
- `status = 1`

Survey area rows:

| distributor_id | area_id |
| --- | --- |
| 0 | 88 |
| 0 | 91 |
| 102 | 88 |
| 103 | 88 |
| 119 | 88 |

Survey salesman rows:

| salesman_id | cust_id |
| --- | --- |
| 415 | C260020001 |
| 421 | C260020001 |
| 435 | C26002 |
| 450 | C26002 |
| 458 | C260020001 |
| 459 | C260020001 |
| 466 | C260020001 |

Survey detail rows:

- `survey_template_id = 53`

Interpretasi:

- Retest local membuktikan fix berjalan end-to-end.
- Payload QA yang muncul dari endpoint salesman local berhasil disimpan di create survey local.
- Tidak ada indikasi partial insert pada header/area/salesman/detail untuk skenario sukses ini.

## Root cause ringkas untuk PR/Jira

- Validator target salesman pada create/update survey sebelumnya belum menyatukan scope principal dengan child distributor saat `distributor_id` mengandung sentinel `0`, sehingga flow select-all principal bisa lebih ketat daripada scope picker FE.
- Saat ada target salesman yang tidak lolos validasi, API hanya mengembalikan error generik tanpa `emp_id`/nama salesman invalid, sehingga QA tidak bisa melihat akar data yang bermasalah.
