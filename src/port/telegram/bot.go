package telegram

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/telebot.v3"

	"github.com/truewebber/secretsantabot/app"
	apperrors "github.com/truewebber/secretsantabot/app/errors"
	"github.com/truewebber/secretsantabot/app/types"
	"github.com/truewebber/secretsantabot/domain/log"
)

type Bot struct {
	bot         *telebot.Bot
	application *app.Application
	builder     builder
	me          *telebot.User
}

func NewTelegramBot(token string, application *app.Application, logger log.Logger) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
		OnError: func(err error, ctx telebot.Context) {
			logger.Error(
				"failed handle handler",
				"message", ctx.Text(),
				"chat", ctx.Chat(),
				"sender", ctx.Sender(),
				"error", err.Error())
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	return &Bot{
		bot:         bot,
		me:          bot.Me,
		application: application,
		builder:     newBuilder(bot),
	}, nil
}

func MustNewTelegramBot(token string, application *app.Application, logger log.Logger) *Bot {
	t, err := NewTelegramBot(token, application, logger)
	if err != nil {
		panic(err)
	}

	return t
}

func (t *Bot) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()

		t.bot.Stop()
	}()

	t.bot.Handle(EnrollCommand, t.Enroll)
	t.bot.Handle(DisEnrollCommand, t.DisEnroll)
	t.bot.Handle(ListCommand, t.List)
	t.bot.Handle(MagicCommand, t.Magic)
	t.bot.Handle(MyCommand, t.My)
	t.bot.Handle(HelpCommand, t.Help)
	t.bot.Handle(StartCommand, t.Start)
	t.bot.Handle(telebot.OnAddedToGroup, t.Start)

	t.bot.Start()
}

const (
	EnrollCommand    = "/enroll"
	DisEnrollCommand = "/disenroll"
	ListCommand      = "/list"
	MagicCommand     = "/magic"
	MyCommand        = "/my"
	HelpCommand      = "/help"
	StartCommand     = "/start"
)

//nolint:dupl // no sense to merge this func
func (t *Bot) Enroll(ctx telebot.Context) error {
	person, err := t.builder.buildPersonFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build chat from message: %w", err)
	}

	handleErr := t.application.Commands.Enroll.Handle(context.TODO(), chat, person)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle enroll: %w", handleErr)
	}

	replyMessage := t.builder.buildEnrollSuccessTextMessage(ctx.Message().Sender)

	if _, err := t.bot.Send(ctx.Chat(), replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

//nolint:dupl // no sense to merge this func
func (t *Bot) DisEnroll(ctx telebot.Context) error {
	person, err := t.builder.buildPersonFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build chat from message: %w", err)
	}

	handleErr := t.application.Commands.DisEnroll.Handle(context.TODO(), chat, person)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle disenroll: %w", handleErr)
	}

	replyMessage := t.builder.buildDisEnrollSuccessTextMessage(ctx.Message().Sender)

	if _, err := t.bot.Send(ctx.Chat(), replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) List(ctx telebot.Context) error {
	chat, err := t.builder.buildChatFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build chat from message: %w", err)
	}

	participants, err := t.application.Queries.ListParticipants.Handle(context.TODO(), chat)
	if err != nil {
		return fmt.Errorf("handle list of participants: %w", err)
	}

	replyMessage, err := t.builder.buildListOfParticipantsTextMessage(chat, participants)
	if err != nil {
		return fmt.Errorf("build list of participants message: %w", err)
	}

	if _, err := t.bot.Send(ctx.Chat(), replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) Magic(ctx telebot.Context) error {
	person, err := t.builder.buildPersonFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build chat from message: %w", err)
	}

	handleErr := t.application.Commands.Magic.Handle(context.TODO(), chat, person)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		replyMessage := t.builder.buildRestartChatTextMessage()

		if _, sendErr := t.bot.Send(ctx.Chat(), replyMessage); sendErr != nil {
			return fmt.Errorf("send message: %w", sendErr)
		}

		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle magic: %w", handleErr)
	}

	magic, err := t.application.Queries.GetMagic.Handle(context.TODO(), chat, person)
	if err != nil {
		return fmt.Errorf("handle get magic: %w", err)
	}

	for _, pair := range magic.Pairs {
		if err := t.notifyGiver(chat, pair.Giver, pair.Receiver); err != nil {
			return fmt.Errorf("notify giver: %w", err)
		}
	}

	replyMessage := t.builder.buildMagicTextMessage()

	if _, err := t.bot.Send(ctx.Chat(), replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) notifyGiver(chat types.Chat, giver, receiver types.Person) error {
	msg, err := t.builder.buildMyReceiverTextMessage(chat, receiver)
	if err != nil {
		return fmt.Errorf("build message: %w", err)
	}

	recipient := &telebot.User{ID: giver.TelegramUserID}

	if _, err := t.bot.Send(recipient, msg); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) My(ctx telebot.Context) error {
	giver, err := t.builder.buildPersonFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build giver person from message: %w", err)
	}

	chat, err := t.builder.buildChatFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build chat from message: %w", err)
	}

	receiver, handleErr := t.application.Queries.GetMyReceiver.Handle(context.TODO(), chat, giver)

	if errors.Is(handleErr, apperrors.ErrNotFound) {
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle get receiver by giver: %w", handleErr)
	}

	replyMessage, err := t.builder.buildMyReceiverTextMessage(chat, receiver)
	if err != nil {
		return fmt.Errorf("build my receiver message: %w", err)
	}

	recipient := &telebot.User{ID: giver.TelegramUserID}

	_, sendErr := t.bot.Send(recipient, replyMessage)

	if sendErr == nil {
		return nil
	}

	if isForbidden(sendErr) {
		if _, err := t.bot.Send(ctx.Chat(), "Please start me in private!"); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
	}

	return fmt.Errorf("send private message: %w", sendErr)
}

func (t *Bot) Help(ctx telebot.Context) error {
	replyMessage := t.builder.buildHelpTextMessage()

	if _, err := t.bot.Send(ctx.Chat(), replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) Start(ctx telebot.Context) error {
	chat, err := t.builder.buildChatFromContext(ctx)
	if err != nil {
		return fmt.Errorf("build chat from message: %w", err)
	}

	handleErr := t.application.Commands.RegisterNewChatAndVersion.Handle(context.TODO(), chat)

	if errors.Is(handleErr, apperrors.ErrAlreadyExists) {
		return nil
	}

	if handleErr != nil {
		return fmt.Errorf("handle register new chat: %w", handleErr)
	}

	replyMessage := t.builder.buildStartTextMessage()

	if _, err := t.bot.Send(ctx.Chat(), replyMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (t *Bot) IsMeANewMemberMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if t.isMeANewMember(c.Message().UsersJoined) {
			return next(c)
		}

		return nil
	}
}

func (t *Bot) isMeANewMember(newUsers []telebot.User) bool {
	for i := range newUsers {
		if t.me.ID == newUsers[i].ID {
			return true
		}
	}

	return false
}
