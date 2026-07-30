# Configuration

`browser-agent` is configured entirely through environment variables. Defaults
are derived from `spec.config.*` in `agent.yaml`; the variables below override
them at runtime. This page lists the settings most relevant to this agent — see
the generated `README.md` for the exhaustive table.

## LLM provider

The agent uses an OpenAI-compatible LLM client.

| Variable | Description | Default |
|----------|-------------|---------|
| `A2A_AGENT_CLIENT_PROVIDER` | Provider: `openai`, `anthropic`, `azure`, `ollama`, `deepseek` | – |
| `A2A_AGENT_CLIENT_MODEL` | Model identifier | – |
| `A2A_AGENT_CLIENT_API_KEY` | API key for the provider | – |
| `A2A_AGENT_CLIENT_BASE_URL` | Custom endpoint (e.g. the Inference Gateway) | – |
| `A2A_AGENT_CLIENT_MAX_TOKENS` | Max tokens per response | `4096` |
| `A2A_AGENT_CLIENT_TEMPERATURE` | Sampling temperature | `0.7` |

## Server

| Variable | Description | Default |
|----------|-------------|---------|
| `A2A_PORT` / `A2A_SERVER_PORT` | Server port | `8080` |
| `A2A_DEBUG` | Enable debug logging | `false` |
| `A2A_STREAMING_STATUS_UPDATE_INTERVAL` | Streaming status update frequency | `1s` |

## Browser

These map to `spec.config.browser` and control the Playwright runtime.

| Variable | Description | Default |
|----------|-------------|---------|
| `BROWSER_ENGINE` | Browser engine: `chromium`, `firefox`, `webkit`, `lightpanda` | `chromium` |
| `BROWSER_CDP_URL` | Drive a browser at this CDP endpoint instead of launching one | _(unset)_ |
| `BROWSER_HEADLESS` | Run headless | `true` |
| `BROWSER_STEALTH_MODE` | Enable stealth patches | `false` |
| `BROWSER_SESSION_TIMEOUT` | Idle session timeout | `2m` |
| `BROWSER_VIEWPORT_WIDTH` | Viewport width | `1920` |
| `BROWSER_VIEWPORT_HEIGHT` | Viewport height | `1080` |
| `BROWSER_USER_AGENT` | User-Agent header | Chrome 131 UA |
| `BROWSER_DATA_DIR` | Scratch/artifacts directory | `/tmp/playwright/artifacts` |
| `BROWSER_XVFB_ENABLED` | Run under Xvfb (for headed mode on a headless host) | `false` |

### Browser engines

`chromium` is the default and covers the full Web platform. `firefox` and
`webkit` launch locally through Playwright the same way.

`lightpanda` is the lean alternative: the image bundles the
[Lightpanda](https://github.com/lightpanda-io/browser) binary and the entrypoint
starts it on `127.0.0.1:9222`, so nothing extra needs to run. Playwright drives
it over CDP instead of launching a browser per session, which cuts cold-start
time and drops Chromium's apt dependency tree — 871MB against 3.03GB.

The tradeoff is coverage. Lightpanda implements a subset of the Web platform, so
JavaScript-heavy pages that work under Chromium can fail or render incompletely,
and it has no graphical rendering engine at all, so `take_screenshot` returns an
error. Every shipped skill uses screenshots, so pick it only for text and DOM
extraction where speed matters.

The engine is baked into the image so nothing has to be downloaded or installed
at runtime — each variant ships exactly one browser and its system libraries.
Pick the tag that matches the engine you want:

| Image tag | Engine | Size |
|-----------|--------|------|
| `ghcr.io/inference-gateway/browser-agent:<version>`, `:latest` | `chromium` | 3.03GB |
| `ghcr.io/inference-gateway/browser-agent:chromium-<version>`, `:chromium` | `chromium` (explicit alias) | 3.03GB |
| `ghcr.io/inference-gateway/browser-agent:firefox-<version>`, `:firefox` | `firefox` | 1.74GB |
| `ghcr.io/inference-gateway/browser-agent:webkit-<version>`, `:webkit` | `webkit` | 1.95GB |
| `ghcr.io/inference-gateway/browser-agent:lightpanda-<version>`, `:lightpanda` | `lightpanda` | 871MB |

Each image sets `BROWSER_ENGINE` to its own engine, so no configuration is
needed. Setting it to an engine the image doesn't ship will fail at startup. To
build a variant locally:

```sh
docker build --build-arg BROWSER_ENGINE=lightpanda -t browser-agent:lightpanda .
```

### Driving a remote browser over CDP

Set `BROWSER_CDP_URL` and the agent connects to that endpoint instead of
launching a browser of its own. This is the Chrome DevTools Protocol, so it
works with `chromium` and `lightpanda`; Playwright cannot drive `firefox` or
`webkit` over CDP and setting the variable for those engines is rejected at
startup.

```sh
# a Chrome started with --remote-debugging-port, a browser pod, a hosted service
docker run -e BROWSER_CDP_URL=ws://chrome.internal:9222 \
  ghcr.io/inference-gateway/browser-agent:latest
```

No local browser is installed when a CDP URL is set, and the endpoint is dialled
once at startup so an unreachable one fails immediately rather than on the first
task.

For `lightpanda` the entrypoint starts the bundled binary and points
`BROWSER_CDP_URL` at it automatically; give it a non-loopback address to use one
running elsewhere instead. See `docker-compose.lightpanda.yaml` for a sidecar
example. Outside Docker (`task run`) there is no bundled binary, so `lightpanda`
needs an endpoint you run yourself.

## Built-in tools

The `read`, `write`, `edit`, and `fetch` tools are toggled and tuned here.

| Variable | Description | Default |
|----------|-------------|---------|
| `TOOLS_READ_ENABLED` | Enable the `read` tool | `true` |
| `TOOLS_READ_MAX_LINES` | Default read slice | `2000` |
| `TOOLS_WRITE_ENABLED` | Enable the `write` tool | `true` |
| `TOOLS_EDIT_ENABLED` | Enable the `edit` tool | `true` |
| `TOOLS_FETCH_ENABLED` | Enable the `fetch` tool | `true` |
| `TOOLS_FETCH_ALLOW_DOWNLOADS` | Allow `fetch` to save response bodies | `true` |
| `TOOLS_FETCH_DOWNLOAD_DIR` | Download directory for `fetch` | `/tmp/playwright/artifacts` |

## Artifacts

Screenshots and extracted data are saved as downloadable artifacts. Enable the
artifacts server with `A2A_ARTIFACTS_ENABLE=true` (`filesystem` and `minio`
backends are supported). See the `README.md` for the full `A2A_ARTIFACTS_*`
table.

## Example `.env`

```bash
A2A_AGENT_CLIENT_PROVIDER=openai
A2A_AGENT_CLIENT_MODEL=gpt-4o
A2A_AGENT_CLIENT_API_KEY=sk-...
A2A_DEBUG=true
BROWSER_HEADLESS=true
```
