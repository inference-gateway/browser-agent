package artifacts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGlobalManagerIntegration(t *testing.T) {
	// Reset global manager before test
	ResetGlobalManager()
	
	logger := zaptest.NewLogger(t)
	
	// Create temporary directory
	tempDir := t.TempDir()
	
	// Create artifact server
	server := NewArtifactServer(logger, 8082, tempDir)
	
	// Initialize global manager
	manager := InitializeGlobalManager(logger, server)
	require.NotNil(t, manager)
	
	// Test that we can get the global manager
	retrievedManager := GetGlobalManager()
	assert.Equal(t, manager, retrievedManager)
	
	// Test helper access
	helper := manager.GetHelper()
	assert.NotNil(t, helper)
	
	// Test registry access
	registry := manager.GetRegistry()
	assert.NotNil(t, registry)
	
	// Test creating artifacts through helper
	testData := []byte("Test file content")
	mimeType := "text/plain"
	
	artifact, err := helper.CreateFileArtifactFromBytesWithRegistry(
		"Test File",
		"Test file description",
		"test.txt",
		testData,
		&mimeType,
		map[string]interface{}{"test": "value"},
	)
	
	require.NoError(t, err)
	assert.NotNil(t, artifact)
	
	// Verify artifact was registered
	artifacts := registry.ListArtifacts()
	assert.Len(t, artifacts, 1)
	assert.Equal(t, "test.txt", artifacts[0].FileName)
	assert.Equal(t, "text/plain", artifacts[0].MimeType)
	assert.Equal(t, int64(len(testData)), artifacts[0].Size)
	
	// URL generation removed - now handled by artifact server directly
}

func TestArtifactFileStorage(t *testing.T) {
	// Reset global manager before test
	ResetGlobalManager()
	
	logger := zaptest.NewLogger(t)
	
	// Create temporary directory
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test.png")
	
	// Create test file
	testData := []byte("fake PNG data")
	err := os.WriteFile(testFilePath, testData, 0644)
	require.NoError(t, err)
	
	// Create server and manager
	server := NewArtifactServer(logger, 8083, tempDir)
	manager := InitializeGlobalManager(logger, server)
	
	helper := manager.GetHelper()
	
	// Test file artifact creation from path
	artifact, err := helper.CreateFileArtifactFromPath(
		"Test Image",
		"Test image description",
		testFilePath,
		"image/png",
		map[string]interface{}{"width": 100, "height": 100},
	)
	
	require.NoError(t, err)
	assert.NotNil(t, artifact)
	
	// Verify registry entry
	registry := manager.GetRegistry()
	artifacts := registry.ListArtifacts()
	
	assert.Len(t, artifacts, 1)
	entry := artifacts[0]
	assert.Equal(t, testFilePath, entry.FilePath)
	assert.Equal(t, "test.png", entry.FileName)
	assert.Equal(t, "image/png", entry.MimeType)
	assert.Equal(t, int64(len(testData)), entry.Size)
	assert.NotNil(t, entry.Metadata)
	assert.Equal(t, 100, entry.Metadata["width"])
}

func TestMultipleArtifactTypes(t *testing.T) {
	// Reset global manager before test
	ResetGlobalManager()
	
	logger := zaptest.NewLogger(t)
	tempDir := t.TempDir()
	
	server := NewArtifactServer(logger, 8084, tempDir)
	manager := InitializeGlobalManager(logger, server)
	helper := manager.GetHelper()
	registry := manager.GetRegistry()
	
	// Create text artifact
	textArtifact, err := helper.CreateTextArtifact(
		"Sample Text",
		"A sample text file",
		"This is sample text content.",
		nil,
	)
	require.NoError(t, err)
	assert.NotNil(t, textArtifact)
	
	// Create data artifact (JSON)
	jsonData := []byte(`{"name": "test", "value": 123}`)
	dataArtifact, err := helper.CreateDataArtifact(
		"Sample Data",
		"A sample JSON data file",
		jsonData,
		map[string]interface{}{"format": "json"},
	)
	require.NoError(t, err)
	assert.NotNil(t, dataArtifact)
	
	// Create binary file artifact
	binaryData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG header
	mimeType := "image/jpeg"
	binaryArtifact, err := helper.CreateFileArtifactFromBytesWithRegistry(
		"Sample Image",
		"A sample binary image",
		"sample.jpg",
		binaryData,
		&mimeType,
		map[string]interface{}{"type": "binary"},
	)
	require.NoError(t, err)
	assert.NotNil(t, binaryArtifact)
	
	// Verify all artifacts are registered
	artifacts := registry.ListArtifacts()
	assert.Len(t, artifacts, 3)
	
	// Find each artifact type
	var textEntry, dataEntry, binaryEntry *ArtifactEntry
	for _, entry := range artifacts {
		switch entry.MimeType {
		case "text/plain":
			textEntry = entry
		case "application/json":
			dataEntry = entry
		case "image/jpeg":
			binaryEntry = entry
		}
	}
	
	require.NotNil(t, textEntry, "Text artifact not found")
	require.NotNil(t, dataEntry, "Data artifact not found")
	require.NotNil(t, binaryEntry, "Binary artifact not found")
	
	// Verify text artifact
	assert.Contains(t, textEntry.FileName, ".txt")
	assert.Equal(t, "Sample Text", textEntry.Title)
	
	// Verify data artifact
	assert.Contains(t, dataEntry.FileName, ".json")
	assert.Equal(t, "Sample Data", dataEntry.Title)
	assert.Equal(t, "json", dataEntry.Metadata["format"])
	
	// Verify binary artifact
	assert.Equal(t, "sample.jpg", binaryEntry.FileName)
	assert.Equal(t, "Sample Image", binaryEntry.Title)
	assert.Equal(t, "binary", binaryEntry.Metadata["type"])
}