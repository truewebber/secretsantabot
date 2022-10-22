package command

import (
	"context"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
)

type DisEnrollHandler struct {
	service storage.Storage
}

func NewDisEnrollHandler(service storage.Storage, logger log.Logger) (*DisEnrollHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &DisEnrollHandler{service: service}, nil
}

func MustNewDisEnrollHandler(service storage.Storage, logger log.Logger) *DisEnrollHandler {
	h, err := NewDisEnrollHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *DisEnrollHandler) Handle(ctx context.Context, appChat types.Chat, participant types.Person) error {
	if appChat.IsNotAGroup() {
		return apperrors.ErrChatTypeIsUnsupported
	}

	domainChat := types.ChatToDomain(appChat)
	domainParticipant := types.PersonToDomain(participant)

	doErr := h.service.DoOperationOnTx(ctx, func(opCtx context.Context, tx storage.Tx) error {
		version, err := tx.GetLatestMagicVersion(opCtx, domainChat)
		if err != nil {
			return fmt.Errorf("get latest magic version: %w", err)
		}

		if deleteErr := tx.DeleteParticipant(opCtx, version, domainParticipant); deleteErr != nil {
			return fmt.Errorf("delete participant: %w", deleteErr)
		}

		return nil
	})

	if doErr != nil {
		return fmt.Errorf("do operation on tx: %w", doErr)
	}

	return nil
}
