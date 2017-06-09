package seed

import "net/http"

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
