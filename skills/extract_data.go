package skills

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	server "github.com/inference-gateway/adk/server"
	playwright "github.com/inference-gateway/playwright-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// ExtractDataSkill struct holds the skill with dependencies
type ExtractDataSkill struct {
	logger     *zap.Logger
	playwright playwright.BrowserAutomation
}

// NewExtractDataSkill creates a new extract_data skill
func NewExtractDataSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	skill := &ExtractDataSkill{
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
					"items":       map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string", "description": "Name for the extracted data field"}, "selector": map[string]any{"type": "string", "description": "CSS selector or XPath to extract data from"}, "attribute": map[string]any{"type": "string", "description": "Attribute to extract (text, href, src, etc.)", "default": "text"}, "multiple": map[string]any{"type": "boolean", "description": "Extract all matching elements or just the first", "default": false}}, "required": []string{"name", "selector"}},
					"type":        "array",
				},
				"format": map[string]any{
					"default":     "json",
					"description": "Output format (json, csv, text)",
					"type":        "string",
				},
			},
			"required": []string{"extractors"},
		},
		skill.ExtractDataHandler,
	)
}

// ExtractDataHandler handles the extract_data skill execution
func (s *ExtractDataSkill) ExtractDataHandler(ctx context.Context, args map[string]any) (string, error) {
	extractors, ok := args["extractors"].([]any)
	if !ok || len(extractors) == 0 {
		s.logger.Error("extractors parameter is required and must be a non-empty array")
		return "", fmt.Errorf("extractors parameter is required and must be a non-empty array")
	}

	format := "json"
	if f, ok := args["format"].(string); ok && f != "" {
		if !s.isValidFormat(f) {
			return "", fmt.Errorf("invalid format: %s. Must be one of: json, csv, text", f)
		}
		format = f
	}

	s.logger.Info("extracting data from page",
		zap.Int("extractors_count", len(extractors)),
		zap.String("format", format))

	session, err := s.getOrCreateSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	playwrightExtractors, err := s.convertExtractors(extractors)
	if err != nil {
		s.logger.Error("failed to convert extractors", zap.Error(err))
		return "", fmt.Errorf("failed to convert extractors: %w", err)
	}

	rawResult, err := s.playwright.ExtractData(ctx, session.ID, playwrightExtractors, format)
	if err != nil {
		s.logger.Error("data extraction failed",
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("data extraction failed: %w", err)
	}

	result, err := s.processExtractedData(rawResult, format, extractors)
	if err != nil {
		s.logger.Error("failed to process extracted data", zap.Error(err))
		return "", fmt.Errorf("failed to process extracted data: %w", err)
	}

	s.logger.Info("data extraction completed successfully",
		zap.String("sessionID", session.ID),
		zap.String("format", format))

	return result, nil
}

// isValidFormat validates the output format parameter
func (s *ExtractDataSkill) isValidFormat(format string) bool {
	validFormats := []string{"json", "csv", "text"}
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	return false
}

// convertExtractors converts extractors from any to the format expected by Playwright service
func (s *ExtractDataSkill) convertExtractors(extractors []any) ([]map[string]any, error) {
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

// processExtractedData processes the raw extracted data and formats it according to the specified format
func (s *ExtractDataSkill) processExtractedData(rawResult, format string, originalExtractors []any) (string, error) {
	switch format {
	case "json":
		return s.formatAsJSON(rawResult, originalExtractors)
	case "csv":
		return s.formatAsCSV(rawResult, originalExtractors)
	case "text":
		return s.formatAsText(rawResult, originalExtractors)
	default:
		return rawResult, nil
	}
}

// formatAsJSON formats the extracted data as JSON with metadata
func (s *ExtractDataSkill) formatAsJSON(rawResult string, extractors []any) (string, error) {
	parsedData, err := s.parseRawResult(rawResult)
	if err != nil {
		return "", err
	}

	result := map[string]any{
		"success":    true,
		"format":     "json",
		"extractors": len(extractors),
		"data":       parsedData,
		"metadata": map[string]any{
			"extraction_time": time.Now().Unix(),
			"total_fields":    len(parsedData),
		},
	}

	cleanedData := s.cleanAndNormalizeData(result["data"])
	result["data"] = cleanedData

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// formatAsCSV formats the extracted data as CSV
func (s *ExtractDataSkill) formatAsCSV(rawResult string, extractors []any) (string, error) {
	parsedData, err := s.parseRawResult(rawResult)
	if err != nil {
		return "", err
	}

	var csvBuilder strings.Builder
	writer := csv.NewWriter(&csvBuilder)

	headers := make([]string, 0, len(parsedData))
	for key := range parsedData {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	rows := s.generateCSVRows(parsedData, headers)
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writing error: %w", err)
	}

	return csvBuilder.String(), nil
}

// formatAsText formats the extracted data as human-readable text
func (s *ExtractDataSkill) formatAsText(rawResult string, extractors []any) (string, error) {
	parsedData, err := s.parseRawResult(rawResult)
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder
	textBuilder.WriteString("Extracted Data:\n")
	textBuilder.WriteString("==============\n\n")

	for key, value := range parsedData {
		textBuilder.WriteString(fmt.Sprintf("%s: ", key))
		switch v := value.(type) {
		case []any:
			textBuilder.WriteString("\n")
			for i, item := range v {
				textBuilder.WriteString(fmt.Sprintf("  [%d] %v\n", i+1, item))
			}
		default:
			textBuilder.WriteString(fmt.Sprintf("%v\n", v))
		}
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}

// parseRawResult parses the raw result string from Playwright service
func (s *ExtractDataSkill) parseRawResult(rawResult string) (map[string]any, error) {

	var parsedData map[string]any

	if err := json.Unmarshal([]byte(rawResult), &parsedData); err == nil {
		return parsedData, nil
	}

	if strings.HasPrefix(rawResult, "map[") && strings.HasSuffix(rawResult, "]") {
		return s.parseGoMapFormat(rawResult)
	}

	return s.extractDataFromString(rawResult)
}

// parseGoMapFormat parses Go map format: map[key1:value1 key2:value2]
func (s *ExtractDataSkill) parseGoMapFormat(mapStr string) (map[string]any, error) {
	data := make(map[string]any)

	content := strings.TrimPrefix(mapStr, "map[")
	content = strings.TrimSuffix(content, "]")

	if content == "" {
		return data, nil
	}

	parts := s.smartSplit(content)

	for _, part := range parts {
		if strings.Contains(part, ":") {
			keyValue := strings.SplitN(part, ":", 2)
			if len(keyValue) == 2 {
				key := strings.TrimSpace(keyValue[0])
				value := strings.TrimSpace(keyValue[1])

				data[key] = s.parseValue(value)
			}
		}
	}

	return data, nil
}

// smartSplit splits the map content intelligently, handling quoted values and brackets
func (s *ExtractDataSkill) smartSplit(content string) []string {
	var parts []string
	var current strings.Builder
	inBrackets := 0
	inQuotes := false
	foundKey := false

	for i, char := range content {
		switch char {
		case '[':
			inBrackets++
			current.WriteRune(char)
		case ']':
			inBrackets--
			current.WriteRune(char)
		case '"', '\'':
			inQuotes = !inQuotes
			current.WriteRune(char)
		case ':':
			if inBrackets == 0 && !inQuotes && !foundKey {
				foundKey = true
				current.WriteRune(char)
			} else {
				current.WriteRune(char)
			}
		case ' ':
			if inBrackets == 0 && !inQuotes && foundKey {
				if s.isNextKeyStart(content, i) {
					if current.Len() > 0 {
						parts = append(parts, current.String())
						current.Reset()
						foundKey = false
					}
				} else {
					current.WriteRune(char)
				}
			} else if inBrackets == 0 && !inQuotes && !foundKey {
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// isNextKeyStart checks if the next non-space characters form a key (word followed by colon)
func (s *ExtractDataSkill) isNextKeyStart(content string, currentPos int) bool {
	i := currentPos + 1
	for i < len(content) && content[i] == ' ' {
		i++
	}

	if i >= len(content) {
		return false
	}

	wordStart := i
	for i < len(content) && (content[i] != ' ' && content[i] != ':' && content[i] != '[' && content[i] != ']') {
		i++
	}

	return i > wordStart && i < len(content) && content[i] == ':'
}

// parseValue attempts to parse a string value into appropriate Go type
func (s *ExtractDataSkill) parseValue(value string) any {
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		arrayContent := strings.TrimPrefix(value, "[")
		arrayContent = strings.TrimSuffix(arrayContent, "]")

		if arrayContent == "" {
			return []any{}
		}

		items := strings.Fields(arrayContent)
		result := make([]any, len(items))
		for i, item := range items {
			result[i] = s.parseScalarValue(item)
		}
		return result
	}

	return s.parseScalarValue(value)
}

// parseScalarValue parses individual scalar values
func (s *ExtractDataSkill) parseScalarValue(value string) any {
	if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
		return value[1 : len(value)-1]
	}

	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}

	if value == "<nil>" || value == "null" {
		return nil
	}

	return value
}

// extractDataFromString extracts data from string representation (fallback method)
func (s *ExtractDataSkill) extractDataFromString(result string) (map[string]any, error) {
	data := make(map[string]any)

	pattern := regexp.MustCompile(`(\w+):\s*([^\n]+)`)
	matches := pattern.FindAllStringSubmatch(result, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])
			data[key] = value
		}
	}

	return data, nil
}

// generateCSVRows generates CSV rows from parsed data, handling arrays appropriately
func (s *ExtractDataSkill) generateCSVRows(data map[string]any, headers []string) [][]string {
	maxRows := 1
	for _, value := range data {
		if arr, ok := value.([]any); ok {
			if len(arr) > maxRows {
				maxRows = len(arr)
			}
		}
	}

	rows := make([][]string, maxRows)
	for i := 0; i < maxRows; i++ {
		rows[i] = make([]string, len(headers))
		for j, header := range headers {
			value := data[header]
			if arr, ok := value.([]any); ok {
				if i < len(arr) {
					rows[i][j] = fmt.Sprintf("%v", arr[i])
				} else {
					rows[i][j] = ""
				}
			} else {
				if i == 0 {
					rows[i][j] = fmt.Sprintf("%v", value)
				} else {
					rows[i][j] = ""
				}
			}
		}
	}

	return rows
}

// cleanAndNormalizeData applies data cleaning and normalization
func (s *ExtractDataSkill) cleanAndNormalizeData(data any) any {
	switch v := data.(type) {
	case map[string]any:
		cleaned := make(map[string]any)
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
		return s.cleanString(v)
	default:
		return v
	}
}

// cleanString performs string cleaning and normalization
func (s *ExtractDataSkill) cleanString(text string) string {
	cleaned := strings.TrimSpace(text)

	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

	controlRegex := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	cleaned = controlRegex.ReplaceAllString(cleaned, "")

	return cleaned
}

// getOrCreateSession gets the shared default session
func (s *ExtractDataSkill) getOrCreateSession(ctx context.Context) (*playwright.BrowserSession, error) {
	return s.playwright.GetOrCreateDefaultSession(ctx)
}
