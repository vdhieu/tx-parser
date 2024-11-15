package storage

import (
	"sync"

	"github.com/vdhieu/tx-parser/internal/models"
)

type memoryStorage struct {
	currentBlock int64
	subscribers  map[string]bool
	transactions map[string][]models.Transaction
	mu           sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		subscribers:  make(map[string]bool),
		transactions: make(map[string][]models.Transaction),
	}
}

func (s *memoryStorage) AddSubscriber(address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers[address] = true
	return nil
}

func (s *memoryStorage) IsSubscribed(address string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.subscribers[address]
}

func (s *memoryStorage) GetSubscribers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addresses := make([]string, 0, len(s.subscribers))
	for addr := range s.subscribers {
		addresses = append(addresses, addr)
	}
	return addresses
}

// SaveTransactions save txn for subscribed address, non-subscribber will not be saved
func (s *memoryStorage) SaveTransactions(address string, txs []models.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subscribers[address] {
		s.transactions[address] = txs
	}
	return nil
}

func (s *memoryStorage) GetTransactions(address string) ([]models.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transactions[address], nil
}

func (s *memoryStorage) SetCurrentBlock(blockNum int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentBlock = blockNum
	return nil
}

func (s *memoryStorage) GetCurrentBlock() (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentBlock, nil
}
