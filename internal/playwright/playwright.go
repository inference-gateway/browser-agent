package playwright

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	config "github.com/inference-gateway/browser-agent/config"
	zap "go.uber.org/zap"

	"github.com/playwright-community/playwright-go"
)

// BrowserEngine represents the browser type
type BrowserEngine string

const (
	Chromium BrowserEngine = "chromium"
	Firefox  BrowserEngine = "firefox"
	WebKit   BrowserEngine = "webkit"
)

const (
	DefaultSessionID = "default"
)

// BrowserConfig holds browser configuration options
type BrowserConfig struct {
	Engine         BrowserEngine
	Headless       bool
	Timeout        time.Duration
	ViewportWidth  int
	ViewportHeight int
	Args           []string
}

// DefaultBrowserConfig returns default browser configuration
func DefaultBrowserConfig() *BrowserConfig {
	return &BrowserConfig{
		Engine:         Chromium,
		Headless:       true,
		Timeout:        30 * time.Second,
		ViewportWidth:  1920,
		ViewportHeight: 1080,
		Args: []string{
			"--disable-dev-shm-usage",
			"--no-sandbox",
			"--disable-blink-features=AutomationControlled",
			"--disable-features=VizDisplayCompositor",
			"--no-first-run",
			"--disable-default-apps",
			"--disable-extensions",
			"--disable-plugins",
			"--disable-sync",
			"--disable-translate",
			"--hide-scrollbars",
			"--mute-audio",
			"--no-zygote",
			"--disable-background-timer-throttling",
			"--disable-backgrounding-occluded-windows",
			"--disable-renderer-backgrounding",
			"--disable-ipc-flooding-protection",
		},
	}
}

// NewBrowserConfigFromConfig creates browser config from app configuration
func NewBrowserConfigFromConfig(cfg *config.Config) *BrowserConfig {
	width, _ := strconv.Atoi(cfg.Browser.ViewportWidth)
	if width == 0 {
		width = 1920
	}

	height, _ := strconv.Atoi(cfg.Browser.ViewportHeight)
	if height == 0 {
		height = 1080
	}

	argsStr := strings.Trim(cfg.Browser.Args, "[]")
	args := []string{"--disable-dev-shm-usage", "--no-sandbox"}
	if argsStr != "" {
		configArgs := strings.Fields(argsStr)
		args = append(args, configArgs...)
	}

	return &BrowserConfig{
		Engine:         Chromium,
		Headless:       cfg.Browser.Headless,
		Timeout:        30 * time.Second,
		ViewportWidth:  width,
		ViewportHeight: height,
		Args:           args,
	}
}

// BrowserSession represents an active browser session
type BrowserSession struct {
	ID       string
	Browser  playwright.Browser
	Context  playwright.BrowserContext
	Page     playwright.Page
	Created  time.Time
	LastUsed time.Time
}

