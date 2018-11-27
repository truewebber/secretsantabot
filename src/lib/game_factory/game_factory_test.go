package game_factory

import (
	"testing"

	"lib/random"
)

func TestGameFactory_Magic(t *testing.T) {
	list := []string{
		"Артем",
		"Aня",
		"Леша",
		"Коля",
		"Настя",
		"Дима",
		"Стася",
		"Илья",
	}

	l := len(list)

	shuffled := make([]string, l)
	copy(shuffled, list)

	for _, man := range list {
		var pair string
		pair, shuffled = getPairFor(man, shuffled)

		t.Log(man, pair)
	}
}

func getPairFor(man string, list []string) (string, []string) {
	var elem string
	for len(elem) == 0 {
		if len(list) != 1 {
			list = shuffle(list)
		}

		if man == list[0] {
			continue
		}

		if man == "Артем" {
			if list[0] == "Аня" {
				continue
			}
		}

		elem, list = list[0], list[1:]
	}

	return elem, list
}

func shuffle(list []string) []string {
	l := len(list)
	entropy := 1

	rnd := random.New()

	for k := 0; k < entropy; k++ {
		for i := 0; i < l; i++ {
			rndI := rnd.Intn(l - 1)

			t := list[i]
			list[i] = list[rndI]
			list[rndI] = t
		}
	}

	return list
}
