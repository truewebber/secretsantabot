package query

import (
	"context"

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

	return giver, nil
}
