package random

import (
	"math/rand"
)

type (
	Rand struct {
		source rand.Source
	}

	Option func(*Rand)
)

func WithNative() Option {
	return func(r *Rand) {
		r.source = NewSourceNative()
	}
}

func WithRandomORG() Option {
	return func(r *Rand) {
		r.source = NewSourceRandomORG()
	}
}

func New(options ...Option) *rand.Rand {
	factory := new(Rand)

	switch len(options) {
	case 1:
		{
		}
	case 0:
		{
			options = []Option{
				WithNative(),
			}
		}
	default:
		{
			panic("Only one source can be set")
		}
	}

	options[0](factory)

	return rand.New(factory.source)
}
