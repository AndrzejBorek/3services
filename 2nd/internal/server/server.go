package server

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"

	types "github.com/AndrzejBorek/3services/1st/pkg/types"
	client "github.com/AndrzejBorek/3services/2nd/internal/client"
	utils "github.com/AndrzejBorek/3services/2nd/internal/utils"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Handled %s %s in %v", r.Method, r.URL.Path, duration)
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Exact-Time", fmt.Sprint(time.Now().Unix()))
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func writeCSV(w http.ResponseWriter, statusCode int, jsonList []types.ExampleJson) error {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Exact-Time", fmt.Sprint(time.Now().Unix()))

	w.WriteHeader(statusCode)
	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write([]string{"Type", "Id", "Name", "Type_", "Latitude", "Longitude"}); err != nil {
		return utils.ErrorWritingCsvHeaders
	}

	for _, record := range jsonList {

		stringRecord := []string{
			record.Type,
			strconv.FormatInt(int64(record.Id), 10),
			record.Name,
			record.Type_,
			strconv.FormatFloat(record.GeoPosition.Latitude, 'f', 2, 64),
			strconv.FormatFloat(record.GeoPosition.Longitude, 'f', 2, 64),
		}

		if err := writer.Write(stringRecord); err != nil {
			return utils.ErrorWritingCsvRecord
		}
	}

	return nil
}

func FirstEndpointHandler(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return utils.ErrorGenericMethodNotAllowed
	}

	size, valid := utils.ValidateUrlFirstEndpoint(r.URL.Path)

	if !valid {
		return utils.ErrorGenericInvalidRequest
	}

	data, err := client.GetListOfJsons(r.Context(), size)

	if err != nil {
		if apiErr, ok := err.(utils.APIError); ok {
			return apiErr
		} else {
			log.Print(utils.ErrorCastingApiError.Msg + err.Error())
			return utils.ErrorGenericInternalServerError
		}
	}

	notFullyRandomJsonIndex := utils.RandRange(0, size-1)
	notFullyRandomJson := data[notFullyRandomJsonIndex]

	err = client.RedisSet(&client.RedisClient, "key", notFullyRandomJson)
	if err != nil {
		log.Println("Error when inserting data into redis. ", err)
		return utils.ErrorGenericInternalServerError
	}

	return writeCSV(w, http.StatusOK, data)

}

func SecondEndpointHandler(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return utils.ErrorGenericMethodNotAllowed
	}

	queryParams := r.URL.Query()

	for key := range queryParams {
		lowerKey := strings.ToLower(key)
		if _, valid := utils.PossibleQueryParams[lowerKey]; !valid {
			return utils.ErrorGenericInvalidRequest
		}
	}

	var data map[string]interface{}
	err := client.RedisGet(&client.RedisClient, "key", &data)
	if err == redis.Nil {
		listData, err := client.GetListOfJsons(r.Context(), 1)
		if err != nil {
			if apiErr, ok := err.(utils.APIError); ok {
				log.Println(apiErr)
				return apiErr
			}
			log.Print(utils.ErrorCastingApiError.Msg + err.Error())
			return utils.ErrorGenericInternalServerError
		}

		if len(listData) > 0 {
			data = listData[0].ConvertToMap()
		} else {
			log.Println("Received empty data from GetListOfJsons.")
			return utils.ErrorGenericInternalServerError
		}
	} else if err != nil {
		log.Printf("Redis error %s. ", err)
		return utils.ErrorGenericInternalServerError
	}

	returnJson := make(map[string]interface{})
	for key := range queryParams {
		if value, exists := data[key]; exists {
			returnJson[key] = value
		}
	}
	return writeJSON(w, http.StatusOK, returnJson)
}

func ThirdEndpointHandler(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return utils.ErrorGenericMethodNotAllowed
	}

	return nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func MakeHandler(handler apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if e, ok := err.(utils.APIError); ok {
				writeJSON(w, e.Status, e)
			} else {
				writeJSON(w, http.StatusInternalServerError, errors.New("Internal error."))
			}
		}
	}
}
