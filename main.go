package main

import (
	"Fetch-Coding-Exercise2024/db"
	"Fetch-Coding-Exercise2024/endpoints"
	"net/http"
)

func main() {
	db.InitDB()
	endpoints.UpdateCache()
	//POST method /add
	http.HandleFunc("/add", endpoints.AddTransaction)
	//POST method /spend
	http.HandleFunc("/spend", endpoints.SpendPoints)
	//GET method /balance
	http.HandleFunc("/balance", endpoints.GetBalance)
	//DELETE method /clear
	http.HandleFunc("/clear", endpoints.ClearDB)

	http.ListenAndServe(":8000", nil)
}
