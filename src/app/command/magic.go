package command

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/magic"
)

type MagicHandler struct {
	service storage.Storage
}

func NewMagicHandler(service storage.Storage) (*MagicHandler, error) {
	if service == nil {
		return nil, errServiceIsNil
	}

	return &MagicHandler{service: service}, nil
}

func MustNewMagicHandler(service storage.Storage) *MagicHandler {
	h, err := NewMagicHandler(service)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *MagicHandler) Handle(ctx context.Context, appChat types.Chat, caller types.Person) error {
	if appChat.IsNotAGroup() {
		return apperrors.ErrChatTypeIsUnsupported
	}

	if appChat.Admin.TelegramUserID != caller.TelegramUserID {
		return apperrors.ErrForbidden
	}

	domainChat := types.ChatToDomain(appChat)

	if doErr := h.service.DoLockedOperationOnTx(
		ctx, appChat.Admin.TelegramUserID, calculateAndInsertMagicOnTx(domainChat),
	); doErr != nil {
		return fmt.Errorf("do locked operation on tx: %w", doErr)
	}

	return nil
}

func calculateAndInsertMagicOnTx(domainChat chat.Chat) func(ctx context.Context, tx storage.Tx) error {
	return func(ctx context.Context, tx storage.Tx) error {
		version, err := tx.GetLatestMagicVersion(ctx, domainChat)
		if err != nil {
			return fmt.Errorf("get latest magic version: %w", err)
		}

		_, magicErr := tx.GetMagic(ctx, version)

		if magicErr == nil {
			return apperrors.ErrAlreadyExists
		}

		if !errors.Is(magicErr, storage.ErrNotFound) {
			return fmt.Errorf("get magic: %w", magicErr)
		}

		participants, err := tx.ListParticipants(ctx, version)
		if err != nil {
			return fmt.Errorf("list participants: %w", err)
		}

		calculatedMagic, err := magic.Calculate(participants)

		if errors.Is(err, magic.ErrNotEnoughParticipants) {
			return apperrors.ErrNotEnoughParticipants
		}

		if err != nil {
			return fmt.Errorf("calculate magic: %w", err)
		}

		if err := tx.InsertMagic(ctx, version, calculatedMagic); err != nil {
			return fmt.Errorf("insert magic: %w", err)
		}

		return nil
	}
}
