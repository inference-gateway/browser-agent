package playwright_test

import (
	"context"
	"testing"
	"time"

	"github.com/inference-gateway/playwright-agent/config"
	"github.com/inference-gateway/playwright-agent/internal/playwright"
	"github.com/inference-gateway/playwright-agent/internal/playwright/mocks"
	"go.uber.org/zap"
)

func TestPlaywrightServiceWithMock(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T, mockService *mocks.FakeBrowserAutomation)
	}{
		{
			name: "LaunchBrowser returns session successfully",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				config := playwright.DefaultBrowserConfig()
				expectedSession := &playwright.BrowserSession{
					ID:       "test-session-id",
					Created:  time.Now(),
					LastUsed: time.Now(),
				}

				mockService.LaunchBrowserReturns(expectedSession, nil)

				session, err := mockService.LaunchBrowser(ctx, config)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if session.ID != expectedSession.ID {
					t.Errorf("Expected session ID %s, got %s", expectedSession.ID, session.ID)
				}

				if mockService.LaunchBrowserCallCount() != 1 {
					t.Errorf("Expected LaunchBrowser to be called once, got %d calls", mockService.LaunchBrowserCallCount())
				}

				argCtx, argConfig := mockService.LaunchBrowserArgsForCall(0)
				if argCtx != ctx {
					t.Error("Expected context to match")
				}
				if argConfig != config {
					t.Error("Expected config to match")
				}
			},
		},
		{
			name: "NavigateToURL succeeds",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				sessionID := "test-session"
				url := "https://example.com"
				waitUntil := "load"
				timeout := 30 * time.Second

				mockService.NavigateToURLReturns(nil)

				err := mockService.NavigateToURL(ctx, sessionID, url, waitUntil, timeout)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if mockService.NavigateToURLCallCount() != 1 {
					t.Errorf("Expected NavigateToURL to be called once, got %d calls", mockService.NavigateToURLCallCount())
				}

				argCtx, argSessionID, argURL, argWaitUntil, argTimeout := mockService.NavigateToURLArgsForCall(0)
				if argCtx != ctx {
					t.Error("Expected context to match")
				}
				if argSessionID != sessionID {
					t.Errorf("Expected session ID %s, got %s", sessionID, argSessionID)
				}
				if argURL != url {
					t.Errorf("Expected URL %s, got %s", url, argURL)
				}
				if argWaitUntil != waitUntil {
					t.Errorf("Expected waitUntil %s, got %s", waitUntil, argWaitUntil)
				}
				if argTimeout != timeout {
					t.Errorf("Expected timeout %v, got %v", timeout, argTimeout)
				}
			},
		},
		{
			name: "ClickElement succeeds",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				sessionID := "test-session"
				selector := "button[type='submit']"
				options := map[string]any{
					"timeout": 10000,
					"force":   false,
				}

				mockService.ClickElementReturns(nil)

				err := mockService.ClickElement(ctx, sessionID, selector, options)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if mockService.ClickElementCallCount() != 1 {
					t.Errorf("Expected ClickElement to be called once, got %d calls", mockService.ClickElementCallCount())
				}

				argCtx, argSessionID, argSelector, argOptions := mockService.ClickElementArgsForCall(0)
				if argCtx != ctx {
					t.Error("Expected context to match")
				}
				if argSessionID != sessionID {
					t.Errorf("Expected session ID %s, got %s", sessionID, argSessionID)
				}
				if argSelector != selector {
					t.Errorf("Expected selector %s, got %s", selector, argSelector)
				}
				if len(argOptions) != len(options) {
					t.Errorf("Expected options length %d, got %d", len(options), len(argOptions))
				}
			},
		},
		{
			name: "FillForm succeeds",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				sessionID := "test-session"
				fields := []map[string]any{
					{
						"selector": "input[name='username']",
						"value":    "testuser",
						"type":     "text",
					},
					{
						"selector": "input[name='password']",
						"value":    "testpass",
						"type":     "password",
					},
				}
				submit := true
				submitSelector := "button[type='submit']"

				mockService.FillFormReturns(nil)

				err := mockService.FillForm(ctx, sessionID, fields, submit, submitSelector)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if mockService.FillFormCallCount() != 1 {
					t.Errorf("Expected FillForm to be called once, got %d calls", mockService.FillFormCallCount())
				}

				argCtx, argSessionID, argFields, argSubmit, argSubmitSelector := mockService.FillFormArgsForCall(0)
				if argCtx != ctx {
					t.Error("Expected context to match")
				}
				if argSessionID != sessionID {
					t.Errorf("Expected session ID %s, got %s", sessionID, argSessionID)
				}
				if len(argFields) != len(fields) {
					t.Errorf("Expected fields length %d, got %d", len(fields), len(argFields))
				}
				if argSubmit != submit {
					t.Errorf("Expected submit %v, got %v", submit, argSubmit)
				}
				if argSubmitSelector != submitSelector {
					t.Errorf("Expected submit selector %s, got %s", submitSelector, argSubmitSelector)
				}
			},
		},
		{
			name: "ExtractData returns expected data",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				sessionID := "test-session"
				extractors := []map[string]any{
					{
						"name":      "title",
						"selector":  "title",
						"attribute": "text",
					},
				}
				format := "json"
				expectedData := `{"title": "Example Domain"}`

				mockService.ExtractDataReturns(expectedData, nil)

				data, err := mockService.ExtractData(ctx, sessionID, extractors, format)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if data != expectedData {
					t.Errorf("Expected data %s, got %s", expectedData, data)
				}

				if mockService.ExtractDataCallCount() != 1 {
					t.Errorf("Expected ExtractData to be called once, got %d calls", mockService.ExtractDataCallCount())
				}
			},
		},
		{
			name: "TakeScreenshot succeeds",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				sessionID := "test-session"
				path := "/tmp/screenshot.png"
				fullPage := false
				selector := ""
				format := "png"
				quality := 80

				mockService.TakeScreenshotReturns(nil)

				err := mockService.TakeScreenshot(ctx, sessionID, path, fullPage, selector, format, quality)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if mockService.TakeScreenshotCallCount() != 1 {
					t.Errorf("Expected TakeScreenshot to be called once, got %d calls", mockService.TakeScreenshotCallCount())
				}

				argCtx, argSessionID, argPath, argFullPage, argSelector, argFormat, argQuality := mockService.TakeScreenshotArgsForCall(0)
				if argCtx != ctx {
					t.Error("Expected context to match")
				}
				if argSessionID != sessionID {
					t.Errorf("Expected session ID %s, got %s", sessionID, argSessionID)
				}
				if argPath != path {
					t.Errorf("Expected path %s, got %s", path, argPath)
				}
				if argFullPage != fullPage {
					t.Errorf("Expected fullPage %v, got %v", fullPage, argFullPage)
				}
				if argSelector != selector {
					t.Errorf("Expected selector %s, got %s", selector, argSelector)
				}
				if argFormat != format {
					t.Errorf("Expected format %s, got %s", format, argFormat)
				}
				if argQuality != quality {
					t.Errorf("Expected quality %d, got %d", quality, argQuality)
				}
			},
		},
		{
			name: "ExecuteScript returns expected result",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				sessionID := "test-session"
				script := "return document.title"
				args := []any{}
				expectedResult := "Example Domain"

				mockService.ExecuteScriptReturns(expectedResult, nil)

				result, err := mockService.ExecuteScript(ctx, sessionID, script, args)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if result != expectedResult {
					t.Errorf("Expected result %v, got %v", expectedResult, result)
				}

				if mockService.ExecuteScriptCallCount() != 1 {
					t.Errorf("Expected ExecuteScript to be called once, got %d calls", mockService.ExecuteScriptCallCount())
				}
			},
		},
		{
			name: "GetSession returns session successfully",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				sessionID := "test-session"
				expectedSession := &playwright.BrowserSession{
					ID:       sessionID,
					Created:  time.Now(),
					LastUsed: time.Now(),
				}

				mockService.GetSessionReturns(expectedSession, nil)

				session, err := mockService.GetSession(sessionID)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if session.ID != expectedSession.ID {
					t.Errorf("Expected session ID %s, got %s", expectedSession.ID, session.ID)
				}

				if mockService.GetSessionCallCount() != 1 {
					t.Errorf("Expected GetSession to be called once, got %d calls", mockService.GetSessionCallCount())
				}
			},
		},
		{
			name: "CloseBrowser succeeds",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()
				sessionID := "test-session"

				mockService.CloseBrowserReturns(nil)

				err := mockService.CloseBrowser(ctx, sessionID)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if mockService.CloseBrowserCallCount() != 1 {
					t.Errorf("Expected CloseBrowser to be called once, got %d calls", mockService.CloseBrowserCallCount())
				}
			},
		},
		{
			name: "GetHealth succeeds",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()

				mockService.GetHealthReturns(nil)

				err := mockService.GetHealth(ctx)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if mockService.GetHealthCallCount() != 1 {
					t.Errorf("Expected GetHealth to be called once, got %d calls", mockService.GetHealthCallCount())
				}
			},
		},
		{
			name: "Shutdown succeeds",
			testFunc: func(t *testing.T, mockService *mocks.FakeBrowserAutomation) {
				ctx := context.Background()

				mockService.ShutdownReturns(nil)

				err := mockService.Shutdown(ctx)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if mockService.ShutdownCallCount() != 1 {
					t.Errorf("Expected Shutdown to be called once, got %d calls", mockService.ShutdownCallCount())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.FakeBrowserAutomation{}
			tt.testFunc(t, mockService)
		})
	}
}

