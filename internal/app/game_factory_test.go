package app

import (
	"testing"

	"github.com/truewebber/secretsantabot/internal/model"
	"github.com/truewebber/secretsantabot/internal/storage"
)

var (
	list = map[int]*model.HellMan{
		10: {FirstName: "Артем"},
		11: {FirstName: "Aня"},
		12: {FirstName: "Леша"},
		13: {FirstName: "Коля"},
		14: {FirstName: "Настя"},
		15: {FirstName: "Дима"},
		16: {FirstName: "Стася"},
		17: {FirstName: "Илья"},
	}
)

func TestGameFactory_MagicDeny(t *testing.T) {
	gFactory := New(nil)
	mapList := list
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
			pairId, shuffled, err = gFactory.GetPairFor(id, shuffled)
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
		t.Error(noResultErr.Error())

		return
	}

	for hellSanta, hellMan := range out {
		t.Log(hellSanta.FirstName, hellMan.FirstName)
	}
}

func TestGameFactory_Magic(t *testing.T) {
	strg, err := storage.NewRedisStorage()
	if err != nil {
		t.Error(err.Error())

		return
	}
	gFactory := New(strg)

	result, err := gFactory.Magic()
	if err != nil {
		t.Error(err.Error())

		return
	}

	for hellSanta, hellMan := range result {
		t.Log(hellSanta.FirstName, hellMan.FirstName)
	}
}

func TestShuffle(t *testing.T) {
	idList := make([]int, 0)
	for id := range list {
		idList = append(idList, id)
	}

	t.Log("REGULAR")
	for _, man := range idList {
		t.Log(man)
	}

	shuffle(idList)

	t.Log("SHUFFLED")
	for _, man := range idList {
		t.Log(man)
	}
}
