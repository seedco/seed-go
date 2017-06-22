package seed

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Balance contains relevant balance amounts for a given checking account
type Balance struct {
	// Checking Account ID specifies the id of the checking account that this balance belongs to
	CheckingAccountID string `json:"checking_account_id"`
	// Total Available refers to the balance that is safely usable
	// this number is calculated in the following way TotalAvailable = Accessible - PendingDebits - ScheduledDebits
	TotalAvailable int64 `json:"total_available"`
	// Settled refers to the total amount of transactions that have settled
	Settled int64 `json:"settled"`
	// PendingCredits refers to credits that are pending
	PendingCredits uint64 `json:"pending_credits"`
	// PendingDebits refers to debits that are pending
	PendingDebits uint64 `json:"pending_debits"`
	// ScheduledDebits refers to debits that are scheduled
	ScheduledDebits uint64 `json:"scheduled_debits"`
	// Accessible refers to the balance is usable
	Accessible int64 `json:"accessible"`
	// Lockbox refers to the amount in the virtual lockbox
	Lockbox uint64 `json:"lockbox"`
}

// BalanceRequest is a request for fetching a balance for a given checking account
type BalanceRequest struct {
	// Checking Account ID specifies the id of the checking account for the balance in question
	// Must be a uuid that corresponds to a valid checking account
	CheckingAccountID string
	// Client is the seed client
	Client *Client
}

// BalanceResponse is the struct that the server response will get unmarshalled into
type BalanceResponse struct {
	// Errors is a list of errors
	Errors ErrorList `json:"errors"`
	// Results is a slice of Balance objects. The size of this slice is expected to be 1
	Results []Balance `json:"results"`
}

// NewBalanceRequest creates a new balance request
func (c *Client) NewBalanceRequest() *BalanceRequest {
	return &BalanceRequest{Client: c}
}

// Get retrieves the balance
func (b *BalanceRequest) Get() (Balance, error) {
	var response BalanceResponse
	var req *http.Request
	var err error

	if req, err = http.NewRequest("GET", fmt.Sprintf("%s/%s", ApiBase, "balance"), nil); err != nil {
		return Balance{}, err
	}
	var resp *http.Response

	if resp, err = b.Client.do(req); err != nil {
		return Balance{}, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Balance{}, err
	}

	balances := response.Results

	if len(balances) == 0 {
		return Balance{}, errors.New("no balance found")
	}

	return balances[0], nil
}
