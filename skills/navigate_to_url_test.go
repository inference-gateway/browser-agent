package skills

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/inference-gateway/playwright-agent/internal/playwright"
	"github.com/inference-gateway/playwright-agent/internal/playwright/mocks"
	"go.uber.org/zap/zaptest"
)

func TestNavigateToURLSkill_NavigateToURLHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	session := &playwright.BrowserSession{
		ID:       "test-session-123",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
	mockPlaywright.GetOrCreateDefaultSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.NavigateToURLReturns(nil)

	skill := &NavigateToURLSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name        string
		args        map[string]any
		expectError bool
		errorMsg    string
		setup       func()
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
			errorMsg:    "url parameter is required",
		},
		{
			name: "invalid URL parameter type",
			args: map[string]any{
				"url": 123,
			},
			expectError: true,
			errorMsg:    "url parameter is required",
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
			name: "negative timeout parameter",
			args: map[string]any{
				"url":     "https://example.com",
				"timeout": -1000,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			ctx := context.Background()
			result, err := skill.NavigateToURLHandler(ctx, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result == "" {
					t.Errorf("expected non-empty result")
				}
			}
		})
	}
}

func TestNavigateToURLSkill_validateAndNormalizeURL(t *testing.T) {
	logger := zaptest.NewLogger(t)
	skill := &NavigateToURLSkill{logger: logger}

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
			name:        "empty URL",
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
			result, err := skill.validateAndNormalizeURL(tt.input)

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

func TestNavigateToURLSkill_isValidWaitCondition(t *testing.T) {
	logger := zaptest.NewLogger(t)
	skill := &NavigateToURLSkill{logger: logger}

	tests := []struct {
		condition string
		expected  bool
	}{
		{"domcontentloaded", true},
		{"load", true},
		{"networkidle", true},
		{"invalid", false},
		{"", false},
		{"LOAD", false},
	}

	for _, tt := range tests {
		t.Run(tt.condition, func(t *testing.T) {
			result := skill.isValidWaitCondition(tt.condition)
			if result != tt.expected {
				t.Errorf("expected %v for condition %q, got %v", tt.expected, tt.condition, result)
			}
		})
	}
}
