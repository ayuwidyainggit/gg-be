# Outlet Management Enhancement - Testing Guide

## Overview
Dokumentasi testing untuk enhancement fitur Outlet Management yang menambahkan support untuk `file_url` field.

## Changes Summary

### 1. Database Changes
**File**: `mobile/migration/add_file_url_to_m_outlet.sql`
- Added column: `file_url TEXT` (nullable)
- No breaking changes - backward compatible

### 2. Code Changes

#### Entity Layer
**File**: `mobile/entity/m_outlet.go`
- Updated `CreateMOutletBody`: Added `FileUrl *string` with validation `omitempty,url`
- Updated `MOutletRespone`: Added `FileUrl *string`
- Updated `OutletRespone`: Added `FileUrl *string`

#### Model Layer
**File**: `mobile/model/m_outlet.go`
- Updated `MOutlet`: Added `FileUrl *string` with db tag
- Updated `OutletReads`: Added `FileUrl *string` with db tag

#### Repository Layer
**File**: `mobile/repository/m_outlet.go`
- Updated `Store()`: Modified INSERT query to include `file_url` column
- Updated `FindOneByOutletIdAndCustId()`: Modified SELECT to include `file_url`

#### Service Layer
**File**: `mobile/service/m_outlet.go`
- No changes needed - `structs.Automapper` handles field mapping automatically

## API Endpoints Status

### ✅ Endpoint 1: Upload File (Existing)
**Status**: Already implemented
```
POST /mobile/v1/files/uploads
```
- Returns `file_url` in response
- No changes needed

### ✅ Endpoint 2: Create Outlet (Enhanced)
**Status**: Enhanced with file_url support
```
POST /mobile/v1/m-outlets
```

**Request Body** (Enhanced):
```json
{
  "outlet_name": "Outlet Test",
  "created_by": 120,
  "address": "Jalan Contoh No. 123",
  "phone_no": "081234567890",
  "building_own": 1,
  "latitude": "-7.7671403",
  "longitude": "110.4243715",
  "file_url": "https://storage.example.com/outlets/photo123.jpg",
  "details": {
    "contact": [
      {
        "contact_name": "John Doe",
        "job_title": "Manager",
        "phone_no": "081234567890",
        "wa_no": "081234567890",
        "email": "john@example.com",
        "identity_no": null
      }
    ]
  }
}
```

**Changes**:
- `file_url` field is now accepted (optional)
- Validates URL format if provided
- Saved to database

### ✅ Endpoint 3: List Outlets (Existing)
**Status**: Already implemented
```
GET /mobile/v1/outlet
```

**Query Parameters**:
- `q`: Search by outlet_name or outlet_code
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 5)
- `sort`: Sort field:direction (default: outlet_code:asc)
- `is_active`: Filter by active status (0=inactive, 1=active)

**Response**:
```json
{
  "message": "success",
  "data": [
    {
      "outlet_id": 123,
      "outlet_code": "OUT001",
      "outlet_name": "Toko ABC",
      "outlet_status": 1,
      "address1": "Jl. Contoh No. 123",
      "latitude": "-6.2",
      "longitude": "106.816666",
      "avg_sales_week": 1500000.00
    }
  ],
  "paging": {
    "total_record": 100,
    "page_current": 1,
    "page_limit": 10,
    "page_total": 10
  },
  "request_id": "req-abc-123456"
}
```

### ✅ Endpoint 4: Delete Outlet (Existing)
**Status**: Already implemented (soft delete)
```
DELETE /mobile/v1/m-outlets/:outlet_id
```

**Behavior**:
- Soft delete: Sets `is_del=true`, `deleted_by`, `deleted_at`
- Deleted outlets won't appear in list queries

## Testing Scenarios

### Scenario 1: Create Outlet WITH file_url
**Steps**:
1. Upload a file via `POST /mobile/v1/files/uploads`
2. Get `file_url` from response
3. Create outlet with the `file_url`
4. Verify outlet created successfully
5. Get outlet detail - verify `file_url` is saved

**Expected Result**: ✅ Outlet created with file_url

### Scenario 2: Create Outlet WITHOUT file_url (Backward Compatibility)
**Steps**:
1. Create outlet WITHOUT `file_url` field in request
2. Verify outlet created successfully
3. Get outlet detail - verify `file_url` is null

**Expected Result**: ✅ Outlet created successfully (backward compatible)

### Scenario 3: Create Outlet with Invalid URL
**Steps**:
1. Create outlet with invalid URL in `file_url` (e.g., "not-a-url")
2. Expect validation error

**Expected Result**: ✅ Validation error returned

### Scenario 4: List Outlets
**Steps**:
1. Create multiple outlets (some with file_url, some without)
2. List outlets with pagination
3. Verify all outlets returned correctly

