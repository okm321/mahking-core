package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/okm321/mahking-go/config"
	pkgerror "github.com/okm321/mahking-go/pkg/error"
	"github.com/okm321/mahking-go/pkg/logger"
)

func Run(cfg *config.Config, router Router) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Address, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 別goroutineでサーバー起動
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.InfofContext(context.Background(), "shutting down the server: %v", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	<-quit

	logger.InfoContext(context.Background(), "Sever is shutting down...")

	// Graceful Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return pkgerror.Errorf("graceful shut down failed: %w", err)
	}

	logger.InfoContext(ctx, "graceful shut down success")
	return nil
}
