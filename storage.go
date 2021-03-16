package main

import (
	"context"
	"sync"
)

type storage struct {
	userTransactions map[string][]Transaction
	transactions     map[string]Transaction
	latestBlock      int64
	txMux            sync.RWMutex

	subscriptions    map[string]Subscription
	subscriptionsMux sync.RWMutex
}

func NewStorage() *storage {
	return &storage{
		userTransactions: make(map[string][]Transaction),
		transactions:     make(map[string]Transaction),
		subscriptions:    make(map[string]Subscription),
	}
}

type Subscription struct {
	Address string
}

func (s *storage) AddTransaction(_ context.Context, t Transaction) error {
	s.txMux.Lock()
	defer s.txMux.Unlock()

	s.userTransactions[t.From] = append(s.userTransactions[t.From], t)
	s.userTransactions[t.To] = append(s.userTransactions[t.To], t)
	s.transactions[t.Hash] = t
	return nil
}

func (s *storage) GetCurrentBlock(_ context.Context) (int64, error) {
	s.txMux.RLock()
	defer s.txMux.RUnlock()

	return s.latestBlock, nil
}

func (s *storage) SetCurrentBlock(_ context.Context, blockNumber int64) error {
	s.txMux.RLock()
	defer s.txMux.RUnlock()

	s.latestBlock = blockNumber

	return nil
}

func (s *storage) GetTransactionsBy(_ context.Context, address string) ([]Transaction, error) {
	s.txMux.RLock()
	defer s.txMux.RUnlock()

	transactions := s.userTransactions[address]

	return transactions, nil
}

func (s *storage) DeleteTransactionsBy(_ context.Context, address string) error {
	transactions := s.userTransactions[address]
	for _, t := range transactions {
		delete(s.transactions, t.Hash)
	}

	delete(s.userTransactions, address)

	return nil
}

func (s *storage) GetTransactionByHash(_ context.Context, hash string) (Transaction, error) {
	s.txMux.RLock()
	defer s.txMux.RUnlock()

	transaction := s.transactions[hash]

	return transaction, nil
}

func (s *storage) AddSubscription(_ context.Context, address string) error {
	s.subscriptionsMux.Lock()
	defer s.subscriptionsMux.Unlock()

	if _, ok := s.subscriptions[address]; !ok {
		s.subscriptions[address] = Subscription{Address: address}
	}

	return nil
}

func (s *storage) RemoveSubscription(_ context.Context, address string) error {
	s.subscriptionsMux.Lock()
	defer s.subscriptionsMux.Unlock()

	delete(s.subscriptions, address)
	return nil
}

func (s *storage) GetSubscription(_ context.Context, address string) (Subscription, error) {
	s.subscriptionsMux.RLock()
	defer s.subscriptionsMux.RUnlock()

	var result Subscription
	for _, subscr := range s.subscriptions {
		if subscr.Address == address {
			result = subscr
			break
		}
	}

	return result, nil
}
