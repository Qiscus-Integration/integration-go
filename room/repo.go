package room

import (
	"context"
	"integration-go/entity"
)

type roomRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.Room, error)
	Save(ctx context.Context, room *entity.Room) error
}
