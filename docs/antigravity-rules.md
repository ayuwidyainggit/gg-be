# Antigravity Rules - Scylla Backend

These rules define specific behavioral requirements for Antigravity when working on this project.

## 1. Docker Service Verification

**BEFORE** performing any API testing or curl request:
1. Run `docker ps` to check if required containers are running
2. Verify health status of the target service (e.g., `curl localhost:[PORT]/ping`)
3. If containers are not running, notify the user and suggest running `./docker/start-dev.sh`

## 2. API Testing with CURL

**AFTER** implementing or modifying any API endpoint:
1. Ensure Docker container is UP (see Rule #1)
2. Perform `curl` request to test the endpoint functionality
3. Include request and response in verification summary

Example:
```bash
# Check service health
curl -s http://localhost:9001/ping

# Test endpoint
curl -X POST http://localhost:9001/api/v1/resource \
  -H "Content-Type: application/json" \
  -d '{"field": "value"}'
```

## 3. Database Verification

**AFTER** any CREATE or UPDATE operation:
1. Query the database directly to verify data persistence
2. Use `docker exec` with PostgreSQL client

Example:
```bash
# Connect to database and verify data
docker exec scylla-postgres-citus psql -U postgres -d scylla_citus_dev -c \
  "SELECT id, column_name FROM table_name WHERE id = 'xxx' LIMIT 1;"
```

## 4. Service Endpoints Reference

| Service    | Port | Health Check           |
|------------|------|------------------------|
| System     | 9001 | localhost:9001/ping    |
| Master     | 9002 | localhost:9002/ping    |
| Inventory  | 9003 | localhost:9003/ping    |
| Sales      | 9004 | localhost:9004/ping    |
| Finance    | 9005 | localhost:9005/ping    |
| TMS        | 9006 | localhost:9006/ping    |
| Mobile     | 9008 | localhost:9008/ping    |
| PJP        | 9009 | localhost:9009/ping    |
| Cronjob    | 9100 | localhost:9100/ping    |

## 5. Database Connection

```
Host: scylla-postgres-citus (container) or localhost:54321 (host)
User: postgres
Database: scylla_citus_dev
```
