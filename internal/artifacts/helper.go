package artifacts

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	server "github.com/inference-gateway/adk/server"
	"go.uber.org/zap"
)

// ArtifactHelper wraps ADK's ArtifactHelper with registry integration for REST endpoints
type ArtifactHelper struct {
	*server.ArtifactHelper
	logger   *zap.Logger
	registry *ArtifactRegistry
}

// NewArtifactHelper creates a new artifact helper with registry support
func NewArtifactHelper(logger *zap.Logger, registry *ArtifactRegistry) *ArtifactHelper {
	return &ArtifactHelper{
		ArtifactHelper: server.NewArtifactHelper(),
		logger:         logger,
		registry:       registry,
	}
}

// CreateFileArtifactFromPath creates a file artifact from a file path and registers it
func (h *ArtifactHelper) CreateFileArtifactFromPath(title, description, filePath, mimeType string, metadata map[string]interface{}) (interface{}, error) {
	// Get file info for registry
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Read file data for ADK artifact creation
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	filename := filepath.Base(filePath)

	// Create ADK artifact
	artifact := h.CreateFileArtifactFromBytes(title, description, filename, data, &mimeType)
	if metadata != nil {
		artifact.Metadata = metadata
	}

	// Register in our registry for REST endpoint access
	h.registerArtifact(artifact.ArtifactID, filePath, filename, mimeType, fileInfo.Size(), fileInfo.ModTime(), metadata, title, description)

	return artifact, nil
}

// CreateFileArtifactFromBytesWithRegistry creates a file artifact from bytes and registers it
func (h *ArtifactHelper) CreateFileArtifactFromBytesWithRegistry(title, description, filename string, data []byte, mimeType *string, metadata map[string]interface{}) (interface{}, error) {
	// Create ADK artifact
	artifact := h.CreateFileArtifactFromBytes(title, description, filename, data, mimeType)
	if metadata != nil {
		artifact.Metadata = metadata
	}

	// Save to disk for REST endpoint access
	artifactPath := filepath.Join("/tmp/artifacts/runtime", artifact.ArtifactID, filename)
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0755); err != nil {
		h.logger.Warn("failed to create runtime artifact directory", zap.Error(err))
		return artifact, nil // Continue without registry entry
	}

	if err := os.WriteFile(artifactPath, data, 0644); err != nil {
		h.logger.Warn("failed to save runtime artifact to disk", zap.Error(err))
		return artifact, nil // Continue without registry entry
	}

	resolvedMimeType := "application/octet-stream"
	if mimeType != nil {
		resolvedMimeType = *mimeType
	}

	// Register for REST endpoint access
	h.registerArtifact(artifact.ArtifactID, artifactPath, filename, resolvedMimeType, int64(len(data)), time.Now(), metadata, title, description)

	return artifact, nil
}

// CreateTextArtifact creates a text artifact
func (h *ArtifactHelper) CreateTextArtifact(title, description, content string, metadata map[string]interface{}) (interface{}, error) {
	data := []byte(content)
	mimeType := "text/plain"
	filename := fmt.Sprintf("%s.txt", sanitizeFilename(title))

	return h.CreateFileArtifactFromBytesWithRegistry(title, description, filename, data, &mimeType, metadata)
}

// CreateDataArtifact creates a data artifact (JSON)
func (h *ArtifactHelper) CreateDataArtifact(title, description string, data []byte, metadata map[string]interface{}) (interface{}, error) {
	mimeType := "application/json"
	filename := fmt.Sprintf("%s.json", sanitizeFilename(title))

	return h.CreateFileArtifactFromBytesWithRegistry(title, description, filename, data, &mimeType, metadata)
}

// registerArtifact is a helper to register artifacts in the registry
func (h *ArtifactHelper) registerArtifact(id, filePath, fileName, mimeType string, size int64, createdAt time.Time, metadata map[string]interface{}, title, description string) {
	entry := &ArtifactEntry{
		ID:          id,
		FilePath:    filePath,
		FileName:    fileName,
		MimeType:    mimeType,
		Size:        size,
		CreatedAt:   createdAt,
		Metadata:    metadata,
		Title:       title,
		Description: description,
	}

	h.registry.RegisterArtifact(entry)
	h.logger.Info("artifact registered",
		zap.String("artifactID", id),
		zap.String("filename", fileName),
		zap.Int64("size", size),
	)
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
	if sanitized == "" || sanitized == "." {
		sanitized = "artifact"
	}
	return sanitized
}