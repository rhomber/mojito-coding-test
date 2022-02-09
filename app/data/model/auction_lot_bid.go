package model

import (
	"gorm.io/gorm"
	"mojito-coding-test/app/data/types"
)

type AuctionLotBid struct {
	gorm.Model

	AuctionLotId uint
	UserId       uint
	Type         types.BidType
	Bid          uint
}
