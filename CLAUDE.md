# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

browser-agent is an A2A (Agent-to-Agent) server implementing the [A2A Protocol](https://github.com/inference-gateway/adk) for agent-to-agent communication. AI agent for browser automation and web testing using Playwright. The project is automatically generated from ADL (Agent Definition Language) specifications defined in `agent.yaml`.

## Core Architecture

### ADL-Generated Structure

The codebase is generated using ADL CLI 0.39.1 and follows a strict generation pattern:
- **Generated Files**: Marked with `DO NOT EDIT` headers - manual changes will be overwritten
- **Configuration Source**: `agent.yaml` - defines agent capabilities, skills, and metadata
- **Server Implementation**: Built on the ADK (Agent Development Kit) framework from `github.com/inference-gateway/adk`

### Key Components

- **Main Entry Point**: `main.go` - A cobra-based CLI. The root command exposes
  `--version` and `--help`; the `start` subcommand boots the A2A server with:
  - OpenAI-compatible LLM client configuration
  - Agent builder with system prompt from `agent.yaml`
  - A2A server with streaming and background task handlers
  - Graceful shutdown handling

- **Agent Configuration**: `.well-known/agent-card.json` - Serves agent metadata at runtime
- **Environment Configuration**: Extensive env vars with `A2A_` prefix (see README for full list)

## Development Commands

```bash
# Generate/regenerate code from ADL specification
task generate

# Run the agent in development mode (debug enabled, port 8080)
task run

# Run tests (note: no tests currently exist)
task test
task test:cover  # with coverage

# Code quality
task lint         # Run golangci-lint
task fmt          # Format code with go fmt

# Build
task build        # Creates bin/browser-agent
task docker:build # Build Docker image

# Clean build artifacts
task clean
```

## Testing Individual Components

```bash
# Run specific test file (when tests are added)
go test -v ./path/to/package -run TestFunctionName

# Debug with A2A Debugger
docker run --rm -it --network host ghcr.io/inference-gateway/a2a-debugger:latest \
  --server-url http://localhost:8080 tasks submit "Your query"
```

## LLM Provider Configuration

The agent uses OpenAI-compatible LLM client. Configure with:
- `A2A_AGENT_CLIENT_PROVIDER`: `openai`, `anthropic`, `azure`, `ollama`, `deepseek`
- `A2A_AGENT_CLIENT_MODEL`: Model identifier
- `A2A_AGENT_CLIENT_API_KEY`: Provider API key
- `A2A_AGENT_CLIENT_BASE_URL`: Custom endpoint (optional)

## Adding New Functionality

### Tools (function-call)
The following tools are currently defined:
- **Read** (built-in): Read a file from disk. Returns its contents, optionally sliced by line offset/limit. Use this to load SKILL.md bodies on demand.
- **Write** (built-in): Write content to a file, creating intermediate directories as needed. Overwrites the file if it already exists.
- **Edit** (built-in): Replace a unique string in a file with a new value. Errors if old_string is not found or appears more than once.
- **Fetch** (built-in): Fetch a URL over HTTP(S). Subject to an allowed-domains whitelist and a max-bytes cap; can optionally save the response body to a file inside the configured download_dir (defaults to /tmp).
- **navigate_to_url**: Navigate to a specific URL and wait for the page to fully load
- **click_element**: Click on an element identified by selector, text, or other locator strategies
- **fill_form**: Fill form fields with provided data, handling various input types
- **extract_data**: Extract data from the page using selectors and return structured information
- **take_screenshot**: Capture a screenshot of the current page or specific element
- **execute_script**: Execute custom JavaScript inside the current page via Playwright's page.evaluate(). The script runs in the browser context, NOT in Node.js: globals like window, document, navigator, fetch and localStorage are available; Node.js built-ins (require, process, __dirname, __filename, fs, path, os, http, https, child_process, etc.) are NOT available and calls to them will be rejected. Use browser/DOM APIs only. The script body is automatically wrapped in an IIFE, so a top-level `return` is valid. Set async=true if the body uses `await`.
- **handle_authentication**: Handle various authentication scenarios including basic auth, OAuth, and custom login forms
- **wait_for_condition**: Wait for specific conditions before proceeding with automation

To modify tools:
1. Update `agent.yaml` `spec.tools` with tool definitions
2. Run `task generate` to regenerate the codebase
3. Implement tool logic in the generated `tools/` files (look for TODO placeholders)
4. Write tests for each tool

### Skills (markdown system-prompt playbooks)
The following skills are currently shipped with the agent:
- **webapp-testing** (bare scaffold): Use this when the user asks to verify, validate, or test a webapp end-to-end. Performs reconnaissance-then-action: navigate, screenshot the rendered DOM, identify selectors, then exercise the flow using navigate_to_url, click_element, fill_form, wait_for_condition, and take_screenshot.
- **web-scraping** (bare scaffold): Use this when the user asks to extract structured data from one or more pages. Drives extract_data across paginated URLs, normalizes results, and writes a JSON/CSV artifact via the write tool.
- **form-automation** (bare scaffold): Use this when the user asks to complete a multi-step form, optionally behind a login. Orchestrates handle_authentication, navigate_to_url, fill_form, click_element, wait_for_condition, and take_screenshot to capture the post-submit confirmation.
- **deep-research** (bare scaffold): Use this when the user asks an open-ended question that needs synthesis from multiple web sources. Plans sub-questions, drives a search engine, visits and cross-references sources via navigate_to_url + extract_data, and writes a cited markdown report with write.

Each skill lives in its own directory at `skills/<id>/SKILL.md` and is
loaded into the system prompt at startup. Bare skills can ship arbitrary
bundled assets (scripts, templates, resources) alongside `SKILL.md` -
the whole `skills/<id>/` directory is protected by `.adl-ignore` against
regeneration overwrites. To modify skills:
1. Update `agent.yaml` `spec.skills` with skill definitions
2. Run `task generate` (registry skills are re-fetched; bare skill
   directories are preserved when listed in `.adl-ignore`)
3. For bare skills, edit `skills/<id>/SKILL.md` directly - frontmatter
   (`name`/`description`/`tags`) shows up on the agent card. Drop helper
   scripts or templates next to it (e.g. `skills/<id>/scripts/foo.py`).

### Modifying Agent Behavior

- **System Prompt**: Edit in `agent.yaml`, then regenerate
- **Capabilities**: Modify in `agent.yaml` (streaming, pushNotifications, stateTransitionHistory)
- **Server Configuration**: Update environment variables or `agent.yaml` server section

## Testing Strategy

When implementing tests:
- Create `*_test.go` files alongside implementation files
- Use table-driven tests for comprehensive coverage
- Mock external dependencies (LLM client, Redis if used)
- Test A2A protocol compliance with integration tests

## Environment Management

### Development Environment
- **Flox Environment**: ✅ Configured via `.flox/env/manifest.toml` providing Go 1.26.2, linter, `go-task`, Docker, and the Claude Code CLI. Activate with `flox activate`.
- **Docker Compose**: ✅ Local service stack defined in `docker-compose.yaml`. Bring up the Inference Gateway and the agent (built from the local `Dockerfile`) with `docker compose up --build`. Opt-in profiles add the `infer` CLI (`docker compose --profile cli run --rm cli`) and the `a2a-debugger` (`docker compose --profile debugger run --rm debugger --server-url http://browser-agent:8080 tasks list`).

## Important Constraints

- **Generated Files**: Never manually edit files with "DO NOT EDIT" headers
- **Configuration Changes**: Always modify `agent.yaml` and regenerate
- **ADL Version**: Ensure ADL CLI 0.39.1 or compatible version for regeneration
- **Port Configuration**: Default 8080, configurable via `A2A_PORT` or `A2A_SERVER_PORT`

## Debugging Tips

- Enable debug mode: `A2A_DEBUG=true`
- Check health: `GET /health`
- View agent metadata: `GET /.well-known/agent-card.json`
- Monitor streaming updates: Set `A2A_STREAMING_STATUS_UPDATE_INTERVAL`
- Use A2A Debugger container for interactive testing
