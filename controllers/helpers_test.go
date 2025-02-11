package controllers

import (
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
			name: "✅ Success - Valid JSON",
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
			name:      "❌ Failure - Invalid JSON",
			inputJSON: `{invalid-json}`,
			expected: map[string]interface{}{
				"status":  "error",
				"message": "Response parsing failed",
			},
			expectError: true,
		},
	}

	// ✅ Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ✅ Call function under test
			result := FormatSuccessResponse(tc.inputJSON)

			// ✅ Assertions
			assert.Equal(t, tc.expected, result, "Response should match expected")
		})
	}
}
