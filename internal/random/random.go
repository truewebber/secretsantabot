package random

import (
	"math/rand"
)

func NewRandomORG() *rand.Rand {
	return rand.New(newSourceRandomORG())
}

func NewNative() *rand.Rand {
	return rand.New(newSourceNative())
}
