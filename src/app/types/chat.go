package types

import "github.com/truewebber/secretsantabot/domain/chat"

type Chat struct {
	Participants   []Person
	Admin          Person
	TelegramChatID int64
	ChatType       ChatType
}

type ChatType uint8

const (
	ChatTypeUnsupported ChatType = iota
	ChatTypeGroup
	ChatTypePrivate
)

func (c *Chat) IsNotAGroup() bool {
	return c.ChatType != ChatTypeGroup
}

func (c *Chat) IsPrivate() bool {
	return c.ChatType == ChatTypePrivate
}

func (c *Chat) IsUnsupported() bool {
	return c.ChatType == ChatTypeUnsupported
}

func ChatToDomain(c Chat) chat.Chat {
	return chat.Chat{
		Admin:          PersonToDomain(c.Admin),
		TelegramChatID: c.TelegramChatID,
	}
}
