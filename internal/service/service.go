package service

import (
	"github.com/truewebber/secretsantabot/internal/app"
	"github.com/truewebber/secretsantabot/internal/app/command"
	"github.com/truewebber/secretsantabot/internal/app/query"
	"github.com/truewebber/secretsantabot/internal/chat/storage/postgres"
	"github.com/truewebber/secretsantabot/internal/log"
)

type (
	Config struct {
		Logger      log.Logger
		ChatService ChatService
	}

	ChatService struct {
		PostgresURI string
	}
)

func NewConfig(logger log.Logger) *Config {
	return &Config{
		Logger: logger,
	}
}

func MustNewApplication(cfg *Config) *app.Application {
	chatService := postgres.MustNewPostgres(cfg.ChatService.PostgresURI)

	return &app.Application{
		Commands: app.Commands{
			RegisterNewChat: command.MustNewRegisterNewChatHandler(chatService, cfg.Logger),
			Enroll:          command.MustNewEnrollHandler(chatService, cfg.Logger),
			DisEnroll:       command.MustNewDisEnrollHandler(chatService, cfg.Logger),
			Magic:           command.MustNewMagicHandler(chatService, cfg.Logger),
		},
		Queries: app.Queries{
			GetMyReceiver: query.MustNewGetMyReceiverHandler(chatService, cfg.Logger),
			List:          query.MustNewListHandler(chatService, cfg.Logger),
		},
	}
}
