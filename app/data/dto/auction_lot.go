package dto

import "time"

type CreateAuctionLot struct {
	Name      string    `json:"name" yaml:"name" validate:"required"`
	StartTime time.Time `json:"start_time" yaml:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" yaml:"end_time" validate:"required"`
}

type AuctionLot struct {
	Id        uint      `json:"id" yaml:"id"`
	Name      string    `json:"name" yaml:"name" validate:"required"`
	StartTime time.Time `json:"start_time" yaml:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" yaml:"end_time" validate:"required"`
}
