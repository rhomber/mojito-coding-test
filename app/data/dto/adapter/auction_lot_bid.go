package adapter

import (
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/model"
)

func AuctionLotBidModelToDTO(entity model.AuctionLotBid) dto.AuctionLotBid {
	res := dto.AuctionLotBid{
		Id:           entity.ID,
		AuctionLotId: entity.AuctionLotId,
		UserId:       entity.UserId,
		Type:         entity.Type,
		Bid:          entity.Bid,
		CreatedAt:    entity.CreatedAt,
	}

	return res
}

func AuctionLotBidModelsToDTOs(entities []model.AuctionLotBid) []dto.AuctionLotBid {
	dtos := make([]dto.AuctionLotBid, 0)

	for _, e := range entities {
		dtos = append(dtos, AuctionLotBidModelToDTO(e))
	}

	return dtos
}

func CreateAuctionLotBidResultModelToDTO(success bool, bidEntity model.AuctionLotBid,
	maxBidEntity model.AuctionLotBidMax) dto.CreateAuctionLotBidResult {

	res := dto.CreateAuctionLotBidResult{
		Success:      success,
		AuctionLotId: bidEntity.AuctionLotId,
		UserId:       bidEntity.UserId,
		Type:         bidEntity.Type,
		Bid:          bidEntity.Bid,
		MaxBid:       maxBidEntity.MaxBid,
		CreatedAt:    bidEntity.CreatedAt,
	}

	return res
}
