package types

import "github.com/truewebber/secretsantabot/domain/chat"

type Person struct {
	TelegramUserID int64
}

func PersonToDomain(p Person) chat.Person {
	return chat.Person{
		TelegramUserID: p.TelegramUserID,
	}
}

func DomainToPerson(p chat.Person) Person {
	return Person{
		TelegramUserID: p.TelegramUserID,
	}
}

func DomainsToPersons(persons []chat.Person) []Person {
	appPersons := make([]Person, 0, len(persons))

	for i := range persons {
		appPerson := DomainToPerson(persons[i])

		appPersons = append(appPersons, appPerson)
	}

	return appPersons
}
