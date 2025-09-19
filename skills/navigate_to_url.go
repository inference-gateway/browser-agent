package skills

import (
	"context"
	"fmt"
	"net/url"
	"time"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/playwright-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// NavigateToURLSkill struct holds the skill with dependencies
type NavigateToURLSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewNavigateToURLSkill creates a new navigate_to_url skill
func NewNavigateToURLSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &NavigateToURLSkill{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"navigate_to_url",
		"Navigate to a specific URL and wait for the page to fully load",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"timeout": map[string]any{
					"default":     30000,
					"description": "Maximum navigation timeout in milliseconds",
					"type":        "integer",
				},
				"url": map[string]any{
					"description": "The URL to navigate to",
					"type":        "string",
				},
				"wait_until": map[string]any{
					"default":     "load",
					"description": "When to consider navigation succeeded (domcontentloaded, load, networkidle)",
					"type":        "string",
				},
			},
			"required": []string{"url"},
		},
		skill.NavigateToURLHandler,
	)
}

// NavigateToURLHandler handles the navigate_to_url skill execution
func (s *NavigateToURLSkill) NavigateToURLHandler(ctx context.Context, args map[string]any) (string, error) {
	// Extract and validate parameters from args
	url, ok := args["url"].(string)
	if !ok || url == "" {
		return "", fmt.Errorf("url parameter is required and must be a non-empty string")
	}

	// Validate and normalize URL format
	normalizedURL, err := s.validateAndNormalizeURL(url)
	if err != nil {
		s.logger.Error("invalid URL provided", zap.String("url", url), zap.Error(err))
		return "", fmt.Errorf("invalid URL: %w", err)
	}
	url = normalizedURL

	// Extract wait_until parameter with default
	waitUntil := "load"
	if wu, ok := args["wait_until"].(string); ok && wu != "" {
		if !s.isValidWaitCondition(wu) {
			return "", fmt.Errorf("invalid wait_until value: %s. Must be one of: domcontentloaded, load, networkidle", wu)
		}
		waitUntil = wu
	}

	// Extract timeout parameter with default
	timeout := 30000
	if t, ok := args["timeout"].(int); ok && t > 0 {
		timeout = t
	} else if tf, ok := args["timeout"].(float64); ok && tf > 0 {
		timeout = int(tf)
	}

	s.logger.Info("navigating to URL",
		zap.String("url", url),
		zap.String("wait_until", waitUntil),
		zap.Int("timeout_ms", timeout))

	// Get or create a browser session
	session, err := s.getOrCreateSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	// Perform navigation
	timeoutDuration := time.Duration(timeout) * time.Millisecond
	err = s.playwright.NavigateToURL(ctx, session.ID, url, waitUntil, timeoutDuration)
	if err != nil {
		s.logger.Error("navigation failed",
			zap.String("url", url),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("navigation failed: %w", err)
	}

	s.logger.Info("navigation completed successfully",
		zap.String("url", url),
		zap.String("sessionID", session.ID))

	// Return success response
	response := map[string]interface{}{
		"success":    true,
		"url":        url,
		"wait_until": waitUntil,
		"timeout_ms": timeout,
		"session_id": session.ID,
		"message":    "Navigation completed successfully",
	}

	return fmt.Sprintf(`%+v`, response), nil
}

// validateAndNormalizeURL validates that the provided URL is well-formed and supported, returning the normalized URL
func (s *NavigateToURLSkill) validateAndNormalizeURL(urlStr string) (string, error) {
	if urlStr == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Parse URL first to check if it already has a scheme
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// If no scheme is provided, add https:// prefix and re-parse
	if parsedURL.Scheme == "" {
		urlStr = "https://" + urlStr
		parsedURL, err = url.Parse(urlStr)
		if err != nil {
			return "", fmt.Errorf("invalid URL format: %w", err)
		}
	}

	// Check for valid scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme: %s. Only http and https are supported", parsedURL.Scheme)
	}

	// Check for valid host
	if parsedURL.Host == "" {
		return "", fmt.Errorf("URL must include a valid host")
	}

	// Return the normalized URL
	return parsedURL.String(), nil
}

// isValidWaitCondition validates the wait_until parameter
func (s *NavigateToURLSkill) isValidWaitCondition(condition string) bool {
	validConditions := []string{"domcontentloaded", "load", "networkidle"}
	for _, valid := range validConditions {
		if condition == valid {
			return true
		}
	}
	return false
}

// getOrCreateSession gets an existing session or creates a new one
func (s *NavigateToURLSkill) getOrCreateSession(ctx context.Context) (*playwright.BrowserSession, error) {
	// For simplicity, we'll create a new session each time
	// In a production system, you might want to reuse sessions or manage them more efficiently
	config := playwright.DefaultBrowserConfig()
	session, err := s.playwright.LaunchBrowser(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	return session, nil
}
