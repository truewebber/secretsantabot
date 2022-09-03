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

func MustRegisterNewChatAndVersionHandler(
	service storage.Storage,
	logger log.Logger,
) *RegisterNewChatAndVersionHandler {
	h, err := NewRegisterNewChatAndVersionHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *RegisterNewChatAndVersionHandler) Handle(appChat *types.Chat) error {
	if !appChat.IsGroup {
		return apperrors.ErrRegisterLocalChatIsRestricted
	}

	chatToSave := types.ChatToDomain(appChat)

	doErr := h.service.DoOperationOnTx(func(ctx context.Context, tx storage.Tx) error {
		if err := tx.InsertPerson(ctx, chatToSave.Admin); err != nil {
			return fmt.Errorf("insert person: %w", err)
		}

		if err := tx.InsertChat(ctx, chatToSave); err != nil {
			return fmt.Errorf("insert chat: %w", err)
		}

		chatVersionToSave := chat.MagicVersion{
			Chat:    chatToSave,
			ID:      0,
			Version: 1,
			Status:  0,
		}

		if err := tx.InsertNewMagicVersion(ctx, chatToSave); err != nil {
			return fmt.Errorf("insert chat: %w", err)
		}

		return nil
	})
	if doErr != nil {
		return fmt.Errorf("do operation on tx: %w", doErr)
	}

	return nil
}

func firstChatVersionFromChat(c chat.Chat) chat.MagicVersion {

}
