package tools

import (
	"context"
	"fmt"

	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

var validAuthTypes = []string{"basic", "form", "oauth"}

// HandleAuthenticationTool struct holds the tool with dependencies
type HandleAuthenticationTool struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewHandleAuthenticationTool creates a new handle_authentication tool
func NewHandleAuthenticationTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	tool := &HandleAuthenticationTool{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"handle_authentication",
		"NOT YET IMPLEMENTED: This tool currently returns an error explaining that authentication is not wired through. For basic auth, configure HTTP credentials at the browser context level. For form login, compose navigate_to_url + fill_form + click_element. For OAuth, run the flow manually with the above primitives.",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"login_url": map[string]any{
					"description": "URL of the login page for form authentication",
					"type":        "string",
				},
				"password": map[string]any{
					"description": "Password for authentication",
					"type":        "string",
				},
				"password_selector": map[string]any{
					"description": "Selector for password field in form authentication",
					"type":        "string",
				},
				"submit_selector": map[string]any{
					"description": "Selector for submit button in form authentication",
					"type":        "string",
				},
				"type": map[string]any{
					"description": "Authentication type (basic, form, oauth)",
					"type":        "string",
				},
				"username": map[string]any{
					"description": "Username or email for authentication",
					"type":        "string",
				},
				"username_selector": map[string]any{
					"description": "Selector for username field in form authentication",
					"type":        "string",
				},
			},
			"required": []string{"type"},
		},
		tool.HandleAuthenticationHandler,
	)
}

// HandleAuthenticationHandler returns an explicit "not implemented" error
// instead of the previous silent no-op TODO that returned a fake-success
// payload. We still validate args["type"] so the caller gets a clearer
// signal when they passed something garbage in addition to hitting an
// unimplemented tool.
//
// The plumbing in internal/playwright.HandleAuthentication exists but the
// wiring between this tool's args (login_url, *_selector, ...) and the
// service's expected layout (selectors map[string]string) was never
// completed. Until it is, fail loudly.
func (s *HandleAuthenticationTool) HandleAuthenticationHandler(ctx context.Context, args map[string]any) (string, error) {
	authType, err := requiredString(args, "type")
	if err != nil {
		return "", err
	}
	if !oneOf(authType, validAuthTypes...) {
		return "", fmt.Errorf("invalid auth type: %s. Must be one of: %v", authType, validAuthTypes)
	}

	s.logger.Warn("handle_authentication invoked but not implemented",
		zap.String("auth_type", authType))

	return "", fmt.Errorf("handle_authentication is not yet implemented (auth_type=%s). "+
		"For basic auth, configure HTTP credentials at the browser context level. "+
		"For form login, compose navigate_to_url + fill_form + click_element instead. "+
		"For OAuth, drive the flow manually with the same primitives", authType)
}
