package tools

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	zap "go.uber.org/zap"

	assert "github.com/stretchr/testify/assert"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

func TestExecuteScriptTool_ExecuteScriptHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name          string
		args          map[string]any
		executeResult any
		executeError  error
		expectError   bool
		expectResult  func(t *testing.T, result string)
	}{
		{
			name: "successful simple script execution",
			args: map[string]any{
				"script": "return 42;",
			},
			executeResult: 42,
			executeError:  nil,
			expectError:   false,
			expectResult: func(t *testing.T, result string) {
				var scriptResult ScriptExecutionResult
				err := json.Unmarshal([]byte(result), &scriptResult)
				assert.NoError(t, err)
				assert.True(t, scriptResult.Success)
				assert.Equal(t, float64(42), scriptResult.Result)
				assert.Equal(t, "number", scriptResult.ResultType)
				assert.Equal(t, "Script executed successfully", scriptResult.Message)
			},
		},
		{
			name: "script execution with arguments",
			args: map[string]any{
				"script": "return arguments[0] + arguments[1];",
				"args":   []any{10, 20},
			},
			executeResult: 30,
			executeError:  nil,
			expectError:   false,
			expectResult: func(t *testing.T, result string) {
				var scriptResult ScriptExecutionResult
				err := json.Unmarshal([]byte(result), &scriptResult)
				assert.NoError(t, err)
				assert.True(t, scriptResult.Success)
				assert.Equal(t, float64(30), scriptResult.Result)
				assert.Equal(t, "number", scriptResult.ResultType)
			},
		},
		{
			name: "async script execution",
			args: map[string]any{
				"script": "const result = await Promise.resolve('async result'); return result;",
				"async":  true,
			},
			executeResult: "async result",
			executeError:  nil,
			expectError:   false,
			expectResult: func(t *testing.T, result string) {
				var scriptResult ScriptExecutionResult
				err := json.Unmarshal([]byte(result), &scriptResult)
				assert.NoError(t, err)
				assert.True(t, scriptResult.Success)
				assert.Equal(t, "async result", scriptResult.Result)
				assert.Equal(t, "string", scriptResult.ResultType)
			},
		},
		{
			name: "script execution without return value",
			args: map[string]any{
				"script":       "console.log('test');",
				"return_value": false,
			},
			executeResult: nil,
			executeError:  nil,
			expectError:   false,
			expectResult: func(t *testing.T, result string) {
				var scriptResult ScriptExecutionResult
				err := json.Unmarshal([]byte(result), &scriptResult)
				assert.NoError(t, err)
				assert.True(t, scriptResult.Success)
				assert.Nil(t, scriptResult.Result)
				assert.Equal(t, "null", scriptResult.ResultType)
			},
		},
		{
			name: "script execution failure",
			args: map[string]any{
				"script": "throw new Error('Script error');",
			},
			executeResult: nil,
			executeError:  errors.New("Script error"),
			expectError:   false,
			expectResult: func(t *testing.T, result string) {
				var scriptResult ScriptExecutionResult
				err := json.Unmarshal([]byte(result), &scriptResult)
				assert.NoError(t, err)
				assert.False(t, scriptResult.Success)
				assert.Equal(t, "Script error", scriptResult.Error)
				assert.Equal(t, "Script execution failed", scriptResult.Message)
			},
		},
		{
			name: "empty script should fail",
			args: map[string]any{
				"script": "",
			},
			expectError: true,
		},
		{
			name:        "missing script parameter should fail",
			args:        map[string]any{},
			expectError: true,
		},
		{
			name: "script with dangerous pattern should fail",
			args: map[string]any{
				"script": "require('fs').readFileSync('/etc/passwd');",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlaywright := &mocks.FakeBrowserAutomation{}

			session := &playwright.BrowserSession{
				ID:       "test-session-123",
				Created:  time.Now(),
				LastUsed: time.Now(),
			}
			mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
			mockPlaywright.GetSessionReturns(session, nil)
			mockPlaywright.ExecuteScriptReturns(tt.executeResult, tt.executeError)

			tool := &ExecuteScriptTool{
				logger:     logger,
				playwright: mockPlaywright,
			}

			result, err := tool.ExecuteScriptHandler(context.Background(), tt.args)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectResult != nil {
					tt.expectResult(t, result)
				}
			}
		})
	}
}

