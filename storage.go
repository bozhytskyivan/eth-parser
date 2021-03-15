package main

import (
	"context"
	"sync"
)

type storage struct {
	transactions []Transaction
	txMux        sync.RWMutex

	subscriptions    map[string]Subscription
	subscriptionsMux sync.RWMutex
}

func NewStorage() *storage {
	return &storage{
		subscriptions: make(map[string]Subscription),
	}
}

type Subscription struct {
	Address string
}

func (s *storage) AddTransaction(_ context.Context, t Transaction) error {
	s.txMux.Lock()
	defer s.txMux.Unlock()

	s.transactions = append(s.transactions, t)
	return nil
}

func (s *storage) GetCurrentBlock(_ context.Context) (int64, error) {
	s.txMux.RLock()
	defer s.txMux.RUnlock()

	size := len(s.transactions)
	if size == 0 {
		return -1, nil
	}

	return int64(s.transactions[size-1].BlockNumber), nil
}

func (s *storage) GetTransactionsBy(_ context.Context, address string) ([]Transaction, error) {
	s.txMux.RLock()
	defer s.txMux.RUnlock()

	var result []Transaction
	for _, t := range s.transactions {
		if t.From == address || t.To == address {
			result = append(result, t)
		}
	}

	return result, nil
}

func (s *storage) GetTransactionByHash(_ context.Context, hash string) (Transaction, error) {
	s.txMux.RLock()
	defer s.txMux.RUnlock()

	var result Transaction
	for _, t := range s.transactions {
		if t.Hash == hash {
			result = t
			break
		}
	}

	return result, nil
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
