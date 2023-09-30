package main

import (
	"Fetch-Coding-Exercise2024/api"
	"Fetch-Coding-Exercise2024/structs"
	"encoding/json"
	"net/http"
	"time"
)

var timeLayout = "2006-01-02T15:04:05Z"
var pq = structs.TransactionPQ{}
var payers = map[string]int{}
var debt = 0

func addTransaction(w http.ResponseWriter, r *http.Request) {
	var request api.AddTransactionJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {

		return
	}
	timeStamp, err := time.Parse(timeLayout, request.Timestamp)
	if err != nil {

		return
	}

	transaction := structs.Transaction{
		Payer:     request.Payer,
		Points:    request.Points,
		Timestamp: timeStamp,
	}

	_, exists := payers[transaction.Payer]
	if !exists {
		payers[transaction.Payer] = 0
	}

	if transaction.Points <= 0 {
		debt -= transaction.Points
	} else {
		payers[transaction.Payer] += transaction.Points
		pq.Push(&transaction)
	}

}

func main() {

	http.HandleFunc("/add", addTransaction)
	http.ListenAndServe(":8000", nil)
}
