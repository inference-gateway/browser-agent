package skills

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	server "github.com/inference-gateway/adk/server"
	types "github.com/inference-gateway/adk/types"
	"github.com/inference-gateway/browser-agent/internal/playwright"
	"github.com/inference-gateway/browser-agent/internal/playwright/mocks"
	zap "go.uber.org/zap"
)

func TestNewScrollSkill(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}

	skill := NewScrollSkill(logger, mockPlaywright)

	if skill == nil {
		t.Errorf("Expected skill to be created, got nil")
	}
}

func TestScrollHandler_ValidPageScroll(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	skill := &ScrollSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	ctx := context.WithValue(context.Background(), server.TaskContextKey, &types.Task{ID: "test-task"})

	session := &playwright.BrowserSession{ID: "test-session"}
	mockPlaywright.GetOrCreateTaskSessionReturns(session, nil)
	mockPlaywright.ScrollReturns(nil)

	args := map[string]any{
		"target":    "page",
		"direction": "top",
	}

	result, err := skill.ScrollHandler(ctx, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result")
	}

	var response map[string]any
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Errorf("Expected no unmarshal error, got %v", err)
	}
	if !response["success"].(bool) {
		t.Errorf("Expected success to be true")
	}
	if response["target"] != "page" {
		t.Errorf("Expected target to be 'page', got %v", response["target"])
	}
	if response["direction"] != "top" {
		t.Errorf("Expected direction to be 'top', got %v", response["direction"])
	}

	if mockPlaywright.GetOrCreateTaskSessionCallCount() != 1 {
		t.Errorf("Expected GetOrCreateTaskSession to be called once")
	}
	if mockPlaywright.ScrollCallCount() != 1 {
		t.Errorf("Expected Scroll to be called once")
	}
}

func TestScrollHandler_MissingTarget(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	skill := &ScrollSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	ctx := context.Background()
	args := map[string]any{}

	result, err := skill.ScrollHandler(ctx, args)

	if err == nil {
		t.Errorf("Expected error for missing target")
	}
	if result != "" {
		t.Errorf("Expected empty result on error")
	}
	if err.Error() != "target parameter is required and must be a non-empty string" {
		t.Errorf("Expected specific error message, got %v", err.Error())
	}
}

func TestScrollHandler_InvalidTarget(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	skill := &ScrollSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	ctx := context.Background()
	args := map[string]any{
		"target": "invalid",
	}

	result, err := skill.ScrollHandler(ctx, args)

	if err == nil {
		t.Errorf("Expected error for invalid target")
	}
	if result != "" {
		t.Errorf("Expected empty result on error")
	}
}

func TestScrollHandler_ElementTargetMissingSelector(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	skill := &ScrollSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	ctx := context.Background()
	args := map[string]any{
		"target": "element",
	}

	result, err := skill.ScrollHandler(ctx, args)

	if err == nil {
		t.Errorf("Expected error for missing selector")
	}
	if result != "" {
		t.Errorf("Expected empty result on error")
	}
}

func TestScrollHandler_SessionError(t *testing.T) {
	logger := zap.NewNop()
	mockPlaywright := &mocks.FakeBrowserAutomation{}
	skill := &ScrollSkill{
		logger:     logger,
		playwright: mockPlaywright,
	}

	ctx := context.WithValue(context.Background(), server.TaskContextKey, &types.Task{ID: "test-task"})

	mockPlaywright.GetOrCreateTaskSessionReturns(nil, errors.New("session error"))

	args := map[string]any{
		"target": "page",
	}

	result, err := skill.ScrollHandler(ctx, args)

	if err == nil {
		t.Errorf("Expected error for session failure")
	}
	if result != "" {
		t.Errorf("Expected empty result on error")
	}

	if mockPlaywright.GetOrCreateTaskSessionCallCount() != 1 {
		t.Errorf("Expected GetOrCreateTaskSession to be called once")
	}
}

func TestIsValidTarget(t *testing.T) {
	skill := &ScrollSkill{}

	tests := []struct {
		target   string
		expected bool
	}{
		{"page", true},
		{"element", true},
		{"coordinates", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := skill.isValidTarget(test.target)
		if result != test.expected {
			t.Errorf("For target '%s', expected %v, got %v", test.target, test.expected, result)
		}
	}
}

func TestIsValidBehavior(t *testing.T) {
	skill := &ScrollSkill{}

	tests := []struct {
		behavior string
		expected bool
	}{
		{"smooth", true},
		{"instant", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := skill.isValidBehavior(test.behavior)
		if result != test.expected {
			t.Errorf("For behavior '%s', expected %v, got %v", test.behavior, test.expected, result)
		}
	}
}

func TestIsValidDirection(t *testing.T) {
	skill := &ScrollSkill{}

	tests := []struct {
		direction string
		expected  bool
	}{
		{"up", true},
		{"down", true},
		{"left", true},
		{"right", true},
		{"top", true},
		{"bottom", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := skill.isValidDirection(test.direction)
		if result != test.expected {
			t.Errorf("For direction '%s', expected %v, got %v", test.direction, test.expected, result)
		}
	}
}
