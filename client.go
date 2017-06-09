package seed

import (
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

func New(accessToken string) *Client {
	return &Client{
		httpClient:  &http.Client{},
		accessToken: accessToken,
	}
}

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
