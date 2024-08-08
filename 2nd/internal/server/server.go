package server

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AndrzejBorek/3services/1st/pkg/types"
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

func writeCSV(w http.ResponseWriter, statusCode int, jsonList []*types.ExampleJson) error {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Exact-Time", fmt.Sprint(time.Now().Unix()))
	w.WriteHeader(statusCode)

	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write([]string{"Type", "Id", "Name", "Type_", "Latitude", "Longitude"}); err != nil {
		return utils.APIError{
			Status: http.StatusInternalServerError,
			Msg:    "Error writing header to csv.",
		}
	}

	for _, record := range jsonList {
		if record == nil {
			continue
		}
		stringRecord := []string{
			record.Type,
			strconv.FormatInt(record.Id, 10),
			record.Name,
			record.Type_,
			strconv.FormatFloat(record.GeoPosition.Latitude, 'f', 2, 64),
			strconv.FormatFloat(record.GeoPosition.Longitude, 'f', 2, 64),
		}
		if err := writer.Write(stringRecord); err != nil {
			return utils.APIError{
				Status: http.StatusInternalServerError,
				Msg:    "Error writing record to csv.",
			}
		}
	}
	return nil
}

func FirstEndpointHandler(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return utils.APIError{
			Status: http.StatusMethodNotAllowed,
			Msg:    "Method not allowed.",
		}
	}

	data, err := client.GetListOfJsons(10000000)
	if err != nil {
		return utils.APIError{
			Status: err.Status,
			Msg:    err.Msg,
		}
	}
	return writeCSV(w, http.StatusOK, data)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

// I could move checking request method here and change name for MakeGetHandler, since I know I won't need to deal with POST method in this task
func MakeHandler(handler apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if e, ok := err.(utils.APIError); ok {
				writeJSON(w, e.Status, e)
			} else {

				writeJSON(w, http.StatusInternalServerError, errors.New("internal error"))
			}
		}
	}
}
