package notification

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConsoleNotifier(t *testing.T) {
	got := NewConsoleNotifier()
	require.NotNil(t, got)
}

func Test_consoleNotifier_Notify(t *testing.T) {
	type args struct {
		address string
		message string
		payload interface{}
	}
	tests := []struct {
		name    string
		n       *consoleNotifier
		args    args
		wantErr bool
	}{
		{
			name: "successful notification with string payload",
			n:    &consoleNotifier{},
			args: args{
				address: "0x123",
				message: "Test message",
				payload: "string payload",
			},
			wantErr: false,
		},
		{
			name: "successful notification with struct payload",
			n:    &consoleNotifier{},
			args: args{
				address: "0x123",
				message: "Test message",
				payload: struct {
					Key   string
					Value int
				}{
					Key:   "test",
					Value: 123,
				},
			},
			wantErr: false,
		},
		{
			name: "successful notification with nil payload",
			n:    &consoleNotifier{},
			args: args{
				address: "0x123",
				message: "Test message",
				payload: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &consoleNotifier{}
			if err := n.Notify(tt.args.address, tt.args.message, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("consoleNotifier.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
