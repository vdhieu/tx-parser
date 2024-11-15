package rpc

type RPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Block represents an Ethereum block structure
type Block struct {
	Number           int64         `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Timestamp        int64         `json:"timestamp"`
	Transactions     []Transaction `json:"transactions"`
	StateRoot        string        `json:"stateRoot"`
	TransactionsRoot string        `json:"transactionsRoot"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	GasUsed          int64         `json:"gasUsed"`
	GasLimit         int64         `json:"gasLimit"`
	BaseFeePerGas    int64         `json:"baseFeePerGas"`
	Miner            string        `json:"miner"`
}

type Transaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Nonce       int64  `json:"nonce"`
	Gas         int64  `json:"gas"`
	GasPrice    int64  `json:"gas_price"`
	Value       int64  `json:"value"`
	Input       string `json:"input"`
	Status      int64  `json:"status"`
	BlockHash   string `json:"block_hash"`
	BlockNumber int64  `json:"block_number"`
}

type TransactionReceipt struct {
	Status      int64  `json:"status"`
	BlockHash   string `json:"block_hash"`
	BlockNumber int64  `json:"block_number"`
}