func TestPlaywrightServiceIntegrationWithMock(t *testing.T) {
	mockService := &mocks.FakeBrowserAutomation{}
	logger := zap.NewNop()
	cfg := &config.Config{}

	ctx := context.Background()

	session := &playwright.BrowserSession{
		ID:       "integration-test-session",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	mockService.LaunchBrowserReturns(session, nil)
	mockService.NavigateToURLReturns(nil)
	mockService.TakeScreenshotReturns(nil)
	mockService.ExecuteScriptReturns("Example Domain", nil)
	mockService.ExtractDataReturns(`{"title": "Example Domain"}`, nil)
	mockService.GetSessionReturns(session, nil)
	mockService.CloseBrowserReturns(nil)

	config := playwright.DefaultBrowserConfig()
	launchedSession, err := mockService.LaunchBrowser(ctx, config)
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}

	if launchedSession == nil {
		t.Fatal("Expected non-nil session")
	}

	if launchedSession.ID == "" {
		t.Fatal("Expected non-empty session ID")
	}

	err = mockService.NavigateToURL(ctx, session.ID, "https://example.com", "load", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to navigate to URL: %v", err)
	}

	screenshotPath := "/tmp/test-screenshot.png"
	err = mockService.TakeScreenshot(ctx, session.ID, screenshotPath, false, "", "png", 80)
	if err != nil {
		t.Fatalf("Failed to take screenshot: %v", err)
	}

	result, err := mockService.ExecuteScript(ctx, session.ID, "return document.title", nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil script result")
	}

	extractors := []map[string]any{
		{
			"name":      "title",
			"selector":  "title",
			"attribute": "text",
		},
	}

	data, err := mockService.ExtractData(ctx, session.ID, extractors, "json")
	if err != nil {
		t.Fatalf("Failed to extract data: %v", err)
	}

	if data == "" {
		t.Error("Expected non-empty extracted data")
	}

	retrievedSession, err := mockService.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve session: %v", err)
	}

	if retrievedSession.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, retrievedSession.ID)
	}

	err = mockService.CloseBrowser(ctx, session.ID)
	if err != nil {
		t.Fatalf("Failed to close browser: %v", err)
	}

	_ = logger
	_ = cfg
}

