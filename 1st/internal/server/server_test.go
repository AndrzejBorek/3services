package server


import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
)

// func TestGenerateJSONHandlerShouldPass(t *testing.T) {
//     req := httptest.NewRequest(http.MethodGet, "/generate/json/10", nil)
//     w := httptest.NewRecorder()
//     GenerateJsonHandler(w, req)
//     res := w.Result()
//     defer res.Body.Close()
//     data, err := io.ReadAll(res.Body)

//     if err != nil {
//         t.Errorf("expected error to be nil got %v", err)
//     }
//     if string(data) != "ABC" {
//         t.Errorf("expected ABC got %v", string(data))
//     }
// }




func testGetRequest(method string, url string, body io.Reader, t *testing.T){
	req := httptest.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	GenerateJsonHandler(w, req)

	res := w.Result()
    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }

	if res.StatusCode != http.StatusBadRequest {
        t.Errorf("expected status %v got %v", http.StatusBadRequest, res.StatusCode)
    }

    if string(data) != "Invalid request\n" { 
        t.Errorf("expected 'Invalid request' got %v", string(data))
    }
}


func TestGenerateJSONHandlerShouldFailTooLargeNumber(t *testing.T) {
	testGetRequest(http.MethodGet, "/generate/json/10000001", nil, t)
}

func TestGenerateJSONHandlerShouldFailTooManyURLParts(t *testing.T) {
	testGetRequest(http.MethodGet, "/generate/json/10/10", nil, t)
}

func TestGenerateJSONHandlerShouldFailNonPositiveInteger(t *testing.T) {
	testGetRequest(http.MethodGet, "/generate/json/-10", nil, t)
}

func TestGenerateJSONHandlerShouldFailNonValidInteger(t *testing.T) {
	testGetRequest(http.MethodGet, "/generate/json/test", nil, t)
}


// func TestGenerateJSONHandlerShouldFailWrongMethod(t *testing.T) {
// 	testRequest(http.MethodPost,  "/generate/json/10", nil, t)
// }


