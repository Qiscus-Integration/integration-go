package cron

import (
	"context"
	"integration-go/domain"
	"integration-go/repository/api"
	"integration-go/repository/cache"
	"integration-go/repository/persist"
	"integration-go/usecase"
	"integration-go/util"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// NewCron creates a new instance of Cron struct.
func NewCron() *Cron {
	dbConn := util.NewDatabase()
	cacheConn := util.NewCache(os.Getenv("REDIS_URL"))

	roomRepo := persist.NewPgsqlRoom(dbConn)
	roomCacheRepo := cache.NewRedisRoom(cacheConn, 10*time.Minute)
	omniRepo := api.NewApiQismo(os.Getenv("QISCUS_APP_ID"), os.Getenv("QISCUS_SECRET_KEY"))

	roomUC := usecase.NewRoom(roomRepo, omniRepo, roomCacheRepo)

	cron := &Cron{
		roomUC: roomUC,
	}

	return cron
}

type Cron struct {
	roomUC domain.RoomUsecase
}

// Run starts the cron job and schedules it to execute every minute.
func (c *Cron) Run() {
	log.Info().Msg("cron is started")

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(func() {
		reqID := uuid.New().String()
		ctx := log.With().Str("request_id", reqID).Logger().WithContext(context.Background())

		c.roomUC.ExecuteResolvedRoom(ctx)
	})

	s.StartBlocking()
}
