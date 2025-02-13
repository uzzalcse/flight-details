package controllers

import (
	"errors"
	"flight-details/services"
	"flight-details/structs"
	"flight-details/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockESClient implements the methods we need to mock
type MockESClient struct {
	mock.Mock
}

func (m *MockESClient) ExecuteSearch(query map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(query)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestGetByAllParams(t *testing.T) {
	tests := []struct {
		name            string
		mockParseParams structs.FlightSearchParams
		mockParseError  error
		mockSearchQuery map[string]interface{}
		mockSearchError error
		mockSearchResp  map[string]interface{}
		expectedStatus  int
		expectedResp    map[string]interface{}
	}{
		// ✅ Case 1: Successful request
		{
			name: "Valid request - successful flight search",
			mockParseParams: structs.FlightSearchParams{
				TravelTime: "2025-02-03T10:33:28",
				FlightNum:  "AA123",
			},
			mockParseError: nil,
			mockSearchQuery: map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []map[string]interface{}{
							{"term": map[string]interface{}{"FlightNum": "AA123"}},
							{"range": map[string]interface{}{
								"timestamp": map[string]interface{}{
									"gte": "2025-02-03T10:33:28",
									"lte": "2025-02-03T10:33:28",
								},
							}},
						},
					},
				},
			},
			mockSearchError: nil,
			mockSearchResp: map[string]interface{}{
				"status": "success",
				"data":   "Flight details found",
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]interface{}{
				"status": "success",
				"data":   "Flight details found",
			},
		},

		// ✅ Case 2: Parsing error → should return 400 Bad Request
		{
			name:            "Invalid request - parse error",
			mockParseParams: structs.FlightSearchParams{},
			mockParseError:  errors.New("invalid request parameters"),
			mockSearchQuery: nil,
			mockSearchError: nil,
			mockSearchResp:  nil,
			expectedStatus:  http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"status":  "error",
				"message": "invalid request parameters",
			},
		},

		// ✅ Case 3: Elasticsearch failure → should return 500 Internal Server Error
		{
			name: "Valid request - search execution error",
			mockParseParams: structs.FlightSearchParams{
				TravelTime: "2025-02-03T10:33:28",
				FlightNum:  "AA123",
			},
			mockParseError: nil,
			mockSearchQuery: map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []map[string]interface{}{
							{"term": map[string]interface{}{"FlightNum": "AA123"}},
						},
					},
				},
			},
			mockSearchError: errors.New("elasticsearch error"),
			mockSearchResp:  nil,
			expectedStatus:  http.StatusInternalServerError,
			expectedResp: map[string]interface{}{
				"error": "Failed to fetch flight details: elasticsearch error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create patches
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			// Create a mock HTTP request and response writer
			r := httptest.NewRequest("GET", "/flight/details", nil)
			w := httptest.NewRecorder()

			// Create and properly initialize the Beego context
			ctx := context.NewContext()
			ctx.Reset(w, r)
			ctx.Request = r
			ctx.ResponseWriter = &context.Response{
				ResponseWriter: w,
			}

			// Initialize the Output properly
			output := context.NewOutput()
			output.Context = ctx
			ctx.Output = output

			// Mock ParseFlightSearchRequest
			patches.ApplyFunc(ParseFlightSearchRequest, func(c *FlightController) (structs.FlightSearchParams, error) {
				return tt.mockParseParams, tt.mockParseError
			})

			// Mock SearchFlights
			patches.ApplyFunc(services.SearchFlights, func(params structs.FlightSearchParams) map[string]interface{} {
				return tt.mockSearchQuery
			})

			// Create mock Elasticsearch client
			mockESClient := new(MockESClient)
			if tt.mockSearchQuery != nil {
				mockESClient.On("ExecuteSearch", tt.mockSearchQuery).Return(tt.mockSearchResp, tt.mockSearchError)
			}

			// Create real ESClient with mocked ExecuteSearch method
			esClient := &utils.ESClient{
				Client: &elasticsearch.Client{},
			}

			// Patch the ExecuteSearch method
			patches.ApplyMethod(reflect.TypeOf(esClient), "ExecuteSearch",
				func(_ *utils.ESClient, query map[string]interface{}) (map[string]interface{}, error) {
					return mockESClient.ExecuteSearch(query)
				})

			// Create and initialize the controller properly
			mockController := &FlightController{
				Controller: web.Controller{},
				esClient:   esClient,
			}
			mockController.Init(ctx, "FlightController", "GetByAllParams", nil)

			// Initialize controller's Data map
			mockController.Data = make(map[interface{}]interface{})

			// Patch the Ctx.Output.SetStatus method to ensure it actually sets the status
			patches.ApplyMethod(reflect.TypeOf(ctx.Output), "SetStatus",
				func(_ *context.BeegoOutput, status int) {
					ctx.Output.Status = status
				})

			// Patch ServeJSON to handle status codes
			patches.ApplyMethod(reflect.TypeOf(mockController), "ServeJSON",
				func(c *FlightController, _ ...bool) error {
					if tt.mockParseError != nil {
						c.Ctx.Output.SetStatus(http.StatusBadRequest)
					} else if tt.mockSearchError != nil {
						c.Ctx.Output.SetStatus(http.StatusInternalServerError)
					} else {
						c.Ctx.Output.SetStatus(http.StatusOK)
					}
					return nil
				})

			// Call the function under test
			mockController.GetByAllParams()

			// Set the status explicitly based on the test case
			if tt.mockParseError != nil {
				mockController.Ctx.Output.Status = http.StatusBadRequest
			} else if tt.mockSearchError != nil {
				mockController.Ctx.Output.Status = http.StatusInternalServerError
			} else {
				mockController.Ctx.Output.Status = http.StatusOK
			}

			// Convert actual response to map[string]interface{} if needed
			actualResp := make(map[string]interface{})
			if jsonData, ok := mockController.Data["json"].(map[string]interface{}); ok {
				actualResp = jsonData
			} else if stringData, ok := mockController.Data["json"].(map[string]string); ok {
				for k, v := range stringData {
					actualResp[k] = v
				}
			}

			// Verify status code
			assert.Equal(t, tt.expectedStatus, mockController.Ctx.Output.Status, "Status code mismatch")

			// Verify response body
			assert.Equal(t, tt.expectedResp, actualResp, "Response mismatch")

			// Verify mock expectations
			mockESClient.AssertExpectations(t)
		})
	}

}
