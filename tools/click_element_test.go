package tools

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	zap "go.uber.org/zap"

	assert "github.com/stretchr/testify/assert"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

func TestClickElementTool_ClickElementHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name          string
		args          map[string]any
		setupMock     func(*mocks.FakeBrowserAutomation)
		expectedError bool
		errorContains string
	}{
		{
			name: "successful basic click",
			args: map[string]any{
				"selector": "#submit-button",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{ID: "test-session"}
				m.GetOrCreateTaskSessionReturns(session, nil)
				m.GetSessionReturns(session, nil)
				m.WaitForConditionReturns(nil)
				m.ClickElementReturns(nil)
			},
			expectedError: false,
		},
		{
			name: "successful click with all parameters",
			args: map[string]any{
				"selector":    "button.primary",
				"button":      "right",
				"click_count": 2,
				"force":       true,
				"timeout":     5000,
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{ID: "test-session"}
				m.GetOrCreateTaskSessionReturns(session, nil)
				m.GetSessionReturns(session, nil)
				m.WaitForConditionReturns(nil)
				m.ClickElementReturns(nil)
			},
			expectedError: false,
		},
		{
			name: "xpath selector",
			args: map[string]any{
				"selector": "//button[@id='submit']",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{ID: "test-session"}
				m.GetOrCreateTaskSessionReturns(session, nil)
				m.GetSessionReturns(session, nil)
				m.WaitForConditionReturns(nil)
				m.ClickElementReturns(nil)
			},
			expectedError: false,
		},
		{
			name: "text-based selector with quotes",
			args: map[string]any{
				"selector": "'Click Me'",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{ID: "test-session"}
				m.GetOrCreateTaskSessionReturns(session, nil)
				m.GetSessionReturns(session, nil)
				m.WaitForConditionReturns(nil)
				m.ClickElementReturns(nil)
			},
			expectedError: false,
		},
		{
			name: "missing selector parameter",
			args: map[string]any{
				"button": "left",
			},
			setupMock:     func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
			errorContains: "selector parameter is required",
		},
		{
			name: "empty selector parameter",
			args: map[string]any{
				"selector": "",
			},
			setupMock:     func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
			errorContains: "non-empty string",
		},
		{
			name: "invalid button parameter",
			args: map[string]any{
				"selector": "#button",
				"button":   "invalid",
			},
			setupMock:     func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
			errorContains: "invalid button value",
		},
		{
			name: "negative timeout rejected",
			args: map[string]any{
				"selector": "#button",
				"timeout":  -1,
			},
			setupMock:     func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
			errorContains: "timeout must be between",
		},
		{
			name: "click_count must be positive",
			args: map[string]any{
				"selector":    "#button",
				"click_count": 0,
			},
			setupMock:     func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
			errorContains: "click_count must be between",
		},
		{
			name: "browser launch failure",
			args: map[string]any{
				"selector": "#button",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				m.GetOrCreateTaskSessionReturns(nil, errors.New("browser launch failed"))
			},
			expectedError: true,
		},
		{
			name: "element not found",
			args: map[string]any{
				"selector": "#nonexistent",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{ID: "test-session"}
				m.GetOrCreateTaskSessionReturns(session, nil)
				m.GetSessionReturns(session, nil)
				m.WaitForConditionReturns(errors.New("element not found"))
			},
			expectedError: true,
		},
		{
			name: "click operation failure",
			args: map[string]any{
				"selector": "#button",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{ID: "test-session"}
				m.GetOrCreateTaskSessionReturns(session, nil)
				m.GetSessionReturns(session, nil)
				m.WaitForConditionReturns(nil)
				m.ClickElementReturns(errors.New("click failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlaywright := &mocks.FakeBrowserAutomation{}
			tt.setupMock(mockPlaywright)

			tool := &ClickElementTool{
				logger:     logger,
				playwright: mockPlaywright,
			}

			result, err := tool.ClickElementHandler(context.Background(), tt.args)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Empty(t, result)
				return
			}

			assert.NoError(t, err)
			var parsed map[string]any
			assert.NoError(t, json.Unmarshal([]byte(result), &parsed), "response should be valid JSON")
			assert.Equal(t, true, parsed["success"])
		})
	}
}

