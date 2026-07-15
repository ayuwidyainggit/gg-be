## check-plan auto-fixes (2026-07-15T09:56Z)

- Appended `evidence_refresh` token to `## Progress Tracking` **Update rules** list (explicit status transition rule for evidence rewrites).
- Appended an explicit `update_rules` and `task_map` YAML block inside `## Progress Tracking`, mapping every worklist id (A1..K3) to its owner and exact `task-progress.py --update <id> ...` command derived from the Execution-ready Worklist.

All edits append-only; no existing plan content rewritten or removed.
