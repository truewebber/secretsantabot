package command

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/internal/app/errors"
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/chat/storage"
	"github.com/truewebber/secretsantabot/internal/log"
)

type RegisterNewChatHandler struct {
	service storage.Storage
}

func NewRegisterNewChatHandler(service storage.Storage, logger log.Logger) (*RegisterNewChatHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &RegisterNewChatHandler{service: service}, nil
}

func MustNewRegisterNewChatHandler(service storage.Storage, logger log.Logger) *RegisterNewChatHandler {
	h, err := NewRegisterNewChatHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *RegisterNewChatHandler) Handle(appChat *types.Chat) error {
	if !appChat.IsGroup {
		return apperrors.ErrRegisterLocalChatIsRestricted
	}

	chatToSave := types.ChatToDomain(appChat)

	doErr := h.service.DoOperationOnTx(func(ctx context.Context, tx storage.Tx) error {
		if err := tx.InsertChat(ctx, chatToSave); err != nil {
			return fmt.Errorf("insert chat: %w", err)
		}

		if err := tx.InsertPerson(ctx, chatToSave.Admin); err != nil {
			return fmt.Errorf("insert person: %w", err)
		}

		version, err := tx.GetLatestMagicVersion(ctx, chatToSave)

		if errors.Is(err, storage.ErrNotFound) {

		}

		if err != nil {
			return fmt.Errorf("get latest magic version: %w", err)
		}

		return nil
	})
	if doErr != nil {
		return fmt.Errorf("do operation on tx: %w", doErr)
	}

	return nil
}
