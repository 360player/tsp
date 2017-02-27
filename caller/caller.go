package caller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const EP_AUTH = "/v1/auth"
const EP_USER_LIST = "/v1/users"

var baseUrl string
var apiKey string
var client *http.Client

var InvalidAuthError = errors.New("Auth token is invalid.")

type QueryParam struct {
	Key   string
	Value string
}

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

func Get(ep string, queryParams ...QueryParam) ([]byte, error) {
	request, _ := http.NewRequest("GET", buildUrl(ep), nil)

	request.Header.Add("Content-Type", "application/json")

	if apiKey != "" {
		request.Header.Add("Authorization", "Bearer "+apiKey)
	}

	if len(queryParams) > 0 {
		q := request.URL.Query()

		for _, param := range queryParams {
			q.Add(param.Key, param.Value)
		}

		request.URL.RawQuery = q.Encode()
	}

	fmt.Println(request.URL.String())

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
