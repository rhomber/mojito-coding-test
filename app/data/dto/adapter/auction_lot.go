package adapter

import (
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/model"
)

func CreateAuctionLotDTOToModel(dtoEntity dto.CreateAuctionLot) model.AuctionLot {
	return model.AuctionLot{
		Name:      dtoEntity.Name,
		StartTime: dtoEntity.StartTime,
		EndTime:   dtoEntity.EndTime,
	}
}

func AuctionLotModelToDTO(entity model.AuctionLot) dto.AuctionLot {
	res := dto.AuctionLot{
		Id:        entity.ID,
		Name:      entity.Name,
		StartTime: entity.StartTime,
		EndTime:   entity.EndTime,
	}

	return res
}

func AuctionLotModelsToDTOs(entities []model.AuctionLot) []dto.AuctionLot {
	dtos := make([]dto.AuctionLot, 0)

	for _, e := range entities {
		dtos = append(dtos, AuctionLotModelToDTO(e))
	}

	return dtos
}
