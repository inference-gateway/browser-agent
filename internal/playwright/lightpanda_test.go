package playwright

import (
	"testing"

	assert "github.com/stretchr/testify/assert"

	zap "go.uber.org/zap"

	config "github.com/inference-gateway/browser-agent/config"
)

func TestNewBrowserConfigFromConfigEngineMapping(t *testing.T) {
	chromiumArgs := []string{"--disable-dev-shm-usage", "--no-sandbox"}

	tests := []struct {
		name       string
		engine     string
		cdpURL     string
		wantEngine BrowserEngine
		wantCDPURL string
		wantArgs   []string
	}{
		{"lightpanda", "lightpanda", "ws://lightpanda:9222", Lightpanda, "ws://lightpanda:9222", nil},
		{"lightpanda is case insensitive", "LIGHTPANDA", "ws://lightpanda:9222", Lightpanda, "ws://lightpanda:9222", nil},
		{"lightpanda without a cdp url", "lightpanda", "", Lightpanda, "", nil},
		{"unknown engine falls back to chromium", "unknown-engine", "", Chromium, "", chromiumArgs},
		{"chromium keeps its launch args", "chromium", "", Chromium, "", chromiumArgs},
		// WebKit rejects Chromium flags outright: "Unknown option --disable-dev-shm-usage".
		{"webkit gets no chromium args", "webkit", "", WebKit, "", nil},
		{"firefox gets no chromium args", "firefox", "", Firefox, "", nil},
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
			assert.Equal(t, tt.wantArgs, browserConfig.Args)
			assert.True(t, browserConfig.Headless)
		})
	}
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
