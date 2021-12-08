//nolint:dupl // while under dev
package command

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/chat/storage"
	"github.com/truewebber/secretsantabot/internal/log"
)

type EnrollHandler struct {
	service storage.Storage
}

func NewEnrollHandler(service storage.Storage, logger log.Logger) (*EnrollHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &EnrollHandler{service: service}, nil
}

func MustNewEnrollHandler(service storage.Storage, logger log.Logger) *EnrollHandler {
	h, err := NewEnrollHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *EnrollHandler) Handle(participant *types.Person) error {
	return nil
}
