package utils

import (
	"unsafe"
    "math/rand"
    "time"
    "sync"
)

const (
    letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits  // # of letter indices fitting in 63 bits
)

type customRandomGenerator struct{
    src rand.Source
    lock sync.Mutex //Since rand.Int63 is not safe for concurrent use, I had to create this struct
}

func(gen *customRandomGenerator) generateRandomInt63() int64{
    gen.lock.Lock()
    defer gen.lock.Unlock()
    return gen.src.Int63()
}

func createNewCustomRandomGenerator(seed int64) *customRandomGenerator{
    return &customRandomGenerator{
        src: rand.NewSource(seed),
    }
}


var gen = createNewCustomRandomGenerator(time.Now().UnixNano())


type ExampleJson struct {
    Type           string  `json:"_type"`
    Id             int64   `json:"_id"`
    Key            *int64  `json:"key"` 
    Name           string  `json:"name"`
    FullName       string  `json:"fullName"`
    LocationID     int64   `json:"location_id"`
    IataAirportCode *string `json:"iata_airport_code"`
    Type_          string  `json:"type"`
    Country        string  `json:"country"`
    GeoPosition    *GeoPosition `json:"geoPosition"`
    InEurope       bool    `json:"inEurope"`
    CountryCode    string  `json:"countryCode"`
    CoreCountry    bool    `json:"coreCountry"`
    Distance       *float64 `json:"distance"`
}


type GeoPosition struct{
	Latitude float64
	Longitude float64
}


func RandStringBytesMaskImpr(n int, gen *customRandomGenerator) string {
	b := make([]byte, n)
    // A rand.Int63() generates 63 random bits, enough for ${letterIdxMax} letters, so to make it faster,
    // my random strings will have length of max 10

    for i, cache, remain := n-1, gen.generateRandomInt63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = gen.generateRandomInt63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }
    return  *(*string)(unsafe.Pointer(&b))
}

func createRandomGeoPosition() *GeoPosition{
    return &GeoPosition{
        Latitude: rand.Float64() * 100, 
        Longitude: rand.Float64() * 100,
    }
}



func CreateRandomJson(id int64) *ExampleJson {
    var name = RandStringBytesMaskImpr(7, gen)
    var country = RandStringBytesMaskImpr(9, gen)
    var iata_airport_code *string 
    var distance *float64 
    var key *int64 

    var check =  rand.Intn(4)
    if check == 0 { 
        d := rand.Float64() * 10000  
        distance = &d
    } 

    if  check % 2 == 0 {  
        k := gen.generateRandomInt63() 
        key = &k
    }

    if check != 3 {
        i := RandStringBytesMaskImpr(3, gen)
        iata_airport_code = &i
    }

    return &ExampleJson{
        Type: "Position",
        Id: id,  
        Name: name,
        Key: key,
        Country: country,
        FullName: name + "," + country, 
        IataAirportCode: iata_airport_code,
        Type_: RandStringBytesMaskImpr(3, gen),
        GeoPosition: createRandomGeoPosition(),
        LocationID:  gen.generateRandomInt63(),
        InEurope: !(id % 2 == 0),
        CountryCode:  RandStringBytesMaskImpr(2, gen),
        CoreCountry: !(id % 3 == 0),
        Distance: distance,
    }
}