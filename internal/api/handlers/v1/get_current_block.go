package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ParserHandler) GetCurrentBlock(c *gin.Context) {
	block := h.parser.GetCurrentBlock()
	c.JSON(http.StatusOK, BlockResponse{Block: block})
}
