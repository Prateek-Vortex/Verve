package restclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type RestClient interface {
	Get(path string) (*http.Response, error)
	Post(path string, body interface{}) (*http.Response, error)
	Put(path string, body interface{}) (*http.Response, error)
	Delete(path string) (*http.Response, error)
}

// RestClient represents a simple HTTP client
type RestHttpClient struct {
	httpClient *http.Client
	headers    map[string]string
}

// NewRestClient creates a new REST client instance
func NewRestClient() RestClient {
	return &RestHttpClient{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// Get performs a GET request
func (c *RestHttpClient) Get(path string) (*http.Response, error) {
	return c.doRequest("GET", path, nil)
}

// Post performs a POST request
func (c *RestHttpClient) Post(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("POST", path, body)
}

// Put performs a PUT request
func (c *RestHttpClient) Put(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("PUT", path, body)
}

// Delete performs a DELETE request
func (c *RestHttpClient) Delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil)
}

// SetHeader sets a custom header
func (c *RestHttpClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// doRequest performs the HTTP request
func (c *RestHttpClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
	url := path
	var reqBody *bytes.Buffer

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	return c.httpClient.Do(req)
}
