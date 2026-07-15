# Agent Routing

- Default entrypoint: `@orchestrator`.
- Unknown scope, repo discovery, pattern search, or cross-service mapping → `@explorer`.
- Multi-step or evidence-heavy plan → `@artifact-planner`.
- Bounded code changes, tests, fixtures, or refactors in one service/module → `@fixer`.
- UI, accessibility, visual parity, or design-system work → `@designer` first.
- Architecture tradeoffs, data-model boundaries, multi-tenant/runtime concerns → `@architect` or `@oracle`.
- External docs or version-sensitive library behavior → `@librarian`.
- Material changes, security-sensitive work, prompt/config updates, or final signoff → `@quality-gate`.

Do not default multi-file implementation or broad discovery back to `@orchestrator`.
