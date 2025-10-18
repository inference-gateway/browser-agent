# Browser Agent Example

This example demonstrates how to use the browser-agent for AI-powered browser automation using Playwright. The agent can navigate webpages, fill forms, take screenshots, extract data, and more.

## Prerequisites

Configure the environment variables:

```bash
cp .env.example .env
```

**Note:** Add at least two LLM provider API keys (e.g., Google and DeepSeek) in the `.env` file.

## Quick Start

### Headless Mode (Default)

Start all containers in headless mode (fastest, most secure):

```bash
docker compose up --build
```

### Headed Mode with VNC (Visual Debugging)

To view the browser in real-time via VNC:

1. **Update docker-compose.yaml agent service:**
   ```yaml
   BROWSER_HEADLESS: false
   BROWSER_XVFB_ENABLED: true
   BROWSER_STEALTH_MODE: true  # Optional: helps avoid bot detection
   ```

2. **Start with VNC profile:**
   ```bash
   docker compose --profile vnc up --build
   ```

3. **Connect to VNC:**
   ```bash
   # macOS
   open vnc://localhost:5900
   # Password: password

   # Or use any VNC client: localhost:5900
   ```

## Usage

Go into the CLI for convenience:

```bash
docker compose run --rm cli
```

Ask the following:

```text
Please visit http://demo-site which is running locally and take a screenshot of the homepage. Use the agent.
```

You would see the CLI (A2A agent client) submitting a task to the A2A agent server and the screenshot will appear in the `screenshots` directory since it's mounted as a volume.

```text
Please visit http://demo-site which is running locally and collect all of the prices, write them to a CSV file. Use the agent.
```

You would see the CLI (A2A agent client) submitting a task to the A2A agent server and the csv file with all of the prices of the website will appear inside of the artifacts directory.

Check the logs to see that the browser indeed went to the demo site and took a screenshot:

```bash
docker compose logs -f demo-site
```

Also you can check the task was successfully submitted to the agent and it's available using the a2a debugger:

```bash
docker compose run --rm a2a-debugger tasks list
```

Finally clean up:

```bash
docker compose down
```

## Configuration Options

### Browser Modes

The browser-agent supports different operational modes:

**Headless Production Mode (Default):**
- `BROWSER_HEADLESS: true`
- `BROWSER_XVFB_ENABLED: false`
- `BROWSER_STEALTH_MODE: false`
- Fastest, most secure, lowest resource usage
- Best for production/CI/CD

**Headed Mode with VNC (Development):**
- `BROWSER_HEADLESS: false`
- `BROWSER_XVFB_ENABLED: true`
- `BROWSER_STEALTH_MODE: true`
- Visual browser viewing via VNC
- Best for development, debugging, demos

**Headless with Extensions:**
- `BROWSER_HEADLESS: true`
- `BROWSER_XVFB_ENABLED: true`
- Required for browser extensions
- Specific rendering features

### Browser Engines

You can choose different browser engines by modifying the build args and environment:

```yaml
build:
  args:
    BROWSER_ENGINE: firefox  # chromium (default), firefox, webkit, or all
environment:
  BROWSER_ENGINE: firefox
```

Or build directly with docker:

```bash
# Build with default browser (chromium)
docker build -t browser-agent ..

# Build with specific browser engine
docker build --build-arg BROWSER_ENGINE=firefox -t browser-agent:firefox ..

# Build with all browsers (larger image)
docker build --build-arg BROWSER_ENGINE=all -t browser-agent:all ..
```

**Available Build Arguments:**
- `VERSION` - Agent version (default: from agent.yaml)
- `AGENT_NAME` - Agent name (default: from agent.yaml)
- `AGENT_DESCRIPTION` - Agent description (default: from agent.yaml)
- `BROWSER_ENGINE` - Browser to install (`chromium`, `firefox`, `webkit`, or `all`) (default: `chromium`)

### Xvfb Configuration

When Xvfb is enabled, you can customize:

```yaml
BROWSER_XVFB_DISPLAY: ":99"                    # X11 display number
BROWSER_XVFB_SCREEN_RESOLUTION: "1920x1080x24" # Resolution and color depth
```

**Security Note:** Xvfb is configured without the `-ac` flag (access control enabled) and uses `-nolisten tcp` to prevent remote network access.

## Troubleshooting

### VNC Connection Issues

If VNC doesn't connect:

1. Check Xvfb is enabled:
   ```bash
   docker compose exec agent env | grep XVFB
   ```

2. Check X11 socket exists:
   ```bash
   docker compose exec agent ls -la /tmp/.X11-unix/
   ```

3. Check VNC logs:
   ```bash
   docker compose logs browser-vnc
   ```

### Browser Not Starting

Check agent logs for errors:
```bash
docker compose logs agent
```
