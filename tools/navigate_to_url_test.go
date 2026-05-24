package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	zaptest "go.uber.org/zap/zaptest"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

func TestNavigateToURLTool_NavigateToURLHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	session := &playwright.BrowserSession{
		ID:       "test-session-123",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.NavigateToURLReturns(nil)

	tool := &NavigateToURLTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name        string
		args        map[string]any
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful navigation with default parameters",
			args: map[string]any{
				"url": "https://example.com",
			},
			expectError: false,
		},
		{
			name: "successful navigation with all parameters",
			args: map[string]any{
				"url":        "https://example.com",
				"wait_until": "domcontentloaded",
				"timeout":    5000,
			},
			expectError: false,
		},
		{
			name: "successful navigation with URL normalization",
			args: map[string]any{
				"url": "example.com",
			},
			expectError: false,
		},
		{
			name: "missing URL parameter",
			args: map[string]any{
				"wait_until": "load",
			},
			expectError: true,
			errorMsg:    "url parameter is required",
		},
		{
			name: "empty URL parameter",
			args: map[string]any{
				"url": "",
			},
			expectError: true,
			errorMsg:    "non-empty string",
		},
		{
			name: "invalid URL parameter type",
			args: map[string]any{
				"url": 123,
			},
			expectError: true,
			errorMsg:    "must be a string",
		},
		{
			name: "invalid wait_until parameter",
			args: map[string]any{
				"url":        "https://example.com",
				"wait_until": "invalid",
			},
			expectError: true,
			errorMsg:    "invalid wait_until value",
		},
		{
			name: "invalid URL scheme",
			args: map[string]any{
				"url": "ftp://example.com",
			},
			expectError: true,
			errorMsg:    "unsupported URL scheme",
		},
		{
			name: "malformed URL",
			args: map[string]any{
				"url": "ht tp://example.com",
			},
			expectError: true,
			errorMsg:    "invalid URL format",
		},
		{
			name: "timeout parameter as float",
			args: map[string]any{
				"url":     "https://example.com",
				"timeout": 5000.0,
			},
			expectError: false,
		},
		{
			name: "negative timeout is rejected with explicit error",
			args: map[string]any{
				"url":     "https://example.com",
				"timeout": -1000,
			},
			expectError: true,
			errorMsg:    "timeout must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := tool.NavigateToURLHandler(ctx, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
			if result == "" {
				t.Fatalf("expected non-empty result")
			}
			var parsed map[string]any
			if err := json.Unmarshal([]byte(result), &parsed); err != nil {
				t.Fatalf("expected valid JSON response, got %q: %v", result, err)
			}
			if got, _ := parsed["success"].(bool); !got {
				t.Errorf("expected success=true in response, got %v", parsed["success"])
			}
		})
	}
}

func TestNavigateToURLTool_validateAndNormalizeURL(t *testing.T) {
	logger := zaptest.NewLogger(t)
	tool := &NavigateToURLTool{logger: logger}

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:     "valid HTTPS URL",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "valid HTTP URL",
			input:    "http://example.com",
			expected: "http://example.com",
		},
		{
			name:     "URL without scheme - should add HTTPS",
			input:    "example.com",
			expected: "https://example.com",
		},
		{
			name:     "URL with path",
			input:    "https://example.com/path",
			expected: "https://example.com/path",
		},
		{
			name:        "empty URL falls through to host check",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid scheme",
			input:       "ftp://example.com",
			expectError: true,
		},
		{
			name:        "malformed URL",
			input:       "ht tp://example.com",
			expectError: true,
		},
		{
			name:        "URL without host",
			input:       "https://",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.validateAndNormalizeURL(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}
