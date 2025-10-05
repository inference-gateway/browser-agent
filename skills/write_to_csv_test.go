package skills

import (
	"context"
	"strings"
	"testing"

	server "github.com/inference-gateway/adk/server"
	types "github.com/inference-gateway/adk/types"
	config "github.com/inference-gateway/browser-agent/config"
	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"
	zap "go.uber.org/zap"
)

func TestWriteToCsvHandler(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	mockPlaywright.GetConfigReturns(&config.Config{
		Browser: config.BrowserConfig{
			DataDir: "/tmp",
		},
	})

	skill := &WriteToCsvSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name           string
		args           map[string]any
		expectedError  bool
		expectedRows   int
		validateOutput func(t *testing.T, result string)
	}{
		{
			name: "basic CSV writing",
			args: map[string]any{
				"data": []any{
					map[string]any{"name": "Alice", "age": 30, "city": "New York"},
					map[string]any{"name": "Bob", "age": 25, "city": "San Francisco"},
				},
				"filename": "basic.csv",
			},
			expectedError: false,
			expectedRows:  2,
			validateOutput: func(t *testing.T, result string) {
				if !strings.Contains(result, "2 rows") {
					t.Errorf("Expected result to mention 2 rows, got: %s", result)
				}
				if !strings.Contains(result, "basic.csv") {
					t.Errorf("Expected result to mention basic.csv, got: %s", result)
				}
				if !strings.Contains(result, "artifact") {
					t.Errorf("Expected result to mention artifact, got: %s", result)
				}
			},
		},
		{
			name: "CSV with custom headers",
			args: map[string]any{
				"data": []any{
					map[string]any{"name": "Alice", "age": 30},
					map[string]any{"name": "Bob", "age": 25},
				},
				"filename": "custom_headers.csv",
				"headers":  []any{"name", "age"},
			},
			expectedError: false,
			expectedRows:  2,
			validateOutput: func(t *testing.T, result string) {
				if !strings.Contains(result, "2 rows") {
					t.Errorf("Expected result to mention 2 rows, got: %s", result)
				}
				if !strings.Contains(result, "custom_headers.csv") {
					t.Errorf("Expected result to mention custom_headers.csv, got: %s", result)
				}
				if !strings.Contains(result, "artifact") {
					t.Errorf("Expected result to mention artifact, got: %s", result)
				}
			},
		},
		{
			name: "CSV without headers",
			args: map[string]any{
				"data": []any{
					map[string]any{"name": "Alice", "age": 30},
					map[string]any{"name": "Bob", "age": 25},
				},
				"filename":        "no_headers.csv",
				"include_headers": false,
			},
			expectedError: false,
			expectedRows:  2,
			validateOutput: func(t *testing.T, result string) {
				if !strings.Contains(result, "2 rows") {
					t.Errorf("Expected result to mention 2 rows, got: %s", result)
				}
				if !strings.Contains(result, "no_headers.csv") {
					t.Errorf("Expected result to mention no_headers.csv, got: %s", result)
				}
				if !strings.Contains(result, "artifact") {
					t.Errorf("Expected result to mention artifact, got: %s", result)
				}
			},
		},
		{
			name: "append to existing file",
			args: map[string]any{
				"data": []any{
					map[string]any{"name": "Charlie", "age": 35},
				},
				"filename": "basic.csv",
				"append":   true,
			},
			expectedError: false,
			expectedRows:  1,
			validateOutput: func(t *testing.T, result string) {
				if !strings.Contains(result, "1 rows") {
					t.Errorf("Expected result to mention 1 rows, got: %s", result)
				}
				if !strings.Contains(result, "basic.csv") {
					t.Errorf("Expected result to mention basic.csv, got: %s", result)
				}
				if !strings.Contains(result, "artifact") {
					t.Errorf("Expected result to mention artifact, got: %s", result)
				}
			},
		},
		{
			name: "invalid data type",
			args: map[string]any{
				"data":     "not an array",
				"filename": "invalid.csv",
			},
			expectedError: true,
		},
		{
			name: "empty file path",
			args: map[string]any{
				"data":     []any{map[string]any{"name": "Alice"}},
				"filename": "",
			},
			expectedError: true,
		},
		{
			name: "empty data array",
			args: map[string]any{
				"data":     []any{},
				"filename": "empty.csv",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, server.ArtifactHelperContextKey, server.NewArtifactHelper())
			ctx = context.WithValue(ctx, server.TaskContextKey, &types.Task{ID: "test-task-123"})
			
			result, err := skill.WriteToCsvHandler(ctx, tt.args)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !strings.Contains(result, "Successfully created CSV") {
				t.Errorf("Expected success message, got: %s", result)
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, result)
			}
		})
	}
}

func TestConvertDataToRows(t *testing.T) {
	logger := zap.NewNop()
	skill := &WriteToCsvSkill{logger: logger}

	tests := []struct {
		name          string
		input         []any
		expectedError bool
		expectedLen   int
	}{
		{
			name: "valid map[string]any data",
			input: []any{
				map[string]any{"name": "Alice", "age": 30},
				map[string]any{"name": "Bob", "age": 25},
			},
			expectedError: false,
			expectedLen:   2,
		},
		{
			name: "mixed map types",
			input: []any{
				map[string]any{"name": "Alice"},
				map[any]any{"name": "Bob", "age": 25},
			},
			expectedError: false,
			expectedLen:   2,
		},
		{
			name: "invalid data type",
			input: []any{
				"not a map",
				map[string]any{"name": "Alice"},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := skill.convertDataToRows(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d rows, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestValueToString(t *testing.T) {
	logger := zap.NewNop()
	skill := &WriteToCsvSkill{logger: logger}

	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"nil", nil, ""},
		{"array", []any{"a", "b", "c"}, "[%!v([]string=[a b c])]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skill.valueToString(tt.input)
			if tt.name != "array" && result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}

			if tt.name == "array" && result == "" {
				t.Error("Expected non-empty string for array")
			}
		})
	}
}