// TestClickElementTool_ForceSkipsActionabilityWait verifies that force=true
// bypasses the visibility wait, which is the documented purpose of the flag.
// Before the fix the flag was accepted but ignored.
func TestClickElementTool_ForceSkipsActionabilityWait(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.ClickElementReturns(nil)

	tool := &ClickElementTool{logger: logger, playwright: mockPlaywright}
	_, err := tool.ClickElementHandler(context.Background(), map[string]any{
		"selector": "#hidden",
		"force":    true,
	})

	assert.NoError(t, err)
	assert.Equal(t, 0, mockPlaywright.WaitForConditionCallCount(),
		"WaitForCondition should be skipped when force=true")
	assert.Equal(t, 1, mockPlaywright.ClickElementCallCount(),
		"ClickElement should still be invoked")
}

func TestClickElementTool_normalizeSelector(t *testing.T) {
	tool := &ClickElementTool{}

	tests := []struct {
		name             string
		selector         string
		expectedSelector string
		expectedType     string
	}{
		{
			name:             "CSS selector",
			selector:         "#submit-button",
			expectedSelector: "#submit-button",
			expectedType:     "css",
		},
		{
			name:             "XPath with //",
			selector:         "//button[@id='submit']",
			expectedSelector: "//button[@id='submit']",
			expectedType:     "xpath",
		},
		{
			name:             "XPath with /",
			selector:         "/html/body/button",
			expectedSelector: "/html/body/button",
			expectedType:     "xpath",
		},
		{
			name:             "XPath with xpath= prefix",
			selector:         "xpath=//button",
			expectedSelector: "//button",
			expectedType:     "xpath",
		},
		{
			name:             "Text with single quotes",
			selector:         "'Click Me'",
			expectedSelector: "text=Click Me",
			expectedType:     "text",
		},
		{
			name:             "Text with double quotes",
			selector:         "\"Submit\"",
			expectedSelector: "text=Submit",
			expectedType:     "text",
		},
		{
			name:             "Text with text= prefix",
			selector:         "text=Submit Form",
			expectedSelector: "text=Submit Form",
			expectedType:     "text",
		},
		{
			name:             "Role selector",
			selector:         "role=button",
			expectedSelector: "role=button",
			expectedType:     "role",
		},
		{
			name:             "Role in CSS selector",
			selector:         "button[role=submit]",
			expectedSelector: "button[role=submit]",
			expectedType:     "role",
		},
		{
			name:             "Data testid selector",
			selector:         "[data-testid=submit-btn]",
			expectedSelector: "[data-testid=submit-btn]",
			expectedType:     "testid",
		},
		{
			name:             "Complex CSS selector",
			selector:         "div.container > button.primary:first-child",
			expectedSelector: "div.container > button.primary:first-child",
			expectedType:     "css",
		},
		{
			name:             "CSS data-text attribute is not misclassified as text",
			selector:         `[data-text="hello"]`,
			expectedSelector: `[data-text="hello"]`,
			expectedType:     "css",
		},
		{
			name:             "CSS attribute selector with text in name stays CSS",
			selector:         `input[name="search-text"]`,
			expectedSelector: `input[name="search-text"]`,
			expectedType:     "css",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector, selectorType := tool.normalizeSelector(tt.selector)
			assert.Equal(t, tt.expectedSelector, selector)
			assert.Equal(t, tt.expectedType, selectorType)
		})
	}
}

func TestClickElementTool_NewClickElementTool(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	tool := NewClickElementTool(logger, mockPlaywright)

	assert.NotNil(t, tool)
}
