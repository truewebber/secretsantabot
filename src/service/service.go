package service

import (
	"github.com/truewebber/secretsantabot/app"
	"github.com/truewebber/secretsantabot/app/command"
	"github.com/truewebber/secretsantabot/app/query"
	"github.com/truewebber/secretsantabot/domain/chat/storage/postgres"
)

type (
	Config struct {
		ChatService ChatService
	}

	ChatService struct {
		PostgresURI string
	}
)

func NewConfig(chatService ChatService) *Config {
	return &Config{
		ChatService: chatService,
	}
}

func MustNewApplication(cfg *Config) *app.Application {
	chatStorage := postgres.MustNewPGX(cfg.ChatService.PostgresURI)

	return &app.Application{
		Commands: app.Commands{
			RegisterNewChatAndVersion: command.MustNewRegisterNewChatAndVersionHandler(chatStorage),
			RegisterMagicVersion:      command.MustNewRegisterMagicVersionHandler(chatStorage),
			Enroll:                    command.MustNewEnrollHandler(chatStorage),
			DisEnroll:                 command.MustNewDisEnrollHandler(chatStorage),
			Magic:                     command.MustNewMagicHandler(chatStorage),
		},
		Queries: app.Queries{
			GetMyReceiver:    query.MustNewGetMyReceiverHandler(chatStorage),
			ListParticipants: query.MustNewListParticipantsHandler(chatStorage),
			GetMagic:         query.MustNewGetMagicHandler(chatStorage),
		},
	}
}