func TestPlaywrightServiceErrorHandlingWithMock(t *testing.T) {
	mockService := &mocks.FakeBrowserAutomation{}
	ctx := context.Background()

	t.Run("LaunchBrowser returns error", func(t *testing.T) {
		mockService.LaunchBrowserReturns(nil, &testError{"launch failed"})

		config := playwright.DefaultBrowserConfig()
		_, err := mockService.LaunchBrowser(ctx, config)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err.Error() != "launch failed" {
			t.Errorf("Expected error message 'launch failed', got %s", err.Error())
		}
	})

	t.Run("NavigateToURL returns error", func(t *testing.T) {
		mockService.NavigateToURLReturns(&testError{"navigation failed"})

		err := mockService.NavigateToURL(ctx, "session-id", "https://example.com", "load", 30*time.Second)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err.Error() != "navigation failed" {
			t.Errorf("Expected error message 'navigation failed', got %s", err.Error())
		}
	})

	t.Run("GetSession returns error for non-existent session", func(t *testing.T) {
		mockService.GetSessionReturns(nil, &testError{"session not found"})

		_, err := mockService.GetSession("non-existent-session")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err.Error() != "session not found" {
			t.Errorf("Expected error message 'session not found', got %s", err.Error())
		}
	})
}

func TestMultipleBrowserSessionsWithMock(t *testing.T) {
	mockService := &mocks.FakeBrowserAutomation{}

	ctx := context.Background()

	session1 := &playwright.BrowserSession{
		ID:       "chromium-session",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	session2 := &playwright.BrowserSession{
		ID:       "firefox-session",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	mockService.LaunchBrowserReturnsOnCall(0, session1, nil)
	mockService.LaunchBrowserReturnsOnCall(1, session2, nil)
	mockService.NavigateToURLReturns(nil)
	mockService.CloseBrowserReturns(nil)

	config1 := playwright.DefaultBrowserConfig()
	config1.Engine = playwright.Chromium

	config2 := playwright.DefaultBrowserConfig()
	config2.Engine = playwright.Firefox

	launchedSession1, err := mockService.LaunchBrowser(ctx, config1)
	if err != nil {
		t.Fatalf("Failed to launch chromium browser: %v", err)
	}

	if launchedSession1.ID != "chromium-session" {
		t.Errorf("Expected session ID 'chromium-session', got %s", launchedSession1.ID)
	}

	launchedSession2, err := mockService.LaunchBrowser(ctx, config2)
	if err != nil {
		t.Fatalf("Failed to launch firefox browser: %v", err)
	}

	if launchedSession2.ID != "firefox-session" {
		t.Errorf("Expected session ID 'firefox-session', got %s", launchedSession2.ID)
	}

	err = mockService.NavigateToURL(ctx, session1.ID, "https://example.com", "load", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to navigate in session 1: %v", err)
	}

	err = mockService.NavigateToURL(ctx, session2.ID, "https://httpbin.org", "load", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to navigate in session 2: %v", err)
	}

	err = mockService.CloseBrowser(ctx, session2.ID)
	if err != nil {
		t.Fatalf("Failed to close browser session 2: %v", err)
	}

	err = mockService.CloseBrowser(ctx, session1.ID)
	if err != nil {
		t.Fatalf("Failed to close browser session 1: %v", err)
	}

	if mockService.LaunchBrowserCallCount() != 2 {
		t.Errorf("Expected LaunchBrowser to be called twice, got %d calls", mockService.LaunchBrowserCallCount())
	}

	if mockService.NavigateToURLCallCount() != 2 {
		t.Errorf("Expected NavigateToURL to be called twice, got %d calls", mockService.NavigateToURLCallCount())
	}

	if mockService.CloseBrowserCallCount() != 2 {
		t.Errorf("Expected CloseBrowser to be called twice, got %d calls", mockService.CloseBrowserCallCount())
	}
}

func TestPlaywrightFormFillingWithMock(t *testing.T) {
	mockService := &mocks.FakeBrowserAutomation{}

	ctx := context.Background()

	session := &playwright.BrowserSession{
		ID:       "form-test-session",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	mockService.LaunchBrowserReturns(session, nil)
	mockService.NavigateToURLReturns(nil)
	mockService.FillFormReturns(nil)
	mockService.ClickElementReturns(nil)
	mockService.CloseBrowserReturns(nil)

	config := playwright.DefaultBrowserConfig()
	launchedSession, err := mockService.LaunchBrowser(ctx, config)
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}

	if launchedSession.ID != "form-test-session" {
		t.Errorf("Expected session ID 'form-test-session', got %s", launchedSession.ID)
	}

	err = mockService.NavigateToURL(ctx, session.ID, "https://httpbin.org/forms/post", "load", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to navigate to forms page: %v", err)
	}

	fields := []map[string]any{
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

	err = mockService.FillForm(ctx, session.ID, fields, false, "")
	if err != nil {
		t.Fatalf("Failed to fill form: %v", err)
	}

	clickOptions := map[string]any{
		"timeout": 10000,
		"force":   false,
	}

	err = mockService.ClickElement(ctx, session.ID, "input[type='submit']", clickOptions)
	if err != nil {
		t.Fatalf("Failed to click submit button: %v", err)
	}

	err = mockService.CloseBrowser(ctx, session.ID)
	if err != nil {
		t.Fatalf("Failed to close browser: %v", err)
	}

	if mockService.FillFormCallCount() != 1 {
		t.Errorf("Expected FillForm to be called once, got %d calls", mockService.FillFormCallCount())
	}

	argCtx, argSessionID, argFields, argSubmit, argSubmitSelector := mockService.FillFormArgsForCall(0)
	if argCtx != ctx {
		t.Error("Expected context to match")
	}
	if argSessionID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, argSessionID)
	}
	if len(argFields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(argFields))
	}
	if argSubmit != false {
		t.Errorf("Expected submit false, got %v", argSubmit)
	}
	if argSubmitSelector != "" {
		t.Errorf("Expected empty submit selector, got %s", argSubmitSelector)
	}

	if mockService.ClickElementCallCount() != 1 {
		t.Errorf("Expected ClickElement to be called once, got %d calls", mockService.ClickElementCallCount())
	}

	clickCtx, clickSessionID, clickSelector, clickOpts := mockService.ClickElementArgsForCall(0)
	if clickCtx != ctx {
		t.Error("Expected context to match")
	}
	if clickSessionID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, clickSessionID)
	}
	if clickSelector != "input[type='submit']" {
		t.Errorf("Expected selector 'input[type='submit']', got %s", clickSelector)
	}
	if len(clickOpts) != 2 {
		t.Errorf("Expected 2 click options, got %d", len(clickOpts))
	}
}

