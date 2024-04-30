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

type APIError struct {
	Status int
	Msg    string
}

func (e APIError) Error() string {
	return e.Msg
}

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

func GenerateJsonHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return APIError{
			Status: http.StatusMethodNotAllowed,
			Msg:    "Method not allowed",
		}
	}

	count, valid := utils.ValidateUrl(r)
	if !valid {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Invalid request",
		}
	}

	return writeJSON(w, http.StatusOK, utils.GenerateRandomJsons(count))
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func MakeHandler(handler apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if e, ok := err.(APIError); ok { // Checking if err from handler is APIError type
				writeJSON(w, e.Status, e)
			} else {
				writeJSON(w, http.StatusInternalServerError, errors.New("internal error"))
			}
		}
	}
}
