package tools

import (
	"encoding/json"
	"fmt"
	"slices"
)

// Shared defaults and bounds used across the manually-implemented tools.
// Promoted from per-file magic numbers so behavior stays consistent.
const (
	defaultTimeoutMs   = 30000
	minTimeoutMs       = 1
	maxTimeoutMs       = 600000
	defaultJPEGQuality = 80
	minJPEGQuality     = 0
	maxJPEGQuality     = 100
	defaultClickCount  = 1
	minClickCount      = 1
	maxClickCount      = 10
	maxScriptSize      = 50000

	screenshotSelectorMaxRunes = 20
)

// requiredString returns args[key] as a non-empty string. Returns an error
// if the key is absent, the value is not a string, or the string is empty.
func requiredString(args map[string]any, key string) (string, error) {
	raw, ok := args[key]
	if !ok {
		return "", fmt.Errorf("%s parameter is required", key)
	}
	s, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string, got %T", key, raw)
	}
	if s == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return s, nil
}

// stringArg returns args[key] as a string, or defaultValue if absent or empty.
// Returns an error if the key is present but the value is not a string.
func stringArg(args map[string]any, key, defaultValue string) (string, error) {
	raw, ok := args[key]
	if !ok {
		return defaultValue, nil
	}
	s, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string, got %T", key, raw)
	}
	if s == "" {
		return defaultValue, nil
	}
	return s, nil
}

// boolArg returns args[key] as a bool, or defaultValue if absent.
// Returns an error if the key is present but the value is not a bool.
func boolArg(args map[string]any, key string, defaultValue bool) (bool, error) {
	raw, ok := args[key]
	if !ok {
		return defaultValue, nil
	}
	b, ok := raw.(bool)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean, got %T", key, raw)
	}
	return b, nil
}

// intArg returns args[key] as an int, or defaultValue if absent. Accepts
// int, int64, and float64 (since JSON unmarshaling produces float64).
// Returns an error if the value cannot be converted to a whole int.
func intArg(args map[string]any, key string, defaultValue int) (int, error) {
	raw, ok := args[key]
	if !ok {
		return defaultValue, nil
	}
	switch v := raw.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		if v != float64(int(v)) {
			return 0, fmt.Errorf("%s must be a whole integer, got %v", key, v)
		}
		return int(v), nil
	default:
		return 0, fmt.Errorf("%s must be an integer, got %T", key, raw)
	}
}

// boundedIntArg returns args[key] as an int validated to be in
// [minInclusive, maxInclusive]. Returns defaultValue if absent.
func boundedIntArg(args map[string]any, key string, defaultValue, minInclusive, maxInclusive int) (int, error) {
	v, err := intArg(args, key, defaultValue)
	if err != nil {
		return 0, err
	}
	if v < minInclusive || v > maxInclusive {
		return 0, fmt.Errorf("%s must be between %d and %d, got %d", key, minInclusive, maxInclusive, v)
	}
	return v, nil
}

// sliceArg returns args[key] as []any. The second return value reports
// whether the key was present (so callers can distinguish "missing" from
// "present-but-empty"). Returns an error if the key is present but the
// value is not an array.
func sliceArg(args map[string]any, key string) ([]any, bool, error) {
	raw, ok := args[key]
	if !ok {
		return nil, false, nil
	}
	s, ok := raw.([]any)
	if !ok {
		return nil, true, fmt.Errorf("%s must be an array, got %T", key, raw)
	}
	return s, true, nil
}

// marshalResponse encodes a tool response as JSON. Centralized so we get
// consistent error messages and avoid scattering identical
// json.Marshal+error-wrap boilerplate across every tool.
func marshalResponse(response any) (string, error) {
	b, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}
	return string(b), nil
}

// oneOf reports whether candidate matches any of allowed. Centralized so
// each tool's per-field validator can collapse to a single call.
func oneOf(candidate string, allowed ...string) bool {
	return slices.Contains(allowed, candidate)
}
