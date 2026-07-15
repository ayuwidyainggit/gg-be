Smoke test manual belum dijalankan di sesi ini.

Alasan:
- Tidak ada token Bearer valid yang bisa saya pakai dari environment tanpa menyalin rahasia.
- Verifikasi staging `best.staging.scyllax.online` butuh akses runtime + target distributor uji yang aman untuk diubah/rollback.
- Saya sengaja tidak mengarang curl success tanpa bukti runtime.

Status:
- A6 local smoke: not-run
- A7 staging verify: not-run
- A8 PR creation: not-run (repo ini bukan git repo lokal; tidak ada remote/branch context untuk PR)

Perintah yang harus dijalankan operator/QA:

Local or staging PATCH:
```bash
BASE_URL="https://best.staging.scyllax.online"
TOKEN="$TOKEN"
DIST_ID="128"
NEW_CODE="DIST-${DIST_ID}-ALPHA"

curl -X PATCH "$BASE_URL/master/v1/distributors/$DIST_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"distributor_code\":\"$NEW_CODE\",\"distributor_name\":\"Dist Sapi Madura\"}" \
  -w "\nHTTP %{http_code}\n"

curl -X GET "$BASE_URL/master/v1/distributors/$DIST_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nHTTP %{http_code}\n"
```

Ekspektasi post-fix:
- PATCH -> HTTP 200
- GET -> distributor_code identik dengan nilai PATCH (`DIST-<id>-ALPHA`)
