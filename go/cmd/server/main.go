package main

import (
	"context"
	"fmt"

	"github.com/okm321/mahking-go/config"
	"github.com/okm321/mahking-go/internal/bootstrap"
	"github.com/okm321/mahking-go/pkg/logger"
	pkgtrace "github.com/okm321/mahking-go/pkg/trace"
)

func main() {
	ctx := context.Background()
	cnf := config.Get()
	logger.Init(cnf.GCP.ProjectID, cnf.Server.Debug)

	traceShutdown, err := pkgtrace.Init(ctx, pkgtrace.Config{
		ServiceName:    cnf.Telemetry.ServiceName,
		ServiceVersion: cnf.Telemetry.ServiceVersion,
		Environment:    cnf.Telemetry.Environment,
		SampleRate:     cnf.Telemetry.SampleRate,
		ProjectID:      cnf.GCP.ProjectID,
		Debug:          cnf.Server.Debug,
	})
	if err != nil {
		logger.FatalContext(ctx, fmt.Sprintf("trace init failed: %v", err))
	}
	defer func() {
		if err := traceShutdown(ctx); err != nil {
			logger.ErrorContext(ctx, fmt.Sprintf("trace shutdown failed: %v", err))
		}
	}()

	logger.InfoContext(ctx, "mahking-go started!")

	app, err := bootstrap.NewApp(ctx, cnf)
	if err != nil {
		logger.FatalContext(ctx, fmt.Sprintf("bootstrap failed: %v", err))
	}
	defer app.Close()

	if err := app.Run(); err != nil {
		logger.FatalContext(ctx, fmt.Sprintf("サーバー起動失敗: %v", err))
	}
}
