package command

import (
	"context"
	"fmt"

	"github.com/truewebber/secretsantabot/internal/app/types"
	chatdomain "github.com/truewebber/secretsantabot/internal/chat"
	"github.com/truewebber/secretsantabot/internal/chat/storage"
	"github.com/truewebber/secretsantabot/internal/log"
)

type EnrollHandler struct {
	service storage.Storage
}

func NewEnrollHandler(service storage.Storage, logger log.Logger) (*EnrollHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &EnrollHandler{service: service}, nil
}

func MustNewEnrollHandler(service storage.Storage, logger log.Logger) *EnrollHandler {
	h, err := NewEnrollHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *EnrollHandler) Handle(participant *types.Person) error {
	doErr := h.service.DoOperationOnTx(func(ctx context.Context, storageTx storage.Tx) error {
		chat, err := storageTx.GetChatByTelegramID(ctx, participant.TelegramChatID)
		if err != nil {
			return fmt.Errorf("get chat by telegramID: %w", err)
		}

		version, err := storageTx.GetLatestMagicVersion(ctx, chat)
		if err != nil {
			return fmt.Errorf("get magic version by chat: %w", err)
		}

		participantToSave := castParticipantToDomain(participant)

		if err := storageTx.InsertParticipant(ctx, version, participantToSave); err != nil {
			return fmt.Errorf("insert new participant: %w", err)
		}

		return nil
	})

	if doErr != nil {
		return fmt.Errorf("do operation on tx: %w", doErr)
	}

	return nil
}

func castParticipantToDomain(p *types.Person) *chatdomain.Person {
	return &chatdomain.Person{
		ID:             0,
		TelegramUserID: p.TelegramUserID,
		TelegramChatID: p.TelegramChatID,
	}
}
