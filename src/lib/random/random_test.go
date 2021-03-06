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

func TestNewDefault(t *testing.T) {
	r := New()

	for i := 0; i < 30; i++ {
		t.Log(r.Intn(1000))
	}
}

func TestNewNative(t *testing.T) {
	r := New(WithNative())

	for i := 0; i < 30; i++ {
		t.Log(r.Intn(1000))
	}
}

func TestNewRandomORG(t *testing.T) {
	r := New(WithRandomORG())

	for i := 0; i < 30; i++ {
		t.Log(r.Intn(1000))
	}
}
