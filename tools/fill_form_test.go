package tools

import (
	"context"
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/assert"
	zap "go.uber.org/zap"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

func TestFillFormTool_FillFormHandler_ValidationTests(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	tool := &FillFormTool{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name          string
		args          map[string]any
		expectedError bool
		errorContains string
	}{
		{
			name:          "missing fields parameter",
			args:          map[string]any{},
			expectedError: true,
			errorContains: "fields parameter is required",
		},
		{
			name: "invalid fields type",
			args: map[string]any{
				"fields": "invalid",
			},
			expectedError: true,
			errorContains: "fields must be an array",
		},
		{
			name: "empty fields array",
			args: map[string]any{
				"fields": []any{},
			},
			expectedError: true,
			errorContains: "at least one field is required",
		},
		{
			name: "missing selector in field",
			args: map[string]any{
				"fields": []any{
					map[string]any{
						"value": "test",
					},
				},
			},
			expectedError: true,
			errorContains: "selector is required",
		},
		{
			name: "missing value in field",
			args: map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#test",
					},
				},
			},
			expectedError: true,
			errorContains: "value is required",
		},
		{
			name: "invalid field type",
			args: map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#test",
						"value":    "test",
						"type":     "invalid_type",
					},
				},
			},
			expectedError: true,
			errorContains: "invalid type 'invalid_type'",
		},
		{
			name: "submit without submit_selector",
			args: map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#test",
						"value":    "test",
					},
				},
				"submit": true,
			},
			expectedError: true,
			errorContains: "submit_selector is required when submit is true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.FillFormHandler(context.Background(), tt.args)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestFillFormTool_FillFormHandler_SuccessTests(t *testing.T) {
	tests := []struct {
		name string
		args map[string]any
	}{
		{
			name: "text input",
			args: map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#username",
						"value":    "testuser",
						"type":     "text",
					},
				},
			},
		},
		{
			name: "multiple field types",
			args: map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#username",
						"value":    "testuser",
						"type":     "text",
					},
					map[string]any{
						"selector": "#agree",
						"value":    "true",
						"type":     "checkbox",
					},
					map[string]any{
						"selector": "#country",
						"value":    "US",
						"type":     "select",
					},
				},
			},
		},
		{
			name: "with form submission",
			args: map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#email",
						"value":    "test@example.com",
						"type":     "text",
					},
				},
				"submit":          true,
				"submit_selector": "#submit-btn",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()
			mockPlaywright := &mocks.FakeBrowserAutomation{}
			session := &playwright.BrowserSession{ID: "test-session"}
			mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
			mockPlaywright.GetSessionReturns(session, nil)
			mockPlaywright.FillFormReturns(nil)

			tool := &FillFormTool{logger: logger, playwright: mockPlaywright}
			result, err := tool.FillFormHandler(context.Background(), tt.args)

			assert.NoError(t, err)
			assert.NotEmpty(t, result)

			var parsed map[string]any
			assert.NoError(t, json.Unmarshal([]byte(result), &parsed), "response should be valid JSON")
			assert.Equal(t, true, parsed["success"])
		})
	}
}

// TestFillFormTool_FillFormHandler_SingleBatchCall confirms the refactor:
// previously the tool re-called FillForm once per field and additionally
// called ClickElement to submit. After the fix, the playwright service
// receives one batched call regardless of field count.
func TestFillFormTool_FillFormHandler_SingleBatchCall(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.FillFormReturns(nil)

	tool := &FillFormTool{logger: logger, playwright: mockPlaywright}

	_, err := tool.FillFormHandler(context.Background(), map[string]any{
		"fields": []any{
			map[string]any{"selector": "#a", "value": "1"},
			map[string]any{"selector": "#b", "value": "2"},
			map[string]any{"selector": "#c", "value": "3"},
		},
		"submit":          true,
		"submit_selector": "#submit",
	})
	assert.NoError(t, err)

	assert.Equal(t, 1, mockPlaywright.FillFormCallCount(),
		"FillForm should be called exactly once with the full batch")
	assert.Equal(t, 0, mockPlaywright.ClickElementCallCount(),
		"submit should be delegated to playwright.FillForm, not a separate ClickElement")

	_, _, gotFields, gotSubmit, gotSelector := mockPlaywright.FillFormArgsForCall(0)
	assert.Len(t, gotFields, 3)
	assert.True(t, gotSubmit)
	assert.Equal(t, "#submit", gotSelector)
}

func TestFillFormTool_ValidateFieldTypes(t *testing.T) {
	validTypes := []string{"text", "textarea", "password", "select", "checkbox", "radio", "file"}

	for _, fieldType := range validTypes {
		t.Run("valid_type_"+fieldType, func(t *testing.T) {
			logger := zap.NewNop()
			mockPlaywright := &mocks.FakeBrowserAutomation{}

			session := &playwright.BrowserSession{ID: "test-session"}
			mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
			mockPlaywright.GetSessionReturns(session, nil)
			mockPlaywright.FillFormReturns(nil)

			tool := &FillFormTool{logger: logger, playwright: mockPlaywright}
			_, err := tool.FillFormHandler(context.Background(), map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#test",
						"value":    "test",
						"type":     fieldType,
					},
				},
			})

			if err != nil {
				assert.NotContains(t, err.Error(), "invalid type")
			}
		})
	}
}

func TestFillFormTool_DefaultFieldType(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.FillFormReturns(nil)

	tool := &FillFormTool{logger: logger, playwright: mockPlaywright}
	_, err := tool.FillFormHandler(context.Background(), map[string]any{
		"fields": []any{
			map[string]any{
				"selector": "#test",
				"value":    "test",
			},
		},
	})

	assert.NoError(t, err)

	_, _, gotFields, _, _ := mockPlaywright.FillFormArgsForCall(0)
	assert.Equal(t, "text", gotFields[0]["type"], "type should default to text")
}
