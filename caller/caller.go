package caller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const EP_AUTH = "/v1/auth"
const EP_USER_LIST = "/v1/users"

var baseUrl string
var apiKey string
var client *http.Client

var InvalidAuthError = errors.New("Auth token is invalid.")

func Post(ep string, data interface{}) ([]byte, error) {
	jsonData, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", buildUrl(ep), bytes.NewReader(jsonData))

	request.Header.Add("Content-Type", "application/json")

	if apiKey != "" {
		request.Header.Add("Authorization", "Bearer "+apiKey)
	}

	resp, respErr := client.Do(request)

	if respErr != nil {
		return nil, respErr
	}

	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	return result, nil
}

func Get(ep string) ([]byte, error) {
	request, _ := http.NewRequest("GET", buildUrl(ep), nil)

	request.Header.Add("Content-Type", "application/json")

	if apiKey != "" {
		request.Header.Add("Authorization", "Bearer "+apiKey)
	}

	resp, respErr := client.Do(request)

	if respErr != nil {
		return nil, respErr
	}

	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)

	return result, nil
}

func SetAuth(key string) {
	apiKey = key
}

func SetBaseUrl(url string) {
	baseUrl = url
}

func buildUrl(ep string) string {
	return baseUrl + ep
}

func init() {
	client = &http.Client{}
}
