// services/flight_service.go
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"flight-details/structs"
	"flight-details/utils"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// SearchFlights queries Elasticsearch based on exact input filters
func SearchFlights(params structs.FlightSearchParams) (string, error) {
	// Retrieve your wrapped ESClient from utils
	es := utils.GetElasticClient()

	// Build the base query (must match "timestamp" exactly)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"term": map[string]interface{}{"timestamp": params.TravelTime}},
				},
			},
		},
	}

	// Dynamically add exact match filters
	addTermQuery(query, "FlightNum", params.FlightNum)
	addTermQuery(query, "DestCountry", params.DestCountry)
	addTermQuery(query, "OriginWeather", params.OriginWeather)
	addTermQuery(query, "OriginCityName", params.OriginCityName)
	addTermQuery(query, "DestWeather", params.DestWeather)
	addTermQuery(query, "Dest", params.Dest)
	addTermQuery(query, "FlightDelayType", params.FlightDelayType)
	addTermQuery(query, "OriginCountry", params.OriginCountry)
	addTermQuery(query, "DestAirportID", params.DestAirportID)
	addTermQuery(query, "Carrier", params.Carrier)
	addTermQuery(query, "Origin", params.Origin)
	addTermQuery(query, "DestRegion", params.DestRegion)
	addTermQuery(query, "OriginAirportID", params.OriginAirportID)
	addTermQuery(query, "OriginRegion", params.OriginRegion)
	addTermQuery(query, "DestCityName", params.DestCityName)

	// Ensure range queries check for exact matches
	addExactRangeQuery(query, "AvgTicketPrice", params.AvgTicketPrice)
	addExactRangeQuery(query, "DistanceMiles", params.DistanceMiles)
	addExactRangeQuery(query, "DistanceKilometers", params.DistanceKilometers)
	addExactRangeQuery(query, "FlightTimeMin", params.FlightTimeMin)
	addExactRangeQuery(query, "FlightTimeHour", params.FlightTimeHour)

	// Boolean filters
	addBoolQuery(query, "FlightDelay", params.FlightDelay)
	addBoolQuery(query, "Cancelled", params.Cancelled)

	// Exact geolocation filter
	addGeoLocQuery(query, params.OriginLocationLat, params.OriginLocationLon)

	// Convert query to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform search request
	req := esapi.SearchRequest{
		Index: []string{"kibana_sample_data_flights"},
		Body:  &buf,
	}

	// Use es.Client (the *elasticsearch.Client) to execute
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonResult), nil
}

// Below are the helper functions for building the query:

func addTermQuery(query map[string]interface{}, field string, value string) {
	if value != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] =
			append(query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
				map[string]interface{}{
					"term": map[string]interface{}{
						field: value,
					},
				},
			)
	}
}

func addExactRangeQuery(query map[string]interface{}, field string, value float64) {
	if value > 0 {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] =
			append(query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
				map[string]interface{}{
					"range": map[string]interface{}{
						field: map[string]interface{}{
							"gte": value,
							"lte": value, // Ensures exact match
						},
					},
				},
			)
	}
}

func addBoolQuery(query map[string]interface{}, field string, value bool) {
	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] =
		append(query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{
					field: value,
				},
			},
		)
}

func addGeoLocQuery(query map[string]interface{}, lat, lon string) {
	if lat != "" && lon != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] =
			append(query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
				map[string]interface{}{
					"geo_distance": map[string]interface{}{
						"distance": "1m", // Very small distance => near-exact geo match
						"OriginLocation": map[string]interface{}{
							"lat": lat,
							"lon": lon,
						},
					},
				},
			)
	}
}
