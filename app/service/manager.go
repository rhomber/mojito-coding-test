package service

import "mojito-coding-test/common/config"

type Manager struct {
	// Facilities
	Config *config.Config `inject:""`

	// Services
}
