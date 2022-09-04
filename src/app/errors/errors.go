package errors

import "errors"

var (
	ErrChatTypeIsUnsupported = errors.New("chat type is unsupported")
	ErrChatIsPrivate         = errors.New("chat is private")
	ErrForbidden             = errors.New("forbidden")
)
