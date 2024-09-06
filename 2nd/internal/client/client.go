package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"

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

var httpClient = createNewHttpClient(30)

func GetListOfJsons(ctx context.Context, size int64) ([]types.ExampleJson, error) {
	// Here some additional algorithm could be used - instead of doing one request with lets say milion jsons (param size = 1m) I could use goroutines to see
	// if maybe 10 concurrent requests with 100k would be faster.
	var url = fmt.Sprintf("%s%d", service1Url, size)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err.Error())
		return nil, utils.ErrorGenericInternalServerError
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error when calling service1: %s", err.Error())
		return nil, utils.ErrorGenericInternalServerError
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body from service1: %v", err)
			return nil, utils.ErrorGenericInternalServerError
		}
		resp.Body.Close()
		log.Printf("Error in response from service1: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
		return nil, utils.ErrorGenericBadGateway
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Printf(utils.ErrorReadingResponseBody.Msg + err.Error())
		return nil, utils.ErrorGenericInternalServerError
	}

	var jsonList []types.ExampleJson
	if err := json.Unmarshal(bodyBytes, &jsonList); err != nil {
		log.Printf("Error unmarshaling JSON: %s, %s", err.Error(), utils.ErrorWritingCsvRecord)
		return nil, utils.ErrorGenericInternalServerError
	}
	return jsonList, nil

}

// Redis client
var redisHost = os.Getenv("REDIS_HOST")
var redisPort = os.Getenv("REDIS_PORT")
var redisAddr = fmt.Sprintf("%s:%s", redisHost, redisPort)

func createNewRedisClient(redisAddr string, redisPasswd string, redisDb int) redis.Client {
	log.Printf("Redis client has been created in service 2.")
	return *redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPasswd,
		DB:       redisDb,
	})
}

var RedisClient = createNewRedisClient(redisAddr, "", 0)
