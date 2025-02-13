package controllers

import (
	"encoding/json"
	"errors"
	"flight-details/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
)

func TestFlightController_Get_Success(t *testing.T) {
	mockFlights := map[string]interface{}{
		"hits": map[string]interface{}{
			"hits": []map[string]interface{}{
				{
					"_source": map[string]interface{}{
						"FlightNum":     "EAYQW69",
						"DestCityName":  "Treviso",
						"timestamp":     "2025-02-03T10:33:28",
						"Carrier":       "Kibana Airlines",
					},
				},
			},
		},
	}

	monkey.Patch(services.SearchFlightDetails, func(destination, date string) (map[string]interface{}, error) {
		return mockFlights, nil
	})
	defer monkey.UnpatchAll()

	req, _ := http.NewRequest("GET", "/v1/api/flights/dest_time/search?DestCityName=Treviso&timestamp=2025-02-03T10:33:28", nil)
	w := httptest.NewRecorder()

	flightController := &FlightController{}
	ctx := context.NewContext()
	ctx.Reset(w, req)
	ctx.Input = &context.BeegoInput{Context: ctx}
	ctx.Output = &context.BeegoOutput{Context: ctx}

	flightController.Init(ctx, "FlightController", "Get", nil)
	flightController.Get()

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 but got %d", w.Code)

	// Parse JSON response correctly as a map
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Expected valid JSON response but got error %v", err)

	// Extract the "hits" field properly
	hits, ok := response["hits"].(map[string]interface{})
	assert.True(t, ok, "Expected 'hits' field in response")

	hitsArray, ok := hits["hits"].([]interface{})
	assert.True(t, ok, "Expected 'hits' field to be an array")

	assert.Len(t, hitsArray, 1, "Expected 1 flight entry")

	flightData, ok := hitsArray[0].(map[string]interface{})["_source"].(map[string]interface{})
	assert.True(t, ok, "Expected _source field in flight entry")

	// Ensure expected values are returned
	assert.Equal(t, mockFlights["hits"].(map[string]interface{})["hits"].([]map[string]interface{})[0]["_source"], flightData, "Flight details do not match expected")
}


func TestFlightController_Get_SearchError(t *testing.T) {
	monkey.Patch(services.SearchFlightDetails, func(destination, date string) (map[string]interface{}, error) {
		return nil, errors.New("search service error")
	})
	defer monkey.UnpatchAll()

	req, _ := http.NewRequest("GET", "/v1/api/flights/dest_time/search?DestCityName=Treviso&timestamp=2025-02-03T10:33:28", nil)
	w := httptest.NewRecorder()

	flightController := &FlightController{}
	ctx := context.NewContext()
	ctx.Reset(w, req)
	ctx.Input = &context.BeegoInput{Context: ctx}
	ctx.Output = &context.BeegoOutput{Context: ctx}

	flightController.Init(ctx, "FlightController", "Get", nil)
	flightController.Get()

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "search service error", response["error"])
}

func TestFlightController_Get_MissingParameters(t *testing.T) {
	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Missing DestCityName",
			url:            "/v1/api/flights/dest_time/search?timestamp=2025-02-03T10:33:28",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Destination city name is required",
		},
		{
			name:           "Missing timestamp",
			url:            "/v1/api/flights/dest_time/search?DestCityName=Treviso",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid timestamp format. Expected format: YYYY-MM-DDTHH:MM:SS",
		},
		{
			name:           "Missing both parameters",
			url:            "/v1/api/flights/dest_time/search",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Destination city name is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()

			flightController := &FlightController{}
			ctx := context.NewContext()
			ctx.Reset(w, req)
			ctx.Input = &context.BeegoInput{Context: ctx}
			ctx.Output = &context.BeegoOutput{Context: ctx}

			flightController.Init(ctx, "FlightController", "Get", nil)
			flightController.Get()

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedError, response["error"])
		})
	}
}
