package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"flight-details/structs"
)

// parseFlightSearchRequest extracts parameters from the request context
func parseFlightSearchRequest(c *FlightController) (structs.FlightSearchParams, error) {
	ctx := c.Ctx
	params := structs.FlightSearchParams{
		FlightNum:         ctx.Input.Query("FlightNum"),
		DestCountry:       ctx.Input.Query("DestCountry"),
		OriginWeather:     ctx.Input.Query("OriginWeather"),
		OriginCityName:    ctx.Input.Query("OriginCityName"),
		DestWeather:       ctx.Input.Query("DestWeather"),
		Dest:              ctx.Input.Query("Dest"),
		FlightDelayType:   ctx.Input.Query("FlightDelayType"),
		OriginCountry:     ctx.Input.Query("OriginCountry"),
		DayOfWeek:         parseInt(ctx.Input.Query("dayOfWeek")),
		TravelTime:        ctx.Input.Query("timestamp"), // Mandatory field
		DestLocationLat:   ctx.Input.Query("DestLocationLat"),
		DestLocationLon:   ctx.Input.Query("DestLocationLon"),
		DestAirportID:     ctx.Input.Query("DestAirportID"),
		Carrier:           ctx.Input.Query("Carrier"),
		Origin:            ctx.Input.Query("Origin"),
		OriginLocationLat: ctx.Input.Query("OriginLocationLat"),
		OriginLocationLon: ctx.Input.Query("OriginLocationLon"),
		DestRegion:        ctx.Input.Query("DestRegion"),
		OriginAirportID:   ctx.Input.Query("OriginAirportID"),
		OriginRegion:      ctx.Input.Query("OriginRegion"),
		DestCityName:      ctx.Input.Query("DestCityName"),
		FlightDelayMin:    parseInt(ctx.Input.Query("FlightDelayMin")),
		Cancelled:         parseBool(ctx.Input.Query("Cancelled")),
		FlightDelay:       parseBool(ctx.Input.Query("FlightDelay")),
	}

	// Parse float values
	params.AvgTicketPrice = parseFloat(ctx.Input.Query("AvgTicketPrice"))
	params.DistanceMiles = parseFloat(ctx.Input.Query("DistanceMiles"))
	params.DistanceKilometers = parseFloat(ctx.Input.Query("DistanceKilometers"))
	params.FlightTimeMin = parseFloat(ctx.Input.Query("FlightTimeMin"))
	params.FlightTimeHour = parseFloat(ctx.Input.Query("FlightTimeHour"))

	fmt.Println(params.TravelTime)

	// The "timestamp" (TravelTime) is required
	if params.TravelTime == "" {
		return params, errors.New("timestamp (TravelTime) is required")
	}

	return params, nil
}

// formatSuccessResponse builds a consistent JSON success structure
func formatSuccessResponse(data string) map[string]interface{} {
	var parsedBody map[string]interface{}
	if err := json.Unmarshal([]byte(data), &parsedBody); err != nil {
		log.Println("Error formatting response:", err)
		return map[string]interface{}{
			"status":  "error",
			"message": "Response parsing failed",
		}
	}
	return map[string]interface{}{
		"status": "success",
		"data":   parsedBody,
	}
}

// Helper parse functions for query params
func parseFloat(value string) float64 {
	if value == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(value, 64)
	return v
}

func parseInt(value string) int {
	if value == "" {
		return 0
	}
	v, _ := strconv.Atoi(value)
	return v
}

func parseBool(value string) bool {
	return value == "true" || value == "1"
}
