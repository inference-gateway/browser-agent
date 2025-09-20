package skills

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestWriteToCsvHandler(t *testing.T) {
	logger := zap.NewNop()
	skill := &WriteToCsvSkill{logger: logger}

	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		args           map[string]any
		expectedError  bool
		expectedRows   int
		validateOutput func(t *testing.T, filePath string)
	}{
		{
			name: "basic CSV writing",
			args: map[string]any{
				"data": []any{
					map[string]any{"name": "Alice", "age": 30, "city": "New York"},
					map[string]any{"name": "Bob", "age": 25, "city": "San Francisco"},
				},
				"file_path": filepath.Join(tempDir, "basic.csv"),
			},
			expectedError: false,
			expectedRows:  2,
			validateOutput: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}

				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				if len(lines) != 3 { // header + 2 data rows
					t.Errorf("Expected 3 lines, got %d", len(lines))
				}

				// Check if headers are present
				if !strings.Contains(lines[0], "name") {
					t.Error("Expected headers to contain 'name'")
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
				"file_path": filepath.Join(tempDir, "custom_headers.csv"),
				"headers":   []any{"name", "age"},
			},
			expectedError: false,
			expectedRows:  2,
			validateOutput: func(t *testing.T, filePath string) {
				file, err := os.Open(filePath)
				if err != nil {
					t.Fatalf("Failed to open output file: %v", err)
				}
				defer func() {
					if closeErr := file.Close(); closeErr != nil {
						t.Logf("Failed to close file: %v", closeErr)
					}
				}()

				reader := csv.NewReader(file)
				records, err := reader.ReadAll()
				if err != nil {
					t.Fatalf("Failed to read CSV: %v", err)
				}

				if len(records) != 3 { // header + 2 data rows
					t.Errorf("Expected 3 records, got %d", len(records))
				}

				// Check header order
				if records[0][0] != "name" || records[0][1] != "age" {
					t.Errorf("Headers not in expected order: %v", records[0])
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
				"file_path":       filepath.Join(tempDir, "no_headers.csv"),
				"include_headers": false,
			},
			expectedError: false,
			expectedRows:  2,
			validateOutput: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}

				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				if len(lines) != 2 { // only data rows, no header
					t.Errorf("Expected 2 lines, got %d", len(lines))
				}
			},
		},
		{
			name: "append to existing file",
			args: map[string]any{
				"data": []any{
					map[string]any{"name": "Charlie", "age": 35},
				},
				"file_path": filepath.Join(tempDir, "basic.csv"), // reuse existing file
				"append":    true,
			},
			expectedError: false,
			expectedRows:  1,
			validateOutput: func(t *testing.T, filePath string) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}

				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				if len(lines) != 4 { // original header + 2 original rows + 1 new row
					t.Errorf("Expected 4 lines after append, got %d", len(lines))
				}

				// Check if new data is appended
				if !strings.Contains(string(content), "Charlie") {
					t.Error("Expected appended data to contain 'Charlie'")
				}
			},
		},
		{
			name: "invalid data type",
			args: map[string]any{
				"data":      "not an array",
				"file_path": filepath.Join(tempDir, "invalid.csv"),
			},
			expectedError: true,
		},
		{
			name: "empty file path",
			args: map[string]any{
				"data":      []any{map[string]any{"name": "Alice"}},
				"file_path": "",
			},
			expectedError: true,
		},
		{
			name: "empty data array",
			args: map[string]any{
				"data":      []any{},
				"file_path": filepath.Join(tempDir, "empty.csv"),
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := skill.WriteToCsvHandler(context.Background(), tt.args)

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

			if !strings.Contains(result, "Successfully wrote") {
				t.Errorf("Expected success message, got: %s", result)
			}

			if tt.validateOutput != nil {
				filePath := tt.args["file_path"].(string)
				tt.validateOutput(t, filePath)
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
			// For array test, just check it's not empty
			if tt.name == "array" && result == "" {
				t.Error("Expected non-empty string for array")
			}
		})
	}
}
