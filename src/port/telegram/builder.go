package telegram

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/truewebber/secretsantabot/app/types"
)

type builder struct {
	bot *tgbotapi.BotAPI
}

func newBuilder(bot *tgbotapi.BotAPI) builder {
	return builder{
		bot: bot,
	}
}

func (b *builder) getTelegramUser(chatID, userID int64) (*tgbotapi.User, error) {
	member, err := b.bot.GetChatMember(tgbotapi.ChatConfigWithUser{ChatID: chatID, UserID: int(userID)})
	if err != nil {
		return nil, fmt.Errorf("get chat member: %w", err)
	}

	return member.User, nil
}

var (
	errChatIsNil = errors.New("chat is nil")
	errUserIsNil = errors.New("user is nil")
)

func (*builder) buildPersonFromMessage(message *tgbotapi.Message) (types.Person, error) {
	if message.From == nil {
		return types.Person{}, errUserIsNil
	}

	if message.Chat == nil {
		return types.Person{}, errChatIsNil
	}

	return types.Person{
		TelegramUserID: int64(message.From.ID),
	}, nil
}

func (b *builder) buildChatFromMessage(message *tgbotapi.Message) (types.Chat, error) {
	person, err := b.buildPersonFromMessage(message)
	if err != nil {
		return types.Chat{}, fmt.Errorf("build person from message: %w", err)
	}

	return types.Chat{
		ChatType:       b.buildChatType(message.Chat),
		TelegramChatID: message.Chat.ID,
		Admin:          person,
	}, nil
}

func (b *builder) buildChatType(chat *tgbotapi.Chat) types.ChatType {
	if chat.IsGroup() || chat.IsSuperGroup() {
		return types.ChatTypeGroup
	}

	if chat.IsPrivate() {
		return types.ChatTypePrivate
	}

	return types.ChatTypeUnsupported
}

const enrollSuccessMessageTemplate = "Congratulations!\n%s %s is having part in Secret Santa."

func (*builder) buildEnrollSuccessMessage(from *tgbotapi.User, chat *tgbotapi.Chat) *tgbotapi.MessageConfig {
	text := fmt.Sprintf(enrollSuccessMessageTemplate, from.FirstName, from.LastName)

	replyMessage := tgbotapi.NewMessage(chat.ID, text)

	return &replyMessage
}

const disEnrollSuccessMessageTemplate = "Sad to see you leaving =(\n%s %s is not in game from now."

func (*builder) buildDisEnrollSuccessMessage(from *tgbotapi.User, chat *tgbotapi.Chat) *tgbotapi.MessageConfig {
	text := fmt.Sprintf(disEnrollSuccessMessageTemplate, from.FirstName, from.LastName)

	replyMessage := tgbotapi.NewMessage(chat.ID, text)

	return &replyMessage
}

func (b *builder) buildListOfParticipantsMessage(
	chat types.Chat, participants []types.Person,
) (*tgbotapi.MessageConfig, error) {
	text, err := b.listOfParticipantsToText(chat, participants)
	if err != nil {
		return nil, fmt.Errorf("list of participants to text: %w", err)
	}

	replyMessage := tgbotapi.NewMessage(chat.TelegramChatID, text)

	return &replyMessage, nil
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

		if user.UserName != "" {
			text += fmt.Sprintf(" (@%s)", user.UserName)
		}

		if index != len(participants)-1 {
			text += "\n"
		}
	}

	return text, nil
}

func (b *builder) buildMyReceiverMessage(
	chat types.Chat,
	recipient, receiver types.Person,
) (*tgbotapi.MessageConfig, error) {
	text, err := b.getMyReceiverToText(chat, receiver)
	if err != nil {
		return nil, fmt.Errorf("receiver to text: %w", err)
	}

	replyMessage := tgbotapi.NewMessage(recipient.TelegramUserID, text)

	return &replyMessage, nil
}

const getMyReceiverMessageTemplate = "Hey! Your target is `%s %s%s`"

func (b *builder) getMyReceiverToText(chat types.Chat, receiver types.Person) (string, error) {
	user, err := b.getTelegramUser(chat.TelegramChatID, receiver.TelegramUserID)
	if err != nil {
		return "", fmt.Errorf("get telegram user: %w", err)
	}

	var usernameText string
	if user.UserName != "" {
		usernameText = fmt.Sprintf(" (@%s)", user.UserName)
	}

	text := fmt.Sprintf(getMyReceiverMessageTemplate, user.FirstName, user.LastName, usernameText)

	return text, nil
}

const magicText = "Ho-ho-ho!\nLet's the Christmas begin 🎁\nAll of you should receive the private message from me!\n" +
	"If not, press @secrethellsantabot and press start or restart. After that you could press the /my command."

func (b *builder) buildMagicMessage(chat types.Chat) *tgbotapi.MessageConfig {
	replyMessage := tgbotapi.NewMessage(chat.TelegramChatID, magicText)

	return &replyMessage
}

func (b *builder) buildStartMessage(chat types.Chat) *tgbotapi.MessageConfig {
	const startText = "Ho-ho-ho!\nWelcome guys and Merry Christmas 🎁\n\nTo start game, every " +
		"one who wants to participate need to send message /enroll to the chat, also, you need " +
		"to allow me to write to you in direct. Press @secrethellsantabot and press start or restart.\n" +
		"After that, my inviter should begin the MAGIC (send message /magic)."

	replyMessage := tgbotapi.NewMessage(chat.TelegramChatID, startText)

	return &replyMessage
}

const helpText = "/enroll - enroll the game\n" +
	"/disenroll - stop your enroll (only before magic starts)\n" +
	"/list - list all enrolling people\n" +
	"/magic - start the game (only admin)\n" +
	"/my - Secret Santa will resend magic info for you (only in private chat with me)\n" +
	"/help - show this message\n" +
	"/start - register new chat (don't work with private messages)\n"

func (b *builder) buildHelpMessage(chat types.Chat) *tgbotapi.MessageConfig {
	replyMessage := tgbotapi.NewMessage(chat.TelegramChatID, helpText)

	return &replyMessage
}

const restartChatText = "Ho-ho-ho!\nMagic already happened!\nIf you wanna to make MAGIC again, do the restart command."

func (b *builder) buildRestartChatMessage(chat types.Chat) *tgbotapi.MessageConfig {
	replyMessage := tgbotapi.NewMessage(chat.TelegramChatID, restartChatText)

	return &replyMessage
}

func (b *builder) notifyGiver(_, _ types.Person) error {
	return nil
}
