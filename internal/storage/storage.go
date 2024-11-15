package storage

import "github.com/vdhieu/tx-parser/internal/models"

type Storage interface {
	AddSubscriber(address string) error
	IsSubscribed(address string) bool
	GetSubscribers() []string

	SaveTransactions(address string, txs []models.Transaction) error
	GetTransactions(address string) ([]models.Transaction, error)

	SetCurrentBlock(blockNum int64) error
	GetCurrentBlock() (int64, error)
}
