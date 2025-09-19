package skills

import (
	"context"
	"errors"
	"testing"

	"github.com/inference-gateway/playwright-agent/internal/playwright"
	"github.com/inference-gateway/playwright-agent/internal/playwright/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestClickElementSkill_ClickElementHandler(t *testing.T) {
	logger := zap.NewNop()
	
	tests := []struct {
		name           string
		args           map[string]any
		setupMock      func(*mocks.FakeBrowserAutomation)
		expectedError  bool
		expectedResult string
	}{
		{
			name: "successful basic click",
			args: map[string]any{
				"selector": "#submit-button",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{
					ID: "test-session",
				}
				m.LaunchBrowserReturns(session, nil)
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
				session := &playwright.BrowserSession{
					ID: "test-session",
				}
				m.LaunchBrowserReturns(session, nil)
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
				session := &playwright.BrowserSession{
					ID: "test-session",
				}
				m.LaunchBrowserReturns(session, nil)
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
				session := &playwright.BrowserSession{
					ID: "test-session",
				}
				m.LaunchBrowserReturns(session, nil)
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
			setupMock:    func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
		},
		{
			name: "empty selector parameter",
			args: map[string]any{
				"selector": "",
			},
			setupMock:    func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
		},
		{
			name: "invalid button parameter",
			args: map[string]any{
				"selector": "#button",
				"button":   "invalid",
			},
			setupMock:    func(m *mocks.FakeBrowserAutomation) {},
			expectedError: true,
		},
		{
			name: "browser launch failure",
			args: map[string]any{
				"selector": "#button",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				m.LaunchBrowserReturns(nil, errors.New("browser launch failed"))
			},
			expectedError: true,
		},
		{
			name: "element not found",
			args: map[string]any{
				"selector": "#nonexistent",
			},
			setupMock: func(m *mocks.FakeBrowserAutomation) {
				session := &playwright.BrowserSession{
					ID: "test-session",
				}
				m.LaunchBrowserReturns(session, nil)
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
				session := &playwright.BrowserSession{
					ID: "test-session",
				}
				m.LaunchBrowserReturns(session, nil)
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

			skill := &ClickElementSkill{
				logger:     logger,
				playwright: mockPlaywright,
			}

			result, err := skill.ClickElementHandler(context.Background(), tt.args)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, result, "success:true")
			}
		})
	}
}

func TestClickElementSkill_isValidButton(t *testing.T) {
	skill := &ClickElementSkill{}

	tests := []struct {
		button   string
		expected bool
	}{
		{"left", true},
		{"right", true},
		{"middle", true},
		{"invalid", false},
		{"", false},
		{"LEFT", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.button, func(t *testing.T) {
			result := skill.isValidButton(tt.button)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClickElementSkill_normalizeSelector(t *testing.T) {
	skill := &ClickElementSkill{}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector, selectorType := skill.normalizeSelector(tt.selector)
			assert.Equal(t, tt.expectedSelector, selector)
			assert.Equal(t, tt.expectedType, selectorType)
		})
	}
}

func TestClickElementSkill_NewClickElementSkill(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	tool := NewClickElementSkill(logger, mockPlaywright)
	
	assert.NotNil(t, tool)
}
