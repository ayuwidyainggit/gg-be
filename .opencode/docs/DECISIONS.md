# Decisions

- 2026-05-13: Initialized `.opencode/docs/` as the repo-local system of record for harness workflow.
- 2026-05-13: Kept project-local `rtk` command posture because the existing repo workflow and runtime guidance depend on it.
- 2026-05-13: Preserved service-layer transaction, multi-tenant, and migration-path rules as first-class repo policy.
- 2026-05-13: Distinguished compose-managed default services from extra repo modules such as `pjp-principle` and `pjp-sales`.
- 2026-05-13: Added repo-evidence service/module matrix and documented env/port/Makefile/migration differences to reduce ambiguity from stale template READMEs.
