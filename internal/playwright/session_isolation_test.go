package playwright

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/inference-gateway/browser-agent/config"
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

	ctx := context.Background()

	session1, err := service.GetOrCreateTaskSession(ctx)
	require.NoError(t, err)
	assert.NotNil(t, session1)

	session2, err := service.GetOrCreateTaskSession(ctx)
	require.NoError(t, err)
	assert.NotNil(t, session2)

	session3, err := service.GetOrCreateTaskSession(ctx)
	require.NoError(t, err)
	assert.NotNil(t, session3)

	// Verify each session has a unique ID
	assert.NotEqual(t, session1.ID, session2.ID, "Session IDs should be unique")
	assert.NotEqual(t, session1.ID, session3.ID, "Session IDs should be unique")
	assert.NotEqual(t, session2.ID, session3.ID, "Session IDs should be unique")

	assert.Contains(t, session1.ID, "task_", "Session should be task-scoped")
	assert.Contains(t, session2.ID, "task_", "Session should be task-scoped")
	assert.Contains(t, session3.ID, "task_", "Session should be task-scoped")

	assert.NotEqual(t, session1.Browser, session2.Browser, "Each session should have its own browser instance")
	assert.NotEqual(t, session1.Context, session2.Context, "Each session should have its own context")
	assert.NotEqual(t, session1.Page, session2.Page, "Each session should have its own page")

	assert.True(t, session1.ExpiresAt.After(time.Now()), "Session should have future expiration")
	assert.True(t, session2.ExpiresAt.After(time.Now()), "Session should have future expiration")
	assert.True(t, session3.ExpiresAt.After(time.Now()), "Session should have future expiration")

	err = service.CloseBrowser(ctx, session1.ID)
	assert.NoError(t, err)
	err = service.CloseBrowser(ctx, session2.ID)
	assert.NoError(t, err)
	err = service.CloseBrowser(ctx, session3.ID)
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

	ctx := context.Background()

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

func TestDefaultSessionStillWorks(t *testing.T) {
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

	ctx := context.Background()

	defaultSession1, err := service.GetOrCreateDefaultSession(ctx)
	require.NoError(t, err)
	assert.Equal(t, "default", defaultSession1.ID)

	defaultSession2, err := service.GetOrCreateDefaultSession(ctx)
	require.NoError(t, err)
	assert.Equal(t, defaultSession1.ID, defaultSession2.ID)
	assert.Equal(t, defaultSession1.Browser, defaultSession2.Browser)

	taskSession, err := service.GetOrCreateTaskSession(ctx)
	require.NoError(t, err)
	assert.NotEqual(t, defaultSession1.ID, taskSession.ID)
	assert.NotEqual(t, defaultSession1.Browser, taskSession.Browser)

	err = service.CloseBrowser(ctx, taskSession.ID)
	assert.NoError(t, err)
}