package query

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
)

type GetMyReceiverHandler struct {
	service storage.Storage
}

func NewGetMyReceiverHandler(service storage.Storage, logger log.Logger) (*GetMyReceiverHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &GetMyReceiverHandler{service: service}, nil
}

func MustNewGetMyReceiverHandler(service storage.Storage, logger log.Logger) *GetMyReceiverHandler {
	h, err := NewGetMyReceiverHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *GetMyReceiverHandler) Handle(
	ctx context.Context,
	appChat types.Chat,
	giver types.Person,
) (types.Person, error) {
	if appChat.IsUnsupported() {
		return types.Person{}, apperrors.ErrChatTypeIsUnsupported
	}

	domainChat := types.ChatToDomain(appChat)
	domainGiver := types.PersonToDomain(giver)

	var receiver types.Person

	if doErr := h.service.DoLockedOperationOnTx(
		ctx, appChat.Admin.TelegramUserID, func(ctx context.Context, tx storage.Tx) error {
			version, err := tx.GetLatestMagicVersion(ctx, domainChat)
			if err != nil {
				return fmt.Errorf("get latest magic version: %w", err)
			}

			domainReceiver, magicErr := tx.GetMagicRecipient(ctx, version, domainGiver)

			if errors.Is(magicErr, storage.ErrNotFound) {
				return apperrors.ErrNotFound
			}

			if magicErr != nil {
				return fmt.Errorf("get magic: %w", magicErr)
			}

			receiver = types.DomainToPerson(domainReceiver)

			return nil
		},
	); doErr != nil {
		return types.Person{}, fmt.Errorf("do locked operation on tx: %w", doErr)
	}

	return receiver, nil
}
