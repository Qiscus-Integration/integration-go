package resolver

import (
	"context"
	"integration-go/entity"
)

type roomRepository interface {
	Fetch(ctx context.Context) ([]*entity.Room, error)
	DeleteBy(ctx context.Context, query map[string]interface{}) error
}
