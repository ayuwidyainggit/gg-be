# Validation Evidence — SX-1906 Survey Principal-only BU

## Red

Added `TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman` before production fix.

Command:

```bash
go test ./service -run TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman
```

Result before fix:

```text
--- FAIL: TestSurveyService_Store_ShouldCreatePrincipalOnlySurveyWithSelectedAreasAndSalesman (0.00s)
    survey_service_test.go:520: expected principal-only survey to be created, got distributor_id is required when area_id is provided
FAIL
FAIL	master/service	0.577s
FAIL
```

## Green / Regression

Command:

```bash
go test ./service -run 'TestSurveyService_Store|TestSurveyService_Update'
```

Result after fix:

```text
ok  	master/service	(cached)
```

Full module command:

```bash
go test ./...
```

Result after fix:

```text
?   	master	[no test files]
ok  	master/adapter	0.456s
ok  	master/controller	0.542s
ok  	master/entity	0.978s
?   	master/model	[no test files]
?   	master/pkg/config	[no test files]
?   	master/pkg/config/env	[no test files]
?   	master/pkg/constant	[no test files]
?   	master/pkg/conversion	[no test files]
?   	master/pkg/errmsg	[no test files]
?   	master/pkg/generator	[no test files]
?   	master/pkg/jwthelper	[no test files]
?   	master/pkg/middleware	[no test files]
?   	master/pkg/rabbitmq	[no test files]
?   	master/pkg/responsebuild	[no test files]
?   	master/pkg/server	[no test files]
?   	master/pkg/sql_helper	[no test files]
?   	master/pkg/str	[no test files]
?   	master/pkg/structs	[no test files]
ok  	master/pkg/texttranslator	2.002s
?   	master/pkg/validation	[no test files]
ok  	master/repository	0.587s
ok  	master/service	0.807s
```

## Notes

- `go test` commands were run from `master/` because this repo is multi-module and has no root `go.mod`.
- `git status` from `scylla-be` and parent `Geekgarden` returned `fatal: not a git repository`, so no Git diff/status evidence is available in this workspace.
- Authenticated API retest was later run against `https://best.scyllax.online/master/v1/survey` with the SX-1906 principal-only payload. The response was still a validation failure from the deployed environment:

```json
{"message":"distributor_id diperlukan ketika area_id disediakan","request_id":"69f7fd1e83167fc901dbfea8"}
```

- Database validation after that failed API call found no created survey row for `cust_id = 'C22001'` and `survey_title = 'Testing Survey Principal'`.
- Interpretation: local code/tests are fixed, but the remote/staging API endpoint has not picked up this workspace change yet, or it is serving a different build than the modified local source.

## Localhost Docker API retest

Local health check:

```bash
curl -i http://localhost:9002/ping
```

Result:

```text
HTTP/1.1 200 OK
It works
```

Initial local URL check with production gateway prefix:

```bash
curl -i http://localhost:9002/master/v1/survey ...
```

Result:

```text
HTTP/1.1 404 Not Found
Cannot POST /master/v1/survey
```

Local service registers routes without the `/master` gateway prefix, so retest used `http://localhost:9002/v1/survey`.

The provided staging token was signed with the local `JWT_SECRET_KEY=secret`, so the same token was valid locally.

Command used a unique title `Testing Survey Principal Local` to avoid title overlap with prior attempts.

Result:

```text
HTTP/1.1 201 Created
{"message":"Survei telah berhasil dibuat","request_id":"69f8021b697fe93aee03b4de"}
```

Database validation for the locally created survey:

```text
 row_type | survey_id | cust_id |          survey_title          | target_type | target_id |             extra             
----------+-----------+---------+--------------------------------+-------------+-----------+-------------------------------
 area     |        98 |         |                                |             | 70        | 0
 area     |        98 |         |                                |             | 82        | 0
 area     |        98 |         |                                |             | 83        | 0
 area     |        98 |         |                                |             | 84        | 0
 area     |        98 |         |                                |             | 85        | 0
 area     |        98 |         |                                |             | 86        | 0
 area     |        98 |         |                                |             | 89        | 0
 salesman |        98 | C22001  |                                |             | 369       | 
 survey   |        98 | C22001  | Testing Survey Principal Local | Specific    | 369       | 2026-05-04 02:19:07.546699+00
(9 rows)
```

Interpretation: localhost Docker using the modified source successfully creates the principal-only survey, persists all selected areas with sentinel `distributor_id = 0`, and persists salesman `369` scoped to principal cust `C22001`.
