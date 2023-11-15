package cron

import (
	"context"
	"integration-go/config"
	"integration-go/pgsql"
	"integration-go/qismo"
	"integration-go/resolver"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// NewCron creates a new instance of Cron struct.
func NewCron() *cron {
	cfg := config.Load()
	db := pgsql.NewDatabase(cfg)

	roomRepo := pgsql.NewRoom(db)
	qismo := qismo.NewClient(cfg.Qiscus.Omnichannel.URL, cfg.Qiscus.AppID, cfg.Qiscus.SecretKey)

	resolverSvc := resolver.NewService(roomRepo, qismo)

	return &cron{
		svc: resolverSvc,
	}
}

type cron struct {
	svc *resolver.Service
}

// Run starts the cron job and schedules it to execute every minute.
func (c *cron) Run() {
	log.Info().Msg("cron is started")

	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Second().Do(func() {
		reqID := uuid.New().String()
		ctx := log.With().Str("request_id", reqID).Logger().WithContext(context.Background())

		err := c.svc.ResolvedOmnichannelRoom(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("error handle resolved room: %s", err.Error())
		}
	})

	s.StartBlocking()
}
