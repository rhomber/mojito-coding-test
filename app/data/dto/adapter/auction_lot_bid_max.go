package adapter

import (
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/model"
)

func CreateAuctionLotBidMaxDTOToModel(dtoEntity dto.CreateAuctionLotBid) model.AuctionLotBidMax {
	return model.AuctionLotBidMax{
		MaxBid: dtoEntity.Bid,
	}
}

func AuctionLotBidMaxModelToDTO(entity model.AuctionLotBidMax) dto.AuctionLotBid {
	res := dto.AuctionLotBid{
		Id:           entity.ID,
		AuctionLotId: entity.AuctionLotId,
		UserId:       entity.UserId,
		Bid:          entity.MaxBid,
		CreatedAt:    entity.CreatedAt,
	}

	return res
}

func AuctionLotBidMaxModelsToDTOs(entities []model.AuctionLotBidMax) []dto.AuctionLotBid {
	dtos := make([]dto.AuctionLotBid, 0)

	for _, e := range entities {
		dtos = append(dtos, AuctionLotBidMaxModelToDTO(e))
	}

	return dtos
}
