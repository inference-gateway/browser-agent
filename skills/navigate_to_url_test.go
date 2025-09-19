package skills

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/inference-gateway/playwright-agent/internal/playwright"
	"go.uber.org/zap/zaptest"
)

// MockBrowserAutomation is a mock implementation of the BrowserAutomation interface
type MockBrowserAutomation struct {
	sessions    map[string]*playwright.BrowserSession
	navigateErr error
	launchErr   error
}

func NewMockBrowserAutomation() *MockBrowserAutomation {
	return &MockBrowserAutomation{
		sessions: make(map[string]*playwright.BrowserSession),
	}
}

func (m *MockBrowserAutomation) LaunchBrowser(ctx context.Context, config *playwright.BrowserConfig) (*playwright.BrowserSession, error) {
	if m.launchErr != nil {
		return nil, m.launchErr
	}

	session := &playwright.BrowserSession{
		ID:       "test-session-123",
		Created:  time.Now(),
		LastUsed: time.Now(),
	}
	m.sessions[session.ID] = session
	return session, nil
}

func (m *MockBrowserAutomation) CloseBrowser(ctx context.Context, sessionID string) error {
	delete(m.sessions, sessionID)
	return nil
}

func (m *MockBrowserAutomation) GetSession(sessionID string) (*playwright.BrowserSession, error) {
	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, nil
	}
	return session, nil
}

func (m *MockBrowserAutomation) NavigateToURL(ctx context.Context, sessionID, url string, waitUntil string, timeout time.Duration) error {
	return m.navigateErr
}

func (m *MockBrowserAutomation) ClickElement(ctx context.Context, sessionID, selector string, options map[string]any) error {
	return nil
}

func (m *MockBrowserAutomation) FillForm(ctx context.Context, sessionID string, fields []map[string]any, submit bool, submitSelector string) error {
	return nil
}

func (m *MockBrowserAutomation) ExtractData(ctx context.Context, sessionID string, extractors []map[string]any, format string) (string, error) {
	return "", nil
}

func (m *MockBrowserAutomation) TakeScreenshot(ctx context.Context, sessionID, path string, fullPage bool, selector string, format string, quality int) error {
	return nil
}

func (m *MockBrowserAutomation) ExecuteScript(ctx context.Context, sessionID, script string, args []any) (any, error) {
	return nil, nil
}

func (m *MockBrowserAutomation) WaitForCondition(ctx context.Context, sessionID, condition, selector, state string, timeout time.Duration, customFunction string) error {
	return nil
}

func (m *MockBrowserAutomation) HandleAuthentication(ctx context.Context, sessionID, authType, username, password, loginURL string, selectors map[string]string) error {
	return nil
}

func (m *MockBrowserAutomation) GetHealth(ctx context.Context) error {
	return nil
}

func (m *MockBrowserAutomation) Shutdown(ctx context.Context) error {
	return nil
}

func TestNavigateToURLSkill_NavigateToURLHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	mockPlaywright := NewMockBrowserAutomation()
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
