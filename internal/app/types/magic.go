package types

import (
	"errors"
	"fmt"

	"github.com/truewebber/secretsantabot/internal/chat"
)

type (
	Magic struct {
		Data   []GiverReceiverPair
		Status MagicStatus
	}

	GiverReceiverPair struct {
		Giver    *Person
		Receiver *Person
	}

	MagicStatus uint8
)

const (
	OpenMagicStatus MagicStatus = iota + 1
	ClosedMagicStatus
)

func MagicToDomain(magic Magic) (chat.Magic, error) {
	status, err := magicStatusToDomain(magic.Status)
	if err != nil {
		return chat.Magic{}, fmt.Errorf("magic status to domain: %w", err)
	}

	return chat.Magic{
		Data:   dataToDomain(magic.Data),
		Status: status,
	}, nil
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

var errUnknownMagicStatus = errors.New("unknown magic status")

func magicStatusToDomain(status MagicStatus) (chat.MagicStatus, error) {
	switch status {
	case OpenMagicStatus:
		return chat.OpenMagicStatus, nil
	case ClosedMagicStatus:
		return chat.ClosedMagicStatus, nil
	}

	return 0, errUnknownMagicStatus
}
