package playwright

import (
	"testing"

	"github.com/inference-gateway/browser-agent/config"
	"github.com/stretchr/testify/assert"
)

func TestLightpandaEngineConstant(t *testing.T) {
	assert.Equal(t, BrowserEngine("lightpanda"), Lightpanda)
}

func TestNewBrowserConfigFromConfigWithLightpanda(t *testing.T) {
	cfg := &config.Config{
		Browser: config.BrowserConfig{
			Engine:         "lightpanda",
			CdpURL:         "ws://lightpanda:9222",
			Headless:       true,
			ViewportWidth:  "1920",
			ViewportHeight: "1080",
		},
	}

	browserConfig := NewBrowserConfigFromConfig(cfg)
	assert.Equal(t, Lightpanda, browserConfig.Engine)
	assert.Equal(t, "ws://lightpanda:9222", browserConfig.CdpURL)
	assert.True(t, browserConfig.Headless)
}

func TestNewBrowserConfigFromConfigWithLightpandaNoCDPURL(t *testing.T) {
	cfg := &config.Config{
		Browser: config.BrowserConfig{
			Engine:         "lightpanda",
			CdpURL:         "",
			Headless:       true,
			ViewportWidth:  "1920",
			ViewportHeight: "1080",
		},
	}

	browserConfig := NewBrowserConfigFromConfig(cfg)
	assert.Equal(t, Lightpanda, browserConfig.Engine)
	assert.Empty(t, browserConfig.CdpURL)
}

func TestNewBrowserConfigFromConfigWithLightpandaCaseInsensitive(t *testing.T) {
	cfg := &config.Config{
		Browser: config.BrowserConfig{
			Engine: "LIGHTPANDA",
			CdpURL: "ws://lightpanda:9222",
		},
	}

	browserConfig := NewBrowserConfigFromConfig(cfg)
	assert.Equal(t, Lightpanda, browserConfig.Engine)
}

func TestDefaultBrowserConfigIsChromium(t *testing.T) {
	config := DefaultBrowserConfig()
	assert.Equal(t, Chromium, config.Engine)
	assert.Empty(t, config.CdpURL)
}

func TestBrowserEngineValues(t *testing.T) {
	assert.Equal(t, BrowserEngine("chromium"), Chromium)
	assert.Equal(t, BrowserEngine("firefox"), Firefox)
	assert.Equal(t, BrowserEngine("webkit"), WebKit)
	assert.Equal(t, BrowserEngine("lightpanda"), Lightpanda)
}

func TestNewBrowserConfigFromConfigDefaultsToChromium(t *testing.T) {
	cfg := &config.Config{
		Browser: config.BrowserConfig{
			Engine:         "unknown-engine",
			Headless:       true,
			ViewportWidth:  "1920",
			ViewportHeight: "1080",
		},
	}

	browserConfig := NewBrowserConfigFromConfig(cfg)
	assert.Equal(t, Chromium, browserConfig.Engine)
	assert.Empty(t, browserConfig.CdpURL)
}
