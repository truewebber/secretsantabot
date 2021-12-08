package query

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/chat/storage"
	"github.com/truewebber/secretsantabot/internal/log"
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

func (h *ListHandler) Handle() ([]types.Person, error) {
	return nil, nil
}
