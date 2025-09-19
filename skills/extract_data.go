package skills

import (
	"context"
	"fmt"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/playwright-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// ExtractDataSkill struct holds the skill with dependencies
type ExtractDataSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewExtractDataSkill creates a new extract_data skill
func NewExtractDataSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &ExtractDataSkill{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"extract_data",
		"Extract data from the page using selectors and return structured information",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"extractors": map[string]any{
					"description": "List of data extractors to run",
					"items":       map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string", "description": "Name for the extracted data field"}, "selector": map[string]any{"type": "string", "description": "CSS selector or XPath to extract data from"}, "attribute": map[string]any{"type": "string", "description": "Attribute to extract (text, href, src, etc.)", "default": "text"}, "multiple": map[string]any{"type": "boolean", "description": "Extract all matching elements or just the first", "default": false}}, "required": []string{"name", "selector"}},
					"type":        "array",
				},
				"format": map[string]any{
					"default":     "json",
					"description": "Output format (json, csv, text)",
					"type":        "string",
				},
			},
			"required": []string{"extractors"},
		},
		skill.ExtractDataHandler,
	)
}

// ExtractDataHandler handles the extract_data skill execution
func (s *ExtractDataSkill) ExtractDataHandler(ctx context.Context, args map[string]any) (string, error) {
	// TODO: Implement extract_data logic
	// Extract data from the page using selectors and return structured information

	// Example of using dependencies:
	// s.logger.SomeMethod(ctx, ...)
	// s.playwright.SomeMethod(ctx, ...)

	// Extract parameters from args
	// extractors := args["extractors"]
	// format := args["format"].(string)

	return fmt.Sprintf(`{"result": "TODO: Implement extract_data logic", "input": %+v}`, args), nil
}
