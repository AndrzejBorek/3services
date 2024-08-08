package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AndrzejBorek/3services/1st/pkg/types"

	utils "github.com/AndrzejBorek/3services/2nd/internal/utils"
)

var service1Url = "http://service1:8080/generate/json/"

func createNewClient(timeout int64) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

var client = createNewClient(10)

func GetListOfJsons(size int32) (jsonList []*types.ExampleJson, utilErr *utils.APIError) {
	var url = fmt.Sprintf("%s%d", service1Url, size)
	resp, err := client.Get(url)

	if err != nil {
		return nil, &utils.APIError{
			Status: resp.StatusCode,
			Msg:    "Error when calling service1:" + err.Error(),
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, &utils.APIError{
			Status: http.StatusInternalServerError,
			Msg:    "Error reading response body:" + err.Error(),
		}
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr utils.APIError
		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			fmt.Println("Error unmarshaling error response:", err.Error())
			return nil, &utils.APIError{Status: resp.StatusCode, Msg: "Error from service."}
		}
		return nil, &apiErr
	}

	err = json.Unmarshal(bodyBytes, &jsonList)

	if err != nil {
		return nil, &utils.APIError{
			Status: http.StatusBadRequest,
			Msg:    "Error unmarshaling JSON:" + err.Error(),
		}
	}

	return jsonList, nil
}
