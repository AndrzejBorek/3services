package types

// In real case, probably this should be in separate repository, so developers working on 1st service would not interfere
// with structures that are being imported by 2nd, but since this is simple implementation, I left it how it is.

type ExampleJson struct {
	Type            string      `json:"_type"`
	Id              int64       `json:"_id"`
	Key             int64       `json:"key"`
	Name            string      `json:"name"`
	FullName        string      `json:"fullName"`
	LocationID      int64       `json:"location_id"`
	IataAirportCode string      `json:"iata_airport_code"`
	Type_           string      `json:"type"`
	Country         string      `json:"country"`
	GeoPosition     GeoPosition `json:"geoPosition"`
	InEurope        bool        `json:"inEurope"`
	CountryCode     string      `json:"countryCode"`
	CoreCountry     bool        `json:"coreCountry"`
	Distance        float64     `json:"distance"`
}

type GeoPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (json ExampleJson) ConvertToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["Type"] = json.Type
	m["Id"] = json.Id
	m["Key"] = json.Key
	m["Name"] = json.Name
	m["FullName"] = json.FullName
	m["LocationID"] = json.LocationID
	m["IataAirportCode"] = json.IataAirportCode
	m["Type_"] = json.Type_
	m["Country"] = json.Country
	m["GeoPosition"] = map[string]interface{}{
		"Latitude":  json.GeoPosition.Latitude,
		"Longitude": json.GeoPosition.Longitude,
	}
	m["InEurope"] = json.InEurope
	m["CountryCode"] = json.CountryCode
	m["CoreCountry"] = json.CoreCountry
	m["Distance"] = json.Distance
	return m
}
