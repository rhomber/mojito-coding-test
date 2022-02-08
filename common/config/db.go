package config

import (
	"fmt"
)

func (c *Config) GetDsn() string {
	if dbFile := c.GetDbFile(); dbFile != "" {
		// Hard coded options (only for test, usually I would make them configurable).
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.GetDbFile())
	}

	return ""
}

func (c *Config) GetDbFile() string {
	return c.GetString("db.file")
}
