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
		RegisterNewChat *command.RegisterNewChatHandler
		Enroll          *command.EnrollHandler
		DisEnroll       *command.DisEnrollHandler
		Magic           *command.MagicHandler
	}

	Queries struct {
		GetMyReceiver *query.GetMyReceiverHandler
		List          *query.ListHandler
	}
)
