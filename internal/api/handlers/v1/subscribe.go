package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ParserHandler) Subscribe(c *gin.Context) {
	var req SubscribeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SubscribeResponse{Error: "address is required"})
		return
	}

	success := h.parser.Subscribe(req.Address)
	if success {
		c.JSON(http.StatusOK, SubscribeResponse{Message: "successfully subscribed"})
		return
	}

	c.JSON(http.StatusInternalServerError, SubscribeResponse{Message: "unable to subscribe"})
}
