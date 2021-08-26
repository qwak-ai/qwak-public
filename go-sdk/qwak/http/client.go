package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

const (
	MAXIMUM_RETRY_ATTEMPTS = 5
	RETRY_DELAY            = 500 * time.Millisecond
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

func GetDefaultHttpClient() HttpClient {
	return http.DefaultClient
}

func DoRequest(client HttpClient, request *http.Request) (responseBody []byte, httpCode int, err error) {

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, (fmt.Errorf("an error occured on authentication request: %v", err.Error()))
	}
	defer response.Body.Close()
	
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	if err != nil {
		return nil, response.StatusCode, errors.New("failed decode authentication request")
	}

	return body, response.StatusCode, nil

}

func DoRequestWithRetry(client HttpClient, request *http.Request) (responseBody []byte, statusCode int, err error) {
	retryAttempt := 0
	lastHttpCode := 500
	var lastErr error = nil
	var body []byte = nil

	for retryAttempt < MAXIMUM_RETRY_ATTEMPTS && lastHttpCode >= 500 && lastErr == nil {
		body, lastHttpCode, lastErr = DoRequest(client, request)

		if lastHttpCode >= 500 && err == nil {
			time.Sleep(RETRY_DELAY * time.Duration(int(math.Pow(2, float64(retryAttempt)))))
		}
	}

	return body, lastHttpCode, nil
}
