package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GTHttpClient struct {
	client    *http.Client
	cookieVal string
}

func NewGTHttpClient() *GTHttpClient {
	return &GTHttpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *GTHttpClient) reRequest(options *http.Request) ([]byte, error) {
	resp, err := c.client.Do(options)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[ERROR] reRequest. Request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *GTHttpClient) Request(method string, host string, path string, qs url.Values, agent *http.Transport) ([]byte, error) {
	urlStr := fmt.Sprintf("https://%s%s?%s", host, path, qs.Encode())

	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return nil, err
	}

	if agent != nil {
		c.client.Transport = agent
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")

	if c.cookieVal != "" {
		req.Header.Set("Cookie", c.cookieVal)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests && resp.Header.Get("Set-Cookie") != "" {
		cookieVal := strings.Split(resp.Header.Get("Set-Cookie"), ";")[0]
		c.cookieVal = cookieVal
		req.Header.Set("Cookie", c.cookieVal)
		return c.reRequest(req)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return body, nil
}

// Backoff, rate-limit, circuit breaker
// func (c *GTHttpClient) Request(method string, host string, path string, qs url.Values, agent *http.Transport) ([]byte, error) {
// 	urlStr := fmt.Sprintf("https://%s%s?%s", host, path, qs.Encode())

// 	req, err := http.NewRequest(method, urlStr, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if agent != nil {
// 		c.client.Transport = agent
// 	}

// 	if c.cookieVal != "" {
// 		req.Header.Set("Cookie", c.cookieVal)
// 	}

// 	var resp *http.Response
// 	var body []byte

// 	retryCount := 0
// 	maxRetries := 5
// 	baseDelay := 1 * time.Second
// 	maxDelay := 10 * time.Second

// 	for retryCount < maxRetries {
// 		resp, err = c.client.Do(req)
// 		if err != nil {
// 			retryCount++
// 			delay := baseDelay * (2 << retryCount)
// 			if delay > maxDelay {
// 				delay = maxDelay
// 			}
// 			time.Sleep(delay)
// 			continue
// 		}
// 		defer resp.Body.Close()

// 		body, err = io.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if resp.StatusCode == http.StatusTooManyRequests && resp.Header.Get("Set-Cookie") != "" {
// 			cookieVal := strings.Split(resp.Header.Get("Set-Cookie"), ";")[0]
// 			c.cookieVal = cookieVal
// 			req.Header.Set("Cookie", c.cookieVal)
// 			continue
// 		}

// 		return body, nil
// 	}

// 	return nil, errors.New("Max retries exceeded")
// }
