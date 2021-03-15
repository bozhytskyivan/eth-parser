package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt)

	s := NewService()

	apis := NewHttpHandlers(func(handlers *httpHandlers) {
		handlers.service = s
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/getCurrentBlock", apis.getCurrentBlockHandler)
	mux.HandleFunc("/getTransactions", apis.getTransactionsHandler)
	mux.HandleFunc("/subscribe", apis.subscribeHandler)
	mux.HandleFunc("/unsubscribe", apis.unsubscribeHandler)

	wg := &sync.WaitGroup{}

	// start api
	apiCtx, apiServerCancelFn := context.WithCancel(context.Background())

	go serveApi(apiCtx, mux, wg)

	wg.Add(1)

	// start parsing blocks
	parserCtx, parserCancelFn := context.WithCancel(context.Background())

	go s.ParseBlocks(parserCtx, wg)

	wg.Add(1)

	// wait until process terminated
	sig := <-osSignalChan
	log.Println("System signal received:", sig)

	// stop serving api
	apiServerCancelFn()

	// stop parsing
	parserCancelFn()

	// wait until workers complete
	wg.Wait()

	log.Println("Gracefully exited")
}

func serveApi(ctx context.Context, handler http.Handler, wg *sync.WaitGroup) {
	port := 8080

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	go func() {
		log.Println("Started serving port", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelFn := context.WithTimeout(ctx, time.Second*5)
	defer cancelFn()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Println(err)
	}

	wg.Done()
}
