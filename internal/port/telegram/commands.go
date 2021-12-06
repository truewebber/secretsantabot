package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

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
	switch command {
	case EnrollCommand:
		t.Enroll(message)
	case DisEnrollCommand:
		t.DisEnroll(message)
	case ListCommand:
		t.List(message)
	case MagicCommand:
		t.Magic(message)
	case MyCommand:
		t.My(message)
	case HelpCommand:
		t.Help(message)
	case StartCommand:
		t.Start(message)
	}
}

func (t *Bot) Enroll(message *tgbotapi.Message) {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		t.logger.Infof("failed build person, message: %#v, error: %v", message, err)

		return
	}

	if err := t.application.Commands.Enroll.Handle(person); err != nil {
		t.logger.Infof("failed handle enroll, message: %#v, error: %v", message, err)

		return
	}

	replyMessage := t.builder.buildEnrollSuccessMessage(message.From, message.Chat)

	t.SendAndLogOnError(replyMessage)
}

func (t *Bot) DisEnroll(message *tgbotapi.Message) {
	person, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		t.logger.Infof("failed build person, message: %#v, error: %v", message, err)

		return
	}

	if err := t.application.Commands.DisEnroll.Handle(person); err != nil {
		t.logger.Infof("failed handle disenroll, message: %#v, error: %v", message, err)

		return
	}

	replyMessage := t.builder.buildDisEnrollSuccessMessage(message.From, message.Chat)

	t.SendAndLogOnError(replyMessage)
}

func (t *Bot) List(message *tgbotapi.Message) {
	participants, err := t.application.Queries.List.Handle()
	if err != nil {
		t.logger.Infof("failed handle disenroll, message: %#v, error: %v", message, err)

		return
	}

	replyMessage, err := t.builder.buildListOfParticipantsMessage(message.Chat, participants)
	if err != nil {
		t.logger.Errorf("failed build list of participants message, message: %#v, error: %v", message, err)

		return
	}

	t.SendAndLogOnError(replyMessage)
}

func (t *Bot) Magic(message *tgbotapi.Message) {
}

func (t *Bot) My(message *tgbotapi.Message) {
	giver, err := t.builder.buildPersonFromMessage(message)
	if err != nil {
		t.logger.Infof("failed build person, message: %#v, error: %v", message, err)

		return
	}

	receiver, err := t.application.Queries.GetMyReceiver.Handle(giver)
	if err != nil {
		t.logger.Errorf("failed get receiver by giver, message: %#v, error: %v", message, err)

		return
	}

	replyMessage, err := t.builder.buildMyReceiverMessage(giver.TelegramUserID, receiver)
	if err != nil {
		t.logger.Errorf("failed build my receiver message, message: %#v, error: %v", message, err)

		return
	}

	t.SendAndLogOnError(replyMessage)
}

const helpText = "/enroll - enroll the game\n" +
	"/end - stop your enroll (only before magic starts)\n" +
	"/list - list all enrolling people\n" +
	"/magic - start the game (only admin)\n" +
	"/my - SecretHelSanta will resend magic info for you (only in private chat wi me)\n" +
	"/help - show this message\n" +
	"/start - register new chat (don't work with private messages)\n"

func (t *Bot) Help(message *tgbotapi.Message) {
	replyMsg := tgbotapi.NewMessage(message.Chat.ID, helpText)

	t.SendAndLogOnError(&replyMsg)
}

const startText = "Dummy start text!"

func (t *Bot) Start(message *tgbotapi.Message) {
	chat, err := t.builder.buildChatFromMessage(message)
	if err != nil {
		t.logger.Infof("failed build chat, message: %#v, error: %v", message, err)

		return
	}

	if err := t.application.Commands.RegisterNewChat.Handle(chat); err != nil {
		t.logger.Errorf("failed register new chat, message: %#v, error: %v", message, err)

		return
	}

	replyMsg := tgbotapi.NewMessage(message.Chat.ID, startText)

	t.SendAndLogOnError(&replyMsg)
}
