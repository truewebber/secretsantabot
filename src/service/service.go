package service

import (
	"github.com/truewebber/secretsantabot/app"
	"github.com/truewebber/secretsantabot/app/command"
	"github.com/truewebber/secretsantabot/app/query"
	"github.com/truewebber/secretsantabot/domain/chat/storage/postgres"
	"github.com/truewebber/secretsantabot/domain/log"
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

func NewConfig(chatService ChatService, logger log.Logger) *Config {
	return &Config{
		ChatService: chatService,
		Logger:      logger,
	}
}

func MustNewApplication(cfg *Config) *app.Application {
	chatService := postgres.MustNewPostgres(cfg.ChatService.PostgresURI)

	return &app.Application{
		Commands: app.Commands{
			RegisterNewChatAndVersion: command.MustNewRegisterNewChatAndVersionHandler(chatService, cfg.Logger),
			RegisterMagicVersion:      command.MustNewRegisterMagicVersionHandler(chatService, cfg.Logger),
			Enroll:                    command.MustNewEnrollHandler(chatService, cfg.Logger),
			DisEnroll:                 command.MustNewDisEnrollHandler(chatService, cfg.Logger),
			Magic:                     command.MustNewMagicHandler(chatService, cfg.Logger),
		},
		Queries: app.Queries{
			GetMyReceiver: query.MustNewGetMyReceiverHandler(chatService, cfg.Logger),
			List:          query.MustNewListHandler(chatService, cfg.Logger),
		},
	}
}
