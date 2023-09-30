package main

import (
	"fmt"
	"net/http"

	"github.com/ahmetahunbay/Fetch-Coding-Exercise-2024/api"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var request api.AddTransactionJSONRequestBody


}

func main() {
	http.HandleFunc("/add", handler)
	http.ListenAndServe(":8000", nil)
}
