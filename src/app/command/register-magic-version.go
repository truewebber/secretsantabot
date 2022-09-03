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

type RegisterMagicVersion struct {
	service storage.Storage
}

func NewRegisterMagicVersionHandler(
	service storage.Storage,
	logger log.Logger,
) (*RegisterMagicVersion, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &RegisterMagicVersion{service: service}, nil
}

func MustNewRegisterMagicVersionHandler(
	service storage.Storage,
	logger log.Logger,
) *RegisterMagicVersion {
	h, err := NewRegisterMagicVersionHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *RegisterMagicVersion) Handle(appChat *types.Chat) error {
	if !appChat.IsGroup {
		return apperrors.ErrRegisterLocalChatIsRestricted
	}

	chatToSave := types.ChatToDomain(appChat)

	chatVersionToSave := &chat.MagicVersion{
		Chat: chatToSave,
	}

	doErr := h.service.DoOperationOnTx(func(ctx context.Context, tx storage.Tx) error {
		if err := tx.InsertNewMagicVersion(ctx, chatVersionToSave); err != nil {
			return fmt.Errorf("insert new chat magic version: %w", err)
		}

		return nil
	})
	if doErr != nil {
		return fmt.Errorf("do operation on tx: %w", doErr)
	}

	return nil
}
