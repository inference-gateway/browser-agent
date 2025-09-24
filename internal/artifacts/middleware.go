package artifacts

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ArtifactMiddleware provides middleware to automatically include artifacts in task responses
type ArtifactMiddleware struct {
	logger   *zap.Logger
	registry *ArtifactRegistry
	baseURL  string
}

// NewArtifactMiddleware creates new artifact middleware
func NewArtifactMiddleware(logger *zap.Logger, registry *ArtifactRegistry, baseURL string) *ArtifactMiddleware {
	return &ArtifactMiddleware{
		logger:   logger,
		registry: registry,
		baseURL:  baseURL,
	}
}

// ResponseInterceptor intercepts HTTP responses to add artifacts when requested
func (m *ArtifactMiddleware) ResponseInterceptor() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Check if this is a task-related endpoint and if artifacts should be included
		if !m.shouldProcessRequest(c) {
			c.Next()
			return
		}

		// Capture the response
		blw := &bodyLogWriter{body: &strings.Builder{}, ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// Process the response to add artifacts if needed
		responseData := blw.body.String()
		if includeArtifacts := m.shouldIncludeArtifacts(c, responseData); includeArtifacts {
			modifiedResponse := m.addArtifactsToResponse(responseData)
			c.Writer = blw.ResponseWriter
			c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(modifiedResponse)))
			_, _ = c.Writer.Write([]byte(modifiedResponse))
		} else {
			c.Writer = blw.ResponseWriter
			_, _ = c.Writer.Write([]byte(responseData))
		}
	})
}

// shouldProcessRequest determines if the request should be processed for artifacts
func (m *ArtifactMiddleware) shouldProcessRequest(c *gin.Context) bool {
	// Check if this is a task-related endpoint
	path := c.Request.URL.Path
	return strings.Contains(path, "task") || strings.Contains(path, "jsonrpc")
}

// shouldIncludeArtifacts checks if artifacts should be included in the response
func (m *ArtifactMiddleware) shouldIncludeArtifacts(c *gin.Context, responseData string) bool {
	// Check query parameters
	if includeParam := c.Query("includeArtifacts"); includeParam == "true" {
		return true
	}

	// Check request body for includeArtifacts parameter
	if c.Request.Method == "POST" {
		var requestBody map[string]interface{}
		if bodyBytes, err := c.GetRawData(); err == nil {
			if json.Unmarshal(bodyBytes, &requestBody) == nil {
				if params, ok := requestBody["params"].(map[string]interface{}); ok {
					if include, ok := params["includeArtifacts"].(bool); ok && include {
						return true
					}
				}
			}
		}
	}

	return false
}

// addArtifactsToResponse adds artifact information to the response
func (m *ArtifactMiddleware) addArtifactsToResponse(responseData string) string {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseData), &response); err != nil {
		m.logger.Warn("failed to parse response for artifact injection", zap.Error(err))
		return responseData
	}

	// Get all artifacts
	artifacts := m.registry.ListArtifacts()
	if len(artifacts) == 0 {
		return responseData // No artifacts to add
	}

	// Create artifact list in A2A format
	artifactList := make([]map[string]interface{}, 0, len(artifacts))
	for _, artifact := range artifacts {
		artifactInfo := map[string]interface{}{
			"artifactId":  artifact.ID,
			"title":       artifact.Title,
			"description": artifact.Description,
			"filename":    artifact.FileName,
			"mimeType":    artifact.MimeType,
			"size":        artifact.Size,
			"createdAt":   artifact.CreatedAt,
			"fileUri":     fmt.Sprintf("%s/artifacts/%s", m.baseURL, artifact.ID),
		}

		// Determine part type
		switch artifact.MimeType {
		case "text/plain":
			artifactInfo["type"] = "TextPart"
		case "application/json":
			artifactInfo["type"] = "DataPart"
		default:
			artifactInfo["type"] = "FilePart"
		}

		if artifact.Metadata != nil {
			artifactInfo["metadata"] = artifact.Metadata
		}

		artifactList = append(artifactList, artifactInfo)
	}

	// Add artifacts to the response
	if result, ok := response["result"]; ok {
		if resultMap, ok := result.(map[string]interface{}); ok {
			resultMap["artifacts"] = artifactList
		} else {
			// If result is not a map, wrap it
			response["result"] = map[string]interface{}{
				"originalResult": result,
				"artifacts":      artifactList,
			}
		}
	} else {
		// Add artifacts at the root level
		response["artifacts"] = artifactList
	}

	modifiedData, err := json.Marshal(response)
	if err != nil {
		m.logger.Warn("failed to marshal response with artifacts", zap.Error(err))
		return responseData
	}

	m.logger.Info("added artifacts to response", zap.Int("artifactCount", len(artifactList)))
	return string(modifiedData)
}

// bodyLogWriter captures response body for modification
type bodyLogWriter struct {
	gin.ResponseWriter
	body *strings.Builder
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return len(b), nil
}