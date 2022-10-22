package magic

import (
	"errors"
	"math/rand"
	"time"

	"github.com/truewebber/secretsantabot/domain/chat"
)

var ErrNotEnoughParticipants = errors.New("not enough participants")

const minAmountOfParticipants = 2

func Calculate(participants []chat.Person) (chat.Magic, error) {
	if len(participants) < minAmountOfParticipants {
		return chat.Magic{}, ErrNotEnoughParticipants
	}

	shuffle(participants)

	pairs := makePairs(participants)

	return chat.Magic{Pairs: pairs}, nil
}

func makePairs(participants []chat.Person) []chat.GiverReceiverPair {
	pairs := make([]chat.GiverReceiverPair, 0, len(participants))

	for i := range participants {
		var (
			giverIdx    = i
			receiverIdx = i + 1
		)

		if receiverIdx >= len(participants) {
			receiverIdx = 0
		}

		pair := chat.GiverReceiverPair{
			Giver:    participants[giverIdx],
			Receiver: participants[receiverIdx],
		}

		pairs = append(pairs, pair)
	}

	return pairs
}

func shuffle(participants []chat.Person) {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(
		len(participants),
		func(i, j int) {
			participants[i], participants[j] = participants[j], participants[i]
		},
	)
}
