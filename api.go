package main

import (
	"context"
	"net/http"
)

type api struct {
	service Service
}

func NewAPI() *api {
	return &api{
		service: NewService(),
	}
}

func (a *api) getCurrentBlockHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("getCurrentBlockHandler"))
}

func (a *api) getTransactionsHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("getTransactionsHandler"))
}

func (a *api) subscribeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("subscribeHandler"))
}

func (a *api) unsubscribeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("unsubscribeHandler"))
}

type Service interface {
	// last parsed block
	GetCurrentBlock() int64
	// add address to observer
	Subscribe(address string) bool
	// remove address and transactions from observer
	Unsubscribe(address string) bool
	// list of inbound or outbound transactions for the address
	GetTransactions(address string) []Transaction

	ParseBlocks(ctx context.Context)
}
