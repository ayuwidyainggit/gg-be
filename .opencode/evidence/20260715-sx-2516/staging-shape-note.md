# Staging shape note

- Existing `import.sales_update_temp` remains header/raw staging table. Original 22 columns unchanged.
- B6 adds header snapshot and reconstruction metadata columns. `staging_key` links header row to detail rows.
- B7 adds `import.sales_update_temp_detail` for grouped `CreateOrderDetBody` snapshots.
- Worker reconstruction must select staged values only. It must not re-query mutable outlet, salesman, product, parent-product, unit, price, conversion, or warehouse master data.
- `cust_id` and `history_id` remain required header/detail correlation fields. Tenant filtering stays `cust_id` on header rows; detail rows are additionally constrained by `history_id` and `staging_key`.
- Existing `status_insert`/`error_message` behavior remains unchanged. Current plan says only valid rows are staged, so worker may rely on `status_insert='SUCCESS'` while still re-validating staged rows.
- Stable constants not stored in B6/B7 are documented in `B6-wide-alter.md`; these values come from current parser/service code and plan §12.

## Residual assumptions

1. Existing `sales_update_temp.id` is accepted as `staging_key`; B7 FK depends on that existing primary key.
2. `parent_cust_id`, `user_id`, and `created_by` widths follow current request/service types, not inferred DB master widths.
3. `CreateOrderDetBody` JSON-array remark fields are omitted because current parser never populates them and current Store path does not consume them for this import. Add JSONB columns only if later implementation proves Store path requires them.
4. Header snapshot columns are intentionally nullable to preserve existing raw-row compatibility. Wiring must populate them before worker processing.
5. B6 rollback drops only B6 columns. B7 rollback drops detail table. Rollback order must be B7 first, then B6, if ever needed.
6. Local DB was reachable at `localhost`, not `host.docker.internal`; first documented command failed DNS resolution. No remote DB accessed.
