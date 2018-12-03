package storage

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"

	"lib/config"
	"lib/model"
)

type (
	RedisStorage struct {
		client *redis.Client
	}
)

const (
	GameKey  = "game"
	MagicKey = "magic"
)

func NewRedisStorage() (*RedisStorage, error) {
	rc := redis.NewClient(&redis.Options{
		Addr: config.Get().GetString("redis.addr"),
		DB:   config.Get().GetInt("redis.db"),
	})

	if cmd := rc.Ping(); cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return &RedisStorage{
		client: rc,
	}, nil
}

func (r *RedisStorage) SaveMagic(save map[*model.HellMan]*model.HellMan) error {
	for santa, man := range save {
		cmd := r.client.HSet(MagicKey, strconv.Itoa(santa.TelegramId), man.TelegramId)
		if cmd.Err() != nil {
			r.DropMagic()

			return cmd.Err()
		}
	}

	return nil
}

func (r *RedisStorage) DropMagic() error {
	cmd := r.client.Del(MagicKey)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *RedisStorage) ListMagic() (map[int]*model.HellMan, error) {
	cmd := r.client.HGetAll(MagicKey)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	list, err := r.ListEnrolled()
	if cmd.Err() != nil {
		return nil, err
	}

	out := make(map[int]*model.HellMan)
	for santaIdStr, manIdStr := range cmd.Val() {
		santaId, err := strconv.Atoi(santaIdStr)
		if err != nil {
			return nil, err
		}

		manId, err := strconv.Atoi(manIdStr)
		if err != nil {
			return nil, err
		}

		out[santaId] = list[manId]
	}

	return out, nil
}

func (r *RedisStorage) ListEnrolled() (map[int]*model.HellMan, error) {
	cmd := r.client.HGetAll(GameKey)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	out := make(map[int]*model.HellMan)
	for _, value := range cmd.Val() {
		obj := new(model.HellMan)

		err := json.Unmarshal([]byte(value), obj)
		if err != nil {
			return nil, cmd.Err()
		}

		out[obj.TelegramId] = obj
	}

	return out, nil
}

func (r *RedisStorage) Enroll(user *model.HellMan) error {
	data, _ := json.Marshal(user)
	key := strconv.Itoa(user.TelegramId)

	cmd := r.client.HSet(GameKey, key, data)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *RedisStorage) DropEnroll(user *model.HellMan) error {
	key := strconv.Itoa(user.TelegramId)

	cmd := r.client.HDel(GameKey, key)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return cmd.Err()
	}

	return nil
}

func (r *RedisStorage) IsEnroll(user *model.HellMan) (bool, error) {
	key := strconv.Itoa(user.TelegramId)

	cmd := r.client.HGet(GameKey, key)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return false, cmd.Err()
	} else if cmd.Err() == redis.Nil {
		return false, nil
	}

	return true, nil
}

func (r *RedisStorage) IsMagicDone() (bool, error) {
	cmd := r.client.Keys(MagicKey)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return false, cmd.Err()
	}

	if len(cmd.Val()) != 1 {
		return false, nil
	}

	return true, nil
}
