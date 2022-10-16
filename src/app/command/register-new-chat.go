package command

import (
	"context"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
)

type RegisterNewChatAndVersionHandler struct {
	service storage.Storage
}

func NewRegisterNewChatAndVersionHandler(
	service storage.Storage,
	logger log.Logger,
) (*RegisterNewChatAndVersionHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &RegisterNewChatAndVersionHandler{service: service}, nil
}

func MustNewRegisterNewChatAndVersionHandler(
	service storage.Storage,
	logger log.Logger,
) *RegisterNewChatAndVersionHandler {
	h, err := NewRegisterNewChatAndVersionHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *RegisterNewChatAndVersionHandler) Handle(ctx context.Context, appChat *types.Chat) error {
	if appChat.IsPrivate() {
		return apperrors.ErrChatIsPrivate
	}

	if appChat.IsNotAGroup() {
		return apperrors.ErrChatTypeIsUnsupported
	}

	chatToSave := types.ChatToDomain(appChat)

	doErr := h.service.DoOperationOnTx(ctx, func(opCtx context.Context, tx storage.Tx) error {
		if err := tx.InsertChat(opCtx, chatToSave); err != nil {
			return fmt.Errorf("insert chat: %w", err)
		}

		chatVersionToSave := &chat.MagicVersion{
			Chat: chatToSave,
		}

		if err := tx.InsertNewMagicVersion(opCtx, chatVersionToSave); err != nil {
			return fmt.Errorf("insert new chat magic version: %w", err)
		}

		return nil
	})
	if doErr != nil {
		return fmt.Errorf("do operation on tx: %w", doErr)
	}

	return nil
}
