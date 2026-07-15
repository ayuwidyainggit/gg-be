# Implementation Discovery & Verification

## Scope Delivery

✓ Endpoint: GET /v1/monitoring_locations/update-locations
✓ DTOs: UpdateLocationsRequest, UpdateLocationsResponse, TimelineItem
✓ Model: UpdateLocationRow
✓ Repository: GetEmployeeRole, GetUpdateLocations (concrete impl in update_locations_repository.go)
✓ Service: GetUpdateLocations method
✓ Controller: GetUpdateLocations handler
✓ Route: Registered in live_monitoring.go
✓ Tests: 2 test files (controller + repository with sqlmock)

## Key Implementation Details

### Request Contract
- emp_id (int, required) via binding:"required"
- date (string, optional, YYYY-MM-DD format) via binding:"omitempty,datetime=2006-01-02"
- Default date: time.Now().In(jakartaLocation).Format("2006-01-02")

### Response Contract
- message: "Success" or "No Data"
- data.timeline: array of TimelineItem (ordered, no pagination)
- request_id: UUID v4

### Role Resolver
- Query: SELECT cust_id FROM mst.m_employee WHERE emp_id = ? AND cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = ? OR parent_cust_id = ?) LIMIT 1
- Branch: pjp (if len(cust_id) > 6) else pjp_principles
- Error: 404 ErrRecordNotFound if employee outside tenant

### Timeline Sources (UNION ALL)
1. mobile.attendances: type 1 → clock_in, else → clock_out
2. sys.user_location: type → gps
3. outlet_visit_list (branch): arrive_at/leave_at → arrive/leave

### Tenant Scope
- All 3 sources filtered: cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = $jwtCust OR parent_cust_id = $jwtCust)
- Employee resolved via mst.m_employee in role resolver
- Consistent with existing PJP patterns

### Error Handling
- 400 Bad Request: emp_id missing
- 401 Unauthorized: JWT middleware (existing)
- 404 Not Found: employee outside tenant
- 500 Internal Server Error: service error
- 200 OK: success or no data (message field differentiates)

### No Pagination
- No page/limit query params
- No paging response key
- Full one-day timeline returned
- Order: recorded_at ASC, type ASC, destination_id ASC

### Context & Timeout
- context.WithTimeout 30 seconds (parity with existing endpoints)
- Request ID: UUID v4

## Test Coverage

### Controller Tests
- TestGetUpdateLocations_ReturnsHappy200ResponseShape: PASS
  - Validates 200 success response shape (message/data.timeline/request_id)
- TestGetUpdateLocations_EmptyTimelineReturnsNoData200: PASS
  - Validates 200 with message = "No Data" for empty timeline
- TestGetUpdateLocations_RecordNotFoundReturns404: PASS
  - Validates gorm.ErrRecordNotFound maps to 404
- TestGetUpdateLocations_MissingEmployeeIDReturns400: PASS
  - Validates binding error handling (400) with request_id in body

### Repository Tests (sqlmock)
- TestGetUpdateLocations_ReturnsTimelineForEmployee: PASS
  - Mock 3 timeline rows from union
  - Validates record scanned correctly

- TestGetEmployeeRole_ReturnsEmployeeCustIDWhenFound: PASS
  - Mock employee role resolver
  - Validates cust_id returned

- TestGetEmployeeRole_ReturnsNotFoundWhenEmployeeOutsideTenant: PASS
  - Mock empty result from resolver
  - Validates 404 behavior

## Lint & Build

✓ gofmt -l . → no output (all formatted)
✓ go vet ./... → no output (no issues)
✓ go build ./... → success
✓ go test ./... → all PASS

## Preserved Invariants

✓ Controller → Service → Repository → DB layering
✓ Service-level transactions pattern (no writes on this endpoint)
✓ Tenant scope self+child pattern (smc.m_customer)
✓ JWT middleware path (401 handled by existing middleware)
✓ snake_case JSON contract
✓ Parameterized SQL (? placeholders)
✓ No panic-recover swallow (service returns errors, controller maps)
✓ context.WithTimeout 30s in controller
✓ Existing LiveMonitoringData/Paging/DetailData untouched
✓ golang naming (snake_case files, CamelCase exports)
✓ No pagination (no LIMIT/OFFSET/count/page/limit params)
✓ Full-day timeline chronological order

## Decisions Applied

- D1: Role resolver via mst.m_employee + smc.m_customer ✓
- D2: No pagination ✓
- D3: Order recorded_at ASC (+ type/destination_id tie-breaker for determinism) ✓
- D4: Missing emp_id → 400 ✓
- D5: Empty date → default today Asia/Jakarta ✓
- D6: Resolver 404 generically (safe behavior) ✓
- D7: LENGTH(cust_id) > 6 for branch selector ✓
- D8: Distributor destination_type null by code (schema has field, not used) ✓
- D9: destination_name join to mst.m_outlet ✓

## Open Assumptions (Deferred)

- D6: 404 vs 200 empty when employee outside tenant — implemented as 404 (security-safe). FE expectation TBD if needed in follow-up.
- RecordedAt as string (RFC3339) in model, not time.Time — avoids sqlmock type mismatch; conversion happens in response mapping.

