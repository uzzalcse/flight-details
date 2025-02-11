package controllers

import (
	"flight-details/utils"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	_ "github.com/elastic/go-elasticsearch/v8"
)

type FlightController struct {
	beego.Controller
	esClient *utils.ESClient
}

// Init initializes the Elasticsearch client
func (f *FlightController) Init(ctx *beecontext.Context, controllerName, actionName string, app interface{}) {
	f.Controller.Init(ctx, controllerName, actionName, app)
	f.esClient = Client
}

// GetFlightDetails handles GET /v1/api/:id
// GetFlightDetails handles retrieving flight details from Elasticsearch
// @Title Get Flight Details
// @Description Fetch flight details by flight ID
// @Param   id      path    string  true  "Flight ID"
// @Success 200 {object} map[string]interface{} "Successful response"
// @Success 200 {object} map[string]interface{} "Example success response"
// @Success 200 {object} map[string]interface{} "Example failure response"
// @Failure 400 {object} map[string]string "Bad Request - Flight ID is required"
// @Failure 500 {object} map[string]string "Internal Server Error - Failed to fetch flight details"
// @Router /v1/api/flights/{id} [get]
func (f *FlightController) GetFlightDetails() {
	flightID := f.Ctx.Input.Param(":id")
	if flightID == "" {
		f.Data["json"] = map[string]string{"error": "Flight ID is required"}
		f.ServeJSON()
		return
	}

	// Create the search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"_id": []string{flightID},
			},
		},
	}

	// Execute search using ESClient
	res, err := f.esClient.ExecuteSearch(query)
	if err != nil {
		f.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to fetch flight details: %v", err)}
		f.ServeJSON()
		return
	}

	// Send response
	f.Data["json"] = res
	f.ServeJSON()
}
