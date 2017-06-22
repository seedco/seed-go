package seed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// MaxBatchSize is the maximum pagination limit for a transaction query
	MaxBatchSize = 1000
)

// Transaction contains relevant information about a transaction
type Transaction struct {
	// Date is the date of the transaction
	Date time.Time `json:"date"`
	// Description is the description of the transaction
	Description string `json:"description"`
	// Amount is the amount of the transaction in cents
	Amount int64 `json:"amount"`
	// Error contains any errors that happened with the transaction
	Error string `json:"error"`
	// Status is either "pending" or "settled"
	Status string `json:"status"`
	// Category is the category of the transaction
	Category string `json:"category"`
}

// TransactionsRequest contains fields for querying transactions
type TransactionsRequest struct {
	// CheckingAccountID is a uuid of the checking account for the transaction request
	CheckingAccountID string
	// Status is either "pending" or "settled"
	Status string
	// From is the start date of the date range, inclusive
	From time.Time
	// To is the end date of the date range, exclusive
	To time.Time
	// Client is the seed client that will send the request
	Client *Client
}

// TransactionsIterator is an iterator to iterate through pages of transaction results
type TransactionsIterator struct {
	request   *TransactionsRequest
	response  *TransactionsResponse
	batchSize int
}

// TransactionsResponse is the response object that the server data unmarshalls into
type TransactionsResponse struct {
	// Errors is a list of errors
	Errors ErrorList `json:"errors"`
	// Results is a slice of transaction objects
	Results []Transaction `json:"results"`
	// Pages contains pagination information
	Pages Pages `json:"pages"`
}

// Get retrieves a list of transactions
func (t *TransactionsRequest) Get() ([]Transaction, error) {
	var resp *TransactionsResponse
	var err error
	if resp, err = t.get(nil); err != nil {
		return []Transaction{}, err
	}

	if len(resp.Errors) > 0 {
		return resp.Results, resp.Errors
	}

	return resp.Results, nil
}

func (t *TransactionsRequest) get(paginationParams *PaginationParams) (*TransactionsResponse, error) {
	var err error
	var req *http.Request
	var response TransactionsResponse

	params := &url.Values{}
	if t.CheckingAccountID != "" {
		params.Set("checking_account_id", t.CheckingAccountID)
	}
	if t.Status != "" {
		params.Set("status", t.Status)
	}
	dateLayout := "2006-01-02"
	if !t.From.IsZero() {
		params.Set("start_date", t.From.Format(dateLayout))
	}
	if !t.To.IsZero() {
		params.Set("end_date", t.To.Format(dateLayout))
	}

	var url *url.URL
	if url, err = url.Parse(fmt.Sprintf("%s/%s", ApiBase, "transactions")); err != nil {
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

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Iterator returns a TransactionIterator
func (t *TransactionsRequest) Iterator() TransactionsIterator {
	return TransactionsIterator{
		request:   t,
		batchSize: MaxBatchSize,
	}
}

// SetBatchSize sets the batch size for paginated results
func (t *TransactionsIterator) SetBatchSize(n int) {
	if n < MaxBatchSize {
		t.batchSize = n
	}
}

// Next will retrieve the next batch of transactions. It returns a slice of Transactions, and any http errors
func (t *TransactionsIterator) Next() ([]Transaction, error) {
	var err error
	var params PaginationParams
	if t.response != nil {
		params = t.response.Pages.Next
	} else {
		params = PaginationParams{Limit: t.batchSize}
	}

	if t.response, err = t.request.get(&params); err != nil {
		return []Transaction{}, fmt.Errorf("error when sending the request to seed: %v", err)
	}

	if len(t.response.Errors) > 0 {
		return t.response.Results, t.response.Errors
	}

	return t.response.Results, nil
}

// Previous will retrieve the previous batch of transactions. It returns a slice of Transactions, and any errors that happen
func (t *TransactionsIterator) Previous() ([]Transaction, error) {
	var err error
	var params PaginationParams
	if t.response != nil {
		params = t.response.Pages.Previous
	} else {
		params = PaginationParams{Limit: t.batchSize}
	}

	if t.response, err = t.request.get(&params); err != nil {
		return []Transaction{}, fmt.Errorf("error when sending the request to seed: %v", err)
	}
	if len(t.response.Errors) > 0 {
		return t.response.Results, t.response.Errors
	}

	return t.response.Results, nil
}
