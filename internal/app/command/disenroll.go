package command

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/log"
)

type (
	DisEnrollHandler struct {
	}
)

func NewDisEnrollHandler(logger log.Logger) (*DisEnrollHandler, error) {
	if logger == nil {
		return nil, errServiceIsNil
	}

	return &DisEnrollHandler{}, nil
}

func MustNewDisEnrollHandler(logger log.Logger) *DisEnrollHandler {
	h, err := NewDisEnrollHandler(logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *DisEnrollHandler) Handle(p *types.Person) error {
	return nil
}
