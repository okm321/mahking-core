package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okm321/mahking-go/config"
	"github.com/okm321/mahking-go/internal/application"
	"github.com/okm321/mahking-go/internal/infrastructure/postgres"
	"github.com/okm321/mahking-go/internal/presentation/api"
	pkgpostgres "github.com/okm321/mahking-go/pkg/postgres"
)

// App wires all layers together and owns their lifecycle.
type App struct {
	cfg    *config.Config
	router api.Router

	pool *pgxpool.Pool
}

// NewApp creates all dependencies and returns a runnable App.
func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	// Infra
	pool, err := pkgpostgres.Connect(cfg.DBPostgres)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	// Repository
	groupRepo := postgres.NewGroupRepository(pool)

	// Usecase
	groupUsecase := application.NewGroupUsecase(&application.NewGroupUsecaseArgs{
		GroupRepo: groupRepo,
	})

	// Handler
	groupHandler := api.NewGroupHandler(groupUsecase)

	// Router
	router := api.NewRouter(api.HandlerSet{
		Group: groupHandler,
	})

	return &App{
		cfg:    cfg,
		router: router,
		pool:   pool,
	}, nil
}

// Run starts the HTTP server. It blocks until shutdown.
func (a *App) Run() error {
	return api.Run(a.cfg, a.router)
}

// Close releases resources.
func (a *App) Close() {
	if a.pool != nil {
		a.pool.Close()
	}
}
