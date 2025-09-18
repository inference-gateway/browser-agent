package skills

import (
	"context"
	"fmt"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/playwright-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// TakeScreenshotSkill struct holds the skill with dependencies
type TakeScreenshotSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewTakeScreenshotSkill creates a new take_screenshot skill
func NewTakeScreenshotSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &TakeScreenshotSkill{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"take_screenshot",
		"Capture a screenshot of the current page or specific element",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"full_page": map[string]any{
					"default":     false,
					"description": "Capture the entire scrollable page",
					"type":        "boolean",
				},
				"path": map[string]any{
					"description": "File path to save the screenshot",
					"type":        "string",
				},
				"quality": map[string]any{
					"default":     80,
					"description": "Quality for jpeg images (0-100)",
					"type":        "integer",
				},
				"selector": map[string]any{
					"description": "Optional selector to screenshot specific element",
					"type":        "string",
				},
				"type": map[string]any{
					"default":     "png",
					"description": "Image format (png, jpeg)",
					"type":        "string",
				},
			},
			"required": []string{"path"},
		},
		skill.TakeScreenshotHandler,
	)
}

// TakeScreenshotHandler handles the take_screenshot skill execution
func (s *TakeScreenshotSkill) TakeScreenshotHandler(ctx context.Context, args map[string]any) (string, error) {
	// TODO: Implement take_screenshot logic
	// Capture a screenshot of the current page or specific element

	// Example of using dependencies:
	// s.logger.SomeMethod(ctx, ...)
	// s.playwright.SomeMethod(ctx, ...)

	// Extract parameters from args
	// full_page := args["full_page"].(bool)
	// path := args["path"].(string)
	// quality := args["quality"].(int)
	// selector := args["selector"].(string)
	// type := args["type"].(string)

	return fmt.Sprintf(`{"result": "TODO: Implement take_screenshot logic", "input": %+v}`, args), nil
}
