package handler

import (
	"reflect"
	"testing"

	"github.com/vdhieu/tx-parser/internal/parser"
	mockParser "github.com/vdhieu/tx-parser/mocks/internal_/parser"
)

func TestNewParserHandler(t *testing.T) {
	mockEthParser := mockParser.NewParser(t)
	type args struct {
		p parser.Parser
	}
	tests := []struct {
		name string
		args args
		want *ParserHandler
	}{
		{
			name: "successful initialization",
			args: args{
				p: mockEthParser,
			},
			want: &ParserHandler{
				parser: mockEthParser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParserHandler(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
