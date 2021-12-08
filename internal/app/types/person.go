package types

import "github.com/truewebber/secretsantabot/internal/chat"

type Person struct {
	TelegramUserID int
	TelegramChatID int64
}

func PersonToDomain(p *Person) *chat.Person {
	return &chat.Person{
		ID:             0,
		TelegramUserID: p.TelegramUserID,
		TelegramChatID: p.TelegramChatID,
	}
}
