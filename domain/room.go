package domain

import (
	"context"
	"time"
)

// Room ...
type Room struct {
	ID                 int64
	MultichannelRoomID string `gorm:"index"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// RoomRepository defines an interface for saving Room data to persistent storage.
type RoomRepository interface {
	Save(ctx context.Context, room *Room) (err error)
	Fetch(ctx context.Context) (rooms []Room, err error)
	DeleteBy(ctx context.Context, query map[string]interface{}) (err error)
}

// OmnichannelRepository defines an interface for interacting with an omnichannel platform.
type OmnichannelRepository interface {
	CreateRoomTag(ctx context.Context, roomID string, tag string) (err error)
	ResolvedRoom(ctx context.Context, roomID string) (err error)
}

// RoomUsecase main application business logic hold room usecases
type RoomUsecase interface {
	CreateRoom(ctx context.Context, room *Room) (err error)
	FetchRoom(ctx context.Context) (rooms []Room, err error)
	ExecuteResolvedRoom(ctx context.Context) (err error)
}

// CacheRepository represents a repository that provides cache related operations.
type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) (err error)
	Get(ctx context.Context, key string) (res string, err error)
	Del(ctx context.Context, key string) (err error)
}
