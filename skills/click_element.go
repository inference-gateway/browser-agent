package skills

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// ClickElementSkill struct holds the skill with dependencies
type ClickElementSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewClickElementSkill creates a new click_element skill
func NewClickElementSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &ClickElementSkill{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"click_element",
		"Click on an element identified by selector, text, or other locator strategies",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"button": map[string]any{
					"default":     "left",
					"description": "Mouse button to use (left, right, middle)",
					"type":        "string",
				},
				"click_count": map[string]any{
					"default":     1,
					"description": "Number of times to click",
					"type":        "integer",
				},
				"force": map[string]any{
					"default":     false,
					"description": "Force click even if element is not visible",
					"type":        "boolean",
				},
				"selector": map[string]any{
					"description": "CSS selector, XPath, or text to identify the element",
					"type":        "string",
				},
				"timeout": map[string]any{
					"default":     30000,
					"description": "Maximum time to wait for element in milliseconds",
					"type":        "integer",
				},
			},
			"required": []string{"selector"},
		},
		skill.ClickElementHandler,
	)
}

// ClickElementHandler handles the click_element skill execution
func (s *ClickElementSkill) ClickElementHandler(ctx context.Context, args map[string]any) (string, error) {
	selector, ok := args["selector"].(string)
	if !ok || selector == "" {
		return "", fmt.Errorf("selector parameter is required and must be a non-empty string")
	}

	button := "left"
	if b, ok := args["button"].(string); ok && b != "" {
		if !s.isValidButton(b) {
			return "", fmt.Errorf("invalid button value: %s. Must be one of: left, right, middle", b)
		}
		button = b
	}

	clickCount := 1
	if c, ok := args["click_count"].(int); ok && c > 0 {
		clickCount = c
	} else if cf, ok := args["click_count"].(float64); ok && cf > 0 {
		clickCount = int(cf)
	}

	force := false
	if f, ok := args["force"].(bool); ok {
		force = f
	}

	timeout := 30000
	if t, ok := args["timeout"].(int); ok && t > 0 {
		timeout = t
	} else if tf, ok := args["timeout"].(float64); ok && tf > 0 {
		timeout = int(tf)
	}

	s.logger.Info("clicking element",
		zap.String("selector", selector),
		zap.String("button", button),
		zap.Int("click_count", clickCount),
		zap.Bool("force", force),
		zap.Int("timeout_ms", timeout))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	normalizedSelector, selectorType := s.normalizeSelector(selector)
	s.logger.Debug("normalized selector",
		zap.String("original", selector),
		zap.String("normalized", normalizedSelector),
		zap.String("type", selectorType))

	options := map[string]any{
		"timeout":     time.Duration(timeout) * time.Millisecond,
		"force":       force,
		"click_count": clickCount,
		"button":      button,
	}

	err = s.waitForElementActionable(ctx, session.ID, normalizedSelector, timeout)
	if err != nil {
		s.logger.Error("element not actionable",
			zap.String("selector", normalizedSelector),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("element not actionable: %w", err)
	}

	err = s.playwright.ClickElement(ctx, session.ID, normalizedSelector, options)
	if err != nil {
		s.logger.Error("click failed",
			zap.String("selector", normalizedSelector),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("click failed: %w", err)
	}

	s.logger.Info("element clicked successfully",
		zap.String("selector", normalizedSelector),
		zap.String("sessionID", session.ID))

	response := map[string]any{
		"success":       true,
		"selector":      selector,
		"button":        button,
		"click_count":   clickCount,
		"force":         force,
		"timeout_ms":    timeout,
		"session_id":    session.ID,
		"selector_type": selectorType,
		"message":       "Element clicked successfully",
	}

	return fmt.Sprintf(`%+v`, response), nil
}

// isValidButton validates the button parameter
func (s *ClickElementSkill) isValidButton(button string) bool {
	validButtons := []string{"left", "right", "middle"}
	for _, valid := range validButtons {
		if button == valid {
			return true
		}
	}
	return false
}

// normalizeSelector processes the selector to support CSS, XPath, and text-based strategies
func (s *ClickElementSkill) normalizeSelector(selector string) (string, string) {
	selector = strings.TrimSpace(selector)

	if strings.HasPrefix(selector, "/") || strings.HasPrefix(selector, "//") || strings.HasPrefix(selector, "xpath=") {
		if strings.HasPrefix(selector, "xpath=") {
			return selector[6:], "xpath"
		}
		return selector, "xpath"
	}

	textRegex := regexp.MustCompile(`^(text=|:text\(|:has-text\(|:text-is\(|:text-matches\()`)
	if textRegex.MatchString(selector) {
		return selector, "text"
	}

	if strings.Contains(selector, "text=") ||
		(strings.HasPrefix(selector, "'") && strings.HasSuffix(selector, "'")) ||
		(strings.HasPrefix(selector, "\"") && strings.HasSuffix(selector, "\"")) {
		if (strings.HasPrefix(selector, "'") && strings.HasSuffix(selector, "'")) ||
			(strings.HasPrefix(selector, "\"") && strings.HasSuffix(selector, "\"")) {
			text := selector[1 : len(selector)-1]
			return fmt.Sprintf("text=%s", text), "text"
		}
		return selector, "text"
	}

	if strings.HasPrefix(selector, "role=") || strings.Contains(selector, "[role=") {
		return selector, "role"
	}

	if strings.Contains(selector, "data-testid") || strings.Contains(selector, "test-id") {
		return selector, "testid"
	}

	return selector, "css"
}

// waitForElementActionable waits for the element to be actionable before clicking
func (s *ClickElementSkill) waitForElementActionable(ctx context.Context, sessionID, selector string, timeoutMs int) error {
	session, err := s.playwright.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	timeoutDuration := time.Duration(timeoutMs) * time.Millisecond

	err = s.playwright.WaitForCondition(ctx, sessionID, "selector", selector, "visible", timeoutDuration, "")
	if err != nil {
		s.logger.Warn("element not visible, attempting force click if enabled",
			zap.String("selector", selector),
			zap.Error(err))
		return s.checkElementInIframes(ctx, session, selector)
	}

	return nil
}

// checkElementInIframes checks if the element exists within any iframes
func (s *ClickElementSkill) checkElementInIframes(ctx context.Context, session *playwright.BrowserSession, selector string) error {
	if session.Page == nil {
		return fmt.Errorf("element not found: %s (page not available)", selector)
	}

	iframeCount, err := session.Page.Locator("iframe").Count()
	if err != nil {
		return fmt.Errorf("element not found and iframe check failed: %w", err)
	}

	if iframeCount == 0 {
		return fmt.Errorf("element not found: %s", selector)
	}

	s.logger.Info("checking for element in iframes",
		zap.String("selector", selector),
		zap.Int("iframe_count", iframeCount))

	return fmt.Errorf("element not found in main frame, %d iframes detected but cross-frame clicking not yet implemented", iframeCount)
}
