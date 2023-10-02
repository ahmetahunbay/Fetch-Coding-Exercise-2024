package endpoints

import (
	"Fetch-Coding-Exercise2024/api"
	"Fetch-Coding-Exercise2024/db"
	"Fetch-Coding-Exercise2024/structs"
	"container/heap"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// priority queue to store the transactions
var pq = structs.TransactionPQ{}

// map to store the balances
var payers = map[string]int{}

// hardcoded time layout
var timeLayout = "2006-01-02T15:04:05Z"

/*
This is the /add route

It updates the cache and database with the transaction.

I add all positive transactions to the DB/cache before taking care of potential debt,
just for readability and simplicity.
*/
func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var request api.AddTransactionJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Payer == nil || request.Points == nil || request.Timestamp == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("/add json structure: {\"payer\": \"DANNON\", \"points\": 300, \"timestamp\": \"2022-10-31T10:00:00Z\"}"))
		return
	}

	//formats timestamp
	castedTimeStamp, err := time.Parse(timeLayout, *request.Timestamp)
	if err != nil || request.Payer == nil || request.Points == nil || request.Timestamp == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("timestamp format: \"2022-10-31T10:00:00Z\""))
		return
	}

	//initializes payers[payer] if not already in map
	_, exist := payers[*request.Payer]
	if !exist {
		(payers)[*request.Payer] = 0
	}

	//pulls balance BEFORE adding next transaction
	balance, err := db.GetBalance()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB Error"))
		return
	}

	//inserts transaction into db (if not a debt transaction) and gets id for later hashing
	id, err := db.InsertTransaction(*request.Payer, *request.Points, *request.Timestamp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB Error"))
		return
	}

	//adds to PQ cache if this isn't a debt transaction
	if *request.Points > 0 {
		heap.Push(&pq, structs.Transaction{
			ID:        id,
			Payer:     *request.Payer,
			Points:    *request.Points,
			Timestamp: castedTimeStamp,
		})
		(payers)[*request.Payer] += *request.Points
	}

	//***TRICKY LOGIC*** just remember that the pq has current positive transactions
	//if balance is negative, checks if new transaction can counteract it
	if balance < 0 {
		//pass in negative balance
		err = HandleDebtBalance(balance)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("DB Error"))
			return
		}
		// if balance is positive and new transaction is negative, pops from pq to try and remove debt
	} else if *request.Points < 0 {
		//pass in negative points
		err = HandleDebtBalance(*request.Points)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("DB Error"))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

/*
This is the /spend route.

It updates the db and cache after points are spent.
*/
func SpendPoints(w http.ResponseWriter, r *http.Request) {
	var request api.SpendPointsJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Points == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("/spend json structure: {\"points\": 300}"))
		return
	}

	pointsSpent := *request.Points

	balance, err := db.GetBalance()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB Error"))
		return
	}

	//checks for invalid balance
	if pointsSpent > balance {
		errorStmt := fmt.Sprintf("INVALID POINTS BALANCE: Have %d, spending %d. %d more points needed", balance, pointsSpent, pointsSpent-balance)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorStmt))
		return
	}

	//updated db balance
	db.UpdateBalance(-pointsSpent)

	//updates cache/db for pq and payers
	paidMap, err := HandlePointsSpent(pointsSpent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB Error"))
		return
	}

	//converts map of payer/point pairs into json list
	var paidList = []api.PayersSpent{}

	for payer, points := range paidMap {
		paidList = append(paidList, api.PayersSpent{
			Payer:  payer,
			Points: points,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paidList)
}

/*
This is the /balance route.

Very simple because of payer cache
*/
func GetBalance(w http.ResponseWriter, r *http.Request) {

	balance, err := db.GetBalance()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB Error"))
		return
	}

	if balance < 0 {
		errorStmt := fmt.Sprintf("INVALID POINTS BALANCE: Account is %d in debt", balance)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorStmt))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payers)
}

/*
This is an extra route: /clear.

Use this to clear the database and cache
*/
func ClearDB(w http.ResponseWriter, r *http.Request) {
	err := db.ClearDB()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DB Error"))
		return
	}

	//clears cache
	pq = structs.TransactionPQ{}
	payers = map[string]int{}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
