package query

import (
	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
)

type ListHandler struct {
	service storage.Storage
}

func NewListHandler(service storage.Storage, logger log.Logger) (*ListHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &ListHandler{service: service}, nil
}

func MustNewListHandler(service storage.Storage, logger log.Logger) *ListHandler {
	h, err := NewListHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *ListHandler) Handle(appChat *types.Chat) ([]types.Person, error) {
	if appChat.IsNotAGroup() {
		return nil, apperrors.ErrChatTypeIsUnsupported
	}

	return nil, nil
}