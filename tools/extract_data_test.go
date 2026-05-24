package tools

import (
	"context"
	"encoding/json"
	"testing"

	zap "go.uber.org/zap"

	assert "github.com/stretchr/testify/assert"

	mocks "github.com/inference-gateway/browser-agent/internal/playwright/mocks"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

func TestExtractDataHandler(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name        string
		args        map[string]any
		mockSetup   func(*mocks.FakeBrowserAutomation)
		expectedErr bool
		contains    string
	}{
		{
			name: "successful extraction with single extractor",
			args: map[string]any{
				"extractors": []any{
					map[string]any{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "json",
			},
			mockSetup: func(m *mocks.FakeBrowserAutomation) {
				m.GetOrCreateTaskSessionReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				m.ExtractDataReturns(`{"title":"Test Title"}`, nil)
			},
			expectedErr: false,
			contains:    "Test Title",
		},
		{
			name: "successful extraction with multiple extractors",
			args: map[string]any{
				"extractors": []any{
					map[string]any{
						"name":      "title",
						"selector":  "h1",
						"attribute": "text",
					},
					map[string]any{
						"name":      "links",
						"selector":  "a",
						"attribute": "href",
						"multiple":  true,
					},
				},
				"format": "json",
			},
			mockSetup: func(m *mocks.FakeBrowserAutomation) {
				m.GetOrCreateTaskSessionReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				m.ExtractDataReturns(`{"title":"Test Title","links":["/page1","/page2"]}`, nil)
			},
			expectedErr: false,
			contains:    "Test Title",
		},
		{
			name: "extraction with CSV format",
			args: map[string]any{
				"extractors": []any{
					map[string]any{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "csv",
			},
			mockSetup: func(m *mocks.FakeBrowserAutomation) {
				m.GetOrCreateTaskSessionReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				m.ExtractDataReturns(`{"title":"Test Title"}`, nil)
			},
			expectedErr: false,
			contains:    "title",
		},
		{
			name: "extraction with text format",
			args: map[string]any{
				"extractors": []any{
					map[string]any{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "text",
			},
			mockSetup: func(m *mocks.FakeBrowserAutomation) {
				m.GetOrCreateTaskSessionReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				m.ExtractDataReturns(`{"title":"Test Title"}`, nil)
			},
			expectedErr: false,
			contains:    "Extracted Data",
		},
		{
			name: "missing extractors parameter",
			args: map[string]any{
				"format": "json",
			},
			mockSetup:   func(m *mocks.FakeBrowserAutomation) {},
			expectedErr: true,
			contains:    "extractors parameter is required",
		},
		{
			name: "empty extractors array",
			args: map[string]any{
				"extractors": []any{},
				"format":     "json",
			},
			mockSetup:   func(m *mocks.FakeBrowserAutomation) {},
			expectedErr: true,
			contains:    "extractors parameter is required",
		},
		{
			name: "invalid format",
			args: map[string]any{
				"extractors": []any{
					map[string]any{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "xml",
			},
			mockSetup:   func(m *mocks.FakeBrowserAutomation) {},
			expectedErr: true,
			contains:    "invalid format",
		},
		{
			name: "extractor missing name",
			args: map[string]any{
				"extractors": []any{
					map[string]any{
						"selector": "h1",
					},
				},
				"format": "json",
			},
			mockSetup:   func(m *mocks.FakeBrowserAutomation) {},
			expectedErr: true,
			contains:    "must have a non-empty 'name' field",
		},
		{
			name: "extractor missing selector",
			args: map[string]any{
				"extractors": []any{
					map[string]any{
						"name": "title",
					},
				},
				"format": "json",
			},
			mockSetup:   func(m *mocks.FakeBrowserAutomation) {},
			expectedErr: true,
			contains:    "must have a non-empty 'selector' field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlaywright := &mocks.FakeBrowserAutomation{}
			tt.mockSetup(mockPlaywright)

			tool := &ExtractDataTool{logger: logger, playwright: mockPlaywright}
			result, err := tool.ExtractDataHandler(context.Background(), tt.args)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.contains)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, result, tt.contains)
		})
	}
}

// TestExtractDataHandler_JSONFormat_ProducesValidJSON checks that the
// json output is itself parseable — guards against future regressions
// to the broken `%+v` envelope.
func TestExtractDataHandler_JSONFormat_ProducesValidJSON(t *testing.T) {
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	mockPlaywright.GetOrCreateTaskSessionReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
	mockPlaywright.ExtractDataReturns(`{"title":"  Spaced  Out  ","count":42}`, nil)

	tool := &ExtractDataTool{logger: zap.NewNop(), playwright: mockPlaywright}
	result, err := tool.ExtractDataHandler(context.Background(), map[string]any{
		"extractors": []any{map[string]any{"name": "title", "selector": "h1"}},
		"format":     "json",
	})
	assert.NoError(t, err)

	var parsed map[string]any
	assert.NoError(t, json.Unmarshal([]byte(result), &parsed))
	assert.Equal(t, true, parsed["success"])
	assert.Equal(t, "json", parsed["format"])

	data, _ := parsed["data"].(map[string]any)
	assert.Equal(t, "Spaced Out", data["title"], "whitespace runs should collapse via cleanString")
	assert.Equal(t, float64(42), data["count"])
}

func TestConvertExtractors(t *testing.T) {
	tool := &ExtractDataTool{}

	tests := []struct {
		name        string
		extractors  []any
		expected    []map[string]any
		expectedErr bool
	}{
		{
			name: "single extractor with defaults",
			extractors: []any{
				map[string]any{
					"name":     "title",
					"selector": "h1",
				},
			},
			expected: []map[string]any{
				{
					"name":      "title",
					"selector":  "h1",
					"attribute": "text",
					"multiple":  false,
				},
			},
			expectedErr: false,
		},
		{
			name: "extractor with all fields",
			extractors: []any{
				map[string]any{
					"name":      "links",
					"selector":  "a",
					"attribute": "href",
					"multiple":  true,
				},
			},
			expected: []map[string]any{
				{
					"name":      "links",
					"selector":  "a",
					"attribute": "href",
					"multiple":  true,
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid extractor type",
			extractors: []any{
				"invalid",
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.convertExtractors(tt.extractors)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCleanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic trimming",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "multiple spaces",
			input:    "hello    world",
			expected: "hello world",
		},
		{
			name:     "mixed whitespace",
			input:    "hello\t\n  world",
			expected: "hello world",
		},
		{
			name:     "control characters",
			input:    "hello\x00\x01world",
			expected: "helloworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, cleanString(tt.input))
		})
	}
}
