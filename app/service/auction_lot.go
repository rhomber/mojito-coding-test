package service

import (
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/dto/adapter"
	"mojito-coding-test/app/data/model"
	"mojito-coding-test/app/errs"
	"mojito-coding-test/common/config"
	"time"
)

type AuctionLot struct {
	// Facilities
	Config    *config.Config      `inject:""`
	Db        *gorm.DB            `inject:""`
	Validator *validator.Validate `inject:""`

	// Services
}

func (s *AuctionLot) GetById(db *gorm.DB, id uint) (model.AuctionLot, error) {
	var item model.AuctionLot

	if err := db.First(&item, id).Error; err != nil {
		return model.AuctionLot{}, err
	}

	return item, nil
}

func (s *AuctionLot) GetByIdAndVerify(db *gorm.DB, id uint) (model.AuctionLot, error) {
	auction, err := s.GetById(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.AuctionLot{}, errs.ErrAuctionLotNotFound
		}

		return model.AuctionLot{}, err
	}

	now := time.Now()

	if auction.StartTime.After(now) {
		return model.AuctionLot{}, errs.ErrAuctionLotInvalid.
			WithDetailsF("start time is in the future")
	}
	if auction.EndTime.Before(now) {
		return model.AuctionLot{}, errs.ErrAuctionLotInvalid.
			WithDetailsF("end time is in the past")
	}

	return auction, nil
}

func (s *AuctionLot) List(db *gorm.DB) ([]dto.AuctionLot, error) {
	var entities []model.AuctionLot

	// Select all AuctionLots
	// TODO: Limit to range.
	if err := db.Order("id ASC").Find(&entities).Error; err != nil {
		return nil, errors.Wrap(err, "error listing auction lots")
	}

	return adapter.AuctionLotModelsToDTOs(entities), nil
}

func (s *AuctionLot) Create(db *gorm.DB, entityDTO dto.CreateAuctionLot) (dto.AuctionLot, error) {
	// Validate
	if err := s.Validator.Struct(entityDTO); err != nil {
		return dto.AuctionLot{}, errors.Wrap(err, "validation failed for creation of new auction lot")
	}

	// Populate
	entity := adapter.CreateAuctionLotDTOToModel(entityDTO)

	// Insert
	if err := db.Create(&entity).Error; err != nil {
		return dto.AuctionLot{}, errors.Wrap(err, "error creating new auction lot")
	}

	return adapter.AuctionLotModelToDTO(entity), nil
}
