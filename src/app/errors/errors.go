package errors

import "errors"

var (
	ErrChatTypeIsUnsupported = errors.New("chat type is unsupported")
	ErrChatIsPrivate         = errors.New("chat is private")
	ErrForbidden             = errors.New("forbidden")
	ErrNotEnoughParticipants = errors.New("not enough participants")
	ErrNotFound              = errors.New("not found")
	ErrAlreadyExists         = errors.New("already exists")
)
