# Usage

Once the agent is running (see [Getting Started](getting-started.md)), you drive
it by submitting A2A tasks in natural language. The agent plans the work,
selects tools, and streams progress back.

## Submitting a task

The quickest way to exercise the agent is the
[A2A Debugger](https://github.com/inference-gateway/a2a-debugger):

```bash
docker run --rm -it --network host \
  ghcr.io/inference-gateway/a2a-debugger:latest \
  --server-url http://localhost:8080 \
  tasks submit "Go to https://example.com and take a full-page screenshot"

# List tasks and inspect a result
docker run --rm -it --network host \
  ghcr.io/inference-gateway/a2a-debugger:latest \
  --server-url http://localhost:8080 tasks list
docker run --rm -it --network host \
  ghcr.io/inference-gateway/a2a-debugger:latest \
  --server-url http://localhost:8080 tasks get <task-id>
```

You can also submit directly over the A2A protocol at `POST /a2a`, or wire the
agent into the `infer` CLI (see the `cli` profile in `docker-compose.yaml`).

## What the agent can do

The agent exposes browser-automation tools and a set of skill playbooks it loads
on demand.

### Tools

| Tool | Purpose |
|------|---------|
| `navigate_to_url` | Load a URL and wait for the page |
| `click_element` | Click by CSS selector, XPath, or text |
| `fill_form` | Fill and optionally submit form fields |
| `extract_data` | Pull structured data out of the DOM |
| `take_screenshot` | Capture the page or a single element |
| `execute_script` | Run JavaScript in the page (browser context only) |
| `handle_authentication` | Basic auth, form login, or OAuth flows |
| `wait_for_condition` | Wait for a selector, navigation, or custom predicate |
| `read` / `write` / `edit` | Work with local files (e.g. save artifacts) |
| `fetch` | Fast HTTP(S) GET/HEAD without a browser session |

> **Tip:** For static content (raw files, JSON/XML APIs, feeds, downloads) the
> agent prefers `fetch` over `navigate_to_url` — it is much faster and opens no
> browser session. Browser tools are reserved for JavaScript-rendered pages and
> stateful, authenticated sessions.

### Skills

| Skill | When it triggers |
|-------|------------------|
| `webapp-testing` | Verify/validate/test a webapp end-to-end |
| `web-scraping` | Extract structured data across one or more pages |
| `form-automation` | Complete a multi-step form, optionally behind a login |
| `deep-research` | Synthesize an answer from multiple web sources |

## Example prompts

- "Test the login flow on `https://staging.example.com` with user `demo` and
  screenshot the dashboard after signing in."
- "Scrape every product title and price from the three catalog pages and save
  them as a CSV."
- "Fill the multi-step signup form at `<url>`, submit it, and capture the
  confirmation page."
- "Research the current state of the WebGPU spec across the official sources and
  write a cited Markdown report."

Screenshots and extracted data are returned as downloadable artifacts when the
artifacts server is enabled (`A2A_ARTIFACTS_ENABLE=true`).
