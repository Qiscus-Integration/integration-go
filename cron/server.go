package cron

import (
	"context"
	"fmt"
	"integration-go/client"
	"integration-go/config"
	"integration-go/qismo"
	"integration-go/resolver"
	"integration-go/room"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServer() *Server {
	cfg := config.Load()
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("unable to open db connection")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get sql db")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

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
