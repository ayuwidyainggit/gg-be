# Release

- Confirm impacted service tests pass from the service directory.
- Check compose/runtime assumptions when the change affects startup, env, ports, or integration behavior.
- Call out migration requirements explicitly before release.
- Escalate risky runtime, security, or data changes to `@quality-gate`.
