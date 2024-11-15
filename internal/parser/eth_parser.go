package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vdhieu/tx-parser/internal/models"
	"github.com/vdhieu/tx-parser/internal/storage"
	"github.com/vdhieu/tx-parser/pkg/logger"
	"github.com/vdhieu/tx-parser/pkg/notification"
	"github.com/vdhieu/tx-parser/pkg/rpc"
	"go.uber.org/zap"
)

type ethParser struct {
	storage  storage.Storage
	client   rpc.Client
	log      *zap.Logger
	notifier notification.Notifier
	running  bool
}

// NewEthParser create new parser instance and start a background process to process eth blocks
func NewEthParser(storage storage.Storage, client rpc.Client, notifier notification.Notifier) Parser {
	log := logger.GetLogger()
	p := &ethParser{
		storage:  storage,
		client:   client,
		notifier: notifier,
		log:      log.With(zap.String("parser", "eth")),
	}

	p.log.Info("Starting ETH parser background process")
	p.running = true
	go p.processBlocks()

	return p
}

// Shutdown stop the parser
func (p *ethParser) Shutdown() {
	p.log.Info("Shutting down ETH parser")
	p.running = false
}

// GetCurrentBlock return current processed block
func (p *ethParser) GetCurrentBlock() int {
	block, _ := p.storage.GetCurrentBlock()
	return int(block)
}

// Subscribe subscribe address to event notification
func (p *ethParser) Subscribe(address string) bool {
	err := p.storage.AddSubscriber(strings.ToLower(address))
	if err != nil {
		p.log.Error("Failed to add subscriber",
			zap.String("address", address),
			zap.Error(err))
		return false
	}
	p.log.Info("New subscriber added", zap.String("address", address))

	return true
}

// GetTransactions return all txn parsed filter by address
func (p *ethParser) GetTransactions(address string) []models.Transaction {
	txs, err := p.storage.GetTransactions(strings.ToLower(address))
	if err != nil {
		p.log.Error("Failed to add subscriber",
			zap.String("address", address),
			zap.Error(err))
		return nil
	}
	return txs
}

func (p *ethParser) processBlocks() {
	log := p.log
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for p.running {
		<-ticker.C
		currentBlock, err := p.client.GetLatestBlockNumber()
		if err != nil {
			p.log.Error("Failed to get latest block number", zap.Error(err))
			continue
		}

		lastProcessed, _ := p.storage.GetCurrentBlock()
		if int64(lastProcessed) >= currentBlock {
			continue
		}
		// make sure when startup process from newest block only
		if lastProcessed == 0 {
			lastProcessed = currentBlock - 1
		}

		log.Info("Processing new blocks",
			zap.Int64("from_block", lastProcessed),
			zap.Int64("to_block", currentBlock))

		for blockNum := lastProcessed + 1; blockNum <= currentBlock; blockNum++ {
			block, err := p.client.GetBlockByNumber(blockNum)
			if err != nil {
				log.Error("Failed to get block",
					zap.Int64("block_number", blockNum),
					zap.Error(err))
				continue
			}

			p.storage.SetCurrentBlock(blockNum)
			p.processTransactions(block)
			log.Debug("Processed block", zap.Int64("block_number", blockNum))
		}
	}
}

func (p *ethParser) processTransactions(block rpc.Block) {
	transactions := block.Transactions
	if transactions == nil {
		p.log.Error("Invalid transactions data in block",
			zap.Int64("block_number", block.Number))
		return
	}

	subscribers := p.storage.GetSubscribers()
	if len(subscribers) == 0 {
		return
	}

	subscriberMap := make(map[string]bool)
	for _, addr := range subscribers {
		subscriberMap[strings.ToLower(addr)] = true
	}

	blockNumber := block.Number
	matchedTxs := 0
	for _, tx := range transactions {
		fromAddr := strings.ToLower(tx.From)
		toAddr := strings.ToLower(tx.To)

		if subscriberMap[fromAddr] || subscriberMap[toAddr] {
			matchedTxs++
			transaction := models.Transaction{
				Hash:        tx.Hash,
				From:        fromAddr,
				To:          tx.To,
				Value:       strconv.FormatInt(tx.Value, 10),
				BlockNumber: strconv.FormatInt(blockNumber, 10),
				Timestamp:   strconv.FormatInt(block.Timestamp, 10),
			}

			p.log.Debug("Found matching transaction",
				zap.String("hash", transaction.Hash),
				zap.String("from", tx.From),
				zap.String("to", tx.To))

			if subscriberMap[fromAddr] {
				existing, _ := p.storage.GetTransactions(fromAddr)
				err := p.storage.SaveTransactions(fromAddr, append(existing, transaction))
				if err != nil {
					p.log.Error(fmt.Sprintf("Unable to save txn for address %v", fromAddr), zap.Error(err))
				}

				p.notifier.Notify(fromAddr, "found a new transactions", transaction)
			}

			if subscriberMap[toAddr] {
				existing, _ := p.storage.GetTransactions(toAddr)
				err := p.storage.SaveTransactions(toAddr, append(existing, transaction))
				if err != nil {
					p.log.Error(fmt.Sprintf("Unable to save txn for address %v", toAddr), zap.Error(err))
				}
				p.notifier.Notify(toAddr, "found a new transactions", transaction)
			}
		}
	}

	if matchedTxs > 0 {
		p.log.Info("Processed transactions for block",
			zap.Int64("block_number", blockNumber),
			zap.Int("matched_transactions", matchedTxs))
	}
}
