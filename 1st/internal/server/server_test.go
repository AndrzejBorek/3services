package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var handler = LoggingMiddleware(MakeHandler(GenerateJsonHandler))
var baseURL = "/generate"

func TestGenerateJSONHandlerShouldFailWrongMethod(t *testing.T) {

	url := baseURL + "/10"
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	_, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %v got %v", http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestGenerateJSONHandlerShouldFailTooManyURLParts(t *testing.T) {

	url := baseURL + "/10/10"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	_, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestGenerateJSONHandlerShouldFailTooLargeNumber(t *testing.T) {
	url := baseURL + "/10000001"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	_, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestGenerateJSONHandlerShouldFailNonPositiveInteger(t *testing.T) {
	url := baseURL + "/-1"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	_, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestGenerateJSONHandlerShouldFailNonValidInteger(t *testing.T) {
	url := baseURL + "/test"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	_, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v got %v", http.StatusBadRequest, res.StatusCode)
	}
}
