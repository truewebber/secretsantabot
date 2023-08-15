package telegram

import (
	"errors"

	"gopkg.in/telebot.v3"
)

const forbiddenMessage = "Forbidden: bot was blocked by the user"

func isForbidden(err error) bool {
	var tgErr *telebot.Error

	ok := errors.As(err, &tgErr)
	if !ok {
		return false
	}

	return tgErr.Message == forbiddenMessage
}
