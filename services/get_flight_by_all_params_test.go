package services_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"flight-details/services"
	"flight-details/structs"
	"flight-details/utils"

	"bou.ke/monkey"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
)

// TestSearchFlights - Table-Driven Test (TDT) for SearchFlights
func TestSearchFlights(t *testing.T) {
	// Define test cases
	tests := []struct {
		name             string
		mockResponse     string
		mockError        error
		inputParams      structs.FlightSearchParams
		expectedSuccess  bool
		expectedResponse map[string]interface{}
	}{
		{
			name:         "Success - Valid Request",
			mockResponse: `{"hits":{"total":1,"hits":[]}}`,
			mockError:    nil,
			inputParams: structs.FlightSearchParams{
				TravelTime:    "2025-02-03T00:00:00",
				FlightNum:     "9HY9SWR",
				DestCountry:   "AU",
				OriginWeather: "Sunny",
			},
			expectedSuccess: true,
			expectedResponse: map[string]interface{}{
				"hits": map[string]interface{}{
					"total": float64(1),
					"hits":  []interface{}{},
				},
			},
		},
		{
			name:         "Success - Triggers addExactRangeQuery",
			mockResponse: `{"hits":{"total":2,"hits":[]}}`,
			mockError:    nil,
			inputParams: structs.FlightSearchParams{
				TravelTime:         "2025-02-03T00:00:00",
				AvgTicketPrice:     500.75,  // ✅ Triggers addExactRangeQuery
				DistanceMiles:      1500.55, // ✅ Triggers addExactRangeQuery
				DistanceKilometers: 2414.56, // ✅ Triggers addExactRangeQuery
				FlightTimeMin:      120,     // ✅ Triggers addExactRangeQuery
				FlightTimeHour:     2,       // ✅ Triggers addExactRangeQuery
			},
			expectedSuccess: true,
			expectedResponse: map[string]interface{}{
				"hits": map[string]interface{}{
					"total": float64(2),
					"hits":  []interface{}{},
				},
			},
		},
		{
			name:         "Success - Triggers addGeoLocQuery",
			mockResponse: `{"hits":{"total":3,"hits":[]}}`,
			mockError:    nil,
			inputParams: structs.FlightSearchParams{
				TravelTime:        "2025-02-03T00:00:00",
				OriginLocationLat: "40.7128",  // ✅ Triggers addGeoLocQuery
				OriginLocationLon: "-74.0060", // ✅ Triggers addGeoLocQuery
			},
			expectedSuccess: true,
			expectedResponse: map[string]interface{}{
				"hits": map[string]interface{}{
					"total": float64(3),
					"hits":  []interface{}{},
				},
			},
		},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ✅ Mock Elasticsearch Client
			mockESClient := &elasticsearch.Client{}

			monkey.Patch(utils.GetElasticClient, func() *utils.ESClient {
				return &utils.ESClient{Client: mockESClient}
			})
			defer monkey.Unpatch(utils.GetElasticClient)

			// ✅ Fully Mock Elasticsearch API Calls
			monkey.Patch((*elasticsearch.Client).Perform, func(_ *elasticsearch.Client, req *http.Request) (*http.Response, error) {
				if tc.mockError != nil {
					// ✅ Return a `nil` response + controlled error
					return nil, tc.mockError
				}

				// ✅ Mock successful Elasticsearch response
				resp := &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(tc.mockResponse)),
				}
				return resp, nil
			})
			defer monkey.Unpatch((*elasticsearch.Client).Perform)

			// ✅ Debugging Output to Ensure Patch is Applied
			t.Logf("Running Test: %s | Expected Success: %v", tc.name, tc.expectedSuccess)

			// Call the actual SearchFlights function
			response, err := services.SearchFlights(tc.inputParams)

			// ✅ Debugging Output to Check Results
			t.Logf("Test: %s | Response: %s | Error: %v", tc.name, response, err)

			// Check success/failure based on expectations
			if tc.expectedSuccess {
				assert.NoError(t, err, "Expected no error")
				var body map[string]interface{}
				err = json.Unmarshal([]byte(response), &body)
				assert.NoError(t, err, "Expected valid JSON response")
				assert.Equal(t, tc.expectedResponse, body, "Response body should match expected")
			} else {
				assert.Error(t, err, "Expected an error")
			}
		})
	}
}
