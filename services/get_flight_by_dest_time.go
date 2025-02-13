package services

import (
	"errors"
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
			"size": 100, 
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

// Build returns the final query as a map
func (qb *QueryBuilder) Build() map[string]interface{} {
	qb.query["size"] = 100
	
	// Add source transformation to flatten nested location fields
	qb.query["_source"] = map[string]interface{}{
		"includes": []string{
			"*",
			"DestLocation.lat",
			"DestLocation.lon",
			"OriginLocation.lat",
			"OriginLocation.lon",
		},
	}
	
	// Add runtime fields to map nested locations to flat fields
	qb.query["runtime_mappings"] = map[string]interface{}{
		"DestLocationLat": map[string]interface{}{
			"type": "keyword",
			"script": map[string]interface{}{
				"source": "emit(doc['DestLocation.lat'].value)",
			},
		},
		"DestLocationLon": map[string]interface{}{
			"type": "keyword",
			"script": map[string]interface{}{
				"source": "emit(doc['DestLocation.lon'].value)",
			},
		},
		"OriginLocationLat": map[string]interface{}{
			"type": "keyword",
			"script": map[string]interface{}{
				"source": "emit(doc['OriginLocation.lat'].value)",
			},
		},
		"OriginLocationLon": map[string]interface{}{
			"type": "keyword",
			"script": map[string]interface{}{
				"source": "emit(doc['OriginLocation.lon'].value)",
			},
		},
	}
	
	return qb.query
}


// SearchFlightDetails searches for flights using the modular query builder
func SearchFlightDetails(destination, date string) (map[string]interface{}, error) {
	// Validate destination (should not be empty)
	if destination == "" {
		return nil, errors.New("destination city name is required")
	}

	// Validate timestamp format (ISO 8601 format: YYYY-MM-DDTHH:MM:SS)
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`, date)
	if !match {
		return nil, errors.New("invalid timestamp format. Expected format: YYYY-MM-DDTHH:MM:SS")
	}

	esClient := utils.Client

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

	return resp, nil
}