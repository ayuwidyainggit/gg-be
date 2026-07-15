# Discovery — SX-2361 / SX-2364 Latest Location Pipeline

## Ringkasan
Task ini adalah maintenance/debug lintas `mobile` + `pjp` untuk memastikan marker user di Live Monitoring Distributor memakai latest actual location, bukan planned outlet coordinate, dan untuk merencanakan pipeline realtime berbasis RabbitMQ + WebSocket sesuai kontrak user/Jira prompt.

## Source strategy
- **Dipakai:** repo-local evidence (`pjp`, `mobile`, `.opencode/plans/*`, `.opencode/docs/*`).
- **Tidak dipakai:** official docs/context7 karena belum ada kebutuhan version-sensitive library behavior; problem ini dominan repo-specific.
- **Tidak dipakai:** GitHub/web search/browser karena kontrak teknis sudah diberikan user dan targetnya perilaku internal repo.
- **Skip reason:** tahap ini cukup dengan discovery lokal untuk menulis implementation-ready plan.

## Files inspected
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/SECURITY.md`
- `.opencode/docs/SERVICE_MATRIX.md`
- `.opencode/plans/20260521-1515-sx2034-extra-call-monitoring.md`
- `.opencode/plans/20260522-1101-sx-2038-monitoring-detail-null.md`
- `.opencode/plans/20260529-1136-sx-2097-monitoring-survey.md`
- `pjp/router/live_monitoring.go`
- `pjp/service/live_monitoring/live_monitoring_service.go`
- `pjp/service/live_monitoring/get_distributor_service.go`
- `pjp/repository/live_monitoring/live_monitoring_repository.go`
- `pjp/repository/live_monitoring/get_distributor_repository.go`
- `pjp/repository/live_monitoring/current_coordinate_selector.go`
- `pjp/model/live_monitoring.go`
- `pjp/data/response/live_monitoring_response.go`
- `mobile/controller/visits.go`

## Project patterns found
- `pjp` adalah owner endpoint monitoring map/detail (`GET /live-monitoring-distributor`, `GET /monitoring_locations/details`).
- Layering konsisten: `router -> controller -> service -> repository -> DB`.
- `pjp` sudah punya response field untuk actual user coordinate:
  - `current_longitude`
  - `current_latitude`
  - `current_coordinate_at`
  - `current_coordinate_source`
- Source selection actual coordinate **sudah ada** di repository, dengan priority:
  1. `mobile.attendances` checkout
  2. `mobile.attendances` any type
  3. `mobile.visits`
  4. `pjp.outlet_visit_list`
- Selector latest coordinate sudah ada di `pjp/repository/live_monitoring/current_coordinate_selector.go`, dan latest dipilih terutama berdasarkan timestamp terbaru, lalu `source_rank` bila timestamp sama.

## Reuse candidates
- Reuse `GetDistributorCurrentCoordinates(...)` di `pjp/repository/live_monitoring/get_distributor_repository.go` untuk initial-load latest location; tidak perlu desain ulang helper latest selector.
- Reuse response contract existing `current_*` daripada menambah field baru bila FE bisa menyesuaikan ke field existing. Bila FE task SX-2363 sudah bergantung pada nama `latest_*`, perlu koordinasi; repo lokal saat ini sudah expose `current_*`.
- Reuse pola test `sqlmock` dan service stub dari plan SX-2097 / test file live monitoring existing.
- Reuse wrapper/pola RabbitMQ dari module `sales` atau `master` bila implementation lane memutuskan menambahkan publish/subscribe di `mobile`/`pjp`; saat ini `pjp` belum punya package RabbitMQ sendiri.

## Key discoveries
1. **Initial load latest location sebenarnya sudah tersedia di response `pjp`, bukan nol dari planned field semata.**
   - Bukti: `pjp/data/response/live_monitoring_response.go` punya `current_*` fields.
2. **Repository actual location sudah benar secara high level dan memakai source actual, bukan planned route.**
   - Bukti: `GetDistributorCurrentCoordinates` membangun union dari attendance/visit/outlet-visit sources dan memilih kandidat terbaru.
3. **Service distributor saat ini menghapus current coordinate bila salesman tidak punya attendance pada requested date.**
   - Bukti: `pjp/service/live_monitoring/get_distributor_service.go` baris sekitar 227-256.
   - Dampak: marker user bisa hilang/0 walau ada data `mobile.visits` atau source actual lain. Ini mismatch terhadap kontrak user yang meminta `last known actual location`.
4. **Mobile service tidak memiliki endpoint `POST /location` atau publisher `location.updated` yang mudah ditemukan.**
   - Route yang ada terkait lokasi paling dekat adalah `/v1/visits/*` dan `/v1/attendances/*`.
5. **Repo belum menunjukkan bukti consumer RabbitMQ `web-monitoring-location` atau WebSocket push `location_updated` di `pjp`.**
   - Search broad untuk `location.exchange`, `location.updated`, `web-monitoring-location`, `websocket` tidak menemukan implementasi target.
6. **Docs file `monitoring_activity_be_doc.txt` yang disebut prompt tidak ditemukan di repo lokal ini.**
   - Maka source of truth plan adalah prompt user + bukti repo lokal.

## Constraints
- Jangan copy secret/token/credential dari env file tracked atau browser/Jira.
- Repo policy project-local untuk shell adalah `rtk`-prefixed, tetapi planner tidak menjalankan runtime/destructive step pada tahap ini.
- Planner lane hanya boleh menulis artifact `.opencode/**`, bukan source code.

## Risks
- Scope user menyebut Mobile BE publisher + broker + Web BE consumer + WebSocket, tetapi repo lokal belum memperlihatkan plumbing yang siap pakai di `mobile` dan `pjp`.
- Ada kemungkinan istilah “Web BE” pada prompt sebenarnya mengacu ke service lain, bukan `pjp`; namun endpoint live monitoring yang aktif jelas ada di `pjp`.
- Menambahkan realtime pipeline mungkin butuh env/config baru untuk `pjp` dan `mobile`, yang belum terkonfirmasi ada di module masing-masing.
- Existing FE mungkin belum memakai `current_*` walau field sudah ada; itu berada di scope SX-2363, bukan planner backend saja.

## Confirmed vs Assumed Audit
| Claim | Status | Basis |
|---|---|---|
| `GET /live-monitoring-distributor` ada di `pjp` | confirmed_repo | `pjp/router/live_monitoring.go` |
| `GET /monitoring_locations/details` ada di `pjp` | confirmed_repo | `pjp/router/live_monitoring.go` |
| Response distributor sudah punya `current_*` fields | confirmed_repo | `pjp/data/response/live_monitoring_response.go` |
| Actual coordinate source precedence attendance/visit/outlet-visit sudah ada | confirmed_repo | `pjp/repository/live_monitoring/get_distributor_repository.go` |
| Latest candidate dipilih berdasar timestamp terbaru | confirmed_repo | `pjp/repository/live_monitoring/current_coordinate_selector.go` |
| Service meng-gate current coordinate dengan attendance pada requested date | confirmed_repo | `pjp/service/live_monitoring/get_distributor_service.go` |
| `mobile` punya endpoint khusus `POST /location` | unverified | search tidak menemukan route itu |
| `mobile` publish `location.updated` ke RabbitMQ | unverified | search tidak menemukan string/publisher itu |
| `pjp` consume queue `web-monitoring-location` | unverified | search tidak menemukan consumer itu |
| `pjp` punya WebSocket hub untuk live monitoring | unverified | search tidak menemukan websocket implementation yang jelas |
| FE existing memakai `current_*` untuk marker user | unverified | FE repo tidak didiscovery di task ini |
| Initial load salah hanya karena FE pakai planned field | assumption | mungkin benar, perlu verifikasi FE/task SX-2363 |

## Commands/docs checked
- local file reads via planner tools
- broad grep for `location.exchange`, `location.updated`, `web-monitoring-location`, `websocket`, `visits`, `attendance`
- prior local planning artifacts for monitoring domain

## Conclusion
Ada dua jalur fix yang perlu dipisahkan dalam plan:
1. **Bugfix initial load yang sangat mungkin dan minimum diff:** perbaiki gating di service `pjp` agar latest actual location tidak dibuang hanya karena tidak ada attendance di hari itu.
2. **Gap feature realtime pipeline:** repo lokal belum menunjukkan publisher/consumer/WebSocket yang diminta kontrak SX-2364, sehingga implementation-ready plan harus memecah pekerjaan lintas `mobile` dan `pjp`, termasuk wiring RabbitMQ dan kanal push FE dengan evidence lebih eksplisit selama implementasi.
