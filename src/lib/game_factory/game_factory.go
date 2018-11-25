package game_factory

import (
	"github.com/pkg/errors"

	"lib/model"
	"lib/storage"
)

type (
	GameFactory struct {
		storage *storage.RedisStorage
	}
)

var (
	ErrorAlreadyEnroll       = errors.New("Already enroll")
	ErrorMagicWasAlreadyDone = errors.New("Magic was already done")
)

func New(strg *storage.RedisStorage) *GameFactory {
	return &GameFactory{
		storage: strg,
	}
}

func (g *GameFactory) ListEnrolled() ([]*model.HellMan, error) {
	return g.storage.ListEnrolled()
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
