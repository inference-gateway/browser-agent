# Security

The agent's A2A endpoint has **no authentication by default** and the
shipped `docker-compose.yaml` deliberately does not publish its port. If
you expose port 8080 beyond a trusted network, enable authentication
first: set `A2A_AUTH_ENABLED=true` together with the `A2A_AUTH_ISSUER_URL`,
`A2A_AUTH_CLIENT_ID`, and `A2A_AUTH_CLIENT_SECRET` OIDC variables.

Built-in tools are sandboxed by default:

- `read`/`write`/`edit` are restricted to `/tmp/playwright/artifacts`
  (read additionally covers `.agents/skills`). Widen per deployment via
  `TOOLS_READ_ALLOWED_ROOTS`, `TOOLS_WRITE_ALLOWED_ROOTS`, and
  `TOOLS_EDIT_ALLOWED_ROOTS`.
- `fetch` is capped at 10 MiB / 30 s per request; restrict reachable
  hosts with `TOOLS_FETCH_ALLOWED_DOMAINS`.
- `navigate_to_url` refuses loopback, private (RFC1918/ULA), and
  link-local addresses (including cloud metadata endpoints). Set
  `BROWSER_ALLOW_INTERNAL_URLS=true` to test internal webapps.

To report a vulnerability, see [SECURITY.md](../SECURITY.md).
