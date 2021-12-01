package query

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/log"
)

type (
	ListHandler struct {
	}
)

func NewListHandler(logger log.Logger) (*ListHandler, error) {
	if logger == nil {
		return nil, errServiceIsNil
	}

	return &ListHandler{}, nil
}

func MustNewListHandler(logger log.Logger) *ListHandler {
	h, err := NewListHandler(logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *ListHandler) Handle() ([]types.Person, error) {
	return nil, nil
}
