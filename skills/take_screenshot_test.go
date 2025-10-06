package skills

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	server "github.com/inference-gateway/adk/server"
	types "github.com/inference-gateway/adk/types"
	config "github.com/inference-gateway/browser-agent/config"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"
	zap "go.uber.org/zap"
)

func createTestSkill() *TakeScreenshotSkill {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	session := &playwright.BrowserSession{
		ID:       "test-session-123",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
	mockPlaywright.GetOrCreateDefaultSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.TakeScreenshotCalls(func(ctx context.Context, sessionID, path string, fullPage bool, selector string, format string, quality int) error {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		return os.WriteFile(path, []byte("mock screenshot data"), 0644)
	})
	mockPlaywright.GetConfigReturns(&config.Config{
		Browser: config.BrowserConfig{
			DataDir: "test_screenshots",
		},
	})

	return &TakeScreenshotSkill{
		logger:        logger,
		playwright:    mockPlaywright,
		screenshotDir: "test_screenshots",
	}
}

func TestTakeScreenshotHandler_BasicFunctionality(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{}

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
	ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{
		ID: "test-task-123",
	})

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

	resultFilename, ok := response["filename"].(string)
	if !ok {
		t.Errorf("Expected filename in response, got: %v", response["filename"])
	}

	if !strings.Contains(resultFilename, "viewport_") {
		t.Errorf("Expected viewport screenshot filename, got: %s", resultFilename)
	}

	if !strings.HasSuffix(resultFilename, ".png") {
		t.Errorf("Expected .png extension in filename, got: %s", resultFilename)
	}

	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_FullPageScreenshot(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"full_page": true,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
	ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{ID: "test-task-123"})

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

	resultFilename, ok := response["filename"].(string)
	if !ok || !strings.Contains(resultFilename, "fullpage_") {
		t.Errorf("Expected fullpage screenshot filename, got: %s", resultFilename)
	}

	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_JPEGWithQuality(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"type":    "jpeg",
		"quality": 95,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
	ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{ID: "test-task-123"})

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

	resultFilename, ok := response["filename"].(string)
	if !ok || !strings.HasSuffix(resultFilename, ".jpeg") {
		t.Errorf("Expected .jpeg extension in filename, got: %s", resultFilename)
	}

	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_ElementSelector(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"selector": "#main-content",
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
	ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{ID: "test-task-123"})

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

	resultFilename, ok := response["filename"].(string)
	if !ok || !strings.Contains(resultFilename, "element_") {
		t.Errorf("Expected element screenshot filename, got: %s", resultFilename)
	}

	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_DeterministicPath(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{}

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
	ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{ID: "test-task-123"})

	result, err := skill.TakeScreenshotHandler(ctx, args)

	if err != nil {
		t.Fatalf("Expected no error for deterministic path generation, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if _, ok := response["filename"]; !ok {
		t.Error("Expected filename to be generated in response")
	}

	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_InvalidImageType(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"type": "gif",
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
	ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{ID: "test-task-123"})

	_, err := skill.TakeScreenshotHandler(ctx, args)

	if err == nil {
		t.Error("Expected error for invalid image type, got nil")
	}

	expectedMsg := "invalid image type: gif. Must be 'png' or 'jpeg'"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got: %v", expectedMsg, err)
	}

	_ = os.RemoveAll("test_screenshots")
}

func TestTakeScreenshotHandler_InvalidQuality(t *testing.T) {
	skill := createTestSkill()

	args := map[string]any{
		"type":    "jpeg",
		"quality": 150,
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
	ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{ID: "test-task-123"})

	_, err := skill.TakeScreenshotHandler(ctx, args)

	if err == nil {
		t.Error("Expected error for invalid quality, got nil")
	}

	expectedMsg := "quality must be between 0 and 100 for JPEG images, got 150"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got: %v", expectedMsg, err)
	}

	_ = os.RemoveAll("test_screenshots")
}

func TestGenerateDeterministicPath(t *testing.T) {
	skill := createTestSkill()

	tests := []struct {
		name     string
		args     map[string]any
		expected string
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

			if !strings.HasPrefix(result, "test_screenshots/") {
				t.Errorf("Expected path to start with test_screenshots/, got: %s", result)
			}

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
