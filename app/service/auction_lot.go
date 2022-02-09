package service

import (
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/dto/adapter"
	"mojito-coding-test/app/data/model"
	"mojito-coding-test/common/config"
)

type AuctionLot struct {
	// Facilities
	Config    *config.Config      `inject:""`
	Db        *gorm.DB            `inject:""`
	Validator *validator.Validate `inject:""`

	// Services
}

func (s *AuctionLot) Create(txn *gorm.DB, entityDTO dto.CreateAuctionLot) (dto.AuctionLot, error) {
	// Validate
	if err := s.Validator.Struct(entityDTO); err != nil {
		return dto.AuctionLot{}, errors.Wrap(err, "validation failed for creation of new auction lot")
	}

	// Populate
	entity := adapter.CreateAuctionLotDTOToModel(entityDTO)

	// Insert
	if err := txn.Create(&entity).Error; err != nil {
		return dto.AuctionLot{}, errors.Wrap(err, "error creating new auction lot")
	}

	return adapter.AuctionLotModelToDTO(entity), nil
}

func (s *AuctionLot) List(txn *gorm.DB) ([]dto.AuctionLot, error) {
	var entities []model.AuctionLot

	// Select all AuctionLots
	// TODO: Limit to range.
	if err := txn.Order("id ASC").Find(&entities).Error; err != nil {
		return nil, errors.Wrap(err, "error listing auction lots")
	}

	return adapter.AuctionLotModelsToDTOs(entities), nil
}
