<div align="center">

# Browser-Agent
[![CI](https://github.com/inference-gateway/browser-agent/workflows/CI/badge.svg)](https://github.com/inference-gateway/browser-agent/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![A2A Protocol](https://img.shields.io/badge/A2A-Protocol-blue?style=flat)](https://github.com/inference-gateway/adk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AI agent for browser automation and web testing using Playwright**

A production-ready [Agent-to-Agent (A2A)](https://github.com/inference-gateway/adk) server that provides AI-powered capabilities through a standardized protocol.

</div>

## Quick Start

```bash
# Run the agent
go run .

# Or with Docker
docker build -t browser-agent .
docker run -p 8080:8080 browser-agent
```

## Features

- ✅ A2A protocol compliant
- ✅ AI-powered capabilities
- ✅ Streaming support
- ✅ Production ready
- ✅ Minimal dependencies

## Endpoints

- `GET /.well-known/agent.json` - Agent metadata and capabilities
- `GET /health` - Health check endpoint
- `POST /a2a` - A2A protocol endpoint

## Available Skills

| Skill | Description | Parameters |
|-------|-------------|------------|
| `navigate_to_url` | Navigate to a specific URL and wait for the page to fully load |timeout, url, wait_until |
| `click_element` | Click on an element identified by selector, text, or other locator strategies |button, click_count, force, selector, timeout |
| `fill_form` | Fill form fields with provided data, handling various input types |fields, submit, submit_selector |
| `extract_data` | Extract data from the page using selectors and return structured information |extractors, format |
| `take_screenshot` | Capture a screenshot of the current page or specific element |full_page, quality, selector, type |
| `execute_script` | Execute custom JavaScript code in the browser context |args, return_value, script |
| `handle_authentication` | Handle various authentication scenarios including basic auth, OAuth, and custom login forms |login_url, password, password_selector, submit_selector, type, username, username_selector |
| `wait_for_condition` | Wait for specific conditions before proceeding with automation |condition, custom_function, selector, state, timeout |
| `write_to_csv` | Write structured data to CSV files with support for custom headers and file paths |append, data, filename, headers, include_headers |

## Configuration

Configure the agent via environment variables:

### Custom Configuration

The following custom configuration variables are available:

| Category | Variable | Description | Default |
|----------|----------|-------------|---------|
| **Browser** | `BROWSER_ARGS` | Args configuration | `[--disable-blink-features=AutomationControlled --disable-features=VizDisplayCompositor --no-first-run --disable-default-apps --disable-extensions --disable-plugins --disable-sync --disable-translate --hide-scrollbars --mute-audio --no-zygote --disable-background-timer-throttling --disable-backgrounding-occluded-windows --disable-renderer-backgrounding --disable-ipc-flooding-protection]` |
| **Browser** | `BROWSER_DATA_DIR` | Data_dir configuration | `/tmp/playwright/artifacts` |
| **Browser** | `BROWSER_HEADER_ACCEPT` | Header_accept configuration | `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7` |
| **Browser** | `BROWSER_HEADER_ACCEPT_ENCODING` | Header_accept_encoding configuration | `gzip, deflate, br` |
| **Browser** | `BROWSER_HEADER_ACCEPT_LANGUAGE` | Header_accept_language configuration | `en-US,en;q=0.9` |
| **Browser** | `BROWSER_HEADER_CONNECTION` | Header_connection configuration | `keep-alive` |
| **Browser** | `BROWSER_HEADER_DNT` | Header_dnt configuration | `1` |
| **Browser** | `BROWSER_HEADER_UPGRADE_INSECURE_REQUESTS` | Header_upgrade_insecure_requests configuration | `1` |
| **Browser** | `BROWSER_USER_AGENT` | User_agent configuration | `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36` |
| **Browser** | `BROWSER_VIEWPORT_HEIGHT` | Viewport_height configuration | `1080` |
| **Browser** | `BROWSER_VIEWPORT_WIDTH` | Viewport_width configuration | `1920` |

| Category | Variable | Description | Default |
|----------|----------|-------------|---------|
| **Server** | `A2A_PORT` | Server port | `8080` |
| **Server** | `A2A_DEBUG` | Enable debug mode | `false` |
| **Server** | `A2A_AGENT_URL` | Agent URL for internal references | `http://localhost:8080` |
| **Server** | `A2A_STREAMING_STATUS_UPDATE_INTERVAL` | Streaming status update frequency | `1s` |
| **Server** | `A2A_SERVER_READ_TIMEOUT` | HTTP server read timeout | `120s` |
| **Server** | `A2A_SERVER_WRITE_TIMEOUT` | HTTP server write timeout | `120s` |
| **Server** | `A2A_SERVER_IDLE_TIMEOUT` | HTTP server idle timeout | `120s` |
| **Server** | `A2A_SERVER_DISABLE_HEALTHCHECK_LOG` | Disable logging for health check requests | `true` |
| **Agent Metadata** | `A2A_AGENT_CARD_FILE_PATH` | Path to agent card JSON file | `.well-known/agent.json` |
| **LLM Client** | `A2A_AGENT_CLIENT_PROVIDER` | LLM provider (`openai`, `anthropic`, `azure`, `ollama`, `deepseek`) |`` |
| **LLM Client** | `A2A_AGENT_CLIENT_MODEL` | Model to use |`` |
| **LLM Client** | `A2A_AGENT_CLIENT_API_KEY` | API key for LLM provider | - |
| **LLM Client** | `A2A_AGENT_CLIENT_BASE_URL` | Custom LLM API endpoint | - |
| **LLM Client** | `A2A_AGENT_CLIENT_TIMEOUT` | Timeout for LLM requests | `30s` |
| **LLM Client** | `A2A_AGENT_CLIENT_MAX_RETRIES` | Maximum retries for LLM requests | `3` |
| **LLM Client** | `A2A_AGENT_CLIENT_MAX_CHAT_COMPLETION_ITERATIONS` | Max chat completion rounds | `10` |
| **LLM Client** | `A2A_AGENT_CLIENT_MAX_TOKENS` | Maximum tokens for LLM responses |`4096` |
| **LLM Client** | `A2A_AGENT_CLIENT_TEMPERATURE` | Controls randomness of LLM output |`0.7` |
| **Capabilities** | `A2A_CAPABILITIES_STREAMING` | Enable streaming responses | `true` |
| **Capabilities** | `A2A_CAPABILITIES_PUSH_NOTIFICATIONS` | Enable push notifications | `false` |
| **Capabilities** | `A2A_CAPABILITIES_STATE_TRANSITION_HISTORY` | Track state transitions | `false` |
| **Task Management** | `A2A_TASK_RETENTION_MAX_COMPLETED_TASKS` | Max completed tasks to keep (0 = unlimited) | `100` |
| **Task Management** | `A2A_TASK_RETENTION_MAX_FAILED_TASKS` | Max failed tasks to keep (0 = unlimited) | `50` |
| **Task Management** | `A2A_TASK_RETENTION_CLEANUP_INTERVAL` | Cleanup frequency (0 = manual only) | `5m` |
| **Storage** | `A2A_QUEUE_PROVIDER` | Storage backend (`memory` or `redis`) | `memory` |
| **Storage** | `A2A_QUEUE_URL` | Redis connection URL (when using Redis) | - |
| **Storage** | `A2A_QUEUE_MAX_SIZE` | Maximum queue size | `100` |
| **Storage** | `A2A_QUEUE_CLEANUP_INTERVAL` | Task cleanup interval | `30s` |
| **Authentication** | `A2A_AUTH_ENABLE` | Enable OIDC authentication | `false` |

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
# Build with default values from ADL
docker build -t browser-agent .

# Build with custom version information
docker build \
  --build-arg VERSION=1.2.3 \
  --build-arg AGENT_NAME="My Custom Agent" \
  --build-arg AGENT_DESCRIPTION="Custom agent description" \
  -t browser-agent:1.2.3 .
```

**Available Build Arguments:**
- `VERSION` - Agent version (default: `0.1.3`)
- `AGENT_NAME` - Agent name (default: `browser-agent`)
- `AGENT_DESCRIPTION` - Agent description (default: `AI agent for browser automation and web testing using Playwright`)

These values are embedded into the binary at build time using linker flags, making them accessible at runtime without requiring environment variables.

## License

MIT License - see LICENSE file for details
