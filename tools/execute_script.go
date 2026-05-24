package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// ExecuteScriptTool struct holds the tool with dependencies
type ExecuteScriptTool struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// ScriptExecutionResult represents the result of script execution
type ScriptExecutionResult struct {
	Success     bool           `json:"success"`
	Result      any            `json:"result"`
	ResultType  string         `json:"result_type"`
	Error       string         `json:"error,omitempty"`
	ExecutionMS int64          `json:"execution_ms"`
	SessionID   string         `json:"session_id"`
	Timestamp   string         `json:"timestamp"`
	ScriptHash  string         `json:"script_hash"`
	Message     string         `json:"message"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// executeScriptToolDescription is the description shown to the LLM when the
// tool is registered. It deliberately spells out the browser/Playwright
// execution context so the model does not waste turns trying Node.js APIs.
const executeScriptToolDescription = "Execute custom JavaScript inside the current page via Playwright's " +
	"page.evaluate(). The script runs in the browser context, NOT in Node.js: " +
	"globals like window, document, navigator, fetch and localStorage are " +
	"available; Node.js built-ins (require, process, __dirname, __filename, " +
	"fs, path, os, http, https, child_process, etc.) are NOT available and " +
	"calls to them will be rejected. Use browser/DOM APIs only. The script " +
	"body is automatically wrapped in an IIFE, so a top-level `return` is " +
	"valid. Set async=true if the body uses `await`."

// NewExecuteScriptTool creates a new execute_script tool
func NewExecuteScriptTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	tool := &ExecuteScriptTool{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"execute_script",
		executeScriptToolDescription,
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"script": map[string]any{
					"description": "JavaScript code to execute in the browser (Playwright page.evaluate context). Use DOM/Web APIs only - Node.js built-ins are unavailable.",
					"type":        "string",
				},
				"args": map[string]any{
					"description": "Arguments to pass to the script (will be available as arguments[0], arguments[1], etc.)",
					"items":       map[string]any{},
					"type":        "array",
					"default":     []any{},
				},
				"return_value": map[string]any{
					"default":     true,
					"description": "Whether to return the script execution result",
					"type":        "boolean",
				},
				"timeout": map[string]any{
					"default":     30000,
					"description": "Maximum script execution timeout in milliseconds",
					"type":        "integer",
				},
				"async": map[string]any{
					"default":     false,
					"description": "Whether the script contains async operations (will wrap in async function)",
					"type":        "boolean",
				},
			},
			"required": []string{"script"},
		},
		tool.ExecuteScriptHandler,
	)
}

// ExecuteScriptHandler handles the execute_script tool execution
func (s *ExecuteScriptTool) ExecuteScriptHandler(ctx context.Context, args map[string]any) (string, error) {
	startTime := time.Now()

	script, ok := args["script"].(string)
	if !ok || script == "" {
		return "", fmt.Errorf("script parameter is required and must be a non-empty string")
	}

	if err := s.validateScriptSecurity(script); err != nil {
		s.logger.Error("script security validation failed", zap.String("script", script), zap.Error(err))
		return "", fmt.Errorf("execute_script rejected the script: %w. execute_script runs inside the browser via Playwright page.evaluate(); Node.js built-ins (require, process, fs, path, os, http, https, child_process, __dirname, __filename, etc.) are unavailable by design. Use browser/DOM APIs only (window, document, fetch, localStorage, ...)", err)
	}

	scriptArgs := []any{}
	if argsVal, ok := args["args"]; ok {
		if argsSlice, ok := argsVal.([]any); ok {
			scriptArgs = argsSlice
		}
	}

	returnValue := true
	if rv, ok := args["return_value"].(bool); ok {
		returnValue = rv
	}

	timeout := 30000
	if t, ok := args["timeout"].(int); ok && t > 0 {
		timeout = t
	} else if tf, ok := args["timeout"].(float64); ok && tf > 0 {
		timeout = int(tf)
	}

	isAsync := false
	if async, ok := args["async"].(bool); ok {
		isAsync = async
	}

	processedScript, err := s.prepareScript(script, isAsync)
	if err != nil {
		s.logger.Error("script preparation failed", zap.Error(err))
		return "", fmt.Errorf("script preparation failed: %w", err)
	}

	s.logger.Info("executing script",
		zap.String("script_hash", s.calculateScriptHash(script)),
		zap.Int("args_count", len(scriptArgs)),
		zap.Bool("return_value", returnValue),
		zap.Int("timeout_ms", timeout),
		zap.Bool("async", isAsync))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	var result any
	if returnValue {
		result, err = s.playwright.ExecuteScript(ctx, session.ID, processedScript, scriptArgs)
	} else {
		_, err = s.playwright.ExecuteScript(ctx, session.ID, processedScript, scriptArgs)
		result = nil
	}

	executionTime := time.Since(startTime)
	scriptHash := s.calculateScriptHash(script)
	timestamp := time.Now().Format(time.RFC3339)

	scriptResult := &ScriptExecutionResult{
		Success:     err == nil,
		Result:      result,
		ResultType:  s.getResultType(result),
		ExecutionMS: executionTime.Milliseconds(),
		SessionID:   session.ID,
		Timestamp:   timestamp,
		ScriptHash:  scriptHash,
		Metadata: map[string]any{
			"args_count":    len(scriptArgs),
			"return_value":  returnValue,
			"timeout_ms":    timeout,
			"async":         isAsync,
			"script_length": len(script),
			"processed":     processedScript != script,
		},
	}

	if err != nil {
		scriptResult.Error = err.Error()
		scriptResult.Message = "Script execution failed"
		s.logger.Error("script execution failed",
			zap.String("sessionID", session.ID),
			zap.String("script_hash", scriptHash),
			zap.Error(err))
	} else {
		scriptResult.Message = "Script executed successfully"
		s.logger.Info("script execution completed successfully",
			zap.String("sessionID", session.ID),
			zap.String("script_hash", scriptHash),
			zap.String("result_type", scriptResult.ResultType),
			zap.Int64("execution_ms", scriptResult.ExecutionMS))
	}

	responseJSON, err := json.Marshal(scriptResult)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(responseJSON), nil
}

// dangerousScriptPattern pairs a regex with an LLM-actionable reason
// explaining why the match is rejected. Patterns marked caseSensitive=true
// are matched against the original script (not the lowercased form); the
// rest are matched against strings.ToLower(script) for cheap normalization.
type dangerousScriptPattern struct {
	pattern       string
	reason        string
	caseSensitive bool

	compiled *regexp.Regexp
}

// dangerousScriptPatterns enumerates patterns the validator rejects together
// with the human-readable reason that is surfaced back to the caller. The
// reasons are written for an LLM consumer: they name the offending construct,
// explain why it cannot work in the browser execution context, and (where
// useful) point at the correct browser/DOM alternative.
//
// Note: the legacy `function\s*\(` pattern was removed because it matched any
// JavaScript function expression (including the IIFE this tool itself emits
// in prepareScript and any user `(function () { ... })()`). The dynamic-code
// hazard it tried to cover - the `Function` constructor - is now handled by
// the case-sensitive `\bFunction\s*\(` pattern below.
var dangerousScriptPatterns = []*dangerousScriptPattern{
	{
		pattern: `require\s*\(\s*['"]fs['"]`,
		reason:  "uses Node.js `require('fs')` for filesystem access; the browser sandbox has no filesystem - there is no equivalent from page.evaluate()",
	},
	{
		pattern: `require\s*\(\s*['"]path['"]`,
		reason:  "uses Node.js `require('path')`; path manipulation must be done with plain string operations or URL APIs in the browser",
	},
	{
		pattern: `require\s*\(\s*['"]os['"]`,
		reason:  "uses Node.js `require('os')`; OS metadata is not exposed to the browser context - inspect `navigator.userAgent` / `navigator.platform` instead",
	},
	{
		pattern: `require\s*\(\s*['"]http['"]`,
		reason:  "uses Node.js `require('http')`; perform network requests with the browser `fetch()` or `XMLHttpRequest` APIs",
	},
	{
		pattern: `require\s*\(\s*['"]https['"]`,
		reason:  "uses Node.js `require('https')`; perform network requests with the browser `fetch()` or `XMLHttpRequest` APIs",
	},
	{
		pattern: `require\s*\(\s*['"]net['"]`,
		reason:  "uses Node.js `require('net')` for raw sockets; the browser sandbox cannot open raw TCP sockets",
	},
	{
		pattern: `require\s*\(\s*['"]child_process['"]`,
		reason:  "uses Node.js `require('child_process')`; the browser context cannot spawn subprocesses",
	},
	{
		pattern: `\bexec\s*\(`,
		reason:  "calls `exec(...)` (child_process); subprocess execution is not available in the browser context",
	},
	{
		pattern: `\bspawn\s*\(`,
		reason:  "calls `spawn(...)` (child_process); subprocess execution is not available in the browser context",
	},
	{
		pattern: `\beval\s*\(`,
		reason:  "calls `eval(...)`; dynamic code evaluation is blocked. Inline the logic directly instead",
	},
	{
		pattern:       `\bFunction\s*\(`,
		reason:        "uses the `Function` constructor for dynamic code generation; this is blocked for the same reason as eval. Inline the logic directly instead. (Regular `function () { ... }` expressions and IIFEs are fine.)",
		caseSensitive: true,
	},
	{
		pattern: `\bsettimeout\s*\(`,
		reason:  "schedules work with `setTimeout`; long-lived timers cannot outlive a single page.evaluate() call. If you need to wait, use the dedicated `wait_for_condition` tool",
	},
	{
		pattern: `\bsetinterval\s*\(`,
		reason:  "schedules work with `setInterval`; recurring timers cannot outlive a single page.evaluate() call. Use the dedicated `wait_for_condition` tool instead",
	},
	{
		pattern: `\bglobal\.`,
		reason:  "accesses the Node.js `global` object, which does not exist in the browser. Use `window` (or `globalThis`) for the page's global scope",
	},
	{
		pattern: `\bprocess\.`,
		reason:  "accesses Node.js `process.*` (cwd, env, argv, ...). The browser has no `process` global. For environment-like info use `navigator`, `location`, or ask the user to expose it via a dedicated tool",
	},
	{
		pattern: `__dirname`,
		reason:  "references the Node.js `__dirname` global; it does not exist in the browser. Use `location.href` / `document.baseURI` for the current page URL",
	},
	{
		pattern: `__filename`,
		reason:  "references the Node.js `__filename` global; it does not exist in the browser. Use `location.href` for the current page URL",
	},
	{
		pattern: `localstorage\.clear`,
		reason:  "calls `localStorage.clear()`; wiping the page's localStorage is blocked because it destroys session state. Remove individual keys with `localStorage.removeItem(key)` if needed",
	},
	{
		pattern: `sessionstorage\.clear`,
		reason:  "calls `sessionStorage.clear()`; wiping sessionStorage is blocked because it destroys session state. Remove individual keys with `sessionStorage.removeItem(key)` if needed",
	},
	{
		pattern: `document\.cookie\s*=`,
		reason:  "writes to `document.cookie`; cookie mutation from a script is blocked to protect the session. Use Playwright-level auth helpers / `handle_authentication` instead",
	},
	{
		pattern: `window\.location\s*=`,
		reason:  "assigns to `window.location` to force navigation; use the dedicated `navigate_to_url` tool so the agent can track the navigation",
	},
}

func init() {
	for _, p := range dangerousScriptPatterns {
		p.compiled = regexp.MustCompile(p.pattern)
	}
}

// validateScriptSecurity performs basic security validation on the script and
// returns an LLM-actionable error describing exactly why the script was
// rejected and how to express the intent in the browser context.
func (s *ExecuteScriptTool) validateScriptSecurity(script string) error {
	scriptLower := strings.ToLower(script)

	for _, p := range dangerousScriptPatterns {
		target := scriptLower
		if p.caseSensitive {
			target = script
		}
		if p.compiled.MatchString(target) {
			return fmt.Errorf("%s", p.reason)
		}
	}

	if len(script) > 50000 {
		return fmt.Errorf("script is too large: %d characters (max 50000). Split the work into smaller execute_script calls", len(script))
	}

	return nil
}

// prepareScript prepares the script for execution, always wrapping in a function to avoid syntax errors
func (s *ExecuteScriptTool) prepareScript(script string, isAsync bool) (string, error) {
	var wrappedScript string

	if isAsync {
		wrappedScript = fmt.Sprintf(`
(async function() {
	try {
		%s
	} catch (error) {
		throw error;
	}
})()`, script)
	} else {
		wrappedScript = fmt.Sprintf(`
(function() {
	try {
		%s
	} catch (error) {
		throw error;
	}
})()`, script)
	}

	return wrappedScript, nil
}

// calculateScriptHash creates a simple hash of the script for logging/tracking
func (s *ExecuteScriptTool) calculateScriptHash(script string) string {
	if len(script) <= 32 {
		return fmt.Sprintf("script_%d_chars", len(script))
	}
	return fmt.Sprintf("script_%d_chars_%x", len(script), script[:32])
}

// getResultType determines the type of the result for metadata
func (s *ExecuteScriptTool) getResultType(result any) string {
	if result == nil {
		return "null"
	}

	switch result.(type) {
	case bool:
		return "boolean"
	case int, int8, int16, int32, int64, float32, float64:
		return "number"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return fmt.Sprintf("unknown:%T", result)
	}
}
