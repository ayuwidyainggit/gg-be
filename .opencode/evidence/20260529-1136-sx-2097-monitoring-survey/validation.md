# SX-2097 Validation

## Commands run

```bash
# from pjp/
rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_(IncludesSurveyData|NoSurveyReturnsEmptyList)'
rtk go test ./repository/live_monitoring -run TestGetSubmittedSurveyData
rtk go test ./service/live_monitoring ./repository/live_monitoring
rtk go test ./...
```

## Results
- `rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_(IncludesSurveyData|NoSurveyReturnsEmptyList)'` -> pass (`2 passed in 1 packages`)
- `rtk go test ./repository/live_monitoring -run TestGetSubmittedSurveyData` -> pass (`1 passed in 1 packages`)
- `rtk go test ./service/live_monitoring ./repository/live_monitoring` -> pass (`45 passed in 2 packages`)
- `rtk go test ./...` -> pass (`56 passed in 48 packages`)
- After tenant-join hardening adjustment, reran:
  - `rtk go test ./repository/live_monitoring -run TestGetSubmittedSurveyData` -> pass
  - `rtk go test ./service/live_monitoring ./repository/live_monitoring` -> pass
  - `rtk go test ./...` -> pass

## Docker/API checks

```bash
# from repo root
rtk docker compose -f docker-compose.yml ps pjp
curl -sS 'http://localhost:9010/api/v1/health'
```

- Compose command returns warning only about obsolete `version` attribute in `docker-compose.yml`.
- Health endpoint returns: `{"code":200,"message":"OK"}`.

## DB checks

```bash
PGPASSWORD="postgres" psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c 'SELECT current_database(), current_user, inet_server_port();'
PGPASSWORD="postgres" psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "SELECT table_schema, table_name, column_name, data_type FROM information_schema.columns WHERE table_schema = 'mst' AND table_name IN ('survey_answer', 'm_survey', 'm_outlet') AND column_name IN ('survey_answer_id', 'survey_id', 'outlet_id', 'emp_id', 'answer_date', 'status', 'survey_title', 'outlet_code', 'outlet_name', 'cust_id') ORDER BY table_name, ordinal_position;"
PGPASSWORD="postgres" psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "SELECT sa.emp_id, sa.answer_date::date AS answer_date, sa.status, COUNT(*) AS total FROM mst.survey_answer sa WHERE sa.emp_id = 210 AND sa.answer_date::date = '2026-05-28' GROUP BY sa.emp_id, sa.answer_date::date, sa.status ORDER BY sa.status;"
PGPASSWORD="postgres" psql -h localhost -p 5432 -U postgres -d ggn_scyllax -c "SELECT COUNT(sa.survey_answer_id) AS submission, ms.survey_title, mo.outlet_code, mo.outlet_name FROM mst.survey_answer sa JOIN mst.m_survey ms ON ms.survey_id = sa.survey_id JOIN mst.m_outlet mo ON mo.outlet_id = sa.outlet_id WHERE sa.answer_date::date = '2026-05-28' AND sa.emp_id = 210 AND sa.status = 'Submitted' GROUP BY ms.survey_title, mo.outlet_code, mo.outlet_name ORDER BY ms.survey_title ASC, mo.outlet_code ASC;"
```

- Local DB access verified: `ggn_scyllax` as user `postgres` on port `5432`.
- Schema verified:
  - `mst.survey_answer` has `cust_id`, `survey_answer_id`, `survey_id`, `emp_id`, `outlet_id`, `answer_date`, `status`.
  - `mst.m_survey` has `cust_id`, `survey_id`, `survey_title`.
  - `mst.m_outlet` has `cust_id`, `outlet_id`, `outlet_code`, `outlet_name`.
- Source data for `emp_id=210`, `date=2026-05-28`: one row with status `Submitted`.
- Runtime tenant check on sampled row:
  - `sa.cust_id = C220010001`
  - `ms.cust_id = C22001`
  - `mo.cust_id = C220010001`
- Expected aggregation query result:
  - `submission=1`
  - `survey_title=Survey (with image) - Optional`
  - `outlet_code=PRA260055`
  - `outlet_name=Toko Depok macet`

## Monitoring detail curl validation

Authenticated local service validation executed with generated local JWT using repo-local secret, without persisting secret/token in artifacts.

Result snippet:

```json
{
  "message": "Success",
  "survey_data": [
    {
      "submission": 1,
      "survey_title": "Survey (with image) - Optional",
      "outlet_code": "PRA260055",
      "outlet_name": "Toko Depok macet"
    }
  ]
}
```

- Endpoint validated: `GET http://localhost:9010/api/v1/monitoring_locations/details?emp_id=210&date=2026-05-28&distributor_id=67`
- `survey_data` exists and is JSON array.
- API response matches local DB aggregation for sampled data.

## Tenant filter decision evidence
- Implemented tenant filter in survey query with `sa.cust_id IN ?` using `targetCustIDs` derived from salesman cust_id in service.
- Runtime schema validation confirmed `mst.survey_answer.cust_id` exists, so tenant filter is valid and safer than unscoped query.
- Final join strategy:
  - `mst.m_outlet` joined by `outlet_id` **and** `cust_id` for tenant-safe outlet mapping.
  - `mst.m_survey` joined by `survey_id` only, because sampled runtime data shows survey master row stored on parent cust (`ms.cust_id = C22001`) while answer row uses child cust (`sa.cust_id = C220010001`). Adding `ms.cust_id = sa.cust_id` incorrectly dropped valid data.
