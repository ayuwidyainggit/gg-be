# A4 SX-2508 validation

## Scope
Local sales build/test validation and classification of two DetailV2 failures. No source edits made.

## Commands

| Command | Exit | Result |
|---|---:|---|
| `rtk go build .` from `sales` | 0 | Passed |
| `rtk go test ./... 2>&1 \| tee /tmp/sx2508-full.log` from `sales` | 1 | 359 passed, 2 failed, 1 skipped across 22 packages |
| `python3 ~/.config/opencode/scripts/pre-gate-smoke-check.py --project-root .` | 0 | Passed; zero-byte assets 0, manifest mismatches 0, empty surfaces 0 |
| `python3 ~/.config/opencode/scripts/runtime-verify.py --help` | 0 | Available; help only, no live runtime call |

Full redacted test capture: `/tmp/sx2508-full.log`. RTK detailed output may exist under local RTK application support logs; no secrets copied here.

## Failure classification

### `TestDetailV2_PostRolloutWithoutSnapshot_UsesConsultV2ByTab`

- Fails at `sales/service/order_service_test.go:1682`.
- Exact failure: `post-rollout rows without snapshot must expose promo remarks from v2 consult`.
- Test does exercise DetailV2 quantity fields (`Qty=52`, `Qty1=2`, `Qty2=0`, `Qty3=1`, conversion units 10/5), but failure occurs after V2 consult: expected promo remarks are not populated. It is not a Qty1/Qty2/Qty3 assertion.
- Runtime log shows V2 consult path executes (`promotion_service.go:2290`, `2315`, `2318`, `2338`, `2344`).
- Classifier: pre-existing/unrelated to A3 remap. A3 changes DetailV2 mapped-row quantity remap at `order_service.go:3079-3095`; this failure is promo response/remarks application after consult. No minimum source fix made.

### `TestDetailV2_PreRolloutWithoutSnapshot_MustConsultV2ByTab`

- Fails at `sales/service/order_service_test.go:6172`.
- Exact failure: `detail without snapshot must inject runtime promo, got { ... Promo1:0 ... }`.
- Test exercises DetailV2 quantity fields and conversion units. `consultCalled` passes its prior assertion, so V2 consult is reached; failure is runtime promo injection/application, not consult dispatch.
- Captured V2 log shows conversion and criteria validation executes, including `Phase 1: Initial Quantity Conversion`, `Phase 1: After Quantity Conversion`, and product criteria validation; test's non-mandatory criteria path does not produce a validated promo response for injection.
- Classifier: pre-existing/unrelated to A3 remap. A3 remap only changes mapped-row display conversion; this failure is V2 promotion qualification/application. No minimum source fix made.

Focused rerun of both tests reproduced both failures. Direct `go test` output also reproduced them; `rtk go test` summarized them as 2 failures.

## Contract and risk review

- Behavior changed: none.
- Quantity contract preserved: no source change; Large/Middle/Small and global converter untouched.
- Migration/cache/dependency changes: none.
- Tenant/transaction safety: no source change.
- Runtime verification: not run against live service; fresh `SCYLLAX_TOKEN` unavailable/was not requested, and handoff forbids live call without it.
- Residual risk: two DetailV2 promotion tests remain red and need a separate promotion-path investigation; A4 cannot claim full-module green.

## Status

A4 validation complete with blocker classification: build and smoke checks pass; full module remains red only on the two listed promotion tests, classified unrelated/pre-existing to A3. Evidence path: `.opencode/evidence/20260713-1400-sx-2508-import-order-values/A4-validation.md`.
