package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	//beego "github.com/beego/beego/v2/server/web"
	"flight-details/utils"
	"net/http/httptest"

	//beego "github.com/beego/beego/v2/server/web"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock transport
			mockTransport := &MockTransport{
				Response:    tt.response,
				StatusCode:  tt.statusCode,
				ErrorString: tt.errorString,
			}

			// Create Elasticsearch client with mock transport
			esClient, _ := elasticsearch.NewClient(elasticsearch.Config{
				Transport: mockTransport,
			})

			// Create controller with mocked ES client
			controller := &FlightController{
				esClient: &utils.ESClient{Client: esClient},
			}

			// Create test context
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/v1/api/"+tt.flightID, nil)
			context := beecontext.NewContext()
			context.Reset(w, r)

			// Initialize controller
			controller.Init(context, "FlightController", "GetFlightDetails", nil)

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
				hits, ok := responseBody["hits"].(map[string]interface{})
				assert.True(t, ok, "expected hits in response")

				total, ok := hits["total"].(map[string]interface{})
				assert.True(t, ok, "expected total in hits")

				value, ok := total["value"].(float64)
				assert.True(t, ok, "expected value to be a number")
				assert.True(t, value >= 0, "expected non-negative hit count")
			}
		})
	}
}

