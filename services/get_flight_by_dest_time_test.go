package services

import (
	"errors"
	"flight-details/utils"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

// Mocked Elasticsearch Response
var mockElasticResponse = map[string]interface{}{
	"hits": map[string]interface{}{
		"total": map[string]interface{}{
			"value": 1,
		},
		"hits": []map[string]interface{}{
			{
				"_source": map[string]interface{}{
					"FlightNum":          "EAYQW69",
					"DestCityName":       "Treviso",
					"timestamp":          "2025-02-03T10:33:28",
					"DestCountry":        "IT",
					"OriginWeather":      "Thunder & Lightning",
					"OriginCityName":     "Naples",
					"DestWeather":        "Clear",
					"Dest":               "Treviso-Sant'Angelo Airport",
					"FlightDelayType":    "Weather Delay",
					"OriginCountry":      "IT",
					"DayOfWeek":          0,
					"DestAirportID":      "TV01",
					"Carrier":            "Kibana Airlines",
					"Origin":             "Naples International Airport",
					"DestRegion":         "IT-34",
					"OriginAirportID":    "NA01",
					"OriginRegion":       "IT-72",
					"FlightDelayMin":     180,
					"Cancelled":          true,
					"FlightDelay":        true,
					"AvgTicketPrice":     181.69421554118,
					"DistanceMiles":      345.31943877289535,
					"DistanceKilometers": 555.7377668725265,
					"FlightTimeMin":      222.74905899019436,
					"FlightTimeHour":     3.712484316503239,
				},
			},
		},
	},
}

func setupTestEnvironment() func() {
	// Create a mock ESClient
	mockClient := &utils.ESClient{}
	
	// Store the original Client
	originalClient := utils.Client
	
	// Set our mock client
	utils.Client = mockClient
	
	// Return cleanup function
	return func() {
		utils.Client = originalClient
	}
}

func TestSearchFlights_Success(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment()
	defer cleanup()

	// Patch ExecuteSearch method
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(utils.Client), "ExecuteSearch", 
		func(_ *utils.ESClient, _ map[string]interface{}) (map[string]interface{}, error) {
			return mockElasticResponse, nil
		})
	defer patch.Unpatch()

	flights, err := SearchFlightDetails("Treviso", "2025-02-03T10:33:28")

	assert.Nil(t, err, "Expected no error, but got: %v", err)
	assert.NotNil(t, flights, "Expected flights to not be nil")
	assert.Len(t, flights["hits"].(map[string]interface{})["hits"].([]map[string]interface{}), 1, 
		"Expected 1 flight result but got a different number")
	assert.Equal(t, "EAYQW69", 
		flights["hits"].(map[string]interface{})["hits"].([]map[string]interface{})[0]["_source"].(map[string]interface{})["FlightNum"], 
		"FlightNum does not match expected")
}
func TestSearchFlights_ExecutionError(t *testing.T) {
	// Setup test environment
	cleanup := setupTestEnvironment()
	defer cleanup()

	// Patch ExecuteSearch to return an error
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(utils.Client), "ExecuteSearch",
		func(_ *utils.ESClient, _ map[string]interface{}) (map[string]interface{}, error) {
			return nil, errors.New("Elasticsearch execution error")
		})
	defer patch.Unpatch()

	flights, err := SearchFlightDetails("Treviso", "2025-02-03T10:33:28")

	assert.NotNil(t, err, "Expected an error but got nil")
	assert.Contains(t, err.Error(), "error executing search", "Error message mismatch")
	assert.Nil(t, flights, "Expected nil flights but got: %v", flights)
}

func TestSearchFlights_EmptyDestination(t *testing.T) {
	flights, err := SearchFlightDetails("", "2025-02-03T10:33:28")

	assert.NotNil(t, err, "Expected an error for empty destination")
	assert.Equal(t, "destination city name is required", err.Error(), "Error message mismatch")
	assert.Nil(t, flights, "Expected nil flights")
}

func TestSearchFlights_InvalidTimestampFormat(t *testing.T) {
	flights, err := SearchFlightDetails("Treviso", "invalid-date")

	assert.NotNil(t, err, "Expected an error for invalid timestamp format")
	assert.Equal(t, "invalid timestamp format. Expected format: YYYY-MM-DDTHH:MM:SS", err.Error(), "Error message mismatch")
	assert.Nil(t, flights, "Expected nil flights")
}

func TestAddTermQuery(t *testing.T) {
	// Initialize query builder
	qb := NewQueryBuilder()
	
	// Add term query
	qb.AddTermQuery("test_field", "test_value")

	// Get built query
	builtQuery := qb.Build()

	// Remove non-relevant fields for assertion
	delete(builtQuery, "_source")
	delete(builtQuery, "runtime_mappings")

	// Define expected query structure
	expectedQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"test_field": "test_value",
						},
					},
				},
			},
		},
		"size": 100,
	}

	// Assert equality
	assert.Equal(t, expectedQuery, builtQuery, "The query built does not match the expected output")
}


