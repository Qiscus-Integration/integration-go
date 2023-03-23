package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"integration-go/domain"
	"time"

	rd "github.com/go-redis/redis"
)

type redisRoom struct {
	client *rd.Client
	exp    time.Duration
}

const keyParent = "room"

// NewRedisRoom returns a new instance of Redis implementation of RoomCacheRepository.
func NewRedisRoom(client *rd.Client, exp time.Duration) *redisRoom {
	return &redisRoom{client, exp}
}

func (r *redisRoom) Save(ctx context.Context, room domain.Room) (err error) {
	key := fmt.Sprintf("%s:%d", keyParent, room.ID)
	dataByte, err := json.Marshal(room)
	if err != nil {
		return
	}

	err = r.client.Set(key, string(dataByte), r.exp).Err()
	return
}

func (r *redisRoom) GetByID(ctx context.Context, id int64) (room domain.Room, err error) {
	key := fmt.Sprintf("%s:%d", keyParent, id)
	cached, err := r.client.Get(key).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(cached), &room)
	return
}

func (r *redisRoom) DeletetByID(ctx context.Context, id int64) (err error) {
	key := fmt.Sprintf("%s:%d", keyParent, id)
	err = r.client.Del(key).Err()
	return
}
