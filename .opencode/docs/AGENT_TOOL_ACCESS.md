# Agent Tool Access

This document is the repo-local role/tool matrix. Project rules in `AGENTS.md` and `.opencode/docs/ARCHITECTURE.md` override generic agent habits.

| Lane | Primary responsibility | May modify | Must not own |
|---|---|---|---|
| `@orchestrator` | route, integrate, plan intake, evidence, execution tracking | tiny reversible edits; `.opencode` runtime/evidence docs | broad multi-file implementation, final self-signoff |
| `@explorer` | read-only service/module/test/pattern discovery | none | source edits, plan status claims |
| `@artifact-planner` | durable plans/evidence under `.opencode/plans`, `.opencode/draft`, `.opencode/evidence` | planning artifacts only | implementation, framework/runtime config |
| `@backend` | Go API/service/repository/migration implementation | target module files within plan boundary | UI direction, final signoff, unrelated modules |
| `@fixer` | bounded implementation/tests/refactor | target files/tests in approved scope | architecture decisions, final signoff |
| `@designer` | UI/API-consumer surface direction, accessibility/review | UI work only when explicitly routed | default backend API changes |
| `@librarian` | official/current docs and document extraction | none | source edits |
| `@architect` / `@oracle` | architecture and tradeoff review | none | implementation and final signoff |
| `@quality-gate` | read-only final conformance/risk review | evidence/review note only | source edits |
| `@devops` | compose/CI/deploy/env/monitoring boundaries | approved infra/config scope | service business logic |

## Tool posture

- Shell: normal project commands are `rtk`-prefixed. Never echo secrets, tokens, `.env` values, or database credentials.
- File inspection: use `read`, `grep`, and `glob` first. Use source edits only in owning implementation lane.
- Database: no ad-hoc production/staging writes. Local target only unless user explicitly authorizes otherwise.
- Browser: backend-only repository; BrowserOS is not default. Use only for a documented browser/UI flow.
- Context7/librarian: required for version-sensitive API or CLI questions not settled by repo docs.
- Semgrep: use for SQL injection, auth, secret, or trust-boundary changes when available.

## Handoff requirements

Every non-trivial worker handoff contains: task ID, plan ID, caller/callee, exact scope, claim boundary, source basis, preserve/do-not-touch list, validation, exit criteria, evidence path, dependencies, and compact context bundle. Validate payload with `subagent-handoff-check.py` where a plan is execution-bound.

## Current active contract

- Default first lane: `@orchestrator`.
- New Go code: `@backend` after plan intake.
- Multi-file tests/refactor without API ownership ambiguity: `@fixer`.
- Final material claim: `@quality-gate`.
