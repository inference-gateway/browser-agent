// +build integration

package playwright

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/inference-gateway/playwright-agent/config"
	"go.uber.org/zap"
)

// TestPlaywrightIntegration tests the full playwright integration
// Run with: go test -tags integration
func TestPlaywrightIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	if os.Getenv("CI") == "true" && os.Getenv("PLAYWRIGHT_BROWSERS_INSTALLED") != "true" {
		t.Skip("Skipping integration test in CI without browser installation")
	}
	
	logger := zap.NewNop()
	cfg := &config.Config{}
	
	service, err := NewPlaywrightService(logger, cfg)
	if err != nil {
		t.Fatalf("Failed to create playwright service: %v", err)
	}
	defer func() {
		if shutdownErr := service.Shutdown(context.Background()); shutdownErr != nil {
			t.Logf("Failed to shutdown service: %v", shutdownErr)
		}
	}()
	
	ctx := context.Background()
	
	config := DefaultBrowserConfig()
	session, err := service.LaunchBrowser(ctx, config)
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}
	
	if session == nil {
		t.Fatal("Expected non-nil session")
	}
	
	if session.ID == "" {
		t.Fatal("Expected non-empty session ID")
	}
	
	err = service.NavigateToURL(ctx, session.ID, "https://example.com", "load", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to navigate to URL: %v", err)
	}
	
	screenshotPath := "/tmp/test-screenshot.png"
	err = service.TakeScreenshot(ctx, session.ID, screenshotPath, false, "", "png", 80)
	if err != nil {
		t.Fatalf("Failed to take screenshot: %v", err)
	}
	
	if _, err := os.Stat(screenshotPath); os.IsNotExist(err) {
		t.Error("Screenshot file was not created")
	} else {
		os.Remove(screenshotPath)
	}
	
	result, err := service.ExecuteScript(ctx, session.ID, "return document.title", nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
	
	if result == nil {
		t.Error("Expected non-nil script result")
	}
	
	extractors := []map[string]interface{}{
		{
			"name":     "title",
			"selector": "title",
			"attribute": "text",
		},
	}
	
	data, err := service.ExtractData(ctx, session.ID, extractors, "json")
	if err != nil {
		t.Fatalf("Failed to extract data: %v", err)
	}
	
	if data == "" {
		t.Error("Expected non-empty extracted data")
	}
	
	retrievedSession, err := service.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve session: %v", err)
	}
	
	if retrievedSession.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, retrievedSession.ID)
	}
	
	err = service.CloseBrowser(ctx, session.ID)
	if err != nil {
		t.Fatalf("Failed to close browser: %v", err)
	}
	
	_, err = service.GetSession(session.ID)
	if err == nil {
		t.Error("Expected error when getting closed session")
	}
}

func TestMultipleBrowserSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	if os.Getenv("CI") == "true" && os.Getenv("PLAYWRIGHT_BROWSERS_INSTALLED") != "true" {
		t.Skip("Skipping integration test in CI without browser installation")
	}
	
	logger := zap.NewNop()
	cfg := &config.Config{}
	
	service, err := NewPlaywrightService(logger, cfg)
	if err != nil {
		t.Fatalf("Failed to create playwright service: %v", err)
	}
	defer func() {
		if shutdownErr := service.Shutdown(context.Background()); shutdownErr != nil {
			t.Logf("Failed to shutdown service: %v", shutdownErr)
		}
	}()
	
	ctx := context.Background()
	
	config1 := DefaultBrowserConfig()
	config1.Engine = Chromium
	
	config2 := DefaultBrowserConfig()
	config2.Engine = Firefox
	
	session1, err := service.LaunchBrowser(ctx, config1)
	if err != nil {
		t.Fatalf("Failed to launch chromium browser: %v", err)
	}
	
	session2, err := service.LaunchBrowser(ctx, config2)
	if err != nil {
		t.Logf("Firefox may not be available, skipping: %v", err)
	} else {
		err = service.NavigateToURL(ctx, session1.ID, "https://example.com", "load", 30*time.Second)
		if err != nil {
			t.Fatalf("Failed to navigate in session 1: %v", err)
		}
		
		err = service.NavigateToURL(ctx, session2.ID, "https://httpbin.org", "load", 30*time.Second)
		if err != nil {
			t.Fatalf("Failed to navigate in session 2: %v", err)
		}
		
		err = service.CloseBrowser(ctx, session2.ID)
		if err != nil {
			t.Fatalf("Failed to close browser session 2: %v", err)
		}
	}
	
	err = service.CloseBrowser(ctx, session1.ID)
	if err != nil {
		t.Fatalf("Failed to close browser session 1: %v", err)
	}
}

func TestPlaywrightFormFilling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	if os.Getenv("CI") == "true" && os.Getenv("PLAYWRIGHT_BROWSERS_INSTALLED") != "true" {
		t.Skip("Skipping integration test in CI without browser installation")
	}
	
	logger := zap.NewNop()
	cfg := &config.Config{}
	
	service, err := NewPlaywrightService(logger, cfg)
	if err != nil {
		t.Fatalf("Failed to create playwright service: %v", err)
	}
	defer func() {
		if shutdownErr := service.Shutdown(context.Background()); shutdownErr != nil {
			t.Logf("Failed to shutdown service: %v", shutdownErr)
		}
	}()
	
	ctx := context.Background()
	
	config := DefaultBrowserConfig()
	session, err := service.LaunchBrowser(ctx, config)
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}
	defer func() {
		if closeErr := service.CloseBrowser(ctx, session.ID); closeErr != nil {
			t.Logf("Failed to close browser: %v", closeErr)
		}
	}()
	
	err = service.NavigateToURL(ctx, session.ID, "https://httpbin.org/forms/post", "load", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to navigate to forms page: %v", err)
	}
	
	fields := []map[string]interface{}{
		{
			"selector": "input[name='custname']",
			"value":    "Test User",
			"type":     "text",
		},
		{
			"selector": "input[name='custtel']",
			"value":    "123-456-7890",
			"type":     "text",
		},
		{
			"selector": "input[name='custemail']",
			"value":    "test@example.com",
			"type":     "text",
		},
	}
	
	err = service.FillForm(ctx, session.ID, fields, false, "")
	if err != nil {
		t.Fatalf("Failed to fill form: %v", err)
	}
	
	clickOptions := map[string]interface{}{
		"timeout": 10000,
		"force":   false,
	}
	
	err = service.ClickElement(ctx, session.ID, "input[type='submit']", clickOptions)
	if err != nil {
		t.Fatalf("Failed to click submit button: %v", err)
	}
}
