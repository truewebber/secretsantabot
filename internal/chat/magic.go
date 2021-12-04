package chat

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
	OpenMagicStatus = iota + 1
	ClosedMagicStatus
)
