package dto

import (
	"mojito-coding-test/app/data/types"
	"time"
)

type CreateAuctionLotBid struct {
	Bid    uint `json:"bid" yaml:"bid" validate:"required"`
	MaxBid uint `json:"max_bid" yaml:"max_bid"`
}

type CreateAuctionLotBidResult struct {
	Success      bool          `json:"success" yaml:"success"`
	AuctionLotId uint          `json:"auction_lot_id,omitempty" yaml:"auction_lot_id,omitempty"`
	UserId       uint          `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Type         types.BidType `json:"type,omitempty" yaml:"type,omitempty"`
	Bid          uint          `json:"bid,omitempty" yaml:"bid,omitempty"`
	MaxBid       uint          `json:"max_bid,omitempty" yaml:"max_bid,omitempty"`
	CreatedAt    time.Time     `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type AuctionLotBid struct {
	Id           uint          `json:"id" yaml:"id"`
	AuctionLotId uint          `json:"auction_lot_id" yaml:"auction_lot_id"`
	UserId       uint          `json:"user_id" yaml:"user_id"`
	Type         types.BidType `json:"type,omitempty" yaml:"type,omitempty"`
	Bid          uint          `json:"bid" yaml:"bid"`
	CreatedAt    time.Time     `json:"created_at" yaml:"created_at"`
}

type AuctionLotBidList struct {
	Bids    []AuctionLotBid `json:"bids" yaml:"bids"`
	MaxBids []AuctionLotBid `json:"max_bids" yaml:"max_bids"`
}
