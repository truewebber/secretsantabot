package errors

import "errors"

var (
	ErrRegisterLocalChatIsRestricted = errors.New("register local chat is restricted")
	ErrNotFound                      = errors.New("not found")
)
