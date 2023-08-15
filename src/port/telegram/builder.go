package telegram

import (
	"errors"
	"fmt"

	"gopkg.in/telebot.v3"

	"github.com/truewebber/secretsantabot/app/types"
)

type builder struct {
	bot *telebot.Bot
}

func newBuilder(bot *telebot.Bot) builder {
	return builder{
		bot: bot,
	}
}

func (b *builder) getTelegramUser(chatID, userID int64) (*telebot.User, error) {
	member, err := b.bot.ChatMemberOf(&telebot.Chat{ID: chatID}, &telebot.User{ID: userID})
	if err != nil {
		return nil, fmt.Errorf("chat member of: %w", err)
	}

	return member.User, nil
}

var (
	errChatIsNil = errors.New("chat is nil")
	errUserIsNil = errors.New("user is nil")
)

func (*builder) buildPersonFromContext(ctx telebot.Context) (types.Person, error) {
	if ctx.Message().Sender == nil {
		return types.Person{}, errUserIsNil
	}

	if ctx.Chat() == nil {
		return types.Person{}, errChatIsNil
	}

	return types.Person{
		TelegramUserID: ctx.Message().Sender.ID,
	}, nil
}

func (b *builder) buildChatFromContext(ctx telebot.Context) (types.Chat, error) {
	person, err := b.buildPersonFromContext(ctx)
	if err != nil {
		return types.Chat{}, fmt.Errorf("build person from message: %w", err)
	}

	return types.Chat{
		ChatType:       b.buildChatType(ctx.Chat()),
		TelegramChatID: ctx.Chat().ID,
		Admin:          person,
	}, nil
}

func (b *builder) buildChatType(chat *telebot.Chat) types.ChatType {
	if chat.Type == telebot.ChatGroup || chat.Type == telebot.ChatSuperGroup {
		return types.ChatTypeGroup
	}

	if chat.Type == telebot.ChatPrivate {
		return types.ChatTypePrivate
	}

	return types.ChatTypeUnsupported
}

const enrollSuccessMessageTemplate = "Congratulations!\n%s %s is having part in Secret Santa."

func (b *builder) buildEnrollSuccessTextMessage(from *telebot.User) string {
	return fmt.Sprintf(enrollSuccessMessageTemplate, from.FirstName, from.LastName)
}

const disEnrollSuccessMessageTemplate = "Sad to see you leaving =(\n%s %s is not in game from now."

func (*builder) buildDisEnrollSuccessTextMessage(from *telebot.User) string {
	return fmt.Sprintf(disEnrollSuccessMessageTemplate, from.FirstName, from.LastName)
}

func (b *builder) buildListOfParticipantsTextMessage(
	chat types.Chat, participants []types.Person,
) (string, error) {
	text, err := b.listOfParticipantsToText(chat, participants)
	if err != nil {
		return "", fmt.Errorf("list of participants to text: %w", err)
	}

	return text, nil
}

const ListIsEmptyMessage = "No one person has enroll yet."

func (b *builder) listOfParticipantsToText(chat types.Chat, participants []types.Person) (string, error) {
	if len(participants) == 0 {
		return ListIsEmptyMessage, nil
	}

	var text string

	for index, participant := range participants {
		user, err := b.getTelegramUser(chat.TelegramChatID, participant.TelegramUserID)
		if err != nil {
			return "", fmt.Errorf("get telegram user: %w", err)
		}

		text += fmt.Sprintf("%s %s", user.FirstName, user.LastName)

		if user.Username != "" {
			text += fmt.Sprintf(" (@%s)", user.Username)
		}

		if index != len(participants)-1 {
			text += "\n"
		}
	}

	return text, nil
}

func (b *builder) buildMyReceiverTextMessage(
	chat types.Chat,
	receiver types.Person,
) (string, error) {
	text, err := b.getMyReceiverToText(chat, receiver)
	if err != nil {
		return "", fmt.Errorf("receiver to text: %w", err)
	}

	return text, nil
}

const getMyReceiverMessageTemplate = "Hey! Your target is `%s %s%s`"

func (b *builder) getMyReceiverToText(chat types.Chat, receiver types.Person) (string, error) {
	user, err := b.getTelegramUser(chat.TelegramChatID, receiver.TelegramUserID)
	if err != nil {
		return "", fmt.Errorf("get telegram user: %w", err)
	}

	var usernameText string
	if user.Username != "" {
		usernameText = fmt.Sprintf(" (@%s)", user.Username)
	}

	text := fmt.Sprintf(getMyReceiverMessageTemplate, user.FirstName, user.LastName, usernameText)

	return text, nil
}

const magicText = "Ho-ho-ho!\nLet's the Christmas begin 🎁\nAll of you should receive the private message from me!\n" +
	"If not, press @secrethellsantabot and press start or restart. After that you could press the /my command."

func (b *builder) buildMagicTextMessage() string {
	return magicText
}

const startText = "Ho-ho-ho!\nWelcome guys and Merry Christmas 🎁\n\nTo start game, every " +
	"one who wants to participate need to send message /enroll to the chat, also, you need " +
	"to allow me to write to you in direct. Press @secrethellsantabot and press start or restart.\n" +
	"After that, my inviter should begin the MAGIC (send message /magic)."

func (b *builder) buildStartTextMessage() string {
	return startText
}

const helpText = "/enroll - enroll the game\n" +
	"/disenroll - stop your enroll (only before magic starts)\n" +
	"/list - list all enrolling people\n" +
	"/magic - start the game (only admin)\n" +
	"/my - Secret Santa will resend magic info for you (only in private chat with me)\n" +
	"/help - show this message\n" +
	"/start - register new chat (don't work with private messages)\n"

func (b *builder) buildHelpTextMessage() string {
	return helpText
}

const restartChatText = "Ho-ho-ho!\nMagic already happened!\nIf you wanna to make MAGIC again, do the restart command."

func (b *builder) buildRestartChatTextMessage() string {
	return restartChatText
}
