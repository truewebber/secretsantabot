package storage

import (
	"context"
	"errors"

	"github.com/truewebber/secretsantabot/domain/chat"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Storage interface {
	DoOperationOnTx(context.Context, func(context.Context, Tx) error) error
	DoLockedOperationOnTx(ctx context.Context, lockID int64, operation func(context.Context, Tx) error) error
}

type Tx interface {
	LockTx(ctx context.Context, lockID int64) error

	ChatTx
	ParticipantTx
	MagicTx
}

type ChatTx interface {
	InsertChat(context.Context, chat.Chat) error
	GetChatByTelegramID(context.Context, int64) (chat.Chat, error)

	InsertNewMagicVersion(context.Context, chat.MagicVersion) error
	GetLatestMagicVersion(context.Context, chat.Chat) (chat.MagicVersion, error)
}

type ParticipantTx interface {
	InsertParticipant(context.Context, chat.MagicVersion, chat.Person) error
	DeleteParticipant(context.Context, chat.MagicVersion, chat.Person) error
	ListParticipants(context.Context, chat.MagicVersion) ([]chat.Person, error)
}

type MagicTx interface {
	InsertMagic(context.Context, chat.MagicVersion, chat.Magic) error
	GetMagic(context.Context, chat.MagicVersion) (chat.Magic, error)
	GetMagicRecipient(context.Context, chat.MagicVersion, chat.Person) (chat.Person, error)
}
