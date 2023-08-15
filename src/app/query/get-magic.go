package query

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
)

type GetMagicHandler struct {
	service storage.Storage
}

func NewGetMagicHandler(service storage.Storage) (*GetMagicHandler, error) {
	if service == nil {
		return nil, errServiceIsNil
	}

	return &GetMagicHandler{service: service}, nil
}

func MustNewGetMagicHandler(service storage.Storage) *GetMagicHandler {
	h, err := NewGetMagicHandler(service)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *GetMagicHandler) Handle(ctx context.Context, appChat types.Chat, caller types.Person) (types.Magic, error) {
	if appChat.IsNotAGroup() {
		return types.Magic{}, apperrors.ErrChatTypeIsUnsupported
	}

	if appChat.Admin.TelegramUserID != caller.TelegramUserID {
		return types.Magic{}, apperrors.ErrForbidden
	}

	domainChat := types.ChatToDomain(appChat)

	var magic types.Magic

	if doErr := h.service.DoLockedOperationOnTx(
		ctx, appChat.Admin.TelegramUserID, func(ctx context.Context, tx storage.Tx) error {
			version, err := tx.GetLatestMagicVersion(ctx, domainChat)
			if err != nil {
				return fmt.Errorf("get latest magic version: %w", err)
			}

			domainMagic, magicErr := tx.GetMagic(ctx, version)

			if errors.Is(magicErr, storage.ErrNotFound) {
				return apperrors.ErrNotFound
			}

			if magicErr != nil {
				return fmt.Errorf("get magic: %w", magicErr)
			}

			magic = types.DomainToMagic(domainMagic)

			return nil
		},
	); doErr != nil {
		return types.Magic{}, fmt.Errorf("do locked operation on tx: %w", doErr)
	}

	return magic, nil
}
