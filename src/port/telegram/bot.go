package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/truewebber/secretsantabot/app"
	"github.com/truewebber/secretsantabot/domain/log"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	application *app.Application
	builder     builder
	logger      log.Logger
	me          *tgbotapi.User
}

func NewTelegramBot(token string, application *app.Application, logger log.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	user, err := bot.GetMe()
	if err != nil {
		return nil, fmt.Errorf("get me: %w", err)
	}

	return &Bot{
		bot:         bot,
		me:          &user,
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

			t.processOneUpdate(ctx, &u)
		}
	}
}

func (t *Bot) processOneUpdate(ctx context.Context, update *tgbotapi.Update) {
	message := messageFromUpdate(update)
	if message == nil {
		return
	}

	t.processMessage(ctx, message)
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

func (t *Bot) processMessage(ctx context.Context, message *tgbotapi.Message) {
	command := t.getCommandFromMessage(message)
	if command == "" {
		return
	}

	t.processCommand(ctx, command, message)
}

func (t *Bot) getCommandFromMessage(message *tgbotapi.Message) string {
	command := message.Command()

	if t.isMeANewMember(message.NewChatMembers) && command == "" {
		command = StartCommand
	}

	return command
}

func (t *Bot) isMeANewMember(newUsersPtr *[]tgbotapi.User) bool {
	if newUsersPtr == nil {
		return false
	}

	newUsers := *newUsersPtr

	for i := range newUsers {
		if t.me.ID == newUsers[i].ID {
			return true
		}
	}

	return false
}
