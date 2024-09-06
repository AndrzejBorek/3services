package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	utils "github.com/AndrzejBorek/3services/1st/internal/utils"
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

	encodedData, err := json.Marshal(data)

	if err != nil {
		return utils.ErrorGenericInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Exact-Time", fmt.Sprint(time.Now().Unix()))
	w.WriteHeader(statusCode)
	_, err = w.Write(encodedData)
	return err
}

func GenerateJsonHandler(w http.ResponseWriter, r *http.Request) error {

	if r.Method != http.MethodGet {
		return utils.ErrorGenericMethodNotAllowed
	}

	count, valid := utils.ValidateUrl(r.URL.Path)

	if !valid {
		log.Printf("Invalid url paths.")
		return utils.ErrorGenericInvalidRequest
	}

	return writeJSON(w, http.StatusOK, utils.GenerateRandomJsons(count))
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func MakeHandler(handler apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if e, ok := err.(utils.APIError); ok { // Checking if err from handler is APIError type
				writeJSON(w, e.Status, e)
			} else {
				writeJSON(w, http.StatusInternalServerError, errors.New("Internal error."))
			}
		}
	}
}
