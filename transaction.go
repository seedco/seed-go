package seed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	Client            *Client
}

type TransactionsIterator struct {
	request  *TransactionsRequest
	response TransactionsResponse
	hasRun   bool
}

type TransactionsResponse struct {
	Errors  []map[string]string `json:"errors"`
	Results []Transaction       `json:"results"`
	Pages   Pages               `json:"pages"`
}

func (t *TransactionsRequest) get(params *url.Values) (TransactionsResponse, error) {
	var err error
	var req *http.Request
	var response TransactionsResponse

	var url *url.URL
	if url, err = url.Parse(fmt.Sprintf("%s/%s", ApiBase, "transactions/")); err != nil {
		return response, err
	}

	if params != nil {
		url.RawQuery = params.Encode()
	}

	if req, err = http.NewRequest("GET", url.String(), nil); err != nil {
		return response, err
	}
	var resp *http.Response

	if resp, err = t.Client.do(req); err != nil {
		return response, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}
	return response, nil
}

func (r *TransactionsRequest) Iterator() TransactionsIterator {
	return TransactionsIterator{
		request: r,
	}
}

func (t *TransactionsIterator) Next() error {
	var err error
	if t.response, err = t.request.get(t.response.Pages.Next); err != nil {
		return err
	}
	t.hasRun = true
	return nil
}

func (t *TransactionsIterator) HasNext() bool {
	return !t.hasRun || t.response.Pages.Next != nil
}

func (t *TransactionsIterator) HasPrevious() bool {
	return t.response.Pages.Previous != nil
}

func (t *TransactionsIterator) Previous() error {
	var err error
	if t.response, err = t.request.get(t.response.Pages.Previous); err != nil {
		return err
	}
	t.hasRun = true
	return nil
}

func (t *TransactionsIterator) Transactions() []Transaction {
	return t.response.Results
}

func (t *TransactionsIterator) Errors() []map[string]string {
	return t.response.Errors
}
