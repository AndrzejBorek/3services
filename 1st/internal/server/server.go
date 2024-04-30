package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	util "github.com/AndrzejBorek/3services/1st/internal/utils"
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

func validateUrl(r *http.Request) (int64, bool) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
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

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Exact-Time", fmt.Sprint(time.Now().Unix()))
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func generateRandomJsons(count int64) (result []*util.ExampleJson) {
	result = make([]*util.ExampleJson, count)
	for i := 0; i < int(count); i++ {
		result[i] = util.CreateRandomJson(int64(i + 1))
	}
	return
}

func GenerateJsonHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return APIError{
			Status: http.StatusMethodNotAllowed,
			Msg:    "Method not allowed",
		}
	}

	count, valid := validateUrl(r)
	if !valid {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return APIError{
			Status: http.StatusBadRequest,
			Msg:    "Invalid request",
		}
	}

	return writeJSON(w, http.StatusOK, generateRandomJsons(count))
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
