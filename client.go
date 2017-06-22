package seed

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// ApiBase is the base url for the seed api
	ApiBase = "https://api.seed.co/v1/public"
)

// Client is a seed client that can be used to create different request objects to fetch data
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

// Pages contains pagination information
type Pages struct {
	// Next is the set of parameters for the next page
	Next PaginationParams `json:"next"`
	// Previous is the set of parameters for the previous page
	Previous PaginationParams `json:"previous"`
}

// PaginationParams encapsulates the two pagination values, offset and limit
type PaginationParams struct {
	// Offset is the pagination offset index
	Offset int `json:"offest"`
	// Limit is the pagination limit
	Limit int `json:"limit"`
}

// MarshalJSON marshalls pagination params
func (p PaginationParams) MarshalJSON() ([]byte, error) {
	return []byte(p.Encode()), nil
}

// UnmarshalJSON unmarshalls pagination params
func (p *PaginationParams) UnmarshalJSON(d []byte) error {
	var err error
	s := string(d)
	params := strings.Split(s, "&")

	for _, param := range params {
		split := strings.Split(param, "=")
		if len(split) != 2 {
			continue
		}
		switch split[0] {
		case "limit":
			var limit int
			if limit, err = strconv.Atoi(split[1]); err != nil {
				return err
			}
			p.Limit = limit
		case "offset":
			var offset int
			if offset, err = strconv.Atoi(split[1]); err != nil {
				return err
			}
			p.Offset = offset
		}
	}
	return nil
}

// Encode encodes offset and limit into a url query string
func (p PaginationParams) Encode() string {
	u := url.Values{}
	if p.Offset > 0 {
		u.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit > 0 {
		u.Set("limit", strconv.Itoa(p.Limit))
	}

	return u.Encode()
}

// ErrorList is a list of maps that contain error message, field pairs
type ErrorList []map[string]string

// Error returns a comma delineated string of errors
func (e ErrorList) Error() string {
	buf := bytes.NewBufferString("")
	for _, e := range e {
		errorString := fmt.Sprintf("field: %s, message: %s", e["field"], e["message"])
		buf.WriteString(errorString)
		buf.WriteString("\n")
	}
	return buf.String()
}