**Expected Result**: ✅ All outlets listed, pagination works

### Scenario 5: Delete Outlet
**Steps**:
1. Create an outlet
2. Delete the outlet
3. List outlets - verify deleted outlet not shown
4. Check database - verify `is_del=true`

**Expected Result**: ✅ Soft delete works correctly

### Scenario 6: Get Outlet Detail
**Steps**:
1. Create outlet with file_url
2. Get outlet detail by ID
3. Verify `file_url` included in response

**Expected Result**: ✅ Detail returns file_url

## Test Cases Checklist

### Unit Tests (if applicable)
- [ ] Entity validation with valid file_url
- [ ] Entity validation with invalid file_url
- [ ] Entity validation without file_url
- [ ] Repository Store() saves file_url correctly
- [ ] Repository FindOne() retrieves file_url

### Integration Tests
- [ ] POST /v1/m-outlets with file_url
- [ ] POST /v1/m-outlets without file_url
- [ ] GET /v1/m-outlets/:id returns file_url
- [ ] GET /v1/outlet filters work correctly
- [ ] DELETE /v1/m-outlets/:id soft deletes

### API Test Examples (cURL)

#### 1. Upload File
```bash
curl --location 'http://localhost:9008/mobile/v1/files/uploads' \
--header 'Authorization: Bearer YOUR_TOKEN' \
--form 'file=@"/path/to/image.jpg"' \
--form 'folder="outlets"' \
--form 'file_type="image"'
```

#### 2. Create Outlet with file_url
```bash
curl --location 'http://localhost:9008/mobile/v1/m-outlets' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YOUR_TOKEN' \
--data '{
  "outlet_name": "Test Outlet",
  "created_by": 120,
  "address": "Jl. Test No. 123",
  "phone_no": "081234567890",
  "building_own": 1,
  "latitude": "-7.7671403",
  "longitude": "110.4243715",
  "file_url": "https://storage.example.com/outlets/photo.jpg",
  "details": {
    "contact": [
      {
        "contact_name": "Test User",
        "job_title": "Manager",
        "phone_no": "081234567890",
        "wa_no": "081234567890",
        "email": "test@example.com"
      }
    ]
  }
}'
```

#### 3. List Outlets
```bash
curl --location 'http://localhost:9008/mobile/v1/outlet?page=1&limit=10&sort=outlet_code:asc' \
--header 'Authorization: Bearer YOUR_TOKEN'
```

#### 4. Get Outlet Detail
```bash
curl --location 'http://localhost:9008/mobile/v1/m-outlets/123' \
--header 'Authorization: Bearer YOUR_TOKEN'
```

#### 5. Delete Outlet
```bash
curl --location --request DELETE 'http://localhost:9008/mobile/v1/m-outlets/123' \
--header 'Authorization: Bearer YOUR_TOKEN'
```

## Database Verification

### Check if file_url column exists
```sql
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_schema = 'mst' 
  AND table_name = 'm_outlet' 
  AND column_name = 'file_url';
```

### Check outlet with file_url
```sql
SELECT outlet_id, outlet_name, file_url, created_at
FROM mst.m_outlet
WHERE file_url IS NOT NULL
ORDER BY created_at DESC
LIMIT 10;
```

### Check soft deleted outlets
```sql
SELECT outlet_id, outlet_name, is_del, deleted_by, deleted_at
FROM mst.m_outlet
WHERE is_del = true
ORDER BY deleted_at DESC
LIMIT 10;
```

## Rollback Plan

If issues occur:

### Code Rollback
```bash
git revert <commit-hash>
```

### Database Rollback (Optional)
```sql
-- Column can remain as it's nullable and won't break anything
-- But if needed to remove:
ALTER TABLE mst.m_outlet DROP COLUMN IF EXISTS file_url;
```

## Performance Considerations

1. **Index**: No index added on `file_url` as it's not used for filtering
2. **Storage**: TEXT column is efficient for URLs
3. **Nullable**: No impact on existing queries
4. **Backward Compatible**: All existing code continues to work

## Notes

- `file_url` is optional - maintains backward compatibility
- URL validation prevents invalid data
- Automapper in service handles field mapping automatically
- No performance impact - nullable column
- List endpoints already implemented (GET /v1/outlet)
- Delete endpoint already implemented (DELETE /v1/m-outlets/:outlet_id)

## Sign-off Checklist

- [ ] Database migration executed successfully
- [ ] All unit tests pass (if applicable)
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Documentation updated
- [ ] Code reviewed
- [ ] Deployed to staging
- [ ] Smoke tests on staging pass
- [ ] Ready for production deployment

---
**Last Updated**: 2026-01-08
**Version**: 1.0.0