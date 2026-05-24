package tools

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

var (
	validExtractFormats = []string{"json", "csv", "text"}

	// Hoisted regexes — previously compiled per call inside cleanString.
	whitespaceRun = regexp.MustCompile(`\s+`)
	controlChars  = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
)

// ExtractDataTool struct holds the tool with dependencies
type ExtractDataTool struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewExtractDataTool creates a new extract_data tool
func NewExtractDataTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	tool := &ExtractDataTool{
		logger:     logger,
		playwright: playwright,
	}
	return server.NewBasicTool(
		"extract_data",
		"Extract data from the page using selectors and return structured information",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"extractors": map[string]any{
					"description": "List of data extractors to run",
					"type":        "array",
					"items": map[string]any{
						"required": []string{"name", "selector"},
						"type":     "object",
						"properties": map[string]any{
							"name":      map[string]any{"type": "string", "description": "Name for the extracted data field"},
							"selector":  map[string]any{"type": "string", "description": "CSS selector or XPath to extract data from"},
							"attribute": map[string]any{"type": "string", "description": "Attribute to extract (text, href, src, etc.)", "default": "text"},
							"multiple":  map[string]any{"type": "boolean", "description": "Extract all matching elements or just the first", "default": false},
						},
					},
				},
				"format": map[string]any{
					"default":     "json",
					"description": "Output format (json, csv, text)",
					"type":        "string",
				},
			},
			"required": []string{"extractors"},
		},
		tool.ExtractDataHandler,
	)
}

// ExtractDataHandler handles the extract_data tool execution.
//
// The playwright service returns canonical JSON for the extracted map;
// previously this tool carried a hand-rolled "Go %+v map" parser as a
// fallback because the service emitted fmt.Sprintf("%+v", results). That
// fallback is now gone.
func (s *ExtractDataTool) ExtractDataHandler(ctx context.Context, args map[string]any) (string, error) {
	rawExtractors, present, err := sliceArg(args, "extractors")
	if err != nil {
		return "", err
	}
	if !present || len(rawExtractors) == 0 {
		return "", fmt.Errorf("extractors parameter is required and must be a non-empty array")
	}

	format, err := stringArg(args, "format", "json")
	if err != nil {
		return "", err
	}
	if !oneOf(format, validExtractFormats...) {
		return "", fmt.Errorf("invalid format: %s. Must be one of: %v", format, validExtractFormats)
	}

	s.logger.Info("extracting data from page",
		zap.Int("extractors_count", len(rawExtractors)),
		zap.String("format", format))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	playwrightExtractors, err := s.convertExtractors(rawExtractors)
	if err != nil {
		s.logger.Error("failed to convert extractors", zap.Error(err))
		return "", fmt.Errorf("failed to convert extractors: %w", err)
	}

	rawResult, err := s.playwright.ExtractData(ctx, session.ID, playwrightExtractors, format)
	if err != nil {
		s.logger.Error("data extraction failed",
			zap.String("sessionID", session.ID),
			zap.Error(err))
		if strings.Contains(err.Error(), "strict mode violation") {
			return "", fmt.Errorf("data extraction failed: %w; the selector matched multiple elements - pass \"multiple\": true to extract all of them, or use a more specific selector to target a single element", err)
		}
		return "", fmt.Errorf("data extraction failed: %w", err)
	}

	parsed, err := s.parseRawResult(rawResult)
	if err != nil {
		s.logger.Error("failed to parse extracted data", zap.Error(err))
		return "", fmt.Errorf("failed to parse extracted data: %w", err)
	}

	cleaned, _ := s.cleanAndNormalizeData(parsed).(map[string]any)
	if cleaned == nil {
		cleaned = parsed
	}

	switch format {
	case "csv":
		return s.formatAsCSV(cleaned)
	case "text":
		return s.formatAsText(cleaned), nil
	default:
		return s.formatAsJSON(cleaned, len(rawExtractors))
	}
}

// parseRawResult parses the canonical JSON document returned by the
// playwright service.
func (s *ExtractDataTool) parseRawResult(rawResult string) (map[string]any, error) {
	var parsed map[string]any
	if err := json.Unmarshal([]byte(rawResult), &parsed); err != nil {
		return nil, fmt.Errorf("playwright service returned non-JSON payload: %w", err)
	}
	return parsed, nil
}

