package controllers

import (
	"flight-details/services"
	"net/http"
	"regexp"
)

// @Summary Search for flights
// @Description Search for available flights based on destination and date
// @Tags Flights
// @Accept json
// @Produce json
// @Param DestCityName query string true "Destination city name" example:"London"
// @Param timestamp query string true "Flight date" example:"2024-02-10T10:33:28"
// @Success 200 {array} map[string]interface{} "List of flights"
// @Failure 400 {object} map[string]string "Bad request - Invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/api/flights/dest_time/search [get]
func (c *FlightController) Get() {
	destination := c.GetString("DestCityName")
	date := c.GetString("timestamp")

	// Validate destination (should not be empty)
	if destination == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Destination city name is required"}
		c.ServeJSON()
		return
	}

	// Validate timestamp format (ISO 8601 format: YYYY-MM-DDTHH:MM:SS)
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`, date)
	if !match {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid timestamp format. Expected format: YYYY-MM-DDTHH:MM:SS"}
		c.ServeJSON()
		return
	}

	// Call service layer
	flights, err := services.SearchFlightDetails(destination, date)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = flights
	}

	c.ServeJSON()
}
