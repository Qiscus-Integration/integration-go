package cron

import (
	"context"
	"integration-go/internal/resolver"
	"integration-go/internal/room"
	"integration-go/internal/pkg/client"
	"integration-go/internal/pkg/config"
	"integration-go/internal/pkg/postgres"
	"integration-go/internal/pkg/qismo"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewServer() *Server {
	cfg := config.Load()

	db := postgres.NewGORM(cfg.Database)

	client := client.New()
	qismo := qismo.New(client, cfg.Qiscus.Omnichannel.URL, cfg.Qiscus.AppID, cfg.Qiscus.SecretKey)

	roomRepo := room.NewRepository(db)
	resolverSvc := resolver.NewService(roomRepo, qismo)

	return &Server{
		svc: resolverSvc,
	}
}

type Server struct {
	svc *resolver.Service
}

// Run starts the cron job and schedules it to execute every minute.
func (c *Server) Run() {
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
