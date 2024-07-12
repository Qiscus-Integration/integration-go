package room

import (
	"context"
	"errors"
	"fmt"
	"integration-go/entity"
	"integration-go/qismo"

	"gorm.io/gorm"
)

//go:generate mockery --case snake --name Omnichannel
type Omnichannel interface {
	CreateRoomTag(ctx context.Context, roomID string, tag string) error
}

//go:generate mockery --case snake --name Repository
type Repository interface {
	FindByID(ctx context.Context, id int64) (*entity.Room, error)
	Save(ctx context.Context, room *entity.Room) error
}

type Service struct {
	repo Repository
	omni Omnichannel
}

func NewService(repo Repository, omni Omnichannel) *Service {
	return &Service{
		repo: repo,
		omni: omni,
	}
}

func (s *Service) GetRoomByID(ctx context.Context, id int64) (*entity.Room, error) {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &roomError{roomErrorNotFound}
		}

		return nil, fmt.Errorf("failed to find room: %w", err)
	}

	return room, nil
}

func (s *Service) CreateRoom(ctx context.Context, req *qismo.WebhookNewSessionRequest) error {
	err := s.omni.CreateRoomTag(ctx, req.Payload.Room.IDStr, req.Payload.Room.IDStr)
	if err != nil {
		return fmt.Errorf("failed to create omnichannel tag: %w", err)
	}

	err = s.repo.Save(ctx, &entity.Room{
		MultichannelRoomID: req.Payload.Room.IDStr,
	})

	if err != nil {
		return fmt.Errorf("failed to save room: %w", err)
	}

	return nil
}
