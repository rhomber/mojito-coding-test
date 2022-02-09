package dto

type CreateUser struct {
	FirstName string `json:"first_name" yaml:"first_name" validate:"required"`
	LastName  string `json:"last_name" yaml:"last_name" validate:"required"`
	Email     string `json:"email" yaml:"email" validate:"required,email"`
}

type User struct {
	Id        uint   `json:"id" yaml:"id"`
	FirstName string `json:"first_name" yaml:"first_name" validate:"required"`
	LastName  string `json:"last_name" yaml:"last_name" validate:"required"`
	Email     string `json:"email" yaml:"email" validate:"required,email"`
}
