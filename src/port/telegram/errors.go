package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const forbiddenMessage = "Forbidden: bot was blocked by the user"

func isForbidden(err error) bool {
	//nolint:errorlint // errors.As won't help here
	tgErr, ok := err.(tgbotapi.Error)
	if !ok {
		return false
	}

	return tgErr.Message == forbiddenMessage
}
