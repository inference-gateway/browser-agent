package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	zap "go.uber.org/zap"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"

	config "github.com/inference-gateway/browser-agent/config"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

// newTestTool returns a TakeScreenshotTool wired with a mock playwright
// service that writes a stub file to a per-test temporary directory.
// Using t.TempDir() instead of a hardcoded "test_screenshots" directory
// avoids parallel-test collisions and removes the need for explicit
// cleanup at the end of each test (the harness deletes the dir).
func newTestTool(t *testing.T) *TakeScreenshotTool {
	t.Helper()
	dir := t.TempDir()

	mockPlaywright := &mocks.FakeBrowserAutomation{}

	session := &playwright.BrowserSession{
		ID:       "test-session-123",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.TakeScreenshotCalls(func(ctx context.Context, sessionID, path string, fullPage bool, selector string, format string, quality int) error {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		return os.WriteFile(path, []byte("mock screenshot data"), 0644)
	})
	mockPlaywright.GetConfigReturns(&config.Config{
		Browser: config.BrowserConfig{
			DataDir: dir,
		},
	})

	return &TakeScreenshotTool{
		logger:        zap.NewNop(),
		playwright:    mockPlaywright,
		screenshotDir: dir,
	}
}

func TestTakeScreenshotHandler_BasicFunctionality(t *testing.T) {
	tool := newTestTool(t)

	result, err := tool.TakeScreenshotHandler(context.Background(), map[string]any{})
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

	resultPath, _ := response["path"].(string)
	if !strings.Contains(resultPath, "viewport_") {
		t.Errorf("Expected viewport screenshot path, got: %s", resultPath)
	}

	resultFilename, _ := response["filename"].(string)
	if !strings.HasSuffix(resultFilename, ".png") {
		t.Errorf("Expected .png extension in filename, got: %s", resultFilename)
	}
}

func TestTakeScreenshotHandler_FullPageScreenshot(t *testing.T) {
	tool := newTestTool(t)

	result, err := tool.TakeScreenshotHandler(context.Background(), map[string]any{"full_page": true})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if fullPage, _ := response["full_page"].(bool); !fullPage {
		t.Errorf("Expected full_page to be true, got: %v", response["full_page"])
	}

	if resultFilename, _ := response["filename"].(string); !strings.Contains(resultFilename, "fullpage_") {
		t.Errorf("Expected fullpage screenshot filename, got: %s", resultFilename)
	}
}

func TestTakeScreenshotHandler_JPEGWithQuality(t *testing.T) {
	tool := newTestTool(t)

	result, err := tool.TakeScreenshotHandler(context.Background(), map[string]any{
		"type":    "jpeg",
		"quality": 95,
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if imageType, _ := response["type"].(string); imageType != "jpeg" {
		t.Errorf("Expected type to be jpeg, got: %v", response["type"])
	}

	if quality, _ := response["quality"].(float64); int(quality) != 95 {
		t.Errorf("Expected quality to be 95, got: %v", response["quality"])
	}

	if resultFilename, _ := response["filename"].(string); !strings.HasSuffix(resultFilename, ".jpeg") {
		t.Errorf("Expected .jpeg extension in filename, got: %s", resultFilename)
	}
}

func TestTakeScreenshotHandler_ElementSelector(t *testing.T) {
	tool := newTestTool(t)

	result, err := tool.TakeScreenshotHandler(context.Background(), map[string]any{
		"selector": "#main-content",
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if selector, _ := response["selector"].(string); selector != "#main-content" {
		t.Errorf("Expected selector to be #main-content, got: %v", response["selector"])
	}

	if resultFilename, _ := response["filename"].(string); !strings.Contains(resultFilename, "element_") {
		t.Errorf("Expected element screenshot filename, got: %s", resultFilename)
	}
}

func TestTakeScreenshotHandler_DeterministicPath(t *testing.T) {
	tool := newTestTool(t)

	result, err := tool.TakeScreenshotHandler(context.Background(), map[string]any{})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if _, ok := response["filename"]; !ok {
		t.Error("Expected filename to be generated in response")
	}
}

func TestTakeScreenshotHandler_InvalidImageType(t *testing.T) {
	tool := newTestTool(t)
	_, err := tool.TakeScreenshotHandler(context.Background(), map[string]any{"type": "gif"})
	if err == nil {
		t.Fatal("Expected error for invalid image type, got nil")
	}
	if !strings.Contains(err.Error(), "invalid image type") {
		t.Errorf("Expected invalid image type error, got: %v", err)
	}
}

func TestTakeScreenshotHandler_InvalidQuality(t *testing.T) {
	tool := newTestTool(t)
	_, err := tool.TakeScreenshotHandler(context.Background(), map[string]any{
		"type":    "jpeg",
		"quality": 150,
	})
	if err == nil {
		t.Fatal("Expected error for invalid quality, got nil")
	}
	if !strings.Contains(err.Error(), "quality must be between") {
		t.Errorf("Expected quality bound error, got: %v", err)
	}
}

func TestGenerateDeterministicPath(t *testing.T) {
	tool := newTestTool(t)

	tests := []struct {
		name      string
		fullPage  bool
		selector  string
		imageType string
		expected  string
	}{
		{name: "viewport screenshot", imageType: "png", expected: "viewport_"},
		{name: "fullpage screenshot", fullPage: true, imageType: "png", expected: "fullpage_"},
		{name: "element screenshot", selector: "#main", imageType: "png", expected: "element_"},
		{name: "jpeg format", imageType: "jpeg", expected: ".jpeg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.generateDeterministicPath(tt.fullPage, tt.selector, tt.imageType)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected path to contain %q, got: %s", tt.expected, result)
			}
		})
	}
}

// TestTruncateRunes_MultibyteSafe confirms that the truncation helper
// does not split multibyte UTF-8 sequences. Slicing the raw byte string
// would produce invalid UTF-8 in the filename.
func TestTruncateRunes_MultibyteSafe(t *testing.T) {
	in := "日本語テスト文字列abcdef" // 13 runes, > 20 bytes
	out := truncateRunes(in, 8)
	if len([]rune(out)) != 8 {
		t.Fatalf("expected 8 runes, got %d (output=%q)", len([]rune(out)), out)
	}
	if out != "日本語テスト文字" {
		t.Errorf("unexpected truncation result: %q", out)
	}

	// Short string should be returned unchanged.
	if got := truncateRunes("abc", 10); got != "abc" {
		t.Errorf("expected pass-through, got %q", got)
	}
}

func TestGetMimeType(t *testing.T) {
	tool := newTestTool(t)

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
			if result := tool.getMimeType(tt.imageType); result != tt.expected {
				t.Errorf("Expected %s for %s, got %s", tt.expected, tt.imageType, result)
			}
		})
	}
}

func TestGetCurrentTimestamp(t *testing.T) {
	tool := newTestTool(t)
	timestamp := tool.getCurrentTimestamp()
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		t.Errorf("Expected valid RFC3339 timestamp, got: %s (error: %v)", timestamp, err)
	}
}
