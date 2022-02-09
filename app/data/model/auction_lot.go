package model

import (
	"gorm.io/gorm"
	"time"
)

type AuctionLot struct {
	gorm.Model

	Name      string
	StartTime time.Time
	EndTime   time.Time
}
