package chat

type (
	Magic struct {
		Data []GiverReceiverPair
	}

	GiverReceiverPair struct {
		Giver    *Person
		Receiver *Person
	}

	MagicVersion struct {
		Chat *Chat
		ID   uint64
	}
)
