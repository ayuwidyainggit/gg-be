Berikut rangkuman bug di endpoint **/master/v1/distributors** dan acceptance‑nya dalam bentuk tabel, dibagi per skenario di SX‑908 dan dokumen “Enhance Distributor _ BE”. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908)

### 1. Masalah & acceptance – Duplicate distributor code

| Aspek | Detail Masalah (Actual) | Acceptance Criteria (Expected) |
| --- | --- | --- |
| Validasi kode distributor unik | Endpoint POST /master/v1/distributors saat create distributor masih mengizinkan `distributor_code` yang sudah ada, sehingga terjadi data dan kode distributor **duplicate**. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | Saat submit, BE wajib cek ke `mst.m_distributor.distributor_code`; jika sudah digunakan, kembalikan error: **"Distributor code already exists. Please use a different distributor code."** dan data tidak tersimpan. |
| Status test | Skenario “Add new distributor will created as duplicate data and duplicate code distributor” dinyatakan **FAILED** pada QA. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | Setelah fix, skenario yang sama harus **PASSED**: create dengan kode yang sama harus ditolak dengan pesan error di atas. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |

### 2. Masalah & acceptance – Validasi Section 1 (Header distributor)

| Field | Masalah di Endpoint/DB | Acceptance Criteria |
| --- | --- | --- |
| distributor_code | Di body sudah ada, tapi di DB belum `NOT NULL`; validasi alphanumeric dan max length 20 karakter belum ditegakkan secara konsisten dari BE/DB. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | Wajib alphanumeric, max 20 karakter, **mandatory**, dan **unique**; di DB `mst.m_distributor.distributor_code` dibuat `NOT NULL` dan ada pengecekan unik, di request body juga wajib (required). |
| distributor_name | Di DB `mst.m_distributor.distributor_name` belum `NOT NULL`, dan validasi max 150 karakter belum ditegakkan penuh. | Wajib alphanumeric, max 150 karakter, **mandatory**; di DB dibuat `NOT NULL`, di request body juga wajib. |
| barcode | Di definisi UX: **Not Mandatory**, numeric; di BE/DB: tipe belum diset ke integer 13 digit, dan perilaku mandatory/optional tidak sepenuhnya align. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | Field **not mandatory**, numeric, panjang max 13 digit; di DB tipe `int(13)` (atau ekuivalen), di request body boleh kosong/null dan tidak men-trigger error required. |
| is_active (Status) | Saat ini di API sudah required, tapi di DB belum benar‑benar mandatory; kemungkinan masih nullable atau tidak divalidasi konsisten. | Tipe boolean, **mandatory**; di DB dibuat mandatory (NOT NULL / default jelas), di API tetap required dan menolak request tanpa nilai valid. |
| region_id | UX: mandatory dropdown, namun di DB belum mandatory meski di API sudah required. | **Mandatory**; di DB dibuat mandatory (NOT NULL / FK wajib), di API sudah required dan harus terus dipertahankan. |
| area_id | Sama dengan region_id: mandatory di UX, DB belum mandatory, API sudah required. | **Mandatory**; di DB mandatory, di API required. |
| channel_id | Mandatory di UX, DB belum mandatory, status di API belum tegas disebut, tapi untuk konsistensi harus required. | **Mandatory**; di DB mandatory, di API required seperti field mandatory lain. |
| sub_distributor_group_id | Mandatory di UX, DB belum mandatory, API sudah required. | **Mandatory**; di DB mandatory, di API required. |
| dist_price_grp_id | Mandatory di UX, DB belum mandatory, API sudah required. | **Mandatory**; di DB mandatory, di API required. |

### 3. Masalah & acceptance – Validasi Section 2 (Alamat)

