package utils

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	types "github.com/AndrzejBorek/3services/1st/pkg/types"
)

// String utils

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index (Cause 2^6 = 64, and since there are 52 letters in alphabet I need 6 of them)
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var maxUrlPaths = 2

type customRandomGenerator struct {
	src  rand.Source
	lock sync.Mutex //Since rand.Int63 is not safe for concurrent use, I had to create this struct
}

func createNewCustomRandomGenerator(seed int64) *customRandomGenerator {
	return &customRandomGenerator{
		src: rand.NewSource(seed),
	}
}

var gen = createNewCustomRandomGenerator(time.Now().UnixNano())

func (gen *customRandomGenerator) generateRandomInt63() int64 {
	gen.lock.Lock()
	defer gen.lock.Unlock()
	return gen.src.Int63()
}

func randStringBytesMaskImpr(n int, gen *customRandomGenerator) string {
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
	return *(*string)(unsafe.Pointer(&b))
}

// Types utils

func createRandomGeoPosition() types.GeoPosition {
	return types.GeoPosition{
		Latitude:  rand.Float64() * 100,
		Longitude: rand.Float64() * 100,
	}
}

func createRandomJson(id int) types.ExampleJson {
	var name = randStringBytesMaskImpr(7, gen)
	var country = randStringBytesMaskImpr(9, gen)
	var iata_airport_code string
	var distance float64
	var key int64

	var check = rand.Intn(4)
	if check != 0 {
		distance = rand.Float64() * 10000
	}

	if check%2 == 0 {
		key = gen.generateRandomInt63()
	}
	// Maybe use Country[:3] as Iata code so it would be much faster without calling randStringBytesMaskImpr
	if check != 3 {
		iata_airport_code = randStringBytesMaskImpr(3, gen)
	}

	// Since in task requirements there is requirement saying that values are random I used totally random values
	// even though iata code is connected to country and should not be random. In real case it would be probably anyway taken
	// from some data storage.
	return types.ExampleJson{
		Type:            "Position",
		Id:              id,
		Name:            name,
		Key:             key,
		Country:         country,
		FullName:        name + "," + country,
		IataAirportCode: iata_airport_code,
		Type_:           "location",
		GeoPosition:     createRandomGeoPosition(),
		LocationID:      gen.generateRandomInt63(),
		InEurope:        !(id%2 == 0),
		CountryCode:     randStringBytesMaskImpr(2, gen),
		CoreCountry:     !(id%3 == 0),
		Distance:        distance,
	}
}

func GenerateRandomJsons(count int) (result []types.ExampleJson) {
	result = make([]types.ExampleJson, count)
	for i := 0; i < int(count); i++ {
		result[i] = createRandomJson(i + 1)
	}
	return
}

// Server utils

func ValidateUrl(url string) (int, bool) {
	trimmedUrl := strings.Trim(url, "/")
	parts := strings.Split(trimmedUrl, "/")

	if len(parts) != maxUrlPaths {
		return 0, false
	}

	count64, err := strconv.ParseInt(parts[maxUrlPaths-1], 10, 64)

	if err != nil {
		return 0, false
	}

	if count64 < 0 || count64 > 1000000 {
		return 0, false
	}
	count := int(count64)

	return count, true
}

// Errors

type APIError struct {
	Status int
	Msg    string
}

func CreateNewApiError(status int, message string) APIError {
	return APIError{
		Status: status,
		Msg:    message,
	}
}

func (e APIError) Error() string {
	return e.Msg
}

var ErrorGenericMethodNotAllowed = CreateNewApiError(http.StatusMethodNotAllowed, "Method not allowed.")
var ErrorGenericInvalidRequest = CreateNewApiError(http.StatusBadRequest, "Invalid request.")
var ErrorGenericInternalServerError = CreateNewApiError(http.StatusInternalServerError, "An unexpected server error occurred.")
