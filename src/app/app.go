package app

import (
	"github.com/truewebber/secretsantabot/app/command"
	"github.com/truewebber/secretsantabot/app/query"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		RegisterNewChatAndVersion *command.RegisterNewChatAndVersionHandler
		RegisterMagicVersion      *command.RegisterMagicVersion
		Enroll                    *command.EnrollHandler
		DisEnroll                 *command.DisEnrollHandler
		Magic                     *command.MagicHandler
	}

	Queries struct {
		GetMagic         *query.GetMagicHandler
		GetMyReceiver    *query.GetMyReceiverHandler
		ListParticipants *query.ListParticipantsHandler
	}
)