// convertExtractors converts extractors from any to the format expected by Playwright service
func (s *ExtractDataTool) convertExtractors(extractors []any) ([]map[string]any, error) {
	converted := make([]map[string]any, len(extractors))

	for i, extractor := range extractors {
		extractorMap, ok := extractor.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("extractor at index %d must be an object", i)
		}

		name, ok := extractorMap["name"].(string)
		if !ok || name == "" {
			return nil, fmt.Errorf("extractor at index %d must have a non-empty 'name' field", i)
		}

		selector, ok := extractorMap["selector"].(string)
		if !ok || selector == "" {
			return nil, fmt.Errorf("extractor at index %d must have a non-empty 'selector' field", i)
		}

		attribute := "text"
		if attr, ok := extractorMap["attribute"].(string); ok && attr != "" {
			attribute = attr
		}

		multiple := false
		if mult, ok := extractorMap["multiple"].(bool); ok {
			multiple = mult
		}

		converted[i] = map[string]any{
			"name":      name,
			"selector":  selector,
			"attribute": attribute,
			"multiple":  multiple,
		}
	}

	return converted, nil
}

// formatAsJSON wraps the extracted data in the canonical envelope with
// metadata. Returns valid JSON.
func (s *ExtractDataTool) formatAsJSON(data map[string]any, extractorCount int) (string, error) {
	return marshalResponse(map[string]any{
		"success":    true,
		"format":     "json",
		"extractors": extractorCount,
		"data":       data,
		"metadata": map[string]any{
			"extraction_time": time.Now().Unix(),
			"total_fields":    len(data),
		},
	})
}

// formatAsCSV writes a header row of extractor names and as many data rows
// as the largest array-valued field. Scalar fields appear only on the first
// row.
func (s *ExtractDataTool) formatAsCSV(data map[string]any) (string, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	headers := make([]string, 0, len(data))
	for key := range data {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, row := range generateCSVRows(data, headers) {
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writing error: %w", err)
	}

	return buf.String(), nil
}

// formatAsText emits a human-readable rendering of the extracted data.
func (s *ExtractDataTool) formatAsText(data map[string]any) string {
	var buf strings.Builder
	buf.WriteString("Extracted Data:\n")
	buf.WriteString("==============\n\n")

	for key, value := range data {
		fmt.Fprintf(&buf, "%s: ", key)
		switch v := value.(type) {
		case []any:
			buf.WriteString("\n")
			for i, item := range v {
				fmt.Fprintf(&buf, "  [%d] %v\n", i+1, item)
			}
		default:
			fmt.Fprintf(&buf, "%v\n", v)
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// generateCSVRows generates CSV rows from parsed data. The number of rows
// equals the longest array-valued field; scalar fields show only on row 0.
func generateCSVRows(data map[string]any, headers []string) [][]string {
	maxRows := 1
	for _, value := range data {
		if arr, ok := value.([]any); ok && len(arr) > maxRows {
			maxRows = len(arr)
		}
	}

	rows := make([][]string, maxRows)
	for i := range maxRows {
		rows[i] = make([]string, len(headers))
		for j, header := range headers {
			switch v := data[header].(type) {
			case []any:
				if i < len(v) {
					rows[i][j] = fmt.Sprintf("%v", v[i])
				}
			default:
				if i == 0 {
					rows[i][j] = fmt.Sprintf("%v", v)
				}
			}
		}
	}
	return rows
}

// cleanAndNormalizeData walks the JSON tree and applies cleanString to
// every string leaf.
func (s *ExtractDataTool) cleanAndNormalizeData(data any) any {
	switch v := data.(type) {
	case map[string]any:
		cleaned := make(map[string]any, len(v))
		for key, value := range v {
			cleaned[key] = s.cleanAndNormalizeData(value)
		}
		return cleaned
	case []any:
		cleaned := make([]any, len(v))
		for i, value := range v {
			cleaned[i] = s.cleanAndNormalizeData(value)
		}
		return cleaned
	case string:
		return cleanString(v)
	default:
		return v
	}
}

// cleanString collapses whitespace runs and strips control characters.
// Regexes are package-level so they compile once.
func cleanString(text string) string {
	cleaned := strings.TrimSpace(text)
	cleaned = whitespaceRun.ReplaceAllString(cleaned, " ")
	cleaned = controlChars.ReplaceAllString(cleaned, "")
	return cleaned
}
