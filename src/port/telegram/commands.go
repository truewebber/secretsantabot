package telegram

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	apperrors "github.com/truewebber/secretsantabot/app/errors"
)

const (
	EnrollCommand    = "enroll"
	DisEnrollCommand = "disenroll"
	ListCommand      = "list"
	MagicCommand     = "magic"
	MyCommand        = "my"
	HelpCommand      = "help"
	StartCommand     = "start"
)

func (t *Bot) processCommand(ctx context.Context, command string, message *tgbotapi.Message) {
	var err error

	switch command {
	case EnrollCommand:
		err = t.Enroll(ctx, message)
	case DisEnrollCommand:
		err = t.DisEnroll(ctx, message)
	case ListCommand:
		err = t.List(ctx, message)
	case MagicCommand:
		err = t.Magic(ctx, message)
	case MyCommand:
		err = t.My(ctx, message)
	case HelpCommand:
		err = t.Help(ctx, message)
	case StartCommand:
		err = t.Start(ctx, message)
	}

	if err != nil {
		t.logger.Errorf("failed process `%s` command, message: %v, error: %v", command, message, err)

		return
	}

	t.logger.Infof("`%s`, from `%#v` in chat `%#v`", command, message.From, message.Chat)
}

//nolint:dupl // no sense to merge this func
func (t *Bot) Enroll(ctx context.Context, message *tgbotapi.Message) error {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	handleErr := t.application.Commands.Enroll.Handle(ctx, chat, person)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle enroll: %w", err)
	}

	replyMessage := t.builder.buildEnrollSuccessMessage(message.From, message.Chat)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

//nolint:dupl // no sense to merge this func
func (t *Bot) DisEnroll(ctx context.Context, message *tgbotapi.Message) error {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	handleErr := t.application.Commands.DisEnroll.Handle(ctx, chat, person)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle disenroll: %w", err)
	}

	replyMessage := t.builder.buildDisEnrollSuccessMessage(message.From, message.Chat)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) List(ctx context.Context, message *tgbotapi.Message) error {
	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	participants, err := t.application.Queries.ListParticipants.Handle(ctx, chat)
	if err != nil {
		return fmt.Errorf("handle list of participants: %w", err)
	}

	replyMessage, err := t.builder.buildListOfParticipantsMessage(chat, participants)
	if err != nil {
		return fmt.Errorf("build list of participants message: %w", err)
	}

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) Magic(ctx context.Context, message *tgbotapi.Message) error {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	handleErr := t.application.Commands.Magic.Handle(ctx, chat, person)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		// send reply to restart game
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle magic: %w", handleErr)
	}

	magic, err := t.application.Queries.GetMagic.Handle(ctx, chat, person)
	if err != nil {
		return fmt.Errorf("handle get magic: %w", err)
	}

	for _, pair := range magic.Pairs {
		if err := t.builder.notifyGiver(pair.Giver, pair.Receiver); err != nil {
			return fmt.Errorf("notify giver: %w", err)
		}
	}

	replyMessage := t.builder.buildMagicMessage(chat)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) My(ctx context.Context, message *tgbotapi.Message) error {
	giver, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	receiver, handleErr := t.application.Queries.GetMyReceiver.Handle(ctx, chat, giver)
	if handleErr != nil {
		return fmt.Errorf("handle get receiver by giver: %w", handleErr)
	}

	replyMessage, err := t.builder.buildMyReceiverMessage(chat, giver, receiver)
	if err != nil {
		return fmt.Errorf("build my receiver message: %w", err)
	}

	_, sendErr := t.bot.Send(replyMessage)

	if sendErr == nil {
		return nil
	}

	if isForbidden(sendErr) {
		errMessage := tgbotapi.NewMessage(message.Chat.ID, "Please start me in private!")

		if _, err := t.bot.Send(errMessage); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
	}

	return fmt.Errorf("send private message: %w", sendErr)
}

func (t *Bot) Help(_ context.Context, message *tgbotapi.Message) error {
	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	replyMessage := t.builder.buildHelpMessage(chat)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) Start(ctx context.Context, message *tgbotapi.Message) error {
	chat, buildErr := t.builder.buildChatFromMessage(message)
	if buildErr != nil {
		return fmt.Errorf("build chat from message: %w", buildErr)
	}

	handleErr := t.application.Commands.RegisterNewChatAndVersion.Handle(ctx, chat)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle register new chat: %w", handleErr)
	}

	replyMessage := t.builder.buildStartMessage(chat)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
