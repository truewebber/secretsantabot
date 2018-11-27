package random

import (
	"encoding/binary"
	"math/rand"
)

type (
	cryptoRandSource struct {
		bytes [][]byte
	}
)

const (
	GetDefaultBytesNum = 128
)

func New() *rand.Rand {
	return rand.New(newCryptoRandSource())
}

func newCryptoRandSource() *cryptoRandSource {
	out := &cryptoRandSource{}
	out.Seed(0)

	return out
}

func (c *cryptoRandSource) Int63() int64 {
	if len(c.bytes) == 0 {
		c.Seed(0)
	}

	var b []byte
	b, c.bytes = c.bytes[0], c.bytes[1:]

	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b) & (1<<63 - 1))
}

func (c *cryptoRandSource) Seed(_ int64) {
	bytes, err := GetRandomBytes(GetDefaultBytesNum)
	if err != nil {
		panic(err.Error())
	}

	c.bytes = bytes
}
