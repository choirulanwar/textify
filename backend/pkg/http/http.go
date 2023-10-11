package http

import (
	"bytes"
	"fmt"

	"io"
	net_http "net/http"
)

func request(requestType, url, token string, cookies string, payload []byte) ([]byte, error) {
	client := &net_http.Client{}

	var request *net_http.Request
	if payload != nil {
		requestBody := bytes.NewReader(payload)
		request, _ = net_http.NewRequest(requestType, url, requestBody)
	} else {
		request, _ = net_http.NewRequest(requestType, url, nil)
	}

	if token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	if cookies != "" {
		request.Header.Set("Cookie", cookies)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if response.StatusCode != net_http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", response.Status)
	}

	body, _ := io.ReadAll(response.Body)
	return body, nil
}

func MakeGetRequest(url string, token string, cookies string) ([]byte, error) {
	return request("GET", url, token, cookies, nil)
}

func MakePostRequest(url, token string, cookies string, payload []byte) ([]byte, error) {
	return request("POST", url, token, cookies, payload)
}
