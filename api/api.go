package api

import (
	"net/http"
)

type PayersSpent struct {
	Payer  string `json:"payer,omitempty"`
	Points int    `json:"points,omitempty"`
}

type Transaction struct {
	Payer     *string `json:"payer,omitempty"`
	Points    *int    `json:"points,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
}

type UserSpent struct {
	Points *int `json:"points,omitempty"`
}

type AddTransactionJSONRequestBody = Transaction

type SpendPointsJSONRequestBody = UserSpent

type AddTransactionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

type ClearDBResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

type GetBalanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *map[string]int
}

type SpendPointsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]PayersSpent
}
