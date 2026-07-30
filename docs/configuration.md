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
| `BROWSER_ENGINE` | Browser engine: `lightpanda`, `chromium`, `firefox`, `webkit` | `lightpanda` |
| `BROWSER_CDP_URL` | CDP endpoint to drive when the engine is `lightpanda` | `ws://127.0.0.1:9222` |
| `BROWSER_HEADLESS` | Run headless | `true` |
| `BROWSER_STEALTH_MODE` | Enable stealth patches | `false` |
| `BROWSER_SESSION_TIMEOUT` | Idle session timeout | `2m` |
| `BROWSER_VIEWPORT_WIDTH` | Viewport width | `1920` |
| `BROWSER_VIEWPORT_HEIGHT` | Viewport height | `1080` |
| `BROWSER_USER_AGENT` | User-Agent header | Chrome 131 UA |
| `BROWSER_DATA_DIR` | Scratch/artifacts directory | `/tmp/playwright/artifacts` |
| `BROWSER_XVFB_ENABLED` | Run under Xvfb (for headed mode on a headless host) | `false` |

### Browser engines

`lightpanda` is the default. The Docker image bundles the
[Lightpanda](https://github.com/lightpanda-io/browser) binary and the entrypoint
starts it on `127.0.0.1:9222`, so nothing extra needs to run. Playwright drives
it over CDP rather than launching a browser per session, which cuts cold-start
time and drops Chromium's apt dependency tree from the image.

The tradeoff is coverage: Lightpanda implements a subset of the Web platform, so
JavaScript-heavy pages that work under Chromium can fail or render incompletely.
It has no graphical rendering engine at all, so `take_screenshot` returns an
error — use a Chromium, Firefox or WebKit image for anything visual.

The engine is baked into the image so that nothing has to be downloaded or
installed at runtime — each variant ships exactly one browser and its system
libraries. Pick the tag that matches the engine you want rather than setting
`BROWSER_ENGINE` on the default image:

| Image tag | Engine |
|-----------|--------|
| `ghcr.io/inference-gateway/browser-agent:<version>`, `:latest` | `lightpanda` |
| `ghcr.io/inference-gateway/browser-agent:lightpanda-<version>`, `:lightpanda` | `lightpanda` (explicit alias) |
| `ghcr.io/inference-gateway/browser-agent:chromium-<version>`, `:chromium` | `chromium` |
| `ghcr.io/inference-gateway/browser-agent:firefox-<version>`, `:firefox` | `firefox` |
| `ghcr.io/inference-gateway/browser-agent:webkit-<version>`, `:webkit` | `webkit` |

Each image sets `BROWSER_ENGINE` to its own engine, so no configuration is
needed. Setting it to an engine the image doesn't ship will fail at startup. To
build a variant locally:

```sh
docker build --build-arg BROWSER_ENGINE=chromium -t browser-agent:chromium .
```

Set `BROWSER_CDP_URL` to a non-loopback address to drive a Lightpanda running
elsewhere; the entrypoint then skips the bundled one. See
`docker-compose.lightpanda.yaml` for a sidecar example.

`chromium`, `firefox` and `webkit` launch locally via Playwright and ignore
`BROWSER_CDP_URL`. Outside Docker (`task run`), the bundled binary does not
exist — either run a Lightpanda yourself or set `BROWSER_ENGINE=chromium`.

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
