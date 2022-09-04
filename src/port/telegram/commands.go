package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/truewebber/secretsantabot/app/types"
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

func (t *Bot) processCommand(command string, message *tgbotapi.Message) {
	var err error

	switch command {
	case EnrollCommand:
		err = t.Enroll(message)
	case DisEnrollCommand:
		err = t.DisEnroll(message)
	case ListCommand:
		err = t.List(message)
	case MagicCommand:
		err = t.Magic(message)
	case MyCommand:
		err = t.My(message)
	case HelpCommand:
		err = t.Help(message)
	case StartCommand:
		err = t.Start(message)
	}

	if err != nil {
		t.logger.Errorf("failed process `%s` command, message: %v, error: %v", command, message, err)

		return
	}

	t.logger.Infof("`%s`, from `%#v` in chat `%#v`", command, message.From, message.Chat)
}

//nolint:dupl // no sense to merge this func
func (t *Bot) Enroll(message *tgbotapi.Message) error {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	if err := t.application.Commands.Enroll.Handle(chat, person); err != nil {
		return fmt.Errorf("handle enroll: %w", err)
	}

	replyMessage := t.builder.buildEnrollSuccessMessage(message.From, message.Chat)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

//nolint:dupl // no sense to merge this func
func (t *Bot) DisEnroll(message *tgbotapi.Message) error {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	if err := t.application.Commands.DisEnroll.Handle(chat, person); err != nil {
		return fmt.Errorf("handle disenroll: %w", err)
	}

	replyMessage := t.builder.buildDisEnrollSuccessMessage(message.From, message.Chat)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) List(message *tgbotapi.Message) error {
	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	participants, err := t.application.Queries.List.Handle(chat)
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

func (t *Bot) Magic(message *tgbotapi.Message) error {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	handleErr := t.application.Commands.Magic.Handle(chat, person, func(p *types.Person) error {
		return nil
	})
	if handleErr != nil {
		return fmt.Errorf("handle magic: %w", handleErr)
	}

	return nil
}

func (t *Bot) My(message *tgbotapi.Message) error {
	giver, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	_, handleErr := t.application.Queries.GetMyReceiver.Handle(chat, giver)
	if handleErr != nil {
		return fmt.Errorf("handle get receiver by giver: %w", handleErr)
	}

	// temp
	receiver := giver

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

const helpText = "/enroll - enroll the game\n" +
	"/disenroll - stop your enroll (only before magic starts)\n" +
	"/list - list all enrolling people\n" +
	"/magic - start the game (only admin)\n" +
	"/my - Secret Santa will resend magic info for you (only in private chat with me)\n" +
	"/help - show this message\n" +
	"/start - register new chat (don't work with private messages)\n"

func (t *Bot) Help(message *tgbotapi.Message) error {
	replyMessage := tgbotapi.NewMessage(message.Chat.ID, helpText)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

const (
	startText = "Ho-ho-ho!\nWelcome guys and Merry Christmas 🎁\n\nTo start game, every " +
		"one who wants to participate need to send message /enroll to the chat, also, you need " +
		"to allow me to write to you in direct. Press @secrethellsantabot and press start or restart.\n" +
		"After that, my inviter should begin the MAGIC (send message /magic)."
)

func (t *Bot) Start(message *tgbotapi.Message) error {
	chat, buildErr := t.builder.buildChatFromMessage(message)
	if buildErr != nil {
		return fmt.Errorf("build chat from message: %w", buildErr)
	}

	if err := t.application.Commands.RegisterNewChatAndVersion.Handle(chat); err != nil {
		return fmt.Errorf("handle register new chat: %w", err)
	}

	replyMessage := tgbotapi.NewMessage(message.Chat.ID, startText)

	if _, err := t.bot.Send(replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