func TestFullWorkflowWithMock(t *testing.T) {
	mockService := &mocks.FakeBrowserAutomation{}

	ctx := context.Background()

	session := &playwright.BrowserSession{
		ID:       "workflow-test-session",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	mockService.LaunchBrowserReturns(session, nil)
	mockService.NavigateToURLReturns(nil)
	mockService.TakeScreenshotReturns(nil)
	mockService.ExecuteScriptReturns("Example Domain", nil)
	mockService.ExtractDataReturns(`{"title": "Example Domain"}`, nil)
	mockService.GetSessionReturns(session, nil)
	mockService.CloseBrowserReturns(nil)

	config := playwright.DefaultBrowserConfig()
	launchedSession, err := mockService.LaunchBrowser(ctx, config)
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}

	if launchedSession == nil {
		t.Fatal("Expected non-nil session")
	}

	if launchedSession.ID == "" {
		t.Fatal("Expected non-empty session ID")
	}

	err = mockService.NavigateToURL(ctx, session.ID, "https://example.com", "load", 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to navigate to URL: %v", err)
	}

	screenshotPath := "/tmp/test-screenshot.png"
	err = mockService.TakeScreenshot(ctx, session.ID, screenshotPath, false, "", "png", 80)
	if err != nil {
		t.Fatalf("Failed to take screenshot: %v", err)
	}

	result, err := mockService.ExecuteScript(ctx, session.ID, "return document.title", nil)
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil script result")
	}

	extractors := []map[string]any{
		{
			"name":      "title",
			"selector":  "title",
			"attribute": "text",
		},
	}

	data, err := mockService.ExtractData(ctx, session.ID, extractors, "json")
	if err != nil {
		t.Fatalf("Failed to extract data: %v", err)
	}

	if data == "" {
		t.Error("Expected non-empty extracted data")
	}

	retrievedSession, err := mockService.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve session: %v", err)
	}

	if retrievedSession.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, retrievedSession.ID)
	}

	err = mockService.CloseBrowser(ctx, session.ID)
	if err != nil {
		t.Fatalf("Failed to close browser: %v", err)
	}

	mockService.GetSessionReturns(nil, &testError{"session not found"})

	_, err = mockService.GetSession(session.ID)
	if err == nil {
		t.Error("Expected error when getting closed session")
	}

	if mockService.LaunchBrowserCallCount() != 1 {
		t.Errorf("Expected LaunchBrowser to be called once, got %d calls", mockService.LaunchBrowserCallCount())
	}

	if mockService.NavigateToURLCallCount() != 1 {
		t.Errorf("Expected NavigateToURL to be called once, got %d calls", mockService.NavigateToURLCallCount())
	}

	if mockService.TakeScreenshotCallCount() != 1 {
		t.Errorf("Expected TakeScreenshot to be called once, got %d calls", mockService.TakeScreenshotCallCount())
	}

	if mockService.ExecuteScriptCallCount() != 1 {
		t.Errorf("Expected ExecuteScript to be called once, got %d calls", mockService.ExecuteScriptCallCount())
	}

	if mockService.ExtractDataCallCount() != 1 {
		t.Errorf("Expected ExtractData to be called once, got %d calls", mockService.ExtractDataCallCount())
	}

	if mockService.GetSessionCallCount() != 2 {
		t.Errorf("Expected GetSession to be called twice, got %d calls", mockService.GetSessionCallCount())
	}

	if mockService.CloseBrowserCallCount() != 1 {
		t.Errorf("Expected CloseBrowser to be called once, got %d calls", mockService.CloseBrowserCallCount())
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
