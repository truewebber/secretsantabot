package query

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/log"
)

type GetMyReceiverHandler struct{}

func NewGetMyReceiverHandler(logger log.Logger) (*GetMyReceiverHandler, error) {
	if logger == nil {
		return nil, errServiceIsNil
	}

	return &GetMyReceiverHandler{}, nil
}

func MustNewGetMyReceiverHandler(logger log.Logger) *GetMyReceiverHandler {
	h, err := NewGetMyReceiverHandler(logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *GetMyReceiverHandler) Handle(giver *types.Person) (*types.Person, error) {
	return &types.Person{}, nil
}
