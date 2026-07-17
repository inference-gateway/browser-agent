# Configuration

Browser Agent is configured via environment variables. Defaults are
derived from `agent.yaml`; the env vars below override them at runtime.

## Custom Configuration

| Category | Variable | Default |
|----------|----------|---------|
| **Browser** | `BROWSER_ARGS` | `[--disable-blink-features=AutomationControlled --disable-features=VizDisplayCompositor --no-first-run --disable-default-apps --disable-extensions --disable-plugins --disable-sync --disable-translate --hide-scrollbars --mute-audio --no-zygote --disable-background-timer-throttling --disable-backgrounding-occluded-windows --disable-renderer-backgrounding --disable-ipc-flooding-protection]` |
| **Browser** | `BROWSER_DATA_DIR` | `/tmp/playwright/artifacts` |
| **Browser** | `BROWSER_ENGINE` | `chromium` |
| **Browser** | `BROWSER_HEADER_ACCEPT` | `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7` |
| **Browser** | `BROWSER_HEADER_ACCEPT_ENCODING` | `gzip, deflate, br` |
| **Browser** | `BROWSER_HEADER_ACCEPT_LANGUAGE` | `en-US,en;q=0.9` |
| **Browser** | `BROWSER_HEADER_CONNECTION` | `keep-alive` |
| **Browser** | `BROWSER_HEADER_DNT` | `1` |
| **Browser** | `BROWSER_HEADER_UPGRADE_INSECURE_REQUESTS` | `1` |
| **Browser** | `BROWSER_HEADLESS` | `true` |
| **Browser** | `BROWSER_SESSION_TIMEOUT` | `2m` |
| **Browser** | `BROWSER_STEALTH_MODE` | `false` |
| **Browser** | `BROWSER_USER_AGENT` | `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36` |
| **Browser** | `BROWSER_VIEWPORT_HEIGHT` | `1080` |
| **Browser** | `BROWSER_VIEWPORT_WIDTH` | `1920` |
| **Browser** | `BROWSER_XVFB_DISPLAY` | `:99` |
| **Browser** | `BROWSER_XVFB_ENABLED` | `false` |
| **Browser** | `BROWSER_XVFB_SCREEN_RESOLUTION` | `1920x1080x24` |
| **Tools** | `TOOLS_EDIT_ENABLED` | `true` |
| **Tools** | `TOOLS_FETCH_ALLOW_DOWNLOADS` | `true` |
| **Tools** | `TOOLS_FETCH_DOWNLOAD_DIR` | `/tmp/playwright/artifacts` |
| **Tools** | `TOOLS_FETCH_ENABLED` | `true` |
| **Tools** | `TOOLS_READ_ENABLED` | `true` |
| **Tools** | `TOOLS_READ_MAX_LINES` | `2000` |
| **Tools** | `TOOLS_WRITE_ENABLED` | `true` |

## Environment Variables

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
| **Agent Metadata** | `A2A_AGENT_CARD_FILE_PATH` | Path to agent card JSON file | `.well-known/agent-card.json` |
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
| **Artifacts** | `A2A_ARTIFACTS_ENABLE` | Enable artifacts support | `false` |
| **Artifacts** | `A2A_ARTIFACTS_SERVER_HOST` | Artifacts server host | `localhost` |
| **Artifacts** | `A2A_ARTIFACTS_SERVER_PORT` | Artifacts server port | `8081` |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_PROVIDER` | Storage backend (`filesystem` or `minio`) | `filesystem` |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_BASE_PATH` | Base path for filesystem storage | `./artifacts` |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_BASE_URL` | Override base URL for direct downloads | (auto-generated) |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_ENDPOINT` | MinIO/S3 endpoint URL | - |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_ACCESS_KEY` | MinIO/S3 access key | - |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_SECRET_KEY` | MinIO/S3 secret key | - |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_BUCKET_NAME` | MinIO/S3 bucket name | `artifacts` |
| **Artifacts** | `A2A_ARTIFACTS_STORAGE_USE_SSL` | Use SSL for MinIO/S3 connections | `true` |
| **Artifacts** | `A2A_ARTIFACTS_RETENTION_MAX_ARTIFACTS` | Max artifacts per task (0 = unlimited) | `5` |
| **Artifacts** | `A2A_ARTIFACTS_RETENTION_MAX_AGE` | Max artifact age (0 = no age limit) | `168h` |
| **Artifacts** | `A2A_ARTIFACTS_RETENTION_CLEANUP_INTERVAL` | Cleanup frequency (0 = manual only) | `24h` |
| **Authentication** | `A2A_AUTH_ENABLE` | Enable OIDC authentication | `false` |
| **Telemetry** | `A2A_TELEMETRY_ENABLE` | Enable OpenTelemetry instrumentation | `true` |
| **Telemetry** | `A2A_OTEL_TRACES_EXPORTER` | Trace exporter (`otlp` or `none`) | `none` |
| **Telemetry** | `A2A_OTEL_METRICS_EXPORTER` | Metrics exporter (`otlp`, `prometheus`, or `none`) | `prometheus` |
| **Telemetry** | `A2A_OTEL_EXPORTER_PROMETHEUS_PORT` | Prometheus metrics exporter port | `9464` |
