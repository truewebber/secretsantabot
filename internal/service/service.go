package service

import (
	"github.com/truewebber/secretsantabot/internal/app"
	"github.com/truewebber/secretsantabot/internal/log"
)

type (
	Config struct {
		Logger      log.Logger
		ChatService ChatService
	}

	ChatService struct {
		RedisURI string
	}
)

func NewConfig(logger log.Logger) *Config {
	return &Config{
		Logger: logger,
	}
}

func MustNewApplication(cfg *Config) *app.Application {
	return &app.Application{
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}
