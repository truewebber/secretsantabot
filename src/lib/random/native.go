package random

import (
	"crypto/rand"
	"encoding/binary"
)

type (
	SourceNative struct{}
)

func NewSourceNative() *SourceNative {
	out := &SourceNative{}
	out.Seed(0)

	return out
}

func (n *SourceNative) Int63() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic("error read random bytes")
	}

	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

func (n *SourceNative) Seed(_ int64) {}
