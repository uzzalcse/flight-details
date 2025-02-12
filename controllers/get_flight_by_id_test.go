package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"flight-details/utils"
	"net/http/httptest"

	beecontext "github.com/beego/beego/v2/server/web/context"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
)

// MockTransport implements http.RoundTripper interface
type MockTransport struct {
	Response    string
	StatusCode  int
	ErrorString string
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ErrorString != "" {
		return nil, &mockError{t.ErrorString}
	}

	header := http.Header{}
	header.Add("X-Elastic-Product", "Elasticsearch")
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode: t.StatusCode,
		Body:       io.NopCloser(strings.NewReader(t.Response)),
		Header:     header,
		Request:    req,
	}, nil
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

type testCase struct {
	name           string
	flightID       string
	response       string
	statusCode     int
	errorString    string
	expectedStatus int
	expectError    bool
	validateResp   bool
}

func setupTestController(t *testing.T, mockTransport *MockTransport) (*FlightController, *httptest.ResponseRecorder) {
	// Create Elasticsearch client with mock transport
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Transport: mockTransport,
	})
	assert.NoError(t, err)

	// Create wrapped ESClient
	wrappedClient := &utils.ESClient{
		Client: esClient,
	}

	// Create controller
	controller := &FlightController{}

	// Create test context
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/api/", nil)
	context := beecontext.NewContext()
	context.Reset(w, r)

	// Initialize controller
	controller.Init(context, "FlightController", "GetFlightDetails", nil)
	
	// Set the ESClient after initialization
	controller.esClient = wrappedClient

	return controller, w
}

func TestGetFlightDetails(t *testing.T) {
	tests := []testCase{
		{
			name:           "Missing flight ID",
			flightID:       "",
			response:       "",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			expectError:    true,
			validateResp:   false,
		},
		{
			name:           "Elasticsearch connection error",
			flightID:       "flight123",
			errorString:    "connection refused",
			expectedStatus: http.StatusOK,
			expectError:    true,
			validateResp:   false,
		},
		{
			name:           "Invalid JSON response",
			flightID:       "flight123",
			response:       "invalid json",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			expectError:    true,
			validateResp:   false,
		},
		{
			name:     "Successful flight details retrieval",
			flightID: "flight123",
			response: `{
				"took": 1,
				"hits": {
					"total": {"value": 1, "relation": "eq"},
					"hits": [{
						"_id": "flight123",
						"_source": {
							"FlightNum": "XY123",
							"Origin": "LAX",
							"Destination": "JFK"
						}
					}]
				}
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResp:   true,
		},
		{
			name:     "No matching flights found",
			flightID: "nonexistent",
			response: `{
				"took": 1,
				"hits": {
					"total": {"value": 0, "relation": "eq"},
					"hits": []
				}
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResp:   true,
		},
		{
			name:     "Multiple matching flights",
			flightID: "flight123",
			response: `{
				"took": 1,
				"hits": {
					"total": {"value": 2, "relation": "eq"},
					"hits": [
						{
							"_id": "flight123",
							"_source": {
								"FlightNum": "XY123",
								"Origin": "LAX",
								"Destination": "JFK"
							}
						},
						{
							"_id": "flight123-2",
							"_source": {
								"FlightNum": "XY123",
								"Origin": "SFO",
								"Destination": "ORD"
							}
						}
					]
				}
			}`,
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResp:   true,
		},
		{
			name:           "Elasticsearch server error",
			flightID:       "flight123",
			response:       `{"error": {"type": "server_error", "reason": "Internal Server Error"}}`,
			statusCode:     http.StatusInternalServerError,
			expectedStatus: http.StatusOK,
			expectError:    false,
			validateResp:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock transport
			mockTransport := &MockTransport{
				Response:    tt.response,
				StatusCode:  tt.statusCode,
				ErrorString: tt.errorString,
			}

			// Setup controller and response recorder
			controller, w := setupTestController(t, mockTransport)

			// Set the flight ID in the context
			controller.Ctx.Input.SetParam(":id", tt.flightID)

			// Execute
			controller.GetFlightDetails()

			// Assert response status
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and validate response
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			if tt.expectError {
				assert.Contains(t, responseBody, "error")
			} else if tt.validateResp {
				// For successful responses, verify the structure
				if tt.statusCode == http.StatusOK {
					// Verify that we got a response that matches the structure we expect
					assert.NotNil(t, responseBody)
					
					if hits, ok := responseBody["hits"].(map[string]interface{}); ok {
						total, hasTotal := hits["total"].(map[string]interface{})
						assert.True(t, hasTotal, "Response should have a total field")
						
						value, hasValue := total["value"].(float64)
						assert.True(t, hasValue, "Total should have a value field")
						assert.True(t, value >= 0, "Total value should be non-negative")
						
						hitsArray, hasHits := hits["hits"].([]interface{})
						assert.True(t, hasHits, "Response should have a hits array")
						assert.Equal(t, int(value), len(hitsArray), "Hits array length should match total value")
					}
				}
			}
		})
	}
}