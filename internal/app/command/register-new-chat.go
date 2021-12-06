package command

import (
	"github.com/truewebber/secretsantabot/internal/app/types"
	"github.com/truewebber/secretsantabot/internal/log"
)

type RegisterNewChatHandler struct{}

func NewRegisterNewChatHandler(logger log.Logger) (*RegisterNewChatHandler, error) {
	if logger == nil {
		return nil, errServiceIsNil
	}

	return &RegisterNewChatHandler{}, nil
}

func MustNewRegisterNewChatHandler(logger log.Logger) *RegisterNewChatHandler {
	h, err := NewRegisterNewChatHandler(logger)
	if err != nil {
		panic(err)
	}

	return h
}

func (h *RegisterNewChatHandler) Handle(chat *types.Chat) error {
	return nil
}
