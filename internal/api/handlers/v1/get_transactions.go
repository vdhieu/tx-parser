package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ParserHandler) GetTransactions(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, TransactionsResponse{Error: "address is required"})
		return
	}

	transactions := h.parser.GetTransactions(address)
	c.JSON(http.StatusOK, TransactionsResponse{
		Data: transactions,
	})
}
