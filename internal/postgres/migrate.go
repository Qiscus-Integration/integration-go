package postgres

import (
	"integration-go/internal/entity"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&entity.Room{})
	if err != nil {
		log.Fatal().Msgf("failed to run migration: %s", err.Error())
	}
}
