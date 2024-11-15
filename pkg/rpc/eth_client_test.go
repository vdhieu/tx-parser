package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEthClient(t *testing.T) {
	got := NewEthClient()
	require.NotNil(t, got)

}

func Test_ethClient_GetLatestBlockNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req RPCRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		require.Equal(t, "eth_blockNumber", req.Method)
		require.Equal(t, "2.0", req.JsonRPC)
		require.Len(t, req.Params, 0)

		response := RPCResponse{
			JsonRPC: "2.0",
			ID:      req.ID,
			Result:  "0x1234",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	type fields struct {
		endpoint string
		client   *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "successful get block number",
			fields: fields{
				endpoint: server.URL,
				client:   &http.Client{},
			},
			want:    4660,
			wantErr: false,
		},
		{
			name: "connection error",
			fields: fields{
				endpoint: "invalid-url",
				client:   &http.Client{},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ethClient{
				endpoint: tt.fields.endpoint,
				client:   tt.fields.client,
			}
			got, err := c.GetLatestBlockNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("ethClient.GetLatestBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_ethClient_GetBlockByNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request method and content type
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Read the request body
		var req RPCRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// Verify the RPC request
		require.Equal(t, "eth_getBlockByNumber", req.Method)
		require.Equal(t, "2.0", req.JsonRPC)
		require.Len(t, req.Params, 2)

		// Return mock response based on block number
		blockNum := req.Params[0].(string)
		if blockNum == "0x0" {
			response := RPCResponse{
				JsonRPC: "2.0",
				ID:      req.ID,
				Result: Block{
					Number:     0,
					Hash:       "0x123",
					ParentHash: "0x123",
					Transactions: []Transaction{
						{
							Hash: "0x123",
							From: "0x123",
							To:   "0x123",
						},
					},
					StateRoot:        "0x123",
					TransactionsRoot: "0x123",
				},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			// Return error for invalid block numbers
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RPCResponse{
				JsonRPC: "2.0",
				ID:      req.ID,
				Error:   &RPCError{Code: -32000, Message: "block not found"},
			})
		}
	}))
	defer server.Close()

	type fields struct {
		endpoint string
		client   *http.Client
	}
	type args struct {
		blockNum int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Block
		wantErr bool
	}{
		{
			name: "successful get block",
			fields: fields{
				endpoint: server.URL,
				client:   &http.Client{},
			},
			args: args{
				blockNum: 0,
			},
			want: Block{
				Number:     0,
				Hash:       "0x123",
				ParentHash: "0x123",
				Timestamp:  0,
				Transactions: []Transaction{
					{
						Hash: "0x123",
						From: "0x123",
						To:   "0x123",
					},
				},
				StateRoot:        "0x123",
				TransactionsRoot: "0x123",
			},
			wantErr: false,
		},
		{
			name: "block not found",
			fields: fields{
				endpoint: server.URL,
				client:   &http.Client{},
			},
			args: args{
				blockNum: -1,
			},
			want:    Block{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ethClient{
				endpoint: tt.fields.endpoint,
				client:   tt.fields.client,
			}
			got, err := c.GetBlockByNumber(tt.args.blockNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ethClient.GetBlockByNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_hexToInt64(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "valid hex",
			args: args{
				hex: "0x1234",
			},
			want:    4660,
			wantErr: false,
		},
		{
			name: "zero hex",
			args: args{
				hex: "0x0",
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "no prefix",
			args: args{
				hex: "1234",
			},
			want:    4660,
			wantErr: false,
		},
		{
			name: "invalid hex characters",
			args: args{
				hex: "0xGHIJ",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hexToInt64(tt.args.hex)
			fmt.Println("got", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("hexToInt64() %v error = %v, wantErr %v", got, err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
