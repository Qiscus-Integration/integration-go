package pgsql

import (
	"context"
	"integration-go/entity"

	"gorm.io/gorm"
)

type Room struct {
	db *gorm.DB
}

// NewRoom creates and returns a new instance of the `room` struct which implements the `RoomRepository` interface.
// It takes a `gorm.DB` object as an argument, which is used to connect to a PostgresSQL database.
func NewRoom(db *gorm.DB) *Room {
	return &Room{
		db: db,
	}
}

func (r *Room) Save(ctx context.Context, room *entity.Room) error {
	err := r.db.WithContext(ctx).Save(room).Error
	return err
}

func (r *Room) Fetch(ctx context.Context) ([]*entity.Room, error) {
	var rooms []*entity.Room
	err := r.db.WithContext(ctx).Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *Room) FindByID(ctx context.Context, id int64) (*entity.Room, error) {
	var room entity.Room
	err := r.db.WithContext(ctx).First(&room, id).Error
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *Room) DeleteBy(ctx context.Context, query map[string]interface{}) error {
	err := r.db.WithContext(ctx).Delete(&entity.Room{}, query).Error
	return err
}
