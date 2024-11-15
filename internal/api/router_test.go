package router

import (
	"testing"

	"github.com/vdhieu/tx-parser/internal/parser"
	mockParser "github.com/vdhieu/tx-parser/mocks/internal_/parser"
)

func TestSetupRouter(t *testing.T) {
	type args struct {
		p parser.Parser
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful router setup",
			args: args{
				p: mockParser.NewParser(t),
			},
			wantErr: false,
		},
	}

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/block/current"},
		{"POST", "/api/v1/subscribe"},
		{"GET", "/api/v1/transactions"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetupRouter(tt.args.p)
			if router == nil {
				t.Error("SetupRouter() returned nil router")
			}

			// Get all routes
			routes := router.Routes()

			// Create a map for easy route checking
			routeMap := make(map[string]bool)
			for _, route := range routes {
				key := route.Method + ":" + route.Path
				routeMap[key] = true
			}

			// Verify all expected routes exist
			for _, expectedRoute := range expectedRoutes {
				routeKey := expectedRoute.method + ":" + expectedRoute.path
				if !routeMap[routeKey] {
					t.Errorf("Expected route %s %s not found",
						expectedRoute.method,
						expectedRoute.path)
				}
			}
		})
	}
}
