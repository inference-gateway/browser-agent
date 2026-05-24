package tools

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

// textSelectorPrefixes recognizes Playwright text-engine selectors.
// Compiled once at package init instead of per call.
var textSelectorPrefixes = regexp.MustCompile(`^(text=|:text\(|:has-text\(|:text-is\(|:text-matches\()`)

var validButtons = []string{"left", "right", "middle"}

func isQuotedString(s string) bool {
	return (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) ||
		(strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\""))
}

// ClickElementTool struct holds the tool with dependencies
type ClickElementTool struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewClickElementTool creates a new click_element tool
func NewClickElementTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	tool := &ClickElementTool{
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
					"default":     defaultClickCount,
					"description": "Number of times to click",
					"type":        "integer",
				},
				"force": map[string]any{
					"default":     false,
					"description": "Force click even if element is not visible or actionable (skips the pre-click visibility wait)",
					"type":        "boolean",
				},
				"selector": map[string]any{
					"description": "CSS selector, XPath, or text to identify the element",
					"type":        "string",
				},
				"timeout": map[string]any{
					"default":     defaultTimeoutMs,
					"description": "Maximum time to wait for element in milliseconds",
					"type":        "integer",
				},
			},
			"required": []string{"selector"},
		},
		tool.ClickElementHandler,
	)
}

// ClickElementHandler handles the click_element tool execution
func (s *ClickElementTool) ClickElementHandler(ctx context.Context, args map[string]any) (string, error) {
	selector, err := requiredString(args, "selector")
	if err != nil {
		return "", err
	}

	button, err := stringArg(args, "button", "left")
	if err != nil {
		return "", err
	}
	if !oneOf(button, validButtons...) {
		return "", fmt.Errorf("invalid button value: %s. Must be one of: %v", button, validButtons)
	}

	clickCount, err := boundedIntArg(args, "click_count", defaultClickCount, minClickCount, maxClickCount)
	if err != nil {
		return "", err
	}

	force, err := boolArg(args, "force", false)
	if err != nil {
		return "", err
	}

	timeout, err := boundedIntArg(args, "timeout", defaultTimeoutMs, minTimeoutMs, maxTimeoutMs)
	if err != nil {
		return "", err
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

	timeoutDuration := time.Duration(timeout) * time.Millisecond
	options := map[string]any{
		"timeout":     timeoutDuration,
		"force":       force,
		"click_count": clickCount,
		"button":      button,
	}

	// When force=true, skip the actionability wait — the caller explicitly
	// asked to click without waiting for visibility (e.g. covered elements,
	// pointer-events:none). The previous behaviour ignored the flag.
	if !force {
		if err := s.waitForElementActionable(ctx, session, normalizedSelector, timeout); err != nil {
			s.logger.Error("element not actionable",
				zap.String("selector", normalizedSelector),
				zap.String("sessionID", session.ID),
				zap.Error(err))
			return "", fmt.Errorf("element not actionable: %w", err)
		}
	}

	if err := s.playwright.ClickElement(ctx, session.ID, normalizedSelector, options); err != nil {
		s.logger.Error("click failed",
			zap.String("selector", normalizedSelector),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("click failed: %w", err)
	}

	s.logger.Info("element clicked successfully",
		zap.String("selector", normalizedSelector),
		zap.String("sessionID", session.ID))

	return marshalResponse(map[string]any{
		"success":       true,
		"selector":      selector,
		"button":        button,
		"click_count":   clickCount,
		"force":         force,
		"timeout_ms":    timeout,
		"session_id":    session.ID,
		"selector_type": selectorType,
		"message":       "Element clicked successfully",
	})
}

// normalizeSelector processes the selector to support CSS, XPath, and text-based strategies.
// The text= classification matches Playwright's text engine PREFIXES exactly, not
// substring "text=" which would false-positive on CSS like [data-text="x"].
func (s *ClickElementTool) normalizeSelector(selector string) (string, string) {
	selector = strings.TrimSpace(selector)

	if strings.HasPrefix(selector, "xpath=") {
		return selector[6:], "xpath"
	}
	if strings.HasPrefix(selector, "/") || strings.HasPrefix(selector, "//") {
		return selector, "xpath"
	}

	if textSelectorPrefixes.MatchString(selector) {
		return selector, "text"
	}

	if isQuotedString(selector) {
		return fmt.Sprintf("text=%s", selector[1:len(selector)-1]), "text"
	}

	if strings.HasPrefix(selector, "role=") || strings.Contains(selector, "[role=") {
		return selector, "role"
	}

	if strings.Contains(selector, "data-testid") || strings.Contains(selector, "test-id") {
		return selector, "testid"
	}

	return selector, "css"
}

// waitForElementActionable waits for the element to be actionable before clicking.
// Called only when force=false; force=true skips this entirely.
func (s *ClickElementTool) waitForElementActionable(ctx context.Context, session *playwright.BrowserSession, selector string, timeoutMs int) error {
	timeoutDuration := time.Duration(timeoutMs) * time.Millisecond

	err := s.playwright.WaitForCondition(ctx, session.ID, "selector", selector, "visible", timeoutDuration, "")
	if err != nil {
		s.logger.Warn("element not visible in main frame, checking iframes",
			zap.String("selector", selector),
			zap.Error(err))
		return s.checkElementInIframes(session, selector)
	}

	return nil
}

// checkElementInIframes reports whether the missing element might live in a
// nested iframe (cross-frame click is not yet supported). It only inspects
// the iframe count; the returned error tells the caller why the click cannot
// proceed.
func (s *ClickElementTool) checkElementInIframes(session *playwright.BrowserSession, selector string) error {
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

	s.logger.Info("element absent from main frame, iframes detected",
		zap.String("selector", selector),
		zap.Int("iframe_count", iframeCount))

	return fmt.Errorf("element not found in main frame, %d iframes detected but cross-frame clicking not yet implemented; pass force=true to click without the visibility wait", iframeCount)
}
