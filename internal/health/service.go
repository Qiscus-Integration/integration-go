package health

import (
	"context"
	"integration-go/internal/entity"

	"github.com/rs/zerolog/log"
)

//go:generate mockery --with-expecter --case snake --name Repository
type Repository interface {
	CheckDatabase(ctx context.Context) error
	CheckRedis(ctx context.Context) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Check(ctx context.Context) (*entity.HealthComponent, bool) {
	healthComponent := &entity.HealthComponent{
		Database: entity.HealthStateOK,
		Redis:    entity.HealthStateOK,
	}

	if err := s.repo.CheckDatabase(ctx); err != nil {
		log.Ctx(ctx).Error().Msgf("check database error: %s", err.Error())
		healthComponent.Database = entity.HealthStateFail
	}

	if err := s.repo.CheckRedis(ctx); err != nil {
		log.Ctx(ctx).Error().Msgf("check redis error: %s", err.Error())
		healthComponent.Redis = entity.HealthStateFail
	}

	isHealthy := healthComponent.Database == entity.HealthStateOK &&
		healthComponent.Redis == entity.HealthStateOK

	return healthComponent, isHealthy
}
