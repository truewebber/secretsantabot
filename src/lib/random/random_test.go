package random

import "testing"

func TestIntInit(t *testing.T) {
	intRndPool, err := IntInit(10, 0, 9)
	if err != nil {
		t.Error(err.Error())

		return
	}

	for i := 0; i < 10; i++ {
		t.Log(intRndPool.Get())
	}
}
