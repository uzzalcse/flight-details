package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/beego/beego/v2/server/web"
	"github.com/elastic/go-elasticsearch/v8"
)

// ESClient wraps the underlying elasticsearch.Client
type ESClient struct {
	Client *elasticsearch.Client
}

// GetElasticClient initializes and returns an ESClient
func GetElasticClient() *ESClient {
	// Retrieve configuration from beego config
	ES_API_KEY, err := web.AppConfig.String("ES_LOCAL_API_KEY")
	if err != nil {
		log.Fatalf("Error getting elasticsearch API key: %s", err)
	}
	ES_ADDRESS, err := web.AppConfig.String("ES_LOCAL_ADDRESS")
	if err != nil {
		log.Fatalf("Error getting elasticsearch address: %s", err)
	}

	// Build the elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{
			ES_ADDRESS,
		},
		APIKey: ES_API_KEY,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Return our wrapped ESClient
	return &ESClient{Client: es}
}

// ExecuteSearch performs a search using the ESClient
func (ec *ESClient) ExecuteSearch(query map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer

	// Encode the query into JSON
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %v", err)
	}

	// Perform the search
	res, err := ec.Client.Search(
		ec.Client.Search.WithContext(context.Background()),
		ec.Client.Search.WithIndex("kibana_sample_data_flights"),
		ec.Client.Search.WithBody(&buf),
		ec.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing search: %v", err)
	}
	defer res.Body.Close()

	// Decode the JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return result, nil
}
