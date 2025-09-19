package skills

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/playwright-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// MockBrowserSession represents a mock browser session for testing
type MockBrowserSession struct {
	ID string
}

// MockPlaywright implements the BrowserAutomation interface for testing
type MockPlaywright struct {
	sessions           map[string]*playwright.BrowserSession
	screenshotResults  map[string]error
	screenshotData     []byte
	shouldFailLaunch   bool
	shouldFailCapture  bool
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

	// Create mock screenshot file
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, m.screenshotData, 0644)
}

// Implement other required methods with no-ops for testing
func (m *MockPlaywright) CloseBrowser(ctx context.Context, sessionID string) error   { return nil }
func (m *MockPlaywright) GetSession(sessionID string) (*playwright.BrowserSession, error) {
	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}
func (m *MockPlaywright) NavigateToURL(ctx context.Context, sessionID, url string, waitUntil string, timeout time.Duration) error { return nil }
func (m *MockPlaywright) ClickElement(ctx context.Context, sessionID, selector string, options map[string]interface{}) error { return nil }
func (m *MockPlaywright) FillForm(ctx context.Context, sessionID string, fields []map[string]interface{}, submit bool, submitSelector string) error { return nil }
func (m *MockPlaywright) ExtractData(ctx context.Context, sessionID string, extractors []map[string]interface{}, format string) (string, error) { return "", nil }
func (m *MockPlaywright) ExecuteScript(ctx context.Context, sessionID, script string, args []interface{}) (interface{}, error) { return nil, nil }
func (m *MockPlaywright) WaitForCondition(ctx context.Context, sessionID, condition, selector, state string, timeout time.Duration, customFunction string) error { return nil }
func (m *MockPlaywright) HandleAuthentication(ctx context.Context, sessionID, authType, username, password, loginURL string, selectors map[string]string) error { return nil }
func (m *MockPlaywright) GetHealth(ctx context.Context) error { return nil }
func (m *MockPlaywright) Shutdown(ctx context.Context) error { return nil }

func createTestSkill() *TakeScreenshotSkill {
	logger := zap.NewNop()
	mockPlaywright := NewMockPlaywright()
	
	return &TakeScreenshotSkill{
		logger:         logger,
		playwright:     mockPlaywright,
		artifactHelper: server.NewArtifactHelper(),
	}
}


func TestTakeScreenshotHandler_BasicFunctionality(t *testing.T) {
	skill := createTestSkill()
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test_screenshot.png")
	
	args := map[string]any{
		"path": path,
	}
	
	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Parse the JSON response
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	
	// Verify response structure
	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got: %v", response["success"])
	}
	
	if resultPath, ok := response["path"].(string); !ok || resultPath != path {
		t.Errorf("Expected path to be %s, got: %v", path, response["path"])
	}
	
	// Verify file was created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected screenshot file to be created at %s", path)
	}
}

func TestTakeScreenshotHandler_FullPageScreenshot(t *testing.T) {
	skill := createTestSkill()
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "fullpage_screenshot.png")
	
	args := map[string]any{
		"path":      path,
		"full_page": true,
	}
	
	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	
	if fullPage, ok := response["full_page"].(bool); !ok || !fullPage {
		t.Errorf("Expected full_page to be true, got: %v", response["full_page"])
	}
}

func TestTakeScreenshotHandler_JPEGWithQuality(t *testing.T) {
	skill := createTestSkill()
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "quality_screenshot.jpg")
	
	args := map[string]any{
		"path":    path,
		"type":    "jpeg",
		"quality": 95,
	}
	
	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	
	if imageType, ok := response["type"].(string); !ok || imageType != "jpeg" {
		t.Errorf("Expected type to be jpeg, got: %v", response["type"])
	}
	
	if quality, ok := response["quality"].(float64); !ok || int(quality) != 95 {
		t.Errorf("Expected quality to be 95, got: %v", response["quality"])
	}
}

func TestTakeScreenshotHandler_ElementSelector(t *testing.T) {
	skill := createTestSkill()
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "element_screenshot.png")
	
	args := map[string]any{
		"path":     path,
		"selector": "#main-content",
	}
	
	ctx := context.Background()
	result, err := skill.TakeScreenshotHandler(ctx, args)
	
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	
	if selector, ok := response["selector"].(string); !ok || selector != "#main-content" {
		t.Errorf("Expected selector to be #main-content, got: %v", response["selector"])
	}
}

func TestTakeScreenshotHandler_InvalidPath(t *testing.T) {
	skill := createTestSkill()
	
	args := map[string]any{
		"path": "",
	}
	
	ctx := context.Background()
	_, err := skill.TakeScreenshotHandler(ctx, args)
	
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
	
	if err.Error() != "path parameter is required and must be a non-empty string" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestTakeScreenshotHandler_InvalidImageType(t *testing.T) {
	skill := createTestSkill()
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.gif")
	
	args := map[string]any{
		"path": path,
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
}

func TestTakeScreenshotHandler_InvalidQuality(t *testing.T) {
	skill := createTestSkill()
	
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.jpg")
	
	args := map[string]any{
		"path":    path,
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
}

func TestValidateAndNormalizePath(t *testing.T) {
	skill := createTestSkill()
	
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errMsg   string
	}{
		{
			name:    "valid relative path",
			input:   "screenshots/test.png",
			wantErr: false,
		},
		{
			name:    "valid absolute path",
			input:   "/tmp/screenshots/test.png",
			wantErr: false,
		},
		{
			name:    "empty path",
			input:   "",
			wantErr: true,
			errMsg:  "path cannot be empty",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := skill.validateAndNormalizePath(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got: %v", tt.errMsg, err)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}
			
			if result == "" {
				t.Error("Expected non-empty result")
			}
			
			// Clean up created directories
			if dir := filepath.Dir(result); dir != "." {
				_ = os.RemoveAll(dir)
			}
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
		{"unknown", "image/png"}, // Default case
		{"", "image/png"},         // Default case
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
	
	// Verify it's a valid RFC3339 timestamp
	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Errorf("Expected valid RFC3339 timestamp, got: %s (error: %v)", timestamp, err)
	}
}