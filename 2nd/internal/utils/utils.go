package utils

import (
	"net/http"
	"strconv"
	"strings"
)

// Server utils
var minRecordsNumber int64 = 0
var maxRecordsNumber int64 = 1000000

var PossibleQueryParams = map[string]struct{}{
	"type":        {},
	"id":          {},
	"key":         {},
	"name":        {},
	"fullName":    {},
	"locationId":  {},
	"iata":        {},
	"type_":       {},
	"country":     {},
	"lat":         {},
	"long":        {},
	"inEurope":    {},
	"countryCode": {},
	"coreCountry": {},
	"dist":        {},
}

// Validation utils

func ValidateUrlFirstEndpoint(url string) (int64, bool) {
	parts := strings.Split(url, "/")

	if len(parts) != 3 {
		return 0, false
	}

	count, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || count < minRecordsNumber || count > maxRecordsNumber {
		return 0, false
	}

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

var ErrorWritingCsvHeaders = CreateNewApiError(http.StatusInternalServerError, "Error writing header to csv. ")
var ErrorWritingCsvRecord = CreateNewApiError(http.StatusInternalServerError, "Error writing record to csv. ")
var ErrorReadingResponseBody = CreateNewApiError(http.StatusInternalServerError, "Error reading response body. ")
var ErrorCastingApiError = CreateNewApiError(http.StatusInternalServerError, "An unexpected server error occurred when casting error to ApiError. ")
