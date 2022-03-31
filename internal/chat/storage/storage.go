package storage

import (
	"context"
	"errors"

	"github.com/truewebber/secretsantabot/internal/chat"
)

var ErrNotFound = errors.New("not found")

type (
	Storage interface {
		DoOperationOnTx(func(context.Context, Tx) error) error
	}

	Tx interface {
		InsertChat(context.Context, *chat.Chat) error
		GetChatByTelegramID(context.Context, int64) (*chat.Chat, error)

		InsertPerson(context.Context, *chat.Person) error
		GetPersonByTelegramID(context.Context, int64) (*chat.Person, error)

		InsertNewMagicVersion(context.Context, *chat.MagicVersion) error
		GetLatestMagicVersion(context.Context, *chat.Chat) (*chat.MagicVersion, error)

		InsertParticipant(context.Context, *chat.MagicVersion, *chat.Person) error
		DeleteParticipant(context.Context, *chat.MagicVersion, *chat.Person) error
		ListParticipants(context.Context, *chat.MagicVersion) ([]chat.Person, error)

		InsertMagic(context.Context, *chat.MagicVersion, chat.Magic) error
		GetMagic(context.Context, *chat.MagicVersion) (chat.Magic, error)
	}
)
