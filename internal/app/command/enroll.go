package command

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/log"
)

type (
	EnrollHandler struct {
	}
)

func NewEnrollHandler(logger log.Logger) (*EnrollHandler, error) {
	if logger == nil {
		return nil, errServiceIsNil
	}

	return &EnrollHandler{}, nil
}

func MustNewEnrollHandler(logger log.Logger) *EnrollHandler {
	h, err := NewEnrollHandler(logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *EnrollHandler) Handle(p *types.Person) error {
	return nil
}
