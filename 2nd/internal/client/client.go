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
var redisHost = os.Getenv("REDIS_HOST")
var redisPort = os.Getenv("REDIS_PORT")
var redisAddr = fmt.Sprintf("%s:%s", redisHost, redisPort)

var redisExpirationTime time.Duration = 30000000000

var RedisClient = createNewRedisClient(redisAddr, "", 0)
var httpClient = createNewHttpClient(30)

// Http client

func createNewHttpClient(timeout int64) http.Client {
	log.Printf("Http client has been created in service 2.")
	return http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func GetListOfJsons(ctx context.Context, size int) ([]types.ExampleJson, error) {
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

func createNewRedisClient(redisAddr string, redisPasswd string, redisDb int) redis.Client {
	log.Printf("Redis client has been created in service 2.")
	return *redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPasswd,
		DB:       redisDb,
	})
}

func RedisSet(redisClient *redis.Client, key string, value types.ExampleJson) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mapValue := value.ConvertToMap()
	jsonValue, err := json.Marshal(mapValue)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, jsonValue, redisExpirationTime).Err()
}

func RedisGet(redisClient *redis.Client, key string, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p := redisClient.Get(ctx, key)

	if err := p.Err(); err != nil {
		if err == redis.Nil {
			log.Printf("Key %s does not exist.", key)
		}
		return err
	}

	val, err := p.Bytes()
	if err != nil {
		log.Printf("Failed to change data into bytes. Err: %s", err)
		return err
	}

	if err := json.Unmarshal(val, dest); err != nil {
		log.Printf("Failed to unmarshal: %s", err)
		return err
	}

	return nil
}
