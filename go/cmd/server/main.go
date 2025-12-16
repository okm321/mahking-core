package main

import (
	"context"
	"fmt"

	"github.com/okm321/mahking-go/config"
	"github.com/okm321/mahking-go/internal/bootstrap"
	"github.com/okm321/mahking-go/pkg/logger"
)

func main() {
	ctx := context.Background()
	cnf := config.Get()
	logger.Init(cnf.GCP.ProjectID)

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
