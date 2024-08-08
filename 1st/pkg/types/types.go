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
	m["type"] = json.Type
	m["id"] = json.Id
	m["key"] = json.Key
	m["name"] = json.Name
	m["fullname"] = json.FullName
	m["locationid"] = json.LocationID
	m["iataairportcode"] = json.IataAirportCode
	m["type_"] = json.Type_
	m["country"] = json.Country
	m["geoposition"] = map[string]interface{}{
		"latitude":  json.GeoPosition.Latitude,
		"longitude": json.GeoPosition.Longitude,
	}
	m["ineurope"] = json.InEurope
	m["countrycode"] = json.CountryCode
	m["corecountry"] = json.CoreCountry
	m["distance"] = json.Distance
	return m
}
