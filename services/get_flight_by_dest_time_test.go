package services_test

import (
	"errors"
	"flight-details/services"
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

func TestSearchFlights_Success(t *testing.T) {
    // Patch GetElasticClient to return a mock ESClient
    monkey.Patch(utils.GetElasticClient, func() *utils.ESClient {
        return &utils.ESClient{}
    })
    defer monkey.UnpatchAll()

    // Patch ExecuteSearch method
    monkey.PatchInstanceMethod(reflect.TypeOf((*utils.ESClient)(nil)), "ExecuteSearch", func(_ *utils.ESClient, _ map[string]interface{}) (map[string]interface{}, error) {
        return mockElasticResponse, nil
    })

    flights, err := services.SearchFlightDetails("Treviso", "2025-02-03T10:33:28")

    assert.Nil(t, err, "Expected no error, but got: %v", err)
    assert.NotNil(t, flights, "Expected flights to not be nil")
    assert.Len(t, flights, 1, "Expected 1 flight result but got %d", len(flights))
    assert.Equal(t, "EAYQW69", flights[0].FlightNum, "FlightNum does not match expected")
}

func TestSearchFlights_ExecutionError(t *testing.T) {
    // Patch GetElasticClient
    monkey.Patch(utils.GetElasticClient, func() *utils.ESClient {
        return &utils.ESClient{}
    })
    defer monkey.UnpatchAll()

    // Patch ExecuteSearch to return an error
    monkey.PatchInstanceMethod(reflect.TypeOf((*utils.ESClient)(nil)), "ExecuteSearch", func(_ *utils.ESClient, _ map[string]interface{}) (map[string]interface{}, error) {
        return nil, errors.New("Elasticsearch execution error")
    })

    flights, err := services.SearchFlightDetails("Treviso", "2025-02-03T10:33:28")

    assert.NotNil(t, err, "Expected an error but got nil")
    assert.Contains(t, err.Error(), "error executing search", "Error message mismatch")
    assert.Nil(t, flights, "Expected nil flights but got: %v", flights)
}

func TestSearchFlights_InvalidJSONResponse(t *testing.T) {
    // Patch GetElasticClient
    monkey.Patch(utils.GetElasticClient, func() *utils.ESClient {
        return &utils.ESClient{}
    })
    defer monkey.UnpatchAll()

    // Patch ExecuteSearch to return an error
    monkey.PatchInstanceMethod(reflect.TypeOf((*utils.ESClient)(nil)), "ExecuteSearch",
        func(_ *utils.ESClient, _ map[string]interface{}) (map[string]interface{}, error) {
            return nil, errors.New("invalid JSON response")
        })

    flights, err := services.SearchFlightDetails("Treviso", "2025-02-03T10:33:28")

    assert.NotNil(t, err, "Expected an error but got nil")
    assert.Contains(t, err.Error(), "invalid JSON response", "Error message mismatch")
    assert.Nil(t, flights, "Expected nil flights but got: %v", flights)
}


func TestAddMatchQuery(t *testing.T) {
    qb := services.NewQueryBuilder()
    qb.AddMatchQuery("test_field", "test_value")

    expectedQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "match": map[string]interface{}{
                            "test_field": "test_value",
                        },
                    },
                },
            },
        },
        "size": 100,
    }

    assert.Equal(t, expectedQuery, qb.Build(), "The query built does not match the expected output")
}


func TestSearchFlights_EmptyDestination(t *testing.T) {
    flights, err := services.SearchFlightDetails("", "2025-02-03T10:33:28")

    assert.NotNil(t, err, "Expected an error for empty destination")
    assert.Equal(t, "destination city name is required", err.Error(), "Error message mismatch")
    assert.Nil(t, flights, "Expected nil flights")
}

func TestSearchFlights_InvalidTimestampFormat(t *testing.T) {
    flights, err := services.SearchFlightDetails("Treviso", "invalid-date")

    assert.NotNil(t, err, "Expected an error for invalid timestamp format")
    assert.Equal(t, "invalid timestamp format. Expected format: YYYY-MM-DDTHH:MM:SS", err.Error(), "Error message mismatch")
    assert.Nil(t, flights, "Expected nil flights")
}

func TestSearchFlights_JSONMarshallingError(t *testing.T) {
    // Patch GetElasticClient
    monkey.Patch(utils.GetElasticClient, func() *utils.ESClient {
        return &utils.ESClient{}
    })
    defer monkey.UnpatchAll()

    // Patch ExecuteSearch to return an invalid type that causes JSON marshalling failure
    monkey.PatchInstanceMethod(reflect.TypeOf((*utils.ESClient)(nil)), "ExecuteSearch", func(_ *utils.ESClient, _ map[string]interface{}) (map[string]interface{}, error) {
        return map[string]interface{}{
            "invalid": make(chan int), // This will cause JSON marshalling to fail
        }, nil
    })

    flights, err := services.SearchFlightDetails("Treviso", "2025-02-03T10:33:28")

    assert.NotNil(t, err, "Expected an error for JSON marshalling failure")
    assert.Contains(t, err.Error(), "error marshalling response", "Error message mismatch")
    assert.Nil(t, flights, "Expected nil flights")
}

func TestSearchFlights_JSONUnmarshallingError(t *testing.T) {
    monkey.Patch(utils.GetElasticClient, func() *utils.ESClient {
        return &utils.ESClient{}
    })
    defer monkey.UnpatchAll()

    // Patch ExecuteSearch to return malformed JSON
    monkey.PatchInstanceMethod(reflect.TypeOf((*utils.ESClient)(nil)), "ExecuteSearch", func(_ *utils.ESClient, _ map[string]interface{}) (map[string]interface{}, error) {
        return map[string]interface{}{
            "hits": "invalid-json-format",
        }, nil
    })

    flights, err := services.SearchFlightDetails("Treviso", "2025-02-03T10:33:28")

    assert.NotNil(t, err, "Expected an error for JSON unmarshalling failure")
    assert.Contains(t, err.Error(), "error unmarshaling response", "Error message mismatch")
    assert.Nil(t, flights, "Expected nil flights")
}
