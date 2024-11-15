package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	mockParser "github.com/vdhieu/tx-parser/mocks/internal_/parser"
)

func TestParserHandler_Subscribe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEthParser := mockParser.NewParser(t)

	tests := []struct {
		name       string
		reqBody    SubscribeRequest
		setupMock  func(m *mockParser.Parser)
		wantStatus int
		wantBody   *SubscribeResponse
	}{
		{
			name: "successful subscription",
			reqBody: SubscribeRequest{
				Address: "0x1234",
			},
			setupMock: func(m *mockParser.Parser) {
				m.On("Subscribe", "0x1234").Return(true)
			},
			wantStatus: http.StatusOK,
			wantBody: &SubscribeResponse{
				Message: "successfully subscribed",
			},
		},
		{
			name: "unable to subscribe address",
			reqBody: SubscribeRequest{
				Address: "0x12345",
			},
			setupMock: func(m *mockParser.Parser) {
				m.On("Subscribe", "0x12345").Return(false)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: &SubscribeResponse{
				Message: "unable to subscribe",
			},
		},
		{
			name:       "invalid request - missing address",
			reqBody:    SubscribeRequest{},
			setupMock:  func(m *mockParser.Parser) {},
			wantStatus: http.StatusBadRequest,
			wantBody: &SubscribeResponse{
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
			router.POST("/subscribe", h.Subscribe)

			jsonBody, _ := json.Marshal(tt.reqBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantBody != nil {
				var got SubscribeResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, *tt.wantBody, got)
			}
		})
	}
}
