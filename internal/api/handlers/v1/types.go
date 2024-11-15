package handler

import "github.com/vdhieu/tx-parser/internal/models"

type BlockResponse struct {
	Block int `json:"block"`
}

type SubscribeRequest struct {
	Address string `json:"address" binding:"required"`
}

type SubscribeResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type TransactionsResponse struct {
	Error string               `json:"error,omitempty"`
	Data  []models.Transaction `json:"data"`
}
