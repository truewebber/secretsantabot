package chat

type (
	Magic struct {
		Pairs []GiverReceiverPair
	}

	GiverReceiverPair struct {
		Giver    Person
		Receiver Person
	}

	MagicVersion struct {
		Chat Chat
		ID   uint64
	}
)
