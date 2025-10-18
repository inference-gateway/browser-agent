package skills

import (
	"context"
	"fmt"
	"time"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// WaitForConditionSkill struct holds the skill with dependencies
type WaitForConditionSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewWaitForConditionSkill creates a new wait_for_condition skill
func NewWaitForConditionSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &WaitForConditionSkill{
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
					"description": "Type of condition (selector, navigation, function, timeout)",
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
					"default":     30000,
					"description": "Maximum time to wait in milliseconds",
					"type":        "integer",
				},
			},
			"required": []string{"condition"},
		},
		skill.WaitForConditionHandler,
	)
}

// WaitForConditionHandler handles the wait_for_condition skill execution
func (s *WaitForConditionSkill) WaitForConditionHandler(ctx context.Context, args map[string]any) (string, error) {
	condition, ok := args["condition"].(string)
	if !ok || condition == "" {
		return "", fmt.Errorf("condition parameter is required and must be a non-empty string")
	}

	if !s.isValidCondition(condition) {
		return "", fmt.Errorf("invalid condition type: %s. Must be one of: selector, navigation, function, timeout, networkidle", condition)
	}

	selector := ""
	if sel, ok := args["selector"].(string); ok {
		selector = sel
	}

	state := "visible"
	if st, ok := args["state"].(string); ok && st != "" {
		if !s.isValidState(st) {
			return "", fmt.Errorf("invalid state: %s. Must be one of: visible, hidden, attached, detached", st)
		}
		state = st
	}

	timeout := 30000
	if t, ok := args["timeout"].(int); ok && t > 0 {
		timeout = t
	} else if tf, ok := args["timeout"].(float64); ok && tf > 0 {
		timeout = int(tf)
	}

	customFunction := ""
	if cf, ok := args["custom_function"].(string); ok {
		customFunction = cf
	}

	if err := s.validateConditionRequirements(condition, selector, customFunction); err != nil {
		return "", err
	}

	s.logger.Info("waiting for condition",
		zap.String("condition", condition),
		zap.String("selector", selector),
		zap.String("state", state),
		zap.Int("timeout_ms", timeout),
		zap.String("custom_function", customFunction))

	session, err := s.getOrCreateSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	timeoutDuration := time.Duration(timeout) * time.Millisecond
	startTime := time.Now()

	err = s.executeWaitCondition(ctx, session.ID, condition, selector, state, timeoutDuration, customFunction)
	if err != nil {
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
		zap.String("selector", selector),
		zap.String("sessionID", session.ID),
		zap.Int64("actual_wait_ms", actualWaitTime))

	response := map[string]any{
		"success":         true,
		"condition":       condition,
		"selector":        selector,
		"state":           state,
		"timeout_ms":      timeout,
		"actual_wait_ms":  actualWaitTime,
		"session_id":      session.ID,
		"message":         "Wait condition completed successfully",
		"custom_function": customFunction,
	}

	return fmt.Sprintf(`%+v`, response), nil
}

// isValidCondition validates the condition type
func (s *WaitForConditionSkill) isValidCondition(condition string) bool {
	validConditions := []string{"selector", "navigation", "function", "timeout", "networkidle"}
	for _, valid := range validConditions {
		if condition == valid {
			return true
		}
	}
	return false
}

// isValidState validates the selector state
func (s *WaitForConditionSkill) isValidState(state string) bool {
	validStates := []string{"visible", "hidden", "attached", "detached"}
	for _, valid := range validStates {
		if state == valid {
			return true
		}
	}
	return false
}

// validateConditionRequirements validates condition-specific requirements
func (s *WaitForConditionSkill) validateConditionRequirements(condition, selector, customFunction string) error {
	switch condition {
	case "selector":
		if selector == "" {
			return fmt.Errorf("selector parameter is required for selector condition")
		}
	case "function":
		if customFunction == "" {
			return fmt.Errorf("custom_function parameter is required for function condition")
		}
	case "navigation", "timeout", "networkidle":
	}
	return nil
}

// executeWaitCondition executes the appropriate wait operation based on condition type
func (s *WaitForConditionSkill) executeWaitCondition(ctx context.Context, sessionID, condition, selector, state string, timeout time.Duration, customFunction string) error {
	switch condition {
	case "selector":
		return s.playwright.WaitForCondition(ctx, sessionID, condition, selector, state, timeout, "")
	case "function":
		return s.playwright.WaitForCondition(ctx, sessionID, condition, "", "", timeout, customFunction)
	case "navigation":
		return s.playwright.WaitForCondition(ctx, sessionID, condition, "", "", timeout, "")
	case "timeout":
		return s.playwright.WaitForCondition(ctx, sessionID, condition, "", "", timeout, "")
	case "networkidle":
		networkIdleFunction := `
			() => {
				return new Promise((resolve) => {
					let timeout;
					let requestCount = 0;
					
					// Monitor fetch requests
					const originalFetch = window.fetch;
					window.fetch = function(...args) {
						requestCount++;
						return originalFetch.apply(this, args).finally(() => {
							requestCount--;
							if (requestCount === 0) {
								clearTimeout(timeout);
								timeout = setTimeout(() => resolve(true), 500);
							}
						});
					};
					
					// Monitor XMLHttpRequest
					const originalXHR = window.XMLHttpRequest;
					window.XMLHttpRequest = function() {
						const xhr = new originalXHR();
						const originalSend = xhr.send;
						xhr.send = function(...args) {
							requestCount++;
							xhr.addEventListener('loadend', () => {
								requestCount--;
								if (requestCount === 0) {
									clearTimeout(timeout);
									timeout = setTimeout(() => resolve(true), 500);
								}
							});
							return originalSend.apply(this, args);
						};
						return xhr;
					};
					
					// Initial check
					if (requestCount === 0) {
						timeout = setTimeout(() => resolve(true), 500);
					}
				});
			}
		`
		return s.playwright.WaitForCondition(ctx, sessionID, "function", "", "", timeout, networkIdleFunction)
	default:
		return fmt.Errorf("unsupported condition type: %s", condition)
	}
}

// getOrCreateSession gets a task-scoped isolated session
func (s *WaitForConditionSkill) getOrCreateSession(ctx context.Context) (*playwright.BrowserSession, error) {
	return s.playwright.GetOrCreateTaskSession(ctx)
}
