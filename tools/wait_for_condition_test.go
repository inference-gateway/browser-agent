package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/inference-gateway/browser-agent/internal/playwright"
	"github.com/inference-gateway/browser-agent/internal/playwright/mocks"
)

func TestWaitForConditionTool_NewWaitForConditionTool(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	tool := NewWaitForConditionTool(logger, mockPlaywright)

	assert.NotNil(t, tool)
}

func TestWaitForConditionTool_ValidateCondition(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name      string
		condition string
		expected  bool
	}{
		{"valid selector", "selector", true},
		{"valid navigation", "navigation", true},
		{"valid function", "function", true},
		{"valid timeout", "timeout", true},
		{"valid networkidle", "networkidle", true},
		{"invalid condition", "invalid", false},
		{"empty condition", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.isValidCondition(tt.condition)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWaitForConditionTool_ValidateState(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name     string
		state    string
		expected bool
	}{
		{"valid visible", "visible", true},
		{"valid hidden", "hidden", true},
		{"valid attached", "attached", true},
		{"valid detached", "detached", true},
		{"invalid state", "invalid", false},
		{"empty state", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.isValidState(tt.state)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWaitForConditionTool_ValidateConditionRequirements(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name           string
		condition      string
		selector       string
		customFunction string
		shouldError    bool
	}{
		{"selector with selector", "selector", ".button", "", false},
		{"selector without selector", "selector", "", "", true},
		{"function with function", "function", "", "() => true", false},
		{"function without function", "function", "", "", true},
		{"navigation no requirements", "navigation", "", "", false},
		{"timeout no requirements", "timeout", "", "", false},
		{"networkidle no requirements", "networkidle", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.validateConditionRequirements(tt.condition, tt.selector, tt.customFunction)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWaitForConditionTool_WaitForConditionHandler_Success(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	session := &playwright.BrowserSession{
		ID: "test-session",
	}

	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.WaitForConditionReturns(nil)

	args := map[string]any{
		"condition": "selector",
		"selector":  ".button",
		"state":     "visible",
		"timeout":   5000,
	}

	result, err := tool.WaitForConditionHandler(context.Background(), args)

	assert.NoError(t, err)
	assert.Contains(t, result, "success:true")
	assert.Contains(t, result, "condition:selector")
	assert.Contains(t, result, "selector:.button")
}

func TestWaitForConditionTool_WaitForConditionHandler_MissingCondition(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	args := map[string]any{
		"selector": ".button",
		"state":    "visible",
	}

	result, err := tool.WaitForConditionHandler(context.Background(), args)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "condition parameter is required")
}

func TestWaitForConditionTool_WaitForConditionHandler_InvalidCondition(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	args := map[string]any{
		"condition": "invalid",
		"selector":  ".button",
	}

	result, err := tool.WaitForConditionHandler(context.Background(), args)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "invalid condition type")
}

func TestWaitForConditionTool_WaitForConditionHandler_InvalidState(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	args := map[string]any{
		"condition": "selector",
		"selector":  ".button",
		"state":     "invalid",
	}

	result, err := tool.WaitForConditionHandler(context.Background(), args)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "invalid state")
}

func TestWaitForConditionTool_WaitForConditionHandler_SelectorWithoutSelector(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	args := map[string]any{
		"condition": "selector",
		"state":     "visible",
	}

	result, err := tool.WaitForConditionHandler(context.Background(), args)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "selector parameter is required")
}

func TestWaitForConditionTool_WaitForConditionHandler_FunctionWithoutFunction(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	args := map[string]any{
		"condition": "function",
	}

	result, err := tool.WaitForConditionHandler(context.Background(), args)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "custom_function parameter is required")
}

func TestWaitForConditionTool_WaitForConditionHandler_DefaultValues(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	session := &playwright.BrowserSession{
		ID: "test-session",
	}

	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.WaitForConditionReturns(nil)

	args := map[string]any{
		"condition": "selector",
		"selector":  ".button",
	}

	result, err := tool.WaitForConditionHandler(context.Background(), args)

	assert.NoError(t, err)
	assert.Contains(t, result, "success:true")
	assert.Contains(t, result, "state:visible")
	assert.Contains(t, result, "timeout_ms:30000")
}
