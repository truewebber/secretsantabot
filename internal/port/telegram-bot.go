package port

import (
	"context"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/truewebber/secretsantabot/internal/app"
	"github.com/truewebber/secretsantabot/internal/log"
)

type TelegramBot struct {
	bot         *tgbotapi.BotAPI
	application *app.Application
	logger      log.Logger
}

func NewTelegramBot(token string, application *app.Application, logger log.Logger) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	return &TelegramBot{
		bot:         bot,
		application: application,
		logger:      logger,
	}, nil
}

func MustNewTelegramBot(token string, application *app.Application, logger log.Logger) *TelegramBot {
	t, err := NewTelegramBot(token, application, logger)
	if err != nil {
		panic(err)
	}

	return t
}

func (t *TelegramBot) Run(ctx context.Context) error {
	const (
		numberOfUpdates = 20
		updatesTimeout  = 5
	)

	u := tgbotapi.NewUpdate(numberOfUpdates)
	u.Timeout = updatesTimeout

	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("get updates: %w", err)
	}

	t.processUpdates(ctx, updates)

	return nil
}

func (t *TelegramBot) processUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case u, ok := <-updates:
			if !ok {
				return
			}

			t.processOneUpdate(&u)
		}
	}
}

func (t *TelegramBot) processOneUpdate(update *tgbotapi.Update) {
	command := update.Message.Command()
	if command == "" {
		return
	}

	t.processCommand(command, update.Message)
}

const (
	EnrollCommand    = "enroll"
	DisEnrollCommand = "disenroll"
	ListCommand      = "list"
	MagicCommand     = "magic"
	MyCommand        = "my"
	HelpCommand      = "help"
)

func (t *TelegramBot) processCommand(command string, message *tgbotapi.Message) {
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
	}
}

func (t *TelegramBot) Enroll(message *tgbotapi.Message) {
}

func (t *TelegramBot) DisEnroll(message *tgbotapi.Message) {
}

func (t *TelegramBot) List(message *tgbotapi.Message) {
}

func (t *TelegramBot) Magic(message *tgbotapi.Message) {
}

func (t *TelegramBot) My(message *tgbotapi.Message) {
}

const helpText = "/enroll - enroll the game\n" +
	"/end - stop your enroll (only before magic starts)\n" +
	"/list - list all enrolling people\n" +
	"/magic - start the game (only admin)\n" +
	"/my - SecretHelSanta will resend magic info for you (only in private chat wi me)\n" +
	"/help - show this message\n"

func (t *TelegramBot) Help(msg *tgbotapi.Message) {
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, helpText)

	t.SendAndLogOnError(&replyMsg)
}

func (t *TelegramBot) SendAndLogOnError(msg *tgbotapi.MessageConfig) {
	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Errorf("send message: %w", err)
	}
}
