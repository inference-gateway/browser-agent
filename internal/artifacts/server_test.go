package artifacts

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestArtifactServer(t *testing.T) {
	logger := zaptest.NewLogger(t)
	
	// Create temporary directory for artifacts
	tempDir := t.TempDir()
	
	// Create test server
	server := NewArtifactServer(logger, 0, tempDir) // Port 0 for random available port
	
	// Test that server starts successfully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := server.Start(ctx)
	require.NoError(t, err)
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Verify server is running
	assert.True(t, server.IsRunning())
	
	// Cleanup
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Stop(ctx)
	}()
}

func TestArtifactRegistry(t *testing.T) {
	registry := NewArtifactRegistry()
	
	// Test empty registry
	artifacts := registry.ListArtifacts()
	assert.Empty(t, artifacts)
	
	// Test registering an artifact
	testArtifact := &ArtifactEntry{
		ID:          "test-id",
		FilePath:    "/tmp/test.png",
		FileName:    "test.png",
		MimeType:    "image/png",
		Size:        1024,
		CreatedAt:   time.Now(),
		Title:       "Test Screenshot",
		Description: "Test screenshot description",
		Metadata: map[string]interface{}{
			"test_key": "test_value",
		},
	}
	
	registry.RegisterArtifact(testArtifact)
	
	// Test retrieval
	artifact, exists := registry.GetArtifact("test-id")
	require.True(t, exists)
	assert.Equal(t, testArtifact.ID, artifact.ID)
	assert.Equal(t, testArtifact.FileName, artifact.FileName)
	assert.Equal(t, testArtifact.MimeType, artifact.MimeType)
	
	// Test listing
	artifacts = registry.ListArtifacts()
	assert.Len(t, artifacts, 1)
	assert.Equal(t, testArtifact.ID, artifacts[0].ID)
	
	// Test non-existent artifact
	_, exists = registry.GetArtifact("non-existent")
	assert.False(t, exists)
}

func TestArtifactServerEndpoints(t *testing.T) {
	logger := zaptest.NewLogger(t)
	
	// Create temporary directory and test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)
	
	// Create server with port 0 (random available port)
	server := NewArtifactServer(logger, 0, tempDir)
	
	// Register test artifact
	testArtifact := &ArtifactEntry{
		ID:          "test-artifact",
		FilePath:    testFile,
		FileName:    "test.txt",
		MimeType:    "text/plain",
		Size:        int64(len(testContent)),
		CreatedAt:   time.Now(),
		Title:       "Test File",
		Description: "Test file description",
	}
	
	server.GetRegistry().RegisterArtifact(testArtifact)
	
	// Start server
	ctx := context.Background()
	err = server.Start(ctx)
	require.NoError(t, err)
	
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Stop(ctx)
	}()
	
	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	
	// Since we're using port 0, we need to get the actual port
	// For now, we'll skip the HTTP tests as we can't easily get the port
	// In a real implementation, you'd expose the actual port
	t.Skip("HTTP endpoint testing requires knowing the actual port - would need server modification to expose port")
	
	// TODO: Add HTTP client tests when server exposes actual port
	// This would test:
	// - GET /artifacts/test-artifact (should return file content)
	// - GET /artifacts/test-artifact/metadata (should return metadata)
	// - GET /artifacts/ (should return artifact list)
	// - GET /health (should return healthy status)
}

func TestArtifactHelper(t *testing.T) {
	logger := zaptest.NewLogger(t)
	registry := NewArtifactRegistry()
	
	helper := NewArtifactHelper(logger, registry)
	require.NotNil(t, helper)
	
	// Test text artifact creation
	textArtifact, err := helper.CreateTextArtifact(
		"Test Text",
		"Test text description",
		"Hello, World!",
		map[string]interface{}{"source": "test"},
	)
	
	assert.NoError(t, err)
	assert.NotNil(t, textArtifact)
	
	// Verify artifact was registered
	artifacts := registry.ListArtifacts()
	assert.Len(t, artifacts, 1)
	assert.Equal(t, "text/plain", artifacts[0].MimeType)
	assert.Contains(t, artifacts[0].FileName, ".txt")
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal_file.txt", "normal_file.txt"},
		{"file with spaces.txt", "file_with_spaces.txt"},
		{"file/with/slashes.txt", "file_with_slashes.txt"},
		{"file\\with\\backslashes.txt", "file_with_backslashes.txt"},
		{"file:with:colons.txt", "file_with_colons.txt"},
		{"file*with*stars.txt", "file_with_stars.txt"},
		{"", "artifact"},
	}
	
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := sanitizeFilename(test.input)
			assert.Equal(t, test.expected, result)
			// Ensure result doesn't contain invalid characters
			invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
			for _, char := range invalidChars {
				assert.False(t, strings.Contains(result, char), "Result contains invalid character: %s", char)
			}
		})
	}
}