# Playwright Service Documentation

The Playwright Service provides comprehensive browser automation capabilities for the playwright-agent. It implements the `BrowserAutomation` interface defined in the agent specification.

## Features

- **Multi-Browser Support**: Chromium, Firefox, and WebKit
- **Concurrent Sessions**: Thread-safe browser session management
- **Configurable Browsers**: Headless/headed modes, viewport settings
- **Comprehensive Automation**: Navigation, form filling, data extraction, screenshots
- **Error Handling**: Robust error handling and recovery mechanisms
- **Resource Management**: Automatic cleanup and browser process management

## Architecture

### Core Components

- **BrowserAutomation Interface**: Defines all browser automation operations
- **BrowserSession**: Represents an active browser session with context and page
- **BrowserConfig**: Configuration options for browser launch parameters
- **playwrightImpl**: Main service implementation

### Session Management

Each browser session includes:
- Unique session ID
- Browser instance
- Browser context (for isolation)
- Active page
- Creation and last-used timestamps

## Configuration

### Browser Configuration Options

```go
type BrowserConfig struct {
    Engine         BrowserEngine  // chromium, firefox, webkit
    Headless       bool           // Run in headless mode
    Timeout        time.Duration  // Browser operation timeout
    ViewportWidth  int            // Browser viewport width
    ViewportHeight int            // Browser viewport height
    Args           []string       // Additional browser arguments
}
```

### Default Configuration

```go
config := &BrowserConfig{
    Engine:         Chromium,
    Headless:       true,
    Timeout:        30 * time.Second,
    ViewportWidth:  1280,
    ViewportHeight: 720,
    Args:           []string{"--disable-dev-shm-usage", "--no-sandbox"},
}
```

### Environment Variables

The service respects standard Playwright environment variables:
- `PLAYWRIGHT_BROWSERS_PATH`: Custom browser installation path
- `PLAYWRIGHT_DOWNLOAD_HOST`: Custom download host for browsers

## API Reference

### Browser Management

#### LaunchBrowser
```go
LaunchBrowser(ctx context.Context, config *BrowserConfig) (*BrowserSession, error)
```
Launches a new browser instance with the specified configuration.

#### CloseBrowser
```go
CloseBrowser(ctx context.Context, sessionID string) error
```
Closes a browser session and cleans up resources.

#### GetSession
```go
GetSession(sessionID string) (*BrowserSession, error)
```
Retrieves an existing browser session by ID.

### Page Operations

#### NavigateToURL
```go
NavigateToURL(ctx context.Context, sessionID, url string, waitUntil string, timeout time.Duration) error
```
Navigates to a URL and waits for the specified condition:
- `load`: Wait for load event
- `domcontentloaded`: Wait for DOM content loaded
- `networkidle`: Wait for network idle

#### ClickElement
```go
ClickElement(ctx context.Context, sessionID, selector string, options map[string]any) error
```
Clicks an element identified by CSS selector or XPath.

Options:
- `timeout`: Maximum wait time for element
- `force`: Force click even if element is not visible
- `click_count`: Number of clicks (default: 1)
- `button`: Mouse button ("left", "right", "middle")

#### FillForm
```go
FillForm(ctx context.Context, sessionID string, fields []map[string]any, submit bool, submitSelector string) error
```
Fills form fields with provided data.

Field structure:
- `selector`: Element selector
- `value`: Value to fill
- `type`: Input type ("text", "select", "checkbox", "radio")

#### ExtractData
```go
ExtractData(ctx context.Context, sessionID string, extractors []map[string]any, format string) (string, error)
```
Extracts data from the page using selectors.

Extractor structure:
- `name`: Name for the extracted field
- `selector`: CSS selector or XPath
- `attribute`: Attribute to extract (default: "text")
- `multiple`: Extract all matching elements

#### TakeScreenshot
```go
TakeScreenshot(ctx context.Context, sessionID, path string, fullPage bool, selector string, format string, quality int) error
```
Captures a screenshot of the page or specific element.

#### ExecuteScript
```go
ExecuteScript(ctx context.Context, sessionID, script string, args []any) (any, error)
```
Executes JavaScript code in the browser context.

#### WaitForCondition
```go
WaitForCondition(ctx context.Context, sessionID, condition, selector, state string, timeout time.Duration, customFunction string) error
```
Waits for specific conditions:
- `selector`: Wait for element state (visible, hidden, attached, detached)
- `navigation`: Wait for navigation (simple timeout)
- `function`: Wait for custom JavaScript function
- `timeout`: Simple timeout wait

#### HandleAuthentication
```go
HandleAuthentication(ctx context.Context, sessionID, authType, username, password, loginURL string, selectors map[string]string) error
```
Handles authentication scenarios:
- `basic`: Basic HTTP authentication
- `form`: Form-based authentication
- `oauth`: OAuth authentication (limited support)

### Service Management

#### GetHealth
```go
GetHealth(ctx context.Context) error
```
Checks the health of the service.

#### Shutdown
```go
Shutdown(ctx context.Context) error
```
Gracefully shuts down the service and closes all sessions.

## Usage Examples

### Basic Navigation
```go
// Launch browser
config := DefaultBrowserConfig()
session, err := service.LaunchBrowser(ctx, config)
if err != nil {
    return err
}
defer service.CloseBrowser(ctx, session.ID)

// Navigate to URL
err = service.NavigateToURL(ctx, session.ID, "https://example.com", "load", 30*time.Second)
```

### Form Automation
```go
fields := []map[string]any{
    {
        "selector": "#username",
        "value":    "user@example.com",
        "type":     "text",
    },
    {
        "selector": "#password",
        "value":    "password123",
        "type":     "text",
    },
}

err = service.FillForm(ctx, session.ID, fields, true, "#submit")
```

### Data Extraction
```go
extractors := []map[string]any{
    {
        "name":      "title",
        "selector":  "h1",
        "attribute": "text",
    },
    {
        "name":      "links",
        "selector":  "a",
        "attribute": "href",
        "multiple":  true,
    },
}

data, err := service.ExtractData(ctx, session.ID, extractors, "json")
```

## Error Handling

The service provides comprehensive error handling:
- Browser launch failures
- Page navigation timeouts
- Element not found errors
- Script execution errors
- Session management errors

All errors include contextual information for debugging.

## Performance Considerations

- **Concurrent Sessions**: The service supports multiple concurrent browser sessions
- **Resource Cleanup**: Automatic cleanup of browser processes and contexts
- **Timeout Management**: Configurable timeouts for all operations
- **Memory Management**: Proper disposal of browser resources

## Testing

The service includes:
- Unit tests for core functionality
- Integration tests for browser automation (run with `go test -tags integration`)
- API compatibility tests
- Error handling tests

Run tests:
```bash
# Unit tests
go test ./internal/playwright

# Integration tests (requires browsers installed)
go test -tags integration ./internal/playwright
```

## Troubleshooting

### Common Issues

1. **Browser Installation**: Ensure Playwright browsers are installed
2. **Permissions**: Check file system permissions for browser downloads
3. **Network**: Verify network access for browser downloads
4. **Resources**: Ensure sufficient memory for browser processes

### Debugging

Enable debug logging:
```bash
A2A_DEBUG=true ./playwright-agent
```

The service logs all browser operations with session IDs for tracing.