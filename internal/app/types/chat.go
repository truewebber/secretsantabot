package types

import "github.com/truewebber/secretsantabot/internal/chat"

type Chat struct {
	Admin          *Person
	Participants   []Person
	Magic          Magic
	TelegramChatID int64
	IsGroup        bool
}

func ChatToDomain(c *Chat) *chat.Chat {
	return &chat.Chat{
		Admin:          PersonToDomain(c.Admin),
		TelegramChatID: c.TelegramChatID,
	}
}