| Field | Masalah di Endpoint/DB | Acceptance Criteria |
| --- | --- | --- |
| address | UX: Mandatory, alphanumeric, max 150 char; di API memang required tapi di DB belum mandatory dan panjang kolom belum di‑align (dicatat perlu jadi varchar 255). [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | **Mandatory**; di DB set sebagai NOT NULL, panjang minimal varchar 150 (disarankan varchar 255), di API required dan validasi panjang + alphanumeric. |
| province_id | UX: Not Mandatory; di API sudah nullable, tapi di skenario QA terlihat BE meng‑require field ini sehingga muncul error saat kosong. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | Field **not mandatory**; di API harus nullable dan tidak divalidasi sebagai required, tidak muncul di list error ketika kosong. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |
| regency_id | Sama seperti province_id: UX Not Mandatory, namun saat kirim kosong, respon error menyatakan field required. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | **Not mandatory**; di API nullable dan tidak ikut validasi required. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |
| sub_district_id | UX: Not Mandatory; di API saat ini masih kena validasi required ketika kosong. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | **Not mandatory**; di API tidak boleh dianggap required dan boleh kosong/null. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |
| ward_id | UX: Not Mandatory; di API masih divalidasi required sehingga QA mendapat error saat kosong. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | **Not mandatory**; di API nullable dan tidak divalidasi required. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |
| zip_code | UX: Not Mandatory, alphanumeric, max 6 char; di API sekarang ikut error list sebagai required ketika kosong. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | **Not mandatory**; di API nullable, validasi hanya ketika diisi (alphanumeric, max 6), tidak muncul di error list kalau dikosongkan. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |
| longitude | UX: Mandatory, alphanumeric, max 50; di API sudah mandatory, tapi di DB kolom belum dipastikan mandatory (dicatat “Mandatory ❌” di dokumen). | **Mandatory**; di API required, di DB kolom dibuat mandatory dengan panjang yang menampung hingga 50 karakter. |
| latitude | Sama kasus dengan longitude: mandatory di UX dan API, DB belum mandatory. | **Mandatory**; di API required, di DB mandatory dengan panjang hingga 50 karakter. |
| phone | UX: Not Mandatory, alphanumeric max 25; di API sudah nullable tapi validasi panjang/format belum sesuai (dicatat “alphanumeric 25 ❌”). | **Not mandatory**; validasi hanya saat diisi: alphanumeric, max 25 karakter; di API nullable. |
| fax_number | UX: Not Mandatory, alphanumeric max 25; kondisi sama seperti phone. | **Not mandatory**; validasi saat diisi: alphanumeric, max 25 karakter; di API nullable. |

### 4. Masalah & acceptance – Validasi Section 3 (Contact)

| Field | Masalah di Endpoint/DB | Acceptance Criteria |
| --- | --- | --- |
| contact_name | UX: mandatory; create sudah ok, tapi di skenario edit, field `email` dan mungkin contact jadi mandatory padahal di spesifikasi `email` tidak wajib. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | `contact_name` tetap **mandatory** dengan max 50 alphanumeric; pada edit maupun create, **email tetap optional** dan tidak boleh mem‑blok save jika kosong. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |
| job_title (Position) | UX: mandatory, max 20 alphanumeric; validasi panjang dan mandatory perlu dipastikan, namun belum dirinci error‑nya di ticket. | **Mandatory**; max 20 karakter, alphanumeric; error jelas jika kosong atau melebihi limit. |
| identity_type | UX: mandatory (National ID, Passport, Others ID); tidak ada catatan bug spesifik, tapi masuk list mandatory yang harus divalidasi per field. | **Mandatory**; hanya menerima salah satu dari opsi: National ID, Passport, Others ID. |
| identity_no | UX: mandatory, alphanumeric max 20; belum diceritakan error, tapi tetap harus ikut validasi per field. | **Mandatory**; alphanumeric, max 20 karakter; error ketika kosong atau melebihi panjang. |
| phone_no | UX: mandatory, numeric max 20; perlu validasi numeric dan panjang. | **Mandatory**; numeric, max 20 karakter, error jika kosong atau non‑numeric. |
| is_wa_no (Set as whatsapp) | UX: Not Mandatory, default No; tidak ada bug spesifik, tapi harus dipastikan tidak memblokir kalau tidak diisi (gunakan default). | **Not mandatory**; default false/No jika tidak dikirim, tetap boolean valid di DB. |
| wa_no | UX: Mandatory, numeric max 20; saat ini bekerja, tapi harus aman ketika is_wa_no diset true. | **Mandatory** ketika contact diisi, numeric, max 20 karakter, konsisten di create dan update. |
| email | UX: Not Mandatory; pada skenario edit distributor, email contact menjadi mandatory sehingga QA tidak bisa save ketika kosong. [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) | **Not mandatory**; BE tidak boleh mengembalikan error ketika email kosong pada create maupun update, validation hanya saat diisi (format email, max 100). [scyllax-pratesis.atlassian](https://scyllax-pratesis.atlassian.net/browse/SX-908) |