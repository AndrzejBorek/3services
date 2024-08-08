package server


import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

)


var handler = LoggingMiddleware(MakeHandler(GenerateJsonHandler))


func TestGenerateJSONHandlerShouldFailWrongMethod(t *testing.T) {

    req := httptest.NewRequest(http.MethodPost, "/generate/json/10", nil)
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

    req := httptest.NewRequest(http.MethodGet, "/generate/json/10/10", nil)
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

    req := httptest.NewRequest(http.MethodGet, "/generate/json/10000001", nil)
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

    req := httptest.NewRequest(http.MethodGet, "/generate/json/-1", nil)
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

    req := httptest.NewRequest(http.MethodGet, "/generate/json/test", nil)
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



