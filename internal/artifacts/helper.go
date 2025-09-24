package artifacts

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	server "github.com/inference-gateway/adk/server"
	"go.uber.org/zap"
)

// EnhancedArtifactHelper extends the ADK ArtifactHelper with registry support
type EnhancedArtifactHelper struct {
	*server.ArtifactHelper
	logger   *zap.Logger
	registry *ArtifactRegistry
	baseURL  string
}

// NewEnhancedArtifactHelper creates a new enhanced artifact helper
func NewEnhancedArtifactHelper(logger *zap.Logger, registry *ArtifactRegistry, baseURL string) *EnhancedArtifactHelper {
	return &EnhancedArtifactHelper{
		ArtifactHelper: server.NewArtifactHelper(),
		logger:         logger,
		registry:       registry,
		baseURL:        baseURL,
	}
}

// CreateFileArtifactFromPath creates a file artifact from a file path and registers it
func (h *EnhancedArtifactHelper) CreateFileArtifactFromPath(title, description, filePath, mimeType string, metadata map[string]interface{}) (interface{}, error) {
	// Read file data
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	filename := filepath.Base(filePath)

	// Create ADK artifact
	artifact := h.CreateFileArtifactFromBytes(
		title,
		description,
		filename,
		data,
		&mimeType,
	)

	// Add metadata if provided
	if metadata != nil {
		artifact.Metadata = metadata
	}

	// Register in our registry
	entry := &ArtifactEntry{
		ID:          artifact.ArtifactID,
		FilePath:    filePath,
		FileName:    filename,
		MimeType:    mimeType,
		Size:        fileInfo.Size(),
		CreatedAt:   fileInfo.ModTime(),
		Metadata:    metadata,
		Title:       title,
		Description: description,
	}

	h.registry.RegisterArtifact(entry)

	h.logger.Info("artifact created and registered",
		zap.String("artifactID", artifact.ArtifactID),
		zap.String("filePath", filePath),
		zap.String("filename", filename),
		zap.Int64("size", fileInfo.Size()),
	)

	return artifact, nil
}

// CreateFileArtifactFromBytesWithRegistry creates a file artifact from bytes and registers it
func (h *EnhancedArtifactHelper) CreateFileArtifactFromBytesWithRegistry(title, description, filename string, data []byte, mimeType *string, metadata map[string]interface{}) (interface{}, error) {
	// Create ADK artifact
	artifact := h.CreateFileArtifactFromBytes(
		title,
		description,
		filename,
		data,
		mimeType,
	)

	// Add metadata if provided
	if metadata != nil {
		artifact.Metadata = metadata
	}

	// For byte-based artifacts, we need to save to disk first to register
	// This maintains consistency with file-based artifacts
	artifactPath := filepath.Join("/tmp/artifacts/runtime", artifact.ArtifactID, filename)
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0755); err != nil {
		h.logger.Warn("failed to create runtime artifact directory", zap.Error(err))
		// Continue without registry entry if we can't save to disk
		return artifact, nil
	}

	if err := os.WriteFile(artifactPath, data, 0644); err != nil {
		h.logger.Warn("failed to save runtime artifact to disk", zap.Error(err))
		// Continue without registry entry if we can't save to disk
		return artifact, nil
	}

	resolvedMimeType := "application/octet-stream"
	if mimeType != nil {
		resolvedMimeType = *mimeType
	}

	// Register in our registry
	entry := &ArtifactEntry{
		ID:          artifact.ArtifactID,
		FilePath:    artifactPath,
		FileName:    filename,
		MimeType:    resolvedMimeType,
		Size:        int64(len(data)),
		CreatedAt:   time.Now(),
		Metadata:    metadata,
		Title:       title,
		Description: description,
	}

	h.registry.RegisterArtifact(entry)

	h.logger.Info("artifact created from bytes and registered",
		zap.String("artifactID", artifact.ArtifactID),
		zap.String("filename", filename),
		zap.Int("size", len(data)),
	)

	return artifact, nil
}

// CreateTextArtifact creates a text artifact
func (h *EnhancedArtifactHelper) CreateTextArtifact(title, description, content string, metadata map[string]interface{}) (interface{}, error) {
	// Create text artifact using ADK helper as bytes
	data := []byte(content)
	mimeType := "text/plain"
	filename := fmt.Sprintf("%s.txt", sanitizeFilename(title))

	return h.CreateFileArtifactFromBytesWithRegistry(title, description, filename, data, &mimeType, metadata)
}

// CreateDataArtifact creates a data artifact (JSON)
func (h *EnhancedArtifactHelper) CreateDataArtifact(title, description string, data []byte, metadata map[string]interface{}) (interface{}, error) {
	mimeType := "application/json"
	filename := fmt.Sprintf("%s.json", sanitizeFilename(title))

	return h.CreateFileArtifactFromBytesWithRegistry(title, description, filename, data, &mimeType, metadata)
}

// GetArtifactURL returns the URL to access an artifact via the artifacts server
func (h *EnhancedArtifactHelper) GetArtifactURL(artifactID string) string {
	return fmt.Sprintf("%s/artifacts/%s", h.baseURL, artifactID)
}

// sanitizeFilename removes invalid characters from filename
func sanitizeFilename(name string) string {
	// Replace common invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	sanitized := name
	for _, invalidChar := range invalidChars {
		for i := 0; i < len(sanitized); i++ {
			if string(sanitized[i]) == invalidChar {
				sanitized = sanitized[:i] + "_" + sanitized[i+1:]
			}
		}
	}
	sanitized = filepath.Clean(sanitized)
	if sanitized == "" {
		sanitized = "artifact"
	}
	return sanitized
}