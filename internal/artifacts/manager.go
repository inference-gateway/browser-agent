package artifacts

import (
	"sync"

	"go.uber.org/zap"
)

// Manager provides a global access point for artifact management
type Manager struct {
	registry *ArtifactRegistry
	helper   *ArtifactHelper
	server   *ArtifactServer
	logger   *zap.Logger
	mu       sync.RWMutex
}

var globalManager *Manager
var managerMutex sync.RWMutex

// InitializeGlobalManager initializes the global artifact manager
func InitializeGlobalManager(logger *zap.Logger, server *ArtifactServer) *Manager {
	managerMutex.Lock()
	defer managerMutex.Unlock()
	
	registry := server.GetRegistry()
	helper := NewArtifactHelper(logger, registry)
	
	globalManager = &Manager{
		registry: registry,
		helper:   helper,
		server:   server,
		logger:   logger,
	}
	
	return globalManager
}

// ResetGlobalManager resets the global manager (for testing)
func ResetGlobalManager() {
	managerMutex.Lock()
	defer managerMutex.Unlock()
	globalManager = nil
}

// GetGlobalManager returns the global artifact manager instance
func GetGlobalManager() *Manager {
	managerMutex.RLock()
	defer managerMutex.RUnlock()
	return globalManager
}

// GetHelper returns the artifact helper
func (m *Manager) GetHelper() *ArtifactHelper {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.helper
}

// GetRegistry returns the artifact registry
func (m *Manager) GetRegistry() *ArtifactRegistry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.registry
}

// GetServer returns the artifact server
func (m *Manager) GetServer() *ArtifactServer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.server
}