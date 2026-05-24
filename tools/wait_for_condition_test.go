package tools

import (
	"context"
	"encoding/json"
	"testing"

	zap "go.uber.org/zap"

	assert "github.com/stretchr/testify/assert"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

func TestWaitForConditionTool_NewWaitForConditionTool(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	tool := NewWaitForConditionTool(logger, mockPlaywright)

	assert.NotNil(t, tool)
}

func TestWaitForConditionTool_ValidateConditionRequirements(t *testing.T) {
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
			err := validateConditionRequirements(tt.condition, tt.selector, tt.customFunction)
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

	session := &playwright.BrowserSession{ID: "test-session"}
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
	var parsed map[string]any
	assert.NoError(t, json.Unmarshal([]byte(result), &parsed), "response should be valid JSON")
	assert.Equal(t, true, parsed["success"])
	assert.Equal(t, "selector", parsed["condition"])
	assert.Equal(t, ".button", parsed["selector"])
}

func TestWaitForConditionTool_WaitForConditionHandler_MissingCondition(t *testing.T) {
	tool := &WaitForConditionTool{
		logger:     zap.NewNop(),
		playwright: &mocks.FakeBrowserAutomation{},
	}
	result, err := tool.WaitForConditionHandler(context.Background(), map[string]any{
		"selector": ".button",
		"state":    "visible",
	})

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "condition parameter is required")
}

func TestWaitForConditionTool_WaitForConditionHandler_InvalidCondition(t *testing.T) {
	tool := &WaitForConditionTool{
		logger:     zap.NewNop(),
		playwright: &mocks.FakeBrowserAutomation{},
	}
	result, err := tool.WaitForConditionHandler(context.Background(), map[string]any{
		"condition": "invalid",
		"selector":  ".button",
	})

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "invalid condition type")
}

func TestWaitForConditionTool_WaitForConditionHandler_InvalidState(t *testing.T) {
	tool := &WaitForConditionTool{
		logger:     zap.NewNop(),
		playwright: &mocks.FakeBrowserAutomation{},
	}
	result, err := tool.WaitForConditionHandler(context.Background(), map[string]any{
		"condition": "selector",
		"selector":  ".button",
		"state":     "invalid",
	})

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "invalid state")
}

func TestWaitForConditionTool_WaitForConditionHandler_SelectorWithoutSelector(t *testing.T) {
	tool := &WaitForConditionTool{
		logger:     zap.NewNop(),
		playwright: &mocks.FakeBrowserAutomation{},
	}
	result, err := tool.WaitForConditionHandler(context.Background(), map[string]any{
		"condition": "selector",
		"state":     "visible",
	})

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "selector parameter is required")
}

func TestWaitForConditionTool_WaitForConditionHandler_FunctionWithoutFunction(t *testing.T) {
	tool := &WaitForConditionTool{
		logger:     zap.NewNop(),
		playwright: &mocks.FakeBrowserAutomation{},
	}
	result, err := tool.WaitForConditionHandler(context.Background(), map[string]any{
		"condition": "function",
	})

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

	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.WaitForConditionReturns(nil)

	result, err := tool.WaitForConditionHandler(context.Background(), map[string]any{
		"condition": "selector",
		"selector":  ".button",
	})

	assert.NoError(t, err)
	var parsed map[string]any
	assert.NoError(t, json.Unmarshal([]byte(result), &parsed))
	assert.Equal(t, true, parsed["success"])
	assert.Equal(t, "visible", parsed["state"])
	assert.Equal(t, float64(defaultTimeoutMs), parsed["timeout_ms"])
}

// TestWaitForConditionTool_NetworkidleDoesNotInjectJS confirms that the
// networkidle branch is delegated to the playwright service as-is, with
// an empty customFunction. Previously the tool emitted a hand-rolled JS
// blob that monkey-patched window.fetch / window.XMLHttpRequest and never
// restored them.
func TestWaitForConditionTool_NetworkidleDoesNotInjectJS(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.WaitForConditionReturns(nil)

	tool := &WaitForConditionTool{logger: logger, playwright: mockPlaywright}
	_, err := tool.WaitForConditionHandler(context.Background(), map[string]any{
		"condition": "networkidle",
		"timeout":   1000,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, mockPlaywright.WaitForConditionCallCount())

	_, _, gotCondition, _, _, _, gotCustomFn := mockPlaywright.WaitForConditionArgsForCall(0)
	assert.Equal(t, "networkidle", gotCondition,
		"networkidle should pass through as the condition name, not be rewritten to 'function'")
	assert.Empty(t, gotCustomFn, "no custom JS should be injected for networkidle")
}
