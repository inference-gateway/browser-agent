# Getting Started

`browser-agent` is an A2A (Agent-to-Agent) server that drives a real
[Playwright](https://playwright.dev/) browser to automate web tasks. This page
covers building and running it locally.

## Prerequisites

- Go 1.26.4+ (only if building/running from source)
- A browser for Playwright to drive. The provided `Dockerfile` bundles Lightpanda
  and starts it for you — this is the simplest way to run the agent. Running from
  source instead needs either your own Lightpanda on `BROWSER_CDP_URL` or
  `BROWSER_ENGINE=chromium` with the Playwright browsers installed. See
  [Browser engines](configuration.md#browser-engines).
- An OpenAI-compatible LLM endpoint and credentials (see
  [Configuration](configuration.md)).

## Run from source

```bash
# Start the A2A server (defaults to port 8080)
go run . start

# Or build the CLI and invoke it directly
task build
./bin/browser-agent --help
./bin/browser-agent --version
./bin/browser-agent start
```

The binary is a [Cobra](https://github.com/spf13/cobra) CLI. The root command
exposes `--help` and `--version`; the `start` subcommand boots the server and
blocks until it receives `SIGINT`/`SIGTERM`.

## Run with Docker

```bash
docker build -t browser-agent .
docker run -p 8080:8080 browser-agent
```

Published images ship one browser each — `:latest` and `:chromium` are Chromium,
with `:firefox`, `:webkit` and a lean `:lightpanda` alongside. See
[Configuration](configuration.md#browser-engines) for the tradeoffs.

## Run the local stack

A `docker-compose.yaml` brings up the Inference Gateway alongside the agent
(built from the local `Dockerfile`, so code changes are picked up):

```bash
docker compose up --build
```

Opt-in profiles add the `infer` CLI and the `a2a-debugger`:

```bash
docker compose --profile cli run --rm cli chat
docker compose --profile debugger run --rm debugger \
  --server-url http://browser-agent:8080 tasks submit "What are your skills?"
```

## Verify it is running

```bash
# Health check
curl http://localhost:8080/health

# Agent metadata (capabilities, tools, skills)
curl http://localhost:8080/.well-known/agent-card.json
```

## Next steps

- [Configuration](configuration.md) — environment variables, LLM provider, and
  browser settings.
- [Usage](usage.md) — submitting tasks and driving the agent.
- [Playwright Service](playwright-service.md) — the browser-automation service
  that backs the tools.
