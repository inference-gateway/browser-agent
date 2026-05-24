package tools

import (
	"context"
	"fmt"
	"time"

	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

var (
	validWaitConditionTypes = []string{"selector", "navigation", "function", "timeout", "networkidle"}
	validSelectorStates     = []string{"visible", "hidden", "attached", "detached"}
)

// WaitForConditionTool struct holds the tool with dependencies
type WaitForConditionTool struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewWaitForConditionTool creates a new wait_for_condition tool
func NewWaitForConditionTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	tool := &WaitForConditionTool{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"wait_for_condition",
		"Wait for specific conditions before proceeding with automation",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"condition": map[string]any{
					"description": "Type of condition (selector, navigation, function, timeout, networkidle)",
					"type":        "string",
				},
				"custom_function": map[string]any{
					"description": "Custom JavaScript function to evaluate for 'function' condition",
					"type":        "string",
				},
				"selector": map[string]any{
					"description": "Selector to wait for if condition is 'selector'",
					"type":        "string",
				},
				"state": map[string]any{
					"default":     "visible",
					"description": "State to wait for (visible, hidden, attached, detached)",
					"type":        "string",
				},
				"timeout": map[string]any{
					"default":     defaultTimeoutMs,
					"description": "Maximum time to wait in milliseconds",
					"type":        "integer",
				},
			},
			"required": []string{"condition"},
		},
		tool.WaitForConditionHandler,
	)
}

// WaitForConditionHandler handles the wait_for_condition tool execution.
//
// Note on 'networkidle': previously this branch injected JavaScript that
// permanently replaced window.fetch and window.XMLHttpRequest on the page
// (and never restored the originals, polluting every subsequent request
// in the same session). It is now delegated to Playwright's native
// page.WaitForLoadState(LoadStateNetworkidle) via the playwright service,
// which has no such side effects.
func (s *WaitForConditionTool) WaitForConditionHandler(ctx context.Context, args map[string]any) (string, error) {
	condition, err := requiredString(args, "condition")
	if err != nil {
		return "", err
	}
	if !oneOf(condition, validWaitConditionTypes...) {
		return "", fmt.Errorf("invalid condition type: %s. Must be one of: %v", condition, validWaitConditionTypes)
	}

	selector, err := stringArg(args, "selector", "")
	if err != nil {
		return "", err
	}

	state, err := stringArg(args, "state", "visible")
	if err != nil {
		return "", err
	}
	if !oneOf(state, validSelectorStates...) {
		return "", fmt.Errorf("invalid state: %s. Must be one of: %v", state, validSelectorStates)
	}

	timeout, err := boundedIntArg(args, "timeout", defaultTimeoutMs, minTimeoutMs, maxTimeoutMs)
	if err != nil {
		return "", err
	}

	customFunction, err := stringArg(args, "custom_function", "")
	if err != nil {
		return "", err
	}

	if err := validateConditionRequirements(condition, selector, customFunction); err != nil {
		return "", err
	}

	s.logger.Info("waiting for condition",
		zap.String("condition", condition),
		zap.String("selector", selector),
		zap.String("state", state),
		zap.Int("timeout_ms", timeout))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	timeoutDuration := time.Duration(timeout) * time.Millisecond
	startTime := time.Now()

	if err := s.playwright.WaitForCondition(ctx, session.ID, condition, selector, state, timeoutDuration, customFunction); err != nil {
		s.logger.Error("wait condition failed",
			zap.String("condition", condition),
			zap.String("selector", selector),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("wait condition failed: %w", err)
	}

	actualWaitTime := time.Since(startTime).Milliseconds()

	s.logger.Info("wait condition completed successfully",
		zap.String("condition", condition),
		zap.String("sessionID", session.ID),
		zap.Int64("actual_wait_ms", actualWaitTime))

	return marshalResponse(map[string]any{
		"success":         true,
		"condition":       condition,
		"selector":        selector,
		"state":           state,
		"timeout_ms":      timeout,
		"actual_wait_ms":  actualWaitTime,
		"session_id":      session.ID,
		"custom_function": customFunction,
		"message":         "Wait condition completed successfully",
	})
}

// validateConditionRequirements validates condition-specific requirements.
// Standalone function (not a method) so it can be called from tests without
// constructing a tool.
func validateConditionRequirements(condition, selector, customFunction string) error {
	switch condition {
	case "selector":
		if selector == "" {
			return fmt.Errorf("selector parameter is required for selector condition")
		}
	case "function":
		if customFunction == "" {
			return fmt.Errorf("custom_function parameter is required for function condition")
		}
	}
	return nil
}
