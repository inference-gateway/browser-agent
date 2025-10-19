package skills

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	server "github.com/inference-gateway/adk/server"
	types "github.com/inference-gateway/adk/types"
	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
	zap "go.uber.org/zap"
)

// TakeScreenshotSkill struct holds the skill with dependencies
type TakeScreenshotSkill struct {
	logger        *zap.Logger
	playwright    playwright.BrowserAutomation
	screenshotDir string
}

// NewTakeScreenshotSkill creates a new take_screenshot skill
func NewTakeScreenshotSkill(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	cfg := playwright.GetConfig()
	skill := &TakeScreenshotSkill{
		logger:        logger,
		playwright:    playwright,
		screenshotDir: cfg.Browser.DataDir,
	}
	return server.NewBasicTool(
		"take_screenshot",
		"Capture a screenshot of the current page or specific element with deterministic file naming",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"full_page": map[string]any{
					"default":     false,
					"description": "Capture the entire scrollable page",
					"type":        "boolean",
				},
				"quality": map[string]any{
					"default":     80,
					"description": "Quality for jpeg images (0-100)",
					"type":        "integer",
				},
				"selector": map[string]any{
					"description": "Optional selector to screenshot specific element",
					"type":        "string",
				},
				"type": map[string]any{
					"default":     "png",
					"description": "Image format (png, jpeg)",
					"type":        "string",
				},
			},
			"required": []string{},
		},
		skill.TakeScreenshotHandler,
	)
}

