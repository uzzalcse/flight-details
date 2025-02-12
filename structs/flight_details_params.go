package structs

type FlightSearchParams struct {
	FlightNum          string
	DestCountry        string
	OriginWeather      string
	OriginCityName     string
	DestWeather        string
	Dest               string
	FlightDelayType    string
	OriginCountry      string
	DayOfWeek          int
	TravelTime         string
	DestLocationLat    string
	DestLocationLon    string
	DestAirportID      string
	Carrier            string
	Origin             string
	OriginLocationLat  string
	OriginLocationLon  string
	DestRegion         string
	OriginAirportID    string
	OriginRegion       string
	DestCityName       string
	FlightDelayMin     int
	Cancelled          bool
	FlightDelay        bool
	AvgTicketPrice     float64
	DistanceMiles      float64
	DistanceKilometers float64
	FlightTimeMin      float64
	FlightTimeHour     float64
}
