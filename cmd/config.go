package cmd

import "github.com/caarlos0/env/v6"

type Config struct {
	Port int    `env:"PORT,notEmpty"`
	Env  string `env:"ENVIRONMENT,notEmpty"`

	DatabaseUrl          string `env:"DATABASE_URL,notEmpty"`
	DbMaxIdleConnections int    `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	DbMaxOpenConnections int    `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"5"`
}

func ParseConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		panic(err)
	}

	return cfg
}
