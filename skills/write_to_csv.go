package skills

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	server "github.com/inference-gateway/adk/server"
	"github.com/inference-gateway/browser-agent/config"
	zap "go.uber.org/zap"
)

type WriteToCsvSkill struct {
	logger       *zap.Logger
	dataFilesDir string
}

func NewWriteToCsvSkill(logger *zap.Logger, cfg *config.Config) server.Tool {
	skill := &WriteToCsvSkill{
		logger:       logger,
		dataFilesDir: cfg.Browser.DataDir,
	}
	return server.NewBasicTool(
		"write_to_csv",
		"Write structured data to CSV files with support for custom headers and file paths",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"data": map[string]any{
					"description": "Array of objects to write to CSV, each object represents a row",
					"items":       map[string]any{"type": "object"},
					"type":        "array",
				},
				"filename": map[string]any{
					"description": "Name of the CSV file (without path, will be saved to configured data directory)",
					"type":        "string",
				},
				"headers": map[string]any{
					"description": "Custom column headers for the CSV file (optional, will use object keys if not provided)",
					"items":       map[string]any{"type": "string"},
					"type":        "array",
				},
				"append": map[string]any{
					"default":     false,
					"description": "Whether to append to existing file or create new file",
					"type":        "boolean",
				},
				"include_headers": map[string]any{
					"default":     true,
					"description": "Whether to include headers in the CSV output",
					"type":        "boolean",
				},
			},
			"required": []string{"data", "filename"},
		},
		skill.WriteToCsvHandler,
	)
}

// WriteToCsvHandler handles the write_to_csv skill execution
func (s *WriteToCsvSkill) WriteToCsvHandler(ctx context.Context, args map[string]any) (string, error) {
	data, ok := args["data"].([]any)
	if !ok || len(data) == 0 {
		s.logger.Error("data parameter is required and must be a non-empty array")
		return "", fmt.Errorf("data parameter is required and must be a non-empty array")
	}

	filename, ok := args["filename"].(string)
	if !ok || filename == "" {
		s.logger.Error("filename parameter is required and must be a non-empty string")
		return "", fmt.Errorf("filename parameter is required and must be a non-empty string")
	}

	filePath := s.generateFilePath(filename)

	var customHeaders []string
	if headers, ok := args["headers"].([]any); ok {
		customHeaders = make([]string, len(headers))
		for i, header := range headers {
			if headerStr, ok := header.(string); ok {
				customHeaders[i] = headerStr
			} else {
				return "", fmt.Errorf("all headers must be strings")
			}
		}
	}

	append := false
	if appendVal, ok := args["append"].(bool); ok {
		append = appendVal
	}

	includeHeaders := true
	if includeVal, ok := args["include_headers"].(bool); ok {
		includeHeaders = includeVal
	}

	s.logger.Info("writing data to CSV file",
		zap.String("filename", filename),
		zap.String("file_path", filePath),
		zap.Int("rows_count", len(data)),
		zap.Bool("append", append),
		zap.Bool("include_headers", includeHeaders))

	rows, err := s.convertDataToRows(data)
	if err != nil {
		s.logger.Error("failed to convert data to rows", zap.Error(err))
		return "", fmt.Errorf("failed to convert data to rows: %w", err)
	}

	headers := customHeaders
	if len(headers) == 0 && len(rows) > 0 {
		headers = s.extractHeadersFromRows(rows)
	}

	rowsWritten, err := s.writeCSVFile(filePath, headers, rows, append, includeHeaders)
	if err != nil {
		s.logger.Error("failed to write CSV file",
			zap.String("file_path", filePath),
			zap.Error(err))
		return "", fmt.Errorf("failed to write CSV file: %w", err)
	}

	result := fmt.Sprintf("Successfully wrote %d rows to %s", rowsWritten, filePath)
	s.logger.Info("CSV file written successfully",
		zap.String("file_path", filePath),
		zap.Int("rows_written", rowsWritten))

	return result, nil
}

func (s *WriteToCsvSkill) generateFilePath(filename string) string {
	if err := os.MkdirAll(s.dataFilesDir, 0755); err != nil {
		s.logger.Warn("failed to create data files directory", zap.String("dir", s.dataFilesDir), zap.Error(err))
	}

	if !filepath.IsAbs(filename) {
		return filepath.Join(s.dataFilesDir, filename)
	}
	return filename
}

func (s *WriteToCsvSkill) convertDataToRows(data []any) ([]map[string]any, error) {
	rows := make([]map[string]any, len(data))

	for i, item := range data {
		switch v := item.(type) {
		case map[string]any:
			rows[i] = v
		case map[any]any:
			converted := make(map[string]any)
			for key, value := range v {
				if keyStr, ok := key.(string); ok {
					converted[keyStr] = value
				} else {
					converted[fmt.Sprintf("%v", key)] = value
				}
			}
			rows[i] = converted
		default:
			return nil, fmt.Errorf("data item at index %d must be an object/map, got %T", i, item)
		}
	}

	return rows, nil
}

func (s *WriteToCsvSkill) extractHeadersFromRows(rows []map[string]any) []string {
	headerSet := make(map[string]bool)
	var headers []string

	for _, row := range rows {
		for key := range row {
			if !headerSet[key] {
				headerSet[key] = true
				headers = append(headers, key)
			}
		}
	}

	return headers
}

func (s *WriteToCsvSkill) writeCSVFile(filePath string, headers []string, rows []map[string]any, append bool, includeHeaders bool) (int, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	flag := os.O_CREATE | os.O_WRONLY
	if append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	fileExists := false
	if append {
		if info, err := os.Stat(filePath); err == nil && info.Size() > 0 {
			fileExists = true
		}
	}

	file, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			s.logger.Error("failed to close file", zap.String("file_path", filePath), zap.Error(closeErr))
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	rowsWritten := 0

	if includeHeaders && (!append || !fileExists) {
		if len(headers) > 0 {
			if err := writer.Write(headers); err != nil {
				return 0, fmt.Errorf("failed to write headers: %w", err)
			}
		}
	}

	for _, row := range rows {
		csvRow := make([]string, len(headers))
		for i, header := range headers {
			if value, exists := row[header]; exists {
				csvRow[i] = s.valueToString(value)
			} else {
				csvRow[i] = ""
			}
		}

		if err := writer.Write(csvRow); err != nil {
			return rowsWritten, fmt.Errorf("failed to write row: %w", err)
		}
		rowsWritten++
	}

	return rowsWritten, nil
}

func (s *WriteToCsvSkill) valueToString(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case []any:
		var items []string
		for _, item := range v {
			items = append(items, s.valueToString(item))
		}
		return fmt.Sprintf("[%s]", fmt.Sprintf("%v", items))
	default:
		return fmt.Sprintf("%v", v)
	}
}