// TakeScreenshotHandler handles the take_screenshot skill execution
func (s *TakeScreenshotSkill) TakeScreenshotHandler(ctx context.Context, args map[string]any) (string, error) {

	generatedPath, err := s.generateDeterministicPath(args)
	if err != nil {
		s.logger.Error("failed to generate screenshot path", zap.Error(err))
		return "", fmt.Errorf("failed to generate screenshot path: %w", err)
	}

	fullPage := false
	if fp, ok := args["full_page"].(bool); ok {
		fullPage = fp
	}

	quality := 80
	if q, ok := args["quality"].(int); ok {
		quality = q
	} else if qf, ok := args["quality"].(float64); ok {
		quality = int(qf)
	}

	imageType := "png"
	if t, ok := args["type"].(string); ok && t != "" {
		if !s.isValidImageType(t) {
			return "", fmt.Errorf("invalid image type: %s. Must be 'png' or 'jpeg'", t)
		}
		imageType = t
	}

	if imageType == "jpeg" && (quality < 0 || quality > 100) {
		return "", fmt.Errorf("quality must be between 0 and 100 for JPEG images, got %d", quality)
	}

	selector := ""
	if s, ok := args["selector"].(string); ok {
		selector = s
	}

	s.logger.Info("taking screenshot",
		zap.String("path", generatedPath),
		zap.Bool("full_page", fullPage),
		zap.String("type", imageType),
		zap.Int("quality", quality),
		zap.String("selector", selector))

	session, err := s.playwright.GetOrCreateTaskSession(ctx)
	if err != nil {
		s.logger.Error("failed to get browser session", zap.Error(err))
		return "", fmt.Errorf("failed to get browser session: %w", err)
	}

	err = s.playwright.TakeScreenshot(ctx, session.ID, generatedPath, fullPage, selector, imageType, quality)
	if err != nil {
		s.logger.Error("screenshot failed",
			zap.String("path", generatedPath),
			zap.String("sessionID", session.ID),
			zap.Error(err))
		return "", fmt.Errorf("screenshot failed: %w", err)
	}

	s.logger.Info("screenshot completed successfully",
		zap.String("sessionID", session.ID),
		zap.String("path", generatedPath))

	artifactURL, artifactID, err := s.createArtifactFromScreenshot(ctx, generatedPath, imageType)
	if err != nil {
		s.logger.Debug("artifact creation skipped or failed, returning file path only",
			zap.Error(err),
			zap.String("path", generatedPath))

		response := map[string]any{
			"success":    true,
			"path":       generatedPath,
			"filename":   filepath.Base(generatedPath),
			"full_page":  fullPage,
			"type":       imageType,
			"quality":    quality,
			"selector":   selector,
			"session_id": session.ID,
			"timestamp":  s.getCurrentTimestamp(),
			"message":    fmt.Sprintf("Screenshot captured successfully and saved to %s", generatedPath),
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return "", fmt.Errorf("failed to marshal response: %w", err)
		}
		return string(responseJSON), nil
	}

	response := map[string]any{
		"success":     true,
		"path":        generatedPath,
		"filename":    filepath.Base(generatedPath),
		"full_page":   fullPage,
		"type":        imageType,
		"quality":     quality,
		"selector":    selector,
		"session_id":  session.ID,
		"timestamp":   s.getCurrentTimestamp(),
		"artifact_id": artifactID,
		"url":         artifactURL,
		"message":     fmt.Sprintf("Screenshot captured successfully. Download URL: %s", artifactURL),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(responseJSON), nil
}

// generateDeterministicPath generates a deterministic file path for the screenshot
func (s *TakeScreenshotSkill) generateDeterministicPath(args map[string]any) (string, error) {
	if err := os.MkdirAll(s.screenshotDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create screenshots directory: %w", err)
	}

	imageType := "png"
	if t, ok := args["type"].(string); ok && t != "" {
		imageType = t
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05.000")

	var filename string
	if fullPage, ok := args["full_page"].(bool); ok && fullPage {
		filename = fmt.Sprintf("fullpage_%s.%s", timestamp, imageType)
	} else if selector, ok := args["selector"].(string); ok && selector != "" {
		safeSelector := filepath.Base(selector)
		if len(safeSelector) > 20 {
			safeSelector = safeSelector[:20]
		}
		filename = fmt.Sprintf("element_%s_%s.%s", safeSelector, timestamp, imageType)
	} else {
		filename = fmt.Sprintf("viewport_%s.%s", timestamp, imageType)
	}

	fullPath := filepath.Join(s.screenshotDir, filename)
	return fullPath, nil
}

// isValidImageType validates the image format
func (s *TakeScreenshotSkill) isValidImageType(imageType string) bool {
	validTypes := []string{"png", "jpeg"}
	for _, valid := range validTypes {
		if imageType == valid {
			return true
		}
	}
	return false
}

// getMimeType returns the MIME type for the given image format
func (s *TakeScreenshotSkill) getMimeType(imageType string) string {
	switch imageType {
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "image/png"
	}
}

// getCurrentTimestamp returns the current timestamp in RFC3339 format
func (s *TakeScreenshotSkill) getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// createArtifactFromScreenshot creates an artifact from the screenshot file
func (s *TakeScreenshotSkill) createArtifactFromScreenshot(ctx context.Context, filePath, imageType string) (url string, artifactID string, err error) {
	task, ok := ctx.Value(server.TaskContextKey).(*types.Task)
	if !ok {
		return "", "", fmt.Errorf("task not found in context")
	}

	artifactService, ok := ctx.Value(server.ArtifactServiceContextKey).(server.ArtifactService)
	if !ok || artifactService == nil {
		return "", "", fmt.Errorf("artifact service not available")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read screenshot file: %w", err)
	}

	mimeType := s.getMimeType(imageType)

	filename := filepath.Base(filePath)
	artifact, err := artifactService.CreateFileArtifact(
		fmt.Sprintf("Screenshot - %s", filename),
		fmt.Sprintf("Screenshot captured at %s", s.getCurrentTimestamp()),
		filename,
		data,
		&mimeType,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create artifact: %w", err)
	}

	artifactService.AddArtifactToTask(task, artifact)

	if len(artifact.Parts) > 0 {
		if filePart, ok := artifact.Parts[0].(types.FilePart); ok {
			if fileWithURI, ok := filePart.File.(types.FileWithUri); ok {
				return fileWithURI.URI, artifact.ArtifactID, nil
			}
		}
	}

	return "", artifact.ArtifactID, nil
}
