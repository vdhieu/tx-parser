package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	router "github.com/vdhieu/tx-parser/internal/api"
	"github.com/vdhieu/tx-parser/internal/parser"
	"github.com/vdhieu/tx-parser/internal/storage"
	"github.com/vdhieu/tx-parser/pkg/logger"
	"github.com/vdhieu/tx-parser/pkg/notification"
	"github.com/vdhieu/tx-parser/pkg/rpc"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.Initialize()
	defer logger.Sync()

	// Initialize components
	p := parser.NewEthParser(
		storage.NewMemoryStorage(),
		rpc.NewEthClient(),
		notification.NewConsoleNotifier(),
	)

	// Setup router
	r := router.SetupRouter(p)

	// Create server
	srv := &http.Server{
		Addr:    ":5005",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatal("Failed to start server: %v", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.GetLogger().Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop the parser
	p.Shutdown()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.GetLogger().Error("Server forced to shutdown:", zap.Error(err))
	}

	logger.GetLogger().Info("Server exited properly")
}
