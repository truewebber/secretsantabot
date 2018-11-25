package storage

import (
	"testing"
	"time"

	"lib/model"
)

var (
	strg *RedisStorage
	me   = &model.HellMan{
		TelegramId: 1111,
		Username:   "truewebber",
		FirstName:  "Aleksey",
		LastName:   "Kish",
		EnrollAt:   time.Now(),
	}
)

func init() {
	strg, _ = NewRedisStorage()
}

func TestRedisStorage_Enroll(t *testing.T) {
	err := strg.Enroll(me)
	if err != nil {
		t.Error(err.Error())

		return
	}

	t.Log("OK")
}

func TestRedisStorage_IsEnroll(t *testing.T) {
	result, err := strg.IsEnroll(me)
	if err != nil {
		t.Error(err.Error())

		return
	}

	t.Log(result)
}

func TestRedisStorage_DropEnroll(t *testing.T) {
	err := strg.DropEnroll(me)
	if err != nil {
		t.Error(err.Error())

		return
	}

	t.Log("OK")
}

func TestRedisStorage_IsMagicDone(t *testing.T) {
	result, err := strg.IsMagicDone()
	if err != nil {
		t.Error(err.Error())

		return
	}

	t.Log(result)
}
