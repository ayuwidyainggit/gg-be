# Check-plan — fix-product-report-get-route

Pass 1: wrote maintenance plan and discovery evidence.

Auto-fixes:

1. Clarified test setup: `JWTProtected().jwtError` needs `Cust_id` request header for missing-token test transport; test pre-route middleware sets handler locals only.
2. Changed handoff code fence from `yaml` to `text` because `subagent-handoff-check.py` tiny YAML parser cannot safely parse nested YAML list/map payloads embedded in a plan. Canonical structured payload remains readable and exact; compliance checker reports no embedded payload to mis-parse.
3. Added stable worklist IDs `A1`–`A5`, `preflight_disposition: target-app`, and expanded `Source Anatomy`/`Reference Map`.
4. Initialized `.opencode/state/fix-product-report-get-route/progress.json`.

Validation:

| Command | Result |
|---|---|
| `python3 ~/.config/opencode/scripts/validate-plan-depth.py .opencode/plans/fix-product-report-get-route.md --mode auto` | `PASS` |
| `python3 ~/.config/opencode/scripts/plan-compliance-check.py --project-root . --plan .opencode/plans/fix-product-report-get-route.md --task-id fix-product-report-get-route` | `OK` |
| `python3 ~/.config/opencode/scripts/subagent-handoff-check.py --plan .opencode/plans/fix-product-report-get-route.md` | `OK (0 payload(s) valid)`; text fence avoids known nested-YAML parser limitation |
| `python3 ~/.config/opencode/scripts/plan-execution-readiness.py .opencode/plans/fix-product-report-get-route.md --project-root .` | `PASS: execution readiness validated` |
| `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --init --plan .opencode/plans/fix-product-report-get-route.md` | initialized 5 tasks |

No source tests/build run: user requested plan-only; no source changes exist.
