package command

import (
	"context"
	"fmt"

	"github.com/truewebber/secretsantabot/app/types"
	chatdomain "github.com/truewebber/secretsantabot/domain/chat"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
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

func (h *EnrollHandler) Handle(chat *types.Chat, participant *types.Person) error {
	doErr := h.service.DoOperationOnTx(func(ctx context.Context, storageTx storage.Tx) error {
		chatToParticipate, err := storageTx.GetChatByTelegramID(ctx, chat.TelegramChatID)
		if err != nil {
			return fmt.Errorf("get chat by telegramID: %w", err)
		}

		version, err := storageTx.GetLatestMagicVersion(ctx, chatToParticipate)
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
		TelegramUserID: p.TelegramUserID,
	}
}
