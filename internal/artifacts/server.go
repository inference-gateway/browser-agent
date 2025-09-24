package artifacts

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ArtifactServer handles serving artifacts over HTTP
type ArtifactServer struct {
	logger      *zap.Logger
	port        int
	artifactDir string
	registry    *ArtifactRegistry
	httpServer  *http.Server
	mu          sync.RWMutex
	running     bool
}

// ArtifactRegistry manages artifact ID to file path mappings
type ArtifactRegistry struct {
	mu        sync.RWMutex
	artifacts map[string]*ArtifactEntry
}

// ArtifactEntry represents an artifact stored in the registry
type ArtifactEntry struct {
	ID          string                 `json:"id"`
	FilePath    string                 `json:"file_path"`
	FileName    string                 `json:"file_name"`
	MimeType    string                 `json:"mime_type"`
	Size        int64                  `json:"size"`
	CreatedAt   time.Time              `json:"created_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// NewArtifactServer creates a new artifacts server
func NewArtifactServer(logger *zap.Logger, port int, artifactDir string) *ArtifactServer {
	return &ArtifactServer{
		logger:      logger,
		port:        port,
		artifactDir: artifactDir,
		registry:    NewArtifactRegistry(),
	}
}

// NewArtifactRegistry creates a new artifact registry
func NewArtifactRegistry() *ArtifactRegistry {
	return &ArtifactRegistry{
		artifacts: make(map[string]*ArtifactEntry),
	}
}

// RegisterArtifact registers an artifact in the registry
func (r *ArtifactRegistry) RegisterArtifact(entry *ArtifactEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.artifacts[entry.ID] = entry
}

// GetArtifact retrieves an artifact from the registry
func (r *ArtifactRegistry) GetArtifact(id string) (*ArtifactEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, exists := r.artifacts[id]
	return entry, exists
}


// Start starts the artifacts server
func (s *ArtifactServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("artifacts server is already running")
	}

	// Create artifacts directory if it doesn't exist
	if err := os.MkdirAll(s.artifactDir, 0755); err != nil {
		return fmt.Errorf("failed to create artifacts directory: %w", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Add middleware
	router.Use(gin.Recovery())
	router.Use(s.loggingMiddleware())

	// Routes
	v1 := router.Group("/artifacts")
	{
		v1.GET("/:id", s.getArtifact)
		v1.GET("/:id/metadata", s.getArtifactMetadata)
	}

	// Health check
	router.GET("/health", s.healthCheck)

	s.httpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(s.port),
		Handler: router,
	}

	s.running = true
	s.logger.Info("starting artifacts server", zap.Int("port", s.port))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("artifacts server failed", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the artifacts server
func (s *ArtifactServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || s.httpServer == nil {
		return nil
	}

	s.logger.Info("stopping artifacts server")
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shutdown artifacts server gracefully", zap.Error(err))
		return err
	}

	s.running = false
	s.httpServer = nil
	s.logger.Info("artifacts server stopped")
	
	return nil
}

// IsRunning returns whether the server is running
func (s *ArtifactServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetRegistry returns the artifact registry
func (s *ArtifactServer) GetRegistry() *ArtifactRegistry {
	return s.registry
}

// getArtifact serves the artifact file
func (s *ArtifactServer) getArtifact(c *gin.Context) {
	artifactID := c.Param("id")
	if artifactID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artifact ID is required"})
		return
	}

	entry, exists := s.registry.GetArtifact(artifactID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "artifact not found"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(entry.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "artifact file not found on disk"})
		return
	}

	// Set appropriate headers
	c.Header("Content-Type", entry.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", entry.FileName))
	
	// Serve the file
	c.File(entry.FilePath)
}

// getArtifactMetadata returns metadata for an artifact
func (s *ArtifactServer) getArtifactMetadata(c *gin.Context) {
	artifactID := c.Param("id")
	if artifactID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artifact ID is required"})
		return
	}

	entry, exists := s.registry.GetArtifact(artifactID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "artifact not found"})
		return
	}

	c.JSON(http.StatusOK, entry)
}


// healthCheck returns server health status
func (s *ArtifactServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"server": "artifacts",
		"port":   s.port,
	})
}

// loggingMiddleware adds request logging
func (s *ArtifactServer) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		s.logger.Info("artifacts server request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("ip", param.ClientIP),
		)
		return ""
	})
}