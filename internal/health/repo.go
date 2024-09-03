package health

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type repo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRepository(db *gorm.DB, rdb *redis.Client) *repo {
	return &repo{
		db:  db,
		rdb: rdb,
	}
}

func (r *repo) CheckDatabase(ctx context.Context) error {
	sqlDB, err := r.db.WithContext(ctx).DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}

func (r *repo) CheckRedis(ctx context.Context) error {
	err := r.rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}

	return nil
}
