package tools

import (
	"context"
	"fmt"
	"net/url"
	"time"

	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

var validWaitConditions = []string{"domcontentloaded", "load", "networkidle"}

// NavigateToURLTool struct holds the tool with dependencies
type NavigateToURLTool struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewNavigateToURLTool creates a new navigate_to_url tool
func NewNavigateToURLTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	tool := &NavigateToURLTool{
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
					"default":     defaultTimeoutMs,
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
		tool.NavigateToURLHandler,
	)
}

// NavigateToURLHandler handles the navigate_to_url tool execution
func (s *NavigateToURLTool) NavigateToURLHandler(ctx context.Context, args map[string]any) (string, error) {
	rawURL, err := requiredString(args, "url")
	if err != nil {
		return "", err
	}

	targetURL, err := s.validateAndNormalizeURL(rawURL)
	if err != nil {
		s.logger.Error("invalid URL provided", zap.String("url", rawURL), zap.Error(err))
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	waitUntil, err := stringArg(args, "wait_until", "load")
	if err != nil {
		return "", err
	}
	if !oneOf(waitUntil, validWaitConditions...) {
		return "", fmt.Errorf("invalid wait_until value: %s. Must be one of: %v", waitUntil, validWaitConditions)
	}

	timeout, err := boundedIntArg(args, "timeout", defaultTimeoutMs, minTimeoutMs, maxTimeoutMs)
	if err != nil {
		return "", err
	}

	s.logger.Info("navigating to URL",
		zap.String("url", targetURL),
		zap.String("wait_until", waitUntil),
		zap.Int("timeout_ms", timeout))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	timeoutDuration := time.Duration(timeout) * time.Millisecond
	if err := s.playwright.NavigateToURL(ctx, session.ID, targetURL, waitUntil, timeoutDuration); err != nil {
		s.logger.Error("navigation failed",
			zap.String("url", targetURL),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("navigation failed: %w", err)
	}

	s.logger.Info("navigation completed successfully",
		zap.String("url", targetURL),
		zap.String("sessionID", session.ID))

	return marshalResponse(map[string]any{
		"success":    true,
		"url":        targetURL,
		"wait_until": waitUntil,
		"timeout_ms": timeout,
		"session_id": session.ID,
		"message":    "Navigation completed successfully",
	})
}

// validateAndNormalizeURL validates that the provided URL is well-formed and supported, returning the normalized URL
func (s *NavigateToURLTool) validateAndNormalizeURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme == "" {
		urlStr = "https://" + urlStr
		parsedURL, err = url.Parse(urlStr)
		if err != nil {
			return "", fmt.Errorf("invalid URL format: %w", err)
		}
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme: %s. Only http and https are supported", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("URL must include a valid host")
	}

	return parsedURL.String(), nil
}
