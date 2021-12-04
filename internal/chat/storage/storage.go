package storage

import "github.com/truewebber/secretsantabot/internal/chat"

type Storage interface {
	SaveChat(*chat.Chat) error
	SavePerson(*chat.Person) error

	SaveMagic(chat.Magic) error
	DropMagic() error
	ListMagic() (chat.Magic, error)
	IsMagicDone() (bool, error)

	ListEnrolled() ([]chat.Person, error)
	Enroll(*chat.Person) error
	DropEnroll(*chat.Person) error
	IsEnroll(*chat.Person) (bool, error)
}
