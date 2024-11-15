package parser

import (
	"github.com/vdhieu/tx-parser/internal/models"
)

type Parser interface {
	// Shutdown stop the parser
	Shutdown()
	// GetCurrentBlock last parsed block
	GetCurrentBlock() int
	// Subscribe add address to observer
	Subscribe(address string) bool
	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) []models.Transaction
}
