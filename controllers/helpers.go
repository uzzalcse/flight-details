package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"flight-details/structs"
)

// Constants for query parameters
const (
	FlightNumKey       = "FlightNum"
	CancelledKey       = "Cancelled"
	DayOfWeekKey       = "dayOfWeek"
	FlightDelayMinKey  = "FlightDelayMin"
	AvgTicketPriceKey  = "AvgTicketPrice"
	DistanceMilesKey   = "DistanceMiles"
	DistanceKMKey      = "DistanceKilometers"
	FlightTimeMinKey   = "FlightTimeMin"
	FlightTimeHourKey  = "FlightTimeHour"
	TimestampKey       = "timestamp"
)

// ParseFlightSearchRequest extracts parameters from the request context
func ParseFlightSearchRequest(c *FlightController) (structs.FlightSearchParams, error) {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	ctx := c.Ctx
	params := structs.FlightSearchParams{
		FlightNum:  ctx.Input.Query(FlightNumKey),
		TravelTime: ctx.Input.Query(TimestampKey), // Mandatory field
	}

	var validationErrors []string

	if params.TravelTime == "" {
		validationErrors = append(validationErrors, "timestamp (TravelTime) is required")
	}

	params.DayOfWeek, _ = ParseInt(ctx.Input.Query(DayOfWeekKey))
	if params.DayOfWeek < 0 || params.DayOfWeek > 6 {
		validationErrors = append(validationErrors, "invalid dayOfWeek: must be between 0 and 6")
	}

	params.FlightDelayMin, _ = ParseInt(ctx.Input.Query(FlightDelayMinKey))
	params.Cancelled, _ = ParseBool(ctx.Input.Query(CancelledKey))
	params.AvgTicketPrice, _ = ParseFloat(ctx.Input.Query(AvgTicketPriceKey))
	params.DistanceMiles, _ = ParseFloat(ctx.Input.Query(DistanceMilesKey))
	params.DistanceKilometers, _ = ParseFloat(ctx.Input.Query(DistanceKMKey))
	params.FlightTimeMin, _ = ParseFloat(ctx.Input.Query(FlightTimeMinKey))
	params.FlightTimeHour, _ = ParseFloat(ctx.Input.Query(FlightTimeHourKey))

	if len(validationErrors) > 0 {
		return params, errors.New("validation errors: " + fmt.Sprintf("%v", validationErrors))
	}

	return params, nil
}

// FormatSuccessResponse builds a consistent JSON success structure
func FormatSuccessResponse(data string) map[string]interface{} {
	var parsedBody map[string]interface{}

	if err := json.Unmarshal([]byte(data), &parsedBody); err != nil {
		log.Printf("Response parsing failed: %v (data: %s)", err, data)
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

// Generic function to parse string values into different types
func ParseValue[T any](value string, parseFunc func(string) (T, error)) (T, error) {
	var zero T
	if value == "" {
		return zero, nil
	}
	return parseFunc(value)
}

// Parsing functions using ParseValue
func ParseFloat(value string) (float64, error) {
    return ParseValue(value, func(s string) (float64, error) {
        return strconv.ParseFloat(s, 64) // Default to 64-bit precision
    })
}


func ParseInt(value string) (int, error) {
	return ParseValue(value, strconv.Atoi)
}

func ParseBool(value string) (bool, error) {
	switch value {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean: %s", value)
	}
}
