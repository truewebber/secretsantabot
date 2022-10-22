package types

import "github.com/truewebber/secretsantabot/domain/chat"

type (
	Magic struct {
		Pairs []GiverReceiverPair
	}

	GiverReceiverPair struct {
		Giver    Person
		Receiver Person
	}
)

func DomainToMagic(domainMagic chat.Magic) Magic {
	return Magic{
		Pairs: DomainToPairs(domainMagic.Pairs),
	}
}

func DomainToPairs(domainPairs []chat.GiverReceiverPair) []GiverReceiverPair {
	pairs := make([]GiverReceiverPair, 0, len(domainPairs))

	for i := range domainPairs {
		pair := DomainToPair(domainPairs[i])

		pairs = append(pairs, pair)
	}

	return pairs
}

func DomainToPair(domainPair chat.GiverReceiverPair) GiverReceiverPair {
	return GiverReceiverPair{
		Giver:    DomainToPerson(domainPair.Giver),
		Receiver: DomainToPerson(domainPair.Receiver),
	}
}
