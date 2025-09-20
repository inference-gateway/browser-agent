# Playwright Service

Browser automation service using Playwright with comprehensive mock-based testing.

## Testing

All tests use mocks for fast execution (~0.18s). No browser downloads required.

```bash
# Run all tests
go test -v ./internal/playwright

# Generate mocks (after interface changes)
go run github.com/maxbrunsfeld/counterfeiter/v6 -o internal/playwright/mocks/browser_automation.go internal/playwright BrowserAutomation
```

## Files

- `playwright.go` - Service implementation (generated, do not edit)
- `playwright_test.go` - Basic unit tests
- `playwright_integration_test.go` - Integration tests
- `mocks/browser_automation.go` - Generated mock