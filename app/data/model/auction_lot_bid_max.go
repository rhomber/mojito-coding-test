package model

import (
	"gorm.io/gorm"
)

type AuctionLotBidMax struct {
	gorm.Model

	AuctionLotId uint
	UserId       uint
	MaxBid       uint
	Active       bool
}
