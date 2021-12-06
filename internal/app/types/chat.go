package types

type Chat struct {
	Admin          *Person
	Participants   []Person
	Magic          Magic
	TelegramChatID int64
}
