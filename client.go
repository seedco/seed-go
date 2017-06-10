package seed

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

const (
	ApiBase = "https://api.seed.co/v1/public"
)

type Client struct {
	accessToken   string
	clientVersion string
	httpClient    *http.Client
}

// New creates a new Seed client given an access token
func New(accessToken string) *Client {
	return &Client{
		httpClient:  &http.Client{},
		accessToken: accessToken,
	}
}

// SetClientVersion sets the client version that each request will be made with
func (c *Client) SetClientVersion(clientVersion string) {
	c.clientVersion = clientVersion
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	if c.clientVersion != "" {
		req.Header.Set("Client-Version-Id", c.clientVersion)
	}
	return c.httpClient.Do(req)
}

type Pages struct {
	Next     *url.Values `json:"next"`
	Previous *url.Values `json:"previous"`
}

type ErrorList []map[string]string

func (e ErrorList) Error() string {
	buf := bytes.NewBufferString("")
	for _, e := range e {
		errorString := fmt.Sprintf("field: %s, message: %s", e["field"], e["message"])
		buf.WriteString(errorString)
		buf.WriteString("\n")
	}
	return buf.String()
}
