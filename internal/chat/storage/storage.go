package storage

import (
	"context"

	"github.com/truewebber/secretsantabot/internal/chat"
)

type (
	Storage interface {
		DoOperationOnTx(func(context.Context, Tx) error) error
	}

	Tx interface {
		InsertChat(context.Context, *chat.Chat) error
		UpdateChat(context.Context, *chat.Chat) error
		GetChatByTelegramID(context.Context, int64) (*chat.Chat, error)

		InsertPerson(context.Context, *chat.Person) error
		GetPersonByTelegramID(context.Context, int64) (*chat.Person, error)

		InsertMagic(context.Context, *chat.Chat, chat.Magic) error
		GetMagic(context.Context, *chat.Chat) (chat.Magic, error)

		ListParticipants(context.Context, *chat.Chat) ([]chat.Person, error)

		InsertNewParticipant(context.Context, *chat.Chat, *chat.Person) error
		DeleteParticipant(context.Context, *chat.Chat, *chat.Person) error
	}
)
