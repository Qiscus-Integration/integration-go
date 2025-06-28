package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog/log"
)

func Load() *Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		log.Fatal().Msgf("unable to parse env: %s", err.Error())
	}

	return &c
}

type Config struct {
	App      App
	Database Database
	Redis    Redis
	Qiscus   Qiscus
}

type App struct {
	SecretKey string `env:"APP_SECRET_KEY,required"`
}

type Database struct {
	Host     string `env:"DATABASE_HOST,required"`
	Port     int    `env:"DATABASE_PORT,required"`
	User     string `env:"DATABASE_USER,required"`
	Password string `env:"DATABASE_PASSWORD,required"`
	Name     string `env:"DATABASE_NAME,required"`
	LogLevel string `env:"DATABASE_LOG_LEVEL,required"`
}

func (d Database) DataSourceName() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

type Redis struct {
	URL string `env:"REDIS_URL,required"`
}

type Qiscus struct {
	AppID       string `env:"QISCUS_APP_ID,required"`
	SecretKey   string `env:"QISCUS_SECRET_KEY,required"`
	Omnichannel Omnichannel
}

type Omnichannel struct {
	URL string `env:"QISCUS_OMNICHANNEL_URL,required"`
}
