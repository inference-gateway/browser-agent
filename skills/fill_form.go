package skills

import (
	"context"
	"fmt"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/playwright-agent/internal/playwright"
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
					"items":       map[string]any{"type": "object", "properties": map[string]any{"value": map[string]any{"type": "string", "description": "Value to fill in the field"}, "type": map[string]any{"type": "string", "description": "Type of input (text, select, checkbox, radio)"}, "selector": map[string]any{"type": "string", "description": "Selector for the form field"}}, "required": []string{"selector", "value"}},
					"type":        "array",
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
	// TODO: Implement fill_form logic
	// Fill form fields with provided data, handling various input types

	// Example of using dependencies:
	// s.logger.SomeMethod(ctx, ...)
	// s.playwright.SomeMethod(ctx, ...)

	// Extract parameters from args
	// fields := args["fields"]
	// submit := args["submit"].(bool)
	// submit_selector := args["submit_selector"].(string)

	return fmt.Sprintf(`{"result": "TODO: Implement fill_form logic", "input": %+v}`, args), nil
}
