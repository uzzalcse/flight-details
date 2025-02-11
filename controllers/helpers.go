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
func ParseFlightSearchRequest(c *FlightController) (structs.FlightSearchParams, error) {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}
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
		DayOfWeek:         ParseInt(ctx.Input.Query("dayOfWeek")),
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
		FlightDelayMin:    ParseInt(ctx.Input.Query("FlightDelayMin")),
		Cancelled:         ParseBool(ctx.Input.Query("Cancelled")),
		FlightDelay:       ParseBool(ctx.Input.Query("FlightDelay")),
	}

	// Parse float values
	params.AvgTicketPrice = ParseFloat(ctx.Input.Query("AvgTicketPrice"))
	params.DistanceMiles = ParseFloat(ctx.Input.Query("DistanceMiles"))
	params.DistanceKilometers = ParseFloat(ctx.Input.Query("DistanceKilometers"))
	params.FlightTimeMin = ParseFloat(ctx.Input.Query("FlightTimeMin"))
	params.FlightTimeHour = ParseFloat(ctx.Input.Query("FlightTimeHour"))

	fmt.Println(params.TravelTime)

	// The "timestamp" (TravelTime) is required
	if params.TravelTime == "" {
		return params, errors.New("timestamp (TravelTime) is required")
	}

	return params, nil
}

// formatSuccessResponse builds a consistent JSON success structure
func FormatSuccessResponse(data string) map[string]interface{} {
	var parsedBody map[string]interface{}
	if err := json.Unmarshal([]byte(data), &parsedBody); err != nil {
		log.Println("Error formatting response:", err)
		return map[string]interface{}{
			"status":  "error",
			"message": "Response parsing failed",
		}
	}

	if _, exists := parsedBody["status"]; exists {
		return parsedBody
	}

	return map[string]interface{}{
		"status": "success",
		"data":   parsedBody,
	}
}

// Helper parse functions for query params
func ParseFloat(value string) float64 {
	if value == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(value, 64)
	return v
}

func ParseInt(value string) int {
	if value == "" {
		return 0
	}
	v, _ := strconv.Atoi(value)
	return v
}

func ParseBool(value string) bool {
	return value == "true" || value == "1"
}
