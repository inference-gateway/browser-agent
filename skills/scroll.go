package skills

import (
	"context"
	"encoding/json"
	"fmt"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// ScrollSkill struct holds the skill with dependencies
type ScrollSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewScrollSkill creates a new scroll skill
func NewScrollSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &ScrollSkill{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"scroll",
		"Scroll the page or element to a specific position or into view",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{
					"type":        "string",
					"description": "What to scroll: 'page', 'element', or 'coordinates'",
					"enum":        []string{"page", "element", "coordinates"},
				},
				"selector": map[string]any{
					"type":        "string",
					"description": "Element selector (required if target=element)",
				},
				"behavior": map[string]any{
					"type":        "string",
					"description": "Scroll behavior: 'smooth' or 'instant'",
					"enum":        []string{"smooth", "instant"},
					"default":     "smooth",
				},
				"block": map[string]any{
					"type":        "string",
					"description": "Vertical alignment: 'start', 'center', 'end', 'nearest'",
					"enum":        []string{"start", "center", "end", "nearest"},
					"default":     "start",
				},
				"inline": map[string]any{
					"type":        "string",
					"description": "Horizontal alignment: 'start', 'center', 'end', 'nearest'",
					"enum":        []string{"start", "center", "end", "nearest"},
					"default":     "nearest",
				},
				"x": map[string]any{
					"type":        "integer",
					"description": "X coordinate for scrolling (if target=coordinates)",
				},
				"y": map[string]any{
					"type":        "integer",
					"description": "Y coordinate for scrolling (if target=coordinates)",
				},
				"direction": map[string]any{
					"type":        "string",
					"description": "Direction to scroll: 'up', 'down', 'left', 'right', 'top', 'bottom'",
					"enum":        []string{"up", "down", "left", "right", "top", "bottom"},
				},
				"amount": map[string]any{
					"type":        "integer",
					"description": "Amount to scroll in pixels (for directional scrolling)",
				},
			},
			"required": []string{"target"},
		},
		skill.ScrollHandler,
	)
}

// ScrollHandler handles the scroll skill execution
func (s *ScrollSkill) ScrollHandler(ctx context.Context, args map[string]any) (string, error) {
	target, ok := args["target"].(string)
	if !ok || target == "" {
		return "", fmt.Errorf("target parameter is required and must be a non-empty string")
	}

	// Validate target value
	if !s.isValidTarget(target) {
		return "", fmt.Errorf("invalid target value: %s. Must be one of: page, element, coordinates", target)
	}

	// Extract optional parameters with defaults
	behavior := "smooth"
	if b, ok := args["behavior"].(string); ok && b != "" {
		if !s.isValidBehavior(b) {
			return "", fmt.Errorf("invalid behavior value: %s. Must be one of: smooth, instant", b)
		}
		behavior = b
	}

	block := "start"
	if bl, ok := args["block"].(string); ok && bl != "" {
		if !s.isValidAlignment(bl) {
			return "", fmt.Errorf("invalid block value: %s. Must be one of: start, center, end, nearest", bl)
		}
		block = bl
	}

	inline := "nearest"
	if in, ok := args["inline"].(string); ok && in != "" {
		if !s.isValidAlignment(in) {
			return "", fmt.Errorf("invalid inline value: %s. Must be one of: start, center, end, nearest", in)
		}
		inline = in
	}

	selector := ""
	if sel, ok := args["selector"].(string); ok {
		selector = sel
	}

	direction := ""
	if dir, ok := args["direction"].(string); ok {
		direction = dir
	}

	amount := 0
	if amt, ok := args["amount"].(int); ok {
		amount = amt
	} else if amtf, ok := args["amount"].(float64); ok {
		amount = int(amtf)
	}

	x := 0
	if xVal, ok := args["x"].(int); ok {
		x = xVal
	} else if xf, ok := args["x"].(float64); ok {
		x = int(xf)
	}

	y := 0
	if yVal, ok := args["y"].(int); ok {
		y = yVal
	} else if yf, ok := args["y"].(float64); ok {
		y = int(yf)
	}

	// Validate based on target
	switch target {
	case "element":
		if selector == "" {
			return "", fmt.Errorf("selector is required when target is 'element'")
		}
	case "coordinates":
		// x and y are optional, but at least one should be provided for coordinates
		if x == 0 && y == 0 {
			s.logger.Warn("both x and y are 0 for coordinates scrolling")
		}
	case "page":
		if direction != "" && !s.isValidDirection(direction) {
			return "", fmt.Errorf("invalid direction value: %s. Must be one of: up, down, left, right, top, bottom", direction)
		}
	}

	s.logger.Info("executing scroll",
		zap.String("target", target),
		zap.String("selector", selector),
		zap.String("behavior", behavior),
		zap.String("block", block),
		zap.String("inline", inline),
		zap.String("direction", direction),
		zap.Int("amount", amount),
		zap.Int("x", x),
		zap.Int("y", y))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	// Call the new scroll method on playwright service
	err = s.playwright.Scroll(ctx, session.ID, target, selector, behavior, block, inline, direction, amount, x, y)
	if err != nil {
		s.logger.Error("scroll failed",
			zap.String("target", target),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("scroll failed: %w", err)
	}

	s.logger.Info("scroll completed successfully",
		zap.String("target", target),
		zap.String("sessionID", session.ID))

	response := map[string]any{
		"success":    true,
		"target":     target,
		"selector":   selector,
		"behavior":   behavior,
		"block":      block,
		"inline":     inline,
		"direction":  direction,
		"amount":     amount,
		"x":          x,
		"y":          y,
		"session_id": session.ID,
		"message":    "Scroll completed successfully",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(responseJSON), nil
}

// isValidTarget validates the target parameter
func (s *ScrollSkill) isValidTarget(target string) bool {
	validTargets := []string{"page", "element", "coordinates"}
	for _, valid := range validTargets {
		if target == valid {
			return true
		}
	}
	return false
}

// isValidBehavior validates the behavior parameter
func (s *ScrollSkill) isValidBehavior(behavior string) bool {
	validBehaviors := []string{"smooth", "instant"}
	for _, valid := range validBehaviors {
		if behavior == valid {
			return true
		}
	}
	return false
}

// isValidAlignment validates block and inline parameters
func (s *ScrollSkill) isValidAlignment(alignment string) bool {
	validAlignments := []string{"start", "center", "end", "nearest"}
	for _, valid := range validAlignments {
		if alignment == valid {
			return true
		}
	}
	return false
}

// isValidDirection validates the direction parameter
func (s *ScrollSkill) isValidDirection(direction string) bool {
	validDirections := []string{"up", "down", "left", "right", "top", "bottom"}
	for _, valid := range validDirections {
		if direction == valid {
			return true
		}
	}
	return false
}
