package app

import (
	"github.com/truewebber/secretsantabot/internal/app/command"
	"github.com/truewebber/secretsantabot/internal/app/query"
)

type (
	Application struct {
		Commands Commands
		Queries  Queries
	}

	Commands struct {
		Enroll    command.EnrollHandler
		DisEnroll command.DisEnrollHandler
		Magic     command.MagicHandler
	}

	Queries struct {
		GetMyPerson query.GetMyPersonHandler
		List        query.ListHandler
	}
)
