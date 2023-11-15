package entity

import (
	"time"
)

// Room ...
type Room struct {
	ID                 int64     `json:"id"`
	MultichannelRoomID string    `json:"multichannel_room_id" gorm:"index"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