func TestExecuteScriptTool_validateScriptSecurity(t *testing.T) {
	logger := zap.NewNop()
	tool := &ExecuteScriptTool{
		logger: logger,
	}

	tests := []struct {
		name          string
		script        string
		expectError   bool
		errorContains string
	}{
		{
			name:        "safe script should pass",
			script:      "return document.querySelector('h1').textContent;",
			expectError: false,
		},
		{
			name:        "function expression should pass (not a Function constructor)",
			script:      "const add = function (a, b) { return a + b; }; return add(1, 2);",
			expectError: false,
		},
		{
			name:        "IIFE should pass (regression: legacy function\\s*\\( pattern)",
			script:      "return (function () { return 42; })();",
			expectError: false,
		},
		{
			name:        "async IIFE should pass",
			script:      "return (async function () { return await Promise.resolve(1); })();",
			expectError: false,
		},
		{
			name:        "arrow function with setTimeout-like identifier in string should pass",
			script:      "return 'documented setTimeout usage';",
			expectError: false,
		},
		{
			name:        "dotted property access process should pass",
			script:      "return obj.process.value;",
			expectError: false,
		},
		{
			name:        "regex .exec on a user array should pass",
			script:      "return /foo/.exec(text);",
			expectError: false,
		},
		{
			name:        "user array .exec should pass",
			script:      "return arr.exec(/re/);",
			expectError: false,
		},
		{
			name:        "user object .global should pass",
			script:      "return el.global.foo;",
			expectError: false,
		},
		{
			name:        "user object .spawn should pass",
			script:      "return queue.spawn(worker);",
			expectError: false,
		},
		{
			name:        "user object .eval method should pass",
			script:      "return calculator.eval(expression);",
			expectError: false,
		},
		{
			name:        "user object .setTimeout should pass",
			script:      "return ticker.setTimeout(cb);",
			expectError: false,
		},
		{
			name:        "user identifier ending in __dirname should pass",
			script:      "return obj.__dirname;",
			expectError: false,
		},
		{
			name:        "user property localStorage.clear on a non-global should pass",
			script:      "return cache.localStorage.clear();",
			expectError: false,
		},
		{
			name:          "script with file system access should fail with actionable reason",
			script:        "require('fs').readFileSync('/etc/passwd');",
			expectError:   true,
			errorContains: "browser sandbox has no filesystem",
		},
		{
			name:          "script with eval should fail with actionable reason",
			script:        "eval('malicious code');",
			expectError:   true,
			errorContains: "eval",
		},
		{
			name:          "script with Function constructor should fail with actionable reason",
			script:        "new Function('return 1')();",
			expectError:   true,
			errorContains: "Function` constructor",
		},
		{
			name:          "script with setTimeout should fail with actionable reason",
			script:        "settimeout(() => { /* malicious */ }, 1000);",
			expectError:   true,
			errorContains: "wait_for_condition",
		},
		{
			name:          "script with global access should fail with actionable reason",
			script:        "global.process.exit(1);",
			expectError:   true,
			errorContains: "Node.js `global`",
		},
		{
			name:          "script with process.cwd should fail with actionable reason",
			script:        "return process.cwd();",
			expectError:   true,
			errorContains: "`process`",
		},
		{
			name:          "script with __dirname should fail with actionable reason",
			script:        "return __dirname;",
			expectError:   true,
			errorContains: "__dirname",
		},
		{
			name:          "script with window.location assignment should suggest navigate_to_url",
			script:        "window.location = 'https://example.com';",
			expectError:   true,
			errorContains: "navigate_to_url",
		},
		{
			name:          "very long script should fail",
			script:        string(make([]byte, 60000)),
			expectError:   true,
			errorContains: "too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.validateScriptSecurity(tt.script)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
					assert.NotContains(t, err.Error(), `\s*\(`,
						"error message must not leak raw regex syntax to the LLM")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExecuteScriptTool_prepareScript(t *testing.T) {
	logger := zap.NewNop()
	tool := &ExecuteScriptTool{
		logger: logger,
	}

	tests := []struct {
		name           string
		script         string
		isAsync        bool
		expectedResult string
	}{
		{
			name:    "sync script should be wrapped in function",
			script:  "return 42;",
			isAsync: false,
			expectedResult: `
(function() {
	try {
		return 42;
	} catch (error) {
		throw error;
	}
})()`,
		},
		{
			name:    "async script should be wrapped in async function",
			script:  "const result = await Promise.resolve(42); return result;",
			isAsync: true,
			expectedResult: `
(async function() {
	try {
		const result = await Promise.resolve(42); return result;
	} catch (error) {
		throw error;
	}
})()`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.prepareScript(tt.script, tt.isAsync)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestExecuteScriptTool_getResultType(t *testing.T) {
	logger := zap.NewNop()
	tool := &ExecuteScriptTool{
		logger: logger,
	}

	tests := []struct {
		name     string
		result   any
		expected string
	}{
		{"nil result", nil, "null"},
		{"boolean result", true, "boolean"},
		{"number result", 42, "number"},
		{"float result", 3.14, "number"},
		{"string result", "hello", "string"},
		{"array result", []any{1, 2, 3}, "array"},
		{"object result", map[string]any{"key": "value"}, "object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.getResultType(tt.result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExecuteScriptTool_calculateScriptHash(t *testing.T) {
	logger := zap.NewNop()
	tool := &ExecuteScriptTool{
		logger: logger,
	}

	tests := []struct {
		name     string
		script   string
		expected string
	}{
		{
			name:     "short script",
			script:   "return 42;",
			expected: "script_10_chars",
		},
		{
			name:     "long script",
			script:   "this is a very long script that exceeds 32 characters and should be hashed",
			expected: "script_74_chars_7468697320697320612076657279206c6f6e6720736372697074207468617420",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.calculateScriptHash(tt.script)
			assert.Equal(t, tt.expected, result)
		})
	}
}
