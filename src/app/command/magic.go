package command

import (
	"context"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
)

type MagicHandler struct {
	service storage.Storage
}

func NewMagicHandler(service storage.Storage, logger log.Logger) (*MagicHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &MagicHandler{service: service}, nil
}

func MustNewMagicHandler(service storage.Storage, logger log.Logger) *MagicHandler {
	h, err := NewMagicHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *MagicHandler) Handle(
	ctx context.Context,
	appChat *types.Chat,
	caller *types.Person,
	notifyPersonFn func(p *types.Person) error,
) error {
	if appChat.IsNotAGroup() {
		return apperrors.ErrChatTypeIsUnsupported
	}

	if appChat.Admin.TelegramUserID != caller.TelegramUserID {
		return apperrors.ErrForbidden
	}

	doErr := h.service.DoLockedOperationOnTx(ctx, appChat.Admin.TelegramUserID,
		func(ctx context.Context, tx storage.Tx) error {
			return nil
		},
	)

	if doErr != nil {
		return fmt.Errorf("do locked operation on tx: %w", doErr)
	}

	return nil
}
