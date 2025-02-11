package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"flight-details/controllers"
	"flight-details/services"
	"flight-details/structs"

	"bou.ke/monkey"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
)

// TestFlightController_GetByAllParams - Table-Driven Test (TDT)
func TestFlightController_GetByAllParams(t *testing.T) {
	// Define test cases
	tests := []struct {
		name             string
		mockResponse     string
		mockError        error
		queryParams      string
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:           "Success - Valid Request",
			mockResponse:   `{"hits":{"total":1,"hits":[]}}`,
			mockError:      nil,
			queryParams:    "timestamp=2025-02-03T00:00:00&FlightNum=9HY9SWR",
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"hits": map[string]interface{}{
						"total": float64(1),
						"hits":  []interface{}{}, // ✅ Add hits array to match the actual response
					},
				},
			},
		},
		{
			name:             "Failure - Missing Required Timestamp",
			mockResponse:     "",
			mockError:        errors.New("timestamp (TravelTime) is required"),
			queryParams:      "FlightNum=9HY9SWR",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: nil, // Expecting a failure response
		},
		{
			name:             "Failure - Service Returns an Error",
			mockResponse:     "",
			mockError:        errors.New("Elasticsearch error"),
			queryParams:      "timestamp=2025-02-03T00:00:00&FlightNum=9HY9SWR",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil, // Expecting a failure response
		},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ✅ Mock `services.SearchFlights`
			monkey.Patch(services.SearchFlights, func(params structs.FlightSearchParams) (string, error) {
				return tc.mockResponse, tc.mockError
			})
			defer monkey.Unpatch(services.SearchFlights)

			// ✅ Create a mock HTTP request
			req := httptest.NewRequest(http.MethodGet, "/v1/api/flights/all_params/search?"+tc.queryParams, nil)
			w := httptest.NewRecorder()

			// ✅ Correctly initialize Beego context
			ctx := context.NewContext()
			ctx.Reset(w, req)

			// ✅ Create the FlightController instance
			c := &controllers.FlightController{}
			c.Ctx = ctx // Attach Beego context

			// ✅ Call the controller method
			c.GetByAllParams()

			// ✅ Debugging Output
			t.Logf("Test: %s | Response Status: %d | Response Body: %s", tc.name, w.Code, w.Body.String())

			// ✅ Check HTTP status code
			assert.Equal(t, tc.expectedStatus, w.Code, "HTTP status code should match")

			// ✅ Parse JSON response if expectedResponse is not nil
			if tc.expectedResponse != nil {
				var body map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &body)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.Equal(t, tc.expectedResponse, body, "Response body should match expected")
			}
		})
	}
}
