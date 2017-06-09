package seed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Transaction struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount"`
	Error       string    `json:"error"`
	Status      string    `json:"status"`
	Category    string    `json:"category"`
}

type TransactionsRequest struct {
	CheckingAccountID string
	Status            string
	From              time.Time
	To                time.Time
	Offset            int
	Limit             int
}

type TransactionsResponse struct {
	Errors  []map[string]string `json:"errors"`
	Results []Transaction       `json:"results"`
	Pages   map[string]string   `json:"pages"`
}

func (c *Client) GetTransactions(transactionsReq TransactionsRequest) ([]Transaction, error) {
	var err error
	var req *http.Request
	var transactions []Transaction

	if req, err = http.NewRequest("GET", fmt.Sprintf("%s/%s", ApiBase, "transactions/"), nil); err != nil {
		return transactions, err
	}
	var resp *http.Response

	if resp, err = c.Do(req); err != nil {
		return transactions, err
	}

	var response TransactionsResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	if c.clientVersion != "" {
		req.Header.Set("Client-Version-Id", c.clientVersion)
	}
	return c.httpClient.Do(req)
}
