package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"
	types "github.com/inference-gateway/adk/types"

	playwright "github.com/inference-gateway/browser-agent/internal/playwright"
)

// TakeScreenshotTool struct holds the tool with dependencies
type TakeScreenshotTool struct {
	logger        *zap.Logger
	playwright    playwright.BrowserAutomation
	screenshotDir string
}

// NewTakeScreenshotTool creates a new take_screenshot tool
func NewTakeScreenshotTool(logger *zap.Logger, playwright playwright.BrowserAutomation) server.Tool {
	cfg := playwright.GetConfig()
	tool := &TakeScreenshotTool{
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
		tool.TakeScreenshotHandler,
	)
}

// TakeScreenshotHandler handles the take_screenshot tool execution
func (s *TakeScreenshotTool) TakeScreenshotHandler(ctx context.Context, args map[string]any) (string, error) {
	fullPage, err := boolArg(args, "full_page", false)
	if err != nil {
		return "", err
	}

	imageType, err := stringArg(args, "type", "png")
	if err != nil {
		return "", err
	}
	if !oneOf(imageType, "png", "jpeg") {
		return "", fmt.Errorf("invalid image type: %s. Must be 'png' or 'jpeg'", imageType)
	}

	quality, err := boundedIntArg(args, "quality", defaultJPEGQuality, minJPEGQuality, maxJPEGQuality)
	if err != nil {
		return "", err
	}

	selector, err := stringArg(args, "selector", "")
	if err != nil {
		return "", err
	}

	generatedPath, err := s.generateDeterministicPath(fullPage, selector, imageType)
	if err != nil {
		s.logger.Error("failed to generate screenshot path", zap.Error(err))
		return "", fmt.Errorf("failed to generate screenshot path: %w", err)
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
	}

	artifactURL, artifactID, err := s.createArtifactFromScreenshot(ctx, generatedPath, imageType)
	if err != nil {
		s.logger.Debug("artifact creation skipped or failed, returning file path only",
			zap.Error(err),
			zap.String("path", generatedPath))
		response["message"] = fmt.Sprintf("Screenshot captured successfully and saved to %s", generatedPath)
		return marshalResponse(response)
	}

	response["artifact_id"] = artifactID
	response["url"] = artifactURL
	response["message"] = fmt.Sprintf("Screenshot captured successfully. Download URL: %s", artifactURL)
	return marshalResponse(response)
}

// generateDeterministicPath generates a deterministic file path for the screenshot.
// safeSelector is sliced by runes (not bytes) so multibyte selectors don't
// produce invalid-UTF-8 filenames.
func (s *TakeScreenshotTool) generateDeterministicPath(fullPage bool, selector, imageType string) (string, error) {
	if err := os.MkdirAll(s.screenshotDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create screenshots directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05.000")

	var filename string
	switch {
	case fullPage:
		filename = fmt.Sprintf("fullpage_%s.%s", timestamp, imageType)
	case selector != "":
		safeSelector := truncateRunes(filepath.Base(selector), screenshotSelectorMaxRunes)
		filename = fmt.Sprintf("element_%s_%s.%s", safeSelector, timestamp, imageType)
	default:
		filename = fmt.Sprintf("viewport_%s.%s", timestamp, imageType)
	}

	return filepath.Join(s.screenshotDir, filename), nil
}

// truncateRunes returns s truncated to at most maxRunes runes (not bytes).
// Slicing a byte string at an arbitrary index can split a multibyte rune
// and produce invalid UTF-8; this helper does it safely.
func truncateRunes(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) > maxRunes {
		return string(runes[:maxRunes])
	}
	return s
}

// getMimeType returns the MIME type for the given image format
func (s *TakeScreenshotTool) getMimeType(imageType string) string {
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
func (s *TakeScreenshotTool) getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// createArtifactFromScreenshot creates an artifact from the screenshot file
func (s *TakeScreenshotTool) createArtifactFromScreenshot(ctx context.Context, filePath, imageType string) (url string, artifactID string, err error) {
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
		if artifact.Parts[0].File != nil && artifact.Parts[0].File.FileWithURI != nil {
			return *artifact.Parts[0].File.FileWithURI, artifact.ArtifactID, nil
		}
	}

	return "", artifact.ArtifactID, nil
}
