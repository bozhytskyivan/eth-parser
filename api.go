package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type httpHandlers struct {
	service Service
}

type HttpHandlersOption func(handlers *httpHandlers)

func NewHttpHandlers(opts ...HttpHandlersOption) *httpHandlers {
	h := &httpHandlers{
		service: NewService(),
	}

	for _, optionFn := range opts {
		optionFn(h)
	}

	return h
}

const (
	queryParamAddress = "address"
)

func (h *httpHandlers) getCurrentBlockHandler(rw http.ResponseWriter, r *http.Request) {
	currentBlock, err := h.service.GetCurrentBlock(nil)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write([]byte(fmt.Sprintf("%d", currentBlock)))
	if err != nil {
		log.Printf("Could not write response body")
	}
}

func (h *httpHandlers) getTransactionsHandler(rw http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get(queryParamAddress)

	transactions, err := h.service.GetTransactions(r.Context(), address)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}

	resBody, err := json.Marshal(transactions)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(resBody)
	if err != nil {
		log.Printf("Could not write response body")
	}
}

func (h *httpHandlers) subscribeHandler(rw http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get(queryParamAddress)

	result, err := h.service.Subscribe(r.Context(), address)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write([]byte(fmt.Sprintf("%v", result)))
	if err != nil {
		log.Printf("Could not write response body")
	}
}

func (h *httpHandlers) unsubscribeHandler(rw http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get(queryParamAddress)

	result, err := h.service.Unsubscribe(r.Context(), address)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write([]byte(fmt.Sprintf("%v", result)))
	if err != nil {
		log.Printf("Could not write response body")
	}
}

type Service interface {
	// last parsed block
	GetCurrentBlock(ctx context.Context) (int64, error)
	// add address to observer
	Subscribe(ctx context.Context, address string) (bool, error)
	// remove address and transactions from observer
	Unsubscribe(ctx context.Context, address string) (bool, error)
	// list of inbound or outbound transactions for the address
	GetTransactions(ctx context.Context, address string) ([]Transaction, error)

	ParseBlocks(ctx context.Context, wg *sync.WaitGroup)
}
