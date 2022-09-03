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

func (b builder) getTelegramUser(chatID, userID int64) (*tgbotapi.User, error) {
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

func (builder) buildPersonFromMessage(message *tgbotapi.Message) (*types.Person, error) {
	if message.From == nil {
		return nil, errUserIsNil
	}

	if message.Chat == nil {
		return nil, errChatIsNil
	}

	return &types.Person{
		TelegramUserID: int64(message.From.ID),
	}, nil
}

func (b builder) buildChatFromMessage(message *tgbotapi.Message) (*types.Chat, error) {
	person, err := b.buildPersonFromMessage(message)
	if err != nil {
		return nil, fmt.Errorf("build person from message: %w", err)
	}

	return &types.Chat{
		IsGroup:        message.Chat.IsGroup() || message.Chat.IsSuperGroup(),
		TelegramChatID: message.Chat.ID,
		Admin:          person,
	}, nil
}

const enrollSuccessMessageTemplate = "Congratulations!\n%s %s is having part in Secret Santa."

func (builder) buildEnrollSuccessMessage(from *tgbotapi.User, chat *tgbotapi.Chat) *tgbotapi.MessageConfig {
	text := fmt.Sprintf(enrollSuccessMessageTemplate, from.FirstName, from.LastName)

	replyMessage := tgbotapi.NewMessage(chat.ID, text)

	return &replyMessage
}

const disEnrollSuccessMessageTemplate = "Sad to see you leaving =(\n%s %s is not in game from now."

func (builder) buildDisEnrollSuccessMessage(from *tgbotapi.User, chat *tgbotapi.Chat) *tgbotapi.MessageConfig {
	text := fmt.Sprintf(disEnrollSuccessMessageTemplate, from.FirstName, from.LastName)

	replyMessage := tgbotapi.NewMessage(chat.ID, text)

	return &replyMessage
}

func (b builder) buildListOfParticipantsMessage(chat *types.Chat, participants []types.Person,
) (*tgbotapi.MessageConfig, error) {
	text, err := b.listOfParticipantsToText(chat, participants)
	if err != nil {
		return nil, fmt.Errorf("list of participants to text: %w", err)
	}

	replyMessage := tgbotapi.NewMessage(chat.TelegramChatID, text)

	return &replyMessage, nil
}

const ListIsEmptyMessage = "No one person has enroll yet."

func (b builder) listOfParticipantsToText(chat *types.Chat, participants []types.Person) (string, error) {
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

func (b builder) buildMyReceiverMessage(
	chat *types.Chat,
	recipient, receiver *types.Person,
) (*tgbotapi.MessageConfig, error) {
	text, err := b.getMyReceiverToText(chat, receiver)
	if err != nil {
		return nil, fmt.Errorf("receiver to text: %w", err)
	}

	replyMessage := tgbotapi.NewMessage(recipient.TelegramUserID, text)

	return &replyMessage, nil
}

const getMyReceiverMessageTemplate = "Hey! Your target is `%s %s%s`"

func (b builder) getMyReceiverToText(chat *types.Chat, receiver *types.Person) (string, error) {
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
