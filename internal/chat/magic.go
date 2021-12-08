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
	OpenMagicStatus MagicStatus = iota + 1
	ClosedMagicStatus
)
