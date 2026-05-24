package tools

import (
	"context"
	"fmt"

	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

var validFieldTypes = []string{"text", "textarea", "password", "select", "checkbox", "radio", "file"}

// FillFormTool struct holds the tool with dependencies
type FillFormTool struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewFillFormTool creates a new fill_form tool
func NewFillFormTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	tool := &FillFormTool{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"fill_form",
		"Fill form fields with provided data, handling various input types",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"fields": map[string]any{
					"description": "List of form fields to fill",
					"items": map[string]any{
						"required": []string{"selector", "value"},
						"type":     "object",
						"properties": map[string]any{
							"selector": map[string]any{
								"type":        "string",
								"description": "Selector for the form field",
							},
							"value": map[string]any{
								"type":        "string",
								"description": "Value to fill in the field. For select with multiple=true, use comma-separated values",
							},
							"type": map[string]any{
								"type":        "string",
								"description": "Type of input: text, textarea, password, select, checkbox, radio, file",
								"default":     "text",
							},
							"multiple": map[string]any{
								"type":        "boolean",
								"description": "For select fields only: whether this is a multi-select dropdown",
								"default":     false,
							},
						},
					},
					"type": "array",
				},
				"submit": map[string]any{
					"default":     false,
					"description": "Whether to submit the form after filling",
					"type":        "boolean",
				},
				"submit_selector": map[string]any{
					"description": "Selector for the submit button if submit is true",
					"type":        "string",
				},
			},
			"required": []string{"fields"},
		},
		tool.FillFormHandler,
	)
}

// FillFormHandler handles the fill_form tool execution.
//
// Unlike the previous implementation, the underlying playwright service is
// called ONCE with the full field batch (and the submit flag), rather than
// per-field. The per-field loop was N round-trips and a separate submit
// click with a mismatched timeout type.
func (s *FillFormTool) FillFormHandler(ctx context.Context, args map[string]any) (string, error) {
	rawFields, present, err := sliceArg(args, "fields")
	if err != nil {
		return "", err
	}
	if !present {
		return "", fmt.Errorf("fields parameter is required")
	}
	if len(rawFields) == 0 {
		return "", fmt.Errorf("at least one field is required")
	}

	fields, err := s.validateAndNormalizeFields(rawFields)
	if err != nil {
		return "", err
	}

	submit, err := boolArg(args, "submit", false)
	if err != nil {
		return "", err
	}

	submitSelector := ""
	if submit {
		submitSelector, err = requiredString(args, "submit_selector")
		if err != nil {
			return "", fmt.Errorf("submit_selector is required when submit is true")
		}
	}

	s.logger.Info("filling form",
		zap.Int("field_count", len(fields)),
		zap.Bool("submit", submit),
		zap.String("submit_selector", submitSelector))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	if err := s.playwright.FillForm(ctx, session.ID, fields, submit, submitSelector); err != nil {
		s.logger.Error("fill_form failed",
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("fill_form failed: %w", err)
	}

	message := fmt.Sprintf("Successfully filled %d fields", len(fields))
	if submit {
		message = fmt.Sprintf("Successfully filled %d fields and submitted form", len(fields))
	}

	s.logger.Info("fill_form completed", zap.String("sessionID", session.ID))

	return marshalResponse(map[string]any{
		"success":         true,
		"session_id":      session.ID,
		"fields_count":    len(fields),
		"submit":          submit,
		"submit_selector": submitSelector,
		"message":         message,
	})
}

// validateAndNormalizeFields converts the raw []any into validated
// []map[string]any, defaulting field type to "text" if absent.
func (s *FillFormTool) validateAndNormalizeFields(rawFields []any) ([]map[string]any, error) {
	fields := make([]map[string]any, 0, len(rawFields))
	for i, raw := range rawFields {
		field, ok := raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("field %d must be an object", i)
		}

		selector, hasSelector := field["selector"].(string)
		if !hasSelector || selector == "" {
			return nil, fmt.Errorf("field %d: selector is required and must be a non-empty string", i)
		}

		if _, hasValue := field["value"].(string); !hasValue {
			return nil, fmt.Errorf("field %d: value is required and must be a string", i)
		}

		fieldType, _ := field["type"].(string)
		if fieldType == "" {
			fieldType = "text"
			field["type"] = fieldType
		}
		if !oneOf(fieldType, validFieldTypes...) {
			return nil, fmt.Errorf("field %d: invalid type '%s'. Must be one of: %v", i, fieldType, validFieldTypes)
		}

		fields = append(fields, field)
	}
	return fields, nil
}
