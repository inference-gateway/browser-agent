package skills

import (
	"context"
	"fmt"
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
	// Extract parameters from args
	url, ok := args["url"].(string)
	if !ok {
		return "", fmt.Errorf("url parameter is required")
	}

	timeout := 30000
	if t, ok := args["timeout"].(int); ok {
		timeout = t
	}

	waitUntil := "load"
	if w, ok := args["wait_until"].(string); ok {
		waitUntil = w
	}

	s.logger.Info("navigating to URL",
		zap.String("url", url),
		zap.String("wait_until", waitUntil),
		zap.Int("timeout", timeout))

	config := &playwright.BrowserConfig{
		Engine:         playwright.Chromium,
		Headless:       true,
		Timeout:        time.Duration(timeout) * time.Millisecond,
		ViewportWidth:  1280,
		ViewportHeight: 720,
		Args:           []string{"--disable-dev-shm-usage", "--no-sandbox"},
	}

	session, err := s.playwright.LaunchBrowser(ctx, config)
	if err != nil {
		return "", fmt.Errorf("failed to launch browser: %w", err)
	}

	defer func() {
		if closeErr := s.playwright.CloseBrowser(ctx, session.ID); closeErr != nil {
			s.logger.Error("failed to close browser session", zap.Error(closeErr))
		}
	}()

	err = s.playwright.NavigateToURL(ctx, session.ID, url, waitUntil, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return "", fmt.Errorf("failed to navigate to URL: %w", err)
	}

	return fmt.Sprintf(`{"result": "Successfully navigated to %s", "sessionID": "%s", "url": "%s", "wait_until": "%s"}`,
		url, session.ID, url, waitUntil), nil
}
