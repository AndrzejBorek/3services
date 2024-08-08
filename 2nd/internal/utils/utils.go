package utils

import (
	"net/http"
	"strconv"
	"strings"

	types "github.com/AndrzejBorek/3services/1st/pkg/types"
)

// Server utils

type APIError struct {
	Status int
	Msg    string
}

func (e APIError) Error() string {
	return e.Msg
}

func ValidateUrl(r *http.Request) (int64, bool) {
	// This needs to be changed to csv handler requirements
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 2 {
		return 0, false
	}

	count, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, false
	}

	if count < 0 || count > 1000000 {
		return 0, false
	}

	return count, true
}

func firstEndpointUtil(jsonList []*types.ExampleJson) {

}

func secondEndpointUtil(jsonList []*types.ExampleJson) {

}

func thirdEndpointUtil(jsonList []*types.ExampleJson) {

}
