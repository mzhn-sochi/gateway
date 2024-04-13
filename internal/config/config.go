package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Host string `env:"APP_HOST" env-default:"0.0.0.0"`
		Port int    `env:"APP_PORT" env-default:"8080"`
	}
	Services struct {
		PriceTagAnalyzer struct {
			Host string `env:"PRICE_TAG_ANALYZER_HOST" env-default:"77.221.158.75"`
			Port int    `env:"PRICE_TAG_ANALYZER_PORT" env-default:"50051"`
		}

		Suggestions struct {
			Host string `env:"SUGGESTIONS_HOST" env-default:"localhost"`
			Port int    `env:"SUGGESTIONS_PORT" env-default:"2020"`
		}

		TicketService struct {
			Host string `env:"TICKET_SERVICE_HOST" env-default:"localhost"`
			Port int    `env:"TICKET_SERVICE_PORT" env-default:"50052"`
		}

		AuthService struct {
			Host string `env:"AUTH_SERVICE_HOST" env-default:"localhost"`
			Port int    `env:"AUTH_SERVICE_PORT" env-default:"50053"`
		}

		S3Service struct {
			Host string `env:"S3_SERVICE_HOST" env-default:"localhost"`
			Port int    `env:"S3_SERVICE_PORT" env-default:"50054"`
		}
	}

	LogLevel string `env:"LOG_LEVEL" env-default:"debug"`
}

func New() *Config {
	config := &Config{}
	if err := cleanenv.ReadEnv(config); err != nil {
		header := "SOCHYA GATEWAY"
		f := cleanenv.FUsage(os.Stdout, config, &header)
		f()
		panic(err)
	}

	return config
}
