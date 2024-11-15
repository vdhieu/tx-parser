package router

import (
	"github.com/gin-gonic/gin"
	handler "github.com/vdhieu/tx-parser/internal/api/handlers/v1"
	"github.com/vdhieu/tx-parser/internal/parser"
)

func SetupRouter(p parser.Parser) *gin.Engine {
	r := gin.Default()
	h := handler.NewParserHandler(p)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/block/current", h.GetCurrentBlock)
		v1.POST("/subscribe", h.Subscribe)
		v1.GET("/transactions", h.GetTransactions)
	}

	return r
}
