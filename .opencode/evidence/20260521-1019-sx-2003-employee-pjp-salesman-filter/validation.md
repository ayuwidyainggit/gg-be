# Validation Evidence — SX-2003

Tanggal: 2026-05-21
Validator: orchestrator + DB read-only

## Service

- system: localhost:9001 (running)
- master: localhost:9002 (running)

---

## Skenario 1 — Distributor (cust_id=C260020001, distributor_id=102)

Request:
```
GET http://localhost:9002/v1/employee-pjp?q=&page=1&limit=20&sort=emp_name:asc&is_active[]=1
Authorization: Bearer <distributor token>
```

Response: 9 employee
```
478 BM200  Arifin
480 BM400  Dimas Maisator
476 BM567  Herman
479 BM300  Herman Lee
415 EMP0021 Jaka
421 EMP0025 Piere Njangka
458 202604  Richard
459 R2026   Rizal
466 2026    Subiwo
```

DB check — semua 9 emp_id punya role salesman: ✅
DB check — tidak ada employee non-salesman aktif di scope C260020001: ✅ (no rows)
DB check — employee aktif tanpa sys.m_user (tidak muncul di API, benar): 433 Adul, 434 komeng, 471 Yabes Roni

---

## Skenario 2 — Principal default (cust_id=C26002 dari JWT)

Request:
```
GET http://localhost:9002/v1/employee-pjp?q=&page=1&limit=20&sort=emp_name:asc&is_active[]=1
Authorization: Bearer <principal token>
```

Response: 4 employee
```
483 MS0987 Ahmad Baihaki
450 MS0001 Bagus Prima
482 MS9990 Jihan Fahira
484 MS0909 Syaiful
```

DB check — semua 4 emp_id punya role salesman: ✅
DB check — tidak ada employee non-salesman aktif di scope C26002: ✅ (no rows)
DB check — employee aktif tanpa sys.m_user (tidak muncul di API, benar): 446 CEO, 447 CFO, 448 CFD

---

## Skenario 3 — Principal pilih cust_id=C26002 eksplisit via query param

Request:
```
GET http://localhost:9002/v1/employee-pjp?q=&page=1&limit=20&sort=emp_name:asc&is_active[]=1&cust_id=C26002
Authorization: Bearer <principal token>
```

Response: 4 employee (sama dengan skenario 2) ✅

---

## Skenario 4 — Principal pilih distributor_id=102

Request:
```
GET http://localhost:9002/v1/employee-pjp?q=&page=1&limit=20&sort=emp_name:asc&is_active[]=1&distributor_id=102
Authorization: Bearer <principal token>
```

Response: 9 employee (sama dengan skenario 1) ✅

---

## Kesimpulan

| Check | Status |
|---|---|
| Semua emp_id API response punya role salesman di DB | ✅ PASS |
| Employee non-salesman aktif tidak muncul di API | ✅ PASS |
| Employee tanpa sys.m_user tidak muncul (expected) | ✅ PASS |
| Distributor scope benar (C260020001) | ✅ PASS |
| Principal scope benar (C26002) | ✅ PASS |
| Principal pilih distributor_id=102 benar | ✅ PASS |
| Response contract tidak berubah (emp_id, emp_code, emp_name) | ✅ PASS |
| Paging total_record akurat | ✅ PASS |

SX-2003 validasi live: **PASS semua skenario**.
