package config

import (
	"mojito-coding-test/common/data/env"
	"strings"
)

func (c *Config) GetEnv() env.Env {
	return env.Env(strings.ToUpper(c.GetString("env")))
}

func (c *Config) IsEnvLocal() bool {
	return c.GetEnv() == env.Local
}

func (c *Config) IsEnvTest() bool {
	return c.GetEnv() == env.Test
}

func (c *Config) IsEnvProd() bool {
	return c.GetEnv() == env.Prod
}
