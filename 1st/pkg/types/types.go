package types

type ExampleJson struct {
	Type            string       `json:"_type"`
	Id              int64        `json:"_id"`
	Key             *int64       `json:"key"`
	Name            string       `json:"name"`
	FullName        string       `json:"fullName"`
	LocationID      int64        `json:"location_id"`
	IataAirportCode *string      `json:"iata_airport_code"`
	Type_           string       `json:"type"`
	Country         string       `json:"country"`
	GeoPosition     *GeoPosition `json:"geoPosition"`
	InEurope        bool         `json:"inEurope"`
	CountryCode     string       `json:"countryCode"`
	CoreCountry     bool         `json:"coreCountry"`
	Distance        *float64     `json:"distance"`
}

type GeoPosition struct {
	Latitude  float64
	Longitude float64
}
