package playwright

import (
	"context"
	"os"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"

	zap "go.uber.org/zap"

	config "github.com/inference-gateway/browser-agent/config"
)

func TestNewBrowserConfigFromConfigEngineMapping(t *testing.T) {
	tests := []struct {
		name       string
		engine     string
		cdpURL     string
		wantEngine BrowserEngine
		wantCDPURL string
	}{
		{"lightpanda", "lightpanda", "ws://lightpanda:9222", Lightpanda, "ws://lightpanda:9222"},
		{"lightpanda is case insensitive", "LIGHTPANDA", "ws://lightpanda:9222", Lightpanda, "ws://lightpanda:9222"},
		{"lightpanda without a cdp url", "lightpanda", "", Lightpanda, ""},
		{"unknown engine falls back to chromium", "unknown-engine", "", Chromium, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			browserConfig := NewBrowserConfigFromConfig(&config.Config{
				Browser: config.BrowserConfig{
					Engine:         tt.engine,
					CDPURL:         tt.cdpURL,
					Headless:       true,
					ViewportWidth:  "1920",
					ViewportHeight: "1080",
				},
			})

			assert.Equal(t, tt.wantEngine, browserConfig.Engine)
			assert.Equal(t, tt.wantCDPURL, browserConfig.CDPURL)
			assert.True(t, browserConfig.Headless)
		})
	}
}

// Drives a real Lightpanda when BROWSER_CDP_URL points at one, e.g.
//
//	docker run -d --rm -p 9222:9222 lightpanda/browser:nightly
//	BROWSER_CDP_URL=ws://127.0.0.1:9222 go test ./internal/playwright/ -run Lightpanda
//
// Guards the assumption the CDP path rests on: Lightpanda accepts
// browser.NewContext(), so it shares the local-launch code path.
func TestLightpandaEndToEnd(t *testing.T) {
	cdpURL := os.Getenv("BROWSER_CDP_URL")
	if cdpURL == "" {
		t.Skip("BROWSER_CDP_URL not set - no Lightpanda endpoint to drive")
	}

	service, err := NewPlaywrightService(zap.NewNop(), &config.Config{
		Browser: config.BrowserConfig{
			Engine:         string(Lightpanda),
			CDPURL:         cdpURL,
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
}

// The guard runs before p.pw is touched, so this needs no live browser.
func TestAcquireBrowserLightpandaRequiresCDPURL(t *testing.T) {
	p := &playwrightImpl{logger: zap.NewNop()}

	_, err := p.acquireBrowser(&BrowserConfig{Engine: Lightpanda})

	assert.ErrorContains(t, err, "BROWSER_CDP_URL is required")
}

func TestAcquireBrowserRejectsUnknownEngine(t *testing.T) {
	p := &playwrightImpl{logger: zap.NewNop()}

	_, err := p.acquireBrowser(&BrowserConfig{Engine: BrowserEngine("nope")})

	assert.ErrorContains(t, err, "unsupported browser engine")
}
