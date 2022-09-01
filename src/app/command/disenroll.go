package command

import (
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/chat/storage"
	"github.com/truewebber/secretsantabot/domain/log"
)

type DisEnrollHandler struct {
	service storage.Storage
}

func NewDisEnrollHandler(service storage.Storage, logger log.Logger) (*DisEnrollHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &DisEnrollHandler{service: service}, nil
}

func MustNewDisEnrollHandler(service storage.Storage, logger log.Logger) *DisEnrollHandler {
	h, err := NewDisEnrollHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *DisEnrollHandler) Handle(participant *types.Person) error {
	return nil
}
