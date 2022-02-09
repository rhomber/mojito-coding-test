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

type User struct {
	// Facilities
	Config    *config.Config      `inject:""`
	Db        *gorm.DB            `inject:""`
	Validator *validator.Validate `inject:""`

	// Services
}

func (s *User) GetByEmail(db *gorm.DB, username string) (model.User, error) {
	var user model.User

	if err := db.First(&user, "email = ?", username).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *User) Create(db *gorm.DB, entityDTO dto.CreateUser) (dto.User, error) {
	// Validate
	if err := s.Validator.Struct(entityDTO); err != nil {
		return dto.User{}, errors.Wrap(err, "validation failed for creation of new user")
	}

	// Populate
	entity := adapter.CreateUserDTOToModel(entityDTO)

	// Insert
	if err := db.Create(&entity).Error; err != nil {
		return dto.User{}, errors.Wrap(err, "error creating new user")
	}

	return adapter.UserModelToDTO(entity), nil
}

func (s *User) List(db *gorm.DB) ([]dto.User, error) {
	var entities []model.User

	// Select all users
	// TODO: Limit to range.
	if err := db.Order("id ASC").Find(&entities).Error; err != nil {
		return nil, errors.Wrap(err, "error listing users")
	}

	return adapter.UserModelsToDTOs(entities), nil
}
