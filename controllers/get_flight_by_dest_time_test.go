package controllers

import (
	"encoding/json"
	"errors"
	"flight-details/services"
	"flight-details/structs"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
)

func TestFlightController_Get(t *testing.T) {
	// Mock Response Data that matches actual API response structure
	mockFlights := []structs.FlightSearchParams{
		{
			FlightNum:          "EAYQW69",
			DestCountry:        "IT",
			OriginWeather:      "Thunder & Lightning",
			OriginCityName:     "Naples",
			DestWeather:        "Clear",
			Dest:               "Treviso-Sant'Angelo Airport",
			FlightDelayType:    "Weather Delay",
			OriginCountry:      "IT",
			DayOfWeek:          0,
			TravelTime:         "",
			DestLocationLat:    "",
			DestLocationLon:    "",
			DestAirportID:      "TV01",
			Carrier:            "Kibana Airlines",
			Origin:             "Naples International Airport",
			OriginLocationLat:  "",
			OriginLocationLon:  "",
			DestRegion:         "IT-34",
			OriginAirportID:    "NA01",
			OriginRegion:       "IT-72",
			DestCityName:       "Treviso",
			FlightDelayMin:     180,
			Cancelled:          true,
			FlightDelay:        true,
			AvgTicketPrice:     181.69421554118,
			DistanceMiles:      345.31943877289535,
			DistanceKilometers: 555.7377668725265,
			FlightTimeMin:      222.74905899019436,
			FlightTimeHour:     3.712484316503239,
		},
	}

	// Monkey Patch `services.SearchFlights` to return mock data
	monkey.Patch(services.SearchFlightDetails, func(destination, date string) ([]structs.FlightSearchParams, error) {
		return mockFlights, nil
	})
	defer monkey.UnpatchAll()

	req, _ := http.NewRequest("GET", "/v1/api/flights/dest_time/search?DestCityName=Treviso&timestamp=2025-02-03T10:33:28", nil)
	w := httptest.NewRecorder()

	// Create controller and properly initialize it
	flightController := &FlightController{}
	ctx := context.NewContext()
	ctx.Reset(w, req)

	ctx.Input = &context.BeegoInput{
		Context:     ctx,
		RequestBody: []byte{},
	}

	ctx.Output = &context.BeegoOutput{
		Context: ctx,
	}

	flightController.Init(ctx, "FlightController", "Get", nil)
	flightController.Get()

	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 but got %d", w.Code)

	// Parse JSON response
	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Expected valid JSON response but got error %v", err)

	// Create expected response based on mock data
	expectedResponse := []map[string]interface{}{
		{
			"FlightNum":          "EAYQW69",
			"DestCountry":        "IT",
			"OriginWeather":      "Thunder & Lightning",
			"OriginCityName":     "Naples",
			"DestWeather":        "Clear",
			"Dest":               "Treviso-Sant'Angelo Airport",
			"FlightDelayType":    "Weather Delay",
			"OriginCountry":      "IT",
			"DayOfWeek":          float64(0),
			"TravelTime":         "",
			"DestLocationLat":    "",
			"DestLocationLon":    "",
			"DestAirportID":      "TV01",
			"Carrier":            "Kibana Airlines",
			"Origin":             "Naples International Airport",
			"OriginLocationLat":  "",
			"OriginLocationLon":  "",
			"DestRegion":         "IT-34",
			"OriginAirportID":    "NA01",
			"OriginRegion":       "IT-72",
			"DestCityName":       "Treviso",
			"FlightDelayMin":     float64(180),
			"Cancelled":          true,
			"FlightDelay":        true,
			"AvgTicketPrice":     181.69421554118,
			"DistanceMiles":      345.31943877289535,
			"DistanceKilometers": 555.7377668725265,
			"FlightTimeMin":      222.74905899019436,
			"FlightTimeHour":     3.712484316503239,
		},
	}

	assert.Equal(t, expectedResponse, response, "Expected response does not match actual response")
}

func TestFlightController_Get_SearchError(t *testing.T) {
	monkey.Patch(services.SearchFlightDetails, func(destination, date string) ([]structs.FlightSearchParams, error) {
		return nil, errors.New("search service error")
	})
	defer monkey.UnpatchAll()

	req, _ := http.NewRequest("GET", "/v1/api/flights/dest_time/search?DestCityName=Treviso&timestamp=2025-02-03T10:33:28", nil)
	w := httptest.NewRecorder()

	flightController := &FlightController{}
	ctx := context.NewContext()
	ctx.Reset(w, req)

	ctx.Input = &context.BeegoInput{
		Context:     ctx,
		RequestBody: []byte{},
	}

	ctx.Output = &context.BeegoOutput{
		Context: ctx,
	}

	flightController.Init(ctx, "FlightController", "Get", nil)
	flightController.Get()

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status 500 but got %d", w.Code)

	// Parse JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err, "Expected valid JSON response but got error %v", err)

	assert.Equal(t, "search service error", response["error"], "Expected error message does not match")
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

			ctx.Input = &context.BeegoInput{
				Context:     ctx,
				RequestBody: []byte{},
			}

			ctx.Output = &context.BeegoOutput{
				Context: ctx,
			}

			flightController.Init(ctx, "FlightController", "Get", nil)
			flightController.Get()

			assert.Equal(t, tc.expectedStatus, w.Code, "Expected status %d but got %d", tc.expectedStatus, w.Code)

			// Parse JSON response
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.Nil(t, err, "Expected valid JSON response but got error %v", err)

			assert.Equal(t, tc.expectedError, response["error"], "Expected error message does not match")
		})
	}
}
