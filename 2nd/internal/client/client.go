package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AndrzejBorek/3services/1st/pkg/types"

	utils "github.com/AndrzejBorek/3services/2nd/internal/utils"
)

var service1Url = os.Getenv("service1Url")

// Http client

func createNewHttpClient(timeout int64) http.Client {
	log.Printf("Http client has been created in service 2.")
	return http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

var client = createNewHttpClient(30)

func GetListOfJsons(size int64) ([]types.ExampleJson, error) {

	var url = fmt.Sprintf("%s%d", service1Url, size)

	resp, err := client.Get(url)

	// Here some additional algorithm could be used - instead of doing one request with lets say milion jsons (param size = 1m) I could use goroutines to see
	// if maybe 10 concurrent requests with 100k would be faster.

	// Look at all those errors to see if they should not be more generic for client
	if err != nil {
		log.Printf("Error when calling service1. " + err.Error())
		return nil, utils.ErrorGenericInternalServerError
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Printf(utils.ErrorReadingResponseBody.Msg + err.Error())
		return nil, utils.ErrorReadingResponseBody
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var jsonList []types.ExampleJson
		if err := json.Unmarshal(bodyBytes, &jsonList); err != nil {
			log.Printf("Error unmarshaling JSON: %s", err.Error())
			return nil, utils.ErrorGenericInternalServerError
		}
		return jsonList, nil
	case http.StatusNotFound:
		log.Print("case http.StatusNotFound. Should return ErrorGenericInternalServerError")
		log.Print(resp.StatusCode)
		log.Print("......................")
		return nil, utils.ErrorGenericInternalServerError

	default:
		var apiErr utils.APIError
		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			log.Printf("Error unmarshaling API error response: %s", err.Error())
			return nil, utils.ErrorGenericInternalServerError
		}
		log.Printf("API Error: %s: %d", apiErr.Msg, apiErr.Status)
		return nil, apiErr
	}
}

// Redis client
