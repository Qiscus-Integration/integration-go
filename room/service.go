//go:generate mockery --all --case snake --output ./mocks --exported
package room

import (
	"context"
	"errors"
	"integration-go/entity"
	"integration-go/qismo"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type omnichannel interface {
	CreateRoomTag(ctx context.Context, roomID string, tag string) error
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

func (s *Service) GetRoomByID(ctx context.Context, id int64) (*entity.Room, error) {
	logCtx := log.Ctx(ctx).With().Str("func", "room.service.GetRoomByID").Logger()

	room, err := s.roomRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logCtx.Error().Msgf("unable to find room: %s", err.Error())
			return nil, entity.ErrNotFound
		}

		logCtx.Error().Msgf("unable to find room: %s", err.Error())
		return nil, entity.ErrDatabase
	}

	return room, nil
}

func (s *Service) CreateRoom(ctx context.Context, req *qismo.WebhookNewSessionRequest) error {
	logCtx := log.Ctx(ctx).With().
		Str("func", "room.service.CreateRoom").
		Str("room_id", req.Payload.Room.IDStr).
		Logger()

	err := s.omni.CreateRoomTag(ctx, req.Payload.Room.IDStr, req.Payload.Room.IDStr)
	if err != nil {
		logCtx.Error().Msgf("unable to create omnichannel tag: %s", err.Error())
		return entity.ErrCantProceed
	}

	err = s.roomRepo.Save(ctx, &entity.Room{
		MultichannelRoomID: req.Payload.Room.IDStr,
	})

	if err != nil {
		logCtx.Error().Msgf("unable to save room data: %s", err.Error())
		return entity.ErrDatabase
	}

	return nil
}
