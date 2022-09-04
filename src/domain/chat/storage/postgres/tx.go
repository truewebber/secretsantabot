package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	chatdomain "github.com/truewebber/secretsantabot/domain/chat"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
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

	selectDeletedChatQuery = "SELECT admin_user_id, deleted FROM chats WHERE id=$1;"
)

func (s *StorageTx) InsertChat(ctx context.Context, chat *chatdomain.Chat) error {
	var deleted bool

	selectErr := s.tx.QueryRow(ctx, selectDeletedChatQuery, chat.TelegramChatID).
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

const selectChatQuery = "SELECT admin_user_id FROM chats WHERE id=$1 AND deleted=false;"

func (s *StorageTx) GetChatByTelegramID(ctx context.Context, id int64) (*chatdomain.Chat, error) {
	var adminUserID int64

	err := s.tx.QueryRow(ctx, selectChatQuery, id).Scan(&adminUserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("query row select chat: %w", err)
	}

	return &chatdomain.Chat{
		Admin: &chatdomain.Person{
			TelegramUserID: adminUserID,
		},
		TelegramChatID: id,
	}, nil
}

const insertMagicVersionQuery = `INSERT INTO magic_chat_history (chat_id, deleted)
VALUES ($1, false) RETURNING id;`

func (s *StorageTx) InsertNewMagicVersion(ctx context.Context, version *chatdomain.MagicVersion) error {
	if err := s.tx.QueryRow(ctx, insertMagicVersionQuery, version.Chat.TelegramChatID).Scan(&version.ID); err != nil {
		return fmt.Errorf("exec insert magic version: %w", err)
	}

	return nil
}

const selectLatestMagicVersionQuery = `SELECT id FROM magic_chat_history WHERE chat_id=$1 AND deleted=false
ORDER BY id DESC LIMIT 1;`

func (s *StorageTx) GetLatestMagicVersion(
	ctx context.Context,
	chat *chatdomain.Chat,
) (*chatdomain.MagicVersion, error) {
	version := &chatdomain.MagicVersion{
		Chat: chat,
	}

	err := s.tx.QueryRow(ctx, selectLatestMagicVersionQuery, chat.TelegramChatID).Scan(&version.ID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("exec select latest magic version: %w", err)
	}

	return version, nil
}

func (s *StorageTx) ListParticipants(ctx context.Context, v *chatdomain.MagicVersion) ([]chatdomain.Person, error) {
	return []chatdomain.Person{}, nil
}

func (s *StorageTx) InsertParticipant(
	ctx context.Context,
	version *chatdomain.MagicVersion,
	person *chatdomain.Person,
) error {
	return nil
}

func (s *StorageTx) DeleteParticipant(
	ctx context.Context,
	version *chatdomain.MagicVersion,
	person *chatdomain.Person,
) error {
	return nil
}

func (s *StorageTx) InsertMagic(ctx context.Context, v *chatdomain.MagicVersion, magic chatdomain.Magic) error {
	return nil
}

func (s *StorageTx) GetMagic(ctx context.Context, v *chatdomain.MagicVersion) (chatdomain.Magic, error) {
	return chatdomain.Magic{
		Data: nil,
	}, nil
}
