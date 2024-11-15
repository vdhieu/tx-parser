package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	mockParser "github.com/vdhieu/tx-parser/mocks/internal_/parser"
)

func TestParserHandler_GetCurrentBlock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEthParser := mockParser.NewParser(t)

	tests := []struct {
		name       string
		setupMock  func(m *mockParser.Parser)
		wantStatus int
		wantBody   *BlockResponse
	}{
		{
			name: "successful get current block",
			setupMock: func(m *mockParser.Parser) {
				m.On("GetCurrentBlock").Return(int(1), nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &BlockResponse{Block: 1},
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
			router.GET("/current-block", h.GetCurrentBlock)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/current-block", nil)

			router.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != nil {
				var got BlockResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, *tt.wantBody, got)
			}
		})
	}
}
