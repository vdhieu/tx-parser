package handler

import (
	"github.com/vdhieu/tx-parser/internal/parser"
)

type ParserHandler struct {
	parser parser.Parser
}

func NewParserHandler(p parser.Parser) *ParserHandler {
	return &ParserHandler{parser: p}
}
