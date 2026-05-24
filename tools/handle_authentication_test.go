package tools

import (
	"context"
	"testing"

	assert "github.com/stretchr/testify/assert"

	zap "go.uber.org/zap"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"
)

func TestHandleAuthenticationTool_HandleAuthenticationHandler(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	tool := &HandleAuthenticationTool{logger: logger, playwright: mockPlaywright}

	tests := []struct {
		name          string
		args          map[string]any
		errorContains string
	}{
		{
			name:          "missing type",
			args:          map[string]any{},
			errorContains: "type parameter is required",
		},
		{
			name:          "invalid type",
			args:          map[string]any{"type": "saml"},
			errorContains: "invalid auth type",
		},
		{
			name:          "valid type still returns not-implemented",
			args:          map[string]any{"type": "form", "username": "u", "password": "p"},
			errorContains: "not yet implemented",
		},
		{
			name:          "basic returns not-implemented",
			args:          map[string]any{"type": "basic"},
			errorContains: "not yet implemented",
		},
		{
			name:          "oauth returns not-implemented",
			args:          map[string]any{"type": "oauth"},
			errorContains: "not yet implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.HandleAuthenticationHandler(context.Background(), tt.args)
			assert.Error(t, err)
			assert.Empty(t, result)
			assert.Contains(t, err.Error(), tt.errorContains)
		})
	}

	// Defense-in-depth: even when called, the tool MUST NOT touch the
	// playwright service. Previously the tool was a no-op that bypassed
	// the service entirely; the new explicit error preserves that
	// no-side-effects property.
	assert.Equal(t, 0, mockPlaywright.HandleAuthenticationCallCount(),
		"handle_authentication must not call the playwright service while unimplemented")
}
