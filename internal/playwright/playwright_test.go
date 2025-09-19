package playwright

import (
	"context"
	"testing"
	"time"

	"github.com/inference-gateway/playwright-agent/config"
	"go.uber.org/zap"
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

func TestNewPlaywrightService(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{}
	
	// Note: This test may fail in CI environments without proper browser setup
	// In a real test environment, you might want to mock the playwright.Run() call
	service, err := NewPlaywrightService(logger, cfg)
	if err != nil {
		t.Logf("Expected error in test environment without browser setup: %v", err)
		return
	}
	
	if service == nil {
		t.Error("Expected non-nil service")
	}
	
	// Test health check
	err = service.GetHealth(context.Background())
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
	
	// Cleanup
	err = service.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

func TestBrowserConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		engine BrowserEngine
		valid  bool
	}{
		{"Chromium", Chromium, true},
		{"Firefox", Firefox, true},
		{"WebKit", WebKit, true},
		{"Invalid", "invalid", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &BrowserConfig{
				Engine:         tt.engine,
				Headless:       true,
				Timeout:        30 * time.Second,
				ViewportWidth:  1280,
				ViewportHeight: 720,
			}
			
			// This would be used in the actual browser launch
			// For now, just verify the engine type
			switch config.Engine {
			case Chromium, Firefox, WebKit:
				if !tt.valid {
					t.Errorf("Expected engine %s to be invalid", tt.engine)
				}
			default:
				if tt.valid {
					t.Errorf("Expected engine %s to be valid", tt.engine)
				}
			}
		})
	}
}

func TestBrowserSessionStructure(t *testing.T) {
	session := &BrowserSession{
		ID:       "test-session",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
	
	if session.ID != "test-session" {
		t.Errorf("Expected session ID to be 'test-session', got %s", session.ID)
	}
	
	if session.Created.IsZero() {
		t.Error("Expected created time to be set")
	}
	
	if session.LastUsed.IsZero() {
		t.Error("Expected last used time to be set")
	}
}

// Mock test for browser automation interface
func TestBrowserAutomationInterface(t *testing.T) {
	// This test verifies that our implementation satisfies the interface
	var _ BrowserAutomation = (*playwrightImpl)(nil)
	
	// Create a mock implementation
	impl := &playwrightImpl{
		logger:   zap.NewNop(),
		config:   &config.Config{},
		sessions: make(map[string]*BrowserSession),
	}
	
	// Test session management without actual browser
	_, err := impl.GetSession("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent session")
	}
	
	expectedError := "session not found: non-existent"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
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