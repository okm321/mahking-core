package config

import (
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	pkgpostgres "github.com/okm321/mahking-go/pkg/postgres"
)

type Config struct {
	Server     Server
	GCP        GCP
	DBPostgres pkgpostgres.DB
}

// Server configurations.
type Server struct {
	Address         string        `env:"ADDRESS" envDefault:""`
	Port            string        `env:"PORT" envDefault:"8080"`
	Debug           bool          `env:"DEBUG"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

// GCP configurations.
type GCP struct {
	ProjectID                       string `envconfig:"GCP_PROJECT_ID" default:"mahking-dev"`
	HTTPLBSourceIP                  string `envconfig:"LITE_CORE_GCP_HTTP_LB_SOURCE_IP"`
	BatchInvokerServiceAccountEmail string `envconfig:"LITE_CORE_GCP_BATCH_INVOKER_SERVICE_ACCOUNT_EMAIL" default:"cloud-run-service-invoker-dev@casting-one-dev.iam.gserviceaccount.com"`
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		if err := env.Parse(&config); err != nil {
			panic(err)
		}
	})

	return &config
}
