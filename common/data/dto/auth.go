package dto

type Auth struct {
	UserId   uint   `json:"user_id" yaml:"user_id"`
	Username string `json:"username" yaml:"username"`
}
