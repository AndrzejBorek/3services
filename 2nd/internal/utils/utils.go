package utils

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
)

// Server utils
var minRecordsNumber int64 = 0
var maxRecordsNumber int64 = 1000000
var maxUrlPathsFirstEndpoint = 2

var PossibleQueryParams = map[string]struct{}{
	"type":        {},
	"id":          {},
	"key":         {},
	"name":        {},
	"fullname":    {},
	"locationid":  {},
	"iata":        {},
	"type_":       {},
	"country":     {},
	"latitude":    {},
	"longitude":   {},
	"ineurope":    {},
	"countrycode": {},
	"corecountry": {},
	"distance":    {},
}

// Validation utils

func ValidateUrlFirstEndpoint(url string) (int, bool) {
	trimmedUrl := strings.Trim(url, "/")
	parts := strings.Split(trimmedUrl, "/")

	if len(parts) != maxUrlPathsFirstEndpoint {
		return 0, false
	}

	count64, err := strconv.ParseInt(parts[maxUrlPathsFirstEndpoint-1], 10, 64)
	if err != nil || count64 < minRecordsNumber || count64 > maxRecordsNumber {
		return 0, false
	}

	count := int(count64)
	return count, true

}

// Error utils

type APIError struct {
	Status int
	Msg    string
}

func (e APIError) Error() string {
	return e.Msg
}

func CreateNewApiError(status int, message string) APIError {
	return APIError{
		Status: status,
		Msg:    message,
	}
}

var ErrorGenericMethodNotAllowed = CreateNewApiError(http.StatusMethodNotAllowed, "Method not allowed.")
var ErrorGenericInvalidRequest = CreateNewApiError(http.StatusBadRequest, "Invalid request.")
var ErrorGenericInternalServerError = CreateNewApiError(http.StatusInternalServerError, "An unexpected server error occurred.")
var ErrorGenericBadGateway = CreateNewApiError(http.StatusBadGateway, "Bad gateway.")

var ErrorWritingCsvHeaders = CreateNewApiError(http.StatusInternalServerError, "Error writing header to csv. ")
var ErrorWritingCsvRecord = CreateNewApiError(http.StatusInternalServerError, "Error writing record to csv. ")
var ErrorReadingResponseBody = CreateNewApiError(http.StatusInternalServerError, "Error reading response body. ")
var ErrorCastingApiError = CreateNewApiError(http.StatusInternalServerError, "An unexpected server error occurred when casting error to ApiError. ")

//

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}
