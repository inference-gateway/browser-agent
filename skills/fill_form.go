package skills

import (
	"context"
	"fmt"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// FillFormSkill struct holds the skill with dependencies
type FillFormSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewFillFormSkill creates a new fill_form skill
func NewFillFormSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &FillFormSkill{
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
		skill.FillFormHandler,
	)
}

// FillFormHandler handles the fill_form skill execution
func (s *FillFormSkill) FillFormHandler(ctx context.Context, args map[string]any) (string, error) {
	fieldsRaw, ok := args["fields"]
	if !ok {
		return "", fmt.Errorf("fields parameter is required")
	}

	fieldsSlice, ok := fieldsRaw.([]any)
	if !ok {
		return "", fmt.Errorf("fields must be an array")
	}

	if len(fieldsSlice) == 0 {
		return "", fmt.Errorf("at least one field is required")
	}

	fields := make([]map[string]any, 0, len(fieldsSlice))
	for i, fieldRaw := range fieldsSlice {
		fieldMap, ok := fieldRaw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("field %d must be an object", i)
		}

		selector, hasSelector := fieldMap["selector"].(string)
		if !hasSelector || selector == "" {
			return "", fmt.Errorf("field %d: selector is required and must be a non-empty string", i)
		}

		_, hasValue := fieldMap["value"].(string)
		if !hasValue {
			return "", fmt.Errorf("field %d: value is required", i)
		}

		if _, hasType := fieldMap["type"]; !hasType {
			fieldMap["type"] = "text"
		}

		fieldType := fieldMap["type"].(string)
		validTypes := []string{"text", "textarea", "password", "select", "checkbox", "radio", "file"}
		isValidType := false
		for _, vt := range validTypes {
			if fieldType == vt {
				isValidType = true
				break
			}
		}
		if !isValidType {
			return "", fmt.Errorf("field %d: invalid type '%s'. Must be one of: %v", i, fieldType, validTypes)
		}

		fields = append(fields, fieldMap)
	}

	submit := false
	if submitRaw, ok := args["submit"]; ok {
		if submitBool, ok := submitRaw.(bool); ok {
			submit = submitBool
		}
	}

	submitSelector := ""
	if submit {
		if submitSelectorRaw, ok := args["submit_selector"]; ok {
			if ss, ok := submitSelectorRaw.(string); ok && ss != "" {
				submitSelector = ss
			} else {
				return "", fmt.Errorf("submit_selector is required when submit is true")
			}
		} else {
			return "", fmt.Errorf("submit_selector is required when submit is true")
		}
	}

	s.logger.Info("filling form",
		zap.Int("field_count", len(fields)),
		zap.Bool("submit", submit),
		zap.String("submit_selector", submitSelector))

	session, err := s.getOrCreateSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	fieldResults := make([]map[string]any, 0, len(fields))

	for i, field := range fields {
		fieldResult := map[string]any{
			"field_index": i,
			"selector":    field["selector"],
			"type":        field["type"],
			"success":     false,
		}

		selector := field["selector"].(string)
		fieldType := field["type"].(string)

		s.logger.Info("processing field",
			zap.Int("index", i),
			zap.String("selector", selector),
			zap.String("type", fieldType))

		err := s.fillSingleField(ctx, session.ID, field)
		if err != nil {
			s.logger.Error("failed to fill field",
				zap.Int("index", i),
				zap.String("selector", selector),
				zap.Error(err))
			fieldResult["error"] = err.Error()
			fieldResults = append(fieldResults, fieldResult)

			errorResponse := map[string]any{
				"success":      false,
				"session_id":   session.ID,
				"fields_count": len(fields),
				"fields":       fieldResults,
				"error":        fmt.Sprintf("Failed to fill field %d (%s): %v", i, selector, err),
			}

			return fmt.Sprintf(`%+v`, errorResponse), fmt.Errorf("failed to fill field %d (%s): %w", i, selector, err)
		}

		fieldResult["success"] = true
		fieldResult["message"] = "Field filled successfully"
		fieldResults = append(fieldResults, fieldResult)

		s.logger.Info("field filled successfully",
			zap.Int("index", i),
			zap.String("selector", selector))
	}

	var submitResult map[string]any
	if submit {
		s.logger.Info("submitting form", zap.String("submit_selector", submitSelector))

		err = s.playwright.ClickElement(ctx, session.ID, submitSelector, map[string]any{
			"timeout": 30000,
		})

		submitResult = map[string]any{
			"submit_selector": submitSelector,
			"success":         err == nil,
		}

		if err != nil {
			s.logger.Error("failed to submit form", zap.String("submit_selector", submitSelector), zap.Error(err))
			submitResult["error"] = err.Error()
			return "", fmt.Errorf("failed to submit form: %w", err)
		} else {
			s.logger.Info("form submitted successfully", zap.String("submit_selector", submitSelector))
			submitResult["message"] = "Form submitted successfully"
		}
	}

	response := map[string]any{
		"success":      true,
		"session_id":   session.ID,
		"fields_count": len(fields),
		"fields":       fieldResults,
		"message":      fmt.Sprintf("Successfully filled %d fields", len(fields)),
	}

	if submit {
		response["submit"] = submitResult
		if submitResult["success"].(bool) {
			response["message"] = fmt.Sprintf("Successfully filled %d fields and submitted form", len(fields))
		}
	}

	return fmt.Sprintf(`%+v`, response), nil
}

// fillSingleField handles filling a single form field with enhanced type support
func (s *FillFormSkill) fillSingleField(ctx context.Context, sessionID string, field map[string]any) error {
	fields := []map[string]any{field}

	selector := field["selector"].(string)
	value := field["value"].(string)
	fieldType := field["type"].(string)

	s.logger.Debug("delegating field filling to playwright service",
		zap.String("selector", selector),
		zap.String("type", fieldType),
		zap.String("value", value))

	return s.playwright.FillForm(ctx, sessionID, fields, false, "")
}

// getOrCreateSession gets the shared default session
func (s *FillFormSkill) getOrCreateSession(ctx context.Context) (*playwright.BrowserSession, error) {
	return s.playwright.GetOrCreateDefaultSession(ctx)
}
