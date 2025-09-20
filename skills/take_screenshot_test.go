package skills

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	server "github.com/inference-gateway/adk/server"
	config "github.com/inference-gateway/playwright-agent/config"
	playwright "github.com/inference-gateway/playwright-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// MockBrowserSession represents a mock browser session for testing
type MockBrowserSession struct {
	ID string
}

// MockPlaywright implements the BrowserAutomation interface for testing
type MockPlaywright struct {
	sessions          map[string]*playwright.BrowserSession
	screenshotResults map[string]error
	screenshotData    []byte
	shouldFailLaunch  bool
	shouldFailCapture bool
}

func NewMockPlaywright() *MockPlaywright {
	return &MockPlaywright{
		sessions:          make(map[string]*playwright.BrowserSession),
		screenshotResults: make(map[string]error),
		screenshotData:    []byte("mock screenshot data"),
	}
}

func (m *MockPlaywright) LaunchBrowser(ctx context.Context, config *playwright.BrowserConfig) (*playwright.BrowserSession, error) {
	if m.shouldFailLaunch {
		return nil, fmt.Errorf("mock browser launch failed")
	}

	session := &playwright.BrowserSession{
		ID:       fmt.Sprintf("mock_session_%d", time.Now().UnixNano()),
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	m.sessions[session.ID] = session
	return session, nil
}

func (m *MockPlaywright) TakeScreenshot(ctx context.Context, sessionID, path string, fullPage bool, selector string, format string, quality int) error {
	if m.shouldFailCapture {
		return fmt.Errorf("mock screenshot capture failed")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, m.screenshotData, 0644)
}

// Implement other required methods with no-ops for testing
func (m *MockPlaywright) CloseBrowser(ctx context.Context, sessionID string) error { return nil }
func (m *MockPlaywright) GetSession(sessionID string) (*playwright.BrowserSession, error) {
	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}
func (m *MockPlaywright) NavigateToURL(ctx context.Context, sessionID, url string, waitUntil string, timeout time.Duration) error {
	return nil
}
func (m *MockPlaywright) ClickElement(ctx context.Context, sessionID, selector string, options map[string]any) error {
	return nil
}
func (m *MockPlaywright) FillForm(ctx context.Context, sessionID string, fields []map[string]any, submit bool, submitSelector string) error {
	return nil
}
func (m *MockPlaywright) ExtractData(ctx context.Context, sessionID string, extractors []map[string]any, format string) (string, error) {
	return "", nil
}
func (m *MockPlaywright) ExecuteScript(ctx context.Context, sessionID, script string, args []any) (any, error) {
	return nil, nil
}
func (m *MockPlaywright) WaitForCondition(ctx context.Context, sessionID, condition, selector, state string, timeout time.Duration, customFunction string) error {
	return nil
}
func (m *MockPlaywright) HandleAuthentication(ctx context.Context, sessionID, authType, username, password, loginURL string, selectors map[string]string) error {
	return nil
}
func (m *MockPlaywright) GetHealth(ctx context.Context) error { return nil }
func (m *MockPlaywright) Shutdown(ctx context.Context) error  { return nil }

func createTestSkill() *TakeScreenshotSkill {
	logger := zap.NewNop()
	mockPlaywright := NewMockPlaywright()
	cfg := &config.Config{
		ScreenshotDir: "test_screenshots",
	}

	return &TakeScreenshotSkill{
		logger:         logger,
		playwright:     mockPlaywright,
		artifactHelper: server.NewArtifactHelper(),
		screenshotDir:  cfg.ScreenshotDir,
	}
}

func TestTakeScreenshotHandler_BasicFunctionality(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{}

	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got: %v", response["success"])
	}

	// Check that path is deterministically generated
	resultPath, ok := response["path"].(string)
	if !ok {
		t.Errorf("Expected path in response, got: %v", response["path"])
	}

	// Verify path contains expected directory and timestamp pattern
	if !filepath.IsAbs(resultPath) && !strings.HasPrefix(resultPath, "test_screenshots/") {
		t.Errorf("Expected path to start with test_screenshots/, got: %s", resultPath)
	}

	if !strings.Contains(resultPath, "viewport_") {
		t.Errorf("Expected viewport screenshot filename, got: %s", resultPath)
	}

	// Verify file was actually created
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		t.Errorf("Expected screenshot file to be created at %s", resultPath)
	}

	// Clean up
	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_FullPageScreenshot(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"full_page": true,
	}

	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if fullPage, ok := response["full_page"].(bool); !ok || !fullPage {
		t.Errorf("Expected full_page to be true, got: %v", response["full_page"])
	}

	// Check that filename indicates fullpage screenshot
	resultPath, ok := response["path"].(string)
	if !ok || !strings.Contains(resultPath, "fullpage_") {
		t.Errorf("Expected fullpage screenshot filename, got: %s", resultPath)
	}

	// Clean up
	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_JPEGWithQuality(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"type":    "jpeg",
		"quality": 95,
	}

	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if imageType, ok := response["type"].(string); !ok || imageType != "jpeg" {
		t.Errorf("Expected type to be jpeg, got: %v", response["type"])
	}

	if quality, ok := response["quality"].(float64); !ok || int(quality) != 95 {
		t.Errorf("Expected quality to be 95, got: %v", response["quality"])
	}

	// Check that filename has .jpeg extension
	resultPath, ok := response["path"].(string)
	if !ok || !strings.HasSuffix(resultPath, ".jpeg") {
		t.Errorf("Expected .jpeg extension in filename, got: %s", resultPath)
	}

	// Clean up
	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_ElementSelector(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"selector": "#main-content",
	}

	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if selector, ok := response["selector"].(string); !ok || selector != "#main-content" {
		t.Errorf("Expected selector to be #main-content, got: %v", response["selector"])
	}

	// Check that filename indicates element screenshot
	resultPath, ok := response["path"].(string)
	if !ok || !strings.Contains(resultPath, "element_") {
		t.Errorf("Expected element screenshot filename, got: %s", resultPath)
	}

	// Clean up
	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_DeterministicPath(t *testing.T) {
	skill := createTestSkill()

	// Test that empty args work (no path required)
	args := map[string]any{}

	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)

	if err != nil {
		t.Fatalf("Expected no error for deterministic path generation, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify path is generated
	if _, ok := response["path"]; !ok {
		t.Error("Expected path to be generated in response")
	}

	// Clean up
	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_InvalidImageType(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"type": "gif",
	}

	ctx := context.Background()
	_, err := skill.TakeScreenshotHandler(ctx, args)

	if err == nil {
		t.Error("Expected error for invalid image type, got nil")
	}

	expectedMsg := "invalid image type: gif. Must be 'png' or 'jpeg'"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got: %v", expectedMsg, err)
	}

	// Clean up
	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_InvalidQuality(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"type":    "jpeg",
		"quality": 150,
	}

	ctx := context.Background()
	_, err := skill.TakeScreenshotHandler(ctx, args)

	if err == nil {
		t.Error("Expected error for invalid quality, got nil")
	}

	expectedMsg := "quality must be between 0 and 100 for JPEG images, got 150"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got: %v", expectedMsg, err)
	}

	// Clean up
	_ = os.RemoveAll("test_screenshots")
}

