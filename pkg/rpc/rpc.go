package rpc

import "net/http"

type ethClient struct {
	endpoint string
	client   *http.Client
}

type Client interface {
	GetLatestBlockNumber() (int64, error)
	GetBlockByNumber(blockNum int64) (Block, error)
}
