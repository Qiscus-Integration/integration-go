package usecase

import (
	"context"
	"errors"
	"integration-go/domain"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type room struct {
	roomRepo      domain.RoomRepository
	omniRepo      domain.OmnichannelRepository
	roomCacheRepo domain.RoomCacheRepository
}

// NewRoom returns a new instance of the Room use case.
func NewRoom(roomRepo domain.RoomRepository, omniRepo domain.OmnichannelRepository, roomCacheRepo domain.RoomCacheRepository) *room {
	return &room{
		roomRepo:      roomRepo,
		omniRepo:      omniRepo,
		roomCacheRepo: roomCacheRepo,
	}
}

func (r *room) GetRoomByID(ctx context.Context, id int64) (room domain.Room, err error) {
	room, err = r.roomCacheRepo.GetByID(ctx, id)
	if err == nil {
		return
	}

	room, err = r.roomRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = domain.ErrNotFound
			return
		}

		return
	}

	go func() {
		if err := r.roomCacheRepo.Save(ctx, room); err != nil {
			log.Ctx(ctx).Error().Msgf("unable to save room cache: %s", err.Error())
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
	log.Ctx(ctx).Info().Msg("ExecuteResolvedRoom job executed!")

	rooms, err := r.roomRepo.Fetch(ctx)
	if err != nil {
		return
	}

	now := time.Now()
	for _, room := range rooms {
		diffMinutes := int(now.Sub(room.CreatedAt).Minutes())
		if diffMinutes < 4 {
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

		if err := r.roomCacheRepo.DeletetByID(ctx, room.ID); err != nil {
			log.Ctx(ctx).Error().Msgf("failed to delete room cache: %s", err.Error())
			continue
		}

	}

	return
}
