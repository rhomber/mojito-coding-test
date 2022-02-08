package config

import (
	"github.com/spf13/viper"
	"strings"
)

const (
	EnvPrefix = "GO_APP"
)

func New(configFile string) (*Config, error) {
	c := viper.New()
	c.SetConfigType("yaml")
	c.SetConfigName(configFile)
	c.AddConfigPath("/etc")
	c.AddConfigPath(".")

	// Defaults
	c.SetDefault("http.addr", ":8080")
	c.SetDefault("http.timeout", "60s")
	c.SetDefault("http.drain.interval", "1s")
	c.SetDefault("panic.debug", true)

	err := c.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c.SetEnvPrefix(EnvPrefix)
	replacer := strings.NewReplacer(".", "_")
	c.SetEnvKeyReplacer(replacer)
	c.AutomaticEnv()

	return &Config{
		Viper: c,
	}, nil
}

type Config struct {
	*viper.Viper
}

func (c *Config) IsVerbose() bool {
	return c.GetBool("verbose")
}
