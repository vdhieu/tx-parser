package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const defaultRPCEndpoint = "https://ethereum-rpc.publicnode.com"

func NewEthClient() Client {
	return &ethClient{
		endpoint: defaultRPCEndpoint,
		client:   &http.Client{},
	}
}

func (c *ethClient) call(method string, params []interface{}) (*RPCResponse, error) {
	request := RPCRequest{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResponse RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return nil, err
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResponse.Error.Message)
	}

	return &rpcResponse, nil
}

func (c *ethClient) GetLatestBlockNumber() (int64, error) {
	response, err := c.call("eth_blockNumber", nil)
	if err != nil {
		return 0, err
	}

	hexBlock, ok := response.Result.(string)
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	// Convert hex string to int64
	var blockNum int64
	fmt.Sscanf(hexBlock, "0x%x", &blockNum)
	return blockNum, nil
}

func (c *ethClient) GetBlockByNumber(blockNum int64) (Block, error) {
	blockHex := fmt.Sprintf("0x%x", blockNum)
	response, err := c.call("eth_getBlockByNumber", []interface{}{blockHex, true})
	if err != nil {
		return Block{}, err
	}

	blockData, ok := response.Result.(map[string]interface{})
	if !ok {
		return Block{}, fmt.Errorf("invalid block format")
	}

	// Convert the generic map to our Block structure
	block := Block{}

	// Convert hex values to integers where needed
	if numStr, ok := blockData["number"].(string); ok {
		block.Number, _ = hexToInt64(numStr)
	}
	if timeStr, ok := blockData["timestamp"].(string); ok {
		block.Timestamp, _ = hexToInt64(timeStr)
	}
	if gasUsedStr, ok := blockData["gasUsed"].(string); ok {
		block.GasUsed, _ = hexToInt64(gasUsedStr)
	}
	if gasLimitStr, ok := blockData["gasLimit"].(string); ok {
		block.GasLimit, _ = hexToInt64(gasLimitStr)
	}
	if baseFeeStr, ok := blockData["baseFeePerGas"].(string); ok {
		block.BaseFeePerGas, _ = hexToInt64(baseFeeStr)
	}

	// Copy string values directly
	block.Hash, _ = blockData["hash"].(string)
	block.ParentHash, _ = blockData["parentHash"].(string)
	block.StateRoot, _ = blockData["stateRoot"].(string)
	block.TransactionsRoot, _ = blockData["transactionsRoot"].(string)
	block.ReceiptsRoot, _ = blockData["receiptsRoot"].(string)
	block.Miner, _ = blockData["miner"].(string)

	// Handle transactions array
	if txs, ok := blockData["transactions"].([]interface{}); ok {
		block.Transactions = make([]Transaction, len(txs))
		for i, tx := range txs {
			if txObj, ok := tx.(map[string]interface{}); ok {
				transaction := Transaction{}

				// Convert hex values to integers
				if nonceStr, ok := txObj["nonce"].(string); ok {
					transaction.Nonce, _ = hexToInt64(nonceStr)
				}
				if gasStr, ok := txObj["gas"].(string); ok {
					transaction.Gas, _ = hexToInt64(gasStr)
				}
				if gasPriceStr, ok := txObj["gasPrice"].(string); ok {
					transaction.GasPrice, _ = hexToInt64(gasPriceStr)
				}
				if valueStr, ok := txObj["value"].(string); ok {
					transaction.Value, _ = hexToInt64(valueStr)
				}

				// Copy string values
				transaction.Hash, _ = txObj["hash"].(string)
				transaction.From, _ = txObj["from"].(string)
				transaction.To, _ = txObj["to"].(string)
				transaction.Input, _ = txObj["input"].(string)

				block.Transactions[i] = transaction
			}
		}
	}

	return block, nil
}

// Helper function to convert hex string to int64
func hexToInt64(hex string) (int64, error) {
	if hex == "" {
		return 0, nil
	}
	// Remove "0x" prefix if present
	hex = strings.TrimPrefix(hex, "0x")

	return strconv.ParseInt(hex, 16, 64)
}
