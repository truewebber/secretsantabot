package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	chatdomain "github.com/truewebber/secretsantabot/domain/chat"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
)

type pgxTx struct {
	tx pgx.Tx
}

func newStorageTx(tx pgx.Tx) *pgxTx {
	return &pgxTx{tx: tx}
}

const (
	insertChatQuery = "INSERT INTO chats (id, admin_user_id, deleted) VALUES ($1, $2, false) ON CONFLICT DO NOTHING;"

	restoreChatQuery = "UPDATE chats SET deleted=false WHERE id=$1;"

	selectDeletedChatQuery = "SELECT admin_user_id, deleted FROM chats WHERE id=$1;"
)

func (s *pgxTx) InsertChat(ctx context.Context, chat *chatdomain.Chat) error {
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

const selectChatByIDQuery = "SELECT admin_user_id FROM chats WHERE id=$1 AND deleted=false;"

func (s *pgxTx) GetChatByTelegramID(ctx context.Context, id int64) (*chatdomain.Chat, error) {
	var adminUserID int64

	err := s.tx.QueryRow(ctx, selectChatByIDQuery, id).Scan(&adminUserID)

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

func (s *pgxTx) InsertNewMagicVersion(ctx context.Context, version *chatdomain.MagicVersion) error {
	if err := s.tx.QueryRow(ctx, insertMagicVersionQuery, version.Chat.TelegramChatID).Scan(&version.ID); err != nil {
		return fmt.Errorf("exec insert magic version: %w", err)
	}

	return nil
}

const selectLatestMagicVersionQuery = `SELECT id FROM magic_chat_history WHERE chat_id=$1 AND deleted=false
ORDER BY id DESC LIMIT 1;`

func (s *pgxTx) GetLatestMagicVersion(
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

const selectEnrolledParticipantsQuery = `SELECT participant_user_id FROM magic_participants
WHERE magic_chat_history_id=$1 AND deleted IS FALSE;`

func (s *pgxTx) ListParticipants(ctx context.Context, v *chatdomain.MagicVersion) ([]chatdomain.Person, error) {
	rows, err := s.tx.Query(ctx, selectEnrolledParticipantsQuery, v.ID)
	if err != nil {
		return nil, fmt.Errorf("query enrolled participants: %w", err)
	}

	defer rows.Close()

	var participants []chatdomain.Person

	for rows.Next() {
		var p chatdomain.Person

		if err := rows.Scan(&p.TelegramUserID); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		participants = append(participants, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return participants, nil
}

const (
	insertParticipantQuery = `INSERT INTO magic_participants (magic_chat_history_id, participant_user_id, deleted)
VALUES ($1, $2, false) ON CONFLICT DO NOTHING;`

	restoreParticipantQuery = "UPDATE magic_participants SET deleted=false WHERE id=$1;"

	selectDeletedParticipantQuery = `SELECT id, deleted FROM magic_participants
WHERE magic_chat_history_id=$1 AND participant_user_id=$2;`
)

func (s *pgxTx) InsertParticipant(
	ctx context.Context,
	version *chatdomain.MagicVersion,
	person *chatdomain.Person,
) error {
	var (
		id      uint32
		deleted bool
	)

	selectErr := s.tx.QueryRow(ctx, selectDeletedParticipantQuery, version.ID, person.TelegramUserID).Scan(&id, &deleted)

	if errors.Is(selectErr, pgx.ErrNoRows) {
		_, insertErr := s.tx.Exec(ctx, insertParticipantQuery, version.ID, person.TelegramUserID)
		if insertErr != nil {
			return fmt.Errorf("exec insert participant: %w", insertErr)
		}

		return nil
	}

	if selectErr != nil {
		return fmt.Errorf("query row select participant: %w", selectErr)
	}

	if deleted {
		if _, err := s.tx.Exec(ctx, restoreParticipantQuery, id); err != nil {
			return fmt.Errorf("exec restore participant: %w", err)
		}
	}

	return nil
}

var errInvalidAmountRowsAffected = errors.New("invalid amount rows affected")

const deleteParticipantQuery = `UPDATE magic_participants SET deleted=true
WHERE magic_chat_history_id=$1 AND participant_user_id=$2;`

func (s *pgxTx) DeleteParticipant(
	ctx context.Context,
	version *chatdomain.MagicVersion,
	person *chatdomain.Person,
) error {
	tag, err := s.tx.Exec(ctx, deleteParticipantQuery, version.ID, person.TelegramUserID)
	if err != nil {
		return fmt.Errorf("exec delete participant: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%w on delete: %v", errInvalidAmountRowsAffected, tag.RowsAffected())
	}

	return nil
}

func (s *pgxTx) InsertMagic(ctx context.Context, v *chatdomain.MagicVersion, magic chatdomain.Magic) error {
	for i := range magic.Data {
		if err := insertMagicPairOnTx(ctx, s.tx, v, magic.Data[i]); err != nil {
			return fmt.Errorf("insert magic pair on tx: %w", err)
		}
	}

	return nil
}

const insertMagicPairQuery = `INSERT INTO magic_results 
(magic_chat_history_id, participant_giver_id, participant_receiver_id, deleted) 
VALUES ($1, $2, $3, false);`

func insertMagicPairOnTx(
	ctx context.Context,
	tx pgx.Tx,
	v *chatdomain.MagicVersion,
	pair chatdomain.GiverReceiverPair,
) error {
	tag, err := tx.Exec(ctx, insertMagicPairQuery, v.ID, pair.Giver.TelegramUserID, pair.Receiver.TelegramUserID)
	if err != nil {
		return fmt.Errorf("insert magic pair: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%w on insert: %v", errInvalidAmountRowsAffected, tag.RowsAffected())
	}

	return nil
}

const selectMagicPairQuery = `SELECT participant_giver_id, participant_receiver_id FROM magic_results
WHERE magic_chat_history_id=$1 AND deleted IS FALSE;`

func (s *pgxTx) GetMagic(ctx context.Context, v *chatdomain.MagicVersion) (chatdomain.Magic, error) {
	rows, err := s.tx.Query(ctx, selectMagicPairQuery, v.ID)
	if err != nil {
		return chatdomain.Magic{}, fmt.Errorf("query select magic pair: %w", err)
	}

	defer rows.Close()

	var pairs []chatdomain.GiverReceiverPair

	for rows.Next() {
		var pair chatdomain.GiverReceiverPair

		if err := rows.Scan(&pair.Giver.TelegramUserID, &pair.Receiver.TelegramUserID); err != nil {
			return chatdomain.Magic{}, fmt.Errorf("scan row: %w", err)
		}

		pairs = append(pairs, pair)
	}

	if err := rows.Err(); err != nil {
		return chatdomain.Magic{}, fmt.Errorf("rows: %w", err)
	}

	if len(pairs) == 0 {
		return chatdomain.Magic{}, storage.ErrNotFound
	}

	return chatdomain.Magic{
		Data: nil,
	}, nil
}

const lockQuery = "SELECT pg_advisory_xact_lock($1);"

func (s *pgxTx) LockTx(ctx context.Context, lockID int64) error {
	if _, err := s.tx.Exec(ctx, lockQuery, lockID); err != nil {
		return fmt.Errorf("exec lock query: %w", err)
	}

	return nil
}
