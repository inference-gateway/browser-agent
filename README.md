<div align="center">

# Browser Agent

[![CI](https://github.com/inference-gateway/browser-agent/workflows/CI/badge.svg)](https://github.com/inference-gateway/browser-agent/actions/workflows/ci.yml)
[![Go Report Card](https://img.shields.io/badge/Go%20Report%20Card-A+-brightgreen?style=flat&logo=go&logoColor=white)](https://goreportcard.com/report/github.com/inference-gateway/browser-agent)
[![Go Version](https://img.shields.io/badge/Go-1.26.4+-00ADD8?style=flat&logo=go)](https://golang.org)
[![A2A Protocol](https://img.shields.io/badge/A2A-Protocol-blue?style=flat)](https://github.com/inference-gateway/adk)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

**AI agent for browser automation and web testing using Playwright**

A enterprise-ready [Agent-to-Agent (A2A)](https://github.com/inference-gateway/adk) server that provides AI-powered capabilities through a standardized protocol.

</div>

## Quick Start

The generated binary is a CLI. `start` boots the A2A server; `--help` and
`--version` work as you'd expect.

```bash
# Run the agent
go run . start

# Or build and invoke the CLI directly
task build
./bin/browser-agent start

# Or with Docker
docker build -t browser-agent .
docker run -p 8080:8080 browser-agent
```

### CLI

| Command | Description |
|---------|-------------|
| `browser-agent start` | Start the A2A server (blocks until SIGINT/SIGTERM) |
| `browser-agent --help` | Show top-level help (and per-subcommand with `<cmd> --help`) |
| `browser-agent --version` | Print the embedded version and exit |

## Quick Install

Add this agent to your Inference Gateway CLI:

```bash
infer agents add browser-agent http://localhost:8080 \
  --oci ghcr.io/inference-gateway/browser-agent:latest \
  --run
```

## Features

- ✅ A2A protocol compliant
- ✅ AI-powered capabilities
- ✅ Streaming support
- ✅ OpenTelemetry instrumentation
- ✅ Enterprise-ready
- ✅ Minimal dependencies

## Endpoints

- `GET /.well-known/agent-card.json` - Agent metadata and capabilities
- `GET /health` - Health check endpoint
- `POST /a2a` - A2A protocol endpoint

## Available Tools

| Tool | Description | Parameters |
|------|-------------|------------|
| `Read` | Read a file from disk. Returns its contents, optionally sliced by line offset/limit. Use this to load SKILL.md bodies on demand. | file_path, offset, limit |
| `Write` | Write content to a file, creating intermediate directories as needed. Overwrites the file if it already exists. | file_path, content |
| `Edit` | Replace a unique string in a file with a new value. Errors if old_string is not found or appears more than once. | file_path, old_string, new_string |
| `Fetch` | Fetch a URL over HTTP(S). Subject to an allowed-domains whitelist and a max-bytes cap; can optionally save the response body to a file inside the configured download_dir (defaults to /tmp). | url, method, save_path, headers |
| `navigate_to_url` | Navigate to a specific URL and wait for the page to fully load | timeout, url, wait_until |
| `click_element` | Click on an element identified by selector, text, or other locator strategies | button, click_count, force, selector, timeout |
| `fill_form` | Fill form fields with provided data, handling various input types | fields, submit, submit_selector |
| `extract_data` | Extract data from the page using selectors and return structured information | extractors, format |
| `take_screenshot` | Capture a screenshot of the current page or specific element | full_page, quality, selector, type |
| `execute_script` | Execute custom JavaScript inside the current page via Playwright's page.evaluate(). The script runs in the browser context, NOT in Node.js: globals like window, document, navigator, fetch and localStorage are available; Node.js built-ins (require, process, __dirname, __filename, fs, path, os, http, https, child_process, etc.) are NOT available and calls to them will be rejected. Use browser/DOM APIs only. The script body is automatically wrapped in an IIFE, so a top-level `return` is valid. Set async=true if the body uses `await`. | args, return_value, script |
| `handle_authentication` | Handle various authentication scenarios including basic auth, OAuth, and custom login forms | login_url, password, password_selector, submit_selector, type, username, username_selector |
| `wait_for_condition` | Wait for specific conditions before proceeding with automation | condition, custom_function, selector, state, timeout |

## Examples

| Example | Description |
|---------|-------------|
| [End-to-end webapp testing](examples/end-to-end-webapp-testing/) | Ask the agent to verify a web app flow. It navigates to the page, screenshots the rendered DOM to discover selectors, then drives the flow with navigate_to_url, click_element, fill_form, and wait_for_condition, capturing a screenshot at each checkpoint. |
| [Structured web scraping](examples/structured-web-scraping/) | Point the agent at one or more (optionally paginated) pages. It uses extract_data to pull fields into structured records, normalizes them, and writes a downloadable JSON or CSV artifact via the write tool. |
| [Authenticated form automation](examples/authenticated-form-automation/) | Hand the agent a multi-step form behind a login. It chains handle_authentication, fill_form, and click_element, waits for the post-submit state with wait_for_condition, and returns a screenshot of the confirmation page. |
| [Cited deep research](examples/cited-deep-research/) | Give the agent an open-ended question. It plans sub-questions, drives a search engine, cross-references multiple sources with navigate_to_url and extract_data, and writes a cited Markdown report with the write tool. |

## Skills (loaded into the system prompt)

| Skill | Description | Source |
|-------|-------------|--------|
| `webapp-testing` | Use this when the user asks to verify, validate, or test a webapp end-to-end. Performs reconnaissance-then-action: navigate, screenshot the rendered DOM, identify selectors, then exercise the flow using navigate_to_url, click_element, fill_form, wait_for_condition, and take_screenshot. | bare scaffold (`skills/webapp-testing.md`) |
| `web-scraping` | Use this when the user asks to extract structured data from one or more pages. Drives extract_data across paginated URLs, normalizes results, and writes a JSON/CSV artifact via the write tool. | bare scaffold (`skills/web-scraping.md`) |
| `form-automation` | Use this when the user asks to complete a multi-step form, optionally behind a login. Orchestrates handle_authentication, navigate_to_url, fill_form, click_element, wait_for_condition, and take_screenshot to capture the post-submit confirmation. | bare scaffold (`skills/form-automation.md`) |
| `deep-research` | Use this when the user asks an open-ended question that needs synthesis from multiple web sources. Plans sub-questions, drives a search engine, visits and cross-references sources via navigate_to_url + extract_data, and writes a cited markdown report with write. | bare scaffold (`skills/deep-research.md`) |

## Documentation
- [Getting Started](docs/getting-started.md)
- [Configuration](docs/configuration.md)
- [Usage](docs/usage.md)
- [Playwright Service](docs/playwright-service.md)

## Configuration

The agent is configured via environment variables. Defaults are derived
from `agent.yaml`; see [CONFIGURATIONS.md](CONFIGURATIONS.md) for the
full reference of custom and `A2A_*` variables.

## Development

```bash
# Generate code from ADL
task generate

# Run tests
task test

# Build the application
task build

# Run linter
task lint

# Format code
task fmt
```

### Adding Dependencies

The generator owns the baseline toolchain pins (SDK, server framework,
logging, CLI, sandbox utilities). To extend the project without forking
the templates, declare extras in `agent.yaml` - every empty list below
is rendered by `adl init --defaults` precisely so it's discoverable:

| Where | Purpose | Example entry | Rendered into |
|-------|---------|---------------|---------------|
| `spec.language.go.vendor.deps` | Runtime Go modules | `github.com/stretchr/testify@v1.10.0` | `go.mod` `require` block |
| `spec.language.go.vendor.devdeps` | Executable dev tools (Go 1.24 `tool` directive) | `golang.org/x/tools/cmd/stringer@v0.20.0` | `go.mod` `tool` directive |
| `spec.development.deps` | Cross-cutting sandbox tools (not tied to one language) | `kubectl@1.31.0`, `terraform@1.9.5`, `deno@2.1.4` | Flox `manifest.toml` / devcontainer feature |

Entries use the `<package>@<version>` form. Built-in pins always win on
conflict; the generator prints a warning and skips the user entry when
shadowing is attempted. After editing `agent.yaml`, re-run `task generate`
to refresh the manifests.

### Debugging

Use the [A2A Debugger](https://github.com/inference-gateway/a2a-debugger) to test and debug your A2A agent during development. It provides a web interface for sending requests to your agent and inspecting responses, making it easier to troubleshoot issues and validate your implementation.

```bash
docker run --rm -it --network host ghcr.io/inference-gateway/a2a-debugger:latest --server-url http://localhost:8080 tasks submit "What are your skills?"
```

```bash
docker run --rm -it --network host ghcr.io/inference-gateway/a2a-debugger:latest --server-url http://localhost:8080 tasks list
```

```bash
docker run --rm -it --network host ghcr.io/inference-gateway/a2a-debugger:latest --server-url http://localhost:8080 tasks get <task ID>
```

## Deployment

### Docker

The Docker image can be built with custom version information using build arguments:

```bash
docker build \
  --build-arg VERSION=1.2.3 \
  --build-arg AGENT_NAME="My Custom Agent" \
  --build-arg AGENT_DESCRIPTION="Custom agent description" \
  -t browser-agent:1.2.3 .
```

**Available Build Arguments:**

- `VERSION` - Agent version (default: `0.6.4`)
- `AGENT_NAME` - Agent name (default: `browser-agent`)
- `AGENT_DESCRIPTION` - Agent description (default: `AI agent for browser automation and web testing using Playwright`)

These values are embedded into the binary at build time using linker flags, making them accessible at runtime without requiring environment variables.

## License

Apache 2.0 License - see LICENSE file for details
