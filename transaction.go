package seed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	MaxBatchSize = 1000
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
	Client            *Client
}

type TransactionsIterator struct {
	request   *TransactionsRequest
	response  *TransactionsResponse
	batchSize int
}

type TransactionsResponse struct {
	Errors  ErrorList     `json:"errors"`
	Results []Transaction `json:"results"`
	Pages   Pages         `json:"pages"`
}

func (t *TransactionsRequest) GetAll() ([]Transaction, error) {
	var resp *TransactionsResponse
	var err error
	if resp, err = t.get(nil); err != nil {
		return []Transaction{}, err
	}

	return resp.Results, nil
}

func (t *TransactionsRequest) get(paginationParams *url.Values) (*TransactionsResponse, error) {
	var err error
	var req *http.Request
	var response *TransactionsResponse

	params := &url.Values{}
	if t.CheckingAccountID != "" {
		params.Set("checking_account_id", t.CheckingAccountID)
	}
	if t.Status != "" {
		params.Set("status", t.Status)
	}
	dateLayout := "2006-01-02"
	if !t.From.IsZero() {
		params.Set("from", t.From.Format(dateLayout))
	}
	if !t.To.IsZero() {
		params.Set("to", t.To.Format(dateLayout))
	}

	var url *url.URL
	if url, err = url.Parse(fmt.Sprintf("%s/%s", ApiBase, "transactions/")); err != nil {
		return nil, err
	}

	if paginationParams != nil {
		url.RawQuery = fmt.Sprintf("%s&%s", params.Encode(), paginationParams.Encode())
	}

	if req, err = http.NewRequest("GET", url.String(), nil); err != nil {
		return nil, err
	}
	var resp *http.Response

	if resp, err = t.Client.do(req); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}

func (r *TransactionsRequest) Iterator() TransactionsIterator {
	return TransactionsIterator{
		request:   r,
		batchSize: MaxBatchSize,
	}
}

func (t *TransactionsIterator) SetBatchSize(n int) {
	if n < MaxBatchSize {
		t.batchSize = n
	}
}

// Next will retrieve the next batch of transactions. It returns a slice of Transactions, and any http errors
func (t *TransactionsIterator) Next() ([]Transaction, error) {
	var err error
	var params *url.Values
	if t.response != nil {
		params = t.response.Pages.Next
	} else {
		params = &url.Values{"limit": []string{strconv.Itoa(t.batchSize)}}
	}

	if t.response, err = t.request.get(params); err != nil {
		return []Transaction{}, fmt.Errorf("error when sending the request to seed: %v", err)
	}
	return t.response.Results, &t.response.Errors
}

// Previous will retrieve the previous batch of transactions. It returns a slice of Transactions, and any errors that happen
func (t *TransactionsIterator) Previous() ([]Transaction, error) {
	var err error
	var params *url.Values
	if t.response != nil {
		params = t.response.Pages.Previous
	} else {
		params = &url.Values{"limit": []string{strconv.Itoa(t.batchSize)}}
	}

	if t.response, err = t.request.get(params); err != nil {
		return []Transaction{}, fmt.Errorf("error when sending the request to seed: %v", err)
	}
	return t.response.Results, &t.response.Errors
}
