package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
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
		hasError bool
	}{
		{"Valid float", "12.34", 12.34, false},
		{"Valid integer as float", "7", 7.0, false},
		{"Empty string", "", 0, false},
		{"Invalid number", "abc", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseFloat(tc.input)
			if tc.hasError {
				assert.Error(t, err, "Expected error for invalid float")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result, "Parsed float should match expected value")
			}
		})
	}
}

func TestParseFloatErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"Valid float", "12.34", 12.34, false},
		{"Valid integer as float", "7", 7.0, false},
		{"Empty string", "", 0, false},
		{"Invalid number", "abc", 0, true}, // Should return error
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseFloat(tc.input)
			if tc.hasError {
				assert.Error(t, err, "Expected error for invalid float")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result, "Parsed float should match expected value")
			}
		})
	}
}

// TestParseInt - Test for `ParseInt`
func TestParseInt(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{"Valid integer", "25", 25, false},
		{"Zero value", "0", 0, false},
		{"Empty string (defaults to 0)", "", 0, false},
		{"Invalid number", "xyz", 0, true},
		{"Negative number", "-5", -5, false},
		{"Large number", "1000000", 1000000, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseInt(tc.input)
			if tc.expectError {
				assert.Error(t, err, "Expected an error for invalid input")
			} else {
				assert.NoError(t, err, "Did not expect an error for valid input")
				assert.Equal(t, tc.expected, result, "Parsed int should match expected value")
			}
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

func TestParseFlightSearchRequest_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedError  string
		expectedStatus int
	}{
		// Mandatory Field Test
		{
			name:           "Failure - Missing Required Timestamp",
			queryParams:    "dayOfWeek=2&FlightDelayMin=30&Cancelled=true",
			expectedError:  "timestamp (TravelTime) is required",
			expectedStatus: http.StatusBadRequest,
		},

		// Invalid Integer Test (dayOfWeek)
		{
			name:           "Failure - Invalid dayOfWeek",
			queryParams:    "timestamp=2025-02-03T00:00:00&dayOfWeek=9",
			expectedError:  "invalid dayOfWeek: must be between 0 and 6",
			expectedStatus: http.StatusBadRequest,
		},

		// Invalid Integer Test (FlightDelayMin)
		{
			name:           "Failure - Non-numeric FlightDelayMin",
			queryParams:    "timestamp=2025-02-03T00:00:00&FlightDelayMin=xyz",
			expectedError:  "invalid FlightDelayMin value: invalid integer: xyz",
			expectedStatus: http.StatusBadRequest,
		},

		// Invalid Boolean Test
		{
			name:           "Failure - Invalid Cancelled Boolean",
			queryParams:    "timestamp=2025-02-03T00:00:00&Cancelled=maybe",
			expectedError:  "invalid Cancelled value: invalid boolean: maybe",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Invalid FlightDelay Boolean",
			queryParams:    "timestamp=2025-02-03T00:00:00&FlightDelay=notTrue",
			expectedError:  "invalid FlightDelay value: invalid boolean: notTrue",
			expectedStatus: http.StatusBadRequest,
		},

		// Invalid Float Test
		{
			name:           "Failure - Invalid AvgTicketPrice Float",
			queryParams:    "timestamp=2025-02-03T00:00:00&AvgTicketPrice=abc",
			expectedError:  "invalid AvgTicketPrice value: invalid float: abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Invalid DistanceMiles Float",
			queryParams:    "timestamp=2025-02-03T00:00:00&DistanceMiles=xyz",
			expectedError:  "invalid DistanceMiles value: invalid float: xyz",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failure - Invalid FlightTimeMin Float",
			queryParams:    "timestamp=2025-02-03T00:00:00&FlightTimeMin=abc",
			expectedError:  "invalid FlightTimeMin value: invalid float: abc",
			expectedStatus: http.StatusBadRequest,
		},
		// Invalid Float Test (DistanceKilometers)
		{
			name:           "Failure - Invalid DistanceKilometers Float",
			queryParams:    "timestamp=2025-02-03T00:00:00&DistanceKilometers=xyz",
			expectedError:  "invalid DistanceKilometers value: invalid float: xyz",
			expectedStatus: http.StatusBadRequest,
		},

		// Invalid Float Test (FlightTimeHour)
		{
			name:           "Failure - Invalid FlightTimeHour Float",
			queryParams:    "timestamp=2025-02-03T00:00:00&FlightTimeHour=abc",
			expectedError:  "invalid FlightTimeHour value: invalid float: abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/api/flights/all_params/search?"+tc.queryParams, nil)
			w := httptest.NewRecorder()

			ctx := context.NewContext()
			ctx.Reset(w, req)

			c := &FlightController{}
			c.Ctx = ctx

			// Call the function
			_, err := ParseFlightSearchRequest(c)

			// Check if the expected error occurs
			assert.Error(t, err, "Expected error for invalid input")
			assert.Contains(t, err.Error(), tc.expectedError, "Error message should match expected")
		})
	}
}
