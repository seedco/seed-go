# Seed-Go
A Go client for the Seed API

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/seedco/seed-go)

## Usage

```go
// to obtain an access token go to https://api.seed.co/v1/public/auth/token in a browser and enter in your seed username/password

accessToken := "1.iap2H-4qQ-WR9sy55555uaytQ.o5A32LYL5-87a_60kcQiX1Lp878GVbx8xfVvTfp5tpc.orsHbAqao-5KfsH8SdglQFltK7Ii8ktL7xo8tls3HAB"

client := seed.New(accessToken)

getTransactionsReq := TransactionRequest{
	Limit: 100,
}

txs, err := seed.GetTransactions(getTransactionsReq)
```
