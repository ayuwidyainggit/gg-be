# OpenCode Docs Index

This directory is the repo-local system of record for agent workflow, validation, and risk handling in `scylla-be`.

- `AGENT_ROUTING.md` — who owns discovery, planning, implementation, review, and signoff
- `ARCHITECTURE.md` — monorepo shape, service boundaries, transactions, multi-tenant rules
- `SERVICE_MATRIX.md` — service/module matrix + README authority audit status
- `QUALITY.md` — validation commands, evidence expectations, done criteria
- `EVALS.md` — replayability and task-evidence expectations
- `SECURITY.md` — secret handling, DB safety, sensitive change posture
- `PROMPT_GATES.md` — assumptions, ambiguity, and stop/continue gates
- `SKILLS.md` — practical role-and-skill ownership notes
- `MCP.md` — MCP/tool usage posture for this repo
- `GOLDEN_PRINCIPLES.md` — concise operating principles
- `AGENT_LEGIBILITY.md` — how to keep ownership clear
- `DECISIONS.md` — durable decisions worth preserving
- `RELEASE.md` — release and runtime-readiness checks
- `QUALITY_SCORE.md` — quality dimensions for work review
- `GC_WORKFLOW.md` — guardrail and conformance workflow
- `PROJECT_STACK.md` — detected multi-module Go stack, version/toolchain, runtime, and staleness notes
- `PROJECT_COMMANDS.md` — safe/default, validation, migration, approval, and destructive command boundaries
- `FRAMEWORK_PLAYBOOK.md` — Go/Fiber/Gin/sqlx/GORM/golang-migrate generator-first and manual-fallback guidance
- `PROJECT_DETECTED_TOOLS.md` — detected frameworks, generators, test tools, runtime tooling, and explicit gaps
- `TOOL_USAGE.md` — tool/authority order and MCP routing for this repo
- `AGENT_TOOL_ACCESS.md` — role/tool boundary matrix for this repo

Read `AGENT_ROUTING.md`, `ARCHITECTURE.md`, and `QUALITY.md` first for most tasks.

Quick usage note:
- `SERVICE_MATRIX.md` is the canonical service/module matrix (compose presence, ports, env/Makefile, migration style) and README authority audit.
- `ARCHITECTURE.md` keeps cross-cutting architecture rules and links to `SERVICE_MATRIX.md` for module inventory details.
- `QUALITY.md` now includes per-module validation focus and reliability caveats for stale/inconsistent READMEs.
