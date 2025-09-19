package skills

import (
	"context"
	"testing"

	"github.com/inference-gateway/playwright-agent/internal/playwright"
	"github.com/inference-gateway/playwright-agent/internal/playwright/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestExtractDataHandler(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	skill := &ExtractDataSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	tests := []struct {
		name        string
		args        map[string]any
		mockSetup   func()
		expectedErr bool
		contains    string
	}{
		{
			name: "successful extraction with single extractor",
			args: map[string]any{
				"extractors": []interface{}{
					map[string]interface{}{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "json",
			},
			mockSetup: func() {
				mockPlaywright.LaunchBrowserReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				mockPlaywright.ExtractDataReturns(`map[title:Test Title]`, nil)
			},
			expectedErr: false,
			contains:    "Test Title",
		},
		{
			name: "successful extraction with multiple extractors",
			args: map[string]any{
				"extractors": []interface{}{
					map[string]interface{}{
						"name":      "title",
						"selector":  "h1",
						"attribute": "text",
					},
					map[string]interface{}{
						"name":      "links",
						"selector":  "a",
						"attribute": "href",
						"multiple":  true,
					},
				},
				"format": "json",
			},
			mockSetup: func() {
				mockPlaywright.LaunchBrowserReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				mockPlaywright.ExtractDataReturns(`map[title:Test Title links:[/page1 /page2]]`, nil)
			},
			expectedErr: false,
			contains:    "Test Title",
		},
		{
			name: "extraction with CSV format",
			args: map[string]any{
				"extractors": []interface{}{
					map[string]interface{}{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "csv",
			},
			mockSetup: func() {
				mockPlaywright.LaunchBrowserReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				mockPlaywright.ExtractDataReturns(`map[title:Test Title]`, nil)
			},
			expectedErr: false,
			contains:    "title",
		},
		{
			name: "extraction with text format",
			args: map[string]any{
				"extractors": []interface{}{
					map[string]interface{}{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "text",
			},
			mockSetup: func() {
				mockPlaywright.LaunchBrowserReturns(&playwright.BrowserSession{ID: "test-session"}, nil)
				mockPlaywright.ExtractDataReturns(`map[title:Test Title]`, nil)
			},
			expectedErr: false,
			contains:    "Extracted Data",
		},
		{
			name: "missing extractors parameter",
			args: map[string]any{
				"format": "json",
			},
			mockSetup:   func() {},
			expectedErr: true,
			contains:    "extractors parameter is required",
		},
		{
			name: "empty extractors array",
			args: map[string]any{
				"extractors": []interface{}{},
				"format":     "json",
			},
			mockSetup:   func() {},
			expectedErr: true,
			contains:    "extractors parameter is required",
		},
		{
			name: "invalid format",
			args: map[string]any{
				"extractors": []interface{}{
					map[string]interface{}{
						"name":     "title",
						"selector": "h1",
					},
				},
				"format": "xml",
			},
			mockSetup:   func() {},
			expectedErr: true,
			contains:    "invalid format",
		},
		{
			name: "extractor missing name",
			args: map[string]any{
				"extractors": []interface{}{
					map[string]interface{}{
						"selector": "h1",
					},
				},
				"format": "json",
			},
			mockSetup:   func() {},
			expectedErr: true,
			contains:    "must have a non-empty 'name' field",
		},
		{
			name: "extractor missing selector",
			args: map[string]any{
				"extractors": []interface{}{
					map[string]interface{}{
						"name": "title",
					},
				},
				"format": "json",
			},
			mockSetup:   func() {},
			expectedErr: true,
			contains:    "must have a non-empty 'selector' field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPlaywright = &mocks.FakeBrowserAutomation{}
			skill.playwright = mockPlaywright

			tt.mockSetup()

			result, err := skill.ExtractDataHandler(context.Background(), tt.args)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.contains)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, result, tt.contains)
			}
		})
	}
}

func TestConvertExtractors(t *testing.T) {
	skill := &ExtractDataSkill{}

	tests := []struct {
		name        string
		extractors  []interface{}
		expected    []map[string]interface{}
		expectedErr bool
	}{
		{
			name: "single extractor with defaults",
			extractors: []interface{}{
				map[string]interface{}{
					"name":     "title",
					"selector": "h1",
				},
			},
			expected: []map[string]interface{}{
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
			extractors: []interface{}{
				map[string]interface{}{
					"name":      "links",
					"selector":  "a",
					"attribute": "href",
					"multiple":  true,
				},
			},
			expected: []map[string]interface{}{
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
			extractors: []interface{}{
				"invalid",
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := skill.convertExtractors(tt.extractors)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseGoMapFormat(t *testing.T) {
	skill := &ExtractDataSkill{}

	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:  "simple map",
			input: "map[title:Test Title count:42]",
			expected: map[string]interface{}{
				"title": "Test Title",
				"count": 42,
			},
		},
		{
			name:  "map with array",
			input: "map[links:[/page1 /page2 /page3]]",
			expected: map[string]interface{}{
				"links": []interface{}{"/page1", "/page2", "/page3"},
			},
		},
		{
			name:     "empty map",
			input:    "map[]",
			expected: map[string]interface{}{},
		},
		{
			name:  "map with boolean and nil",
			input: "map[active:true missing:<nil>]",
			expected: map[string]interface{}{
				"active":  true,
				"missing": nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := skill.parseGoMapFormat(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanString(t *testing.T) {
	skill := &ExtractDataSkill{}

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
			result := skill.cleanString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidFormat(t *testing.T) {
	skill := &ExtractDataSkill{}

	tests := []struct {
		format   string
		expected bool
	}{
		{"json", true},
		{"csv", true},
		{"text", true},
		{"xml", false},
		{"", false},
		{"JSON", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			result := skill.isValidFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}
