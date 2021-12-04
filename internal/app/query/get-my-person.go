package query

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/log"
)

type (
	GetMyPersonHandler struct{}
)

func NewGetMyPersonHandler(logger log.Logger) (*GetMyPersonHandler, error) {
	if logger == nil {
		return nil, errServiceIsNil
	}

	return &GetMyPersonHandler{}, nil
}

func MustNewGetMyPersonHandler(logger log.Logger) *GetMyPersonHandler {
	h, err := NewGetMyPersonHandler(logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *GetMyPersonHandler) Handle() (*types.Person, error) {
	return &types.Person{}, nil
}
