package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	chatdomain "github.com/truewebber/secretsantabot/internal/chat"
)

type StorageTx struct {
	tx pgx.Tx
}

func newStorageTx(tx pgx.Tx) *StorageTx {
	return &StorageTx{tx: tx}
}

const (
	insertChatQuery = "INSERT INTO chats (tg_chat_id, tg_admin_id, deleted) VALUES ($1, $2, false) returning id;"

	restoreChatQuery = "UPDATE chats SET deleted=false WHERE id=$1;"

	selectChatQuery = "SELECT id, deleted FROM chats WHERE tg_chat_id=$1 AND tg_admin_id=$2;"
)

func (s StorageTx) InsertChat(ctx context.Context, chat *chatdomain.Chat) error {
	var deleted bool

	selectErr := s.tx.QueryRow(ctx, selectChatQuery, chat.TelegramChatID, chat.Admin.TelegramUserID).
		Scan(&chat.ID, &deleted)

	if errors.Is(selectErr, pgx.ErrNoRows) {
		insertErr := s.tx.QueryRow(ctx, insertChatQuery, chat.TelegramChatID, chat.Admin.TelegramUserID).Scan(&chat.ID)
		if insertErr != nil {
			return fmt.Errorf("query row insert chat: %w", insertErr)
		}

		return nil
	}

	if selectErr != nil {
		return fmt.Errorf("query row select chat: %w", selectErr)
	}

	if deleted {
		if _, err := s.tx.Exec(ctx, restoreChatQuery, chat.ID); err != nil {
			return fmt.Errorf("exec restore chat: %w", err)
		}
	}

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
