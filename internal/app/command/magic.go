package command

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/chat/storage"
	"github.com/truewebber/secretsantabot/internal/log"
)

type MagicHandler struct {
	service storage.Storage
}

func NewMagicHandler(service storage.Storage, logger log.Logger) (*MagicHandler, error) {
	if service == nil || logger == nil {
		return nil, errServiceIsNil
	}

	return &MagicHandler{service: service}, nil
}

func MustNewMagicHandler(service storage.Storage, logger log.Logger) *MagicHandler {
	h, err := NewMagicHandler(service, logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *MagicHandler) Handle(notifyPersonFn func(p *types.Person) error) error {
	return nil
}
