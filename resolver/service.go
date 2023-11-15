//go:generate mockery --all --case snake --output ./mocks --exported
package resolver

import (
	"context"
	"integration-go/entity"
	"time"

	"github.com/rs/zerolog/log"
)

type omnichannel interface {
	ResolvedRoom(ctx context.Context, roomID string) error
}

type Service struct {
	roomRepo roomRepository
	omni     omnichannel
}

func NewService(roomRepo roomRepository, omni omnichannel) *Service {
	return &Service{
		roomRepo: roomRepo,
		omni:     omni,
	}
}

func (s *Service) ResolvedOmnichannelRoom(ctx context.Context) error {
	logCtx := log.Ctx(ctx).With().Str("func", "resolver.service.ResolveOmnichannelRoom").Logger()

	rooms, err := s.roomRepo.Fetch(ctx)
	if err != nil {
		logCtx.Error().Msgf("unable to fetch room data: %s", err.Error())
		return entity.ErrDatabase
	}

	now := time.Now()
	for _, room := range rooms {
		diffMinutes := int(now.Sub(room.CreatedAt).Minutes())
		if diffMinutes < 10 {
			return nil
		}

		if err := s.omni.ResolvedRoom(ctx, room.MultichannelRoomID); err != nil {
			logCtx.Error().Msgf("failed to resolved room: %s", err.Error())
			continue
		}

		err := s.roomRepo.DeleteBy(ctx, map[string]interface{}{
			"multichannel_room_id": room.MultichannelRoomID,
		})

		if err != nil {
			logCtx.Error().Msgf("failed to delete room: %s", err.Error())
			continue
		}

	}

	return nil
}
