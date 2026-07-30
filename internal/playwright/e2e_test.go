package playwright

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"

	zap "go.uber.org/zap"

	config "github.com/inference-gateway/browser-agent/config"
)

// Skipped unless BROWSER_ENGINE is set. Every published image sets it, so
// cross-compiling this test into one verifies the browser it actually ships:
//
//	GOOS=linux GOARCH=arm64 go test -c ./internal/playwright/ -o pw.test
//	docker cp pw.test verify:/tmp/pw.test
//	docker exec verify /tmp/pw.test -test.run TestEngineEndToEnd -test.v
func TestEngineEndToEnd(t *testing.T) {
	engine := os.Getenv("BROWSER_ENGINE")
	if engine == "" {
		t.Skip("BROWSER_ENGINE not set - no browser to drive")
	}

	outDir := os.Getenv("BROWSER_DATA_DIR")
	if outDir == "" {
		outDir = t.TempDir()
	}

	service, err := NewPlaywrightService(zap.NewNop(), &config.Config{
		Browser: config.BrowserConfig{
			Engine:         engine,
			CDPURL:         os.Getenv("BROWSER_CDP_URL"),
			Headless:       true,
			ViewportWidth:  "1920",
			ViewportHeight: "1080",
		},
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, service.Shutdown(context.Background())) }()

	session, err := service.LaunchBrowser(context.Background(), nil)
	require.NoError(t, err)
	defer func() { _ = service.CloseBrowser(context.Background(), session.ID) }()

	require.NoError(t, service.NavigateToURL(context.Background(), session.ID, "https://example.com", "load", 30*time.Second))

	h1, err := session.Page.Locator("h1").TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Example Domain", h1)

	shot := filepath.Join(outDir, engine+".png")
	err = service.TakeScreenshot(context.Background(), session.ID, shot, false, "", "png", 0)

	if engine == string(Lightpanda) {
		assert.ErrorContains(t, err, "not supported by the lightpanda engine")
		return
	}

	require.NoError(t, err)
	info, err := os.Stat(shot)
	require.NoError(t, err)
	assert.NotZero(t, info.Size())
	t.Logf("screenshot: %s (%d bytes)", shot, info.Size())
}
