package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	chatdomain "github.com/truewebber/secretsantabot/domain/chat"
)

type StorageTx struct {
	tx pgx.Tx
}

func newStorageTx(tx pgx.Tx) *StorageTx {
	return &StorageTx{tx: tx}
}

const (
	insertChatQuery = "INSERT INTO chats (id, admin_user_id, deleted) VALUES ($1, $2, false) ON CONFLICT DO NOTHING;"

	restoreChatQuery = "UPDATE chats SET deleted=false WHERE id=$1;"

	selectChatQuery = "SELECT admin_user_id, deleted FROM chats WHERE id=$1;"
)

func (s StorageTx) InsertChat(ctx context.Context, chat *chatdomain.Chat) error {
	var deleted bool

	selectErr := s.tx.QueryRow(ctx, selectChatQuery, chat.TelegramChatID).
		Scan(&chat.Admin.TelegramUserID, &deleted)

	if errors.Is(selectErr, pgx.ErrNoRows) {
		_, insertErr := s.tx.Exec(ctx, insertChatQuery, chat.TelegramChatID, chat.Admin.TelegramUserID)
		if insertErr != nil {
			return fmt.Errorf("exec insert chat: %w", insertErr)
		}

		return nil
	}

	if selectErr != nil {
		return fmt.Errorf("query row select chat: %w", selectErr)
	}

	if deleted {
		if _, err := s.tx.Exec(ctx, restoreChatQuery, chat.TelegramChatID); err != nil {
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

const (
	insertPersonQuery = "INSERT INTO users (id, deleted) VALUES ($1, false) ON CONFLICT DO NOTHING;"

	restorePersonQuery = "UPDATE users SET deleted=false WHERE id=$1;"

	selectPersonQuery = "SELECT deleted FROM users WHERE id=$1;"
)

func (s StorageTx) InsertPerson(ctx context.Context, person *chatdomain.Person) error {
	var deleted bool

	selectErr := s.tx.QueryRow(ctx, selectPersonQuery, person.TelegramUserID).
		Scan(&deleted)

	if errors.Is(selectErr, pgx.ErrNoRows) {
		_, insertErr := s.tx.Exec(ctx, insertPersonQuery, person.TelegramUserID)
		if insertErr != nil {
			return fmt.Errorf("exec insert person: %w", insertErr)
		}

		return nil
	}

	if selectErr != nil {
		return fmt.Errorf("query row select person: %w", selectErr)
	}

	if deleted {
		if _, err := s.tx.Exec(ctx, restorePersonQuery, person.TelegramUserID); err != nil {
			return fmt.Errorf("exec restore person: %w", err)
		}
	}

	return nil
}

func (s StorageTx) GetPersonByTelegramID(ctx context.Context, id int64) (*chatdomain.Person, error) {
	return &chatdomain.Person{
		TelegramUserID: id,
	}, nil
}

func (s StorageTx) InsertNewMagicVersion(ctx context.Context, v *chatdomain.MagicVersion) error {
	return nil
}

func (s StorageTx) GetLatestMagicVersion(ctx context.Context, chat *chatdomain.Chat) (*chatdomain.MagicVersion, error) {
	return &chatdomain.MagicVersion{Chat: chat}, nil
}

func (s StorageTx) ListParticipants(ctx context.Context, v *chatdomain.MagicVersion) ([]chatdomain.Person, error) {
	return []chatdomain.Person{}, nil
}

func (s StorageTx) InsertParticipant(ctx context.Context, v *chatdomain.MagicVersion, person *chatdomain.Person) error {
	return nil
}

func (s StorageTx) DeleteParticipant(ctx context.Context, v *chatdomain.MagicVersion, person *chatdomain.Person) error {
	return nil
}

func (s StorageTx) InsertMagic(ctx context.Context, v *chatdomain.MagicVersion, magic chatdomain.Magic) error {
	return nil
}

func (s StorageTx) GetMagic(ctx context.Context, v *chatdomain.MagicVersion) (chatdomain.Magic, error) {
	return chatdomain.Magic{
		Data: nil,
	}, nil
}
