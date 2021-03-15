package main

import (
	"context"
	"log"
)

type service struct {
	storage   Storage
	ethClient EthClient

	id int64
}

type serviceOption func(s *service)

const DefaultID = 83

func NewService(opts ...serviceOption) *service {
	s := &service{
		storage:   NewStorage(),
		ethClient: NewEthClient(),
		id:        DefaultID,
	}

	for _, optFn := range opts {
		optFn(s)
	}

	return s
}

type Event struct {
	Type        string
	Address     string
	Transaction Transaction
}

const (
	EventTypeTransactionSent     = "TRANSACTION_SENT"
	EventTypeTransactionReceived = "TRANSACTION_RECEIVED"
)

// last parsed block
func (s *service) GetCurrentBlock() int64 {
	blockNumber, err := s.ethClient.GetCurrentBlock(s.id)
	if err != nil {
		log.Fatal(err)
	}

	return blockNumber
}

// add address to observer
func (s *service) Subscribe(address string) bool {
	err := s.storage.AddSubscription(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func (s *service) Unsubscribe(address string) bool {
	err := s.storage.RemoveSubscription(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	//todo: remove unnecessary transactions

	return true
}

// list of inbound or outbound transactions for an address
func (s *service) GetTransactions(address string) []Transaction {
	transactions, err := s.storage.GetTransactionsBy(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	return transactions
}

func (s *service) ParseBlocks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// todo: use ticker
			ctx := context.Background()
			transactions, err := s.getNewTransactions(ctx)
			if err != nil {
				log.Printf("Error while fetching transactions: %v", err)
				continue
			}

			for _, transaction := range transactions {
				err = s.processNewTransaction(ctx, transaction)
				if err != nil {
					log.Printf("Error while processing new transaction: %v", err)
				}
			}
		}
	}
}

func (s *service) onEvent(event Event) error {
	log.Printf("New %s event in block %d for subsriber %s", event.Type, event.Transaction.BlockNumber, event.Address)

	return nil
}

func (s *service) processNewTransaction(ctx context.Context, tx Transaction) error {
	senderSubscription, err := s.storage.GetSubscription(ctx, tx.From)
	if err != nil {
		return err
	}

	if len(senderSubscription.Address) > 0 {
		err = s.onEvent(Event{
			Type:        EventTypeTransactionSent,
			Address:     tx.From,
			Transaction: tx,
		})
		if err != nil {
			return err
		}
	}

	receiverSubscription, err := s.storage.GetSubscription(ctx, tx.To)
	if err != nil {
		return err
	}

	if len(receiverSubscription.Address) > 0 {
		err = s.onEvent(Event{
			Type:        EventTypeTransactionReceived,
			Address:     tx.To,
			Transaction: tx,
		})
		if err != nil {
			return err
		}
	}

	if len(receiverSubscription.Address) > 0 || len(senderSubscription.Address) > 0 {
		err := s.storage.AddTransaction(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) getNewTransactions(ctx context.Context) ([]Transaction, error) {
	lastBlockNumber, err := s.storage.GetCurrentBlock(ctx)
	if err != nil {
		return nil, err
	}

	currentBlockNumber, err := s.ethClient.GetCurrentBlock(s.id)
	if err != nil {
		return nil, err
	}

	var result []Transaction

	lastBlock, err := s.ethClient.GetBlockByNumber(s.id, lastBlockNumber)
	if err != nil {
		return nil, err
	}

	// get unprocessed transactions from last block
	for _, tx := range lastBlock.Transactions {
		processedTx, err := s.storage.GetTransactionByHash(ctx, tx.Hash)
		if err != nil {
			return nil, err
		}
		if len(processedTx.Hash) == 0 {
			result = append(result, processedTx)
		}
	}

	lastBlockNumber++

	// get all transactions from new blocks
	for lastBlockNumber <= currentBlockNumber {
		block, err := s.ethClient.GetBlockByNumber(s.id, lastBlockNumber)
		if err != nil {
			return nil, err
		}

		result = append(result, block.Transactions...)

		lastBlockNumber++
	}

	return result, nil
}

type Storage interface {
	// Transactions processing
	AddTransaction(ctx context.Context, t Transaction) error
	GetCurrentBlock(ctx context.Context) (int64, error)
	GetTransactionsBy(ctx context.Context, address string) ([]Transaction, error)
	GetTransactionByHash(ctx context.Context, hash string) (Transaction, error)

	// Subscriptions processing
	AddSubscription(ctx context.Context, address string) error
	RemoveSubscription(ctx context.Context, address string) error
	GetSubscription(ctx context.Context, address string) (Subscription, error)
}

type EthClient interface {
	GetCurrentBlock(id int64) (int64, error)
	GetBlockByNumber(id int64, number int64) (Block, error)
}
