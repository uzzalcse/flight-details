// controllers/flight_controller.go
package controllers

import (
	"flight-details/services"
	"net/http"
)

// @Summary Search Flights
// @Description Searches for flights based on provided filters
// @Tags Flights
// @Accept  json
// @Produce  json
// @Param timestamp query string true "Flight departure timestamp (ISO format)"
// @Param FlightNum query string false "Flight Number"
// @Param DestCountry query string false "Destination Country"
// @Param OriginWeather query string false "Origin Weather Condition"
// @Param OriginCityName query string false "Origin City Name"
// @Param AvgTicketPrice query number false "Average Ticket Price"
// @Param DistanceMiles query number false "Flight Distance in Miles"
// @Param FlightDelay query boolean false "Whether the flight was delayed"
// @Param DestWeather query string false "Destination Weather Condition"
// @Param Dest query string false "Destination Airport Name"
// @Param FlightDelayType query string false "Type of Flight Delay"
// @Param OriginCountry query string false "Origin Country"
// @Param dayOfWeek query integer false "Day of the week (0-6)"
// @Param DistanceKilometers query number false "Flight Distance in Kilometers"
// @Param DestLocationLat query string false "Destination Location Latitude"
// @Param DestLocationLon query string false "Destination Location Longitude"
// @Param DestAirportID query string false "Destination Airport ID"
// @Param Carrier query string false "Airline Carrier"
// @Param Cancelled query boolean false "Whether the flight was canceled"
// @Param FlightTimeMin query number false "Flight Duration in Minutes"
// @Param Origin query string false "Origin Airport Name"
// @Param OriginLocationLat query string false "Origin Location Latitude"
// @Param OriginLocationLon query string false "Origin Location Longitude"
// @Param DestRegion query string false "Destination Region"
// @Param OriginAirportID query string false "Origin Airport ID"
// @Param OriginRegion query string false "Origin Region"
// @Param DestCityName query string false "Destination City Name"
// @Param FlightTimeHour query number false "Flight Duration in Hours"
// @Param FlightDelayMin query integer false "Flight Delay in Minutes"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/api/flights/all_params/search [get]
func (c *FlightController) GetByAllParams() {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	// Parse incoming query parameters
	params, err := ParseFlightSearchRequest(c)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.Ctx.Output.SetStatus(http.StatusBadRequest) // ✅ Return `400 Bad Request`
		c.ServeJSON()
		return
	}

	// Call the service layer to get flight data
	result, err := services.SearchFlights(params)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": "Error fetching flight data: " + err.Error(),
		}
		c.Ctx.Output.SetStatus(http.StatusInternalServerError) // ✅ Return `500 Internal Server Error`
		c.ServeJSON()
		return
	}

	// Format and return success response
	c.Data["json"] = FormatSuccessResponse(result)
	c.ServeJSON()
}
