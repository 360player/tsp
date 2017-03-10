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
const EP_RATING_ITEMS = "/v1/rating-items"
const EP_POSITIONS = "/v1/positions"
const EP_USER = "/v1/users/%d"
const EP_GROUP = "/v1/groups/%d"
const EP_USER_RATINGS = "/v1/users/%d/ratings"

var baseUrl string
var authToken string
var apiKey string
var client *http.Client

var InvalidAuthError = errors.New("Auth token is invalid.")

type QueryParam struct {
	Key   string
	Value string
}

func addHeaders(request *http.Request) {
	request.Header.Add("Content-Type", "application/json")

	if authToken != "" {
		request.Header.Add("Authorization", "Bearer "+authToken)
	}

	if apiKey != "" {
		request.Header.Add("X-API-Key", apiKey)
	}
}

func Post(ep string, data interface{}) ([]byte, error) {
	jsonData, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", buildUrl(ep), bytes.NewReader(jsonData))
	addHeaders(request)

	resp, respErr := client.Do(request)

	if respErr != nil {
		return nil, respErr
	}

	if resp.StatusCode == 401 {
		return nil, InvalidAuthError
	}

	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	return result, nil
}

func Get(ep string, queryParams ...QueryParam) ([]byte, error) {
	request, _ := http.NewRequest("GET", buildUrl(ep), nil)
	addHeaders(request)

	if len(queryParams) > 0 {
		q := request.URL.Query()

		for _, param := range queryParams {
			q.Add(param.Key, param.Value)
		}

		request.URL.RawQuery = q.Encode()
	}

	resp, respErr := client.Do(request)

	if respErr != nil {
		return nil, respErr
	}

	if resp.StatusCode == 401 {
		return nil, InvalidAuthError
	}

	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return result, errors.New(resp.Status)
	}

	return result, nil
}

func SetAuth(token string) {
	authToken = token
}

func SetApiKey(key string) {
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
