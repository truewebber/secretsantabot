package command

import (
	"context"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
)

type RegisterMagicVersion struct {
	service storage.Storage
}

func NewRegisterMagicVersionHandler(service storage.Storage) (*RegisterMagicVersion, error) {
	if service == nil {
		return nil, errServiceIsNil
	}

	return &RegisterMagicVersion{service: service}, nil
}

func MustNewRegisterMagicVersionHandler(service storage.Storage) *RegisterMagicVersion {
	h, err := NewRegisterMagicVersionHandler(service)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *RegisterMagicVersion) Handle(ctx context.Context, appChat types.Chat) error {
	if appChat.IsNotAGroup() {
		return apperrors.ErrChatTypeIsUnsupported
	}

	chatToSave := types.ChatToDomain(appChat)

	chatVersionToSave := chat.MagicVersion{
		Chat: chatToSave,
	}

	doErr := h.service.DoOperationOnTx(ctx, func(opCtx context.Context, tx storage.Tx) error {
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
