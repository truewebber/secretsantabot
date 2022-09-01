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

	MagicVersion struct {
		Chat    *Chat
		ID      uint64
		Version uint8
		Status  MagicStatus
	}
)

type MagicStatus uint8

const (
	OpenMagicStatus MagicStatus = iota + 1
	ClosedMagicStatus
)
