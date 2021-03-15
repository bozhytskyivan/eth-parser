package main

import (
	"log"
	"net/http"
)

func main() {
	apis := NewAPI()

	http.HandleFunc("/getCurrentBlock", apis.getCurrentBlockHandler)
	http.HandleFunc("/getTransactions", apis.getTransactionsHandler)
	http.HandleFunc("/subscribe", apis.subscribeHandler)
	http.HandleFunc("/unsubscribe", apis.unsubscribeHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
