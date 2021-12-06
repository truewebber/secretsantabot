package service

import (
	"github.com/truewebber/secretsantabot/internal/app"
	"github.com/truewebber/secretsantabot/internal/app/command"
	"github.com/truewebber/secretsantabot/internal/app/query"
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
		Commands: app.Commands{
			RegisterNewChat: command.MustNewRegisterNewChatHandler(cfg.Logger),
			Enroll:          command.MustNewEnrollHandler(cfg.Logger),
			DisEnroll:       command.MustNewDisEnrollHandler(cfg.Logger),
			Magic:           command.MustNewMagicHandler(cfg.Logger),
		},
		Queries: app.Queries{
			GetMyReceiver: query.MustNewGetMyReceiverHandler(cfg.Logger),
			List:          query.MustNewListHandler(cfg.Logger),
		},
	}
}
