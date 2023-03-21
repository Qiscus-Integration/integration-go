package usecase

import (
	"context"
	"encoding/json"
	"integration-go/domain"
	"time"

	"github.com/rs/zerolog/log"
)

type room struct {
	roomRepo  domain.RoomRepository
	omniRepo  domain.OmnichannelRepository
	cacheRepo domain.CacheRepository
}

const roomsCachedKey = "rooms"

// NewRoom returns a new instance of the Room use case.
func NewRoom(roomRepo domain.RoomRepository, omniRepo domain.OmnichannelRepository, cacheRepo domain.CacheRepository) *room {
	return &room{
		roomRepo:  roomRepo,
		omniRepo:  omniRepo,
		cacheRepo: cacheRepo,
	}
}

func (r *room) FetchRoom(ctx context.Context) (rooms []domain.Room, err error) {
	cached, _ := r.cacheRepo.Get(ctx, roomsCachedKey)
	if err = json.Unmarshal([]byte(cached), &rooms); err == nil {
		return
	}

	rooms, err = r.roomRepo.Fetch(ctx)
	if err != nil {
		return
	}

	go func() {
		roomsByte, _ := json.Marshal(rooms)
		if err := r.cacheRepo.Set(ctx, roomsCachedKey, string(roomsByte), 10*time.Second); err != nil {
			log.Ctx(ctx).Error().Msgf("unable to set cache key: %s err: %s", roomsCachedKey, err.Error())
		}
	}()

	return
}

func (r *room) CreateRoom(ctx context.Context, room *domain.Room) (err error) {
	err = r.omniRepo.CreateRoomTag(ctx, room.MultichannelRoomID, room.MultichannelRoomID)
	if err != nil {
		return
	}

	err = r.roomRepo.Save(ctx, room)
	return
}

func (r *room) ExecuteResolvedRoom(ctx context.Context) (err error) {
	rooms, err := r.FetchRoom(ctx)
	if err != nil {
		return
	}

	now := time.Now()
	for _, room := range rooms {
		diffMinutes := int(now.Sub(room.CreatedAt).Minutes())
		if diffMinutes < 5 {
			return
		}

		if err := r.omniRepo.ResolvedRoom(ctx, room.MultichannelRoomID); err != nil {
			log.Ctx(ctx).Error().Msgf("failed to resolved room: %s", err.Error())
			continue
		}

		err := r.roomRepo.DeleteBy(ctx, map[string]interface{}{
			"multichannel_room_id": room.MultichannelRoomID,
		})

		if err != nil {
			log.Ctx(ctx).Error().Msgf("failed to delete room: %s", err.Error())
			continue
		}

		if err := r.cacheRepo.Del(ctx, roomsCachedKey); err != nil {
			log.Ctx(ctx).Error().Msgf("failed to clear cache key: %s err: %s", roomsCachedKey, err.Error())
			continue
		}

	}

	return
}
