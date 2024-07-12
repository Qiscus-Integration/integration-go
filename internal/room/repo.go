package room

import (
	"context"
	"integration-go/internal/entity"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Save(ctx context.Context, room *entity.Room) error {
	err := r.db.WithContext(ctx).Save(room).Error
	return err
}

func (r *repo) Fetch(ctx context.Context) ([]*entity.Room, error) {
	var rooms []*entity.Room
	err := r.db.WithContext(ctx).Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *repo) FindByID(ctx context.Context, id int64) (*entity.Room, error) {
	var room entity.Room
	err := r.db.WithContext(ctx).First(&room, id).Error
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *repo) DeleteBy(ctx context.Context, query map[string]interface{}) error {
	err := r.db.WithContext(ctx).Delete(&entity.Room{}, query).Error
	return err
}
