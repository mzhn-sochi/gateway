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
