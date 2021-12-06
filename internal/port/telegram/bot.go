package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/truewebber/secretsantabot/internal/app"
	"github.com/truewebber/secretsantabot/internal/log"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	application *app.Application
	builder     builder
	logger      log.Logger
}

func NewTelegramBot(token string, application *app.Application, logger log.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	return &Bot{
		bot:         bot,
		application: application,
		builder:     newBuilder(bot),
		logger:      logger,
	}, nil
}

func MustNewTelegramBot(token string, application *app.Application, logger log.Logger) *Bot {
	t, err := NewTelegramBot(token, application, logger)
	if err != nil {
		panic(err)
	}

	return t
}

func (t *Bot) Run(ctx context.Context) error {
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

func (t *Bot) processUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
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

func (t *Bot) processOneUpdate(update *tgbotapi.Update) {
	message := messageFromUpdate(update)
	if message == nil {
		return
	}

	t.processMessage(message)
}

func messageFromUpdate(update *tgbotapi.Update) *tgbotapi.Message {
	if update.Message != nil {
		return update.Message
	}

	if update.EditedMessage != nil {
		return update.EditedMessage
	}

	return nil
}

func (t *Bot) processMessage(message *tgbotapi.Message) {
	command := message.Command()
	if command == "" {
		return
	}

	t.processCommand(command, message)
}

func (t *Bot) SendAndLogOnError(msg *tgbotapi.MessageConfig) {
	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Errorf("send message: %w", err)
	}
}
