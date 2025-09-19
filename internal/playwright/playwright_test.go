package playwright

import (
	"testing"
	"time"
)

func TestDefaultBrowserConfig(t *testing.T) {
	config := DefaultBrowserConfig()
	
	if config.Engine != Chromium {
		t.Errorf("Expected engine to be Chromium, got %s", config.Engine)
	}
	
	if !config.Headless {
		t.Error("Expected headless to be true")
	}
	
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", config.Timeout)
	}
	
	if config.ViewportWidth != 1280 {
		t.Errorf("Expected viewport width to be 1280, got %d", config.ViewportWidth)
	}
	
	if config.ViewportHeight != 720 {
		t.Errorf("Expected viewport height to be 720, got %d", config.ViewportHeight)
	}
}

func TestBrowserEngineConstants(t *testing.T) {
	if Chromium != "chromium" {
		t.Errorf("Expected Chromium constant to be 'chromium', got %s", Chromium)
	}
	
	if Firefox != "firefox" {
		t.Errorf("Expected Firefox constant to be 'firefox', got %s", Firefox)
	}
	
	if WebKit != "webkit" {
		t.Errorf("Expected WebKit constant to be 'webkit', got %s", WebKit)
	}
}
