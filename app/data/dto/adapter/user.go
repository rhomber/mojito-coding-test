package adapter

import (
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/model"
)

func CreateUserDTOToModel(dtoEntity dto.CreateUser) model.User {
	return model.User{
		FirstName: dtoEntity.FirstName,
		LastName:  dtoEntity.LastName,
		Email:     dtoEntity.Email,
	}
}

func UserModelToDTO(entity model.User) dto.User {
	res := dto.User{
		Id:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
	}

	return res
}

func UserModelsToDTOs(entities []model.User) []dto.User {
	dtos := make([]dto.User, 0)

	for _, e := range entities {
		dtos = append(dtos, UserModelToDTO(e))
	}

	return dtos
}
