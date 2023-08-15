package query

import (
	"context"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
)

type ListParticipantsHandler struct {
	service storage.Storage
}

func NewListParticipantsHandler(service storage.Storage) (*ListParticipantsHandler, error) {
	if service == nil {
		return nil, errServiceIsNil
	}

	return &ListParticipantsHandler{service: service}, nil
}

func MustNewListParticipantsHandler(service storage.Storage) *ListParticipantsHandler {
	h, err := NewListParticipantsHandler(service)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *ListParticipantsHandler) Handle(ctx context.Context, appChat types.Chat) ([]types.Person, error) {
	if appChat.IsNotAGroup() {
		return nil, apperrors.ErrChatTypeIsUnsupported
	}

	domainChat := types.ChatToDomain(appChat)

	var domainPersons []chat.Person

	doErr := h.service.DoOperationOnTx(ctx, func(opCtx context.Context, tx storage.Tx) error {
		version, err := tx.GetLatestMagicVersion(ctx, domainChat)
		if err != nil {
			return fmt.Errorf("get latest magic version: %w", err)
		}

		domainPersons, err = tx.ListParticipants(ctx, version)
		if err != nil {
			return fmt.Errorf("list participants: %w", err)
		}

		return nil
	})

	if doErr != nil {
		return nil, fmt.Errorf("do operation on tx: %w", doErr)
	}

	return types.DomainsToPersons(domainPersons), nil
}
