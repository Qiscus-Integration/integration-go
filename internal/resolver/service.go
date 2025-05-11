package resolver

import (
	"context"
	"fmt"
	"integration-go/internal/entity"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockery --with-expecter --case snake --name RoomRepository
type RoomRepository interface {
	Fetch(ctx context.Context) ([]*entity.Room, error)
	DeleteBy(ctx context.Context, query map[string]any) error
}

//go:generate mockery --with-expecter --case snake --name Omnichannel
type Omnichannel interface {
	ResolvedRoom(ctx context.Context, roomID string) error
}

type Service struct {
	roomRepo RoomRepository
	omni     Omnichannel
}

func NewService(roomRepo RoomRepository, omni Omnichannel) *Service {
	return &Service{
		roomRepo: roomRepo,
		omni:     omni,
	}
}

func (s *Service) ResolvedOmnichannelRoom(ctx context.Context) error {
	rooms, err := s.roomRepo.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch rooms: %w", err)
	}

	now := time.Now()
	for _, room := range rooms {
		diffMinutes := int(now.Sub(room.CreatedAt).Minutes())
		if diffMinutes < 10 {
			return nil
		}

		if err := s.omni.ResolvedRoom(ctx, room.MultichannelRoomID); err != nil {
			log.Ctx(ctx).Error().Msgf("failed to resolved room: %s", err.Error())
			continue
		}

		err := s.roomRepo.DeleteBy(ctx, map[string]any{
			"multichannel_room_id": room.MultichannelRoomID,
		})

		if err != nil {
			log.Ctx(ctx).Error().Msgf("failed to delete room: %s", err.Error())
			continue
		}

	}

	return nil
}
