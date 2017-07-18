package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Error represents a error from the bitbucket api.
type Error struct {
	APIError struct {
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
	Type       string `json:"type,omitempty"`
	StatusCode int
	Endpoint   string
}

func (e Error) Error() string {
	return fmt.Sprintf("API Error: %d %s %s", e.StatusCode, e.Endpoint, e.APIError.Message)
}

const (
	// BitbucketEndpoint is the fqdn used to talk to bitbucket
	BitbucketEndpoint string = "https://api.bitbucket.org/"
)

type BitbucketClient struct {
	Username   string
	Password   string
	HTTPClient *http.Client
}

func (c *BitbucketClient) Do(method, endpoint, contentType string, payload *bytes.Buffer) (*http.Response, error) {
	absoluteendpoint := BitbucketEndpoint + endpoint
	log.Printf("[DEBUG] Sending request to %s %s", method, absoluteendpoint)

	var bodyreader io.Reader

	if payload != nil {
		log.Printf("[DEBUG] With payload %s", payload.String())
		bodyreader = payload
	}

	req, err := http.NewRequest(method, absoluteendpoint, bodyreader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	if payload != nil {
		// Can cause bad request when putting default reviews if set.
		req.Header.Add("Content-Type", contentType)
	}

	req.Close = true

	resp, err := c.HTTPClient.Do(req)
	log.Printf("[DEBUG] Resp: %v Err: %v", resp, err)
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		apiError := Error{
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Resp Body: %s", string(body))

		err = json.Unmarshal(body, &apiError)
		if err != nil {
			apiError.APIError.Message = string(body)
		}

		return resp, error(apiError)

	}
	return resp, err
}

func (c *BitbucketClient) DoJson(method, endpoint string, payload *bytes.Buffer) (*http.Response, error) {
	return c.Do(method, endpoint, "application/json", payload)
}

func (c *BitbucketClient) Get(endpoint string) (*http.Response, error) {
	return c.DoJson("GET", endpoint, nil)
}

func (c *BitbucketClient) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.DoJson("POST", endpoint, jsonpayload)
}

func (c *BitbucketClient) Put(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.DoJson("PUT", endpoint, jsonpayload)
}

func (c *BitbucketClient) PutOnly(endpoint string) (*http.Response, error) {
	return c.DoJson("PUT", endpoint, nil)
}

func (c *BitbucketClient) Delete(endpoint string) (*http.Response, error) {
	return c.DoJson("DELETE", endpoint, nil)
}

// API version 1.0 does not use/support "application/json" on all endpoints
func (c *BitbucketClient) PostFormEncoded(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, "application/x-www-form-urlencoded", jsonpayload)
}
