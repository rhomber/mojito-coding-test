package service

import (
	"gorm.io/gorm"
	"mojito-coding-test/common/config"
)

type Manager struct {
	// Facilities
	Config *config.Config `inject:""`
	Db     *gorm.DB       `inject:""`

	// Services
	User       *User       `inject:""`
	AuctionLot *AuctionLot `inject:""`
}
