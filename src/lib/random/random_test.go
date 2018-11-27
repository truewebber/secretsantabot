package random

import (
	"fmt"
	"testing"
)

func TestGetRandomBytes(t *testing.T) {
	rBytes, err := GetRandomBytes(16)
	if err != nil {
		t.Error(err.Error())

		return
	}

	for _, chunkBytes := range rBytes {
		for _, b := range chunkBytes {
			fmt.Printf("%x ", b)
		}
		fmt.Printf("\n")
	}
}

func TestNew(t *testing.T) {
	r := New()

	for i := 0; i < 3; i++ {
		t.Log(r.Intn(10))
	}
}
