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
	var err error
	params := structs.FlightSearchParams{
		FlightNum:       ctx.Input.Query("FlightNum"),
		DestCountry:     ctx.Input.Query("DestCountry"),
		OriginWeather:   ctx.Input.Query("OriginWeather"),
		OriginCityName:  ctx.Input.Query("OriginCityName"),
		DestWeather:     ctx.Input.Query("DestWeather"),
		Dest:            ctx.Input.Query("Dest"),
		FlightDelayType: ctx.Input.Query("FlightDelayType"),
		OriginCountry:   ctx.Input.Query("OriginCountry"),
		TravelTime:      ctx.Input.Query("timestamp"), // Mandatory field
		DestLocationLat: ctx.Input.Query("DestLocationLat"),
		DestLocationLon: ctx.Input.Query("DestLocationLon"),
		DestAirportID:   ctx.Input.Query("DestAirportID"),
		Carrier:         ctx.Input.Query("Carrier"),
		Origin:          ctx.Input.Query("Origin"),
		DestRegion:      ctx.Input.Query("DestRegion"),
		OriginAirportID: ctx.Input.Query("OriginAirportID"),
		OriginRegion:    ctx.Input.Query("OriginRegion"),
		DestCityName:    ctx.Input.Query("DestCityName"),
	}

	// Mandatory timestamp check
	if params.TravelTime == "" {
		return params, errors.New("timestamp (TravelTime) is required")
	}

	// Parse integer values with validation
	params.DayOfWeek, err = ParseInt(ctx.Input.Query("dayOfWeek"))
	if err != nil || params.DayOfWeek < 0 || params.DayOfWeek > 6 {
		return params, fmt.Errorf("invalid dayOfWeek: must be between 0 and 6")
	}

	params.FlightDelayMin, err = ParseInt(ctx.Input.Query("FlightDelayMin"))
	if err != nil {
		return params, fmt.Errorf("invalid FlightDelayMin value: %v", err)
	}

	// Parse boolean values with error handling
	params.Cancelled, err = ParseBool(ctx.Input.Query("Cancelled"))
	if err != nil {
		return params, fmt.Errorf("invalid Cancelled value: %v", err)
	}

	params.FlightDelay, err = ParseBool(ctx.Input.Query("FlightDelay"))
	if err != nil {
		return params, fmt.Errorf("invalid FlightDelay value: %v", err)
	}

	// Parse float values with error handling
	params.AvgTicketPrice, err = ParseFloat(ctx.Input.Query("AvgTicketPrice"))
	if err != nil {
		return params, fmt.Errorf("invalid AvgTicketPrice value: %v", err)
	}

	params.DistanceMiles, err = ParseFloat(ctx.Input.Query("DistanceMiles"))
	if err != nil {
		return params, fmt.Errorf("invalid DistanceMiles value: %v", err)
	}

	params.DistanceKilometers, err = ParseFloat(ctx.Input.Query("DistanceKilometers"))
	if err != nil {
		return params, fmt.Errorf("invalid DistanceKilometers value: %v", err)
	}

	params.FlightTimeMin, err = ParseFloat(ctx.Input.Query("FlightTimeMin"))
	if err != nil {
		return params, fmt.Errorf("invalid FlightTimeMin value: %v", err)
	}

	params.FlightTimeHour, err = ParseFloat(ctx.Input.Query("FlightTimeHour"))
	if err != nil {
		return params, fmt.Errorf("invalid FlightTimeHour value: %v", err)
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
func ParseFloat(value string) (float64, error) {
	if value == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float: %s", value)
	}
	return v, nil
}

func ParseInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %s", value)
	}
	return v, nil
}

func ParseBool(value string) (bool, error) {
	if value == "" { // Allow missing values to be treated as false
		return false, nil
	}
	if value == "true" || value == "1" {
		return true, nil
	} else if value == "false" || value == "0" {
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean: %s", value)
}
