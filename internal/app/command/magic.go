package command

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/log"
)

type (
	MagicHandler struct {
	}
)

func NewMagicHandler(logger log.Logger) (*MagicHandler, error) {
	if logger == nil {
		return nil, errServiceIsNil
	}

	return &MagicHandler{}, nil
}

func MustNewMagicHandler(logger log.Logger) *MagicHandler {
	h, err := NewMagicHandler(logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *MagicHandler) Handle(notifyPersonFn func(p *types.Person) error) error {
	return nil
}
