package service

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/dto/adapter"
	"mojito-coding-test/app/data/model"
	"mojito-coding-test/app/data/types"
	"mojito-coding-test/app/errs"
	"mojito-coding-test/common/config"
	cdto "mojito-coding-test/common/data/dto"
)

type AuctionLotBid struct {
	// Facilities
	Config    *config.Config      `inject:""`
	Db        *gorm.DB            `inject:""`
	Validator *validator.Validate `inject:""`

	// Services
	AuctionLot *AuctionLot `inject:""`
}

func (s *AuctionLotBid) FindCurrentBid(db *gorm.DB, auctionLotId uint) (model.AuctionLotBid, error) {
	var item model.AuctionLotBid

	if err := db.Where("auction_lot_id=?", auctionLotId).
		Order("bid DESC").
		First(&item).Error; err != nil {
		return model.AuctionLotBid{}, err
	}

	return item, nil
}

func (s *AuctionLotBid) FindMaxBid(db *gorm.DB, auctionLotId, userId uint) (model.AuctionLotBidMax, error) {
	var item model.AuctionLotBidMax

	if err := db.Where("auction_lot_id=? AND user_id=? AND active=true", auctionLotId, userId).
		First(&item).Error; err != nil {
		return model.AuctionLotBidMax{}, err
	}

	return item, nil
}

func (s *AuctionLotBid) ListBids(db *gorm.DB, auctionLotId uint) ([]dto.AuctionLotBid, error) {
	var entities []model.AuctionLotBid

	// Select all AuctionLotBids
	// TODO: Limit to range.
	if err := db.Where("auction_lot_id=?", auctionLotId).
		Order("id ASC").
		Find(&entities).Error; err != nil {
		return nil, errors.Wrap(err, "error listing auction lot bids")
	}

	return adapter.AuctionLotBidModelsToDTOs(entities), nil
}

func (s *AuctionLotBid) ListMaxBids(db *gorm.DB, auctionLotId uint) ([]dto.AuctionLotBid, error) {
	var entities []model.AuctionLotBidMax

	// Select all AuctionLotBidMaxs
	// TODO: Limit to range.
	if err := db.Where("auction_lot_id=? AND active=true", auctionLotId).
		Order("id ASC").
		Find(&entities).Error; err != nil {
		return nil, errors.Wrap(err, "error listing auction lot max bids")
	}

	return adapter.AuctionLotBidMaxModelsToDTOs(entities), nil
}

func (s *AuctionLotBid) List(db *gorm.DB, auctionLotId uint) (dto.AuctionLotBidList, error) {
	bids, err := s.ListBids(db, auctionLotId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AuctionLotBidList{}, err
		}
	}

	maxBids, err := s.ListMaxBids(db, auctionLotId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AuctionLotBidList{}, err
		}
	}

	return dto.AuctionLotBidList{
		Bids:    bids,
		MaxBids: maxBids,
	}, nil
}

