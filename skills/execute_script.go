package skills

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

// ExecuteScriptSkill struct holds the skill with dependencies
type ExecuteScriptSkill struct {
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

// NewExecuteScriptSkill creates a new execute_script skill
func NewExecuteScriptSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &ExecuteScriptSkill{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"execute_script",
		"Execute custom JavaScript code in the browser context",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"script": map[string]any{
					"description": "JavaScript code to execute",
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
		skill.ExecuteScriptHandler,
	)
}

// ExecuteScriptHandler handles the execute_script skill execution
func (s *ExecuteScriptSkill) ExecuteScriptHandler(ctx context.Context, args map[string]any) (string, error) {
	startTime := time.Now()

	script, ok := args["script"].(string)
	if !ok || script == "" {
		return "", fmt.Errorf("script parameter is required and must be a non-empty string")
	}

	if err := s.validateScriptSecurity(script); err != nil {
		s.logger.Error("script security validation failed", zap.String("script", script), zap.Error(err))
		return "", fmt.Errorf("script security validation failed: %w", err)
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

	session, err := s.getOrCreateSession(ctx)
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

// validateScriptSecurity performs basic security validation on the script
func (s *ExecuteScriptSkill) validateScriptSecurity(script string) error {
	dangerousPatterns := []string{
		// File system access
		"require\\s*\\(\\s*['\"]fs['\"]",
		"require\\s*\\(\\s*['\"]path['\"]",
		"require\\s*\\(\\s*['\"]os['\"]",
		// Network access
		"require\\s*\\(\\s*['\"]http['\"]",
		"require\\s*\\(\\s*['\"]https['\"]",
		"require\\s*\\(\\s*['\"]net['\"]",
		// Process execution
		"require\\s*\\(\\s*['\"]child_process['\"]",
		"exec\\s*\\(",
		"spawn\\s*\\(",
		// Eval and dynamic code execution
		"eval\\s*\\(",
		"function\\s*\\(",
		"settimeout\\s*\\(",
		"setinterval\\s*\\(",
		// Global object access
		"global\\.",
		"process\\.",
		"__dirname",
		"__filename",
		// Sensitive browser APIs that could be misused
		"localStorage\\.clear",
		"sessionStorage\\.clear",
		"document\\.cookie\\s*=",
		"window\\.location\\s*=",
	}

	scriptLower := strings.ToLower(script)

	for _, pattern := range dangerousPatterns {
		matched, err := regexp.MatchString(pattern, scriptLower)
		if err != nil {
			s.logger.Warn("regex pattern validation failed", zap.String("pattern", pattern), zap.Error(err))
			continue
		}
		if matched {
			return fmt.Errorf("script contains potentially dangerous pattern: %s", pattern)
		}
	}

	if len(script) > 50000 {
		return fmt.Errorf("script too large: %d characters (max 50000)", len(script))
	}

	return nil
}

// prepareScript prepares the script for execution, wrapping async scripts if needed
func (s *ExecuteScriptSkill) prepareScript(script string, isAsync bool) (string, error) {
	if !isAsync {
		return script, nil
	}

	wrappedScript := fmt.Sprintf(`
(async function() {
	try {
		%s
	} catch (error) {
		throw error;
	}
})()`, script)

	return wrappedScript, nil
}

// calculateScriptHash creates a simple hash of the script for logging/tracking
func (s *ExecuteScriptSkill) calculateScriptHash(script string) string {
	if len(script) <= 32 {
		return fmt.Sprintf("script_%d_chars", len(script))
	}
	return fmt.Sprintf("script_%d_chars_%x", len(script), script[:32])
}

// getResultType determines the type of the result for metadata
func (s *ExecuteScriptSkill) getResultType(result any) string {
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

// getOrCreateSession gets a task-scoped isolated session
func (s *ExecuteScriptSkill) getOrCreateSession(ctx context.Context) (*playwright.BrowserSession, error) {
	return s.playwright.GetOrCreateTaskSession(ctx)
}