func TestGenerateDeterministicPath(t *testing.T) {
	skill := createTestSkill()

	tests := []struct {
		name     string
		args     map[string]any
		expected string // substring that should be in the path
	}{
		{
			name:     "viewport screenshot",
			args:     map[string]any{},
			expected: "viewport_",
		},
		{
			name:     "fullpage screenshot",
			args:     map[string]any{"full_page": true},
			expected: "fullpage_",
		},
		{
			name:     "element screenshot",
			args:     map[string]any{"selector": "#main"},
			expected: "element_",
		},
		{
			name:     "jpeg format",
			args:     map[string]any{"type": "jpeg"},
			expected: ".jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := skill.generateDeterministicPath(tt.args)

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if result == "" {
				t.Error("Expected non-empty result")
				return
			}

			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected path to contain '%s', got: %s", tt.expected, result)
			}

			// Verify path starts with screenshot directory
			if !strings.HasPrefix(result, "test_screenshots/") {
				t.Errorf("Expected path to start with test_screenshots/, got: %s", result)
			}

			// Clean up created directory
			_ = os.RemoveAll("test_screenshots")
		})
	}
}

func TestIsValidImageType(t *testing.T) {
	skill := createTestSkill()

	tests := []struct {
		imageType string
		expected  bool
	}{
		{"png", true},
		{"jpeg", true},
		{"jpg", false},
		{"gif", false},
		{"webp", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.imageType, func(t *testing.T) {
			result := skill.isValidImageType(tt.imageType)
			if result != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.imageType, result)
			}
		})
	}
}

func TestGetMimeType(t *testing.T) {
	skill := createTestSkill()

	tests := []struct {
		imageType string
		expected  string
	}{
		{"png", "image/png"},
		{"jpeg", "image/jpeg"},
		{"unknown", "image/png"},
		{"", "image/png"},
	}

	for _, tt := range tests {
		t.Run(tt.imageType, func(t *testing.T) {
			result := skill.getMimeType(tt.imageType)
			if result != tt.expected {
				t.Errorf("Expected %s for %s, got %s", tt.expected, tt.imageType, result)
			}
		})
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	skill := createTestSkill()

	timestamp := skill.getCurrentTimestamp()

	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Errorf("Expected valid RFC3339 timestamp, got: %s (error: %v)", timestamp, err)
	}
}
