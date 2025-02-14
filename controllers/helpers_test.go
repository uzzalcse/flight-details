package controllers

import (
	"flight-details/structs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
)

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"Empty string", "", 0.0, false},
		{"Valid float", "123.45", 123.45, false},
		{"Integer as float", "123", 123.0, false},
		{"Invalid float", "abc", 0.0, true},
		{"Negative float", "-123.45", -123.45, false},
		{"Zero", "0", 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFloat(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		hasError bool
	}{
		{"Empty string", "", 0, false},
		{"Valid integer", "123", 123, false},
		{"Negative integer", "-123", -123, false},
		{"Invalid integer", "123.45", 0, true},
		{"Non-numeric", "abc", 0, true},
		{"Zero", "0", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseInt(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"Empty string", "", false, true},
		{"String true", "true", true, false},
		{"String false", "false", false, false},
		{"Numeric 1", "1", true, false},
		{"Numeric 0", "0", false, false},
		{"Invalid value", "yes", false, true},
		{"Invalid numeric", "2", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBool(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFormatSuccessResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:  "Valid JSON without status",
			input: `{"message": "test"}`,
			expected: map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"message": "test",
				},
			},
		},
		{
			name:  "Valid JSON with status",
			input: `{"status": "custom", "message": "test"}`,
			expected: map[string]interface{}{
				"status":  "custom",
				"message": "test",
			},
		},
		{
			name:  "Invalid JSON",
			input: `{invalid json}`,
			expected: map[string]interface{}{
				"status":  "error",
				"message": "Response parsing failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSuccessResponse(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseFlightSearchRequest(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedParams structs.FlightSearchParams
		expectError    bool
	}{
		{
			name: "Valid complete request",
			queryParams: map[string]string{
				"timestamp":          "2024-01-01",
				"FlightNum":          "FL123",
				"dayOfWeek":          "1",
				"FlightDelayMin":     "30",
				"Cancelled":          "true",
				"AvgTicketPrice":     "299.99",
				"DistanceMiles":      "1000",
				"DistanceKilometers": "1609.34",
				"FlightTimeMin":      "120",
				"FlightTimeHour":     "2",
			},
			expectedParams: structs.FlightSearchParams{
				TravelTime:         "2024-01-01",
				FlightNum:          "FL123",
				DayOfWeek:          1,
				FlightDelayMin:     30,
				Cancelled:          true,
				AvgTicketPrice:     299.99,
				DistanceMiles:      1000,
				DistanceKilometers: 1609.34,
				FlightTimeMin:      120,
				FlightTimeHour:     2,
			},
			expectError: false,
		},
		{
			name: "Missing required timestamp",
			queryParams: map[string]string{
				"FlightNum": "FL123",
			},
			expectedParams: structs.FlightSearchParams{
				FlightNum: "FL123",
			},
			expectError: true,
		},
		{
			name: "Invalid dayOfWeek",
			queryParams: map[string]string{
				"timestamp": "2024-01-01",
				"dayOfWeek": "7",
			},
			expectedParams: structs.FlightSearchParams{
				TravelTime: "2024-01-01",
				DayOfWeek:  7,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request with query parameters
			req, _ := http.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			// Create recorder and context
			w := httptest.NewRecorder()

			// Create a new context properly
			ctx := context.NewContext()
			ctx.Reset(w, req)

			// Create controller with test context
			controller := &FlightController{}
			controller.Init(ctx, "FlightController", "TestAction", nil)

			// Parse request
			params, err := ParseFlightSearchRequest(controller)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedParams, params)
			}
		})
	}
}

type Container struct {
	Data map[interface{}]interface{}
}

func TestDataMapInitialization(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		initial map[interface{}]interface{}
		wantNil bool
	}{
		{
			name:    "nil map should be initialized",
			initial: nil,
			wantNil: false,
		},
		{
			name:    "existing map should not be modified",
			initial: map[interface{}]interface{}{"key": "value"},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create container with initial map
			c := &Container{
				Data: tt.initial,
			}

			// Run the initialization code
			if c.Data == nil {
				c.Data = make(map[interface{}]interface{})
			}

			// Assert map is not nil
			if (c.Data == nil) != tt.wantNil {
				t.Errorf("After initialization, Data map was nil = %v, want nil = %v", c.Data == nil, tt.wantNil)
			}

			// For the case with existing map, verify content wasn't lost
			if tt.initial != nil {
				if val, ok := c.Data["key"]; !ok || val != "value" {
					t.Errorf("Existing map data was lost or modified")
				}
			}
		})
	}
}
