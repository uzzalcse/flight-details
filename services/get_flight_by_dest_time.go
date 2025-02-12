package services

import (
	"encoding/json"
	"errors"
	"flight-details/structs"
	"flight-details/utils"
	"fmt"
	"regexp"
)

// QueryBuilder handles the construction of Elasticsearch queries
type QueryBuilder struct {
	query map[string]interface{}
}

// NewQueryBuilder creates a new instance of QueryBuilder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		query: map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []interface{}{},
				},
			},
		},
	}
}

// AddTermQuery adds a term query for exact matching
func (qb *QueryBuilder) AddTermQuery(field, value string) *QueryBuilder {
	mustClauses := qb.query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})
	termQuery := map[string]interface{}{
		"term": map[string]interface{}{
			field: value,
		},
	}
	qb.query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(mustClauses, termQuery)
	return qb
}

// AddMatchQuery adds a match query for text fields
func (qb *QueryBuilder) AddMatchQuery(field, value string) *QueryBuilder {
	mustClauses := qb.query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})
	matchQuery := map[string]interface{}{
		"match": map[string]interface{}{
			field: value,
		},
	}
	qb.query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(mustClauses, matchQuery)
	return qb
}

// Build returns the final query as a map
func (qb *QueryBuilder) Build() map[string]interface{} {
	qb.query["size"] = 100
	return qb.query
}

// SearchResult represents the structure of Elasticsearch response
type SearchResult struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source structs.FlightSearchParams `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// SearchFlights searches for flights using the modular query builder
func SearchFlightDetails(destination, date string) ([]structs.FlightSearchParams, error) {
	// Validate destination (should not be empty)
	if destination == "" {
		return nil, errors.New("destination city name is required")
	}

	// Validate timestamp format (ISO 8601 format: YYYY-MM-DDTHH:MM:SS)
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`, date)
	if !match {
		return nil, errors.New("invalid timestamp format. Expected format: YYYY-MM-DDTHH:MM:SS")
	}

	esClient := utils.GetElasticClient()
	queryBuilder := NewQueryBuilder()

	// Build the query using method chaining
	query := queryBuilder.
		AddTermQuery("DestCityName", destination).
		AddTermQuery("timestamp", date).
		Build()

	// Execute the search
	resp, err := esClient.ExecuteSearch(query)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %v", err)
	}

	// Parse the response
	var results SearchResult
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshalling response: %v", err)
	}

	err = json.Unmarshal(jsonData, &results)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Extract flights from results
	flights := make([]structs.FlightSearchParams, len(results.Hits.Hits))
	for i, hit := range results.Hits.Hits {
		flights[i] = hit.Source
	}

	return flights, nil
}
