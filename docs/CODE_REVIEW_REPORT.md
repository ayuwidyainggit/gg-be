# Code Review Report - Outlet Management Enhancement

**Date**: 2026-01-08
**Reviewer**: AI Code Assistant
**Scope**: Review kode outlet management untuk memastikan kesesuaian dengan docs dan idiomatic Go

---

## Summary

Kode outlet management telah direview dan diperbaiki untuk memastikan kesesuaian dengan dokumentasi `docs/Add Outlet & Manage Outlet_BE.md` serta best practices idiomatic Go.

## Issues Fixed

### 1. ✅ Endpoint URL tidak sesuai docs

**Before:**
- GET `/v1/outlet`
- DELETE `/v1/m-outlets/:outlet_id`

**After (sesuai docs):**
- GET `/v1/outlet-list`
- DELETE `/v1/outlet-list/:outlet_id`

**File Modified:** [`mobile/controller/m_outlet.go`](../mobile/controller/m_outlet.go:39)

### 2. ✅ Query Parameter tidak sesuai docs

**Before:**
- `is_active` (pointer int)

**After (sesuai docs):**
- `outlet_status` (Array<int>) - 1=Active, 2=Covered, 3=Non Active, 4=Closed
- `is_active` tetap dipertahankan untuk backward compatibility

**Files Modified:**
- [`mobile/entity/m_outlet.go`](../mobile/entity/m_outlet.go:335)
- [`mobile/repository/m_outlet.go`](../mobile/repository/m_outlet.go:718)

### 3. ✅ SQL Injection Prevention

**Before:**
```go
qWhere := ` WHERE o.is_del = false AND o.cust_id = '` + custId + `' `
```

**After (parameterized queries):**
```go
qWhere := ` WHERE o.is_del = false AND o.cust_id = $1 `
args := []interface{}{custId}
```

**File Modified:** [`mobile/repository/m_outlet.go`](../mobile/repository/m_outlet.go:718)

### 4. ✅ Magic Numbers replaced with Constants

**Before:**
```go
switch row.OutletStatus {
case 1:
    vResp.OutletStatusName = "Active"
case 2:
    vResp.OutletStatusName = "Covered"
// ...
}
```

**After:**
```go
// Constants defined in entity
const (
    OutletStatusActive   int = 1
    OutletStatusCovered  int = 2
    OutletStatusNonActive int = 3
    OutletStatusClosed   int = 4
)

// Usage
vResp.OutletStatusName = entity.GetOutletStatusName(int(row.OutletStatus))
```

**Files Modified:**
- [`mobile/entity/m_outlet.go`](../mobile/entity/m_outlet.go:1) - Added constants
- [`mobile/service/m_outlet.go`](../mobile/service/m_outlet.go:226) - Using constants

### 5. ✅ Default Values sesuai Docs

**Defaults applied in service layer:**
- `page` default: 1
- `limit` default: 9999
- `sort` default: `outlet_code:asc`

**File Modified:** [`mobile/service/m_outlet.go`](../mobile/service/m_outlet.go:384)

### 6. ✅ Response Messages sesuai Docs

**Empty state response:**
```json
{
  "message": "No Data",
  "data": null,
  "paging": {...}
}
```

**File Modified:** [`mobile/controller/m_outlet.go`](../mobile/controller/m_outlet.go:296)

---

## Files Modified

| File | Changes |
|------|---------|
| `mobile/entity/m_outlet.go` | Added outlet status & verification constants, updated query filter |
| `mobile/model/m_outlet.go` | No changes needed |
| `mobile/repository/m_outlet.go` | Parameterized queries, outlet_status filter |
| `mobile/service/m_outlet.go` | Using constants, fixed defaults |
| `mobile/controller/m_outlet.go` | New endpoint routes, response messages |
| `docs/Outlet_Management_Postman_Collection.json` | Updated with correct endpoints |

---

## Known Typos (Not Fixed - Breaking Change Risk)

Typos yang tidak diperbaiki karena berpotensi breaking change:
- `MOutletRespone` → seharusnya `MOutletResponse`
- `IsAddiitional` → seharusnya `IsAdditional`
- `BuldingOwn` → seharusnya `BuildingOwn`

**Recommendation:** Fix typos dalam refactor terpisah dengan migration plan.

---

## Build Status

```bash
$ cd mobile && go build -o /dev/null ./...
# Exit code: 0 (Success)
```

---

## API Endpoints (Final)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/mobile/v1/files/uploads` | Upload file (existing) |
| POST | `/mobile/v1/m-outlets` | Create outlet with file_url |
| GET | `/mobile/v1/outlet-list` | List outlets with filters |
| DELETE | `/mobile/v1/outlet-list/:outlet_id` | Soft delete outlet |
| GET | `/mobile/v1/m-outlets/:outlet_id` | Get outlet detail |
| DELETE | `/mobile/v1/m-outlets/:outlet_id` | Soft delete (legacy) |

---

## Testing

Import `docs/Outlet_Management_Postman_Collection.json` ke Postman untuk testing.