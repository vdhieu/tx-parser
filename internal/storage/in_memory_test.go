package storage

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdhieu/tx-parser/internal/models"
)

func TestNewMemoryStorage(t *testing.T) {
	got := NewMemoryStorage()
	require.NotNil(t, got)

}

func Test_memoryStorage_AddSubscriber(t *testing.T) {
	type fields struct {
		currentBlock int64
		subscribers  map[string]bool
		transactions map[string][]models.Transaction
		mu           sync.RWMutex
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add new subscriber",
			fields: fields{
				subscribers:  make(map[string]bool),
				transactions: make(map[string][]models.Transaction),
			},
			args: args{
				address: "0x123",
			},
			wantErr: false,
		},
		{
			name: "add existing subscriber",
			fields: fields{
				subscribers: map[string]bool{
					"0x123": true,
				},
				transactions: make(map[string][]models.Transaction),
			},
			args: args{
				address: "0x123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &memoryStorage{
				currentBlock: tt.fields.currentBlock,
				subscribers:  tt.fields.subscribers,
				transactions: tt.fields.transactions,
				mu:           tt.fields.mu,
			}
			if err := s.AddSubscriber(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.AddSubscriber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memoryStorage_IsSubscribed(t *testing.T) {
	type fields struct {
		currentBlock int64
		subscribers  map[string]bool
		transactions map[string][]models.Transaction
		mu           sync.RWMutex
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
			name: "check existing subscriber",
			fields: fields{
				subscribers: map[string]bool{
					"0x123": true,
				},
			},
			args: args{
				address: "0x123",
			},
			want: true,
		},
		{
			name: "check non-existing subscriber",
			fields: fields{
				subscribers: map[string]bool{},
			},
			args: args{
				address: "0x123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &memoryStorage{
				currentBlock: tt.fields.currentBlock,
				subscribers:  tt.fields.subscribers,
				transactions: tt.fields.transactions,
				mu:           tt.fields.mu,
			}
			if got := s.IsSubscribed(tt.args.address); got != tt.want {
				t.Errorf("memoryStorage.IsSubscribed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryStorage_GetSubscribers(t *testing.T) {
	type fields struct {
		currentBlock int64
		subscribers  map[string]bool
		transactions map[string][]models.Transaction
		mu           sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "get empty subscribers",
			fields: fields{
				subscribers: map[string]bool{},
			},
			want: []string{},
		},
		{
			name: "get multiple subscribers",
			fields: fields{
				subscribers: map[string]bool{
					"0x123": true,
					"0x456": true,
				},
			},
			want: []string{"0x123", "0x456"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &memoryStorage{
				currentBlock: tt.fields.currentBlock,
				subscribers:  tt.fields.subscribers,
				transactions: tt.fields.transactions,
				mu:           tt.fields.mu,
			}
			if got := s.GetSubscribers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryStorage.GetSubscribers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryStorage_SaveTransactions(t *testing.T) {
	type fields struct {
		currentBlock int64
		subscribers  map[string]bool
		transactions map[string][]models.Transaction
		mu           sync.RWMutex
	}
	type args struct {
		address string
		txs     []models.Transaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "save transactions for subscribed address",
			fields: fields{
				subscribers: map[string]bool{
					"0x123": true,
				},
				transactions: make(map[string][]models.Transaction),
			},
			args: args{
				address: "0x123",
				txs: []models.Transaction{
					{Hash: "tx1", From: "0x123", To: "0x456"},
				},
			},
			wantErr: false,
		},
		{
			name: "save transactions for non-subscribed address",
			fields: fields{
				subscribers:  make(map[string]bool),
				transactions: make(map[string][]models.Transaction),
			},
			args: args{
				address: "0x123",
				txs: []models.Transaction{
					{Hash: "tx1", From: "0x123", To: "0x456"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &memoryStorage{
				currentBlock: tt.fields.currentBlock,
				subscribers:  tt.fields.subscribers,
				transactions: tt.fields.transactions,
				mu:           tt.fields.mu,
			}
			if err := s.SaveTransactions(tt.args.address, tt.args.txs); (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.SaveTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memoryStorage_GetTransactions(t *testing.T) {
	sampleTx := []models.Transaction{{Hash: "tx1", From: "0x123", To: "0x456"}}

	type fields struct {
		currentBlock int64
		subscribers  map[string]bool
		transactions map[string][]models.Transaction
		mu           sync.RWMutex
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Transaction
		wantErr bool
	}{
		{
			name: "get transactions for subscribed address",
			fields: fields{
				subscribers: map[string]bool{
					"0x123": true,
				},
				transactions: map[string][]models.Transaction{
					"0x123": sampleTx,
				},
			},
			args: args{
				address: "0x123",
			},
			want:    sampleTx,
			wantErr: false,
		},
		{
			name: "get transactions for non-subscribed address",
			fields: fields{
				subscribers:  make(map[string]bool),
				transactions: make(map[string][]models.Transaction),
			},
			args: args{
				address: "0x123",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &memoryStorage{
				currentBlock: tt.fields.currentBlock,
				subscribers:  tt.fields.subscribers,
				transactions: tt.fields.transactions,
				mu:           tt.fields.mu,
			}
			got, err := s.GetTransactions(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, got, tt.want)
		})
	}
}

func Test_memoryStorage_SetCurrentBlock(t *testing.T) {
	type fields struct {
		currentBlock int64
		subscribers  map[string]bool
		transactions map[string][]models.Transaction
		mu           sync.RWMutex
	}
	type args struct {
		blockNum int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set block number",
			fields: fields{
				currentBlock: 0,
			},
			args: args{
				blockNum: 100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &memoryStorage{
				currentBlock: tt.fields.currentBlock,
				subscribers:  tt.fields.subscribers,
				transactions: tt.fields.transactions,
				mu:           tt.fields.mu,
			}
			if err := s.SetCurrentBlock(tt.args.blockNum); (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.SetCurrentBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memoryStorage_GetCurrentBlock(t *testing.T) {
	type fields struct {
		currentBlock int64
		subscribers  map[string]bool
		transactions map[string][]models.Transaction
		mu           sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "get current block",
			fields: fields{
				currentBlock: 100,
			},
			want:    100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &memoryStorage{
				currentBlock: tt.fields.currentBlock,
				subscribers:  tt.fields.subscribers,
				transactions: tt.fields.transactions,
				mu:           tt.fields.mu,
			}
			got, err := s.GetCurrentBlock()
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryStorage.GetCurrentBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, got, tt.want)
		})
	}
}
