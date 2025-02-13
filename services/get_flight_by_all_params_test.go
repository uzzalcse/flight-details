package services

import (
	"testing"

	"flight-details/structs"

	"github.com/stretchr/testify/assert"
)

func TestSearchFlights(t *testing.T) {
	// Define the sample params for the flight search
	params := structs.FlightSearchParams{
		FlightNum:          "AA123",
		DestCountry:        "USA",
		OriginWeather:      "Clear",
		OriginCityName:     "New York",
		DestWeather:        "Rain",
		Dest:               "LAX",
		FlightDelayType:    "Weather",
		OriginCountry:      "USA",
		Carrier:            "American Airlines",
		Origin:             "JFK",
		DestRegion:         "California",
		OriginAirportID:    "JFK",
		OriginRegion:       "Northeast",
		DestCityName:       "Los Angeles",
		TravelTime:         "2025-02-03T10:33:28", // Example timestamp
		AvgTicketPrice:     200.50,
		DistanceMiles:      3000,
		DistanceKilometers: 4800,
		FlightTimeMin:      180,
		FlightTimeHour:     3,
		FlightDelay:        true,
		Cancelled:          false,
	}

	// Call the SearchFlights function
	query := SearchFlights(params)

	// Test that the query map has been correctly built
	assert.NotNil(t, query)
	assert.Contains(t, query, "query")
	assert.Contains(t, query["query"], "bool")
	assert.Contains(t, query["query"].(map[string]interface{})["bool"], "must")

	// Test for the timestamp range query
	timestampQuery := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})[0]["range"]
	rangeQuery := timestampQuery.(map[string]interface{})["timestamp"].(map[string]interface{})
	assert.Equal(t, rangeQuery["gte"], params.TravelTime)
	assert.Equal(t, rangeQuery["lte"], params.TravelTime)

	// Test for FlightNum term query
	flightNumQuery := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})[1]["term"]
	flightNumValue := flightNumQuery.(map[string]interface{})["FlightNum"]
	assert.Equal(t, flightNumValue, params.FlightNum)

	// Test for the other dynamic filters (example: DestCountry)
	destCountryQuery := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})[2]["term"]
	destCountryValue := destCountryQuery.(map[string]interface{})["DestCountry"]
	assert.Equal(t, destCountryValue, params.DestCountry)
}
