package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/vdhieu/tx-parser/internal/models"
	mockParser "github.com/vdhieu/tx-parser/mocks/internal_/parser"
)

func TestParserHandler_GetTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEthParser := mockParser.NewParser(t)
	sampleTxs := []models.Transaction{
		{
			Hash:  "0xabc123",
			From:  "0x111",
			To:    "0x222",
			Value: "1000000000000000000",
		},
	}

	tests := []struct {
		name       string
		address    string
		setupMock  func(m *mockParser.Parser)
		wantStatus int
		wantBody   *TransactionsResponse
	}{
		{
			name:    "successful get transactions",
			address: "0x1234",
			setupMock: func(m *mockParser.Parser) {
				m.On("GetTransactions", "0x1234").Return(sampleTxs, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: &TransactionsResponse{
				Data: sampleTxs,
			},
		},
		{
			name:    "invalid address",
			address: "",
			setupMock: func(m *mockParser.Parser) {
			},
			wantStatus: http.StatusBadRequest,
			wantBody: &TransactionsResponse{
				Error: "address is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setupMock != nil {
				tt.setupMock(mockEthParser)
			}

			h := &ParserHandler{
				parser: mockEthParser,
			}

			router := gin.New()
			router.GET("/transactions", h.GetTransactions)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/transactions?address="+tt.address, nil)

			router.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != nil {
				var got TransactionsResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, *tt.wantBody, got)
			}
		})
	}
}
