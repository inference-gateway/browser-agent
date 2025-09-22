package skills

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"
)

func TestFillFormSkill_FillFormHandler_ValidationTests(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	skill := &FillFormSkill{
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
			result, err := skill.FillFormHandler(context.Background(), tt.args)

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

func TestFillFormSkill_FillFormHandler_SuccessTests(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	skill := &FillFormSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateDefaultSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.FillFormReturns(nil)
	mockPlaywright.FillFormReturns(nil)
	mockPlaywright.ClickElementReturns(nil)

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
			result, err := skill.FillFormHandler(context.Background(), tt.args)

			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "success")
		})
	}
}

func TestFillFormSkill_ValidateFieldTypes(t *testing.T) {
	validTypes := []string{"text", "textarea", "password", "select", "checkbox", "radio", "file"}

	for _, fieldType := range validTypes {
		t.Run("valid_type_"+fieldType, func(t *testing.T) {
			logger := zap.NewNop()
			mockPlaywright := &mocks.FakeBrowserAutomation{}

			skill := &FillFormSkill{
				logger:     logger,
				playwright: mockPlaywright,
			}

			args := map[string]any{
				"fields": []any{
					map[string]any{
						"selector": "#test",
						"value":    "test",
						"type":     fieldType,
					},
				},
			}

			session := &playwright.BrowserSession{ID: "test-session"}
			mockPlaywright.GetOrCreateDefaultSessionReturns(session, nil)
			mockPlaywright.GetSessionReturns(session, nil)
			mockPlaywright.FillFormReturns(nil)

			_, err := skill.FillFormHandler(context.Background(), args)

			if err != nil {
				assert.NotContains(t, err.Error(), "invalid type")
			}
		})
	}
}

func TestFillFormSkill_DefaultFieldType(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	skill := &FillFormSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	args := map[string]any{
		"fields": []any{
			map[string]any{
				"selector": "#test",
				"value":    "test",
			},
		},
	}

	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateDefaultSessionReturns(session, nil)
	mockPlaywright.GetSessionReturns(session, nil)
	mockPlaywright.FillFormReturns(nil)

	_, err := skill.FillFormHandler(context.Background(), args)

	if err != nil {
		assert.NotContains(t, err.Error(), "invalid type")
		assert.NotContains(t, err.Error(), "type is required")
	}
}
