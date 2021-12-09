package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"

	chatdomain "github.com/truewebber/secretsantabot/internal/chat"
)

type StorageTx struct {
	tx pgx.Tx
}

func newStorageTx(tx pgx.Tx) *StorageTx {
	return &StorageTx{tx: tx}
}

func (s StorageTx) InsertChat(ctx context.Context, chat *chatdomain.Chat) error {
	return nil
}

func (s StorageTx) UpdateChat(ctx context.Context, chat *chatdomain.Chat) error {
	return nil
}

func (s StorageTx) GetChatByTelegramID(ctx context.Context, id int64) (*chatdomain.Chat, error) {
	return &chatdomain.Chat{
		TelegramChatID: id,
	}, nil
}

func (s StorageTx) InsertPerson(ctx context.Context, person *chatdomain.Person) error {
	return nil
}

func (s StorageTx) GetPersonByTelegramID(ctx context.Context, id int64) (*chatdomain.Person, error) {
	return &chatdomain.Person{
		TelegramUserID: id,
	}, nil
}

func (s StorageTx) InsertMagic(ctx context.Context, chat *chatdomain.Chat, magic chatdomain.Magic) error {
	return nil
}

func (s StorageTx) GetMagic(ctx context.Context, chat *chatdomain.Chat) (chatdomain.Magic, error) {
	return chatdomain.Magic{
		Data:   nil,
		Status: chatdomain.OpenMagicStatus,
	}, nil
}

func (s StorageTx) ListParticipants(ctx context.Context, chat *chatdomain.Chat) ([]chatdomain.Person, error) {
	return []chatdomain.Person{}, nil
}

func (s StorageTx) InsertNewParticipant(ctx context.Context, chat *chatdomain.Chat, person *chatdomain.Person) error {
	return nil
}

func (s StorageTx) DeleteParticipant(ctx context.Context, chat *chatdomain.Chat, person *chatdomain.Person) error {
	return nil
}
