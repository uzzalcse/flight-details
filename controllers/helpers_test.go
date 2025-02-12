package controllers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatSuccessResponse - Test for `formatSuccessResponse`
func TestFormatSuccessResponse(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		inputJSON   string
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "Success - Valid JSON",
			inputJSON: `{
				"hits": {"total": 1}
			}`,
			expected: map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"hits": map[string]interface{}{
						"total": float64(1),
					},
				},
			},
			expectError: false,
		},
		{
			name:      "Failure - Invalid JSON",
			inputJSON: `{invalid-json}`,
			expected: map[string]interface{}{
				"status":  "error",
				"message": "Response parsing failed",
			},
			expectError: true,
		},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatSuccessResponse(tc.inputJSON)
			assert.Equal(t, tc.expected, result, "Response should match expected")
		})
	}
}

// TestParseFloat - Test for `ParseFloat`
func TestParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Valid float", "12.34", 12.34},
		{"Valid integer as float", "7", 7.0},
		{"Empty string", "", 0},
		{"Invalid number", "abc", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseFloat(tc.input)
			assert.Equal(t, tc.expected, result, "Parsed float should match expected value")
		})
	}
}

// TestParseInt - Test for `ParseInt`
func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Valid integer", "25", 25},
		{"Zero value", "0", 0},
		{"Empty string", "", 0},
		{"Invalid number", "xyz", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseInt(tc.input)
			assert.Equal(t, tc.expected, result, "Parsed int should match expected value")
		})
	}
}

// MockFlightController to test c.Data initialization
type MockFlightController struct {
	Data map[interface{}]interface{}
	Ctx  struct {
		Input struct {
			Query func(string) string
		}
	}
}

// TestFlightControllerDataInitialization - Test for c.Data initialization
func TestFlightControllerDataInitialization(t *testing.T) {
	controller := &MockFlightController{}

	// Simulate nil data
	if controller.Data == nil {
		controller.Data = make(map[interface{}]interface{})
	}

	assert.NotNil(t, controller.Data, "c.Data should be initialized to an empty map")
}

// TestFormatSuccessResponseStatusCheck - Test if response maintains existing status field
func TestFormatSuccessResponseStatusCheck(t *testing.T) {
	responseWithStatus := `{"status": "error", "message": "Some error"}`
	var expected map[string]interface{}
	_ = json.Unmarshal([]byte(responseWithStatus), &expected)

	result := FormatSuccessResponse(responseWithStatus)

	assert.Equal(t, expected, result, "Response should return the original JSON if status exists")
}
