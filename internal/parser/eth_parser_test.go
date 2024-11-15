package parser

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdhieu/tx-parser/internal/models"
	"github.com/vdhieu/tx-parser/internal/storage"
	mockStorage "github.com/vdhieu/tx-parser/mocks/internal_/storage"
	mockNoti "github.com/vdhieu/tx-parser/mocks/pkg/notification"
	mockClient "github.com/vdhieu/tx-parser/mocks/pkg/rpc"
	"github.com/vdhieu/tx-parser/pkg/notification"
	"github.com/vdhieu/tx-parser/pkg/rpc"
	"go.uber.org/zap"
)

func setupMocks(t *testing.T) (*mockStorage.Storage, *mockClient.Client, *mockNoti.Notifier) {
	mockStorage := mockStorage.NewStorage(t)
	mockClient := mockClient.NewClient(t)
	mockNotifier := mockNoti.NewNotifier(t)
	return mockStorage, mockClient, mockNotifier
}

func TestNewEthParser(t *testing.T) {
	mockStorage, mockClient, mockNotifier := setupMocks(t)

	got := NewEthParser(mockStorage, mockClient, mockNotifier)
	require.NotNil(t, got)
}

func Test_ethParser_Shutdown(t *testing.T) {
	mockStorage, mockClient, mockNotifier := setupMocks(t)
	type fields struct {
		storage  storage.Storage
		client   rpc.Client
		log      *zap.Logger
		notifier notification.Notifier
		running  bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "successful shutdown",
			fields: fields{
				storage:  mockStorage,
				client:   mockClient,
				log:      zap.NewNop(),
				notifier: mockNotifier,
				running:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ethParser{
				storage:  tt.fields.storage,
				client:   tt.fields.client,
				log:      tt.fields.log,
				notifier: tt.fields.notifier,
				running:  tt.fields.running,
			}
			p.Shutdown()
			require.False(t, p.running)
		})
	}
}

func Test_ethParser_GetCurrentBlock(t *testing.T) {
	mockStorage, mockClient, mockNotifier := setupMocks(t)

	mockStorage.On("GetCurrentBlock").Return(int64(100), nil)

	type fields struct {
		storage  storage.Storage
		client   rpc.Client
		log      *zap.Logger
		notifier notification.Notifier
		running  bool
	}

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "get current block",
			fields: fields{
				storage:  mockStorage,
				client:   mockClient,
				log:      zap.NewNop(),
				notifier: mockNotifier,
				running:  true,
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ethParser{
				storage:  tt.fields.storage,
				client:   tt.fields.client,
				log:      tt.fields.log,
				notifier: tt.fields.notifier,
				running:  tt.fields.running,
			}
			got := p.GetCurrentBlock()
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_ethParser_Subscribe(t *testing.T) {
	mockStorage, mockClient, mockNotifier := setupMocks(t)

	address := "0x123"

	mockStorage.On("AddSubscriber", address).Return(nil)

	type fields struct {
		storage  storage.Storage
		client   rpc.Client
		log      *zap.Logger
		notifier notification.Notifier
		running  bool
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "successful subscription",
			fields: fields{
				storage:  mockStorage,
				client:   mockClient,
				log:      zap.NewNop(),
				notifier: mockNotifier,
				running:  true,
			},
			args: args{
				address: address,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ethParser{
				storage:  tt.fields.storage,
				client:   tt.fields.client,
				log:      tt.fields.log,
				notifier: tt.fields.notifier,
				running:  tt.fields.running,
			}
			got := p.Subscribe(tt.args.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_ethParser_GetTransactions(t *testing.T) {
	type fields struct {
		storage  storage.Storage
		client   rpc.Client
		log      *zap.Logger
		notifier notification.Notifier
		running  bool
	}
	type args struct {
		address string
	}
	mockStorage, mockClient, mockNotifier := setupMocks(t)

	address := "0x123"
	expectedTxs := []models.Transaction{
		{
			Hash:      "0xabc",
			From:      "0x123",
			To:        "0x456",
			Value:     "1000000000000000000",
			Timestamp: "123456",
		},
	}

	mockStorage.On("GetTransactions", address).Return(expectedTxs, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.Transaction
	}{
		{
			name: "get transactions",
			fields: fields{
				storage:  mockStorage,
				client:   mockClient,
				log:      zap.NewNop(),
				notifier: mockNotifier,
				running:  true,
			},
			args: args{
				address: address,
			},
			want: expectedTxs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ethParser{
				storage:  tt.fields.storage,
				client:   tt.fields.client,
				log:      tt.fields.log,
				notifier: tt.fields.notifier,
				running:  tt.fields.running,
			}
			got := p.GetTransactions(tt.args.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_ethParser_processTransactions(t *testing.T) {
	mockStorage, mockClient, mockNotifier := setupMocks(t)
	// Mock data
	subscribedAddr := "0x123"
	txHash := "0xabc"
	block := rpc.Block{
		Number: 101,
		Transactions: []rpc.Transaction{
			{
				Hash:  txHash,
				From:  subscribedAddr,
				To:    "0x456",
				Value: 1000,
			},
		},
		Timestamp: 123456,
	}

	txn := models.Transaction{
		BlockNumber: "101",
		Hash:        txHash,
		From:        subscribedAddr,
		To:          "0x456",
		Value:       "1000",
		Timestamp:   strconv.FormatInt(block.Timestamp, 10),
	}
	// Mock the necessary calls
	mockStorage.On("GetSubscribers").Return([]string{subscribedAddr}, nil)
	mockStorage.On("GetTransactions", subscribedAddr).Return([]models.Transaction{}, nil)
	mockStorage.On("SaveTransactions", subscribedAddr, []models.Transaction{txn}).Return(nil)
	mockNotifier.On("Notify", subscribedAddr, "found a new transactions", txn).Return(nil)

	type fields struct {
		storage  storage.Storage
		client   rpc.Client
		log      *zap.Logger
		notifier notification.Notifier
		running  bool
	}
	type args struct {
		block rpc.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "process transaction for subscribed address",
			fields: fields{
				storage:  mockStorage,
				client:   mockClient,
				log:      zap.NewNop(),
				notifier: mockNotifier,
				running:  true,
			},
			args: args{
				block: block,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ethParser{
				storage:  tt.fields.storage,
				client:   tt.fields.client,
				log:      tt.fields.log,
				notifier: tt.fields.notifier,
				running:  tt.fields.running,
			}
			p.processTransactions(tt.args.block)
			mockStorage.AssertExpectations(t)
			mockNotifier.AssertExpectations(t)
		})
	}
}
