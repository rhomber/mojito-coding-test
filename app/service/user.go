package service

import (
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"mojito-coding-test/app/data/dto"
	"mojito-coding-test/app/data/dto/adapter"
	"mojito-coding-test/common/config"
)

type User struct {
	// Facilities
	Config    *config.Config      `inject:""`
	Db        *gorm.DB            `inject:""`
	Validator *validator.Validate `inject:""`

	// Services
}

func (s *User) Create(txn *gorm.DB, entityDTO dto.CreateUser) (dto.User, error) {
	// Validate
	if err := s.Validator.Struct(entityDTO); err != nil {
		return dto.User{}, errors.Wrap(err, "validating data for new user")
	}

	// Populate
	entity := adapter.CreateUserDTOToModel(entityDTO)

	// Insert
	if err := txn.Create(&entity).Error; err != nil {
		return dto.User{}, errors.Wrap(err, "inserting new user")
	}

	return adapter.UserModelToDTO(entity), nil
}
