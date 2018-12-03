package game_factory

import (
	"github.com/pkg/errors"

	"lib/config"
	"lib/model"
	"lib/random"
	"lib/storage"
)

type (
	GameFactory struct {
		storage *storage.RedisStorage
		deny    map[int]int
	}
)

var (
	ErrorAlreadyEnroll       = errors.New("Already enroll")
	ErrorMagicWasAlreadyDone = errors.New("Magic was already done")
	ErrorMagicIsNotProceed   = errors.New("Magic is not proceed for now")
	ErrorYouAreNotInThisGame = errors.New("You are not in this game")
	ErrorNotEnough           = errors.New("Not enough users enrolled to start game")
)

const (
	DefaultEntropy = 3
	Retries        = 10000
)

func New(storageInstance *storage.RedisStorage) *GameFactory {
	out := &GameFactory{
		storage: storageInstance,
		deny:    make(map[int]int),
	}

	deny := make(map[int]int)
	err := config.Get().UnmarshalKey("rules.deny", &deny)
	if err == nil {
		out.deny = deny
	}

	return out
}

func (g *GameFactory) GetMyMagic(santa *model.HellMan) (*model.HellMan, error) {
	isMagicDone, err := g.storage.IsMagicDone()
	if err != nil {
		return nil, err
	}

	if !isMagicDone {
		return nil, ErrorMagicIsNotProceed
	}

	santaManMap, err := g.storage.ListMagic()
	if err != nil {
		return nil, err
	}

	if man, ok := santaManMap[santa.TelegramId]; ok {
		return man, nil
	}

	return nil, ErrorYouAreNotInThisGame
}

func (g *GameFactory) Magic() (map[*model.HellMan]*model.HellMan, error) {
	isMagicDone, err := g.storage.IsMagicDone()
	if err != nil {
		return nil, err
	}

	if isMagicDone {
		return nil, ErrorMagicWasAlreadyDone
	}

	mapList, err := g.storage.ListEnrolled()
	if err != nil {
		return nil, err
	}

	if len(mapList) < 2 {
		return nil, ErrorNotEnough
	}

	var out map[*model.HellMan]*model.HellMan

	var noResultErr error
	for i := 0; i <= Retries; i++ {
		out = make(map[*model.HellMan]*model.HellMan)

		idList := make([]int, 0)
		for id := range mapList {
			idList = append(idList, id)
		}

		shuffled := make([]int, len(idList))
		copy(shuffled, idList)

		var fuckError error
		for _, id := range idList {
			var (
				pairId int
				err    error
			)
			pairId, shuffled, err = g.GetPairFor(id, shuffled)
			if err != nil {
				fuckError = err

				break
			}

			out[mapList[id]] = mapList[pairId]
		}

		if fuckError == nil {
			noResultErr = nil

			break
		} else {
			noResultErr = fuckError
		}
	}
	if noResultErr != nil {
		return nil, noResultErr
	}

	err = g.storage.SaveMagic(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (g *GameFactory) ListEnrolled() ([]*model.HellMan, error) {
	mapList, err := g.storage.ListEnrolled()
	if err != nil {
		return nil, err
	}

	out := make([]*model.HellMan, 0)
	for _, hellMan := range mapList {
		out = append(out, hellMan)
	}

	return out, nil
}

func (g *GameFactory) Enroll(user *model.HellMan) error {
	isMagicDone, err := g.storage.IsMagicDone()
	if err != nil {
		return err
	}

	if isMagicDone {
		return ErrorMagicWasAlreadyDone
	}

	isEnroll, err := g.storage.IsEnroll(user)
	if err != nil {
		return err
	}

	if isEnroll {
		return ErrorAlreadyEnroll
	}

	return g.storage.Enroll(user)
}

func (g *GameFactory) DropEnroll(user *model.HellMan) error {
	isMagicDone, err := g.storage.IsMagicDone()
	if err != nil {
		return err
	}

	if isMagicDone {
		return ErrorMagicWasAlreadyDone
	}

	isEnroll, err := g.storage.IsEnroll(user)
	if err != nil {
		return err
	}

	if !isEnroll {
		return ErrorAlreadyEnroll
	}

	return g.storage.DropEnroll(user)
}

// ---------------------------------------------------------------------------------------------------------------------

func (g *GameFactory) GetPairFor(man int, list []int) (int, []int, error) {
	noResults := 0

	var elem int
	for elem == 0 {
		if noResults > 10 {
			return 0, nil, errors.New("FUCK IT")
		}

		if len(list) != 1 {
			shuffle(list)
		}

		if man == list[0] {
			noResults++

			continue
		}

		if denyMan, ok := g.deny[man]; ok {
			if denyMan == list[0] {
				noResults++

				continue
			}
		}

		elem, list = list[0], list[1:]
	}

	return elem, list, nil
}

func shuffle(list []int) {
	l := len(list)
	rnd := random.New()

	for e := 0; e < DefaultEntropy; e++ {
		for i := 0; i < l; i++ {
			rndI := rnd.Intn(l - 1)

			t := list[i]
			list[i] = list[rndI]
			list[rndI] = t
		}
	}
}
