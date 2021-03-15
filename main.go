package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	apis := NewHttpHandlers()

	http.HandleFunc("/getCurrentBlock", apis.getCurrentBlockHandler)
	http.HandleFunc("/getTransactions", apis.getTransactionsHandler)
	http.HandleFunc("/subscribe", apis.subscribeHandler)
	http.HandleFunc("/unsubscribe", apis.unsubscribeHandler)

	port := 8080
	addr := fmt.Sprintf("0.0.0.0:%d", port)

	log.Println("Started serving port", port)

	log.Fatal(http.ListenAndServe(addr, nil))
}
