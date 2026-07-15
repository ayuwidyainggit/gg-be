# Runtime parent-null defect — SX-2513

- Date: 2026-07-14
- Claim level: `confirmed_runtime`
- Scope: plan amendment only; no source, compose, env, migration, or data change.

## Confirmed runtime report

Production `best.scyllax.online` returned for `cust_id=C260020001`:

```text
gagal memecahkan kode: skema: kesalahan mengonversi nilai untuk pro_id
```

Production fixture contains `mapping_enabled` row. Its parent is missing or fails active/undeleted eligibility. Existing normalized `CASE` selected `parent.pro_id`; `LEFT JOIN` therefore produced `NULL`, which failed scan into non-nullable `pro_id`.

Local `psql` did not reproduce production fixture. For `C260020001`, local `mapping_enabled parent_present` check returned two rows, both parent present and flag false. This is not evidence against production branch.

## Required remediation

When row is mapping-enabled but eligible parent is absent:

- primary `cust_id`, `pro_id`, `pro_code`, `pro_name`, and `parent_pro_id` use `mp` fields;
- `original_cust_id`, `original_pro_id`, `original_pro_code`, and `original_parent_pro_id` remain populated from `mp`;
- `type='Product Mapping'` remains;
- no output primary scan field may be NULL;
- eligible-parent branch remains normalized to parent fields.

## Validation proof blocks

1. SQLMock: mapping-enabled + eligible parent returns parent primary fields and `mp` original fields.
2. SQLMock: mapping-enabled + missing/inactive parent returns `mp` primary and `mp` original fields; scans without conversion error.
3. Production runtime: authorized read-only curl for `C260020001`, capture status/body redacted of token, verify HTTP 200, no conversion error, non-null primary fields, preserved `original_*` fields.

Bearer token rotation reminder exists in runtime operational context; token is never pasted into this evidence or plan.
