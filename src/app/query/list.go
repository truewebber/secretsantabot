package query

import (
	"context"
	"fmt"
	"github.com/truewebber/secretsantabot/domain/chat"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
)

type ListHandler struct {
	service storage.Storage
}

func NewListHandler(service storage.Storage, logger log.Logger) (*ListHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &ListHandler{service: service}, nil
}

func MustNewListHandler(service storage.Storage, logger log.Logger) *ListHandler {
	h, err := NewListHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *ListHandler) Handle(ctx context.Context, appChat *types.Chat) ([]types.Person, error) {
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