// BrowserAutomation represents the playwright dependency interface
// Playwright service for browser automation and web testing
type BrowserAutomation interface {
	// Browser management
	LaunchBrowser(ctx context.Context, config *BrowserConfig) (*BrowserSession, error)
	CloseBrowser(ctx context.Context, sessionID string) error
	GetSession(sessionID string) (*BrowserSession, error)
	GetOrCreateDefaultSession(ctx context.Context) (*BrowserSession, error)

	// Page operations
	NavigateToURL(ctx context.Context, sessionID, url string, waitUntil string, timeout time.Duration) error
	ClickElement(ctx context.Context, sessionID, selector string, options map[string]any) error
	FillForm(ctx context.Context, sessionID string, fields []map[string]any, submit bool, submitSelector string) error
	ExtractData(ctx context.Context, sessionID string, extractors []map[string]any, format string) (string, error)
	TakeScreenshot(ctx context.Context, sessionID, path string, fullPage bool, selector string, format string, quality int) error
	ExecuteScript(ctx context.Context, sessionID, script string, args []any) (any, error)
	WaitForCondition(ctx context.Context, sessionID, condition, selector, state string, timeout time.Duration, customFunction string) error
	HandleAuthentication(ctx context.Context, sessionID, authType, username, password, loginURL string, selectors map[string]string) error

	// Service management
	GetHealth(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetConfig() *config.Config
}

// playwrightImpl is the implementation of BrowserAutomation
type playwrightImpl struct {
	logger      *zap.Logger
	config      *config.Config
	pw          *playwright.Playwright
	sessions    map[string]*BrowserSession
	sessionsMux sync.RWMutex
	isInstalled bool
}

// NewPlaywrightService creates a new instance of BrowserAutomation
func NewPlaywrightService(logger *zap.Logger, cfg *config.Config) (BrowserAutomation, error) {
	logger.Info("initializing playwright dependency")

	service := &playwrightImpl{
		logger:   logger,
		config:   cfg,
		sessions: make(map[string]*BrowserSession),
	}

	if err := service.ensurePlaywrightInstalled(); err != nil {
		return nil, fmt.Errorf("failed to ensure playwright installation: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to start playwright: %w", err)
	}
	service.pw = pw

	browserConfig := NewBrowserConfigFromConfig(cfg)
	logger.Info("playwright service initialized successfully",
		zap.String("engine", string(browserConfig.Engine)),
		zap.Bool("headless", browserConfig.Headless),
		zap.Int("viewport_width", browserConfig.ViewportWidth),
		zap.Int("viewport_height", browserConfig.ViewportHeight))

	return service, nil
}

// ensurePlaywrightInstalled checks and installs playwright browsers if needed
func (p *playwrightImpl) ensurePlaywrightInstalled() error {
	if p.isInstalled {
		return nil
	}

	p.logger.Info("checking playwright browser installation")

	if _, err := os.Stat(os.Getenv("PLAYWRIGHT_BROWSERS_PATH")); err != nil {
		p.logger.Info("installing playwright browsers")
		if err := playwright.Install(); err != nil {
			return fmt.Errorf("failed to install playwright browsers: %w", err)
		}
		p.logger.Info("playwright browsers installed successfully")
	}

	p.isInstalled = true
	return nil
}

// LaunchBrowser launches a new browser instance with the given configuration
func (p *playwrightImpl) LaunchBrowser(ctx context.Context, config *BrowserConfig) (*BrowserSession, error) {
	if config == nil {
		config = NewBrowserConfigFromConfig(p.config)
	}

	p.logger.Info("launching browser",
		zap.String("engine", string(config.Engine)),
		zap.Bool("headless", config.Headless))

	var browserType playwright.BrowserType
	switch config.Engine {
	case Chromium:
		browserType = p.pw.Chromium
	case Firefox:
		browserType = p.pw.Firefox
	case WebKit:
		browserType = p.pw.WebKit
	default:
		return nil, fmt.Errorf("unsupported browser engine: %s", config.Engine)
	}

	timeoutMs := float64(config.Timeout.Milliseconds())
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: &config.Headless,
		Args:     config.Args,
		Timeout:  &timeoutMs,
	}

	browser, err := browserType.Launch(launchOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	contextOptions := p.createContextOptions(config)

	context, err := browser.NewContext(contextOptions)
	if err != nil {
		if closeErr := browser.Close(); closeErr != nil {
			p.logger.Error("failed to close browser after context creation error", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("failed to create browser context: %w", err)
	}

	page, err := context.NewPage()
	if err != nil {
		if closeErr := context.Close(); closeErr != nil {
			p.logger.Error("failed to close context after page creation error", zap.Error(closeErr))
		}
		if closeErr := browser.Close(); closeErr != nil {
			p.logger.Error("failed to close browser after page creation error", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
	session := &BrowserSession{
		ID:       sessionID,
		Browser:  browser,
		Context:  context,
		Page:     page,
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	p.sessionsMux.Lock()
	p.sessions[sessionID] = session
	p.sessionsMux.Unlock()

	p.logger.Info("browser session created", zap.String("sessionID", sessionID))
	return session, nil
}

// CloseBrowser closes a browser session
func (p *playwrightImpl) CloseBrowser(ctx context.Context, sessionID string) error {
	p.sessionsMux.Lock()
	defer p.sessionsMux.Unlock()

	session, exists := p.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	if session.Context != nil {
		if err := session.Context.Close(); err != nil {
			p.logger.Error("failed to close context", zap.Error(err))
		}
	}
	if session.Browser != nil {
		if err := session.Browser.Close(); err != nil {
			p.logger.Error("failed to close browser", zap.Error(err))
		}
	}

	delete(p.sessions, sessionID)
	p.logger.Info("browser session closed", zap.String("sessionID", sessionID))
	return nil
}

// GetSession returns a browser session by ID
func (p *playwrightImpl) GetSession(sessionID string) (*BrowserSession, error) {
	p.sessionsMux.RLock()
	defer p.sessionsMux.RUnlock()

	session, exists := p.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	session.LastUsed = time.Now()
	return session, nil
}

// GetOrCreateDefaultSession gets the default shared session or creates it if it doesn't exist
func (p *playwrightImpl) GetOrCreateDefaultSession(ctx context.Context) (*BrowserSession, error) {
	p.sessionsMux.RLock()
	if session, exists := p.sessions[DefaultSessionID]; exists {
		session.LastUsed = time.Now()
		p.sessionsMux.RUnlock()
		p.logger.Debug("reusing existing default session", zap.String("sessionID", DefaultSessionID))
		return session, nil
	}
	p.sessionsMux.RUnlock()

	p.sessionsMux.Lock()
	defer p.sessionsMux.Unlock()

	if session, exists := p.sessions[DefaultSessionID]; exists {
		session.LastUsed = time.Now()
		p.logger.Debug("reusing existing default session (double-check)", zap.String("sessionID", DefaultSessionID))
		return session, nil
	}

	config := NewBrowserConfigFromConfig(p.config)
	p.logger.Info("creating new default browser session", zap.String("sessionID", DefaultSessionID))

	var browserType playwright.BrowserType
	switch config.Engine {
	case Chromium:
		browserType = p.pw.Chromium
	case Firefox:
		browserType = p.pw.Firefox
	case WebKit:
		browserType = p.pw.WebKit
	default:
		return nil, fmt.Errorf("unsupported browser engine: %s", config.Engine)
	}

	timeoutMs := float64(config.Timeout.Milliseconds())
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: &config.Headless,
		Args:     config.Args,
		Timeout:  &timeoutMs,
	}

	browser, err := browserType.Launch(launchOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	contextOptions := p.createContextOptions(config)

	context, err := browser.NewContext(contextOptions)
	if err != nil {
		if closeErr := browser.Close(); closeErr != nil {
			p.logger.Error("failed to close browser after context creation error", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("failed to create browser context: %w", err)
	}

	page, err := context.NewPage()
	if err != nil {
		if closeErr := context.Close(); closeErr != nil {
			p.logger.Error("failed to close context after page creation error", zap.Error(closeErr))
		}
		if closeErr := browser.Close(); closeErr != nil {
			p.logger.Error("failed to close browser after page creation error", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	session := &BrowserSession{
		ID:       DefaultSessionID,
		Browser:  browser,
		Context:  context,
		Page:     page,
		Created:  time.Now(),
		LastUsed: time.Now(),
	}

	p.sessions[DefaultSessionID] = session
	p.logger.Info("default browser session created successfully", zap.String("sessionID", DefaultSessionID))
	return session, nil
}

// NavigateToURL navigates to a URL in the specified session
func (p *playwrightImpl) NavigateToURL(ctx context.Context, sessionID, url string, waitUntil string, timeout time.Duration) error {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return err
	}

	var waitOption *playwright.WaitUntilState
	switch waitUntil {
	case "domcontentloaded":
		waitOption = playwright.WaitUntilStateDomcontentloaded
	case "networkidle":
		waitOption = playwright.WaitUntilStateNetworkidle
	default:
		waitOption = playwright.WaitUntilStateLoad
	}

	timeoutMs := float64(timeout.Milliseconds())
	options := playwright.PageGotoOptions{
		WaitUntil: waitOption,
		Timeout:   &timeoutMs,
	}

	p.logger.Info("navigating to URL", zap.String("sessionID", sessionID), zap.String("url", url))
	_, err = session.Page.Goto(url, options)
	return err
}

// ClickElement clicks an element in the specified session
func (p *playwrightImpl) ClickElement(ctx context.Context, sessionID, selector string, options map[string]any) error {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return err
	}

	clickOptions := playwright.PageClickOptions{}

	if timeout, ok := options["timeout"].(time.Duration); ok {
		timeoutMs := float64(timeout.Milliseconds())
		clickOptions.Timeout = &timeoutMs
	}
	if force, ok := options["force"].(bool); ok {
		clickOptions.Force = &force
	}
	if clickCount, ok := options["click_count"].(int); ok {
		clickOptions.ClickCount = &clickCount
	}
	if button, ok := options["button"].(string); ok {
		switch button {
		case "right":
			clickOptions.Button = playwright.MouseButtonRight
		case "middle":
			clickOptions.Button = playwright.MouseButtonMiddle
		default:
			clickOptions.Button = playwright.MouseButtonLeft
		}
	}

	p.logger.Info("clicking element", zap.String("sessionID", sessionID), zap.String("selector", selector))
	return session.Page.Locator(selector).Click(playwright.LocatorClickOptions{
		Timeout:    clickOptions.Timeout,
		Force:      clickOptions.Force,
		ClickCount: clickOptions.ClickCount,
		Button:     clickOptions.Button,
	})
}

// FillForm fills form fields in the specified session
func (p *playwrightImpl) FillForm(ctx context.Context, sessionID string, fields []map[string]any, submit bool, submitSelector string) error {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return err
	}

	p.logger.Info("filling form", zap.String("sessionID", sessionID), zap.Int("fields", len(fields)))

	for _, field := range fields {
		selector, ok := field["selector"].(string)
		if !ok {
			return fmt.Errorf("field selector is required")
		}

		value, ok := field["value"].(string)
		if !ok {
			return fmt.Errorf("field value is required")
		}

		fieldType, _ := field["type"].(string)

		switch fieldType {
		case "select":
			_, err = session.Page.Locator(selector).SelectOption(playwright.SelectOptionValues{Values: &[]string{value}}, playwright.LocatorSelectOptionOptions{})
		case "checkbox", "radio":
			if value == "true" || value == "1" {
				err = session.Page.Locator(selector).Check()
			} else {
				err = session.Page.Locator(selector).Uncheck()
			}
		default:
			err = session.Page.Locator(selector).Fill(value)
		}

		if err != nil {
			return fmt.Errorf("failed to fill field %s: %w", selector, err)
		}
	}

	if submit && submitSelector != "" {
		p.logger.Info("submitting form", zap.String("sessionID", sessionID), zap.String("submitSelector", submitSelector))
		err = session.Page.Locator(submitSelector).Click()
		if err != nil {
			return fmt.Errorf("failed to submit form: %w", err)
		}
	}

	return nil
}

// ExtractData extracts data from the page using selectors
func (p *playwrightImpl) ExtractData(ctx context.Context, sessionID string, extractors []map[string]any, format string) (string, error) {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return "", err
	}

	p.logger.Info("extracting data", zap.String("sessionID", sessionID), zap.Int("extractors", len(extractors)))

	results := make(map[string]any)

	for _, extractor := range extractors {
		name, ok := extractor["name"].(string)
		if !ok {
			return "", fmt.Errorf("extractor name is required")
		}

		selector, ok := extractor["selector"].(string)
		if !ok {
			return "", fmt.Errorf("extractor selector is required")
		}

		attribute, _ := extractor["attribute"].(string)
		if attribute == "" {
			attribute = "text"
		}

		multiple, _ := extractor["multiple"].(bool)

		if multiple {
			locator := session.Page.Locator(selector)
			count, err := locator.Count()
			if err != nil {
				return "", fmt.Errorf("failed to count elements for %s: %w", name, err)
			}

			var values []any
			for i := 0; i < count; i++ {
				elementLocator := locator.Nth(i)
				var value any
				if attribute == "text" {
					value, err = elementLocator.InnerText()
				} else {
					value, err = elementLocator.GetAttribute(attribute)
				}
				if err == nil {
					values = append(values, value)
				}
			}
			results[name] = values
		} else {
			locator := session.Page.Locator(selector)
			var value any
			if attribute == "text" {
				value, err = locator.InnerText()
			} else {
				value, err = locator.GetAttribute(attribute)
			}
			if err != nil {
				return "", fmt.Errorf("failed to extract %s: %w", name, err)
			}
			results[name] = value
		}
	}

	switch format {
	case "json":
		return fmt.Sprintf("%+v", results), nil
	case "csv":
		return fmt.Sprintf("%+v", results), nil
	default:
		return fmt.Sprintf("%+v", results), nil
	}
}

// TakeScreenshot captures a screenshot
func (p *playwrightImpl) TakeScreenshot(ctx context.Context, sessionID, path string, fullPage bool, selector string, format string, quality int) error {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return err
	}

	p.logger.Info("taking screenshot", zap.String("sessionID", sessionID), zap.String("path", path))

	options := playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(fullPage),
	}

	if format == "jpeg" {
		options.Type = playwright.ScreenshotTypeJpeg
		options.Quality = playwright.Int(quality)
	} else {
		options.Type = playwright.ScreenshotTypePng
	}

	if selector != "" {
		locator := session.Page.Locator(selector)
		_, err = locator.Screenshot(playwright.LocatorScreenshotOptions{
			Path: &path,
			Type: options.Type,
		})
		return err
	}

	_, err = session.Page.Screenshot(options)
	return err
}

// ExecuteScript executes JavaScript in the browser context
func (p *playwrightImpl) ExecuteScript(ctx context.Context, sessionID, script string, args []any) (any, error) {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	p.logger.Info("executing script", zap.String("sessionID", sessionID))

	result, err := session.Page.Evaluate(script, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute script: %w", err)
	}

	return result, nil
}

// WaitForCondition waits for specific conditions
func (p *playwrightImpl) WaitForCondition(ctx context.Context, sessionID, condition, selector, state string, timeout time.Duration, customFunction string) error {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return err
	}

	p.logger.Info("waiting for condition", zap.String("sessionID", sessionID), zap.String("condition", condition))

	switch condition {
	case "selector":
		var waitState *playwright.WaitForSelectorState
		switch state {
		case "hidden":
			waitState = playwright.WaitForSelectorStateHidden
		case "attached":
			waitState = playwright.WaitForSelectorStateAttached
		case "detached":
			waitState = playwright.WaitForSelectorStateDetached
		default:
			waitState = playwright.WaitForSelectorStateVisible
		}

		timeoutMs := float64(timeout.Milliseconds())
		options := playwright.LocatorWaitForOptions{
			State:   waitState,
			Timeout: &timeoutMs,
		}

		err = session.Page.Locator(selector).WaitFor(options)
		return err

	case "navigation":
		time.Sleep(timeout)
		return nil

	case "function":
		if customFunction == "" {
			return fmt.Errorf("custom function is required for function condition")
		}

		timeoutMs := float64(timeout.Milliseconds())
		options := playwright.PageWaitForFunctionOptions{
			Timeout: &timeoutMs,
		}

		_, err = session.Page.WaitForFunction(customFunction, options)
		return err

	case "timeout":
		time.Sleep(timeout)
		return nil

	default:
		return fmt.Errorf("unsupported condition type: %s", condition)
	}
}

// HandleAuthentication handles various authentication scenarios
func (p *playwrightImpl) HandleAuthentication(ctx context.Context, sessionID, authType, username, password, loginURL string, selectors map[string]string) error {
	session, err := p.GetSession(sessionID)
	if err != nil {
		return err
	}

	p.logger.Info("handling authentication", zap.String("sessionID", sessionID), zap.String("type", authType))

	switch authType {
	case "basic":
		if loginURL != "" {
			_, err = session.Page.Goto(loginURL)
			return err
		}
		return fmt.Errorf("basic auth requires loginURL")

	case "form":
		if loginURL != "" {
			_, err = session.Page.Goto(loginURL)
			if err != nil {
				return fmt.Errorf("failed to navigate to login URL: %w", err)
			}
		}

		if usernameSelector, ok := selectors["username_selector"]; ok && usernameSelector != "" {
			err = session.Page.Locator(usernameSelector).Fill(username)
			if err != nil {
				return fmt.Errorf("failed to fill username: %w", err)
			}
		}

		if passwordSelector, ok := selectors["password_selector"]; ok && passwordSelector != "" {
			err = session.Page.Locator(passwordSelector).Fill(password)
			if err != nil {
				return fmt.Errorf("failed to fill password: %w", err)
			}
		}

		if submitSelector, ok := selectors["submit_selector"]; ok && submitSelector != "" {
			err = session.Page.Locator(submitSelector).Click()
			if err != nil {
				return fmt.Errorf("failed to submit form: %w", err)
			}
		}

		return nil

	case "oauth":
		if loginURL != "" {
			_, err = session.Page.Goto(loginURL)
			return err
		}
		return fmt.Errorf("OAuth implementation not yet supported")

	default:
		return fmt.Errorf("unsupported authentication type: %s", authType)
	}
}

// GetHealth checks the health of the service
func (p *playwrightImpl) GetHealth(ctx context.Context) error {
	if p.pw == nil {
		return fmt.Errorf("playwright not initialized")
	}

	p.sessionsMux.RLock()
	activeSessions := len(p.sessions)
	p.sessionsMux.RUnlock()

	p.logger.Info("playwright service health check", zap.Int("activeSessions", activeSessions))
	return nil
}

// Shutdown gracefully shuts down the service
func (p *playwrightImpl) Shutdown(ctx context.Context) error {
	p.logger.Info("shutting down playwright service")

	p.sessionsMux.Lock()
	for sessionID := range p.sessions {
		if session := p.sessions[sessionID]; session != nil {
			if session.Context != nil {
				if err := session.Context.Close(); err != nil {
					p.logger.Error("failed to close context during shutdown", zap.Error(err))
				}
			}
			if session.Browser != nil {
				if err := session.Browser.Close(); err != nil {
					p.logger.Error("failed to close browser during shutdown", zap.Error(err))
				}
			}
		}
	}
	p.sessions = make(map[string]*BrowserSession)
	p.sessionsMux.Unlock()

	if p.pw != nil {
		err := p.pw.Stop()
		if err != nil {
			p.logger.Error("failed to stop playwright", zap.Error(err))
			return err
		}
	}

	p.logger.Info("playwright service shutdown complete")
	return nil
}

// GetConfig returns the configuration
func (p *playwrightImpl) GetConfig() *config.Config {
	return p.config
}

// createContextOptions creates browser context options from configuration
func (p *playwrightImpl) createContextOptions(browserConfig *BrowserConfig) playwright.BrowserNewContextOptions {
	return playwright.BrowserNewContextOptions{
		Viewport: &playwright.Size{
			Width:  browserConfig.ViewportWidth,
			Height: browserConfig.ViewportHeight,
		},
		UserAgent: playwright.String(p.config.Browser.UserAgent),
		ExtraHttpHeaders: map[string]string{
			"Accept":                    p.config.Browser.HeaderAccept,
			"Accept-Language":           p.config.Browser.HeaderAcceptLanguage,
			"Accept-Encoding":           p.config.Browser.HeaderAcceptEncoding,
			"DNT":                       p.config.Browser.HeaderDnt,
			"Connection":                p.config.Browser.HeaderConnection,
			"Upgrade-Insecure-Requests": p.config.Browser.HeaderUpgradeInsecureRequests,
		},
		JavaScriptEnabled: playwright.Bool(true),
		BypassCSP:         playwright.Bool(true),
	}
}
