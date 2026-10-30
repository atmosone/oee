package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/atmosone/oee/internal/adapters/openapi"
	"github.com/atmosone/oee/internal/adapters/postgres"
	"github.com/atmosone/oee/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(
		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
	))

	db, err := postgres.Open("")
	if err != nil {
		logger.Error("failed to connect postgres",
			slog.Any("error", err))
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close postgres",
				slog.Any("error", err))
			return
		}
		logger.Debug("postgres closed")
	}()

	app := app.New(logger, db)

	var wg sync.WaitGroup
	edge := openapi.NewEdgeServer(":8880", app)
	qualification := openapi.NewQualificationServer(":8882", app)
	wg.Go(func() { edge.ListenAndServe() })
	wg.Go(func() { qualification.ListenAndServe() })
	<-ctx.Done()
	app.Close()
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer shutdownCancel()
	if err := edge.Shutdown(shutdownCtx); err != nil {
		logger.With(slog.Any("error", err)).
			Error("failed to shutdown edge server")
		return
	}
	logger.Info("edge server closed")
	if err := qualification.Shutdown(shutdownCtx); err != nil {
		logger.With(slog.Any("error", err)).
			Error("failed to shutdown qualification server")
		return
	}
	logger.Info("qualification server closed")
	wg.Wait()
}