func (s *AuctionLotBid) Create(db *gorm.DB, auth cdto.Auth, auctionLotId uint,
	entityDTO dto.CreateAuctionLotBid) (dto.CreateAuctionLotBidResult, error) {

	bidIncrement := s.Config.GetUint("bid.increment")

	// Validate
	if err := s.Validator.Struct(entityDTO); err != nil {
		return dto.CreateAuctionLotBidResult{},
			errors.Wrap(err, "validation failed for creation of new auction lot bid")
	}

	if entityDTO.Bid == 0 || entityDTO.Bid%bidIncrement != 0 {
		return dto.CreateAuctionLotBidResult{}, errs.ErrAuctionLotBidInvalid.
			WithDetailsF("bid must be increment of %d", bidIncrement)
	}
	if entityDTO.MaxBid > 0 {
		if entityDTO.MaxBid%bidIncrement != 0 {
			return dto.CreateAuctionLotBidResult{}, errs.ErrAuctionLotBidInvalid.
				WithDetailsF("max bid must be increment of %d", bidIncrement)
		}
	}

	// Verify Auction Lot
	_, err := s.AuctionLot.GetByIdAndVerify(db, auctionLotId)
	if err != nil {
		return dto.CreateAuctionLotBidResult{}, err
	}

	// Verify Bid
	curBid, err := s.FindCurrentBid(db, auctionLotId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No bids, proceed.
			return s.doCreate(db, auth.UserId, auctionLotId, entityDTO.Bid, entityDTO.MaxBid,
				types.BidTypeUser, true, true)
		}

		return dto.CreateAuctionLotBidResult{}, errors.Wrap(err, "error encountered fetching current bid")
	}

	maxBid, err := s.FindMaxBid(db, auctionLotId, curBid.UserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.CreateAuctionLotBidResult{}, errors.Wrap(err, "error encountered fetching max bid")
		}
	}

	if entityDTO.Bid > curBid.Bid {
		if entityDTO.Bid > maxBid.MaxBid {
			// Greater than both the current bid and the max bid.
			return s.doCreate(db, auth.UserId, auctionLotId, entityDTO.Bid, entityDTO.MaxBid,
				types.BidTypeUser, true, true)
		} else {
			// Greater than current bid but not the max bid.
			return s.doCreate(db, curBid.UserId, auctionLotId, entityDTO.Bid, maxBid.MaxBid,
				types.BidTypeMaxBid, false, false)
		}
	} else if entityDTO.MaxBid > curBid.Bid {
		if entityDTO.MaxBid > maxBid.MaxBid {
			if entityDTO.MaxBid >= maxBid.MaxBid+bidIncrement {
				// Max bid is greater than the current bid and the max bid.
				return s.doCreate(db, auth.UserId, auctionLotId, maxBid.MaxBid+bidIncrement, entityDTO.MaxBid,
					types.BidTypeMaxBid, true, true)
			} else {
				// I would normally add a lot more detail to this error message.
				return dto.CreateAuctionLotBidResult{}, fmt.Errorf("expected max bid to be sufficient")
			}
		} else {
			// Not greater than the current bidders max bid.
			return s.doCreate(db, curBid.UserId, auctionLotId, entityDTO.MaxBid, maxBid.MaxBid,
				types.BidTypeMaxBid, false, false)
		}
	}

	return adapter.CreateAuctionLotBidResultModelToDTO(false, curBid, maxBid), nil
}

func (s *AuctionLotBid) doCreate(db *gorm.DB, userId, auctionLotId, bid, maxBid uint,
	bidType types.BidType, setMax, success bool) (dto.CreateAuctionLotBidResult, error) {

	var bidEntity model.AuctionLotBid
	var maxBidEntity = model.AuctionLotBidMax{
		MaxBid: maxBid,
	}

	/// Insert Bid
	bidEntity = model.AuctionLotBid{
		AuctionLotId: auctionLotId,
		UserId:       userId,
		Type:         bidType,
		Bid:          bid,
	}

	if err := db.Create(&bidEntity).Error; err != nil {
		return dto.CreateAuctionLotBidResult{}, err
	}

	/// Insert Max Bid
	if setMax {
		// Reset Active
		if err := resetMaxBid(db, auctionLotId, userId); err != nil {
			return dto.CreateAuctionLotBidResult{}, err
		}

		if maxBid > 0 {
			// Insert
			maxBidEntity = model.AuctionLotBidMax{
				AuctionLotId: auctionLotId,
				UserId:       userId,
				MaxBid:       maxBid,
				Active:       true,
			}

			if err := db.Create(&maxBidEntity).Error; err != nil {
				return dto.CreateAuctionLotBidResult{}, err
			}
		}
	}

	return adapter.CreateAuctionLotBidResultModelToDTO(success, bidEntity, maxBidEntity), nil
}

// Util

func resetMaxBid(db *gorm.DB, auctionLotId, userId uint) error {
	return db.Exec("UPDATE auction_lot_bid_max SET active=NULL "+
		"WHERE auction_lot_id=? AND user_id=?", auctionLotId, userId).Error
}
