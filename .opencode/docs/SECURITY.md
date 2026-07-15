# Security

- Never commit `.env`, backup dumps, credentials, tokens, or copied secret values.
- This repo already contains plaintext credentials in tracked infra files; do not duplicate, rotate, or surface them in summaries.
- Treat DB clone/restore scripts as sensitive operations and verify target database names before running them.
- Use extra caution for auth, tenant isolation, upload, and external-integration changes; route final signoff to `@quality-gate`.
