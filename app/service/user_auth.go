package service

import (
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"mojito-coding-test/common/config"
	"mojito-coding-test/common/data/dto"
	"mojito-coding-test/common/errs"
)

type UserAuth struct {
	// Facilities
	Config    *config.Config      `inject:""`
	Db        *gorm.DB            `inject:""`
	Validator *validator.Validate `inject:""`

	// Services
	User *User `inject:""`
}

// Very basic authentication (no password).
func (s *UserAuth) Authenticate(db *gorm.DB, username string) (dto.Auth, error) {
	user, err := s.User.GetByEmail(db, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.Auth{}, errs.ErrUserAuthFailed
		}

		return dto.Auth{}, errs.ErrUserAuthFailed.
			WithDetailsF("encountered error while authenticating - %+v", err)
	}

	return dto.Auth{
		UserId:   user.ID,
		Username: user.Email,
	}, nil
}
