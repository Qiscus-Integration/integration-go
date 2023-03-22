package cron

import (
	"context"
	"integration-go/common"
	"integration-go/domain"
	"integration-go/repository/api"
	"integration-go/repository/cache"
	"integration-go/repository/pgsql"
	"integration-go/usecase"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jasonlvhit/gocron"
	"github.com/rs/zerolog/log"
)

// NewCron creates a new instance of Cron struct.
func NewCron() *Cron {
	dbConn := common.NewDatabase()
	cacheConn := common.NewCache(os.Getenv("REDIS_URL"))

	roomRepo := pgsql.NewPgsqlRoom(dbConn)
	omniRepo := api.NewApiQismo(os.Getenv("QISCUS_APP_ID"), os.Getenv("QISCUS_SECRET_KEY"))
	cacheRepo := cache.NewCacheRedis(cacheConn)

	roomUC := usecase.NewRoom(roomRepo, omniRepo, cacheRepo)

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

	gocron.Every(uint64(1 * time.Minute)).Do(func() {
		reqID := uuid.New().String()
		ctx := log.With().Str("request_id", reqID).Logger().WithContext(context.Background())

		c.roomUC.ExecuteResolvedRoom(ctx)
	})

	<-gocron.Start()
}
