package playwright

import (
	"context"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	zap "go.uber.org/zap"

	server "github.com/inference-gateway/adk/server"
	types "github.com/inference-gateway/adk/types"

	config "github.com/inference-gateway/browser-agent/config"
)

func TestMultiTenantSessionIsolation(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		Browser: config.BrowserConfig{
			Headless:       true,
			Engine:         "chromium",
			ViewportWidth:  "1920",
			ViewportHeight: "1080",
		},
	}

	service, err := NewPlaywrightService(logger, cfg)
	require.NoError(t, err)
	defer func() {
		err := service.Shutdown(context.Background())
		assert.NoError(t, err)
	}()

	ctx1 := context.WithValue(context.Background(), server.TaskContextKey, &types.Task{ID: "task-1"})
	ctx2 := context.WithValue(context.Background(), server.TaskContextKey, &types.Task{ID: "task-2"})
	ctx3 := context.WithValue(context.Background(), server.TaskContextKey, &types.Task{ID: "task-3"})

	session1, err := service.GetOrCreateTaskSession(ctx1)
	require.NoError(t, err)
	assert.NotNil(t, session1)

	session2, err := service.GetOrCreateTaskSession(ctx2)
	require.NoError(t, err)
	assert.NotNil(t, session2)

	session3, err := service.GetOrCreateTaskSession(ctx3)
	require.NoError(t, err)
	assert.NotNil(t, session3)

	assert.NotEqual(t, session1.ID, session2.ID, "Session IDs should be unique")
	assert.NotEqual(t, session1.ID, session3.ID, "Session IDs should be unique")
	assert.NotEqual(t, session2.ID, session3.ID, "Session IDs should be unique")

	assert.Equal(t, "task-1", session1.ID, "Session ID should match task ID")
	assert.Equal(t, "task-2", session2.ID, "Session ID should match task ID")
	assert.Equal(t, "task-3", session3.ID, "Session ID should match task ID")

	assert.NotEqual(t, session1.Browser, session2.Browser, "Each session should have its own browser instance")
	assert.NotEqual(t, session1.Context, session2.Context, "Each session should have its own context")
	assert.NotEqual(t, session1.Page, session2.Page, "Each session should have its own page")

	assert.True(t, session1.ExpiresAt.After(time.Now()), "Session should have future expiration")
	assert.True(t, session2.ExpiresAt.After(time.Now()), "Session should have future expiration")
	assert.True(t, session3.ExpiresAt.After(time.Now()), "Session should have future expiration")

	err = service.CloseBrowser(ctx1, session1.ID)
	assert.NoError(t, err)
	err = service.CloseBrowser(ctx2, session2.ID)
	assert.NoError(t, err)
	err = service.CloseBrowser(ctx3, session3.ID)
	assert.NoError(t, err)
}

func TestSessionExpiration(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		Browser: config.BrowserConfig{
			Headless:       true,
			Engine:         "chromium",
			ViewportWidth:  "1920",
			ViewportHeight: "1080",
		},
	}

	service, err := NewPlaywrightService(logger, cfg)
	require.NoError(t, err)
	defer func() {
		err := service.Shutdown(context.Background())
		assert.NoError(t, err)
	}()

	ctx := context.WithValue(context.Background(), server.TaskContextKey, &types.Task{ID: "task-expiration-test"})

	session, err := service.GetOrCreateTaskSession(ctx)
	require.NoError(t, err)

	playwrightService := service.(*playwrightImpl)
	playwrightService.sessionsMux.Lock()
	playwrightService.sessions[session.ID].ExpiresAt = time.Now().Add(-1 * time.Minute)
	playwrightService.sessionsMux.Unlock()

	_, err = service.GetSession(session.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session expired")

	err = service.CloseExpiredSessions(ctx)
	assert.NoError(t, err)

	_, err = service.GetSession(session.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}
