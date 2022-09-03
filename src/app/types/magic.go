package types

import "github.com/truewebber/secretsantabot/domain/chat"

type (
	Magic struct {
		Data []GiverReceiverPair
	}

	GiverReceiverPair struct {
		Giver    *Person
		Receiver *Person
	}
)

func MagicToDomain(magic *Magic) *chat.Magic {
	return &chat.Magic{
		Data: dataToDomain(magic.Data),
	}
}

func dataToDomain(data []GiverReceiverPair) []chat.GiverReceiverPair {
	chatData := make([]chat.GiverReceiverPair, 0, len(data))

	for i := range data {
		chatData = append(chatData, chat.GiverReceiverPair{
			Giver:    PersonToDomain(data[i].Giver),
			Receiver: PersonToDomain(data[i].Receiver),
		})
	}

	return chatData
}
